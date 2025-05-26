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

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*entity.User)
	if ok {
		c.authService.Logout(user.UserID)
	}

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
