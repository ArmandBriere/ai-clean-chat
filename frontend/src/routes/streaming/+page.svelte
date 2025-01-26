<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { PUBLIC_SERVER_WS_URL } from '$env/static/public';
  import { v4 } from 'uuid';

  let roomID: string = 'example-room'; // Will be replaced with dynamic room ID
  let userVideo: HTMLVideoElement;
  let otherVideo: HTMLVideoElement;

  let ws: WebSocket;
  let peerConnection: RTCPeerConnection;
  let localStream: MediaStream;

  const ICE_CANDIDATE = 'iceCandidate';
  const OFFER = 'offer';
  const ANSWER = 'answer';
  const STREAMING = 'streaming';
  const TRANSCRIPTION = 'transcription';

  onMount(() => {
    // Connect to the signaling server
    ws = new WebSocket(`${PUBLIC_SERVER_WS_URL}/join?roomID=${roomID}&userID=${v4()}`);

    ws.onopen = () => {
      console.log('WebSocket connection established');
    };

    ws.onmessage = async (event) => {
      const message = JSON.parse(event.data);
      console.log('Received signaling message:', message);

      if (message.type === OFFER) {
        await handleOffer(message.candidate);
      } else if (message.type == ANSWER) {
        console.log('Received answer:', message);
        await handleAnswer(message.sdp);
      } else if (message.type == ICE_CANDIDATE) {
        console.log('Received ICE candidate:', message.candidate);
        await handleIceCandidate(message.candidate);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    // Open the camera and set up local video
    openCamera();

    return () => {
      onDestroy(() => {
        ws?.close();
        peerConnection?.close();
        localStream?.getTracks().forEach((track) => track.stop());
      });
    };
  });

  async function openCamera() {
    try {
      localStream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: true
      });
      userVideo.srcObject = localStream;

      // Create the PeerConnection
      createPeerConnection();

      // Add local stream tracks to the PeerConnection
      localStream.getTracks().forEach((track) => {
        peerConnection.addTrack(track, localStream);
      });
    } catch (error) {
      console.error('Error accessing media devices:', error);
    }
  }

  function createPeerConnection() {
    peerConnection = new RTCPeerConnection();

    // Handle ICE candidates
    peerConnection.onicecandidate = (event) => {
      if (event.type === 'icecandidate' && event.candidate) {
        const message: WebSocketMessage = {
          type: ICE_CANDIDATE,
          candidate: event.candidate
        };
        ws.send(JSON.stringify(message));
      }
    };

    // Handle incoming tracks
    peerConnection.ontrack = (event) => {
      console.log('Received remote track:', event.streams);
      if (event.streams.length > 0) {
        otherVideo.srcObject = event.streams[0];
      }
    };

    // Handle negotiation
    peerConnection.onnegotiationneeded = async () => {
      const offer = await peerConnection.createOffer();
      await peerConnection.setLocalDescription(offer);
      if (offer.type === OFFER && offer) {
        const offerMessage: WebSocketMessage = {
          type: OFFER,
          candidate: offer
        };
        ws.send(JSON.stringify(offerMessage));
      }
    };
  }

  async function handleOffer(offer: RTCSessionDescriptionInit) {
    await peerConnection.setRemoteDescription(offer);
    const answer = await peerConnection.createAnswer();
    await peerConnection.setLocalDescription(answer);

    if (answer.type === ANSWER && answer.sdp) {
      const answerMessage: any = {
        type: ANSWER,
        sdp: answer
      };
      ws.send(JSON.stringify(answerMessage));
    }
  }

  async function handleAnswer(answer: RTCSessionDescriptionInit) {
    try {
      await peerConnection.setRemoteDescription(answer);
    } catch (error) {
      console.error('Error handling answer:', error);
    }
  }

  async function handleIceCandidate(candidate: any) {
    try {
      await peerConnection.addIceCandidate(candidate);
    } catch (error) {
      console.error('Error adding received ICE candidate:', error);
    }
  }
</script>

<main class="container mx-auto p-4">
  <h1 class="mb-4 text-xl font-bold">Svelte WebRTC Video Chat</h1>

  <div class="grid grid-cols-2 gap-4">
    <div>
      <h2 class="text-lg font-semibold">Your Video</h2>
      <video
        playsInline
        autoPlay
        muted
        class="w-full rounded-lg border shadow"
        bind:this={userVideo}
      ></video>
    </div>

    <div>
      <h2 class="text-lg font-semibold">Other User's Video</h2>
      <video
        playsInline
        autoPlay
        muted
        class="w-full rounded-lg border shadow"
        bind:this={otherVideo}
      ></video>
    </div>
  </div>
</main>

<style>
  main {
    max-width: 800px;
  }
</style>
