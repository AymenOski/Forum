package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"forum/usecase"

	"github.com/google/uuid"
)

type PostController struct {
	postService     *usecase.PostService
	commentService  *usecase.CommentService
	categoryService *usecase.CategoryService
	templates       *template.Template
}

func NewPostController(postService *usecase.PostService, commentService *usecase.CommentService,
	categoryService *usecase.CategoryService, templates *template.Template,
) *PostController {
	return &PostController{
		postService:     postService,
		commentService:  commentService,
		categoryService: categoryService,
		templates:       templates,
	}
}

func (pc *PostController) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Unexpected error while reading cookie",
		})
		return
	}

	if r.Method != http.MethodPost {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}

	content := r.FormValue("content")
	categories := r.Form["categories"]

	posts, err := pc.postService.GetPosts()
	if err != nil {
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      usecase.ErrPostNotFound,
			"username":        nil,
			"isAuthenticated": nil,
		})
		return
	}

	if content == "" {
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      usecase.ErrEmptyPostContent,
			"username":        nil,
			"isAuthenticated": nil,
		})
		return
	}

	// verify if the categories exist
	categoriesIDs := make([]*uuid.UUID, 0, len(categories))
	for _, cat := range categories {
		c, err := pc.categoryService.GetCategoryByName(cat)
		if err != nil {
			pc.renderTemplate(w, "layout.html", map[string]interface{}{
				"posts":           posts,
				"form_error":      usecase.ErrCategoryNotFound,
				"username":        nil,
				"isAuthenticated": nil,
			})
			return
		}
		categoriesIDs = append(categoriesIDs, &c.ID)
	}

	_, err = pc.postService.CreatePost(token.Value, content, categoriesIDs)
	if err != nil {
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"form_error":      err.Error(),
			"Content":         content,
			"posts":           posts,
			"username":        nil,
			"isAuthenticated": nil,
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

func (pc *PostController) HandleReactToPost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session_token")
	// Token, err := uuid.Parse(token.Value)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Unexpected error while reading cookie",
		})
		return
	}

	if r.Method != http.MethodPost {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}
	
	id:=strings.Split(r.URL.Query().Get("id"), "/")
	
	ID, err := uuid.Parse(id[0])
	if err != nil {
		fmt.Printf("Failed to parse this ID %v to UUID: %v\n",id, err)
		return
	}
	like := true
	if id[1] == "0" {
		like = false
	}
	pc.postService.ReactToPost(&ID, token.Value, like)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
