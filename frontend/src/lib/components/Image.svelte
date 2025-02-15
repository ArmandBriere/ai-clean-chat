<script lang="ts">
  let {
    class: className = '',
    wrapperClasses = '',
    src,
    alt,
    srcLarge = undefined,
    srcMedium = undefined,
    loading = 'lazy',
    fetchpriority = 'auto',
    width = undefined,
    height = undefined
  }: {
    class: string;
    wrapperClasses: string;
    src: string;
    alt: string;
    srcLarge: string | undefined;
    srcMedium: string | undefined;
    loading: 'eager' | 'lazy' | null | undefined;
    fetchpriority: 'auto' | 'high' | 'low' | null | undefined;
    width: string | undefined;
    height: string | undefined;
  } = $props();
</script>

<picture class={wrapperClasses}>
  <!-- Fallback image for browsers that don't support the <picture> element -->
  <source srcset={srcLarge ?? src} media="(min-width: 800px)" />
  <source srcset={srcMedium ?? src} media="(min-width: 500px)" />
  <img {src} {alt} {loading} {fetchpriority} class={className} {width} {height} />
</picture>
