package controller

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"forum/domain/entity"
	"forum/usecase"

	"github.com/google/uuid"
)

type PostController struct {
	postService     *usecase.PostService
	commentService  *usecase.CommentService
	categoryService *usecase.CategoryService
	authService     *usecase.AuthService
	templates       *template.Template
}

func NewPostController(postService *usecase.PostService, commentService *usecase.CommentService,
	categoryService *usecase.CategoryService, authService *usecase.AuthService, templates *template.Template,
) *PostController {
	return &PostController{
		postService:     postService,
		commentService:  commentService,
		categoryService: categoryService,
		authService:     authService,
		templates:       templates,
	}
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

	if len(categories) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      "Please select at least one category",
			"username":        user.UserName,
			"isAuthenticated": true,
		})
		return
	}

	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"posts":           posts,
			"form_error":      usecase.ErrEmptyPostContent,
			"username":        username,
			"isAuthenticated": isAuthenticated,
		})
		return
	}

	categoriesIDs := make([]*uuid.UUID, 0, len(categories))
	for _, category := range categories {
		c, err := pc.categoryService.GetCategoryByName(category)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "wait a bit") {
			statusCode = http.StatusTooManyRequests // 429 for rate limiting
		} else if strings.Contains(err.Error(), "content") {
			statusCode = http.StatusBadRequest // 400 for validation errors
		}
		w.WriteHeader(statusCode)
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

func (pc PostController) HandleReactToPost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session_token")
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
	id := r.FormValue("postId")
	ID, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Failed to parse this ID %v to UUID: %v\n", id, err)
		return
	}
	like := true
	if r.FormValue("isLike") == "0" {
		like = false
	}
	pc.postService.ReactToPost(ID, token.Value, like)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (pc *PostController) HandleFilteredPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method not allowed",
		})
		return
	}

	var userID *uuid.UUID
	var isAuthenticated bool
	var username string

	cookie, err := r.Cookie("session_token")
	if err == nil {
		user, err := pc.postService.GetUserFromSessionToken(cookie.Value)
		if err == nil && user != nil {
			isAuthenticated = true
			userID = &user.ID
			username = user.UserName
		}
	}

	posts, err := pc.postService.GetPosts()
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Something went wrong while loading posts",
		})
		return
	}

	// Get query parameters
	selectedCategoryNames := r.URL.Query()["category-filter"]
	Radio := r.URL.Query().Get("postFilter")
	var likedPosts, myPosts bool = false, false
	if Radio == "myPosts" {
		myPosts = true
	}
	if Radio == "likedPosts" {
		likedPosts = true
	}

	hasFilters := len(selectedCategoryNames) > 0 || myPosts || likedPosts
	if !hasFilters {
		pc.renderTemplate(w, "layout.html", map[string]interface{}{
			"form_error":      errors.New("No filter is selected"),
			"posts":           posts,
			"username":        username,
			"isAuthenticated": isAuthenticated,
		})
		return
	}

	// Get all categories for name-to-ID conversion
	categories, err := pc.categoryService.GetAllCategories()
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Could not fetch categories",
		})
		return
	}

	// Convert category names to IDs
	var selectedIDs []uuid.UUID
	selectedMap := make(map[string]bool)
	for _, selected := range selectedCategoryNames {
		for _, cat := range categories {
			if cat.Name == selected {
				selectedIDs = append(selectedIDs, cat.ID)
				selectedMap[selected] = true
				break // Found the category, no need to continue inner loop
			}
		}
	}

	// Build filter
	filter := &entity.PostFilter{
		CategoryIDs: selectedIDs,
		MyPosts:     myPosts,
		LikedPosts:  likedPosts,
		AuthorID:    userID,
	}

	// Get filtered posts using the service
	filteredPosts, err := pc.postService.GetFilteredPostsWithDetails(*filter)
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Error filtering posts: " + err.Error(),
		})
		return
	}

	pc.renderTemplate(w, "layout.html", map[string]interface{}{
		"username":           username,
		"isAuthenticated":    isAuthenticated,
		"posts":              filteredPosts,
		"selectedCategories": selectedMap,
	})
}

func (c *PostController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(data.StatusCode)
	err := c.templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, data.Error, data.StatusCode)
	}
}
