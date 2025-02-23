<script lang="ts">
  import { PUBLIC_SERVER_URL } from '$env/static/public';
  import type { MouseEventHandler } from 'svelte/elements';
  import Toast from './Toast.svelte';
  let {
    showInfoPanel = false,
    roomID,
    analyses = [],
    handleClose
  }: {
    showInfoPanel: boolean;
    roomID: string;
    analyses: any[];
    handleClose: MouseEventHandler<HTMLButtonElement>;
  } = $props();

  let meetingUrl = `${PUBLIC_SERVER_URL}/chat/${roomID}`;

  let isCopied = $state(false);

  const shareMeetingURI = () => {
    const copyText = `${PUBLIC_SERVER_URL}/chat/${roomID}`;
    const theClipboard = navigator.clipboard;
    theClipboard.writeText(copyText).then(() => console.log('copied to clipboard'));
    isCopied = true;
    setTimeout(() => {
      isCopied = false;
    }, 2000);
  };
</script>

<div
  class={`absolute bottom-[80px] ${showInfoPanel ? 'right-4 opacity-100' : '-right-[300px] opacity-0'} top-4 flex w-[300px] transform rounded-md bg-gray-200 transition-all`}
>
  <div class="relative h-full w-full">
    <div class="flex h-full flex-col">
      <!-- header -->
      <div class="flex h-16 items-center px-4">
        <div class="flex-grow overflow-hidden text-ellipsis whitespace-nowrap">
          <h2 class="text-xl font-light">Meeting details</h2>
        </div>
        <div class="flex h-full items-center justify-end"
          ><button onclick={handleClose} class="inline-flex items-center justify-center font-light">
            <span class="material-symbols-outlined"> close </span>
          </button>
        </div>
      </div>
      <div class="flex flex-col space-y-2 px-4">
        <h3>Joining info</h3>
        <p class="text-sm font-light">{meetingUrl}</p>
        <button
          onclick={shareMeetingURI}
          class="my-1 inline-flex w-fit rounded-full px-3 py-2 align-middle text-blue-600 hover:bg-blue-200 hover:bg-opacity-50"
          ><span class="material-symbols-outlined mr-2">content_copy</span>Copy joining info</button
        >
      </div>
      <div class="my-2 border-t border-gray-300"></div>
      <!-- content -->
      <div class="flex-grow overflow-scroll px-4">
        <!-- TODO: add analysis from LLM -->
        <!-- NOTE: need to manage maximum message list with First in first out -->
        {#each analyses as analysis}
          <div class="mb-3 flex flex-col">
            <div class="text-sm font-light text-gray-500">At timestamp</div>
            <div>keyword</div>
            <div class="text-sm font-light">description</div>
          </div>
        {/each}
      </div>
    </div>
  </div>
</div>
{#if isCopied}
  <Toast content="Copied meeting link." />
{/if}
