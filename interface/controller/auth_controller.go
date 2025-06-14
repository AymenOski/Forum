package controller

import (
	"html/template"
	"net/http"

	"forum/domain/entity"
	"forum/usecase"
)

type AuthController struct {
	authService *usecase.AuthService
	postService *usecase.PostService
	templates   *template.Template
}

func NewAuthController(authService *usecase.AuthService, postService *usecase.PostService, templates *template.Template) *AuthController {
	return &AuthController{
		authService: authService,
		postService: postService,
		templates:   templates,
	}
}

func (c *AuthController) HandleSignup(w http.ResponseWriter, r *http.Request) {
	// If the method is GET that means loading the html page
	if r.Method == http.MethodGet {
		c.renderTemplate(w, "register.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}

	// if the method is POST that means the user is creating an account
	name := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := c.authService.Signup(name, email, password)
	if err != nil {
		c.renderTemplate(w, "register.html", map[string]interface{}{
			"registerError": err.Error(),
			"username":      name,
			"email":         email, // roll-back values when re-rendering so that the user doesn't have to re-enter it
		})
		return
	}

	_ = user // User created successfully

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.renderTemplate(w, "login.html", nil)
		return
	}
	if r.Method != http.MethodPost {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	token, user, err := c.authService.Login(email, password)
	if err != nil {
		c.renderTemplate(w, "login.html", map[string]interface{}{
			"loginError": err.Error(),
			"email":      email,
		})
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
	})

	_ = user

	// Redirect to home page
	http.Redirect(w, r, "/layout", http.StatusSeeOther)
}

func (c *AuthController) HandleMainPage(w http.ResponseWriter, r *http.Request) {
	c.ShowMainPage(w, r)
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*entity.User)
	if ok {
		c.authService.Logout(user.ID)
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

// New method to validate session (can be used by middleware)
func (c *AuthController) ValidateSessionToken(token string) (*entity.UserSession, error) {
	return c.authService.ValidateSession(token)
}
