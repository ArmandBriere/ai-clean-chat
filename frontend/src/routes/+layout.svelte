<script>
	import '@/app.css';
	const months = [
		'January',
		'February',
		'March',
		'April',
		'May',
		'June',
		'July',
		'August',
		'September',
		'October',
		'November',
		'December'
	];
	const date = new Date();
	const year = date.getFullYear().toString();
	const month = months[date.getMonth()];
	const day = date.getDate().toString();
</script>

<div class="flex h-svh flex-col">
	<header>
		<!-- wrapper -->
		<div class="flex w-full justify-between p-8">
			<!-- left nav -->
			<div>Logo Oseadea Confoo</div>
			<!-- right nav -->
			<div>Sale gueule de PD</div>
		</div>
	</header>

	<div class="flex flex-grow">
		<slot />
	</div>
	<footer>
		<div class="flex justify-between p-8">
			<div>{month} {day} {year}</div>
			<div>Montreal, Canada</div>
		</div>
	</footer>
</div>
