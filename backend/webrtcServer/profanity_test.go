package webrtcserver

import (
	"testing"
)

func init() {
	userSession := UserSession{}
	userSession.startNewSession("roomTest", "userTest")
}

// TestAppendToBuffer tests the appendToBuffer function.
// A buffer should never exceed 15 words. (space separated)
func TestAppendToBuffer(t *testing.T) {
	tests := []struct {
		initialBuffer  string
		appendText     string
		expectedBuffer string
	}{
		{"", "hello", "hello"},
		{"hello", " world", "hello world"},
		{"1234567890", "12345", "123456789012345"},
		{"a  a", "b  b", "a  ab  b"},
		{"abcdefghijklmnopqrstuvwxyz", "", "abcdefghijklmnopqrstuvwxyz"},
		{"abcdefghijklmnopqrstuvwxyz", "a", "abcdefghijklmnopqrstuvwxyza"},
		{"abcdefghijklmnopqrstuvwxyz", " a", "abcdefghijklmnopqrstuvwxyz a"},
		{"a b c d e f g h i j k l m n o p q r s t u v w x y z", " a", "m n o p q r s t u v w x y z a"},
	}
	for _, tt := range tests {
		t.Run(tt.appendText, func(t *testing.T) {
			userSession := UserSession{sentenceBuffer: tt.initialBuffer}
			userSession.appendToBuffer(tt.appendText)
			if userSession.getBuffer() != tt.expectedBuffer {
				t.Errorf("expected %s, got %s", tt.expectedBuffer, userSession.getBuffer())
			}
		})
	}
}

func TestKeepXWords(t *testing.T) {
	test := []struct {
		sentence string
		x        int
		expected string
	}{
		{"", 0, ""},
		{"", 1, ""},
		{"a", 0, ""},
		{"a", 1, "a"},
		{"a b c d e", 2, "d e"},
	}

	for _, tt := range test {
		t.Run(tt.sentence, func(t *testing.T) {
			if keepXWords(tt.sentence, tt.x) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, keepXWords(tt.sentence, tt.x))
			}
		})
	}
}
