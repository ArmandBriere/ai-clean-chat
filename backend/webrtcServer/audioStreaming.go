package webrtcserver

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/oggwriter"
)

// getFileAndWriter creates a new file and oggwriter
func getFileAndWriter(writer *oggwriter.OggWriter) (*os.File, *oggwriter.OggWriter, error) {
	filename := fmt.Sprintf("audio_%d.ogg", time.Now().UnixNano())
	f, err := os.Create(filename)
	if err != nil {
		slog.Error("Error creating file", "Error", err)
		return nil, nil, err
	}
	writer, err = oggwriter.New(filename, 48000, 1)
	if err != nil {
		slog.Error("Error creating oggwriter", "Error", err)
		return nil, nil, err
	}

	return f, writer, nil
}

// handleAudioStream handles the audio stream by writing it to file
func handleAudioStream(track *webrtc.TrackRemote, isStreaming *bool, quit chan bool, wsConn *websocket.Conn) {

	// This take the audio stream for ever
	transcribe(track, isStreaming, quit, wsConn)
}
