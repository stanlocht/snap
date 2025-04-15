package web

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stanlocht/snap/pkg/issue"
	"github.com/stanlocht/snap/pkg/user"
)

// PageData represents the common data for all pages
type PageData struct {
	Title       string
	RepoName    string
	CurrentPage string
	Data        interface{}
}

// HomeData represents the data for the home page
type HomeData struct {
	CommitCount      int
	IssueCount       int
	ContributorCount int
	RecentCommits    []*CommitListItem
}

// CommitListItem represents a commit in the list
type CommitListItem struct {
	ID        string
	ShortID   string
	Message   string
	Author    string
	Timestamp string
	Emoji     string
}

// CommitDetailData represents the data for the commit detail page
type CommitDetailData struct {
	Commit *CommitListItem
	Files  []string
}

// IssueListItem represents an issue in the list
type IssueListItem struct {
	ID         int
	Title      string
	Status     string
	CreatedBy  string
	CreatedAt  string
	AssignedTo string
}

// IssueDetailData represents the data for the issue detail page
type IssueDetailData struct {
	Issue       *IssueListItem
	Description string
	Comments    []string
}

// UserListItem represents a user in the list
type UserListItem struct {
	Name    string
	Points  int
	Commits int
	Issues  int
}

// UserDetailData represents the data for the user detail page
type UserDetailData struct {
	User          *UserListItem
	RecentCommits []*CommitListItem
	RecentActions []string
}

// handleHome handles the home page
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get repository name
	repoName := filepath.Base(s.Repo.Path)

	// Get commit history
	history, err := s.Repo.GetCommitHistory("")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting commit history: %v", err), http.StatusInternalServerError)
		return
	}

	// Get recent commits
	recentCommits := make([]*CommitListItem, 0, 5)
	for i, commit := range history {
		if i >= 5 {
			break
		}

		emoji := extractEmoji(commit.Message)

		recentCommits = append(recentCommits, &CommitListItem{
			ID:        commit.ID,
			ShortID:   truncateID(commit.ID),
			Message:   strings.TrimPrefix(strings.TrimPrefix(commit.Message, emoji), " "),
			Author:    commit.Author,
			Timestamp: formatTime(commit.Timestamp),
			Emoji:     emoji,
		})
	}

	// Get issue count
	issueManager := issue.NewIssueManager(s.Repo.Path)
	issues, err := issueManager.ListIssues(true) // Show all issues, including closed ones
	issueCount := 0
	if err == nil {
		issueCount = len(issues)
	}

	// Get contributor count
	userManager := user.NewUserManager(s.Repo.Path)
	users, err := userManager.GetLeaderboard()
	contributorCount := 0
	if err == nil {
		contributorCount = len(users)
	}

	// Prepare data
	data := &PageData{
		Title:       "Home",
		RepoName:    repoName,
		CurrentPage: "home",
		Data: &HomeData{
			CommitCount:      len(history),
			IssueCount:       issueCount,
			ContributorCount: contributorCount,
			RecentCommits:    recentCommits,
		},
	}

	// Render template
	s.Templates.Execute(w, data)
}

// handleCommits handles the commits page
func (s *Server) handleCommits(w http.ResponseWriter, r *http.Request) {
	// Get repository name
	repoName := filepath.Base(s.Repo.Path)

	// Get commit history
	history, err := s.Repo.GetCommitHistory("")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting commit history: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare commit list
	commits := make([]*CommitListItem, 0, len(history))
	for _, commit := range history {
		emoji := extractEmoji(commit.Message)

		commits = append(commits, &CommitListItem{
			ID:        commit.ID,
			ShortID:   truncateID(commit.ID),
			Message:   strings.TrimPrefix(strings.TrimPrefix(commit.Message, emoji), " "),
			Author:    commit.Author,
			Timestamp: formatTime(commit.Timestamp),
			Emoji:     emoji,
		})
	}

	// Prepare data
	data := &PageData{
		Title:       "Commits",
		RepoName:    repoName,
		CurrentPage: "commits",
		Data:        commits,
	}

	// Render template
	s.Templates.Execute(w, data)
}

// handleCommitDetail handles the commit detail page
func (s *Server) handleCommitDetail(w http.ResponseWriter, r *http.Request) {
	// Get commit ID from URL
	commitID := strings.TrimPrefix(r.URL.Path, "/commit/")
	if commitID == "" {
		http.NotFound(w, r)
		return
	}

	// Get repository name
	repoName := filepath.Base(s.Repo.Path)

	// Get commit
	commit, err := s.Repo.GetCommit(commitID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting commit: %v", err), http.StatusInternalServerError)
		return
	}

	// Get tree
	tree, err := s.Repo.GetTree(commit.TreeID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting tree: %v", err), http.StatusInternalServerError)
		return
	}

	// Get files
	files := make([]string, 0, len(tree.Entries))
	for file := range tree.Entries {
		files = append(files, file)
	}

	// Extract emoji
	emoji := extractEmoji(commit.Message)

	// Prepare commit data
	commitData := &CommitListItem{
		ID:        commit.ID,
		ShortID:   truncateID(commit.ID),
		Message:   strings.TrimPrefix(strings.TrimPrefix(commit.Message, emoji), " "),
		Author:    commit.Author,
		Timestamp: formatTime(commit.Timestamp),
		Emoji:     emoji,
	}

	// Prepare data
	data := &PageData{
		Title:       "Commit Detail",
		RepoName:    repoName,
		CurrentPage: "commits",
		Data: &CommitDetailData{
			Commit: commitData,
			Files:  files,
		},
	}

	// Render template
	s.Templates.Execute(w, data)
}

// handleIssues handles the issues page
func (s *Server) handleIssues(w http.ResponseWriter, r *http.Request) {
	// Get repository name
	repoName := filepath.Base(s.Repo.Path)

	// Get issues
	issueManager := issue.NewIssueManager(s.Repo.Path)
	issues, err := issueManager.ListIssues(true) // Show all issues, including closed ones
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting issues: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare issue list
	issueList := make([]*IssueListItem, 0, len(issues))
	for _, issue := range issues {
		status := "Open"
		if string(issue.Status) == "closed" {
			status = "Closed"
		}

		issueList = append(issueList, &IssueListItem{
			ID:         issue.ID,
			Title:      issue.Title,
			Status:     status,
			CreatedBy:  issue.CreatedBy,
			CreatedAt:  formatTime(issue.CreatedAt),
			AssignedTo: issue.AssignedTo,
		})
	}

	// Prepare data
	data := &PageData{
		Title:       "Issues",
		RepoName:    repoName,
		CurrentPage: "issues",
		Data:        issueList,
	}

	// Render template
	s.Templates.Execute(w, data)
}

// handleIssueDetail handles the issue detail page
func (s *Server) handleIssueDetail(w http.ResponseWriter, r *http.Request) {
	// Get issue ID from URL
	issueIDStr := strings.TrimPrefix(r.URL.Path, "/issue/")
	if issueIDStr == "" {
		http.NotFound(w, r)
		return
	}

	// Parse issue ID
	issueID, err := strconv.Atoi(issueIDStr)
	if err != nil {
		http.Error(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	// Get repository name
	repoName := filepath.Base(s.Repo.Path)

	// Get issue
	issueManager := issue.NewIssueManager(s.Repo.Path)
	issue, err := issueManager.GetIssue(issueID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting issue: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare issue data
	status := "Open"
	if string(issue.Status) == "closed" {
		status = "Closed"
	}

	issueData := &IssueListItem{
		ID:         issue.ID,
		Title:      issue.Title,
		Status:     status,
		CreatedBy:  issue.CreatedBy,
		CreatedAt:  formatTime(issue.CreatedAt),
		AssignedTo: issue.AssignedTo,
	}

	// Prepare data
	data := &PageData{
		Title:       "Issue Detail",
		RepoName:    repoName,
		CurrentPage: "issues",
		Data: &IssueDetailData{
			Issue:       issueData,
			Description: issue.Description,
			Comments:    []string{}, // No comments implementation yet
		},
	}

	// Render template
	s.Templates.Execute(w, data)
}

// handleUsers handles the users page
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	// Get repository name
	repoName := filepath.Base(s.Repo.Path)

	// Get users
	userManager := user.NewUserManager(s.Repo.Path)
	users, err := userManager.GetLeaderboard()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting users: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare user list
	userList := make([]*UserListItem, 0, len(users))
	for _, u := range users {
		userList = append(userList, &UserListItem{
			Name:    u.Name,
			Points:  u.Points,
			Commits: u.Commits,
			Issues:  u.IssuesOpen + u.IssuesClosed,
		})
	}

	// Prepare data
	data := &PageData{
		Title:       "Users",
		RepoName:    repoName,
		CurrentPage: "users",
		Data:        userList,
	}

	// Render template
	s.Templates.Execute(w, data)
}

// handleUserDetail handles the user detail page
func (s *Server) handleUserDetail(w http.ResponseWriter, r *http.Request) {
	// Get user name from URL
	userName := strings.TrimPrefix(r.URL.Path, "/user/")
	if userName == "" {
		http.NotFound(w, r)
		return
	}

	// Get repository name
	repoName := filepath.Base(s.Repo.Path)

	// Get user
	userManager := user.NewUserManager(s.Repo.Path)
	u, err := userManager.GetUser(userName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting user: %v", err), http.StatusInternalServerError)
		return
	}

	// Format actions from user's action log
	formattedActions := make([]string, 0, len(u.ActionLog))
	for _, action := range u.ActionLog {
		formattedActions = append(formattedActions, fmt.Sprintf("%s: %s (+%d points)", action.Timestamp, action.Description, action.Points))
	}

	// Get commit history
	history, err := s.Repo.GetCommitHistory("")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting commit history: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user's recent commits
	recentCommits := make([]*CommitListItem, 0, 5)
	count := 0
	for _, commit := range history {
		if commit.Author == userName {
			emoji := extractEmoji(commit.Message)

			recentCommits = append(recentCommits, &CommitListItem{
				ID:        commit.ID,
				ShortID:   truncateID(commit.ID),
				Message:   strings.TrimPrefix(strings.TrimPrefix(commit.Message, emoji), " "),
				Author:    commit.Author,
				Timestamp: formatTime(commit.Timestamp),
				Emoji:     emoji,
			})

			count++
			if count >= 5 {
				break
			}
		}
	}

	// Prepare user data
	userData := &UserListItem{
		Name:    u.Name,
		Points:  u.Points,
		Commits: u.Commits,
		Issues:  u.IssuesOpen + u.IssuesClosed,
	}

	// Prepare data
	data := &PageData{
		Title:       "User Detail",
		RepoName:    repoName,
		CurrentPage: "users",
		Data: &UserDetailData{
			User:          userData,
			RecentCommits: recentCommits,
			RecentActions: formattedActions,
		},
	}

	// Render template
	s.Templates.Execute(w, data)
}

// extractEmoji extracts the emoji from a commit message
func extractEmoji(message string) string {
	// Check if message starts with an emoji (Unicode character)
	if len(message) > 0 {
		firstRune := []rune(message)[0]
		if firstRune > 127 { // Non-ASCII character
			// Find the end of the emoji
			for i, c := range message {
				if i > 0 && c == ' ' {
					return message[:i]
				}
			}
			// If no space found, return the first character
			return string([]rune(message)[0])
		}
	}

	// Check if message starts with an emoji code (e.g., :sparkles:)
	if strings.HasPrefix(message, ":") {
		endIndex := strings.Index(message[1:], ":")
		if endIndex != -1 {
			return message[:endIndex+2]
		}
	}

	return ""
}
