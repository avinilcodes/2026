package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

// shared limiter for the server
var limiter = rate.NewLimiter(5, 10)

func init() {
	blogService = service.NewBlogService()
	profileService = service.NewProfileService()
	homeService = service.NewHomeService()
	if err := service.InitDefaultWorkerPool(5, 100); err != nil {
		log.Fatalf("failed to initialize worker pool: %v", err)
	}
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// use signal.NotifyContext to listen for interrupt/terminate signals
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("Server starting on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// wait for shutdown signal
	<-sigCtx.Done()
	// allow 15s to finish ongoing requests and drain workers
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fmt.Println("Shutting down server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	// stop the worker pool, give it the same timeout
	if err := service.ShutdownDefaultPool(shutdownCtx); err != nil {
		log.Printf("worker pool shutdown error: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}

// getHome returns home page data
func getHome(w http.ResponseWriter, r *http.Request) {
	// create per-request timeout and derived context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := limiter.Wait(ctx); err != nil {
		handler.WriteError(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	home, err := homeService.GetHomeData(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			handler.WriteError(w, http.StatusGatewayTimeout, "request timed out")
			return
		}
		if errors.Is(err, context.Canceled) {
			// client cancelled
			w.WriteHeader(499)
			handler.WriteError(w, 499, "client canceled request")
			return
		}
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := limiter.Wait(ctx); err != nil {
		handler.WriteError(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	profile, err := profileService.GetProfile(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			handler.WriteError(w, http.StatusGatewayTimeout, "request timed out")
			return
		}
		if errors.Is(err, context.Canceled) {
			w.WriteHeader(499)
			handler.WriteError(w, 499, "client canceled request")
			return
		}
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := limiter.Wait(ctx); err != nil {
		handler.WriteError(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	blogs, err := blogService.GetAllBlogs(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			handler.WriteError(w, http.StatusGatewayTimeout, "request timed out")
			return
		}
		if errors.Is(err, context.Canceled) {
			w.WriteHeader(499)
			handler.WriteError(w, 499, "client canceled request")
			return
		}
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := limiter.Wait(ctx); err != nil {
		handler.WriteError(w, http.StatusTooManyRequests, "Too many requests")
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
		// service returns error when not found
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := limiter.Wait(ctx); err != nil {
		handler.WriteError(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	id, err := blogService.CreateBlog(ctx, blog)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			handler.WriteError(w, http.StatusGatewayTimeout, "request timed out")
			return
		}
		if errors.Is(err, context.Canceled) {
			w.WriteHeader(499)
			handler.WriteError(w, 499, "client canceled request")
			return
		}
		handler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Blog created successfully",
	}
	handler.WriteJSON(w, http.StatusCreated, response)
}
