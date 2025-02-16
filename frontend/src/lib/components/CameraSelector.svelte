<script lang="ts">
  import { onMount } from 'svelte';

  let {
    showCameraModal = false,
    selectedCamera = $bindable('default'),
    isCamOn = true,
    displayTop = false,
    closeCamera = $bindable()
  }: {
    showCameraModal: boolean;
    selectedCamera: string;
    isCamOn: boolean;
    displayTop: boolean;
    closeCamera: () => void;
  } = $props();

  let cameraList: MediaDeviceInfo[] = $state([]);

  onMount(async () => {
    navigator.mediaDevices.enumerateDevices().then((devices) => {
      cameraList = devices.filter((device) => device.kind === 'videoinput');
    });
  });

  function setSelectedCamera(newCamera: MediaDeviceInfo) {
    selectedCamera = newCamera?.deviceId || 'default';
    showCameraModal = false;
  }

  const openModal = () => {
    showCameraModal = !showCameraModal;
  };
</script>

<div class="relative inline-block">
  <div class="flex rounded-full bg-white pl-1">
    <button
      onclick={openModal}
      class="my-auto flex select-none items-center justify-center rounded-full p-1 no-underline hover:brightness-75 dark:text-black"
    >
      <span class="material-symbols-outlined">
        {showCameraModal ? 'arrow_drop_up' : 'arrow_drop_down'}
      </span>
    </button>
    <button
      onclick={closeCamera}
      class={`my-auto flex select-none items-center justify-center rounded-full p-3 no-underline hover:brightness-75 dark:text-black ${isCamOn ? 'bg-gray-200 dark:text-black' : 'bg-red-500'}`}
    >
      <span class="material-symbols-outlined"> videocam </span>
    </button>
  </div>
  {#if showCameraModal}

    <div
      class={`absolute left-1/2 ${displayTop ? 'bottom-16' : 'top-16'} z-10 w-max -translate-x-1/2 transform rounded-lg bg-gray-200 p-2 shadow-lg`}
    >
      <div
        class={`absolute ${displayTop ? '-bottom-2' : '-top-2'} left-1/2 z-10 h-4 w-4 -translate-x-1/2 rotate-45 transform bg-gray-200`}
      ></div>
      {#each cameraList as camera}
        <div class="relative z-20 mx-2 cursor-pointer select-none rounded p-1 hover:bg-gray-500">
          <button class="m-0 w-full text-left text-sm" onclick={() => setSelectedCamera(camera)}>
            {camera.label}
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
