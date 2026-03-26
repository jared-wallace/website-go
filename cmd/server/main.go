package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	db        *sql.DB
	templates map[string]*template.Template
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Excerpt   string    `json:"excerpt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Published bool      `json:"published"`
	Slug      string    `json:"slug"`
}

type Project struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	URL          string    `json:"url"`
	ImageURL     string    `json:"image_url"`
	Technologies string    `json:"technologies"`
	CreatedAt    time.Time `json:"created_at"`
}

type NewsItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	PubDate     string `json:"pub_date"`
}

func main() {
	// Initialize database
	dbPath := "/data/blog.db"
	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal(err)
	}

	server := &Server{db: db}

	// Load templates
	if err := server.loadTemplates(); err != nil {
		log.Fatal(err)
	}

	// Setup routes
	r := mux.NewRouter()

	// Static files with caching
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.HandlerFunc(server.staticHandler)))

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/news", server.newsHandler).Methods("GET")
	api.HandleFunc("/posts", server.postsHandler).Methods("GET")
	api.HandleFunc("/projects", server.projectsHandler).Methods("GET")

	// Page routes
	r.HandleFunc("/", server.homeHandler).Methods("GET")
	r.HandleFunc("/blog", server.blogHandler).Methods("GET")
	r.HandleFunc("/blog/{slug}", server.postHandler).Methods("GET")
	r.HandleFunc("/projects", server.projectsPageHandler).Methods("GET")
	r.HandleFunc("/about", server.aboutHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func (s *Server) loadTemplates() error {
	s.templates = make(map[string]*template.Template)

	// Define custom template functions
	funcMap := template.FuncMap{
		"split": strings.Split,
		"len": func(v interface{}) int {
			switch s := v.(type) {
			case string:
				return len(s)
			case []string:
				return len(s)
			default:
				return 0
			}
		},
		"wordCount": func(text string) int {
			if text == "" {
				return 0
			}
			words := strings.Fields(text)
			return len(words)
		},
		"truncate": func(text string, length int) string {
			if len(text) <= length {
				return text
			}
			return text[:length] + "..."
		},
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
		// Mathematical functions
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		// Comparison functions
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			return a != b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"le": func(a, b int) bool {
			return a <= b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"ge": func(a, b int) bool {
			return a >= b
		},
	}

	templateFiles := []string{
		"templates/layout.html",
		"templates/home.html",
		"templates/blog.html",
		"templates/post.html",
		"templates/projects.html",
		"templates/about.html",
	}

	for _, file := range templateFiles {
		name := filepath.Base(file)
		name = strings.TrimSuffix(name, filepath.Ext(name))

		// Create template with custom functions
		tmpl := template.New("layout.html").Funcs(funcMap)
		tmpl, err := tmpl.ParseFiles("templates/layout.html", file)
		if err != nil {
			return err
		}
		s.templates[name] = tmpl
	}

	return nil
}

func (s *Server) staticHandler(w http.ResponseWriter, r *http.Request) {
	// Set caching headers
	w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year

	// Determine content type
	ext := filepath.Ext(r.URL.Path)
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	}

	http.ServeFile(w, r, "static/"+r.URL.Path)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	// Get recent posts
	posts, _ := s.getRecentPosts(3)
	projects, _ := s.getProjects()

	data := struct {
		Posts    []Post    `json:"posts"`
		Projects []Project `json:"projects"`
		Title    string    `json:"title"`
	}{
		Posts:    posts,
		Projects: projects,
		Title:    "Jared Wallace - Software Developer",
	}

	s.renderTemplate(w, "home", data)
}

func (s *Server) blogHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	posts, _ := s.getPosts(page, 10)

	data := struct {
		Posts []Post `json:"posts"`
		Title string `json:"title"`
		Page  int    `json:"page"`
	}{
		Posts: posts,
		Title: "Blog - Jared Wallace",
		Page:  page,
	}

	s.renderTemplate(w, "blog", data)
}

func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	post, err := s.getPostBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Post  Post   `json:"post"`
		Title string `json:"title"`
	}{
		Post:  post,
		Title: post.Title + " - Jared Wallace",
	}

	s.renderTemplate(w, "post", data)
}

func (s *Server) projectsPageHandler(w http.ResponseWriter, r *http.Request) {
	projects, _ := s.getProjects()

	data := struct {
		Projects []Project `json:"projects"`
		Title    string    `json:"title"`
	}{
		Projects: projects,
		Title:    "Projects - Jared Wallace",
	}

	s.renderTemplate(w, "projects", data)
}

func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string `json:"title"`
	}{
		Title: "About - Jared Wallace",
	}

	s.renderTemplate(w, "about", data)
}

func (s *Server) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := s.templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// API Handlers
func (s *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	posts, _ := s.getRecentPosts(10)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (s *Server) projectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, _ := s.getProjects()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// Database methods
func (s *Server) getRecentPosts(limit int) ([]Post, error) {
	query := `SELECT id, title, content, excerpt, created_at, updated_at, published, slug 
			  FROM posts WHERE published = 1 ORDER BY created_at DESC LIMIT ?`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Excerpt,
			&post.CreatedAt, &post.UpdatedAt, &post.Published, &post.Slug)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Server) getPosts(page, limit int) ([]Post, error) {
	offset := (page - 1) * limit
	query := `SELECT id, title, content, excerpt, created_at, updated_at, published, slug 
			  FROM posts WHERE published = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Excerpt,
			&post.CreatedAt, &post.UpdatedAt, &post.Published, &post.Slug)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Server) getPostBySlug(slug string) (Post, error) {
	var post Post
	query := `SELECT id, title, content, excerpt, created_at, updated_at, published, slug 
			  FROM posts WHERE slug = ? AND published = 1`

	err := s.db.QueryRow(query, slug).Scan(&post.ID, &post.Title, &post.Content,
		&post.Excerpt, &post.CreatedAt, &post.UpdatedAt, &post.Published, &post.Slug)

	return post, err
}

func (s *Server) getProjects() ([]Project, error) {
	query := `SELECT id, title, description, url, image_url, technologies, created_at 
			  FROM projects ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.Title, &project.Description,
			&project.URL, &project.ImageURL, &project.Technologies, &project.CreatedAt)
		if err != nil {
			continue
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			excerpt TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			published BOOLEAN DEFAULT 0,
			slug TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			url TEXT,
			image_url TEXT,
			technologies TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}
	}

	return nil
}
