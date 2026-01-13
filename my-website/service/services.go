package service

import (
	"context"
	"errors"
	"sync"
)

// NOTE: Default worker pool must be initialized by the application
// via InitDefaultWorkerPool or SetDefaultWorkerPool to make ownership explicit.

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
	task := func(c context.Context) (interface{}, error) {
		bs.mu.RLock()
		defer bs.mu.RUnlock()
		blogs := make([]Blog, 0, len(bs.blogs))
		for _, blog := range bs.blogs {
			blogs = append(blogs, blog)
		}
		return blogs, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return nil, err
	}
	res, err := wp.Enqueue(ctx, task)
	if err != nil {
		return nil, err
	}
	return res.([]Blog), nil
}

// GetBlogByID returns a blog by ID
func (bs *BlogService) GetBlogByID(ctx context.Context, id int) (*Blog, error) {
	task := func(c context.Context) (interface{}, error) {
		bs.mu.RLock()
		defer bs.mu.RUnlock()
		blog, exists := bs.blogs[id]
		if !exists {
			return nil, errors.New("blog not found")
		}
		return &blog, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return nil, err
	}
	res, err := wp.Enqueue(ctx, task)
	if err != nil {
		return nil, err
	}
	return res.(*Blog), nil
}

// CreateBlog creates a new blog
func (bs *BlogService) CreateBlog(ctx context.Context, blog Blog) (int, error) {
	task := func(c context.Context) (interface{}, error) {
		bs.mu.Lock()         // Lock for writing
		defer bs.mu.Unlock() // Unlock after writing
		blog.ID = bs.nextID
		bs.blogs[blog.ID] = blog
		bs.nextID++
		return blog.ID, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return 0, err
	}
	res, err := wp.Enqueue(ctx, task)
	if err != nil {
		return 0, err
	}
	return res.(int), nil
}

// UpdateBlog updates an existing blog
func (bs *BlogService) UpdateBlog(ctx context.Context, blog Blog) error {
	task := func(c context.Context) (interface{}, error) {
		bs.mu.Lock()
		defer bs.mu.Unlock()
		if _, exists := bs.blogs[blog.ID]; !exists {
			return nil, errors.New("blog not found")
		}
		bs.blogs[blog.ID] = blog
		return nil, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return err
	}
	_, err = wp.Enqueue(ctx, task)
	return err
}

// DeleteBlog deletes a blog by ID
func (bs *BlogService) DeleteBlog(ctx context.Context, id int) error {
	task := func(c context.Context) (interface{}, error) {
		bs.mu.Lock()
		defer bs.mu.Unlock()
		if _, exists := bs.blogs[id]; !exists {
			return nil, errors.New("blog not found")
		}
		delete(bs.blogs, id)
		return nil, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return err
	}
	_, err = wp.Enqueue(ctx, task)
	return err
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
			Name:  "Avinil Bedarkar",
			Title: "Member of Technical Staff 2",
			Bio:   "I am a passionate learner and want to make an impact in the Information Age using my knowledge and skills.",
			Email: "bedarkar.avinil01@gmail.com",
			Skills: []string{
				"Aerospike",
				"Angular",
				"AWS",
				"BDD",
				"Caching",
				"Data Structure",
				"Design Patterns",
				"Docker",
				"Elasticsearch",
				"Encryption",
				"Git",
				"Github",
				"gRPC",
				"Java",
				"JavaScript",
				"Jenkins",
				"LDAP",
				"Integration testing",
				"Microservice",
				"Oracle",
				"PKI",
				"Postgresql",
				"Python",
				"Redis",
				"Restful",
				"TDD",
				"Unit Testing",
				"Web Services",
				"LLM",
				"GenAI",
				"RADIUS",
			},
			GitHub:   "https://github.com/avinilcodes",
			LinkedIn: "https://www.linkedin.com/in/avinil-bedarkar/",
		},
	}
}

// GetProfile returns the profile
func (ps *ProfileService) GetProfile(ctx context.Context) (*Profile, error) {
	task := func(c context.Context) (interface{}, error) {
		ps.mu.RLock()
		defer ps.mu.RUnlock()
		return &ps.profile, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return nil, err
	}
	res, err := wp.Enqueue(ctx, task)
	if err != nil {
		return nil, err
	}
	return res.(*Profile), nil
}

// UpdateProfile updates the profile
func (ps *ProfileService) UpdateProfile(ctx context.Context, profile Profile) error {
	task := func(c context.Context) (interface{}, error) {
		ps.mu.Lock()
		defer ps.mu.Unlock()
		ps.profile = profile
		return nil, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return err
	}
	_, err = wp.Enqueue(ctx, task)
	return err
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
			Title:       "Avinil Bedarkar — Member of Technical Staff 2",
			Description: "Engineer focused on encryption, PKI, and integrating LLMs with RAG for product features.",
			Name:        "Avinil Bedarkar",
		},
	}
}

// GetHomeData returns home page data
func (hs *HomeService) GetHomeData(ctx context.Context) (*Home, error) {
	task := func(c context.Context) (interface{}, error) {
		hs.mu.RLock()         // Read lock
		defer hs.mu.RUnlock() // Unlock after reading
		return &hs.home, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return nil, err
	}
	res, err := wp.Enqueue(ctx, task)
	if err != nil {
		return nil, err
	}
	return res.(*Home), nil
}

// UpdateHomeData updates home page data
func (hs *HomeService) UpdateHomeData(ctx context.Context, home Home) error {
	task := func(c context.Context) (interface{}, error) {
		hs.mu.Lock()
		defer hs.mu.Unlock()
		hs.home = home
		return nil, nil
	}
	wp, err := GetDefaultPool()
	if err != nil {
		return err
	}
	_, err = wp.Enqueue(ctx, task)
	return err
}
