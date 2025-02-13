<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { PUBLIC_SERVER_WS_URL } from '$env/static/public';
  import { ANSWER, ICE_CANDIDATE, STREAMING, TRANSCRIPTION } from '@/lib/constants/constants';
  import type { AnalyzedMessage } from '@/lib/constants/types';

  import { page } from '$app/state';
  let roomID = page.params.roomID;

  // Transcription
  let wsTranscription: WebSocket | null = null;
  let pcTranscription: RTCPeerConnection | null = null;
  let streamTranscription: MediaStream | null = null;

  let messages: AnalyzedMessage[] = $state([]);
  let isStreaming = false;

  let microphoneList: MediaDeviceInfo[] = $state([]);
  let selectedMicrophone: string = $state('default');

  onMount(async () => {
    navigator.mediaDevices.enumerateDevices().then((devices) => {
      microphoneList = devices.filter((device) => device.kind === 'audioinput');
    });

    startTranscriptionConnection();
  });

  function setSelectedMicrophone() {
    console.log(selectedMicrophone);
    restartConnection();
  }

  onDestroy(() => {
    cleanup();
  });

  async function startTranscriptionConnection() {
    try {
      wsTranscription = new WebSocket(`${PUBLIC_SERVER_WS_URL}/ws`);

      wsTranscription.onopen = async () => {
        console.log('wsTranscription connected');

        try {
          pcTranscription = new RTCPeerConnection();

          pcTranscription.onicecandidate = (event) => {
            if (event.candidate) {
              let answer: IceCandidateMessage = {
                type: 'iceCandidate',
                candidate: event.candidate
              };
              wsTranscription?.send(JSON.stringify(answer));
            }
          };

          streamTranscription = await navigator.mediaDevices.getUserMedia({
            audio: {
              channelCount: 1,
              sampleRate: 48000,
              deviceId: selectedMicrophone
            }
          });
          if (streamTranscription) {
            streamTranscription.getTracks().forEach((track) => {
              pcTranscription?.addTrack(track);
            });
          }

          pcTranscription.onconnectionstatechange = async () => {};

          pcTranscription.oniceconnectionstatechange = () => {
            if (pcTranscription?.iceConnectionState === 'failed') {
              console.log('ICE connection failed, restarting');
              restartConnection();
            }
          };

          const offerDescription = await pcTranscription.createOffer();
          await pcTranscription.setLocalDescription(offerDescription);

          if (offerDescription.sdp) {
            let offer: OfferMessage = {
              type: 'offer',
              sdp: offerDescription.sdp
            };
            wsTranscription?.send(JSON.stringify(offer));
          }
        } catch (getUserMediaError) {
          console.error('Error getting user media:', getUserMediaError);
          cleanup();
          return;
        }

        if (!wsTranscription) {
          return;
        }
        wsTranscription.onmessage = async (event) => {
          const message = JSON.parse(event.data);
          if (message.type === ANSWER) {
            let answer: AnswerMessage = {
              type: ANSWER,
              sdp: message.sdp
            };
            // @ts-ignore
            const remoteDesc = new RTCSessionDescription(answer);
            try {
              await pcTranscription?.setRemoteDescription(remoteDesc);
            } catch (setRemoteDescError) {
              console.error('Error setting remote description:', setRemoteDescError);
              cleanup();
            }
          } else if (message.type === ICE_CANDIDATE) {
            try {
              if (message.candidate) {
                // Correctly parse the JSON string to an object
                const candidate = JSON.parse(message.candidate);
                await pcTranscription?.addIceCandidate(new RTCIceCandidate(candidate));
              }
            } catch (e) {
              console.error('Error adding ICE candidate:', e);
            }
          } else if (message.type === STREAMING) {
            if (message.isStreaming) {
              console.log('Starting streaming');
            } else {
              console.log('Stopping streaming');
            }
          } else if (message.type === TRANSCRIPTION) {
            var newMessage: AnalyzedMessage = {
              uuid: message.uuid,
              text: message.text,
              profanityScore: message.profanity_score
            };

            const updatedMessages = [...messages, newMessage];
            messages = updatedMessages.slice(-25);
          }
        };

        wsTranscription.onclose = () => {
          console.log('WebSocket closed');
        };

        wsTranscription.onerror = (error) => {
          console.error('WebSocket error:', error);
        };
      };

      wsTranscription.onerror = (error) => {
        console.error('WebSocket connection error:', error);
      };
    } catch (error) {
      console.error('Error starting connection:', error);
      cleanup();
    }
  }

  async function startStreaming() {
    if (wsTranscription) {
      isStreaming = !isStreaming;
      let message: StreamingMessage = {
        type: 'streaming',
        isStreaming: isStreaming
      };
      wsTranscription.send(JSON.stringify(message));
    } else {
      console.error('WebSocket not connected');
    }
  }

  function restartConnection() {
    cleanup();
    startTranscriptionConnection();
  }

  function cleanup() {
    wsTranscription?.close();
    pcTranscription?.close();
    streamTranscription?.getTracks().forEach((track) => track.stop());
  }
</script>

<div class="w-full justify-center text-center">
  <div class="m-4 rounded bg-[darkgray] p-4">
    <button onclick={startStreaming} aria-label="Start WebRTC">Start WebRTC</button>
  </div>
  <div class="m-4 text-left text-gray-600">
    <span class="font-normal"> Transcription: </span>
    {#each messages as message}
      <span class={message.profanityScore > 0.9 ? 'text-red-500 line-through' : ''}>
        {message.text}
      </span>
    {/each}
  </div>
  <div>
    <select bind:value={selectedMicrophone} onchange={setSelectedMicrophone}>
      {#each microphoneList as microphone}
        <option value={microphone.deviceId}>{microphone.label}</option>
      {/each}
    </select>
  </div>
  <div>
    {selectedMicrophone}
  </div>
</div>
