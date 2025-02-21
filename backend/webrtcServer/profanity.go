package webrtcserver

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openai/openai-go"
)

type UserSession struct {
	RoomID          string
	UserID          string
	sentenceBuffer  string
	client          *openai.Client
	bufferCounter   int
	timeToProfanity int64
	tokenCounter    int
}

// startNewSession starts a new session with the given roomID and userID
func (s *UserSession) startNewSession(roomID string, userID string) {
	s.sentenceBuffer = ""
	s.RoomID = roomID
	s.UserID = userID
	s.client = openai.NewClient()
	s.bufferCounter = 0
	s.timeToProfanity = 0.0
	s.tokenCounter = 0
}

// appendToBuffer appends the sentence to the sentence buffer
func (s *UserSession) appendToBuffer(sentence string) {
	s.sentenceBuffer += sentence

	if s.bufferCounter >= PROFANITY_ANALYSIS_BUFFER_SIZE {
		s.bufferCounter = 0
	}

	s.bufferCounter++
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
func (s *UserSession) analyzeBuffer(wsConn *websocket.Conn, mu *sync.Mutex) (float64, error) {

	startTime := time.Now()
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

	if responseData.ProfanityScore > 0.9 {
		go s.llmAnalysis(wsConn, mu)
	}

	endTime := time.Now()
	s.tokenCounter += 1
	s.timeToProfanity = s.timeToProfanity + endTime.Sub(startTime).Milliseconds()
	slog.Info("Profanity analysis", "profanityScore", responseData.ProfanityScore, "svgTime", s.timeToProfanity/int64(s.tokenCounter))
	return responseData.ProfanityScore, nil
}

// llmAnalysis sends the sentence buffer to the LLM API and returns the analysis
func (s *UserSession) llmAnalysis(wsConn *websocket.Conn, mu *sync.Mutex) error {

	// Only analyze every PROFANITY_ANALYSIS_BUFFER_SIZE tokens
	if s.bufferCounter < PROFANITY_ANALYSIS_BUFFER_SIZE {
		return nil
	}

	userBuffer := s.sentenceBuffer

	ctx := context.Background()
	userMessage := openai.UserMessage(userBuffer)
	systemMessage := openai.SystemMessage(LLM_PROMPT)

	completion, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			systemMessage,
			userMessage,
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})

	if err != nil {
		slog.Error("Error creating completion", "err", err)
		return err
	}
	slog.Info("LLM answer", "content", completion.Choices[0].Message.Content)

	mu.Lock()
	defer mu.Unlock()
	data := LLMAnalysis{
		Type:        "llmAnalysis",
		LLMMessage:  completion.Choices[0].Message.Content,
		UserMessage: userBuffer,
	}
	if err := wsConn.WriteJSON(data); err != nil {
		slog.Error("Error writing LLM analysis", "error", err)
		return err
	}
	return nil
}
