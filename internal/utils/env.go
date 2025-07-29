package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseEnvFile parses a .env file and returns a map of environment variables
func ParseEnvFile(filePath string) (map[string]string, error) {
	envVars := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open env file '%s': %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format at line %d: %s (expected KEY=VALUE)", lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
				(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
				value = value[1 : len(value)-1]
			}
		}

		if key == "" {
			return nil, fmt.Errorf("empty key at line %d", lineNumber)
		}

		envVars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading env file: %w", err)
	}

	return envVars, nil
}

// ValidateResourceValue validates CPU and memory resource values
func ValidateResourceValue(value, resourceType string) error {
	if value == "" {
		return nil // Optional values
	}

	// Basic validation for CPU and memory formats
	switch resourceType {
	case "cpu":
		if !strings.HasSuffix(value, "m") && !isNumeric(value) {
			return fmt.Errorf("invalid CPU format '%s' (examples: 100m, 0.5, 1)", value)
		}
	case "memory":
		if !strings.HasSuffix(value, "Mi") && !strings.HasSuffix(value, "Gi") && !strings.HasSuffix(value, "Ki") {
			return fmt.Errorf("invalid memory format '%s' (examples: 128Mi, 1Gi, 512Ki)", value)
		}
	}

	return nil
}

// isNumeric checks if a string represents a number
func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			if char != '.' {
				return false
			}
		}
	}
	return true
}
