<script setup>
import { ref } from 'vue'

const props = defineProps({
  messages: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  nextCursor: { type: String, default: null },
  submitting: { type: Boolean, default: false },
  submitError: { type: String, default: null },
  guestName: { type: String, default: 'Fulan' }
})

const emit = defineEmits(['submit', 'load-more'])

const message = ref('')
const localError = ref('')
const MAX = 500

function onSubmit() {
  const trimmed = message.value.trim()
  if (!trimmed) {
    localError.value = 'Please enter a message.'
    return
  }
  if (trimmed.length > MAX) {
    localError.value = `Message must be ${MAX} characters or fewer.`
    return
  }
  localError.value = ''
  emit('submit', trimmed)
  message.value = ''
}

function formatDate(iso) {
  try {
    return new Date(iso).toLocaleString(undefined, {
      dateStyle: 'medium',
      timeStyle: 'short'
    })
  } catch {
    return iso
  }
}
</script>

<template>
  <div class="guestbook">
    <form class="guestbook-form" @submit.prevent="onSubmit">
      <label class="guestbook-label" for="gb-msg">
        Leave a message as <strong>{{ guestName }}</strong>
      </label>
      <textarea
        id="gb-msg"
        v-model="message"
        :maxlength="MAX"
        rows="3"
        placeholder="Write your wishes..."
        :disabled="submitting"
      ></textarea>
      <div class="guestbook-form-footer">
        <span class="guestbook-count">{{ message.length }} / {{ MAX }}</span>
        <button type="submit" :disabled="submitting || !message.trim()">
          {{ submitting ? 'Sending...' : 'Send' }}
        </button>
      </div>
      <p v-if="localError" class="guestbook-error">{{ localError }}</p>
      <p v-if="submitError" class="guestbook-error">{{ submitError }}</p>
    </form>

    <div class="guestbook-list">
      <p v-if="loading && !messages.length" class="guestbook-empty">Loading messages…</p>
      <p v-else-if="!messages.length" class="guestbook-empty">
        No messages yet — be the first!
      </p>
      <article v-for="m in messages" :key="m.message_id" class="guestbook-item">
        <header class="guestbook-item-head">
          <strong>{{ m.display_name }}</strong>
          <time>{{ formatDate(m.created_at) }}</time>
        </header>
        <p>{{ m.message }}</p>
      </article>
    </div>

    <div v-if="nextCursor" class="guestbook-more">
      <button type="button" :disabled="loading" @click="emit('load-more')">
        {{ loading ? 'Loading…' : 'Load more' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.guestbook {
  max-width: 640px;
  margin: 0 auto;
  width: 100%;
  text-align: left;
}
.guestbook-form {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  margin-bottom: 2rem;
}
.guestbook-label {
  font-size: 0.9rem;
  opacity: 0.9;
}
textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 4px;
  font-family: inherit;
  font-size: 1rem;
  resize: vertical;
  background: rgba(255, 255, 255, 0.9);
  color: inherit;
}
.guestbook-form-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}
.guestbook-count {
  font-size: 0.75rem;
  opacity: 0.55;
}
button {
  padding: 0.55rem 1.4rem;
  background: var(--color-primary);
  color: #fff;
  border: none;
  border-radius: 4px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-size: 0.78rem;
}
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.guestbook-error {
  color: #b00020;
  font-size: 0.85rem;
  margin: 0;
}
.guestbook-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.guestbook-item {
  padding: 1rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}
.guestbook-item-head {
  display: flex;
  justify-content: space-between;
  gap: 0.75rem;
  flex-wrap: wrap;
  margin-bottom: 0.5rem;
  font-size: 0.85rem;
}
.guestbook-item-head time {
  opacity: 0.55;
  font-size: 0.75rem;
}
.guestbook-item p {
  margin: 0;
  line-height: 1.6;
  word-break: break-word;
}
.guestbook-empty {
  text-align: center;
  opacity: 0.6;
  padding: 2rem 0;
}
.guestbook-more {
  text-align: center;
  margin-top: 1.5rem;
}
</style>
