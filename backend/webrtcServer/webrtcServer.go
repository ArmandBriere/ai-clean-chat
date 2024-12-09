package webrtcserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/pion/webrtc/v4/pkg/media/oggwriter"
)

// Create the MediaEngine
var m = &webrtc.MediaEngine{}

// NewWebRTCServer creates a new WebRTC server and start listening
func NewWebRTCServer(messageChan chan string) {
	// Create a MediaEngine object to configure the supported codec

	// Register the VP8 and Opus codecs
	registerCodecs(m)

	// Register the default interceptor for audio
	localInterceptor := registerInterceptor(m)

	// Create a new OggWriter
	oggFile := getOggFile("output.ogg")

	// Setup the PeerConnection
	peerConnection := setupPeerConnection(localInterceptor, oggFile)

	// Listen for messages from the frontend
	listenForWebsocketMessage(messageChan, peerConnection)
}

// encode returns a base64 JSON of a SessionDescription
func encode(obj *webrtc.SessionDescription) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// decode a base64 and unmarshal JSON into a SessionDescription
func decode(in string, obj *webrtc.SessionDescription) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, obj); err != nil {
		panic(err)
	}
}

// getOggFile creates a new OggWriter
func getOggFile(fileName string) *oggwriter.OggWriter {
	oggFile, err := oggwriter.New(fileName, 48000, 2)
	if err != nil {
		panic(err)
	}
	return oggFile
}

// saveToDisk copies data from a webrtc.TrackRemote and writes it to disk.
func saveToDisk(i media.Writer, track *webrtc.TrackRemote) {
	defer func() {
		if err := i.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		rtpPacket, _, err := track.ReadRTP()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := i.WriteRTP(rtpPacket); err != nil {
			fmt.Println(err)
			return
		}
	}
}

// registerInterceptor registers a intervalpli interceptor with the MediaEngine.
func registerInterceptor(m *webrtc.MediaEngine) *interceptor.Registry {
	interceptor := &interceptor.Registry{}

	// Use the default set of Interceptors
	err := webrtc.RegisterDefaultInterceptors(m, interceptor)
	if err != nil {
		panic(err)
	}

	return interceptor
}

// Setup the PeerConnection
func setupPeerConnection(localInterceptor *interceptor.Registry, oggFile *oggwriter.OggWriter) *webrtc.PeerConnection {
	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(localInterceptor))

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Allow us to receive 1 audio track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	// Set a handler for when a new remote track starts, this handler saves buffers to disk as
	// an ivf file, since we could have multiple video tracks we provide a counter.
	// In your application this is where you would handle/process video
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) { //nolint: revive
		onTrack(track, oggFile)
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		onICEConnectionStateChange(connectionState, oggFile, peerConnection)
	})

	return peerConnection
}

// onICEConnectionStateChange is called when the ICE connection state changes
func onICEConnectionStateChange(connectionState webrtc.ICEConnectionState, oggFile *oggwriter.OggWriter, peerConnection *webrtc.PeerConnection) {
	fmt.Printf("Connection State has changed %s \n", connectionState.String())

	if connectionState == webrtc.ICEConnectionStateConnected {
		fmt.Println("Ctrl+C the remote client to stop the demo")
	} else if connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateClosed {
		if closeErr := oggFile.Close(); closeErr != nil {
			panic(closeErr)
		}
		fmt.Println("Done writing media files")

		if closeErr := peerConnection.Close(); closeErr != nil {
			panic(closeErr)
		}

		// TODO: Handle user disconnection that currently causes the server to exit
		os.Exit(0)
	}
}

// registerCodecs adds Opus codecs to the MediaEngine
func registerCodecs(m *webrtc.MediaEngine) {
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}
}

// onTrack is called when a new track is added
func onTrack(track *webrtc.TrackRemote, oggFile *oggwriter.OggWriter) {
	codec := track.Codec()
	if strings.EqualFold(codec.MimeType, webrtc.MimeTypeOpus) {
		fmt.Println("Got Opus track, saving to disk as output.opus (48 kHz, 2 channels)")
		saveToDisk(oggFile, track)
	}
}

// listenForWebsocketMessage listens for messages from the frontend and sends them to the peer connection
func listenForWebsocketMessage(messageChan chan string, peerConnection *webrtc.PeerConnection) {
	for {
		msg := <-messageChan
		log.Println("Received message:", msg)

		offer := webrtc.SessionDescription{}
		decode(msg, &offer)

		log.Println("Received offer from remote peer")
		log.Println(offer)
		err := peerConnection.SetRemoteDescription(offer)
		if err != nil {
			panic(err)
		}

		answer := generateICEAnswer(peerConnection)
		messageChan <- answer
	}
}

// generateICEAnswer generates an ICE answer
func generateICEAnswer(peerConnection *webrtc.PeerConnection) string {
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	<-gatherComplete

	return encode(peerConnection.LocalDescription())
}
