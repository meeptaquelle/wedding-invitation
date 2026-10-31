<script setup>
import { computed } from 'vue'
import Guestbook from '../../components/Guestbook.vue'
import Countdown from '../../components/Countdown.vue'
import Gallery from '../../components/Gallery.vue'
import './style.css'

const props = defineProps({
  guestName: { type: String, default: 'John Doe' },
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
  <div class="template-b">
    <div class="tb-border">
      <div class="tb-floral tb-floral-tl">❀</div>
      <div class="tb-floral tb-floral-tr">❀</div>

      <header class="tb-hero">
        <div class="tb-hero-image" :style="{ backgroundImage: `url(${c.images.hero})` }"></div>
        <div class="tb-hero-text">
          <p class="tb-small">Together with joy</p>
          <h1 class="tb-names">
            {{ c.couple.partnerOne }} <span>&amp;</span> {{ c.couple.partnerTwo }}
          </h1>
          <p class="tb-date">{{ c.event.date }}</p>
          <div class="tb-greet">
            <p>Dear</p>
            <p class="tb-guest">{{ guestName }}</p>
          </div>
        </div>
      </header>

      <main class="tb-main">
        <section v-if="c.sections.introduction" class="tb-intro">
          <p>{{ c.introduction }}</p>
        </section>

        <section class="tb-section">
          <h2>— The Celebration —</h2>
          <div class="tb-details">
            <p><span>{{ c.event.date }}</span> · <span>{{ c.event.time }}</span></p>
            <p class="tb-venue">{{ c.event.venue }}</p>
            <p class="tb-address">{{ c.event.address }}</p>
            <a v-if="c.event.mapsUrl" :href="c.event.mapsUrl" target="_blank" rel="noopener">Find the venue →</a>
          </div>
        </section>

        <section v-if="c.sections.countdown" class="tb-section tb-countdown-sec">
          <h2>— The Countdown —</h2>
          <Countdown :target-date="c.event.date" :time="c.event.time" :timezone="c.event.timezone" />
        </section>

        <section v-if="c.sections.gallery" class="tb-section">
          <h2>— Our Moments —</h2>
          <Gallery :images="c.images.gallery" />
        </section>

        <section v-if="c.sections.guestbook" class="tb-section">
          <h2>— Wish &amp; Blessings —</h2>
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
      </main>

      <footer class="tb-footer">
        <p>{{ c.footer.closing }}</p>
        <p class="tb-sig">{{ c.footer.signature }}</p>
      </footer>

      <div class="tb-floral tb-floral-bl">❀</div>
      <div class="tb-floral tb-floral-br">❀</div>
    </div>
  </div>
</template>
