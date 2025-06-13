package middleware

import (
	"context"
	"net/http"

	"forum/usecase"
)

type AuthMiddleware struct {
	authService *usecase.AuthService
}

func NewAuthMiddleware(authService *usecase.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate session
		user, err := m.authService.ValidateSession(cookie.Value)
		if err != nil {
			// Clear invalid cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RedirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No session cookie, proceed
			next.ServeHTTP(w, r)
			return
		}

		// Check if session is valid
		_, err = m.authService.ValidateSession(cookie.Value)
		if err == nil {
			// Valid session exists, redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Invalid session, proceed to login/register
		next.ServeHTTP(w, r)
	})
}
