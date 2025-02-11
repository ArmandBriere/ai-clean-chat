package webrtcserver

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

type UserSession struct {
	RoomID         string
	UserID         string
	sentenceBuffer string
	previousScore  float64
}

// startNewSession starts a new session with the given roomID and userID
func (s *UserSession) startNewSession(roomID string, userID string) {
	s.sentenceBuffer = ""
	s.RoomID = roomID
	s.UserID = userID
	s.previousScore = 0
}

// appendToBuffer appends the sentence to the sentence buffer
func (s *UserSession) appendToBuffer(sentence string) {
	s.sentenceBuffer += sentence

	s.sentenceBuffer = keepXWords(s.sentenceBuffer, 8)
}

// keepXWords returns the last x words of the sentence
func keepXWords(s string, x int) string {
	words := strings.Fields(s)
	if len(words) <= x {
		return s
	}
	return strings.Join(words[len(words)-x:], " ")
}

// getBuffer returns the sentence buffer
func (s *UserSession) getBuffer() string {
	return s.sentenceBuffer
}

// getBuffer returns the sentence buffer
func (s *UserSession) getBufferLen() int {
	return len(s.sentenceBuffer)
}

// clearBuffer clears the sentence buffer
func (s *UserSession) clearBuffer() {
	s.sentenceBuffer = ""
}

// analyzeBuffer sends the sentence buffer to the profanity API and returns the profanity score
func (s *UserSession) analyzeBuffer() (float64, error) {

	url := "http://profanity:8080/profanity"
	data := PostData{
		Text: s.sentenceBuffer,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("Error marshaling JSON", "err", err)
		return 0, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error sending post request", "err", err)
		return 0, err
	}
	defer resp.Body.Close()

	var responseData PostResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		slog.Error("Error decoding response", "err", err)
		return 0, err
	}

	slog.Info("Profanity score", "score", responseData.ProfanityScore, "UserID", s.UserID, "RoomID", s.RoomID)
	return responseData.ProfanityScore, nil
}
