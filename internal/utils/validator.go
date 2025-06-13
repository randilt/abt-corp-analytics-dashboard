package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidatePositiveInt validates and converts string to positive integer
func ValidatePositiveInt(value string, fieldName string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("%s is required", fieldName)
	}

	num, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", fieldName)
	}

	if num <= 0 {
		return 0, fmt.Errorf("%s must be positive", fieldName)
	}

	return num, nil
}

// ValidateStringNotEmpty validates string is not empty
func ValidateStringNotEmpty(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// SanitizeString sanitizes string input
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}
