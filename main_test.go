package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard string",
			input:    "Man",
			expected: "TWFu",
		},
		{
			name:     "1 byte padding",
			input:    "any car",
			expected: "YW55IGNhcg==",
		},
		{
			name:     "2 byte padding",
			input:    "any car.",
			expected: "YW55IGNhci4=",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "emoji multi-byte UTF-8",
			input:    "😊",
			expected: "8J+Yig==",
		},
	}

	for _, tt := range tests {
		t.Run("encodeToBase64 - "+tt.name, func(t *testing.T) {
			inputReader := strings.NewReader(tt.input)
			var outputBuffer bytes.Buffer

			encodeToBase64(inputReader, &outputBuffer)
			got := outputBuffer.String()
			if got != tt.expected {
				t.Errorf("Encode(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})

		t.Run("decodeFromBase64 - "+tt.name, func(t *testing.T) {
			inputReader := strings.NewReader(tt.expected)
			var outputBuffer bytes.Buffer

			decodeFromBase64(inputReader, &outputBuffer)
			got := fmt.Sprint(string(outputBuffer.Bytes()))
			if got != tt.input {
				t.Errorf("Decode(%q) = %q, want %q", tt.expected, got, tt.input)
			}
		})
	}
}
