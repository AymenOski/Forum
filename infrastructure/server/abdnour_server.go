package server

import (
	"database/sql"
	"html/template"
	"net/http"

	infra_repository "forum/infrastructure/repository"
	"forum/interface/controller"
	"forum/usecase"
)

func MyServer(db *sql.DB, templates *template.Template) *http.Server {
	router := http.NewServeMux()

	// Static files handler
	staticFileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// Domain layer
	userRepo := infra_repository.NewSQLiteUserRepository(db)
	sessionRepo := infra_repository.NewSQLiteUserSessionRepository(db)
	postRepo := infra_repository.NewSQLitePostRepository(db)
	postCategoryRepo := infra_repository.NewSQLitePostCategoryRepository(db)
	categoryRepo := infra_repository.NewSQLiteCategoryRepository(db)
	postReactionRepo := infra_repository.NewSQLitePostReactionRepository(db)
	commentRepo := infra_repository.NewSQLiteCommentRepository(db)
	commentReactionRepo := infra_repository.NewSQLiteCommentReactionRepository(db)

	postAggregateRepo := infra_repository.NewSQLitePostAggregateRepository(
		db,
		&postRepo,
		&postCategoryRepo,
		&userRepo,
		&postReactionRepo,
	)

	// Service layer (use cases)
	authService := usecase.NewAuthService(userRepo, sessionRepo)
	postService := usecase.NewPostService(
		&postRepo,
		&userRepo,
		&categoryRepo,
		&postAggregateRepo,
		&postReactionRepo,
		&sessionRepo,
	)
	commentService := usecase.NewCommentService(
		userRepo,
		commentRepo,
		postRepo,
		commentReactionRepo,
	)
	categoryService := usecase.NewCategoryService(
		categoryRepo,
		postCategoryRepo,
		sessionRepo,
		userRepo,
	)

	// Presentation layer (controllers)
	authController := controller.NewAuthController(authService, postService, templates)
	postController := controller.NewPostController(
		postService,
		commentService,
		categoryService,
		templates,
	)

	// Route handlers
	router.HandleFunc("/signup", authController.HandleSignup)
	router.HandleFunc("/login", authController.HandleLogin)
	router.HandleFunc("/post/create", postController.HandleCreatePost)
	router.HandleFunc("/", authController.HandleMainPage)

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
