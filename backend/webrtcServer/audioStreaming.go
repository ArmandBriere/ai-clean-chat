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
func getFileAndWriter(writer *oggwriter.OggWriter) (*os.File, *oggwriter.OggWriter, bool) {
	filename := fmt.Sprintf("audio_%d.ogg", time.Now().UnixNano())
	f, err := os.Create(filename)
	if err != nil {
		slog.Error("Error creating file", "Error", err)
		return nil, nil, true
	}
	writer, err = oggwriter.New(filename, 48000, 2)
	if err != nil {
		slog.Error("Error creating oggwriter", "Error", err)
		return nil, nil, true
	}

	return f, writer, false
}

// handleAudioStream handles the audio stream by writing it to file
func handleAudioStream(track *webrtc.TrackRemote, isStreaming *bool) {
	func() {
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
					var shouldReturn bool
					var f *os.File
					f, writer, shouldReturn = getFileAndWriter(writer)
					if shouldReturn {
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
	}()
}
