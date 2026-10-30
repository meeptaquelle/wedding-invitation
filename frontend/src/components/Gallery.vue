<script setup>
defineProps({
  images: { type: Array, default: () => [] }
})
</script>

<template>
  <div class="gallery">
    <p v-if="!images.length" class="gallery-empty">No photos yet.</p>
    <div v-else class="gallery-grid">
      <figure v-for="(src, i) in images" :key="i" class="gallery-item">
        <img :src="src" :alt="`Wedding photograph ${i + 1}`" loading="lazy" />
      </figure>
    </div>
  </div>
</template>

<style scoped>
.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 0.75rem;
}
.gallery-item {
  margin: 0;
  overflow: hidden;
  border-radius: 4px;
  aspect-ratio: 4 / 5;
  background: var(--color-secondary);
}
.gallery-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}
.gallery-item:hover img {
  transform: scale(1.04);
}
.gallery-empty {
  text-align: center;
  opacity: 0.6;
  padding: 1rem;
}
</style>
