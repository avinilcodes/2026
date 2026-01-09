package dao

// Blog represents a blog post in the database
type Blog struct {
	ID      int
	Title   string
	Content string
	Author  string
	Date    string
	Tags    []string
}

// Profile represents user profile in the database
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

// BlogStore interface defines blog data operations
type BlogStore interface {
	GetAllBlogs() ([]Blog, error)
	GetBlogByID(id int) (*Blog, error)
	CreateBlog(blog Blog) (int, error)
	UpdateBlog(blog Blog) error
	DeleteBlog(id int) error
}

// ProfileStore interface defines profile data operations
type ProfileStore interface {
	GetProfile() (*Profile, error)
	UpdateProfile(profile Profile) error
}

// HomeStore interface defines home data operations
type HomeStore interface {
	GetHomeData() (*Home, error)
}
