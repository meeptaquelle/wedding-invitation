<template>
  <div class="state state-loading" role="status" aria-live="polite">
    <div class="spinner" aria-hidden="true"></div>
    <p>Loading your invitation...</p>
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
  color: var(--color-text);
}
.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
