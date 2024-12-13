<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { v4 as uuidv4 } from 'uuid';

	let ws: WebSocket | null = null;
	let pc: RTCPeerConnection | null = null;
	let localStream: MediaStream | null = null;
	let userId: string = uuidv4();

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
							ws?.send(JSON.stringify({ type: 'iceCandidate', candidate: event.candidate }));
						}
					};

					pc.onconnectionstatechange = async () => {};

					pc.oniceconnectionstatechange = () => {
						if (pc?.iceConnectionState === 'failed') {
							console.log('ICE connection failed, restarting');
							restartConnection();
						}
					};

					const offer = await pc.createOffer();
					await pc.setLocalDescription(offer);
					ws?.send(JSON.stringify({ type: 'offer', sdp: offer.sdp }));
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
						const remoteDesc = new RTCSessionDescription({ type: 'answer', sdp: message.sdp });
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
		// Add code to start streaming and stop it on demand
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

<div>
	{#if userId}
		<div class="userId">
			<p>
				{userId}
			</p>
		</div>
	{/if}
	<div class="startWebRTC">
		<button on:click={startStreaming} aria-label="Start WebRTC">Start WebRTC</button>
	</div>
</div>

<style>
	.userId,
	.startWebRTC {
		margin: 1rem;
		background-color: darkgray;
		padding: 1rem;
		border-radius: 5px;
		justify-content: center;
		text-align: center;
	}
</style>
