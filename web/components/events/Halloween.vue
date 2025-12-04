<template>
  <ul class="-mb-8">
    <li>
      <div class="relative pb-8">
        <span class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true" />
        <div class="relative flex space-x-3">
          <div>
            <span class="h-8 w-8 rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white">
              <Lottie ref="halloween" :lottie-settings="lottieSettings" />
            </span>
          </div>
          <div class="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
            <div>
              <p class="text-sm font-bold text-gray-500">{{ event.payload.name }}</p>
              <div v-if="event.payload.urls.length > 0">
                <div v-for="url in event.payload.urls" :key="url">
                  <p class="text-sm text-gray-400 truncate"><a :href="url" target="_new">{{ url }}</a></p>
                </div>
              </div>
            </div>
            <div class="text-right text-sm whitespace-nowrap text-gray-500">
              {{ $dayjs(event.createdAt.toDate()).fromNow() }}
            </div>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
  props: {
    event: {
      type: Object,
      required: true,
    },
  },
  data () {
    return {
      lottieSettings: { path: '/json/36606-halloween.json', class: 'h-8 w-8' },
    }
  },
})
</script>
