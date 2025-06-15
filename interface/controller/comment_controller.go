package controller

import (
	"html/template"
	"net/http"

	"forum/usecase"

	"github.com/google/uuid"
)

type CommentController struct {
	// flag-1: next field is temperoraly until we have a proper middleware
	commentService *usecase.CommentService
	postService    *usecase.PostService
	templates      *template.Template
}

func NewCommentController(commentService *usecase.CommentService, postService *usecase.PostService, templates *template.Template) *CommentController {
	return &CommentController{
		commentService: commentService,
		postService:    postService,
		templates:      templates,
	}
}

func (cc *CommentController) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	var username string
	var isAuthenticated bool
	// Check session token for authentication
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if err != nil {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Unexpected Error While Reading Cookie",
		})
		return
	}
	// flag-1: next field is temperoraly until we have a proper middleware
	user, err := cc.commentService.GetUserFromSessionToken(cookie.Value)
	if err == nil && user != nil {
		username = user.UserName
		isAuthenticated = true
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

	posts, err := cc.postService.GetPosts()
	if err != nil {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Something Went Wrong While Loading Posts",
		})
		return
	}

	content := r.FormValue("content")
	if content == "" {
		cc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      "Comment cannot be empty",
			"username":        username,
			"isAuthenticated": isAuthenticated,
		})
		return
	}

	// Create comment using the service
	_, err = cc.commentService.CreateComment(&postID, &user.ID, content)
	if err != nil {
		cc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      err.Error(),
			"username":        username,
			"isAuthenticated": isAuthenticated,
		})
		return
	}

	http.Redirect(w, r, "/?succed=true", http.StatusSeeOther)
}

func (cc *CommentController) renderTemplate(w http.ResponseWriter, template string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := cc.templates.ExecuteTemplate(w, template, data)
	if err != nil {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error Rendering Page",
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
