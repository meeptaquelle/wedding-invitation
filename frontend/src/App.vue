<script setup>
import { computed } from 'vue'
import TemplateA from './templates/TemplateA/index.vue'
import TemplateB from './templates/TemplateB/index.vue'
import TemplateC from './templates/TemplateC/index.vue'
import InvitationController from './app/InvitationController.vue'

// DESIGN EDIT:
// Change this value to select the active invitation template.
// Valid values: 'TemplateA' | 'TemplateB' | 'TemplateC'
const activeTemplate = 'TemplateB'

const templates = {
  TemplateA,
  TemplateB,
  TemplateC
}

const SelectedTemplate = computed(() => templates[activeTemplate])
</script>

<template>
  <InvitationController :template="SelectedTemplate" />
</template>
