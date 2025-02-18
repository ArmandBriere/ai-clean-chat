<script lang="ts">
  import { onMount } from 'svelte';

  let showTooltip = $state(false);
  let {
    displayText = '',
    children
  }: {
    displayText: string;
    children: Snippet<[]>;
  } = $props();
</script>

<div class="group relative inline-block" role="tooltip">
  <div
    class="pointer-events-none absolute -top-16 left-1/2 z-10 w-max -translate-x-1/2 transform select-none rounded-lg bg-gray-200 p-2 opacity-0 shadow-lg transition-opacity delay-500 group-hover:opacity-100"
  >
    {displayText}
  </div>
  {@render children()}
</div>
