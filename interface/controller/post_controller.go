package controller

import (
	"html/template"
	"net/http"

	"forum/domain/entity"
	"forum/usecase"

	"github.com/google/uuid"
)

type PostController struct {
	postService    *usecase.PostService
	commentService *usecase.CommentService
	templates      *template.Template
}

func NewPostController(postService *usecase.PostService, commentService *usecase.CommentService, templates *template.Template) *PostController {
	return &PostController{
		postService:    postService,
		commentService: commentService,
		templates:      templates,
	}
}

func (c *PostController) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*entity.User)
	if !ok || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}

	content := r.FormValue("content")
	categoryIDs := r.Form["categories"]

	if content == "" {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusBadRequest,
			Error:      "Title and content are required",
		})
		return
	}

	var categories []*uuid.UUID
	for _, catID := range categoryIDs {
		id, err := uuid.Parse(catID)
		if err == nil {
			categories = append(categories, &id)
		}
	}

	_, err := c.postService.CreatePost(&user.ID, content, categories)
	if err != nil {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *PostController) renderTemplate(w http.ResponseWriter, template string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.templates.ExecuteTemplate(w, template, data)
	if err != nil {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error rendering page",
		})
	}
}

func (c *PostController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(data.StatusCode)

	err := c.templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, data.Error, data.StatusCode)
	}
}
