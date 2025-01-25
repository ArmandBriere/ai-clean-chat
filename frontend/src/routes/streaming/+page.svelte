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

  // STUN/TURN servers for WebRTC
  const iceServers = [{ urls: 'stun:stun.l.google.com:19302' }];

  onMount(() => {
    // Connect to the signaling server
    ws = new WebSocket(`${PUBLIC_SERVER_WS_URL}/join?roomID=${roomID}&userID=${v4()}`);

    ws.onopen = () => {
      console.log('WebSocket connection established');
    };

    ws.onmessage = async (event) => {
      const message = JSON.parse(event.data);
      console.log('Received signaling message:', message);

      if (message.offer) {
        await handleOffer(message.offer);
      } else if (message.answer) {
        await handleAnswer(message.answer);
      } else if (message.iceCandidate) {
        await handleIceCandidate(message.iceCandidate);
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
    peerConnection = new RTCPeerConnection({ iceServers });

    // Handle ICE candidates
    peerConnection.onicecandidate = (event) => {
      if (event.candidate) {
        ws.send(JSON.stringify({ iceCandidate: event.candidate }));
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
      try {
        const offer = await peerConnection.createOffer();
        await peerConnection.setLocalDescription(offer);
        ws.send(JSON.stringify({ offer }));
      } catch (error) {
        console.error('Error during negotiation:', error);
      }
    };
  }

  async function handleOffer(offer: RTCSessionDescriptionInit) {
    try {
      await peerConnection.setRemoteDescription(offer);
      const answer = await peerConnection.createAnswer();
      await peerConnection.setLocalDescription(answer);
      ws.send(JSON.stringify({ answer }));
    } catch (error) {
      console.error('Error handling offer:', error);
    }
  }

  async function handleAnswer(answer: RTCSessionDescriptionInit) {
    try {
      await peerConnection.setRemoteDescription(answer);
    } catch (error) {
      console.error('Error handling answer:', error);
    }
  }

  async function handleIceCandidate(candidate: RTCIceCandidate) {
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
      <video playsInline autoPlay class="w-full rounded-lg border shadow" bind:this={otherVideo}
      ></video>
    </div>
  </div>
</main>

<style>
  main {
    max-width: 800px;
  }
</style>
