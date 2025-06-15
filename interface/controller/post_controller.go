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
	var username string
	var isAuthenticated bool
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Unexpected Error While Reading Cookie",
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
	// flag-1: next field is temperoraly until we have a proper middleware
	user, err := pc.postService.GetUserFromSessionToken(cookie.Value)
	if err == nil && user != nil {
		username = user.UserName
		isAuthenticated = true
	}
	posts, err := pc.postService.GetPosts()
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Something went wrong while loading posts",
		})
		return
	}

	if content == "" {
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      usecase.ErrEmptyPostContent,
			"username":        username,
			"isAuthenticated": isAuthenticated,
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
				"username":        username,
				"isAuthenticated": isAuthenticated,
			})
			return
		}
		categoriesIDs = append(categoriesIDs, &c.ID)
	}

	_, err = pc.postService.CreatePost(cookie.Value, content, categoriesIDs)
	if err != nil {
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"form_error":      err.Error(),
			"Content":         content,
			"posts":           posts,
			"username":        username,
			"isAuthenticated": isAuthenticated,
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



func (pc PostController) HandleReactToPost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session_token")
	// Token, err := uuid.Parse(token.Value)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Unexpected Error While Reading Cookie",
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
	
	id := strings.Split(r.URL.Query().Get("id"), "/")

	ID, err := uuid.Parse(id[0])
	if err != nil {
		fmt.Printf("Failed to parse this ID %v to UUID: %v\n", id, err)
		return
	}
	like := true
	if id[1] == "0" {
		like = false
	}
	pc.postService.ReactToPost(ID, token.Value, like)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}



func (c *PostController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(data.StatusCode)

	err := c.templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, data.Error, data.StatusCode)
	}
}


