<template>
  <div v-if="youtube">
    <div class="relative">
      <LazyYoutubeVideo
        ref="youtube"
        :src="youtube"
        enablejsapi
        inject-player-script
        :autoplay="live"
        class="lg:mb-4"
      />
    </div>
  </div>
  <div v-else>
    <div v-if="mp4">
    <IconLink :url="mp4" />
    <!-- <vue-plyr>
      <video ref="videoStreaming" controls crossorigin playsinline>
        <source src="">
      </video>
    </vue-plyr>  -->
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import 'vue-lazy-youtube-video/dist/style.css'
import LazyYoutubeVideo from 'vue-lazy-youtube-video'
import {connect} from "getstream";

export default Vue.extend({
  components: {
    LazyYoutubeVideo,
  },
  props: {
    url: {
      type: String,
      required: true,
    },
    id: {
      type: String,
      required: true,
    },
    live: {
      type: Boolean,
      required: true,
    },
    id: {
      type: String,
      required: true,
    },
  },
  data () {
    return {
      canvas: null,
      youtube: '',
      mp4: '',
      loading: true,
      showTheater: true,
      /*
      animation: {
        actor: String,
        animstate: String,
        animtype: String,
        createdAt: Date,
        endpoint: String,
        label: Number,
        id: String,
        time: Date,
        verb: String,
      },
       */
      rive: '',
    }
  },
  mounted () {
    this.openYoutube(this.url)
    // this.getStreamFeeds()
  },
  methods: {
    openYoutube (val: string) {
      // Grab text
      const matches = /(?:http:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=)?([^& \n]+)/g.exec(val)
      if (!matches)
        return true

      this.youtube = 'https://www.youtube.com/embed/' + matches[1]
    },
  },
})
</script>

<style scoped>
.noHover{
  pointer-events: none;
}

</style>
