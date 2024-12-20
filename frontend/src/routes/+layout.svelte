<script>
  import '@/app.css';
  import Image from '@/lib/components/Image.svelte';
  import logo1 from '$lib/assets/logo/osedea-logo.png';
  import logo2 from '$lib/assets/logo/X.png';
  import logo3 from '$lib/assets/logo/confoo-logo.png';
  import avatar from '$lib/assets/avatar.jpeg';
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
      <a href="/" class="flex items-center">
        <div class="flex items-center space-x-4">
          <Image src={logo1} alt="logo" class="w-full object-contain" />
          <Image src={logo2} alt="logo" class="w-full object-contain" />
          <Image src={logo3} alt="logo" class="w-full object-contain" />
        </div>
      </a>
      <!-- right nav -->
      <div>
        <div class="w-14 overflow-hidden rounded-full">
          <a href="/health">
            <Image src={avatar} alt="avatar" class="w-full object-contain" />
          </a>
        </div>
      </div>
    </div>
  </header>

  <div class="flex flex-grow">
    <slot />
  </div>
  <footer>
    <div class="flex justify-between p-8 font-light">
      <div>{month} {day} {year}</div>
      <div>Montreal, Canada</div>
    </div>
  </footer>
</div>
