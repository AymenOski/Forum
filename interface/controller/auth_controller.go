package controller

import (
	"html/template"
	"net/http"
	"strings"

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
			"email":         email,
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

	http.Redirect(w, r, "/layout", http.StatusSeeOther)
}

func (c *AuthController) HandleGlobal(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodGet {
		c.ShowMainPage(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/static/") {
		switch r.URL.Path {
		case "/static/css/layout.css", "/static/css/login.css", "/static/css/posts.css", "/static/css/register.css", "/static/images/background.jpg":
			http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)

		case "/static/", "/static/css/", "/static/images/":
			c.ShowErrorPage(w, ErrorMessage{
				StatusCode: http.StatusForbidden,
				Error:      "StatusForbidden",
			})

		default:
			c.ShowErrorPage(w, ErrorMessage{
				StatusCode: http.StatusForbidden,
				Error:      "Page Not Found.",
			})
		}
	} else {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusNotFound,
			Error:      "Page Not Found.",
		})
	}
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get session token from cookie
	// cookie, err := r.Cookie("session_token")
	// if err == nil && cookie.Value != "" {
	// 	// Use the LogoutByToken method to invalidate the specific session
	// 	c.authService.Logout(cookie.Value)
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
