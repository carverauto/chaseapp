<template>
  <div v-if="chase.Votes > 0" class="flex items-center text-xs md:text-sm ml-2">
    <button @click="changeCount(1)" aria-label="Click to vote">
      <Donut />
      <span class="sr-only">Add a vote for a chase</span>
    </button>
    <div class="ml-3 text-donut-pink">
      {{ chase.Votes }} Donuts
    </div>
  </div>
  <div v-else class="flex items-center text-xs md:text-sm ml-2 text-gray-300 hover:text-gray-600">
    <button @click="changeCount(1)" aria-label="No votes, Click to vote">
      <DonutEmpty primary="#FFFFFF" secondary="#000000" class="h-6 w-6 md:w-8 md:h-8 fill-current" @click="changeCount(1)" />
    </button>
    <div class="ml-3 transition-colors">
      No Donuts..
    </div>
  </div>
</template>

<script lang="ts">

// eslint-disable-next-line import/named
import Vue, { PropType } from 'vue'
import { doc, setDoc, increment } from 'firebase/firestore'
import { mapGetters } from 'vuex'
import { Chase } from '@/types/'
import { db } from '~/plugins/firebase'

export default Vue.extend({
  props: {
    chase: {
      type: Object as PropType<Chase>,
      required: true,
    },
  },
  computed: {
    ...mapGetters(['count']),
  },
  methods: {
    async changeCount (amount: number) {
      const chaseRef = doc(db, 'chases', this.chase.ID)
      await setDoc(chaseRef, { Votes: increment(amount) }, { merge: true })
    },
  },
})
</script>
