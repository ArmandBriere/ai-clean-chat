<script lang="ts">
	import Icon from '@iconify/svelte';
</script>

<div class="grid-[1fr] m-auto grid w-full max-w-xl justify-center text-center">
	<h1 class="mb-8 text-2xl">AI-Powered Clean Talk:<br />Advanced Profanity Detection Patterns</h1>
	<div class="mb-8">
		<p>
			Advanced profanity detection goes beyond simple keyword searches. This session will cover data
			collection and training of a profanity detection model, and its integration with live
			Speech-to-Text AI technologies. We will also look at the challenges around handling live data.
		</p>
	</div>
	<div class="flex w-full justify-center space-x-4">
		<a href="/ws" class="flex items-center rounded bg-blue-600 px-5 py-3 text-white"
			><Icon
				icon="material-symbols-light:video-call-outline-rounded"
				width="1.5rem"
				height="1.5rem"
				class="mr-2"
			/>New meeting</a
		>
		<button class="flex items-center rounded border px-5 py-3 text-gray-600"
			><Icon
				icon="material-symbols-light:keyboard"
				width="1.5rem"
				height="1.5rem"
				class="mr-2"
			/>Enter code</button
		>
		<button class="rounded-sm px-4 py-2 text-gray-600">Join</button>
	</div>
</div>
