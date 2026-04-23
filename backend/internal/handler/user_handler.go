package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
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
	logger      *slog.Logger
}

func NewUserHandler(userService *service.UserService, logger *slog.Logger) *UserHandler {
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
		h.logger.Error("decoding register request ", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", err.Error())
		return
	}

	existingUser, err := h.userService.FindByEmail(r.Context(), req.Email)
	if err != nil {
		h.logger.Error("findByEmail", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	if existingUser != nil {
		middleware.WriteError(w, http.StatusConflict, "conflict", "Email Id already present")
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
		h.logger.Error("setUserPassword", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	err = h.userService.Create(r.Context(), user)
	if err != nil {
		h.logger.Error("create", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	token, err := middleware.GenerateToken(user.ID, 24*time.Hour, middleware.ScopeAuth)
	if err != nil {
		h.logger.Error("generateToken", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	err = h.userService.CreateToken(r.Context(), token)
	if err != nil {
		h.logger.Error("createToken failed", "error", err, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	h.logger.Info("user registered", "user_id", user.ID, "email", user.Email, "name", user.Name)

	middleware.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token.PlainText,
		"user":  user,
	})

	// middleware.WriteJSON(w, http.StatusCreated, map[string]any{"user": user})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req model.LoginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("decoding login request ", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	user, err := h.userService.FindByEmail(r.Context(), req.Email)
	if err != nil {
		h.logger.Error("findByEmail", "error", err)
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
		h.logger.Error("generateToken", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	err = h.userService.CreateToken(r.Context(), token)
	if err != nil {
		h.logger.Error("createToken failed", "error", err, "user_id", user.ID)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	h.logger.Info("user logged in", "user_id", user.ID, "email", user.Email, "name", user.Name)

	middleware.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token.PlainText,
		"user":  user,
	})
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("readIDParam", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
		return
	}

	user, err := h.userService.FindByID(r.Context(), userID)
	if err != nil {
		h.logger.Error("findByID", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Error("readIDParam", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == model.AnonymousUser {
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "you must be logged in to update")
		return
	}

	existingUser, err := h.userService.FindByID(r.Context(), userID)
	if err != nil {
		h.logger.Error("findByID", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if existingUser == nil {
		http.NotFound(w, r)
		return
	}

	if existingUser.ID != userID {
		middleware.WriteError(w, http.StatusForbidden, "Forbidden", "You cannot update this user profile")
		return
	}

	// At this point we can assume we are able to find existing user
	var input model.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("decodingUpdateRequest", "error", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid request payload")
		return
	}

	if input.Name != nil {
		existingUser.Name = *input.Name
	}
	if input.AvatarURL != nil {
		existingUser.AvatarURL = *input.AvatarURL
	}
	if input.Password != nil && *input.Password != "" {
		err = SetUserPassword(&existingUser.PasswordHash, *input.Password)
		if err != nil {
			h.logger.Error("setUserPassword", "error", err)
			middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
			return
		}
	}
	existingUser.UpdatedAt = time.Now()

	err = h.userService.Update(r.Context(), existingUser)
	if err != nil {
		h.logger.Error("update", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"user": existingUser})
}

func (h *UserHandler) AllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())

	if err != nil {
		h.logger.Error("getAllUsers", "error", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "could not retrieve users")
		return
	}

	if len(users) == 0 {
		middleware.WriteJSON(w, http.StatusOK, map[string]any{"users": []model.User{}})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"users": users})
}
