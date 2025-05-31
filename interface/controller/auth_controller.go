package controller

import (
	"html/template"
	"net/http"
	"forum/domain/entity"
	"forum/usecase"
)

type AuthController struct {
	authService *usecase.AuthService
	templates   *template.Template
}

func NewAuthController(authService *usecase.AuthService, templates *template.Template) *AuthController {
	return &AuthController{
		authService: authService,
		templates:   templates,
	}
}

func (c *AuthController) GetAuthService() *usecase.AuthService {
	return c.authService
}

func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	err := c.templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, user, err := c.authService.Login(email, password)
	if err != nil {
		// Render login page with error
		data := map[string]interface{}{
			"Error": err.Error(),
		}
		c.templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	_ = user // You can use user data if needed

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	err := c.templates.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := c.authService.Register(name, email, password)
	if err != nil {
		// Render register page with error
		data := map[string]interface{}{
			"Error": err.Error(),
		}
		c.templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	_ = user // User created successfully

	// Redirect to login page with success message
	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err == nil && cookie.Value != "" {
		// Use the LogoutByToken method to invalidate the specific session
		c.authService.LogoutByToken(cookie.Value)
	}

	// Alternative: Get user from context and logout all sessions
	// user, ok := r.Context().Value("user").(*entity.User)
	// if ok {
	// 	c.authService.Logout(user.ID)
	// }

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// New method to handle session refresh
func (c *AuthController) HandleRefreshSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session found", http.StatusUnauthorized)
		return
	}

	newToken, err := c.authService.RefreshSession(cookie.Value)
	if err != nil {
		// Clear invalid cookie and redirect to login
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

	// Update cookie with new token (if different)
	if newToken != cookie.Value {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    newToken,
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
			Secure:   false,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session refreshed"))
}

// New method to validate session (can be used by middleware)
func (c *AuthController) ValidateSessionToken(token string) (*entity.UserSession, error) {
	return c.authService.ValidateSession(token)
}

// New method to cleanup expired sessions (can be called periodically)
func (c *AuthController) CleanupExpiredSessions() error {
	return c.authService.CleanupExpiredSessions()
}