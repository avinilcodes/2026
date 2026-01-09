package models

// HomeResponse represents the home page data
type HomeResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

// ProfileResponse represents user profile information
type ProfileResponse struct {
	Name     string   `json:"name"`
	Title    string   `json:"title"`
	Bio      string   `json:"bio"`
	Email    string   `json:"email"`
	Skills   []string `json:"skills"`
	GitHub   string   `json:"github"`
	LinkedIn string   `json:"linkedin"`
}

// BlogResponse represents a blog post response
type BlogResponse struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Author  string   `json:"author"`
	Date    string   `json:"date"`
	Tags    []string `json:"tags"`
}

// BlogRequest represents a blog post creation request
type BlogRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Author  string   `json:"author"`
	Tags    []string `json:"tags"`
}
