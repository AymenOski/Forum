package middleware

import (
	"context"
	"net/http"

	"forum/domain/entity"
	"forum/domain/repository"
	"forum/usecase"
)

type AuthMiddleware struct {
	authService *usecase.AuthService
	userRepo    repository.UserRepository
}

func NewAuthMiddleware(authService *usecase.AuthService, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		userRepo:    userRepo,
	}
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
		session, err := m.authService.ValidateSession(cookie.Value)
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

		// Get the actual user from the session
		user, err := m.userRepo.GetByID(session.UserID)
		if err != nil {
			// User not found, clear cookie and redirect
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

		// Add both user and session to request context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "session", session)
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

		// Invalid session, clear the cookie and proceed
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		// Invalid session, proceed to login/register
		next.ServeHTTP(w, r)
	})
}

// Optional: Middleware that adds user to context but doesn't require authentication
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No session cookie, proceed without user
			next.ServeHTTP(w, r)
			return
		}

		// Validate session
		session, err := m.authService.ValidateSession(cookie.Value)
		if err != nil {
			// Invalid session, clear cookie and proceed without user
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			next.ServeHTTP(w, r)
			return
		}

		// Get the actual user from the session
		user, err := m.userRepo.GetByID(session.UserID)
		if err != nil {
			// User not found, clear cookie and proceed without user
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			next.ServeHTTP(w, r)
			return
		}

		// Add both user and session to request context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper method to get user from context
func GetUserFromContext(r *http.Request) (*entity.User, bool) {
	user, ok := r.Context().Value("user").(*entity.User)
	return user, ok
}

// Helper method to get session from context
func GetSessionFromContext(r *http.Request) (*entity.UserSession, bool) {
	session, ok := r.Context().Value("session").(*entity.UserSession)
	return session, ok
}
