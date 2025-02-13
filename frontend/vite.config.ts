import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    allowedHosts: ['ai-clean-chat.home.site']
  }

  // test: {
  //   include: ['src/**/*.{test,spec}.{js,ts}']
  // }
});
