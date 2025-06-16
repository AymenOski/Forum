package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"forum/usecase"

	"github.com/google/uuid"
)

type CommentController struct {
	postService     *usecase.PostService
	commentService  *usecase.CommentService
	categoryService *usecase.CategoryService
	templates       *template.Template
}

func NewCommentController(postService *usecase.PostService, commentService *usecase.CommentService,
	categoryService *usecase.CategoryService, templates *template.Template,
) *CommentController {
	return &CommentController{
		postService:     postService,
		commentService:  commentService,
		categoryService: categoryService,
		templates:       templates,
	}
}

func (cc *CommentController) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	var username string
	var isAuthenticated bool

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
	if r.Method != http.MethodPost {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}

	user, err := cc.postService.GetUserFromSessionToken(cookie.Value)
	if err == nil && user != nil {
		username = user.UserName
		isAuthenticated = true
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
		w.WriteHeader(http.StatusBadRequest)
		cc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      "Comment cannot be empty",
			"username":        username,
			"isAuthenticated": isAuthenticated,
		})
		return
	}

	// Create comment using the service
	_, err = cc.commentService.CreateComment(&postID, cookie.Value, content)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "character") || strings.Contains(err.Error(), "empty") {
			statusCode = http.StatusBadRequest
		}
		w.WriteHeader(statusCode)
		cc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      err.Error(),
			"username":        username,
			"isAuthenticated": isAuthenticated,
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

func (cc *CommentController) HandleReactToComment(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session_token")
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

	if r.Method != http.MethodPost {
		cc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}
	id := r.FormValue("CommentID")
	ID, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Failed to parse this ID %v to UUID: %v\n", id, err)
		return
	}
	like := true
	if r.FormValue("isLike") == "0" {
		like = false
	}
	cc.commentService.ReactToComment(&ID, token.Value, like)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
