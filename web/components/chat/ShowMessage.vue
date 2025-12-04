<template>
  <div>
    <div v-if="message.type !== 'join'">
      <div class="flex-shrink-0 group block mb-3">
        <div class="flex items-center">
          <div v-if="message.user.image">
            <span class="inline-block relative">
              <img class="inline-block h-6 w-6 rounded-full mr-1" :src="message.user.image" alt="user avatar">
              <span v-if="message.user.online" class="absolute top-0 right-0 block h-2 w-2 rounded-full ring-2 ring-white bg-green-400"></span>
            </span>
          </div>
          <div v-else>
            <client-only>
              <identicon :value="message.user.id" class="h-6 w-6 rounded-full mr-1" />
              <span v-if="message.user.online" class="absolute top-0 right-0 block h-2 w-2 rounded-full ring-2 ring-white bg-green-400"></span>
            </client-only>
          </div>
          <span class="text-xs text-gray-600 mr-1.5">
            {{ message.user.name ? message.user.name : userData.userName }}
          </span>
        </div>
        <div class="mb-2">
          <span class="inline-block text-xs text-gray-500">
            <time class="whitespace-nowrap" :datetime="$dayjs(message.created_at).tz('America/Los_Angeles').format('h:mm A z')"><i class="far fa-clock" /> {{ $dayjs(message.created_at).tz('America/Los_Angeles').format('h:mm A z') }}</time>
          </span>
          <span class="inline-block relative text-xs ml-1.5 mt-1.5">
            <span class="px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
              <client-only>
                <span class="inline-block relative">
                  <!-- <button v-long-press="300" @long-press-start="onLongPressStart" @long-press-stop="onLongPressStop"> -->
                    <div v-html="$md.render(message.text)" />
                  <!-- </button>
                  <span v-if="message.latest_reactions">
                    <span class="inline-flex absolute bottom-0">
                      <span v-for="(value, name) in message.latest_reactions">
                        <span class="inline-flex inline-block relative">
                          <button @click="addReaction(value.type, message.id)" class="block">
                            {{ reactions[value.type] }}
                          </button>
                          <span v-if="value.score > 1">
                           <span class="block h-2.5 w-2.5 text-xs opacity-75 mr-1">25</span>
                          </span>
                        </span>
                      </span>
                    </span>
                  </span>
                  -->
                </span>
                <!-- TODO: This should be fixed so it toggles based on viewpoints -->
                <!--
                <span v-if="reactionNext" v-on-clickaway="reactionClickAway" v-bind:class="[message.text.length >= 20 ? myClassLeft : myClass, myClass]">
                  <span class="absolute inset-y-4 -left-4 block h-1.5 w-1.5 transform -translate-y-1/2 translate-x-1/2 rounded-full ring-2 ring-white bg-gray-300"></span>
                  <span class="absolute inset-y-4 -left-2 block h-2.5 w-2.5 transform -translate-y-1/2 translate-x-1/2 rounded-full ring-2 ring-white bg-gray-300"></span>
                  <span class="inline-flex border-bottom-0">
                    <div v-for="(value, name) in reactions">
                      <button @click="addReaction(name, message.id)" class="block m-1">
                        <p>{{ value }}</p>
                      </button>
                    </div>
                  </span>
                </span>
                -->
              </client-only>
            </span>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
// import LongPress from 'vue-directive-long-press'
const mila = require('markdown-it-link-attributes')

export default {
  props: {
    message: {
      type: Object,
      required: true,
    },
    userData: {
      type: Object,
      required: true,
    },
  },
  data () {
    return {
      reactions: {
        hot_take: "💥",
        alert: "🚨",
        like: "👍",
        dislike: "👎",
      },
      score: 0,
      askReaction: false,
      reactionNext: false,
      myClass: "absolute border-bottom-0 items-center px-3 py-0.5 rounded-full text-sm font-medium bg-green-100 text-green-800",
      myClassLeft: "absolute left-0 border-bottom-0 items-center px-3 py-0.5 rounded-full text-sm font-medium bg-green-100 text-green-800",
    }
  },
  mounted () {
    this.$md.use(mila, {
      attrs: {
        target: '_blank',
        rel: 'noopener',
      },
    })
  },
  methods: {
    addScores (score) {
      return this.score + score
    },
    addReaction (reaction, messageId) {
      console.log(`Reaction: ${reaction} and msgID: ${messageId}`)
      const completeReaction = {
        messageId: messageId,
        reaction,
      }
      this.$emit('reaction', completeReaction)
      this.reactionNext = false
    },
    reactionClickAway () {
      this.reactionNext = false
    },
    onLongPressStart () {
      this.askReaction = true
    },
    onLongPressStop () {
      this.reactionNext = true
    },
  },
}
</script>

<style scoped>

</style>
