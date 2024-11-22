package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{"int pointer", 42, 42},
		{"string pointer", "hello", "hello"},
		{"float64 pointer", 3.14, 3.14},
		{"bool pointer", true, true},
		{"struct pointer", struct{ A int }{A: 1}, struct{ A int }{A: 1}},

		// Default values
		{"zero int pointer", 0, 0},
		{"empty string pointer", "", ""},
		{"false bool pointer", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := Ptr(tt.input)

			assert.NotNil(t, ptr)

			assert.Equal(t, tt.expected, *ptr)
		})
	}
}
