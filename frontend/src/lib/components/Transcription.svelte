<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { PUBLIC_SERVER_WS_URL } from '$env/static/public';
  import {
    ANSWER,
    ICE_CANDIDATE,
    STREAMING,
    TRANSCRIPTION,
    LLM_ANALYSIS
  } from '@/lib/constants/constants';
  import type { AnalyzedMessage, LLMAnalysis } from '@/lib/constants/types';

  let {
    roomID,
    selectedMicrophone,
    messages = $bindable([]),
    llmAnalysis = $bindable([]),
    micStatus = $bindable(true)
  }: {
    roomID: string;
    selectedMicrophone: string;
    messages: AnalyzedMessage[];
    llmAnalysis: LLMAnalysis[];
    micStatus: boolean;
  } = $props();

  // Transcription
  let wsTranscription: WebSocket | null = null;
  let pcTranscription: RTCPeerConnection | null = null;
  let streamTranscription: MediaStream | null = null;

  let isStreaming = micStatus;

  onMount(async () => {
    console.log('Transcription component mounted');

    startTranscriptionConnection();
  });

  // Restart connection when selected microphone changes
  $effect(function restartConnectionOnMicUpdate() {
    console.log('Updated selected microphone', selectedMicrophone);
    if (selectedMicrophone) {
      restartConnection();
    }
  });

  // Restart connection when selected microphone changes
  $effect(function stopStreamingOnMicMute() {
    console.log('Update microphone status', micStatus);
    if (micStatus != isStreaming) {
      toggleStreaming();
    }
  });

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
          console.log('Received message:', event.data);
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

            const updatedMessages = [...messages, newMessage].map((msg, index, array) => {
              if (
                newMessage.profanityScore > 0.9 &&
                index >= array.length - 8 &&
                index < array.length
              ) {
                return {
                  ...msg,
                  profanityScore: newMessage.profanityScore
                };
              }
              return msg;
            });

            messages = updatedMessages.slice(-25);
          } else if (message.type === LLM_ANALYSIS) {
            var newLLMAnalysis: LLMAnalysis = {
              analysis: message.llm_analysis,
              userMessage: message.user_message
            };

            const updatedLLMAnalysis = [...llmAnalysis, newLLMAnalysis];
            llmAnalysis = updatedLLMAnalysis.slice(-25);
          }
        };
        toggleStreaming(true);
      };

      wsTranscription.onclose = () => {
        console.log('WebSocket closed');
      };

      wsTranscription.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (error) {
      console.error('Error starting connection:', error);
      cleanup();
    }
  }

  async function toggleStreaming(forceStart?: boolean) {
    if (wsTranscription) {
      console.log('Toggling streaming', isStreaming);
      isStreaming = !isStreaming;
      if (forceStart) {
        isStreaming = true;
      }
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
    console.log('Cleaning up transcription');
    wsTranscription?.close();
    pcTranscription?.close();
    streamTranscription?.getTracks().forEach((track) => track.stop());
  }
</script>
