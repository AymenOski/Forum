package controller

import (
	"fmt"
	"net/http"
)

type ErrorMessage struct {
	StatusCode int
	Error      string
}

func (c *AuthController) renderTemplate(w http.ResponseWriter, template string, data interface{}) {
	w.Header().Set("Content-type", "text/html")
	err := c.templates.ExecuteTemplate(w, template, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading %s: %v", template, err), http.StatusInternalServerError)
	}
}

func (c *AuthController) ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate(w, "register.html", nil)
}

func (c *AuthController) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate(w, "login.html", nil)
}

func (c *AuthController) ShowMainPage(w http.ResponseWriter, r *http.Request) {
	posts, err := c.postService.GetPosts()
	if err != nil {
		// Showing the error page temporarily
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	c.renderTemplate(w, "layout.html", posts)
}

func (c *AuthController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(data.StatusCode)

	err := c.templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%d - %s", data.StatusCode, data.Error), data.StatusCode)
	}
}
