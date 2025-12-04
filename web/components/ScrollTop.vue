<template>
  <client-only>
    <a v-show="visible" class="bottom-right" @click="scrollTop">
      <slot />
    </a>
  </client-only>
</template>

<script>
export default {
  data () {
    return {
      visible: false,
    }
  },
  mounted () {
    window.addEventListener('scroll', this.scrollListener, { passive: true })
  },
  beforeDestroy () {
    window.removeEventListener('scroll', this.scrollListener)
  },
  methods: {
    scrollTop () {
      this.intervalId = setInterval(() => {
        if (window.pageYOffset === 0)
          clearInterval(this.intervalId)

        window.scroll(0, window.pageYOffset - 50)
      }, 20)
    },
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    scrollListener (e) {
      this.visible = window.scrollY > 150
    },
  },
}
</script>

<style scoped>
.bottom-right {
  position: fixed;
  bottom: 20px;
  right: 20px;
  cursor: pointer;
}
</style>
