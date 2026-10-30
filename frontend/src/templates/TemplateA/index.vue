<script setup>
import { computed } from 'vue'
import Guestbook from '../../components/Guestbook.vue'
import Countdown from '../../components/Countdown.vue'
import Gallery from '../../components/Gallery.vue'
import './style.css'

const props = defineProps({
  guestName: { type: String, default: 'Anonymous' },
  weddingDetails: { type: Object, required: true },
  messages: { type: Array, default: () => [] },
  messagesLoading: { type: Boolean, default: false },
  nextCursor: { type: String, default: null },
  submitting: { type: Boolean, default: false },
  submitError: { type: String, default: null }
})

const emit = defineEmits(['submit-message', 'load-more'])

const c = computed(() => props.weddingDetails)
</script>

<template>
  <div class="template-a">
    <!-- Hero -->
    <section class="ta-hero">
      <div class="ta-hero-bg" :style="{ backgroundImage: `url(${c.images.hero})` }"></div>
      <div class="ta-hero-overlay"></div>
      <div class="ta-hero-content">
        <p class="ta-eyebrow">The Wedding Of</p>
        <h1 class="ta-names">
          {{ c.couple.partnerOne }}
          <span class="ta-amp">&amp;</span>
          {{ c.couple.partnerTwo }}
        </h1>
        <p class="ta-date">{{ c.event.date }}</p>
        <div class="ta-greeting">
          <p>Kepada Yth.</p>
          <p class="ta-guest">{{ guestName }}</p>
        </div>
      </div>
    </section>

    <!-- Introduction -->
    <section v-if="c.sections.introduction" class="ta-section ta-intro">
      <p>{{ c.introduction }}</p>
    </section>

    <!-- Event -->
    <section class="ta-section ta-event">
      <h2 class="ta-title">The Event</h2>
      <dl class="ta-event-grid">
        <div>
          <dt>Date</dt>
          <dd>{{ c.event.date }}</dd>
        </div>
        <div>
          <dt>Time</dt>
          <dd>{{ c.event.time }} ({{ c.event.timezone }})</dd>
        </div>
        <div>
          <dt>Venue</dt>
          <dd>{{ c.event.venue }}</dd>
        </div>
        <div>
          <dt>Address</dt>
          <dd>{{ c.event.address }}</dd>
        </div>
      </dl>
      <a v-if="c.event.mapsUrl" :href="c.event.mapsUrl" target="_blank" rel="noopener" class="ta-btn">
        View on Maps
      </a>
    </section>

    <!-- Countdown -->
    <section v-if="c.sections.countdown" class="ta-section ta-countdown">
      <h2 class="ta-title">Counting Down</h2>
      <Countdown :target-date="c.event.date" :time="c.event.time" :timezone="c.event.timezone" />
    </section>

    <!-- Gallery -->
    <section v-if="c.sections.gallery" class="ta-section">
      <h2 class="ta-title">Moments</h2>
      <Gallery :images="c.images.gallery" />
    </section>

    <!-- Guestbook -->
    <section v-if="c.sections.guestbook" class="ta-section">
      <h2 class="ta-title">Guestbook</h2>
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
    </section>

    <footer class="ta-footer">
      <p class="ta-footer-close">{{ c.footer.closing }}</p>
      <p class="ta-footer-sig">{{ c.footer.signature }}</p>
      <p class="ta-footer-credit">{{ c.footer.credit }}</p>
    </footer>
  </div>
</template>
