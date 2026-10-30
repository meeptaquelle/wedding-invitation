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
  <div class="template-c">
    <section class="tc-hero">
      <div class="tc-frame">
        <div class="tc-ornament">✦</div>
        <p class="tc-small">The Wedding Celebration of</p>
        <h1 class="tc-names">{{ c.couple.partnerOne }} &amp; {{ c.couple.partnerTwo }}</h1>
        <div class="tc-divider"><span></span>❦<span></span></div>
        <p class="tc-date">{{ c.event.date }}</p>
        <div class="tc-greeting">
          <p>To Our Beloved</p>
          <p class="tc-guest">{{ guestName }}</p>
        </div>
        <div class="tc-ornament">✦</div>
      </div>
    </section>

    <div class="tc-body">
      <section v-if="c.sections.introduction" class="tc-section">
        <p class="tc-intro">{{ c.introduction }}</p>
      </section>

      <section class="tc-section">
        <h2>The Details</h2>
        <div class="tc-details">
          <div>
            <p class="tc-label">When</p>
            <p class="tc-value">{{ c.event.date }} · {{ c.event.time }}</p>
          </div>
          <div>
            <p class="tc-label">Where</p>
            <p class="tc-value">{{ c.event.venue }}</p>
            <p class="tc-address">{{ c.event.address }}</p>
          </div>
        </div>
        <a v-if="c.event.mapsUrl" :href="c.event.mapsUrl" target="_blank" rel="noopener" class="tc-btn">
          Get Directions
        </a>
      </section>

      <section v-if="c.sections.countdown" class="tc-section tc-countdown">
        <h2>Time Remaining</h2>
        <Countdown :target-date="c.event.date" :time="c.event.time" :timezone="c.event.timezone" />
      </section>

      <section v-if="c.sections.gallery" class="tc-section">
        <h2>Gallery</h2>
        <Gallery :images="c.images.gallery" />
      </section>

      <section v-if="c.sections.guestbook" class="tc-section">
        <h2>Guestbook</h2>
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
    </div>

    <footer class="tc-footer">
      <div class="tc-divider"><span></span>❦<span></span></div>
      <p>{{ c.footer.closing }}</p>
      <p class="tc-sig">{{ c.footer.signature }}</p>
    </footer>
  </div>
</template>
