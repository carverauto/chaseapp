<template>
  <div class="w-48">
    <div class="flex items-center justify-between">
      <div>
        <div v-if="completedSteps === 'width:0%'">
          <img class="h-7 w-7" src="/emoji-crying.png" alt="Crying emoji">
        </div>
        <div v-else>
          <img class="h-8 w-8" src="/emoji-sad.png" alt="Sad emoji">
        </div>
      </div>
      <div>
        <img class="h-8 w-8" src="/emoji-happy.png" alt="Happy emoji">
      </div>
    </div>
    <div class="grid grid-cols-2 h-3 relative max-w-48 rounded-full overflow-hidden">
      <div class="w-full h-full bg-gray-200 absolute" />
      <div class="h-full bg-green-500 absolute" :style="completedSteps" />
    </div>
  </div>
</template>

<script lang="ts">
// eslint-disable-next-line import/named
import Vue, { PropType } from 'vue'
import { Sentiment } from '~/types'

export default Vue.extend({
  props: {
    sentiment: {
      type: Object as PropType<Sentiment>,
      required: true,
    },
  },
  data () {
    return {
      totalSteps: 100,
      completedSteps: 'width:' + Math.abs(Math.round(this.sentiment?.score * 100 / 2)) + '%',
    }
  },
})
</script>

<style>
</style>
