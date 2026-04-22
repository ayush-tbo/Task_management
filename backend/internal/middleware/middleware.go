package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/floqast/task-management/backend/internal/model"
	"github.com/floqast/task-management/backend/internal/service"
)

type UserMiddleware struct {
	UserService service.UserService
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(UserContextKey).(*model.User)
	if !ok {
		panic("missing user in request") // bad actor call
	}
	// fmt.Println(user)
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Within this anonymous function
		// We can interject any incoming request to our server

		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ") //Bearer <TOKEN>
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			WriteError(w, http.StatusUnauthorized, "Not authorized", "invalid authorization header")
			return
		}

		tokenIncoming := headerParts[1]

		token, err := um.UserService.GetToken(r.Context(), ScopeAuth, tokenIncoming)
		if err != nil {
			WriteError(w, http.StatusUnauthorized, "Token not found", "Invalid Token")
			return
		}

		if token == nil {
			WriteError(w, http.StatusUnauthorized, "unauthorized", "token expired")
			return
		}

		user, err := um.UserService.FindByID(r.Context(), token.UserID)
		if err != nil || user == nil {
			WriteError(w, http.StatusUnauthorized, "unauthorized", "user not found")
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if service.UserIsAnonymous(user) {
			WriteError(w, http.StatusUnauthorized, "invalid user", "you must be logged in to access this route")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLogger logs every incoming request with method, path, user, status, and duration.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			// try to get user from context (set by Authenticate middleware)
			userID := "anonymous"
			userName := "anonymous"
			if user, ok := r.Context().Value(UserContextKey).(*model.User); ok && user != nil && user != model.AnonymousUser {
				userID = user.ID
				userName = user.Name
			}

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration_ms", duration.Milliseconds(),
				"user_id", userID,
				"user_name", userName,
				"remote_addr", r.RemoteAddr,
			}

			if wrapped.statusCode >= 500 {
				logger.Error("request completed", attrs...)
			} else if wrapped.statusCode >= 400 {
				logger.Warn("request completed", attrs...)
			} else {
				logger.Info("request completed", attrs...)
			}
		})
	}
}
