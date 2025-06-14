package controller

import (
	"fmt"
	"html/template"
	"net/http"

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
	id := r.FormValue("postId")
	content := r.FormValue("content")
	ID, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Failed to parse this ID %v to UUID: %v\n", id, err)
		return
	}
	cc.commentService.CreateComment(&ID,token.Value,content)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (cc *CommentController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(data.StatusCode)

	err := cc.templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, data.Error, data.StatusCode)
	}
}
