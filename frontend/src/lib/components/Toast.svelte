<script lang="ts">
  let { content } = $props();
</script>

<div
  id="toast-bottom-left"
  class="fixed bottom-5 left-5 flex w-full max-w-xs items-center space-x-4 divide-x divide-gray-200 rounded-lg bg-white p-4 text-gray-500 shadow-sm rtl:divide-x-reverse dark:divide-gray-700 dark:bg-gray-800 dark:text-gray-400"
  role="alert"
>
  <div class="text-sm font-normal">{content}</div>
</div>
