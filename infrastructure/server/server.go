package server

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"html/template"
// 	"log"
// 	"net/http"

// 	"forum/domain/entity"
// 	infra_repository "forum/infrastructure/repository"
// 	"forum/interface/controller"
// 	"forum/interface/middleware"
// 	"forum/usecase"
// )

// var (
// 	tmpl *template.Template
// 	Gdb  *sql.DB
// )

// type ErrorResponse struct {
// 	Success bool   `json:"success"`
// 	Error   string `json:"error"`
// 	Message string `json:"message,omitempty"`
// }

// type SuccessResponse struct {
// 	Success bool        `json:"success"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data,omitempty"`
// 	Token   string      `json:"token,omitempty"`
// 	User    interface{} `json:"user,omitempty"`
// }

// func init() {
// 	var err error
// 	tmpl, err = template.ParseGlob("./templates/*.html")
// 	if err != nil {
// 		log.Printf("Warning: Failed to initialize templates: %v", err)
// 	}
// }

// func Forum_server(db *sql.DB) *http.Server {
// 	Gdb = db
// 	mux := http.NewServeMux()

// 	// Initialize services and controllers
// 	sessionRepo := infra_repository.NewSQLiteUserSessionRepository(Gdb)
// 	userRepo := infra_repository.NewSQLiteUserRepository(Gdb)
// 	authService := usecase.NewAuthService(userRepo, sessionRepo)
// 	authController := controller.NewAuthController(authService, tmpl)
// 	authMiddleware := middleware.NewAuthMiddleware(authService)

// 	// Static files
// 	fileServer := http.FileServer(http.Dir("./static"))
// 	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

// 	// Authentication routes
// 	mux.HandleFunc("/auth/register", handleRegister(authController))
// 	mux.HandleFunc("/auth/login", handleLogin(authController))
// 	mux.HandleFunc("/auth/logout", requireAuth(authMiddleware, handleLogout(authController)))
// 	mux.HandleFunc("/auth/me", requireAuth(authMiddleware, handleMe(authService)))

// 	// Legacy routes for compatibility
// 	mux.HandleFunc("/register", handleRegister(authController))
// 	mux.HandleFunc("/login", handleLogin(authController))
// 	mux.HandleFunc("/logout", requireAuth(authMiddleware, handleLogout(authController)))

// 	// Root endpoint
// 	mux.HandleFunc("/", handleRoot)

// 	// Health check
// 	mux.HandleFunc("/health", handleHealth)

// 	server := &http.Server{
// 		Addr:    ":8080",
// 		Handler: LogMiddleware(corsMiddleware(mux)),
// 	}

// 	log.Println("Auth server starting on :8080")
// 	log.Println("Available endpoints:")
// 	log.Println("  POST /auth/register - Register new user")
// 	log.Println("  POST /auth/login    - Login user")
// 	log.Println("  POST /auth/logout   - Logout user (requires auth)")
// 	log.Println("  GET  /auth/me       - Get current user (requires auth)")
// 	log.Println("  GET  /health        - Health check")

// 	return server
// }

// func handleRoot(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		sendJSONError(w, http.StatusNotFound, "Endpoint not found")
// 		return
// 	}

// 	info := map[string]interface{}{
// 		"service": "Authentication Server",
// 		"version": "1.0.0",
// 		"endpoints": map[string]string{
// 			"POST /auth/register": "Register new user",
// 			"POST /auth/login":    "Login user",
// 			"POST /auth/logout":   "Logout user (requires auth)",
// 			"GET /auth/me":        "Get current user (requires auth)",
// 			"GET /health":         "Health check",
// 		},
// 		"usage": map[string]interface{}{
// 			"register": map[string]string{
// 				"method": "POST",
// 				"url":    "/auth/register",
// 				"body":   "name=John&email=john@example.com&password=securepass",
// 			},
// 			"login": map[string]string{
// 				"method": "POST",
// 				"url":    "/auth/login",
// 				"body":   "email=john@example.com&password=securepass",
// 			},
// 		},
// 	}

// 	sendJSONResponse(w, http.StatusOK, info)
// }

// func handleHealth(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		sendJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
// 		return
// 	}

// 	health := map[string]interface{}{
// 		"status":    "healthy",
// 		"service":   "auth-server",
// 		"database":  "connected",
// 		"timestamp": "2024-01-01T00:00:00Z",
// 	}

// 	// Test database connection
// 	if err := Gdb.Ping(); err != nil {
// 		health["status"] = "unhealthy"
// 		health["database"] = "disconnected"
// 		health["error"] = err.Error()
// 		sendJSONResponse(w, http.StatusServiceUnavailable, health)
// 		return
// 	}

// 	sendJSONResponse(w, http.StatusOK, health)
// }

// // Middleware

// func requireAuth(authMiddleware *middleware.AuthMiddleware, next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Get session token from cookie
// 		cookie, err := r.Cookie("session_token")
// 		if err != nil {
// 			sendJSONError(w, http.StatusUnauthorized, "Authentication required")
// 			return
// 		}

// 		// Validate session directly using auth service
// 		sessionRepo := infra_repository.NewSQLiteUserSessionRepository(Gdb)
// 		userRepo := infra_repository.NewSQLiteUserRepository(Gdb)
// 		authService := usecase.NewAuthService(userRepo, sessionRepo)
// 		user, err := authService.ValidateSession(cookie.Value)
// 		if err != nil {
// 			// Clear invalid cookie
// 			http.SetCookie(w, &http.Cookie{
// 				Name:     "session_token",
// 				Value:    "",
// 				Path:     "/",
// 				MaxAge:   -1,
// 				HttpOnly: true,
// 			})
// 			sendJSONError(w, http.StatusUnauthorized, "Invalid session")
// 			return
// 		}

// 		// Add user to request context
// 		ctx := r.Context()
// 		ctx = setUserInContext(ctx, user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }

// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func LogMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
// 		next.ServeHTTP(w, r)
// 	})
// }

// // Helper functions

// func getUserFromContext(r *http.Request) *entity.User {
// 	user, ok := r.Context().Value("user").(*entity.User)
// 	if !ok {
// 		return nil
// 	}
// 	return user
// }

// func setUserInContext(ctx context.Context, user *entity.UserSession) context.Context {
// 	return context.WithValue(ctx, "user", user)
// }

// func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)

// 	if err := json.NewEncoder(w).Encode(data); err != nil {
// 		log.Printf("JSON encoding error: %v", err)
// 	}
// }

// func sendJSONError(w http.ResponseWriter, statusCode int, message string) {
// 	response := ErrorResponse{
// 		Success: false,
// 		Error:   http.StatusText(statusCode),
// 		Message: message,
// 	}
// 	sendJSONResponse(w, statusCode, response)
// }
