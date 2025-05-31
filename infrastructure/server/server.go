package server

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	infra_repository "forum/infrastructure/repository"
	"forum/interface/controller"
	"forum/interface/middleware"
	"forum/usecase"
)

type Server struct {
	authController *controller.AuthController
	authMiddleware *middleware.AuthMiddleware
	templates      *template.Template
	router         *http.ServeMux
	httpServer     *http.Server
}

func Forum_server(db *sql.DB) *Server {
	// Initialize repositories
	userRepo := infra_repository.NewSQLiteUserRepository(db)
	sessionRepo := infra_repository.NewSQLiteUserSessionRepository(db)

	// Initialize services
	authService := usecase.NewAuthService(userRepo, sessionRepo)

	// Load templates
	templates := loadTemplates()

	// Initialize controllers
	authController := controller.NewAuthController(authService, templates)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, userRepo)

	// Create server
	server := &Server{
		authController: authController,
		authMiddleware: authMiddleware,
		templates:      templates,
		router:         http.NewServeMux(),
	}

	// Setup routes
	server.setupRoutes()

	// Configure HTTP server
	server.httpServer = &http.Server{
		Addr:         ":8080",
		Handler:      server.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start background tasks
	server.startBackgroundTasks()

	return server
}

func (s *Server) setupRoutes() {
	// Static files
	fs := http.FileServer(http.Dir("./static/"))
	s.router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public routes (redirect if authenticated)
	//s.router.Handle("/login", s.authMiddleware.RedirectIfAuthenticated(http.HandlerFunc(s.authController.ShowLogin)))
	//s.router.Handle("/register", s.authMiddleware.RedirectIfAuthenticated(http.HandlerFunc(s.authController.ShowRegister)))

	// Auth routes
	s.router.HandleFunc("/login", s.authController.HandleLogin)
	s.router.HandleFunc("/register", s.authController.HandleRegister)
	s.router.HandleFunc("/logout", s.authController.HandleLogout)
	s.router.HandleFunc("/refresh-session", s.authController.HandleRefreshSession)

	// Protected routes
	s.router.Handle("/", s.authMiddleware.OptionalAuth(http.HandlerFunc(s.handleHome)))
	s.router.Handle("/profile", s.authMiddleware.RequireAuth(http.HandlerFunc(s.handleProfile)))
	s.router.Handle("/create-post", s.authMiddleware.RequireAuth(http.HandlerFunc(s.handleCreatePost)))

	// Admin routes (if needed)
	s.router.Handle("/admin", s.authMiddleware.RequireAuth(http.HandlerFunc(s.handleAdmin)))

	// API routes (if needed)
	s.router.HandleFunc("/api/validate-session", s.handleValidateSession)
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() error {
	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// SetPort allows changing the server port
func (s *Server) SetPort(port string) {
	s.httpServer.Addr = ":" + port
}

// Handler functions
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	user, isAuthenticated := middleware.GetUserFromContext(r)

	data := map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"User":            user,
		"Title":           "Forum Home",
	}

	err := s.templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"User":  user,
		"Title": "User Profile",
	}

	err := s.templates.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"User":  user,
			"Title": "Create Post",
		}

		err := s.templates.ExecuteTemplate(w, "create-post.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		// Handle post creation logic here
		title := r.FormValue("title")
		content := r.FormValue("content")

		// TODO: Implement post creation logic
		log.Printf("Creating post: %s by user %s", title, user.UserName)
		log.Printf("Post content: %s", content)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// TODO: Add admin role checking logic
	// if !user.IsAdmin {
	//     http.Error(w, "Forbidden", http.StatusForbidden)
	//     return
	// }

	data := map[string]interface{}{
		"User":  user,
		"Title": "Admin Panel",
	}

	err := s.templates.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleValidateSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session found", http.StatusUnauthorized)
		return
	}

	_, err = s.authController.ValidateSessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"valid": true}`))
}

// Background task to clean up expired sessions
func (s *Server) startBackgroundTasks() {
	// Clean up expired sessions every hour
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			err := s.authController.CleanupExpiredSessions()
			if err != nil {
				log.Printf("Error cleaning up expired sessions: %v", err)
			} else {
				log.Println("Cleaned up expired sessions")
			}
		}
	}()
}

// Load templates from templates directory
func loadTemplates() *template.Template {
	templateDir := "templates"

	// Create template with helper functions
	tmpl := template.New("").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"since": func(t time.Time) string {
			duration := time.Since(t)
			if duration < time.Minute {
				return "just now"
			} else if duration < time.Hour {
				return time.Since(t).Round(time.Minute).String() + " ago"
			} else if duration < 24*time.Hour {
				return time.Since(t).Round(time.Hour).String() + " ago"
			}
			return t.Format("2006-01-02")
		},
	})

	// Walk through templates directory
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				log.Printf("Error parsing template %s: %v", path, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}

	return tmpl
}