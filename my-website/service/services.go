package service

import (
	"context"
	"errors"
	"sync"
)

// Blog represents a blog post
type Blog struct {
	ID      int
	Title   string
	Content string
	Author  string
	Date    string
	Tags    []string
}

// Profile represents user profile
type Profile struct {
	Name     string
	Title    string
	Bio      string
	Email    string
	Skills   []string
	GitHub   string
	LinkedIn string
}

// Home represents home page data
type Home struct {
	Title       string
	Description string
	Name        string
}

// BlogService handles blog business logic
type BlogService struct {
	blogs  map[int]Blog
	mu     sync.RWMutex
	nextID int
}

// NewBlogService creates a new blog service with sample data
func NewBlogService() *BlogService {
	service := &BlogService{
		blogs:  make(map[int]Blog),
		nextID: 3,
	}

	// Initialize with sample data
	service.blogs[1] = Blog{
		ID:      1,
		Title:   "Getting Started with Go",
		Content: "Learn the basics of Go programming language...",
		Author:  "Your Name",
		Date:    "2026-01-09",
		Tags:    []string{"go", "programming"},
	}
	service.blogs[2] = Blog{
		ID:      2,
		Title:   "Angular Best Practices",
		Content: "Explore best practices when developing Angular applications...",
		Author:  "Your Name",
		Date:    "2026-01-08",
		Tags:    []string{"angular", "frontend"},
	}

	return service
}

// GetAllBlogs returns all blogs
func (bs *BlogService) GetAllBlogs(ctx context.Context) ([]Blog, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	blogs := make([]Blog, 0, len(bs.blogs))
	for _, blog := range bs.blogs {
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

// GetBlogByID returns a blog by ID
func (bs *BlogService) GetBlogByID(ctx context.Context, id int) (*Blog, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	blog, exists := bs.blogs[id]
	if !exists {
		return nil, errors.New("blog not found")
	}
	return &blog, nil
}

// CreateBlog creates a new blog
func (bs *BlogService) CreateBlog(ctx context.Context, blog Blog) (int, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	blog.ID = bs.nextID
	bs.blogs[blog.ID] = blog
	bs.nextID++

	return blog.ID, nil
}

// UpdateBlog updates an existing blog
func (bs *BlogService) UpdateBlog(ctx context.Context, blog Blog) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if _, exists := bs.blogs[blog.ID]; !exists {
		return errors.New("blog not found")
	}
	bs.blogs[blog.ID] = blog
	return nil
}

// DeleteBlog deletes a blog by ID
func (bs *BlogService) DeleteBlog(ctx context.Context, id int) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if _, exists := bs.blogs[id]; !exists {
		return errors.New("blog not found")
	}
	delete(bs.blogs, id)
	return nil
}

// ProfileService handles profile business logic
type ProfileService struct {
	profile Profile
	mu      sync.RWMutex
}

// NewProfileService creates a new profile service
func NewProfileService() *ProfileService {
	return &ProfileService{
		profile: Profile{
			Name:     "Your Name",
			Title:    "Software Developer",
			Bio:      "Passionate about building scalable web applications",
			Email:    "your.email@example.com",
			Skills:   []string{"Go", "Angular", "Docker", "Kubernetes"},
			GitHub:   "https://github.com/yourname",
			LinkedIn: "https://linkedin.com/in/yourname",
		},
	}
}

// GetProfile returns the profile
func (ps *ProfileService) GetProfile(ctx context.Context) (*Profile, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return &ps.profile, nil
}

// UpdateProfile updates the profile
func (ps *ProfileService) UpdateProfile(ctx context.Context, profile Profile) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.profile = profile
	return nil
}

// HomeService handles home page business logic
type HomeService struct {
	home Home
	mu   sync.RWMutex
}

// NewHomeService creates a new home service
func NewHomeService() *HomeService {
	return &HomeService{
		home: Home{
			Title:       "Welcome to My Portfolio",
			Description: "Showcasing my work, projects, and blog posts",
			Name:        "Your Name",
		},
	}
}

// GetHomeData returns home page data
func (hs *HomeService) GetHomeData(ctx context.Context) (*Home, error) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	return &hs.home, nil
}

// UpdateHomeData updates home page data
func (hs *HomeService) UpdateHomeData(ctx context.Context, home Home) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.home = home
	return nil
}
