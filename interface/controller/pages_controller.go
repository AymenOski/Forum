package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

type ErrorMessage struct {
	StatusCode int
	Error      string
}

func (c *AuthController) renderTemplate(w http.ResponseWriter, TmplName string, data interface{}) {

	w.Header().Set("Content-type", "text/html")
	err := c.templates.ExecuteTemplate(w, TmplName, data)
	if err != nil {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error: fmt.Sprintf("Internal Server Error"),
		})
	}
}

func (c *AuthController) ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate(w, "register.html", nil)
}

func (c *AuthController) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate(w, "login.html", nil)
}

func (c *AuthController) ShowMainPage(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate(w, "layout.html", nil)
}

func (c *AuthController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {

	TmplStatus, _ := template.ParseFiles("templates/error.html")
	if TmplStatus == nil {	http.Error(w, fmt.Sprintf("%d - %s",data.StatusCode, data.Error), data.StatusCode); return }
	
	w.WriteHeader(data.StatusCode)
	TmplStatus.Execute(w, data)

}