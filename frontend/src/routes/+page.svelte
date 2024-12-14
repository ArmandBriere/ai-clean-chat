<script lang="ts">
	import Keyboard from 'svelte-material-icons/Keyboard.svelte';
	import VideoPlusOutline from 'svelte-material-icons/VideoPlusOutline.svelte';
</script>

<div class="flex flex-grow">
	<div class="grid-[1fr] m-auto grid w-full max-w-xl justify-center text-center">
		<h1 class="mb-8 text-2xl">AI-Powered Clean Talk:<br />Advanced Profanity Detection Patterns</h1>
		<div class="mb-8">
			<p>
				Advanced profanity detection goes beyond simple keyword searches. This session will cover
				data collection and training of a profanity detection model, and its integration with live
				Speech-to-Text AI technologies. We will also look at the challenges around handling live
				data.
			</p>
		</div>
		<div class="flex w-full justify-center space-x-4">
			<a href="/ws" class="flex items-center rounded bg-blue-600 px-4 py-2 text-white"
				><VideoPlusOutline size="1.5rem" />New meeting</a
			>
			<button class="flex items-center rounded border px-4 py-2"
				><Keyboard size="1.5rem" />Enter code</button
			>
			<button class="rounded-sm px-4 py-2">Join</button>
		</div>
	</div>
</div>
