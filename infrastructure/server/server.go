package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	infra_repository "forum/infrastructure/repository"
	"forum/interface/controller"
	"forum/usecase"
)

var (
	tmpl     *template.Template
	database *sql.DB
)

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

func Froum_server(db *sql.DB) *http.Server {
	database = db
	mux := http.NewServeMux()

	// entities
	postRepo := infra_repository.NewSQLitePostRepository(database)
	userRepo := infra_repository.NewSQLiteUserRepository(database)
	// usecase
	postService := usecase.NewPostService(postRepo, userRepo)
	authService := usecase.NewAuthService(userRepo)
	// controller where the handlers should live
	postController := controller.NewPostController(*postService)
	authController := controller.NewAuthController(authService, tmpl)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/signup", authController.HandleRegister)
	mux.HandleFunc("/layout", layoutHandler)
	mux.HandleFunc("/", loginHandler)

	serve := &http.Server{
		Addr:    ":8080",
		Handler: LogMiddleware(notFoundMiddleware(mux)),
	}
	return serve
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func notFoundMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedPaths := map[string]bool{
			"/":                        true,
			"/signup":                  true,
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
}

func layoutHandler(wr http.ResponseWriter, r *http.Request) {
	postRepo := infra_repository.NewSQLitePostRepository(Gdb)
	posts, err := postRepo.GetAll()
	if err != nil {
		// Handle this error properly
		http.Error(wr, "Failed to fetch posts", http.StatusInternalServerError)
		log.Printf("Error fetching posts: %v", err)
		return
	}
	renderTemplate(wr, posts, "layout.html")
}

func registerHandler(wr http.ResponseWriter, r *http.Request) {
	renderTemplate(wr, "", "register.html")
	fmt.Fprintf(wr, "this is test")
}

func renderTemplate(wr http.ResponseWriter, data any, template string) {
	if isFileAvailable(template) {
		err := tmpl.ExecuteTemplate(wr, template, data)
		if err != nil {
			RenderError(wr, http.StatusInternalServerError, "Template Rendering Error")
		}
	} else {
		RenderError(wr, http.StatusInternalServerError, template+" not available")
	}
}

func RenderError(wr http.ResponseWriter, statusCode int, msg string) {
	if isFileAvailable("errorPage.html") {
		wr.WriteHeader(statusCode)
		RenderTemplate(wr, &Err{Message: msg, Value: statusCode}, "errorPage.html")
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
