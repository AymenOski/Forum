package controller

import (
	"html/template"
	"net/http"
	"time"

	"forum/domain/entity"
	"forum/usecase"

	"github.com/google/uuid"
)

type PostController struct {
	postService     *usecase.PostService
	commentService  *usecase.CommentService
	categoryService *usecase.CategoryService
	templates       *template.Template
	authservice     *usecase.AuthService
}

func NewPostController(postService *usecase.PostService, commentService *usecase.CommentService, categoryService *usecase.CategoryService, templates *template.Template, authservice *usecase.AuthService) *PostController {
	return &PostController{
		postService:     postService,
		commentService:  commentService,
		categoryService: categoryService,
		templates:       templates,
		authservice:     authservice,
	}
}

func (pc *PostController) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session_token")
	if err != nil || token == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	if content == "" {
		allCategories, _ := pc.categoryService.GetAllCategories()
		posts, _ := pc.postService.GetPosts()

		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"form_error": "Title and content are required",
			"Content":    content,
			"posts":      posts,
			"categories": allCategories,
		})
		return
	}

	// Verify if the categories exist
	categoriesIDs := make([]*uuid.UUID, 0, len(categories))
	for _, cat := range categories {
		c, err := pc.categoryService.GetCategoryByName(cat)
		if err != nil {
			allCategories, _ := pc.categoryService.GetAllCategories()
			posts, _ := pc.postService.GetPosts()

			pc.renderTemplate(w, "layout.html", map[string]interface{}{
				"form_error": "Invalid category: " + cat,
				"Content":    content,
				"posts":      posts,
				"categories": allCategories,
			})
			return
		}
		categoriesIDs = append(categoriesIDs, &c.ID)
	}

	// Try to create post
	_, err = pc.postService.CreatePost(token.Value, content, categoriesIDs)
	if err != nil {
		allCategories, _ := pc.categoryService.GetAllCategories()
		posts, _ := pc.postService.GetPosts()

		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"form_error": err.Error(),
			"Content":    content,
			"posts":      posts,
			"categories": allCategories,
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
func (pc *PostController) ShowPostsPage(w http.ResponseWriter, r *http.Request) {
	categories, err := pc.categoryService.GetAllCategories()
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error fetching categories",
		})
		return
	}

	posts, err := pc.postService.GetPosts()
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error fetching posts",
		})
		return
	}

	pc.renderTemplate(w, "layout.html", map[string]interface{}{
		"posts":      posts,
		"categories": categories,
	})
}
func (pc *PostController) GetCurrentUserID(r *http.Request) (uuid.UUID, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return uuid.UUID{}, false
	}

	session, err := pc.authservice.ValidateSession(cookie.Value)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		return uuid.UUID{}, false
	}

	return session.UserID, true
}

func (pc *PostController) HandleFilteredPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := entity.PostFilter{}

	// ✅ Category filter (multi-select)
	if categories := query["category"]; len(categories) > 0 {
		for _, cat := range categories {
			id, err := uuid.Parse(cat)
			if err == nil {
				filter.CategoryIDs = append(filter.CategoryIDs, id)
			}
		}
	}

	// ✅ "My Posts" filter
	if query.Get("myPosts") == "on" {
		userID, ok := pc.GetCurrentUserID(r)
		if !ok {
			pc.ShowErrorPage(w, ErrorMessage{
				StatusCode: http.StatusUnauthorized,
				Error:      "You must be logged in to view your posts.",
			})
			return
		}
		filter.MyPosts = true
		filter.AuthorID = &userID
	}

	// ✅ "Liked Posts" filter
	if query.Get("likedPosts") == "on" {
		userID, ok := pc.GetCurrentUserID(r)
		if !ok {
			pc.ShowErrorPage(w, ErrorMessage{
				StatusCode: http.StatusUnauthorized,
				Error:      "You must be logged in to view liked posts.",
			})
			return
		}
		filter.LikedPosts = true
		filter.AuthorID = &userID
	}

	// 🚀 Fetch posts
	posts, err := pc.postService.GetFilteredPosts(filter)
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error fetching filtered posts",
		})
		return
	}

	// ✅ Fetch all categories for the form
	categories, _ := pc.categoryService.GetAllCategories()

	// 🎯 Render
	pc.renderTemplate(w, "layout.html", map[string]interface{}{
		"posts":      posts,
		"filter":     filter,
		"categories": categories,
	})
}
