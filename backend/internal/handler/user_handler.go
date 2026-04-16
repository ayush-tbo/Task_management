package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/floqast/task-management/backend/internal/middleware"
	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
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

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
		return
	}

	userIDStr := strconv.Itoa(int(userID))

	user, err := h.userService.FindByID(r.Context(), userIDStr)
	if err != nil {
		h.logger.Printf("ERROR: findByID: %v", err)
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (h *UserHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		middleware.WriteError(w, http.StatusBadRequest, "bad request", "invalid user id")
		return
	}
	userIDStr := strconv.Itoa(int(userID))

	existingUser, err := h.userService.FindByID(r.Context(), userIDStr)
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
