<script lang="ts">
  import { onMount } from 'svelte';

  let {
    showModal = false,
    selectedDevice = $bindable('default'),
    isDeviceOn = true,
    displayTop = false,
    closeDevice = $bindable(),
    changeMediaSource = $bindable(),
    kind = 'audioinput'
  }: {
    showModal: boolean;
    selectedDevice: string;
    isDeviceOn: boolean;
    displayTop: boolean;
    closeDevice: () => void;
    changeMediaSource: (deviceID: string) => void;
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
    if (selectedDevice != newDevice?.deviceId) {
      selectedDevice = newDevice?.deviceId || 'default';
      changeMediaSource(selectedDevice);
    }
    showModal = false;
  }

  const openModal = () => {
    showModal = !showModal;
  };
</script>

<div class="relative inline-block">
  <div class="flex rounded-full bg-[#262729] pl-2">
    <button
      onclick={openModal}
      class="flex items-center justify-center rounded-full p-1 text-gray-200"
    >
      <span class="material-symbols-outlined font-light">
        {showModal ? 'keyboard_arrow_up' : 'keyboard_arrow_down'}
      </span>
    </button>
    <button
      onclick={closeDevice}
      class={`flex items-center justify-center rounded-full p-3 transition-colors ${!isDeviceOn ? 'bg-red-500 text-black' : 'bg-[#333537] text-gray-200 hover:bg-[#444746]'}`}
    >
      <span class="material-symbols-outlined"> {icon} </span>
    </button>
  </div>
  {#if showModal}
    <div
      class={`absolute left-1/2 ${displayTop ? 'bottom-16' : 'top-16'} z-10 flex w-max -translate-x-1/2 transform flex-col rounded-lg bg-gray-200 p-2 shadow-lg`}
    >
      <div
        class={`absolute ${displayTop ? '-bottom-2' : '-top-2'} left-1/2 z-10 h-4 w-4 -translate-x-1/2 rotate-45 transform bg-gray-200`}
      ></div>
      {#each deviceList as device}
        <button
          class="relative z-20 mx-2 cursor-pointer select-none rounded p-1 hover:bg-gray-500"
          onclick={() => setSelectedDevice(device)}
        >
          <div class="m-0 text-left text-sm">
            {device.label}
          </div>
        </button>
      {/each}
    </div>
  {/if}
</div>
