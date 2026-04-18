package middleware

import (
	"context"
	"net/http"
	"strings"

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

		user, err := um.UserService.FindByID(r.Context(), token.UserID)
		if user == nil {
			WriteError(w, http.StatusUnauthorized, "Token Expired", "token expired or invalid")
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
