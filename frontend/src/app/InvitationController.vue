<script setup>
import { useInvitation } from '../composables/useInvitation'
import LoadingState from '../components/LoadingState.vue'
import ErrorState from '../components/ErrorState.vue'

defineProps({
  template: { type: Object, required: true }
})

const {
  guestName,
  loading,
  error,
  weddingDetails,
  messages,
  messagesLoading,
  nextCursor,
  submitting,
  submitError,
  onSubmitMessage,
  loadMore,
  retry
} = useInvitation()
</script>

<template>
  <LoadingState v-if="loading" />
  <ErrorState v-else-if="error" :message="error" @retry="retry" />
  <component
    v-else
    :is="template"
    :guest-name="guestName"
    :wedding-details="weddingDetails"
    :messages="messages"
    :messages-loading="messagesLoading"
    :next-cursor="nextCursor"
    :submitting="submitting"
    :submit-error="submitError"
    @submit-message="onSubmitMessage"
    @load-more="loadMore"
  />
</template>
