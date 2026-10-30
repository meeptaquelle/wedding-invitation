<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  targetDate: { type: String, required: true }, // YYYY-MM-DD
  time: { type: String, default: '00:00' },     // HH:MM
  timezone: { type: String, default: 'UTC' }    // IANA tz
})

const now = ref(Date.now())
let timer = null

function offsetMs(tz, date) {
  const dtf = new Intl.DateTimeFormat('en-US', {
    timeZone: tz,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
  const parts = dtf.formatToParts(date)
  const get = (t) => Number(parts.find((p) => p.type === t)?.value || 0)
  const asUTC = Date.UTC(
    get('year'),
    get('month') - 1,
    get('day'),
    get('hour'),
    get('minute'),
    get('second')
  )
  return asUTC - date.getTime()
}

function parseInTimezone(dateStr, timeStr, tz) {
  const [y, m, d] = dateStr.split('-').map(Number)
  const [hh, mm] = timeStr.split(':').map(Number)
  let ts = Date.UTC(y, m - 1, d, hh, mm, 0)
  // Adjust: UTC + offset = local
  ts -= offsetMs(tz, new Date(ts))
  return ts
}

const targetTs = computed(() => parseInTimezone(props.targetDate, props.time, props.timezone))

const remaining = computed(() => {
  const diff = targetTs.value - now.value
  if (diff <= 0) return null
  const s = Math.floor(diff / 1000)
  return {
    days: Math.floor(s / 86400),
    hours: Math.floor((s % 86400) / 3600),
    minutes: Math.floor((s % 3600) / 60),
    seconds: s % 60
  }
})

onMounted(() => {
  timer = setInterval(() => {
    now.value = Date.now()
  }, 1000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="countdown">
    <p v-if="!remaining" class="countdown-complete">Today is the day 🤍</p>
    <div v-else class="countdown-grid">
      <div class="countdown-item">
        <span class="countdown-value">{{ remaining.days }}</span>
        <span class="countdown-label">Days</span>
      </div>
      <div class="countdown-item">
        <span class="countdown-value">{{ remaining.hours }}</span>
        <span class="countdown-label">Hours</span>
      </div>
      <div class="countdown-item">
        <span class="countdown-value">{{ remaining.minutes }}</span>
        <span class="countdown-label">Minutes</span>
      </div>
      <div class="countdown-item">
        <span class="countdown-value">{{ remaining.seconds }}</span>
        <span class="countdown-label">Seconds</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.countdown {
  text-align: center;
}
.countdown-grid {
  display: flex;
  justify-content: center;
  gap: clamp(0.75rem, 3vw, 2rem);
  flex-wrap: wrap;
}
.countdown-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 68px;
}
.countdown-value {
  font-family: var(--font-heading);
  font-size: clamp(2rem, 5.5vw, 3.25rem);
  line-height: 1;
}
.countdown-label {
  margin-top: 0.4rem;
  font-size: 0.7rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  opacity: 0.65;
}
.countdown-complete {
  font-family: var(--font-heading);
  font-size: 1.4rem;
}
</style>
