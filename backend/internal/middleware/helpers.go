package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/go-chi/chi/v5"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err, msg string) {
	WriteJSON(w, status, model.ErrorResponse{Error: err, Message: msg})
}

func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func ReadIDParam(r *http.Request) (string, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return "", errors.New("invalid id parameter")
	}

	return idParam, nil
}
