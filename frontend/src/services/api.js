// API service layer. All backend communication flows through here.
const API_BASE = import.meta.env.VITE_API_BASE || '/api'

async function handle(res) {
  let data = null
  try {
    data = await res.json()
  } catch {
    /* no body */
  }
  if (!res.ok) {
    const message = (data && data.error) || `Request failed (${res.status})`
    const err = new Error(message)
    err.status = res.status
    throw err
  }
  return data
}

export async function resolveGuest(slug) {
  const res = await fetch(`${API_BASE}/guests/${encodeURIComponent(slug)}`)
  if (res.status === 404) return null
  return handle(res)
}

export async function fetchMessages({ cursor = null, limit = 20 } = {}) {
  const params = new URLSearchParams()
  if (cursor) params.set('cursor', cursor)
  params.set('limit', String(limit))
  const res = await fetch(`${API_BASE}/messages?${params.toString()}`)
  return handle(res)
}

export async function submitMessage({ message, slug = null }) {
  const url = slug
    ? `${API_BASE}/guests/${encodeURIComponent(slug)}/messages`
    : `${API_BASE}/messages`
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message })
  })
  return handle(res)
}
