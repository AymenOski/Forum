package controller

import (
	"html/template"
	"net/http"

	"forum/domain/entity"
	"forum/domain/repository"
	"forum/usecase"

	"github.com/google/uuid"
)

type CommentController struct {
	commentService *usecase.CommentService
	postService    *usecase.PostService       // To fetch posts for template context
	userRepo       *repository.UserRepository // To fetch user by ID or session
	templates      *template.Template
}

func NewCommentController(commentService *usecase.CommentService, postService *usecase.PostService, userRepo *repository.UserRepository, templates *template.Template) *CommentController {
	return &CommentController{
		commentService: commentService,
		postService:    postService,
		userRepo:       userRepo,
		templates:      templates,
	}
}

func (cc *CommentController) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	// Check session token for authentication
	_, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if err != nil {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Unexpected error while reading cookie",
		})
		return
	}

	// Validate HTTP method
	if r.Method != http.MethodPost {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}

	// Parse form data
	postIDStr := r.FormValue("postId")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid post ID",
		})
		return
	}

	content := r.FormValue("content")
	if content == "" {
		posts, err := cc.postService.GetPosts()
		if err != nil {
			posts = []*entity.PostWithDetails{} // Fallback to empty slice instead of nil
		}
		cc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      "Comment cannot be empty",
			"username":        "userNamessssssssssssssssssssssssssssss",
			"isAuthenticated": true,
		})
		return
	}

	// Get authenticated user from context
	// userID, ok := r.Context().Value("currentUserID").(uuid.UUID)
	// if !ok {
	// 	cc.ShowErrorPage(w, ErrorMessage{
	// 		StatusCode: http.StatusUnauthorized,
	// 		Error:      "Invalid session",
	// 	})
	// 	return
	// }

	// Hardcoded userID for testing (replace with proper auth logic later)
	// Example UUID - replace with a valid user ID from the session
	hardcodedUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	userID := hardcodedUserID

	// Create comment using the service
	_, err = cc.commentService.CreateComment(&postID, &userID, content)
	if err != nil {
		posts, err := cc.postService.GetPosts()
		if err != nil {
			posts = []*entity.PostWithDetails{} // Fallback to empty slice instead of nil
		}
		cc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			// "form_error":      err.Error(),
			"username":        "userNamessssssssssssssssssssssssssssss",
			"isAuthenticated": true,
			// "username": func() any {
			// 	if isAuthenticated {
			// 		return user.UserName
			// 	}
			// 	return nil
			// }(),
			// "isAuthenticated": isAuthenticated,
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (cc *CommentController) renderTemplate(w http.ResponseWriter, template string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := cc.templates.ExecuteTemplate(w, template, data)
	if err != nil {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error rendering page",
		})
	}
}

func (cc *CommentController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(data.StatusCode)
	err := cc.templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, data.Error, data.StatusCode)
	}
}
