<script lang="ts">
  import { goto } from '$app/navigation';
  import { PUBLIC_SERVER_URL } from '$env/static/public';
  import Image from '@/lib/components/Image.svelte';
  import logo1 from '$lib/assets/logo/osedea-logo.png';
  import logo2 from '$lib/assets/logo/X.png';
  import logo3 from '$lib/assets/logo/confoo-logo.png';
  import avatar from '$lib/assets/avatar.jpeg';
  const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December'
  ];
  const date = new Date();
  const year = date.getFullYear().toString();
  const month = months[date.getMonth()];
  const day = date.getDate().toString();

  let meetingCode: string = $state('');

  // getRoomID function to fetch the room ID from the backend
  async function createRoom(): Promise<string> {
    let resp = await fetch(`${PUBLIC_SERVER_URL}/api/backend/create`);
    const { room_id } = await resp.json();

    return room_id;
  }

  // Check if the input is empty
  function isEmpty(item: string): boolean {
    return item === '';
  }

  // onKeydown event handler function to submit form on Enter key press
  function onKeyDown(e: KeyboardEvent): void {
    if (e.key === 'Enter' || e.code === 'Enter') {
      goToMeeting();
    }
  }

  // Redirect to the meeting page
  async function goToMeeting(): Promise<void> {
    if (isEmpty(meetingCode)) {
      console.debug('Meeting code is empty');
      return;
    }
    goto(`/chat/${meetingCode}`);
  }

  // Redirect to the streaming page
  async function newMeeting(): Promise<void> {
    const roomId = await createRoom();
    goto(`/chat/${roomId}`);
  }
</script>

<div class="flex h-svh flex-col">
  <header>
    <!-- wrapper -->
    <div class="flex w-full justify-between p-8">
      <!-- left nav -->
      <a href="/" class="flex items-center">
        <div class="flex items-center space-x-4">
          <Image src={logo1} alt="logo" class="w-full object-contain" />
          <Image src={logo2} alt="logo" class="w-full object-contain" />
          <Image src={logo3} alt="logo" class="w-full object-contain" />
        </div>
      </a>
      <!-- right nav -->
      <div>
        <div class="w-14 overflow-hidden rounded-full">
          <a href="/health">
            <Image src={avatar} alt="avatar" class="w-full object-contain" />
          </a>
        </div>
      </div>
    </div>
  </header>

  <div class="flex flex-grow">
    <div class="grid-[1fr] m-auto grid w-full max-w-xl justify-center text-center">
      <h1 class="mb-8 text-2xl"
        >AI-Powered Clean Talk:<br />Advanced Profanity Detection Patterns</h1
      >
      <div class="mb-8">
        <p>
          Advanced profanity detection goes beyond simple keyword searches. This session will cover
          data collection and training of a profanity detection model, and its integration with live
          Speech-to-Text AI technologies. We will also look at the challenges around handling live
          data.
        </p>
      </div>
      <div class="flex w-full justify-center space-x-4">
        <button
          onclick={newMeeting}
          class="flex items-center rounded bg-blue-600 px-5 py-3 font-medium text-white"
        >
          <span class="material-symbols-outlined mr-2"> videocam </span>
          New meeting
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
  </div>
  <footer>
    <div class="flex justify-between p-8 font-light">
      <div>{month} {day} {year}</div>
      <div>Montreal, Canada</div>
    </div>
  </footer>
</div>
<svelte:document on:keydown={onKeyDown} />
