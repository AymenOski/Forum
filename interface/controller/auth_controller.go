package controller

import (
	"fmt"
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
		// Showing the error page temporarily
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusUnauthorized,
			Error:      err.Error(),
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, user, err := c.authService.Login(email, password)
	if err != nil {
		// Showing the error page temporarily
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusUnauthorized,
			Error:      err.Error(),
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}



func (c *AuthController) HandleGlobal(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path == "/" && r.Method == http.MethodGet){
		c.ShowMainPage(w, r)
	}else{
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusNotFound,
			Error:      "Page Not Found.",
		})
	}
}


func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
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


func (c *AuthController) StaticFileServer(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	switch r.URL.Path {

		case  "/static/css/*.css", "/static/images/background.jpg":
				http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
		
		case  "/static/", "/static/css/", "/static/images/":
			c.ShowErrorPage(w, ErrorMessage{
				StatusCode: http.StatusForbidden,
				Error:      "StatusForbidden",
				})	
		
		default :
			c.ShowErrorPage(w, ErrorMessage{
					StatusCode: http.StatusForbidden,
					Error:      "Page Not Found.",
				})
	}
}