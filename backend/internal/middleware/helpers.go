package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

func ReadURLParam(r *http.Request, name string) (string, error) {
	param := chi.URLParam(r, name)
	if param == "" {
		return "", errors.New("invalid " + name + " parameter")
	}
	return param, nil
}

func GetPaginationParams(r *http.Request) (int, int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
