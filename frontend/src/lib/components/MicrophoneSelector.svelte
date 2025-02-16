<script lang="ts">
  import { onMount } from 'svelte';

  let {
    showMicrophoneModal = false,
    selectedMicrophone = $bindable('default'),
    isMicOn = true,
    displayTop = false,
    closeMic = $bindable()
  }: {
    showMicrophoneModal: boolean;
    selectedMicrophone: string;
    isMicOn: boolean;
    displayTop: boolean;
    closeMic: () => void;
  } = $props();

  let microphoneList: MediaDeviceInfo[] = $state([]);

  onMount(async () => {
    navigator.mediaDevices.enumerateDevices().then((devices) => {
      microphoneList = devices.filter((device) => device.kind === 'audioinput');
    });
  });

  function setSelectedMicrophone(newMicrophone: MediaDeviceInfo) {
    selectedMicrophone = newMicrophone?.deviceId || 'default';
    showMicrophoneModal = false;
  }

  const openModal = () => {
    showMicrophoneModal = !showMicrophoneModal;
  };
</script>

<div class="relative inline-block">
  <div class="flex rounded-full bg-white pl-1">
    <button
      onclick={openModal}
      class="my-auto flex select-none items-center justify-center rounded-full p-1 no-underline hover:brightness-75 dark:text-black"
    >
      <span class="material-symbols-outlined">
        {showMicrophoneModal ? 'arrow_drop_up' : 'arrow_drop_down'}
      </span>
    </button>
    <button
      onclick={closeMic}
      class={`my-auto flex select-none items-center justify-center rounded-full p-3 no-underline hover:brightness-75 dark:text-black ${isMicOn ? 'bg-gray-200 dark:text-black' : 'bg-red-500'}`}
    >
      <span class="material-symbols-outlined"> mic </span>
    </button>
  </div>
  {#if showMicrophoneModal}
    <div
      class={`absolute left-1/2 ${displayTop ? 'bottom-16' : 'top-16'} z-10 w-max -translate-x-1/2 transform rounded-lg bg-gray-200 p-2 shadow-lg`}
    >
      <div
        class={`absolute ${displayTop ? '-bottom-2' : '-top-2'} left-1/2 z-10 h-4 w-4 -translate-x-1/2 rotate-45 transform bg-gray-200`}
      ></div>
      {#each microphoneList as microphone}
        <div class="relative z-20 mx-2 cursor-pointer select-none rounded p-1 hover:bg-gray-500">
          <button
            class="m-0 w-full text-left text-sm"
            onclick={() => setSelectedMicrophone(microphone)}
          >
            {microphone.label}
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
