<script lang="ts">
  import { onMount } from 'svelte';

  let {
    showModal = false,
    selectedDevice = $bindable('default'),
    isDeviceOn = true,
    displayTop = false,
    closeDevice = $bindable(),
    kind // 'audioinput' or 'videoinput'
  }: {
    showModal: boolean;
    selectedDevice: string;
    isDeviceOn: boolean;
    displayTop: boolean;
    closeDevice: () => void;
    kind: 'audioinput' | 'videoinput';
  } = $props();

  let deviceList: MediaDeviceInfo[] = $state([]);
  let icon: string = $state('');
  let deviceName: string = $state('');

  onMount(async () => {
    icon = kind === 'videoinput' ? 'videocam' : 'mic';
    deviceName = kind === 'videoinput' ? 'Camera' : 'Microphone';

    navigator.mediaDevices.enumerateDevices().then((devices) => {
      deviceList = devices.filter((device) => device.kind === kind);
    });
  });

  function setSelectedDevice(newDevice: MediaDeviceInfo) {
    selectedDevice = newDevice?.deviceId || 'default';
    showModal = false;
  }

  const openModal = () => {
    showModal = !showModal;
  };
</script>

<div class="relative inline-block">
  <div class="flex rounded-full bg-white pl-1">
    <button
      onclick={openModal}
      class="my-auto flex select-none items-center justify-center rounded-full p-1 no-underline hover:brightness-75 dark:text-black"
    >
      <span class="material-symbols-outlined">
        {showModal ? 'arrow_drop_up' : 'arrow_drop_down'}
      </span>
    </button>
    <button
      onclick={closeDevice}
      class={`my-auto flex select-none items-center justify-center rounded-full p-3 no-underline hover:brightness-75 dark:text-black ${isDeviceOn ? 'bg-gray-200 dark:text-black' : 'bg-red-500'}`}
    >
      <span class="material-symbols-outlined"> {icon} </span>
    </button>
  </div>
  {#if showModal}
    <div
      class={`absolute left-1/2 ${displayTop ? 'bottom-16' : 'top-16'} z-10 w-max -translate-x-1/2 transform rounded-lg bg-gray-200 p-2 shadow-lg`}
    >
      <div
        class={`absolute ${displayTop ? '-bottom-2' : '-top-2'} left-1/2 z-10 h-4 w-4 -translate-x-1/2 rotate-45 transform bg-gray-200`}
      ></div>
      {#each deviceList as device}
        <div class="relative z-20 mx-2 cursor-pointer select-none rounded p-1 hover:bg-gray-500">
          <button class="m-0 w-full text-left text-sm" onclick={() => setSelectedDevice(device)}>
            {device.label}
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
