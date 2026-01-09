package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"my-website/handler"

	"my-website/models"
	"my-website/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

var (
	blogService    *service.BlogService
	profileService *service.ProfileService
	homeService    *service.HomeService
)

func init() {
	blogService = service.NewBlogService()
	profileService = service.NewProfileService()
	homeService = service.NewHomeService()
}

func main() {
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/api/home", getHome).Methods(http.MethodGet)
	router.HandleFunc("/api/profile", getProfile).Methods(http.MethodGet)
	router.HandleFunc("/api/blogs", getBlogs).Methods(http.MethodGet)
	router.HandleFunc("/api/blogs/{id}", getBlogByID).Methods(http.MethodGet)
	router.HandleFunc("/api/blogs", createBlog).Methods(http.MethodPost)

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(router)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// getHome returns home page data
func getHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limiter := rate.NewLimiter(5, 10)

	if err := limiter.Wait(ctx); err != nil {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	home, err := homeService.GetHomeData(ctx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	homeResponse := models.HomeResponse{
		Title:       home.Title,
		Description: home.Description,
		Name:        home.Name,
	}
	handler.WriteJSON(w, http.StatusOK, homeResponse)
}

// getProfile returns profile information
func getProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limiter := rate.NewLimiter(5, 10)

	if err := limiter.Wait(ctx); err != nil {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	profile, err := profileService.GetProfile(ctx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	profileResponse := models.ProfileResponse{
		Name:     profile.Name,
		Title:    profile.Title,
		Bio:      profile.Bio,
		Email:    profile.Email,
		Skills:   profile.Skills,
		GitHub:   profile.GitHub,
		LinkedIn: profile.LinkedIn,
	}
	handler.WriteJSON(w, http.StatusOK, profileResponse)
}

// getBlogs returns all blogs
func getBlogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limiter := rate.NewLimiter(5, 10)

	if err := limiter.Wait(ctx); err != nil {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	blogs, err := blogService.GetAllBlogs(ctx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	blogResponses := make([]models.BlogResponse, len(blogs))
	for i, blog := range blogs {
		blogResponses[i] = models.BlogResponse{
			ID:      blog.ID,
			Title:   blog.Title,
			Content: blog.Content,
			Author:  blog.Author,
			Date:    blog.Date,
			Tags:    blog.Tags,
		}
	}

	handler.WriteJSON(w, http.StatusOK, blogResponses)
}

// getBlogByID returns a single blog by ID
func getBlogByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limiter := rate.NewLimiter(5, 10)

	if err := limiter.Wait(ctx); err != nil {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	vars := mux.Vars(r)
	blogID, err := strconv.Atoi(vars["id"])
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	blog, err := blogService.GetBlogByID(ctx, blogID)
	if err != nil {
		handler.WriteError(w, http.StatusNotFound, "Blog not found")
		return
	}

	blogResponse := models.BlogResponse{
		ID:      blog.ID,
		Title:   blog.Title,
		Content: blog.Content,
		Author:  blog.Author,
		Date:    blog.Date,
		Tags:    blog.Tags,
	}

	handler.WriteJSON(w, http.StatusOK, blogResponse)
}

// createBlog creates a new blog post
func createBlog(w http.ResponseWriter, r *http.Request) {
	var blogReq models.BlogRequest
	if err := json.NewDecoder(r.Body).Decode(&blogReq); err != nil {
		handler.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	blog := service.Blog{
		Title:   blogReq.Title,
		Content: blogReq.Content,
		Author:  blogReq.Author,
		Tags:    blogReq.Tags,
	}

	ctx := r.Context()

	limiter := rate.NewLimiter(5, 10)

	if err := limiter.Wait(ctx); err != nil {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	id, err := blogService.CreateBlog(ctx, blog)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Blog created successfully",
	}
	handler.WriteJSON(w, http.StatusCreated, response)
}
