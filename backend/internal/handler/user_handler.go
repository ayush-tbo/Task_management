package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService *service.UserService
	logger      *log.Logger
}

func NewUserHandler(userService *service.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func SetUserPassword(p *model.Password, plaintextpassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextpassword), 12)
	if err != nil {
		return err
	}

	p.PlainText = &plaintextpassword
	p.Hash = hash
	return nil
}

func MatchUserPassword(p *model.Password, plaintextpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextpassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err // internal server error
		}
	}

	return true, nil
}

func (h *UserHandler) validateRegisterRequest(req *model.RegisterUserRequest) error {
	if req.Name == "" {
		return errors.New("username is required")
	}
	if len(req.Name) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request : %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", err.Error())
		return
	}

	user := &model.User{
		ID:        primitive.NewObjectID().Hex(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	err = SetUserPassword(&user.PasswordHash, req.Password)
	if err != nil {
		h.logger.Printf("ERROR: setUserPassword %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	err = h.userService.Create(r.Context(), user)
	if err != nil {
		h.logger.Printf("ERROR: create %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"user": user})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req model.LoginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding login request : %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	user, err := h.userService.FindByEmail(r.Context(), req.Email)
	if err != nil {
		h.logger.Printf("ERROR: findByEmail: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if user == nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid email or password")
		return
	}

	match, err := MatchUserPassword(&user.PasswordHash, req.Password)
	if err != nil || !match {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid email or password")
		return
	}

	token, err := middleware.GenerateToken(user.ID, 24*time.Hour, middleware.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: generateToken: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	err = h.userService.CreateToken(r.Context(), token)
	if err != nil {
		h.logger.Printf("ERROR: createToken: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token.PlainText,
		"user":  user,
	})
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
		return
	}

	user, err := h.userService.FindByID(r.Context(), userID)
	if err != nil {
		h.logger.Printf("ERROR: findByID: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
		return
	}

	existingUser, err := h.userService.FindByID(r.Context(), userID)
	if err != nil {
		h.logger.Printf("ERROR: findByID: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if existingUser == nil {
		http.NotFound(w, r)
		return
	}

	// At this point we can assume we are able to find existing user
	var input model.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	if input.Name != nil {
		existingUser.Name = *input.Name
	}
	if input.AvatarURL != nil {
		existingUser.AvatarURL = *input.AvatarURL
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to update")
		return
	}

	err = h.userService.Update(r.Context(), existingUser)
	if err != nil {
		h.logger.Printf("ERROR: update: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"user": existingUser})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}

func (h *UserHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	middleware.WriteError(w, http.StatusNotImplemented, "not_implemented", "endpoint stub")
}
