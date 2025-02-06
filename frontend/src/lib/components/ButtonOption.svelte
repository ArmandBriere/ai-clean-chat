<script lang="ts">
  let className = '';
  export { className as class };
  export let isActive = false;
</script>

<button
  onclick={() => (isActive = !isActive)}
  class={`${className} my-auto flex items-center justify-center rounded-full ${isActive ? 'bg-red-500' : 'bg-gray-200 dark:text-black'}`}
>
  <span class="material-symbols-outlined"> <slot /> </span>
</button>
