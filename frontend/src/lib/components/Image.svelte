<script lang="ts">
  let className = '';
  export { className as class };
  export let wrapperClasses = '';
  export let src: string;
  export let alt: string;
  export let srcLarge: string | undefined = undefined;
  export let srcMedium: string | undefined = undefined;
  export let loading: 'eager' | 'lazy' | null | undefined = 'lazy';
  export let fetchpriority: 'auto' | 'high' | 'low' | null | undefined = 'auto';
  export let width: string | undefined = undefined;
  export let height: string | undefined = undefined;
</script>

<picture class={wrapperClasses}>
  <!-- Fallback image for browsers that don't support the <picture> element -->
  <source srcset={srcLarge ?? src} media="(min-width: 800px)" />
  <source srcset={srcMedium ?? src} media="(min-width: 500px)" />
  <img {src} {alt} {loading} {fetchpriority} class={className} {width} {height} />
</picture>
