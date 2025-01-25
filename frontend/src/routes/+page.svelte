<script lang="ts">
  import { goto } from '$app/navigation';
  import { PUBLIC_SERVER_URL } from '$env/static/public';

  let meetingCode: string = $state('');

  // getRoomID function to fetch the room ID from the backend
  async function setRoomID() {
    let resp = await fetch(`${PUBLIC_SERVER_URL}/api/backend/create`);
    const { room_id } = await resp.json();

    meetingCode = room_id;
  }

  // Check if the input is empty
  function isEmpty(item: string) {
    return item === '';
  }

  // onKeydown event handler function to submit form on Enter key press
  function onKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' || e.code === 'Enter') {
      goToMeeting();
    }
  }

  // Redirect to the meeting page
  async function goToMeeting() {
    if (isEmpty(meetingCode)) {
      await setRoomID();
    }
    goto(`/${meetingCode}`);
  }

  // Redirect to the streaming page
  async function goToStreaming() {
    goto('/streaming');
  }
</script>

<div class="grid-[1fr] m-auto grid w-full max-w-xl justify-center text-center">
  <h1 class="mb-8 text-2xl">AI-Powered Clean Talk:<br />Advanced Profanity Detection Patterns</h1>
  <div class="mb-8">
    <p>
      Advanced profanity detection goes beyond simple keyword searches. This session will cover data
      collection and training of a profanity detection model, and its integration with live
      Speech-to-Text AI technologies. We will also look at the challenges around handling live data.
    </p>
  </div>
  <div class="flex w-full justify-center space-x-4">
    <button
      onclick={goToMeeting}
      class="flex items-center rounded bg-blue-600 px-5 py-3 font-medium text-white"
    >
      <span class="material-symbols-outlined mr-2"> videocam </span>
      New meeting
    </button>
    <button
      onclick={goToStreaming}
      class="flex items-center rounded bg-blue-600 px-5 py-3 font-medium text-white"
    >
      <span class="material-symbols-outlined mr-2"> videocam </span>
      Streaming
    </button>
    <div class="flex items-center rounded border px-5 py-3 text-gray-600">
      <span class="material-symbols-outlined mr-2"> keyboard </span>
      <input
        bind:value={meetingCode}
        class="border-none outline-none"
        type="text"
        placeholder="Enter code"
      />
    </div>
    <button
      disabled={isEmpty(meetingCode)}
      onclick={goToMeeting}
      class="w-20 rounded-sm bg-blue-600 px-4 py-2 font-medium text-white disabled:bg-gray-300 disabled:text-gray-600"
    >
      Join
    </button>
  </div>
</div>

<svelte:document on:keydown={onKeyDown} />
