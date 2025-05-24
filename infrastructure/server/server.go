package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"forum/infrastructure/infra_repository"
)

var tmpl *template.Template

type Err struct {
	Message string
	Value   int
}

func init() {
	var err error
	tmpl, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}
}

func Froum_server() *http.Server {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/signup", registerHandler)
	mux.HandleFunc("/layout", layoutHandler)
	mux.HandleFunc("/", loginHandler)

	serve := &http.Server{
		Addr:    ":8080",
		Handler: logMiddleware(notFoundMiddleware(mux)),
	}
	return serve
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func notFoundMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedPaths := map[string]bool{
			"/":                        true,
			"/signup":                  true, // Add this line
			"/layout":                  true,
			"/artist/":                 true,
			"/static/css/login.css":    true,
			"/static/css/register.css": true,
			"/static/css/layout.css":   true,
		}
		if !allowedPaths[r.URL.Path] {
			// notFoundHandler(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loginHandler(wr http.ResponseWriter, r *http.Request) {
	renderTemplate(wr, "", "login.html")
	fmt.Fprintf(wr, "this is test")
}

func layoutHandler(wr http.ResponseWriter, r *http.Request) {
	postRepo := infra_repository.NewSQLitePostRepository(&sql.DB{})
	posts, err := postRepo.GetAll()
	if err != nil {
		// log.Fatal("posts error")
		//	return
	}
	renderTemplate(wr, posts, "layout.html")
	// fmt.Fprintf(wr, "this is test")
}

func registerHandler(wr http.ResponseWriter, r *http.Request) {
	renderTemplate(wr, "", "register.html")
	fmt.Fprintf(wr, "this is test")
}

func renderTemplate(wr http.ResponseWriter, data interface{}, template string) {
	if isFileAvailable(template) {
		err := tmpl.ExecuteTemplate(wr, template, data)
		if err != nil {
			renderError(wr, http.StatusInternalServerError, "Template Rendering Error")
		}
	} else {
		renderError(wr, http.StatusInternalServerError, template+" not available")
	}
}

func renderError(wr http.ResponseWriter, statusCode int, msg string) {
	if isFileAvailable("errorPage.html") {
		wr.WriteHeader(statusCode)
		renderTemplate(wr, &Err{Message: msg, Value: statusCode}, "errorPage.html")
	} else {
		fallbackErrorMessage(wr)
	}
}

func isFileAvailable(file string) bool {
	_, err := os.Stat("./templates/" + file)
	return err == nil
}

func fallbackErrorMessage(wr http.ResponseWriter) {
	wr.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(wr, "Error 500: Website is under maintenance for security issues.")
}
