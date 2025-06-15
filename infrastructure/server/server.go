package server

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	infra_repository "forum/infrastructure/repository"
	"forum/interface/controller"
	"forum/interface/middleware"
	"forum/usecase"
)

var tmpl1 *template.Template


func init() {

	var err error
	tmpl1, err = template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Printf("Warning: Failed to initialize templates: %v", err)
	}
}

func MyServer(db *sql.DB) *http.Server {
	mux := http.NewServeMux()

	// Entity layer
	user_infra_repo := infra_repository.NewSQLiteUserRepository(db)
	session_infra_repo := infra_repository.NewSQLiteUserSessionRepository(db)
	post_infra_repo := infra_repository.NewSQLitePostRepository(db)
	postCategory_infra_repo := infra_repository.NewSQLitePostCategoryRepository(db)
	category_infra_repo := infra_repository.NewSQLiteCategoryRepository(db)
	post_reaction_infra_repo := infra_repository.NewSQLitePostReactionRepository(db)

	comment_reaction_infra_repo := infra_repository.NewSQLiteCommentReactionRepository(db)

	comment_infra_repo := infra_repository.NewSQLiteCommentRepository(db, &user_infra_repo, &comment_reaction_infra_repo)

	post_category_infra_repo := infra_repository.NewSQLitePostAggregateRepository(db, &post_infra_repo, &postCategory_infra_repo,
		&user_infra_repo, &post_reaction_infra_repo, &comment_infra_repo)

	auth_usecase := usecase.NewAuthService(user_infra_repo, session_infra_repo)
	rate_limiter := usecase.NewPostRateLimiter()
	post_usecase := usecase.NewPostService(&post_infra_repo, &user_infra_repo, &category_infra_repo, &post_category_infra_repo, &post_reaction_infra_repo, &session_infra_repo, rate_limiter)
	comment_usecase := usecase.NewCommentService(user_infra_repo, comment_infra_repo, post_infra_repo, session_infra_repo, comment_reaction_infra_repo)
	category_usecase := usecase.NewCategoryService(category_infra_repo, postCategory_infra_repo, session_infra_repo, user_infra_repo)
	auth_controller := controller.NewAuthController(auth_usecase, post_usecase, tmpl1)

	post_controller := controller.NewPostController(post_usecase, comment_usecase, category_usecase, auth_usecase, tmpl1)

	comment_controller := controller.NewCommentController(post_usecase, comment_usecase, category_usecase, tmpl1)

	auth_middleware := middleware.NewAuthMiddleware(auth_usecase)

	mux.HandleFunc("/signup", auth_controller.HandleSignup)
	mux.HandleFunc("/login", auth_controller.HandleLogin)
	mux.HandleFunc("/logout", auth_controller.HandleLogout)
	mux.HandleFunc("/post/create", auth_middleware.VerifiedAuth(post_controller.HandleCreatePost))
	mux.HandleFunc("/filter", post_controller.HandleFilteredPosts)
	mux.HandleFunc("/likeposts", post_controller.HandleReactToPost)
	mux.HandleFunc("/likecomment", comment_controller.HandleReactToComment)
	mux.HandleFunc("/comment/create", comment_controller.HandleCreateComment)
	mux.HandleFunc("/", auth_controller.HandleRoot)

	server := &http.Server{
		Addr:    ":8080",
		Handler: auth_middleware.LoggerMiddleware(mux),
	}

	return server
}
