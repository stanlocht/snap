package gitmoji

import (
	"errors"
	"fmt"
	"strings"
)

// Gitmoji represents a gitmoji with code and description
type Gitmoji struct {
	Emoji       string
	Code        string
	Description string
}

// List of supported gitmojis
var Gitmojis = []Gitmoji{
	{Emoji: "✨", Code: ":sparkles:", Description: "Introduce new features"},
	{Emoji: "🐛", Code: ":bug:", Description: "Fix a bug"},
	{Emoji: "📚", Code: ":books:", Description: "Add or update documentation"},
	{Emoji: "♻️", Code: ":recycle:", Description: "Refactor code"},
	{Emoji: "🔧", Code: ":wrench:", Description: "Add or update configuration files"},
	{Emoji: "✅", Code: ":white_check_mark:", Description: "Add, update, or pass tests"},
	{Emoji: "🚀", Code: ":rocket:", Description: "Deploy stuff"},
	{Emoji: "💄", Code: ":lipstick:", Description: "Add or update the UI and style files"},
	{Emoji: "🔥", Code: ":fire:", Description: "Remove code or files"},
	{Emoji: "🚑", Code: ":ambulance:", Description: "Critical hotfix"},
	{Emoji: "🎨", Code: ":art:", Description: "Improve structure / format of the code"},
	{Emoji: "⚡️", Code: ":zap:", Description: "Improve performance"},
	{Emoji: "🔒", Code: ":lock:", Description: "Fix security issues"},
	{Emoji: "🚧", Code: ":construction:", Description: "Work in progress"},
	{Emoji: "📝", Code: ":memo:", Description: "Add or update documentation"},
	{Emoji: "🚚", Code: ":truck:", Description: "Move or rename resources"},
	{Emoji: "👷", Code: ":construction_worker:", Description: "Add or update CI build system"},
	{Emoji: "➕", Code: ":heavy_plus_sign:", Description: "Add a dependency"},
	{Emoji: "➖", Code: ":heavy_minus_sign:", Description: "Remove a dependency"},
	{Emoji: "🔖", Code: ":bookmark:", Description: "Release / Version tags"},
}

// ValidateCommitMessage checks if a commit message starts with a valid gitmoji
func ValidateCommitMessage(message string) error {
	if message == "" {
		return errors.New("commit message cannot be empty")
	}

	// Check if the message starts with an emoji or emoji code
	for _, gitmoji := range Gitmojis {
		if strings.HasPrefix(message, gitmoji.Emoji) || strings.HasPrefix(message, gitmoji.Code) {
			return nil
		}
	}

	return errors.New("commit message must start with a gitmoji (e.g., ✨ or :sparkles:)")
}

// GetGitmojiList returns a formatted list of available gitmojis
func GetGitmojiList() string {
	var builder strings.Builder
	builder.WriteString("Available Gitmojis:\n\n")

	for _, gitmoji := range Gitmojis {
		builder.WriteString(gitmoji.Emoji)
		builder.WriteString(" ")
		builder.WriteString(gitmoji.Code)
		builder.WriteString(" - ")
		builder.WriteString(gitmoji.Description)
		builder.WriteString("\n")
	}

	return builder.String()
}

// GetNumberedGitmojiList returns a numbered list of available gitmojis for selection
func GetNumberedGitmojiList() string {
	var builder strings.Builder
	builder.WriteString("Select a Gitmoji by number:\n\n")

	for i, gitmoji := range Gitmojis {
		builder.WriteString(fmt.Sprintf("%2d. %s %s - %s\n", i+1, gitmoji.Emoji, gitmoji.Code, gitmoji.Description))
	}

	return builder.String()
}

// GetGitmojiByNumber returns a gitmoji by its number in the list (1-based)
func GetGitmojiByNumber(number int) (Gitmoji, error) {
	if number < 1 || number > len(Gitmojis) {
		return Gitmoji{}, fmt.Errorf("invalid gitmoji number: %d (valid range: 1-%d)", number, len(Gitmojis))
	}

	return Gitmojis[number-1], nil
}

// AutoConvertKeywordsToEmoji converts keywords in a commit message to emojis
// It looks for keywords like "feature:", "fix:", "docs:" at the beginning of the message
func AutoConvertKeywordsToEmoji(message string) string {
	// Map of keywords to gitmoji codes
	keywordMap := map[string]string{
		"feature:":  ":sparkles:",
		"feat:":     ":sparkles:",
		"fix:":      ":bug:",
		"docs:":     ":books:",
		"doc:":      ":books:",
		"refactor:": ":recycle:",
		"perf:":     ":zap:",
		"test:":     ":white_check_mark:",
		"chore:":    ":wrench:",
		"style:":    ":lipstick:",
		"remove:":   ":fire:",
		"delete:":   ":fire:",
		"hotfix:":   ":ambulance:",
		"ui:":       ":lipstick:",
		"security:": ":lock:",
		"wip:":      ":construction:",
		"deploy:":   ":rocket:",
		"release:":  ":bookmark:",
		"add:":      ":heavy_plus_sign:",
		"subtract:": ":heavy_minus_sign:",
		"move:":     ":truck:",
		"rename:":   ":truck:",
		"ci:":       ":construction_worker:",
	}

	// Check if message starts with any of the keywords
	for keyword, emojiCode := range keywordMap {
		if strings.HasPrefix(strings.ToLower(message), strings.ToLower(keyword)) {
			// Replace the keyword with the emoji code
			return emojiCode + " " + strings.TrimSpace(message[len(keyword):])
		}
	}

	return message
}
