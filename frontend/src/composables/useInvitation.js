import { ref, onMounted } from 'vue'
import { resolveGuest, fetchMessages, submitMessage } from '../services/api'
import { weddingConfig } from '../config/wedding.config'

export function useInvitation() {
  const slug = ref(null)
  const guestName = ref('Anonymous')
  const loading = ref(true)
  const error = ref(null)

  const messages = ref([])
  const messagesLoading = ref(false)
  const nextCursor = ref(null)

  const submitting = ref(false)
  const submitError = ref(null)

  function extractSlug() {
    const path = window.location.pathname.replace(/^\/+|\/+$/g, '')
    if (!path || path === 'index.html') return null
    return path
  }

  async function loadMessages(cursor = null) {
    messagesLoading.value = true
    try {
      const result = await fetchMessages({ cursor })
      if (cursor) {
        messages.value = [...messages.value, ...result.messages]
      } else {
        messages.value = result.messages
      }
      nextCursor.value = result.next_cursor || null
    } catch (e) {
      submitError.value = e.message
    } finally {
      messagesLoading.value = false
    }
  }

  async function init() {
    loading.value = true
    error.value = null
    const s = extractSlug()
    slug.value = s

    if (s) {
      try {
        const guest = await resolveGuest(s)
        if (!guest) {
          error.value = 'Invalid invitation link. This guest could not be found.'
          loading.value = false
          return
        }
        guestName.value = guest.guest_name
      } catch (e) {
        error.value = e.message || 'Failed to load invitation.'
        loading.value = false
        return
      }
    }

    loading.value = false
    await loadMessages()
  }

  async function loadMore() {
    if (!nextCursor.value) return
    await loadMessages(nextCursor.value)
  }

  async function onSubmitMessage(message) {
    submitting.value = true
    submitError.value = null
    try {
      await submitMessage({ message, slug: slug.value })
      await loadMessages()
      return true
    } catch (e) {
      submitError.value = e.message || 'Failed to submit message.'
      return false
    } finally {
      submitting.value = false
    }
  }

  onMounted(init)

  return {
    slug,
    guestName,
    loading,
    error,
    weddingDetails: weddingConfig,
    messages,
    messagesLoading,
    nextCursor,
    submitting,
    submitError,
    onSubmitMessage,
    loadMore,
    retry: init
  }
}
