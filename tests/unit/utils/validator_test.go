package utils_test

import (
	"testing"

	"analytics-dashboard-api/internal/utils"
)

func TestValidatePositiveInt(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		want      int
		wantErr   bool
	}{
		{
			name:      "valid positive integer",
			value:     "123",
			fieldName: "test_field",
			want:      123,
			wantErr:   false,
		},
		{
			name:      "valid single digit",
			value:     "5",
			fieldName: "test_field",
			want:      5,
			wantErr:   false,
		},
		{
			name:      "valid with whitespace",
			value:     "  42  ",
			fieldName: "test_field",
			want:      42,
			wantErr:   false,
		},
		{
			name:      "empty string",
			value:     "",
			fieldName: "test_field",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "zero value",
			value:     "0",
			fieldName: "test_field",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "negative value",
			value:     "-5",
			fieldName: "test_field",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "non-numeric string",
			value:     "abc",
			fieldName: "test_field",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "float value",
			value:     "12.34",
			fieldName: "test_field",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "mixed alphanumeric",
			value:     "12abc",
			fieldName: "test_field",
			want:      0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ValidatePositiveInt(tt.value, tt.fieldName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePositiveInt() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidatePositiveInt() unexpected error: %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("ValidatePositiveInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateStringNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid non-empty string",
			value:     "hello",
			fieldName: "test_field",
			wantErr:   false,
		},
		{
			name:      "valid string with spaces",
			value:     "hello world",
			fieldName: "test_field",
			wantErr:   false,
		},
		{
			name:      "valid string with leading/trailing spaces",
			value:     "  hello  ",
			fieldName: "test_field",
			wantErr:   false,
		},
		{
			name:      "empty string",
			value:     "",
			fieldName: "test_field",
			wantErr:   true,
		},
		{
			name:      "whitespace only",
			value:     "   ",
			fieldName: "test_field",
			wantErr:   true,
		},
		{
			name:      "tab and newline only",
			value:     "\t\n",
			fieldName: "test_field",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateStringNotEmpty(tt.value, tt.fieldName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateStringNotEmpty() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateStringNotEmpty() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no trimming needed",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "leading spaces",
			input: "   hello",
			want:  "hello",
		},
		{
			name:  "trailing spaces",
			input: "hello   ",
			want:  "hello",
		},
		{
			name:  "leading and trailing spaces",
			input: "   hello   ",
			want:  "hello",
		},
		{
			name:  "tabs and newlines",
			input: "\t\nhello\t\n",
			want:  "hello",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "whitespace only",
			input: "   \t\n  ",
			want:  "",
		},
		{
			name:  "preserve internal spaces",
			input: "  hello world  ",
			want:  "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.SanitizeString(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeString() = %q, want %q", got, tt.want)
			}
		})
	}
}
