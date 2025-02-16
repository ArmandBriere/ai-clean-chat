<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { PUBLIC_SERVER_WS_URL, PUBLIC_SERVER_URL } from '$env/static/public';
  import { v4 } from 'uuid';
  import { OFFER, ANSWER, ICE_CANDIDATE, HANG_UP } from '@/lib/constants/constants';
  import type {
    StreamingOfferMessage,
    StreamingAnswerMessage,
    StreamingIceCandidateMessage,
    AnalyzedMessage
  } from '@/lib/constants/types';

  import HangUp from '@/lib/components/HangUp.svelte';
  import Emojis from '@/lib/components/Emojis.svelte';
  import Transcription from '@/lib/components/Transcription.svelte';
  import MicrophoneSelector from '@/lib/components/MicrophoneSelector.svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';

  let roomID = page.params.roomID;
  let userVideo: HTMLVideoElement;
  let otherVideo: HTMLVideoElement;

  let ws: WebSocket;
  let peerConnection: RTCPeerConnection;
  let localStream: MediaStream;

  let connectedUsers = $state(1);

  let selectedMicrophone = $state('default');
  let isMicOn = $state(true);
  let isVideoOn = $state(true);
  let isClosedCaptionOn = $state(false);
  let isBackHandOn = $state(true);

  let showHangUpModal = $state(false);
  let showMicrophoneModal = $state(false);
  let showEmojiModal = $state(false);
  let receivedEmoji = $state('');

  // Transcription
  let messages: AnalyzedMessage[] = $state([]);

  onMount(() => {
    console.log(otherVideo?.srcObject);
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
        connectedUsers = 1;
      } else if (message.type == 'emoji') {
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
        otherVideo.srcObject = event.streams[0];
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

  const toggleCamera = () => {
    const tracks = localStream.getVideoTracks();
    tracks[0].enabled = !tracks[0].enabled;
    isVideoOn = tracks[0].enabled;
  };

  $effect(() => {
    console.log('Mic status:', isMicOn);
    if (localStream) {
      const tracks = localStream.getAudioTracks();
      tracks[0].enabled = !tracks[0].enabled;
      isMicOn = tracks[0].enabled;
    }
  });

  $effect(() => {
    console.log('Selected microphone:', selectedMicrophone);
  });

  const stopMediaTracks = () => {
    ws?.close();
    peerConnection?.close();
    localStream?.getTracks().forEach((track) => track.stop());
  };

  const handleHangUp = () => {
    ws.send(JSON.stringify({ type: HANG_UP }));
    showHangUpModal = false;
    stopMediaTracks();
    goto('/');
  };

  const shareEmoji = (emoji: string) => {
    ws.send(JSON.stringify({ type: 'emoji', payload: emoji }));
  };

  const handleClosedCaption = () => {
    isClosedCaptionOn = !isClosedCaptionOn;
  };

  const shareMeetingURI = () => {
    const copyText = `${PUBLIC_SERVER_URL}/chat/${roomID}`;
    const theClipboard = navigator.clipboard;
    // write text returns a promise, so use `then` to react to success
    theClipboard.writeText(copyText).then(() => console.log('copied to clipboard'));
    // TODO: animate tooltip to show copied to clipboard with data
  };
</script>

<main
  class:p-16={isClosedCaptionOn}
  class:py-8={!isClosedCaptionOn}
  class="flex min-h-screen flex-col items-center justify-center bg-[rgb(25,25,25)] transition-all"
>
  <div class="relative w-full max-w-screen-2xl">
    <div class="relative aspect-video w-full overflow-hidden rounded-lg bg-gray-900">
      <div
        class={`h-full w-full ${connectedUsers === 2 ? 'opacity-100' : 'opacity-0'} transition-opacity duration-500`}
      >
        <div class="h-full w-full overflow-hidden rounded-lg border-2 border-white">
          <video
            class="h-full w-full object-cover"
            autoPlay
            playsInline
            muted
            bind:this={otherVideo}
          >
            Your browser does not support the video tag.
          </video>
        </div>
      </div>
      <div
        class={`absolute ${connectedUsers === 2 ? 'bottom-4 right-4 z-30 aspect-video w-1/4' : 'inset-0 z-10'} transition-all duration-500 ease-in-out`}
      >
        <div
          class={`h-full w-full overflow-hidden rounded-lg ${
            connectedUsers === 2 ? 'border-2 border-white' : 'border-2 border-white shadow-lg'
          }`}
        >
          <video
            class="h-full w-full object-cover"
            autoPlay
            muted
            playsInline
            bind:this={userVideo}
          >
            Your browser does not support the video tag.
          </video>
        </div>
      </div>
    </div>
    <div class="m-4 text-left text-gray-600">
      <span class="font-normal"> Transcription: </span>
      {#each messages as message}
        <span class={message.profanityScore > 0.9 ? 'text-red-500 line-through' : ''}>
          {message.text}
        </span>
      {/each}
    </div>
    <!-- TODO: add profanity detection -->
    <!-- {#if isClosedCaptionOn} -->
    <!--   <div -->
    <!--     class="relative mt-2 flex w-full max-w-5xl flex-grow flex-col overflow-scroll border border-white p-2 text-white" -->
    <!--   > -->
    <!--     <div class="flex w-full justify-end"> -->
    <!--       <label for="hide-profanity" class="flex cursor-pointer items-center"> -->
    <!--         <div class="mr-3 text-gray-200"> Hide profanities </div> -->
    <!--         <div class="relative"> -->
    <!--           <input type="checkbox" id="hide-profanity" class="sr-only" /> -->
    <!--           <div class="block h-8 w-14 rounded-full bg-gray-200"></div> -->
    <!--           <div class="dot absolute left-1 top-1 h-6 w-6 rounded-full bg-gray-600 transition"> -->
    <!--             <span class="material-symbols-outlined font-light"> close </span> -->
    <!--           </div> -->
    <!--         </div> -->
    <!--       </label> -->
    <!--     </div> -->
    <!--     <p>> profanity section</p> -->
    <!--     <p>> profanity section</p> -->
    <!--     <p>> profanity section</p> -->
    <!--     <p>> profanity section</p> -->
    <!--     <p>> profanity section</p> -->
    <!--     <p>> profanity section</p> -->
    <!--     <p>> profanity section</p> -->
    <!--   </div> -->
    <!-- {/if} -->
    <div class="relative mt-4 flex w-full justify-center">
      <div
        class="ml-3 flex max-w-full flex-[1_1_25%] items-center overflow-hidden overflow-ellipsis text-start"
      >
        <span class="text-white">
          <button title="Copy room id" onclick={shareMeetingURI}> {roomID}</button>
        </span>
      </div>
      <div class="flex flex-[1_1_25%] justify-center space-x-4">
        <MicrophoneSelector
          bind:isMicOn
          {showMicrophoneModal}
          bind:selectedMicrophone
          displayTop={true}
        />

        <button
          title="Toggle Camera"
          onclick={toggleCamera}
          class={`my-auto flex select-none items-center justify-center rounded-full p-3 no-underline hover:brightness-75 ${isVideoOn ? 'bg-gray-200 dark:text-black' : 'bg-red-500'}`}
        >
          <span class="material-symbols-outlined"> videocam </span>
        </button>
        <button
          onclick={handleClosedCaption}
          class={`my-auto flex select-none items-center justify-center rounded-full p-3 no-underline hover:brightness-75 ${isClosedCaptionOn ? 'bg-green-500' : ' bg-gray-200 dark:text-black'}`}
        >
          <span class="material-symbols-outlined"> closed_caption </span>
        </button>

        <Emojis {showEmojiModal} {shareEmoji} {receivedEmoji} />

        <!-- Close meeting -->
        <HangUp {handleHangUp} {showHangUpModal} displayTop={true} />

        {#if isClosedCaptionOn}
          <Transcription {roomID} {selectedMicrophone} bind:messages></Transcription>
        {/if}
      </div>
      <div class="flex flex-[1_1_25%] items-center justify-end">
        <nav class="flex">
          <div>
            <button
              class="flex items-center justify-center rounded-full p-3 text-white transition-colors hover:bg-gray-200/10"
              ><span class="material-symbols-outlined">info</span></button
            ></div
          >
        </nav>
      </div>
    </div>
  </div>
</main>
