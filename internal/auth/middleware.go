package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type UserContext struct {
	ID   uint
	Role string
}

type contextKey string

const UserContextKey = contextKey("user")

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			authHeader = r.URL.Query().Get("token")
		}

		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Malformed Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Missing Authorization Token", http.StatusUnauthorized)
			return
		}

		userID, role, err := ParseJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, UserContext{ID: userID, Role: role})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) (UserContext, error) {
	user, ok := ctx.Value(UserContextKey).(UserContext)
	if !ok {
		return UserContext{}, errors.New("user not found in context")
	}
	return user, nil
}
