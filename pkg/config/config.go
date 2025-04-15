package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetValue gets a configuration value from the config file
func GetValue(configPath, key string) (string, error) {
	// Read config file
	file, err := os.Open(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Parse key into section and name
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid key format, expected 'section.name'")
	}
	section := parts[0]
	name := parts[1]

	// Read config file line by line
	scanner := bufio.NewScanner(file)
	inSection := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check if we're entering a section
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := line[1 : len(line)-1]
			inSection = (sectionName == section)
			continue
		}

		// If we're in the right section, check for the key
		if inSection && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if strings.TrimSpace(parts[0]) == name {
				return strings.TrimSpace(parts[1]), nil
			}
		}

		// If we're in the right section and it's a tab-indented line
		if inSection && strings.HasPrefix(line, "\t") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[0]) == name {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading config file: %w", err)
	}

	return "", nil
}

// SetValue sets a configuration value in the config file
func SetValue(configPath, key, value string) error {
	// Parse key into section and name
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format, expected 'section.name'")
	}
	section := parts[0]
	name := parts[1]

	// Read config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Split content into lines
	lines := strings.Split(string(content), "\n")

	// Find section and update value
	inSection := false
	updated := false
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check if we're entering a section
		if strings.HasPrefix(trimmedLine, "[") && strings.HasSuffix(trimmedLine, "]") {
			sectionName := trimmedLine[1 : len(trimmedLine)-1]
			inSection = (sectionName == section)
			continue
		}

		// If we're in the right section, check for the key
		if inSection && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if strings.TrimSpace(parts[0]) == name {
				lines[i] = fmt.Sprintf("\t%s = %s", name, value)
				updated = true
				break
			}
		}

		// If we're in the right section and it's a tab-indented line
		if inSection && strings.HasPrefix(line, "\t") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[0]) == name {
				lines[i] = fmt.Sprintf("\t%s = %s", name, value)
				updated = true
				break
			}
		}
	}

	// If we didn't update an existing value, add it to the section
	if !updated {
		// Find the section again
		newLines := []string{}
		inSection = false
		for _, line := range lines {
			newLines = append(newLines, line)
			trimmedLine := strings.TrimSpace(line)

			// Check if we're entering the section
			if strings.HasPrefix(trimmedLine, "[") && strings.HasSuffix(trimmedLine, "]") {
				sectionName := trimmedLine[1 : len(trimmedLine)-1]
				if sectionName == section {
					inSection = true
					// Add the new value after the section header
					newLines = append(newLines, fmt.Sprintf("\t%s = %s", name, value))
				}
			}
		}
		lines = newLines
	}

	// Write updated content back to file
	err = os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
