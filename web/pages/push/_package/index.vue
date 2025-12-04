<template>
  <div v-if="package">
    <p>{{ package }}</p>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data () {
    return {
      package: this.$route.params.package,
    }
  },
  mounted () {
    console.log(this.package)
  },
  methods: {
  },
})
</script>
