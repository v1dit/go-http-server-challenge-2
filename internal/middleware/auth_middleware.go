package middleware

import (
	"challenge2/internal/models"
	"challenge2/internal/service"
	"context"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"

// AuthMiddleware validates the session token and attaches the user to context.
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Session-Token")
			if token == "" {
				http.Error(w, "missing session token", http.StatusUnauthorized)
				return
			}

			user, err := authService.GetProfileByToken(token)
			if err != nil {
				switch err {
				case service.ErrInvalidToken:
					http.Error(w, "invalid session token", http.StatusUnauthorized)
				default:
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext returns the authenticated user if middleware ran successfully.
func UserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(userContextKey).(models.User)
	return user, ok
}
