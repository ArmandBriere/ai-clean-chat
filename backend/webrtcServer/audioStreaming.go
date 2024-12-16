package webrtcserver

import (
	"fmt"
	"log/slog"
	"os"
	"time"

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
func handleAudioStream(track *webrtc.TrackRemote, isStreaming *bool) {

	// WIP: Use transcribe that is defined in sherpa.go
	// This take the audio stream for ever
	Transcribe(track)

	var writer *oggwriter.OggWriter
	for {
		// Read RTP packets to flush the buffer
		rtpPacket, _, err := track.ReadRTP()
		if err != nil {
			slog.Error("rtpPacket setup", "Error", err)
			return
		}
		if *isStreaming {
			if writer == nil {
				var f *os.File
				f, writer, err = getFileAndWriter(writer)
				if err != nil {
					slog.Error("Error creating file and oggwriter", "Error", err)
					return
				}
				defer func() {
					writer.Close()
					f.Close()
				}()
			}

			if err := writer.WriteRTP(rtpPacket); err != nil {
				slog.Error("write RTP", "Error", err)
				return
			}
		}
	}
}
