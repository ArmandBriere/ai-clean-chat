//go:build !profanity

package webrtcserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hraban/opus"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"github.com/pion/webrtc/v4"
)

var stream *sherpa.OnlineStream
var recognizer *sherpa.OnlineRecognizer

// init initializes the recognizer and stream
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	config := sherpa.OnlineRecognizerConfig{}
	config.FeatConfig = sherpa.FeatureConfig{SampleRate: MODEL_SAMPLE_RATE, FeatureDim: 80}

	config.ModelConfig.NumThreads = 8
	config.ModelConfig.Debug = 0
	config.ModelConfig.Provider = "cpu"

	defaultPath := "sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/"
	config.ModelConfig.Transducer.Encoder = "./" + defaultPath + "encoder-epoch-99-avg-1.onnx"
	config.ModelConfig.Transducer.Decoder = "./" + defaultPath + "decoder-epoch-99-avg-1.onnx"
	config.ModelConfig.Transducer.Joiner = "./" + defaultPath + "joiner-epoch-99-avg-1.onnx"
	config.ModelConfig.Tokens = "./" + defaultPath + "tokens.txt"

	slog.Info("Initializing recognizer (may take several seconds)")
	recognizer = sherpa.NewOnlineRecognizer(&config)
	slog.Info("Recognizer created!")

	stream = sherpa.NewOnlineStream(recognizer)
	slog.Info("Stream created!")
}

// Transcribe transcribes the audio stream
func transcribe(ctx context.Context, track *webrtc.TrackRemote, isStreaming *bool, wsConn *websocket.Conn, mu *sync.Mutex) {

	var last_text string
	segment_idx := 0

	// Create an Opus decoder
	decoder, err := opus.NewDecoder(INPUT_SAMPLE_RATE, 1) // Mono channel
	if err != nil {
		slog.Error("Failed to create Opus decoder", "error", err)
	}

	userSession := UserSession{}
	userSession.startNewSession("roomID", "userID")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Transcription stopped by the context")
			return
		default:
			rtpPacket, _, err := track.ReadRTP()
			if err != nil {
				slog.Error("Failed to read RTP packet", "packet", err)
				recognizer.Reset(stream)
				continue
			}

			// Skip if user is not streaming
			if !*isStreaming {
				continue
			}

			payload := rtpPacket.Payload
			if len(payload) == 0 {
				continue
			}

			// Decode RTP payload into PCM samples (depends on your audio codec)
			pcmSamples, err := decodeRTPPayload(decoder, payload)
			if err != nil {
				pcmSamples = []int16{0}
			}

			// Convert PCM samples ([]int16) to []float32
			samples := PcmToFloat32(pcmSamples)

			// Process samples
			stream.AcceptWaveform(int(INPUT_SAMPLE_RATE), samples)

			for recognizer.IsReady(stream) {
				recognizer.Decode(stream)
			}

			text := recognizer.GetResult(stream).Text
			if len(text) != 0 && last_text != text {
				last_text = strings.ToLower(text)
				slog.Info("Transcription", "text", last_text)
				userSession.appendToBuffer(text)
				profanityScore, err := userSession.analyzeBuffer()
				if err != nil {
					slog.Error("Error analyzing buffer", "error", err)
				}

				slog.Info("Profanity score", "score", profanityScore)
				uuid := uuid.New().String()
				mu.Lock()
				wsConn.WriteJSON(WebSocketTranscription{
					Type:           "transcription",
					Text:           last_text,
					Uuid:           uuid,
					ProfanityScore: profanityScore,
				})
				mu.Unlock()

				recognizer.Reset(stream)
			}

			if recognizer.IsEndpoint(stream) {
				if len(text) != 0 {
					segment_idx++
					fmt.Println()
					recognizer.Reset(stream)
				}
			}
		}
	}
}

// decodeRTPPayload decodes the RTP payload into PCM samples
func decodeRTPPayload(decoder *opus.Decoder, payload []byte) ([]int16, error) {
	// Allocate space for PCM samples
	// Max Opus frame size is 120ms: 48,000 Hz * 0.06 seconds = 5760 samples
	pcm := make([]int16, 5760)

	// Decode the Opus payload into PCM
	n, err := decoder.Decode(payload, pcm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Opus payload: %v", err)
	}

	// Return the decoded PCM samples
	return pcm[:n], nil
}
