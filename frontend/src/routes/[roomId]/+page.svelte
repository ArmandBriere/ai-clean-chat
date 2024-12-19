<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  const servers = {
    iceServers: [{ urls: ['stun:stun1.1.google.com:19302', 'stun:stun2.1.google.com:19302'] }],
    iceCandidatePoolSize: 10
  };
  let userStream: MediaStream;
  let remoteStream: MediaStream;
  onDestroy(() => {
    stopMediaTracks();
  });
  onMount(async () => {
    const stream = await navigator.mediaDevices.getUserMedia({
      video: true,
      audio: true
    });
    userStream = stream;
  });
  onMount(() => {
    let pc = new RTCPeerConnection(servers);

    // TODO: Push tracks from userStream to peer connection

    // TODO: Pull tracks from remoteStream, add to video stream
  });

  function stopMediaTracks() {
    if (remoteStream?.active) {
      const tracks = remoteStream.getTracks();
      for (const track of tracks) {
        track.stop();
      }
    }

    if (userStream?.active) {
      const tracks = userStream.getTracks();
      for (const track of tracks) {
        track.stop();
      }
    }
  }
</script>

<video id="localVideo" autoplay muted playsinline></video>
<video id="remoteVideo" autoplay playsinline></video>

<!-- <script> -->
<!-- 	import { onMount } from 'svelte'; -->
<!---->
<!-- 	let localStream; -->
<!-- 	let remoteStream; -->
<!-- 	let peerConnection; -->
<!-- 	const servers = { iceServers: [{ urls: 'stun:stun1.l.google.com:19302' }] }; -->
<!---->
<!-- 	onMount(async () => { -->
<!-- 		// Get local media stream -->
<!-- 		localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true }); -->
<!-- 		document.querySelector('#localVideo').srcObject = localStream; -->
<!---->
<!-- 		// Setup peer connection -->
<!-- 		peerConnection = new RTCPeerConnection(servers); -->
<!---->
<!-- 		// Add local stream tracks to the connection -->
<!-- 		localStream.getTracks().forEach((track) => peerConnection.addTrack(track, localStream)); -->
<!---->
<!-- 		// Setup remote stream -->
<!-- 		remoteStream = new MediaStream(); -->
<!-- 		document.querySelector('#remoteVideo').srcObject = remoteStream; -->
<!---->
<!-- 		// Handle received remote stream tracks -->
<!-- 		peerConnection.ontrack = (event) => { -->
<!-- 			event.streams[0].getTracks().forEach((track) => remoteStream.addTrack(track)); -->
<!-- 		}; -->
<!---->
<!-- 		// Handle ICE candidates -->
<!-- 		peerConnection.onicecandidate = (event) => { -->
<!-- 			if (event.candidate) { -->
<!-- 				// Send the ICE candidate to the remote peer -->
<!-- 				// This requires signaling server setup -->
<!-- 				console.log('New ICE candidate: ', event.candidate); -->
<!-- 			} -->
<!-- 		}; -->
<!---->
<!-- 		// Create an offer to connect -->
<!-- 		const offer = await peerConnection.createOffer(); -->
<!-- 		await peerConnection.setLocalDescription(offer); -->
<!---->
<!-- 		// Send offer to remote peer through a signaling server -->
<!-- 	}); -->
<!---->
<!-- 	const peerConnection = new RTCPeerConnection(); -->
<!---->
<!-- 	// Event to send ICE candidates to the remote peer -->
<!-- 	peerConnection.onicecandidate = ({ candidate }) => { -->
<!-- 		if (candidate) { -->
<!-- 			signalingServer.send(JSON.stringify({ type: 'ice-candidate', candidate })); -->
<!-- 		} -->
<!-- 	}; -->
<!---->
<!-- 	// Event to listen for remote stream -->
<!-- 	peerConnection.ontrack = (event) => { -->
<!-- 		remoteVideo.srcObject = event.streams[0]; -->
<!-- 	}; -->
<!---->
<!-- 	// Create an offer and set it as the local description -->
<!-- 	async function startCall() { -->
<!-- 		const offer = await peerConnection.createOffer(); -->
<!-- 		await peerConnection.setLocalDescription(offer); -->
<!-- 		signalingServer.send(JSON.stringify({ type: 'offer', offer })); -->
<!-- 	} -->
<!---->
<!-- 	signalingServer.onmessage = async (message) => { -->
<!-- 		const { type, offer, answer, candidate } = JSON.parse(message.data); -->
<!-- 		if (type === 'offer') { -->
<!-- 			await peerConnection.setRemoteDescription(new RTCSessionDescription(offer)); -->
<!-- 			const answer = await peerConnection.createAnswer(); -->
<!-- 			await peerConnection.setLocalDescription(answer); -->
<!-- 			signalingServer.send(JSON.stringify({ type: 'answer', answer })); -->
<!-- 		} else if (type === 'answer') { -->
<!-- 			await peerConnection.setRemoteDescription(new RTCSessionDescription(answer)); -->
<!-- 		} else if (type === 'ice-candidate') { -->
<!-- 			peerConnection.addIceCandidate(new RTCIceCandidate(candidate)); -->
<!-- 		} -->
<!-- 	}; -->
<!-- </script> -->
<!---->
<!-- <video id="localVideo" autoplay muted playsinline></video> -->
<!-- <video id="remoteVideo" autoplay playsinline></video> -->
