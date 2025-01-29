<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { PUBLIC_SERVER_WS_URL } from '$env/static/public';
  import { v4 } from 'uuid';
  import { OFFER, ANSWER, ICE_CANDIDATE } from '@/lib/constants/constants';
  import type {
    StreamingOfferMessage,
    StreamingAnswerMessage,
    StreamingIceCandidateMessage
  } from '@/lib/constants/types';

  let roomID: string = 'example-room'; // Will be replaced with dynamic room ID
  let userVideo: HTMLVideoElement;
  let otherVideo: HTMLVideoElement;

  let ws: WebSocket;
  let peerConnection: RTCPeerConnection;
  let localStream: MediaStream;

  onMount(() => {
    // Connect to the signaling server
    ws = new WebSocket(`${PUBLIC_SERVER_WS_URL}/join?roomID=${roomID}&userID=${v4()}`);

    ws.onopen = () => {
      console.log('WebSocket connection established');
    };

    ws.onmessage = async (event) => {
      const message = JSON.parse(event.data);
      console.log('Received signaling message:', message);

      // Handle different message types
      if (message.type === OFFER) {
        let data: StreamingOfferMessage = message;
        await handleOffer(data.payload);
      } else if (message.type == ANSWER) {
        let data: StreamingAnswerMessage = message;
        await handleAnswer(data.payload);
      } else if (message.type == ICE_CANDIDATE) {
        let data: StreamingIceCandidateMessage = message;
        await handleIceCandidate(data.payload);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    // Open the camera and set up local video
    openCamera();
  });

  onDestroy(() => {
    ws?.close();
    peerConnection?.close();
    localStream?.getTracks().forEach((track) => track.stop());
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
        const message: any = {
          type: ICE_CANDIDATE,
          payload: event.candidate
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
        const offerMessage: any = {
          type: OFFER,
          payload: offer
        };
        ws.send(JSON.stringify(offerMessage));
      }
    };
  }

  // handleOffer sets the remote description of the peer connection and creates an answer
  async function handleOffer(offer: RTCSessionDescriptionInit) {
    await peerConnection.setRemoteDescription(offer);
    const answer = await peerConnection.createAnswer();
    await peerConnection.setLocalDescription(answer);

    if (answer.type === ANSWER && answer.sdp) {
      const answerMessage: any = {
        type: ANSWER,
        payload: answer
      };
      ws.send(JSON.stringify(answerMessage));
    }
  }

  // handleAnswer sets the remote description of the peer connection
  async function handleAnswer(answer: RTCSessionDescriptionInit) {
    try {
      await peerConnection.setRemoteDescription(answer);
    } catch (error) {
      console.error('Error handling answer:', error);
    }
  }

  // handleIceCandidate adds the received ICE candidate to the peer connection
  async function handleIceCandidate(candidate: RTCIceCandidateInit) {
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
