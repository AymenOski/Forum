package controller

import (
	"fmt"
	"html/template"
	"net/http"

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

func (c *PostController) ShowErrorPage(w http.ResponseWriter, data ErrorMessage) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := c.templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, data.Error, data.StatusCode)
	}
}

func (pc PostController) HandleReactToPost(w http.ResponseWriter, r *http.Request) {
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

	token, err := r.Cookie("session_token")
	isAuthenticated := err == nil

	// Collect filters
	selectedCategoryNames := r.URL.Query()["category-filter"]
	wantMyPosts := r.URL.Query().Get("myPosts") != ""
	wantLikedPosts := r.URL.Query().Get("likedPosts") != ""

	// Fetch categories
	categories, err := pc.categoryService.GetAllCategories()
	if err != nil {
		pc.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusInternalServerError,
			Error:      "Could not fetch categories",
		})
		return
	}

	// Convert category names → IDs
	var selectedIDs []uuid.UUID
	selectedMap := make(map[string]bool)
	for _, selected := range selectedCategoryNames {
		for _, cat := range categories {
			if cat.Name == selected {
				selectedIDs = append(selectedIDs, cat.ID)
				selectedMap[selected] = true
			}
		}
	}

	var filteredPosts []*entity.PostWithDetails
	seen := make(map[uuid.UUID]bool)

	// Filter by category
	for _, catID := range selectedIDs {
		posts, err := pc.postService.GetPostsWithDetailsByCategoryID(catID)
		if err != nil {
			continue
		}
		for _, post := range posts {
			if !seen[post.ID] {
				filteredPosts = append(filteredPosts, post)
				seen[post.ID] = true
			}
		}
	}

	// Filter by My Posts or Liked Posts
	if (wantMyPosts || wantLikedPosts) && isAuthenticated {
		user, err := pc.authService.GetUserFromSessionToken(token.Value)
		if err == nil {
			if wantMyPosts {
				userPosts, err := pc.postService.GetPostsByUser(user.ID)
				if err == nil {
					for _, post := range userPosts {
						if !seen[post.ID] {
							filteredPosts = append(filteredPosts, post)
							seen[post.ID] = true
						}
					}
				}
			}
			if wantLikedPosts {
				likedPosts, err := pc.postService.GetLikedPostsByUser(user.ID)
				if err == nil {
					for _, post := range likedPosts {
						if !seen[post.ID] {
							filteredPosts = append(filteredPosts, post)
							seen[post.ID] = true
						}
					}
				}
			}
		}
	}

	pc.renderTemplate(w, "layout.html", map[string]interface{}{
		"posts":              filteredPosts,
		"selectedCategories": selectedMap,
		"isAuthenticated":    isAuthenticated,
		"wantMyPosts":        wantMyPosts,
		"wantLikedPosts":     wantLikedPosts,
	})
}
