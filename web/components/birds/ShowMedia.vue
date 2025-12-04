<template>
  <div class="flex -space-x-2 overflow-hidden">
    <button @click="updateCoords(media.geometry.coordinates)">
      <img class="inline-block h-8 w-8 rounded-full ring-2 ring-white" :src="media.properties.imageUrl" alt="Media agency logo">
    </button>
  </div>
</template>

<script lang="ts">

// eslint-disable-next-line import/named
import Vue, { PropType } from 'vue'
import { GeojsonFeature } from '~/types'

export default Vue.extend({
  props: {
    media: {
      type: Object as PropType<GeojsonFeature>,
      required: true,
    },
    index: {
      type: Number,
      required: true,
    },
  },
  computed: {
    zNum () {
      const zNum: number = this.index
      return 'z-' + zNum * 10
    },
  },
  methods: {
    updateCoords (coords: number[]) {
      // this.$store.commit('SET_COORDS', coords)
      this.$store.dispatch('updateCoords', coords)
    },
  },
})
</script>
