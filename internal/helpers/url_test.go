package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		ex    error
	}{
		{"valid URL", "https://google.com", nil},
		{"invalid string", "not-a-url", ErrInvalidURL},
		{"empty string", "", ErrInvalidURL},
		{"missing TLD", "http://localhost", ErrInvalidURL},
		{"missing scheme", "google.com", ErrInvalidURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.input)

			if tt.ex != nil {
				assert.ErrorIs(t, err, tt.ex)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
