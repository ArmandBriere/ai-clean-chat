<script lang="ts">
  import { PUBLIC_SERVER_URL } from '$env/static/public';
  import Toast from './Toast.svelte';

  let {
    roomID,
    connectedUsers = 1
  }: {
    roomID: string;
    connectedUsers: number;
  } = $props();

  let showCopiedMessage = $state(false);

  const shareMeetingURI = () => {
    const copyText = `${PUBLIC_SERVER_URL}/chat/${roomID}`;
    const theClipboard = navigator.clipboard;
    theClipboard.writeText(copyText).then(() => console.log('copied to clipboard'));
  };
</script>

{#if showCopiedMessage}
  <Toast content="Copied meeting link." />
{/if}
<div class="flex-ro mx-1 flex items-center justify-center gap-1">
  <button
    title="Copy room id"
    onclick={() => {
      shareMeetingURI();
      showCopiedMessage = true;
      setTimeout(() => {
        showCopiedMessage = false;
      }, 1500);
    }}
    class="flex items-center justify-center gap-1 rounded-full p-2 text-white transition-colors hover:bg-gray-200/10"
    ><span class="material-symbols-outlined">content_copy</span>
    {roomID}
  </button>
  <div class="text-white">
    ({connectedUsers}
    {connectedUsers === 1 ? 'user' : 'users'})
  </div>
</div>
