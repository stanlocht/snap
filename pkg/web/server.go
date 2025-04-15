package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/stanlocht/snap/pkg/repository"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

// Server represents the web server
type Server struct {
	Repo      *repository.Repository
	Templates *template.Template
}

// NewServer creates a new web server
func NewServer(repo *repository.Repository) (*Server, error) {
	// Parse templates
	templates, err := template.ParseFS(templateFS, "templates/combined.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	return &Server{
		Repo:      repo,
		Templates: templates,
	}, nil
}

// StartServer starts the web server
func StartServer(repo *repository.Repository, port int) error {
	// Create server
	server, err := NewServer(repo)
	if err != nil {
		return err
	}

	// Set up static file server
	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		return fmt.Errorf("error setting up static file server: %w", err)
	}

	// Register handlers
	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/commits", server.handleCommits)
	http.HandleFunc("/commit/", server.handleCommitDetail)
	http.HandleFunc("/issues", server.handleIssues)
	http.HandleFunc("/issue/", server.handleIssueDetail)
	http.HandleFunc("/users", server.handleUsers)
	http.HandleFunc("/user/", server.handleUserDetail)
	http.HandleFunc("/quest", server.handleQuest)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		data, err := templateFS.ReadFile("templates/test.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading test template: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))

	// Start server
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}

// formatTime formats a time.Time for display
func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// truncateID truncates a commit ID
func truncateID(id string) string {
	if len(id) > 7 {
		return id[:7]
	}
	return id
}
