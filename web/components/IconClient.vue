<template>
  <client-only>
    <Icon :icon="icon" :class="iconClass" />
    <template #fallback>
      <div :class="iconClass" />
    </template>
  </client-only>
</template>

<script lang="ts">
import { Icon } from '@iconify/vue2'
import Vue from 'vue'

export default Vue.extend({
  components: {
    Icon,
  },
  props: {
    icon: {
      type: String,
      required: true,
    },
    iconClass: {
      type: String,
      required: false,
    },
  },
})
</script>
