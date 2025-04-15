package snapmoji

import (
	"errors"
	"fmt"
	"strings"
)

// Snapmoji represents a snapmoji with code and description
type Snapmoji struct {
	Emoji       string
	Code        string
	Description string
}

// List of supported snapmojis
var Snapmojis = []Snapmoji{
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

// ValidateCommitMessage checks if a commit message starts with a valid snapmoji
func ValidateCommitMessage(message string) error {
	if message == "" {
		return errors.New("commit message cannot be empty")
	}

	// Check if the message starts with an emoji or emoji code
	for _, snapmoji := range Snapmojis {
		if strings.HasPrefix(message, snapmoji.Emoji) || strings.HasPrefix(message, snapmoji.Code) {
			return nil
		}
	}

	return errors.New("commit message must start with a snapmoji (e.g., ✨ or :sparkles:)")
}

// GetSnapmojiList returns a formatted list of available snapmojis
func GetSnapmojiList() string {
	var builder strings.Builder
	builder.WriteString("Available Snapmojis:\n\n")

	for _, snapmoji := range Snapmojis {
		builder.WriteString(snapmoji.Emoji)
		builder.WriteString(" ")
		builder.WriteString(snapmoji.Code)
		builder.WriteString(" - ")
		builder.WriteString(snapmoji.Description)
		builder.WriteString("\n")
	}

	return builder.String()
}

// GetNumberedSnapmojiList returns a numbered list of available snapmojis for selection
func GetNumberedSnapmojiList() string {
	var builder strings.Builder
	builder.WriteString("Select a Snapmoji by number:\n\n")

	for i, snapmoji := range Snapmojis {
		builder.WriteString(fmt.Sprintf("%2d. %s %s - %s\n", i+1, snapmoji.Emoji, snapmoji.Code, snapmoji.Description))
	}

	return builder.String()
}

// GetSnapmojiByNumber returns a snapmoji by its number in the list (1-based)
func GetSnapmojiByNumber(number int) (Snapmoji, error) {
	if number < 1 || number > len(Snapmojis) {
		return Snapmoji{}, fmt.Errorf("invalid snapmoji number: %d (valid range: 1-%d)", number, len(Snapmojis))
	}

	return Snapmojis[number-1], nil
}

// AutoConvertKeywordsToEmoji converts keywords in a commit message to emojis
// It looks for keywords like "feature:", "fix:", "docs:" at the beginning of the message
func AutoConvertKeywordsToEmoji(message string) string {
	// Map of keywords to snapmoji codes
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
