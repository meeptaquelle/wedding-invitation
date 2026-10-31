<script setup>
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import Guestbook from '../../components/Guestbook.vue'
import Countdown from '../../components/Countdown.vue'
import Gallery from '../../components/Gallery.vue'
import './style.css'

const props = defineProps({
  guestName: { type: String, default: 'Fulan' },
  weddingDetails: { type: Object, required: true },
  messages: { type: Array, default: () => [] },
  messagesLoading: { type: Boolean, default: false },
  nextCursor: { type: String, default: null },
  submitting: { type: Boolean, default: false },
  submitError: { type: String, default: null }
})

const emit = defineEmits(['submit-message', 'load-more'])

const c = computed(() => props.weddingDetails)

// Floating petals (pure CSS animation, positions are randomized once)
const petals = Array.from({ length: 14 }, (_, i) => ({
  id: i,
  left: Math.round(Math.random() * 100),
  delay: +(Math.random() * 12).toFixed(2),
  duration: +(10 + Math.random() * 10).toFixed(2),
  size: +(10 + Math.random() * 14).toFixed(0),
  drift: Math.round(-60 + Math.random() * 120)
}))

// Scroll reveal
const root = ref(null)
let observer = null

onMounted(() => {
  if (!root.value) return
  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((e) => {
        if (e.isIntersecting) {
          e.target.classList.add('tb-in')
          observer.unobserve(e.target)
        }
      })
    },
    { threshold: 0.15 }
  )
  root.value.querySelectorAll('.tb-reveal').forEach((el) => observer.observe(el))
})

onBeforeUnmount(() => observer && observer.disconnect())
</script>

<template>
  <div class="template-d" ref="root">
    <!-- Floating petals -->
    <div class="tb-petals" aria-hidden="true">
      <span
        v-for="p in petals"
        :key="p.id"
        class="tb-petal"
        :style="{
          left: p.left + '%',
          width: p.size + 'px',
          height: p.size * 1.3 + 'px',
          animationDelay: p.delay + 's',
          animationDuration: p.duration + 's',
          '--drift': p.drift + 'px'
        }"
      ></span>
    </div>

    <!-- Hero -->
    <section class="tb-hero">
      <div class="tb-hero-bg" :style="{ backgroundImage: `url(${c.images.hero})` }"></div>
      <div class="tb-hero-overlay"></div>

      <div class="tb-arch" aria-hidden="true">
        <svg viewBox="0 0 400 600" preserveAspectRatio="none">
          <path
            d="M10 590 V200 C10 90 100 10 200 10 C300 10 390 90 390 200 V590"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          />
          <path
            d="M24 590 V204 C24 100 108 26 200 26 C292 26 376 100 376 204 V590"
            fill="none"
            stroke="currentColor"
            stroke-width="0.6"
            stroke-dasharray="3 6"
          />
        </svg>
      </div>

      <div class="tb-hero-content">
        <p class="tb-eyebrow"><span></span>The Wedding Of<span></span></p>
        <h1 class="tb-names">
          <span class="tb-name tb-name-1">{{ c.couple.partnerOne }}</span>
          <span class="tb-amp">&amp;</span>
          <span class="tb-name tb-name-2">{{ c.couple.partnerTwo }}</span>
        </h1>
        <p class="tb-date">{{ c.event.date }}</p>

        <div class="tb-greeting">
          <p class="tb-greeting-label">Kepada Yth.</p>
          <p class="tb-guest">{{ guestName }}</p>
        </div>

        <div class="tb-scroll" aria-hidden="true">
          <span></span>
        </div>
      </div>
    </section>

    <!-- Divider -->
    <div class="tb-divider" aria-hidden="true">
      <svg viewBox="0 0 200 24" width="160" height="20">
        <path d="M0 12 H80" stroke="currentColor" stroke-width="1" />
        <path d="M120 12 H200" stroke="currentColor" stroke-width="1" />
        <path
          d="M100 3 C106 8 106 16 100 21 C94 16 94 8 100 3 Z M100 12 C108 6 116 6 120 12 C116 18 108 18 100 12 Z M100 12 C92 6 84 6 80 12 C84 18 92 18 100 12 Z"
          fill="currentColor"
        />
      </svg>
    </div>

    <!-- Introduction -->
    <section v-if="c.sections.introduction" class="tb-section tb-intro tb-reveal">
      <span class="tb-quote" aria-hidden="true">&ldquo;</span>
      <p>{{ c.introduction }}</p>
    </section>

    <!-- Event -->
    <section class="tb-section tb-event">
      <h2 class="tb-title tb-reveal">The Event</h2>
      <div class="tb-event-grid">
        <article class="tb-card tb-reveal" style="--d: 0s">
          <svg class="tb-ico" viewBox="0 0 24 24" aria-hidden="true">
            <rect x="3" y="5" width="18" height="16" rx="2" />
            <path d="M3 10h18M8 3v4M16 3v4" />
          </svg>
          <dt>Date</dt>
          <dd>{{ c.event.date }}</dd>
        </article>
        <article class="tb-card tb-reveal" style="--d: 0.12s">
          <svg class="tb-ico" viewBox="0 0 24 24" aria-hidden="true">
            <circle cx="12" cy="12" r="9" />
            <path d="M12 7v5l3 2" />
          </svg>
          <dt>Time</dt>
          <dd>{{ c.event.time }} ({{ c.event.timezone }})</dd>
        </article>
        <article class="tb-card tb-reveal" style="--d: 0.24s">
          <svg class="tb-ico" viewBox="0 0 24 24" aria-hidden="true">
            <path d="M3 21h18M5 21V9l7-5 7 5v12M10 21v-6h4v6" />
          </svg>
          <dt>Venue</dt>
          <dd>{{ c.event.venue }}</dd>
        </article>
        <article class="tb-card tb-reveal" style="--d: 0.36s">
          <svg class="tb-ico" viewBox="0 0 24 24" aria-hidden="true">
            <path d="M12 21s7-6.2 7-11a7 7 0 10-14 0c0 4.8 7 11 7 11z" />
            <circle cx="12" cy="10" r="2.5" />
          </svg>
          <dt>Address</dt>
          <dd>{{ c.event.address }}</dd>
        </article>
      </div>
      <a
        v-if="c.event.mapsUrl"
        :href="c.event.mapsUrl"
        target="_blank"
        rel="noopener"
        class="tb-btn tb-reveal"
      >
        <span>View on Maps</span>
      </a>
    </section>

    <!-- Countdown -->
    <section v-if="c.sections.countdown" class="tb-section tb-countdown">
      <div class="tb-countdown-inner tb-reveal">
        <h2 class="tb-title">Counting Down</h2>
        <Countdown :target-date="c.event.date" :time="c.event.time" :timezone="c.event.timezone" />
      </div>
    </section>

    <!-- Gallery -->
    <section v-if="c.sections.gallery" class="tb-section tb-gallery">
      <h2 class="tb-title tb-reveal">Moments</h2>
      <div class="tb-gallery-frame tb-reveal">
        <Gallery :images="c.images.gallery" />
      </div>
    </section>

    <!-- Guestbook -->
    <section v-if="c.sections.guestbook" class="tb-section tb-guestbook">
      <h2 class="tb-title tb-reveal">Guestbook</h2>
      <div class="tb-guestbook-paper tb-reveal">
        <Guestbook
          :messages="messages"
          :loading="messagesLoading"
          :next-cursor="nextCursor"
          :submitting="submitting"
          :submit-error="submitError"
          :guest-name="guestName"
          @submit="emit('submit-message', $event)"
          @load-more="emit('load-more')"
        />
      </div>
    </section>

    <footer class="tb-footer">
      <div class="tb-footer-wave" aria-hidden="true">
        <svg viewBox="0 0 1440 80" preserveAspectRatio="none">
          <path d="M0 80 V40 C240 0 480 0 720 30 C960 60 1200 60 1440 20 V80 Z" fill="currentColor" />
        </svg>
      </div>
      <p class="tb-footer-close">{{ c.footer.closing }}</p>
      <p class="tb-footer-sig">{{ c.footer.signature }}</p>
      <p class="tb-footer-credit">{{ c.footer.credit }}</p>
    </footer>
  </div>
</template>