<template>
  <div ref="foo" />
</template>

<script>
// https://publish.twitter.com/oembed?url=https://twitter.com/jack/status/20&omit_script=1

export default {
  head () {
    return {
      script: [
        {
          type: 'text/javascript',
          src: '//platform.twitter.com/widgets.js',
          defer: false,
          async: true,
          callback: () => { this.loadTweets() },
        },
      ],
    }
  },
  props: {
    oembedHtml: {
      type: String,
    },
    tweetId: {
      type: String,
      required: false,
    },
  },
  data () {
    return {
      polling: null,
    }
  },
  methods: {
    pollTweet () {
      this.polling = setInterval(() => {
        console.log(`${this.tweetId} Tweet taking too long to load`)
        clearInterval(this.polling)
      }, 10000)
    },
    loadTweets (tweetId) {
      // this.pollTweet()
      twttr.ready((twttr) => {
        twttr.widgets.createTweet(this.tweetId, this.$refs.foo)
      })
    },
  },
}
</script>
