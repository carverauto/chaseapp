<template>
  <a :href="link" :class="foo" @click="linkPopout" />
</template>

<script>
export default {
  props: {
    foo: {
      type: String,
      default: '',
    },
    link: {
      type: String,
      default: '',
    },
  },
  methods: {
    linkPopout () {
      window.open(this.link, '_blank')
    },
  },
}
</script>
