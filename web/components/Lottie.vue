<template>
  <div ref="lottieRef" :class="lottieSettings.class" />
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
  props: {
    lottieSettings: {
      type: Object,
      required: true,
    },
  },
  data () {
    return {
      animation: undefined as undefined|any,
    }
  },
  mounted () {
    this.lottie()
  },
  methods: {
    lottie (): void {
      // @ts-ignore
      if (!process.browser || !window.lottie) return
      // @ts-ignore
      const lottie = window.lottie
      const container = this.$refs.lottieRef as HTMLElement
      this.animation = lottie.loadAnimation({
        container,
        renderer: 'svg',
        path: this.lottieSettings.path,
        loop: true,
        autoplay: true,
      })
    },
  },
})
</script>
