<template>
  <div class="flex -space-x-2 overflow-hidden">
    <button @click="updateCoords(leo.geometry.coordinates)">
      <img class="inline-block h-8 w-8 rounded-full ring-2 ring-white" :src="leo.properties.imageUrl" alt="LEO agency logo">
    </button>
  </div>
</template>

<script lang="ts">

// eslint-disable-next-line import/named
import Vue, { PropType } from 'vue'
import { GeojsonFeature } from '~/types'

export default Vue.extend({
  props: {
    leo: {
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
      this.$store.commit('SET_COORDS', coords)
    },
  },
})
</script>
