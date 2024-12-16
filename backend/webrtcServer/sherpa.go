package webrtcserver

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/hraban/opus"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"github.com/pion/webrtc/v4"
)

var stream *sherpa.OnlineStream
var recognizer *sherpa.OnlineRecognizer

const (
	sampleRate = 48000
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	config := sherpa.OnlineRecognizerConfig{}
	config.FeatConfig = sherpa.FeatureConfig{SampleRate: MODEL_SAMPLE_RATE, FeatureDim: 80}

	// Configure the model paths
	// config.ModelConfig.Zipformer2Ctc.Model = "./sherpa-onnx-streaming-zipformer-ctc-small-2024-03-18/ctc-epoch-30-avg-3-chunk-16-left-128.int8.onnx"
	// config.ModelConfig.Tokens = "./sherpa-onnx-streaming-zipformer-ctc-small-2024-03-18/tokens.txt"
	// config.CtcFstDecoderConfig.Graph = "./sherpa-onnx-streaming-zipformer-ctc-small-2024-03-18/HLG.fst"

	config.ModelConfig.NumThreads = 8
	config.ModelConfig.Debug = 0
	config.ModelConfig.Provider = "cpu"

	defaultPath := "sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/"
	config.ModelConfig.Transducer.Encoder = "./" + defaultPath + "encoder-epoch-99-avg-1.onnx"
	config.ModelConfig.Transducer.Decoder = "./" + defaultPath + "decoder-epoch-99-avg-1.onnx"
	config.ModelConfig.Transducer.Joiner = "./" + defaultPath + "joiner-epoch-99-avg-1.onnx"
	config.ModelConfig.Tokens = "./" + defaultPath + "tokens.txt"
	// config.ModelConfig.ModelType = "zipformer2"

	slog.Info("Initializing recognizer (may take several seconds)")
	recognizer = sherpa.NewOnlineRecognizer(&config)
	slog.Info("Recognizer created!")

	stream = sherpa.NewOnlineStream(recognizer)
	slog.Info("Stream created!")
}

func Transcribe(track *webrtc.TrackRemote) {

	var last_text string
	segment_idx := 0

	// Create an Opus decoder
	decoder, err := opus.NewDecoder(48000, 1) // Mono channel
	if err != nil {
		slog.Error("Failed to create Opus decoder", "error", err)
	}

	for {
		rtpPacket, _, err := track.ReadRTP()
		if err != nil {
			log.Fatalf("Failed to read RTP packet: %v", err)
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
		samples := pcmToFloat32(pcmSamples)

		// Process samples
		stream.AcceptWaveform(int(INPUT_SAMPLE_RATE), samples)

		for recognizer.IsReady(stream) {
			recognizer.Decode(stream)
		}

		text := recognizer.GetResult(stream).Text
		if len(text) != 0 && last_text != text {
			last_text = strings.ToLower(text)
			fmt.Printf("\r%d: %s", segment_idx, last_text)
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

func pcmToFloat32(pcm []int16) []float32 {
	floatSamples := make([]float32, len(pcm))
	for i, sample := range pcm {
		// Normalize from int16 (-32768 to 32767) to float32 (-1.0 to 1.0)
		floatSamples[i] = float32(sample) / 32768.0
	}
	return floatSamples
}

func intToFloat32(arr []int) []float32 {
	result := make([]float32, len(arr))
	for i, v := range arr {
		result[i] = float32(v)
	}
	return result
}
