<template>
  <div>
    <div v-if="showTweet" v-on-clickaway="showTweets">
      <Tweet :tweet-id="eventProp.payload.id" />
    </div>
    <div>
      <ul class="-mb-8">
        <li>
          <div class="relative pb-8">
            <span class="absolute top-4 left-4 -ml-px h-full bg-gray-200" aria-hidden="true" />
            <div class="flex items-center justify-between">
              <div v-if="eventProp.payload.image_url">
                <img class="inline-block h-8 w-8 rounded-full" :src="eventProp.payload.image_url" alt="Twitter Profile Picture">
              </div>
              <div v-else>
                <img class="inline-block h-8 w-8 rounded-full mr-1" src="twitter.png" alt="Twitter Profile Picture">
              </div>
              <div class="hidden lg:block ml-3">
                <p class="text-sm font-sans font-medium text-gray-800 group-hover:text-gray-900">{{ eventProp.payload.name }}</p>
                <p class="text-xs font-sans font-medium text-gray-500 group-hover:text-gray-700">@{{ eventProp.payload.username }}</p>

                <p class="text-xs font-sans font-medium text-gray-400 group-hover:text-gray-700">{{ $dayjs(eventProp.created_at).fromNow() }}</p>
              </div>
              <div class="min-w-0 flex-1 pt-1.5 justify-between space-x-4">
                <div>
                  <p class="font-medium font-sans text-gray-900 sm:truncate ml-2" @click="seeTweet">{{ eventProp.payload.text }}</p>
                  <div class="absolute top-0 right-0">
                    <img class="h-4 w-4" src="Twitter_logo_blue_32.png">
                  </div>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">

import Vue from 'vue'

export default Vue.extend({
  props: {
    oembedHtml: {
      type: String,
    },
    eventProp: {
      type: Object,
      required: false,
    },
  },
  data () {
    return {
      showTweet: false,
      tweetId: '',
    }
  },
  methods: {
    seeTweet () {
      this.showTweet = true
    },
    showTweets () {
      this.showTweet = !this.showTweet
    },
  },
})
</script>

<style>
.font-sans {
  font-family: "Helvetica Neue"
}
</style>