<template>
  <div>
    <span class="text-left lg:text-2xl sm:text-3xl mb-4 mr-2 text-gray-900 sm:text-2xl">
        FireHose <span class="mb-0.5 text-red-800 sm:text-sm">Beta</span>
    </span>
    <div class="m-4">
      <div class="flex items-center justify-between">
        <div>
          <Tracker />
        </div>
      </div>
    </div>
    <div v-if="!loading">
      <div v-if="firehoseData.length > 0" class="border overflow-y-scroll h-48 max-w-3xl px-4 py-4 sm:px-6 lg:rounded-lg z-20 lg:shadow dark:bg-gray-800 bg-white lg:max-w-4xl lg:mb-4">
        <div class="flow-root">
          <ul class="-mb-8">
            <div v-if="firehoseData.length > 0" class="flex-col truncate items-center justify-center font-sans lg:space-y-4 divide-y divide-gray-300 lg:divide-y-0">
              <ShowEvent v-for="(event, $index) in orderedFirehoseData" :key="$index" :myevent="event" />
            </div>
          </ul>
        </div>
      </div>
    </div>
    <div v-else>
      <ListEventsSkeleton />
    </div>
  </div>
</template>

<script lang="ts">

import Vue from 'vue'
import {connect} from 'getstream'
// eslint-disable-next-line @typescript-eslint/no-var-requires,@typescript-eslint/no-unsafe-assignment
const orderBy = require('lodash.orderby')

export default Vue.extend({
  data() {
    return {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      firehoseData: [] as any[]| any,
      token: undefined as undefined | string,
      feedUser: 'firehose',
      loading: true,
    }
  },
  mounted() {
    void this.getStreamFeeds(this.feedUser)
  },
  computed: {
    // eslint-disable-next-line vue/return-in-computed-property
    orderedFirehoseData: function (): any {
      if (this.firehoseData)
        // eslint-disable-next-line @typescript-eslint/no-unsafe-return,@typescript-eslint/no-unsafe-call
        return orderBy(this.firehoseData, 'created_at', 'desc')
    },
  },
  methods: {
    async getStreamFeeds(userId: string) {
      try {
        await this.$axios.post('https://us-central1-chaseapp-8459b.cloudfunctions.net/GetStreamToken',
            {
              user_id: userId,
            },
        ).then((res) => {
          // eslint-disable-next-line @typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-member-access
          const client = connect('uq7mwraum8nu', res.data.data.message, '102359')
          const firehose = client.feed('events', this.feedUser)
          firehose.get({limit: 5}).then((body) => {
            this.firehoseData = body.results
            const subscription = firehose.subscribe((data) => {
              // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
              this.firehoseData = this.firehoseData.concat(data.new)
            })
            this.loading = false
          }).catch((error) => {
            // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
            this.$toast.warning(error)
            console.error(error)
          })
          this.$store.commit('SET_STREAM_TOKEN', res.data)
        })
      } catch (e) {
        console.error(e)
      }
    },
  },
})
</script>
