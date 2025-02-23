<script lang="ts">
  let { showInfoPanel = false }: { showInfoPanel: boolean } = $props();
</script>

<div
  class={`absolute bottom-[80px] ${showInfoPanel ? 'right-4 opacity-100' : '-right-[300px] opacity-0'} top-4 flex w-[300px] transform rounded-md bg-gray-200 transition-all`}
></div>
