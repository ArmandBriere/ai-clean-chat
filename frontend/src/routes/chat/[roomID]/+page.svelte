<script lang="ts">
  import MeetingID from '../../../lib/components/MeetingID.svelte';

  import { onMount, onDestroy } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { PUBLIC_SERVER_WS_URL } from '$env/static/public';
  import { v4 } from 'uuid';
  import { OFFER, ANSWER, ICE_CANDIDATE, HANG_UP, EMOJI } from '@/lib/constants/constants';
  import type {
    StreamingOfferMessage,
    StreamingAnswerMessage,
    StreamingIceCandidateMessage,
    AnalyzedMessage,
    LLMAnalysis
  } from '@/lib/constants/types';

  import HangUp from '@/lib/components/HangUp.svelte';
  import Emojis from '@/lib/components/Emojis.svelte';
  import Transcription from '@/lib/components/Transcription.svelte';
  import Selector from '@/lib/components/Selector.svelte';
  import ToolTip from '@/lib/components/ToolTip.svelte';
  import InfoPanel from '@/lib/components/InfoPanel.svelte';

  let roomID: string = page.params.roomID;
  let userVideo: HTMLVideoElement;
  let otherVideo: HTMLVideoElement;

  let ws: WebSocket;
  let peerConnection: RTCPeerConnection;
  let localStream: MediaStream;

  let connectedUsers = $state(1);

  let selectedMicrophone = $state('default');
  let selectedCamera = $state('default');
  let isMicOn = $state(true);
  let isVideoOn = $state(true);
  let isClosedCaptionOn = $state(false);

  let showHangUpModal = $state(false);
  let showMicrophoneModal = $state(false);
  let showCameraModal = $state(false);
  let showEmojiModal = $state(false);
  let receivedEmoji = $state('');
  let showInfoPanel = $state(false);

  // Transcription
  let messages: AnalyzedMessage[] = $state([]);

  // LLM Analysis
  let llmAnalysis: LLMAnalysis[] = $state([]);

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
      } else if (message.type == HANG_UP) {
        connectedUsers = Math.max(connectedUsers - 1, 1);
      } else if (message.type == EMOJI) {
        receivedEmoji = message.payload;
      }
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    // Open the camera and set up local video
    openCamera();
  });

  onDestroy(() => {
    stopMediaTracks();
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
        if (otherVideo) {
          otherVideo.srcObject = event.streams[0];
        }
        connectedUsers = 2;
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

    // Handle connection state changes
    peerConnection.onconnectionstatechange = () => {
      console.log('Connection state:', peerConnection.connectionState);
      if (peerConnection.connectionState === 'disconnected') {
        connectedUsers = 1;
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

  // Change media track locally and for others
  const changeMediaDevice = async (deviceID: string, kind: 'audio' | 'video') => {
    console.log(`Changing ${kind} device to:`, deviceID);
    const tracks = localStream.getTracks().filter((track) => track.kind === kind);
    if (tracks.length === 0) {
      console.warn(`No ${kind} track found in localStream.`);
      return;
    }

    tracks[0].stop();
    localStream.removeTrack(tracks[0]);

    const constraints: MediaStreamConstraints =
      kind === 'audio' ? { audio: { deviceId: deviceID } } : { video: { deviceId: deviceID } };

    try {
      const newStream = await navigator.mediaDevices.getUserMedia(constraints);
      const newTrack = newStream.getTracks().find((track) => track.kind === kind);

      if (!newTrack) {
        console.error(`No ${kind} track found in new stream.`);
        return;
      }

      localStream.addTrack(newTrack);

      // Update track for others
      peerConnection.getSenders().forEach((sender) => {
        if (sender?.track?.kind === kind) {
          sender.replaceTrack(newTrack);
        }
      });
    } catch (error) {
      console.error(`Error changing ${kind} device:`, error);
    }
  };

  // Toggle On/Off
  const toggleMic = () => {
    const tracks = localStream.getAudioTracks();
    tracks[0].enabled = !tracks[0].enabled;
    isMicOn = tracks[0].enabled;
  };

  const toggleCamera = () => {
    const tracks = localStream.getVideoTracks();
    tracks[0].enabled = !tracks[0].enabled;
    isVideoOn = tracks[0].enabled;
  };

  const stopMediaTracks = () => {
    ws?.close();
    peerConnection?.close();
    localStream?.getTracks().forEach((track) => track.stop());
  };

  const handleHangUp = () => {
    ws.send(JSON.stringify({ type: HANG_UP }));
    showHangUpModal = false;
    goto('/');
  };

  const shareEmoji = (emoji: string) => {
    ws.send(JSON.stringify({ type: EMOJI, payload: emoji }));
  };

  const handleClosedCaption = () => {
    isClosedCaptionOn = !isClosedCaptionOn;
  };

  const handleInfoPanel = () => {
    showInfoPanel = !showInfoPanel;
  };
</script>

<div class="fixed left-0 top-0 flex h-full min-h-full w-full flex-col bg-[rgb(25,25,25)]">
  <div class="relative h-full w-full">
    <main
      class={`ease-[.5s cubic-bezier(0.4,0,0.2,1)] absolute ${!isClosedCaptionOn ? (showInfoPanel ? 'inset-[16px_332px_80px_16px]' : 'inset-[16px_16px_80px]') : showInfoPanel ? 'inset-[16px_332px_300px_16px]' : 'inset-[16px_16px_300px]'} transition-[bottom,right]`}
    >
      <div class="h-full transition-all duration-500 ease-in-out">
        <div class="flex h-full justify-center">
          <div
            class={`h-full transform transition-all duration-500 ease-in-out ${connectedUsers === 2 ? 'absolute -right-0 bottom-0 max-h-52' : ''}`}
          >
            <video
              class="aspect-video h-full rounded-md object-cover"
              class:rounded-lg={connectedUsers === 2}
              class:border={connectedUsers === 2}
              autoPlay
              muted
              playsInline
              bind:this={userVideo}
            >
              Your browser does not support the video tag.
            </video>
          </div>
          <video
            class:hidden={connectedUsers === 1}
            class="aspect-video rounded-md object-cover"
            autoPlay
            muted
            playsInline
            bind:this={otherVideo}
          >
            Your browser does not support the video tag.
          </video>
        </div>
      </div>
    </main>
  </div>

  <!-- transcription -->
  {#if isClosedCaptionOn}
    <div class="absolute bottom-[5rem] left-0 right-0">
      <div
        class="relative left-[16px] right-[16px] flex h-[calc(12.5rem+16px)] w-[calc(100%-32px)] flex-col flex-nowrap items-center justify-end overflow-hidden"
      >
        <div class="absolute left-0 right-0 top-0 mx-4 my-2 flex h-8 items-center">
          <h2 class="text-gray-200">Transcription</h2>
        </div>
        <div
          class="absolute bottom-0 left-0 right-0 flex h-[9.75rem] translate-y-0 transform flex-nowrap justify-start overflow-y-auto whitespace-pre pb-4 pl-[20vw] pr-4 pt-[0.875rem] text-left text-white"
        >
          {#each messages as message}
            <span class={message.profanityScore > 0.9 ? 'text-red-500 line-through' : ''}>
              {message.text}
            </span>
          {/each}
        </div>
        <div>
          {#each llmAnalysis as llmText}
            <div class="w-1/2 pt-1 text-left">
              <span class="font-bold text-white">{llmText.userMessage.toLowerCase()}:</span>
              <span class="text-gray-200">
                {llmText.analysis}
              </span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}

  <!-- footer -->
  <div class="absolute bottom-0 left-0 right-0">
    <div class="relative mt-auto flex h-20 w-full items-center justify-center">
      <div
        class="relative ml-3 flex h-full max-w-full flex-[1_1_25%] items-center overflow-ellipsis text-start"
      >
        <MeetingID {roomID} {connectedUsers} />
      </div>

      <div class="relative flex flex-[1_1_25%] justify-center space-x-4">
        <ToolTip displayText="Microphone settings">
          <Selector
            showModal={showMicrophoneModal}
            kind="audioinput"
            bind:selectedDevice={selectedMicrophone}
            isDeviceOn={isMicOn}
            closeDevice={toggleMic}
            changeMediaSource={changeMediaDevice}
            displayTop={true}
          />
        </ToolTip>

        <ToolTip displayText="Camera settings">
          <Selector
            showModal={showCameraModal}
            kind="videoinput"
            bind:selectedDevice={selectedCamera}
            isDeviceOn={isVideoOn}
            closeDevice={toggleCamera}
            changeMediaSource={changeMediaDevice}
            displayTop={true}
          />
        </ToolTip>

        <ToolTip displayText="Activate live transcription">
          <button
            onclick={handleClosedCaption}
            class={`flex items-center justify-center rounded-full p-3 transition-colors ${isClosedCaptionOn ? 'bg-green-500 text-black' : 'bg-[#333537] text-gray-200 hover:bg-[#3f4245]'}`}
          >
            <span class="material-symbols-outlined"> closed_caption </span>
          </button>
        </ToolTip>

        <ToolTip displayText="💩">
          <Emojis {showEmojiModal} {shareEmoji} {receivedEmoji} />
        </ToolTip>

        <ToolTip displayText="Bye 👋">
          <HangUp {handleHangUp} {showHangUpModal} displayTop={true} />
        </ToolTip>
        {#if isClosedCaptionOn}
          <Transcription
            {roomID}
            {selectedMicrophone}
            bind:messages
            bind:llmAnalysis
            bind:micStatus={isMicOn}
          />
        {/if}
      </div>
      <div class="mr-3 flex h-full flex-[1_1_25%] items-center justify-end">
        <nav class="flex">
          <div>
            <button
              onclick={handleInfoPanel}
              class="flex items-center justify-center rounded-full p-3 text-white transition-colors hover:bg-gray-200/10"
              ><span class="material-symbols-outlined" class:text-blue-600={showInfoPanel}
                >info</span
              ></button
            ></div
          >
        </nav>
      </div>
    </div>
  </div>
  <InfoPanel {showInfoPanel} {roomID} handleClose={handleInfoPanel} />
</div>
