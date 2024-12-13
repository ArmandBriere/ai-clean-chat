<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { v4 as uuidv4 } from 'uuid';

	let ws: WebSocket | null = null;
	let pc: RTCPeerConnection | null = null;
	let localStream: MediaStream | null = null;
	let userId: string = uuidv4();

	let isStreaming = false;

	onMount(async () => {
		startConnection();
	});

	onDestroy(() => {
		cleanup();
	});

	async function startConnection() {
		try {
			ws = new WebSocket(`ws://localhost:8081/ws`);

			ws.onopen = async () => {
				console.log('WebSocket connected');

				try {
					pc = new RTCPeerConnection();

					localStream = await navigator.mediaDevices.getUserMedia({ audio: true });
					if (localStream) {
						localStream.getTracks().forEach((track) => pc?.addTrack(track));
					}

					pc.onicecandidate = (event) => {
						if (event.candidate) {
							let answer: IceCandidateMessage = {
								type: 'iceCandidate',
								candidate: event.candidate
							};
							ws?.send(JSON.stringify(answer));
						}
					};

					pc.onconnectionstatechange = async () => {};

					pc.oniceconnectionstatechange = () => {
						if (pc?.iceConnectionState === 'failed') {
							console.log('ICE connection failed, restarting');
							restartConnection();
						}
					};

					const offerDescription = await pc.createOffer();
					await pc.setLocalDescription(offerDescription);

					if (offerDescription.sdp) {
						let offer: OfferMessage = {
							type: 'offer',
							sdp: offerDescription.sdp
						};
						ws?.send(JSON.stringify(offer));
					}
				} catch (getUserMediaError) {
					console.error('Error getting user media:', getUserMediaError);
					cleanup();
					return; // Important: Exit the onopen handler if getUserMedia fails
				}

				if (!ws) {
					return;
				}
				ws.onmessage = async (event) => {
					const message = JSON.parse(event.data);
					if (message.type === 'answer') {
						let answer: AnswerMessage = {
							type: 'answer',
							sdp: message.sdp
						};
						const remoteDesc = new RTCSessionDescription(answer);
						try {
							await pc?.setRemoteDescription(remoteDesc);
						} catch (setRemoteDescError) {
							console.error('Error setting remote description:', setRemoteDescError);
							cleanup();
						}
					} else if (message.type === 'iceCandidate') {
						try {
							if (message.candidate) {
								// Correctly parse the JSON string to an object
								const candidate = JSON.parse(message.candidate);
								await pc?.addIceCandidate(new RTCIceCandidate(candidate));
							}
						} catch (e) {
							console.error('Error adding ICE candidate:', e);
						}
					}
				};

				ws.onclose = () => {
					console.log('WebSocket closed');
					cleanup();
				};

				ws.onerror = (error) => {
					console.error('WebSocket error:', error);
					cleanup();
				};
			};

			ws.onerror = (error) => {
				console.error('WebSocket connection error:', error);
				cleanup();
			};
		} catch (error) {
			console.error('Error starting connection:', error);
			cleanup();
		}
	}

	async function startStreaming() {
		if (ws) {
			isStreaming = !isStreaming;
			let message: StreamingMessage = {
				type: 'streaming',
				isStreaming: isStreaming
			};
			ws.send(JSON.stringify(message));
		} else {
			console.error('WebSocket not connected');
		}
	}

	function restartConnection() {
		cleanup();
		startConnection();
	}

	function cleanup() {
		if (ws) {
			ws.close();
			ws = null;
		}
		if (pc) {
			pc.close();
			pc = null;
		}
		if (localStream) {
			localStream.getTracks().forEach((track) => track.stop());
			localStream = null;
		}
	}
</script>

<div class="justify-center text-center">
	{#if userId}
		<div class="m-4 p-4 bg-[darkgray] rounded">
			<p>
				{userId}
			</p>
		</div>
	{/if}
	<div class="m-4 p-4 bg-[darkgray] rounded">
		<button on:click={startStreaming} aria-label="Start WebRTC">Start WebRTC</button>
	</div>
</div>
