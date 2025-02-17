<script lang="ts">
  import JSConfetti from 'js-confetti';
  import { onMount } from 'svelte';

  let {
    showEmojiModal = false,
    shareEmoji = () => {},
    receivedEmoji = ''
  }: {
    showEmojiModal: boolean;
    shareEmoji: (emoji: string) => void;
    receivedEmoji: string;
  } = $props();

  let confetti: JSConfetti | null = null;

  onMount(() => {
    confetti = new JSConfetti();
  });

  $effect(() => {
    if (confetti && receivedEmoji) {
      confetti.addConfetti({ emojis: [receivedEmoji] });
    }
  });

  const emojis = ['😃', '👍', '🤣', '❤️', '😭', '🦀'];

  const openModal = () => {
    showEmojiModal = !showEmojiModal;
  };

  const triggerConfetti = (emoji: string) => {
    if (confetti) {
      confetti.addConfetti({ emojis: [emoji] });
      shareEmoji(emoji);
    }
  };
</script>

<div class="relative inline-block">
  <button
    onclick={openModal}
    class={`flex items-center justify-center rounded-full p-3 transition-colors ${showEmojiModal ? 'bg-green-500 text-black' : 'bg-[#333537] text-gray-200 hover:bg-[#3f4245]'}`}
  >
    <span class="material-symbols-outlined"> mood </span>
  </button>
  {#if showEmojiModal}
    <div
      class="absolute -top-16 left-1/2 z-10 w-max -translate-x-1/2 transform rounded-lg bg-gray-200 p-2 shadow-lg"
    >
      <div class="absolute left-1/2 top-10 h-4 w-4 -translate-x-1/2 rotate-45 transform bg-gray-200"
      ></div>
      {#each emojis as emoji}
        <span class="emoji mx-5 cursor-pointer select-none rounded p-1 text-2xl hover:bg-gray-500">
          <button class="m-0 w-auto" onclick={() => triggerConfetti(emoji)}>
            {emoji}
          </button>
        </span>
      {/each}
    </div>
  {/if}
</div>
