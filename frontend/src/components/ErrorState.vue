<script setup>
defineProps({
  message: { type: String, default: 'Something went wrong.' }
})
defineEmits(['retry'])
</script>

<template>
  <div class="state state-error" role="alert">
    <h1>Oops!</h1>
    <p>{{ message }}</p>
    <button @click="$emit('retry')">Try Again</button>
  </div>
</template>

<style scoped>
.state {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 2rem;
  text-align: center;
}
h1 {
  font-size: clamp(2rem, 5vw, 3rem);
  margin: 0;
}
button {
  margin-top: 1rem;
  padding: 0.75rem 1.75rem;
  background: var(--color-primary);
  color: #fff;
  border: none;
  border-radius: 4px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  font-size: 0.85rem;
}
button:hover {
  opacity: 0.9;
}
</style>
