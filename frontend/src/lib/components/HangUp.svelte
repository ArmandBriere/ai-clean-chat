<script lang="ts">
  let {
    handleHangUp,
    showHangUpModal = false,
    displayTop = false
  }: {
    handleHangUp: () => void;
    showHangUpModal: boolean;
    displayTop: boolean;
  } = $props();

  const openModal = () => {
    showHangUpModal = !showHangUpModal;
  };
</script>

<!-- Close meeting -->
<div class="relative inline-block">
  <button
    onclick={openModal}
    class="my-auto flex select-none items-center justify-center rounded-full bg-red-500 p-3 no-underline hover:opacity-70"
  >
    <span class="material-symbols-outlined"> call_end </span>
  </button>

  {#if showHangUpModal}
    <div
      class={`absolute left-1/2 ${displayTop ? 'bottom-16' : 'top-16'} z-10 w-40 -translate-x-1/2 transform rounded-lg bg-gray-200 p-3 shadow-lg`}
    >
      <div
        class={`absolute ${displayTop ? '-bottom-2' : '-top-2'} left-1/2 h-4 w-4 -translate-x-1/2 rotate-45 transform bg-gray-200`}
      ></div>
      <p class="mb-2 select-none text-center text-black no-underline">Leaving already?</p>

      <div class="flex justify-center space-x-2">
        <button
          class="w-full select-none rounded bg-green-500 px-3 py-1 text-black no-underline hover:bg-green-600"
          onclick={handleHangUp}>Yes</button
        >
        <button
          class="w-full select-none rounded bg-red-500 px-3 py-1 text-black no-underline hover:bg-red-600"
          onclick={() => (showHangUpModal = false)}>No</button
        >
      </div>
    </div>
  {/if}
</div>
