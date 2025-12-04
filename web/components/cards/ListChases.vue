<template>
  <div v-if="chases[0]">
    <div class="mx-auto max-w-7xl">
      <div v-if="mp4">
        <IconLink :url="videoLink" />
      </div>
      <div v-else-if="video">
        <Video
          :id="id"
          :url="videoLink"
          :live="live"
        />
      </div>
      <div v-else>
        <Map :chase="chases[0]" />
      </div>
    </div>
    <ListEvents v-if="!$route.params.id" />
    <div class="flex-col items-center justify-center font-sans lg:space-y-4 divide-y divide-gray-300 lg:divide-y-0">
      <show-chase
        v-for="(chase, $index) in chases"
        :key="$index"
        :chase="chase"
      />
    </div>
    <div v-if="!id">
      <infinite-loading
        v-if="chases.length"
        spinner="spiral"
        @infinite="getMoreChases"
      />
      <ScrollTopArrow />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import {
  doc,
  collection,
  query,
  orderBy,
  onSnapshot,
  limit,
  startAfter,
} from 'firebase/firestore'
// import { getAuth, onAuthStateChanged } from 'firebase/auth'
import { db } from '~/plugins/firebase'
import { Chase, Network } from '@/types'
import Logo from "~/components/Logo.vue";
import IconLink from "~/components/IconLink.vue";
import Hls from "hls.js";
import Plyr from "plyr";

const delay = (ms: number | undefined) => new Promise(res => setTimeout(res, ms));

export default Vue.extend({
  components: {
    Logo,
    IconLink,
  },
  middleware: 'auth',
  props: {
    id: {
      type: String,
      required: false,
    },
  },
  data () {
    return {
      // networks: undefined as any | Chase[],
      nextChases: [],
      lastVisible: {},
      chaseLimit: 25,
      refresh: false,
      /*
      rive: {
        wheel: {
          src: '/rive/wheel.riv',
          layout: {
            fit: 'cover',
            alignment: 'BottomCenter',
          },
        },
      },
       */
      chases: [] as Chase[],
      video: false,
      mp4: false,
      videoLink: '',
      live: false,
      votes: new Number() as number,
      playerOptions: {
        controls: [
          "play-large",
          "current-time",
          "play",
          "mute",
          "volume",
          "progress",
          "settings",
          "fullscreen",
        ],
        settings: ["quality", "speed", "loop"],
      },
    }
  },
  mounted () {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    if (this.chases[0]?.Live)
      this.getChase(this.chases[0].ID)
    else if (this.id)
      this.getChase(this.id)
    else
      this.getChases(this.chaseLimit)
      // this.fetchFromBundle() // fetchBundle shit not working with firebase v9
  },
  methods: {
    // toggleConfetti takes an argument that is a number in milliseconds
    async toggleConfetti (myDelay: number) {
      this.showConfetti = !this.showConfetti
      await delay(myDelay)
      this.showConfetti = !this.showConfetti
    },
    isMobile () {
      if ('maxTouchPoints' in navigator) return navigator.maxTouchPoints > 0 ? true : false
      else return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
    },
    navigationType () {
      if (window.performance.getEntriesByType('navigation')) {
        const p = window.performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming

        switch (p.type) {
          case 'navigate':
            return false
          case 'reload':
            return true
          case 'back_forward':
            return false
          case 'prerender':
            return false
        }
      }
      return false
    },
    getChase (id: string) {
      onSnapshot(doc(db, 'chases', id), (doc) => {
        if (doc.data()) {
          // this.toggleConfetti(1000)
          this.live = doc.data()?.Live as boolean
          // if (this.live)
          //    this.toggleConfetti(3000)
          const networks = doc.data()?.Networks as Network[]
          networks.forEach((network: { URL: string, Streams: any[] }) => {
            if (network.Streams)  {
              network.Streams.forEach((stream: { URL: string, Tier: number }) => {
                if (stream.URL !== "") {
                  this.mp4 = true
                  this.videoLink = stream.URL
                } else if (network.URL.includes("youtube")) {
                  this.video = true
                  this.videoLink = network.URL
                }
              })
            } else if (network.URL.includes("youtube")) {
              this.video = true
              this.videoLink = network.URL
            }
          })
          this.chases = []
          const chase = {
            ID: id,
            ...(doc.data()),
          } as Chase

          if (this.votes != chase.Votes) {
            if (chase.Votes)
              this.votes = chase.Votes
          }
          this.chases.push(chase)
        }
      })
    },
    getChases (myLimit: number) {
      const q = query(collection(db, 'chases'), orderBy('CreatedAt', 'desc'), limit(myLimit))
      onSnapshot(q, (querySnapshot) => {
        this.chases = querySnapshot.docs.map(doc => ({ ID: doc.id, ...doc.data() }))
        this.lastVisible = querySnapshot.docs[querySnapshot.docs?.length - 1]
      }, (error) => {
        console.error(error)
      })
    },
    getMoreChases ($state: { complete: () => void; loaded: () => void }) {
      if (this.lastVisible === undefined)
        $state.complete()
      else
        setTimeout(() => {
          const q = query(
            collection(db, 'chases'),
            orderBy('CreatedAt', 'desc'),
            startAfter(this.lastVisible),
            limit(this.chaseLimit),
          )
          onSnapshot(q, (querySnapshot) => {
            this.chases.push(...querySnapshot.docs.map(doc => ({ ID: doc.id, ...doc.data() })))
            $state.loaded()
            this.lastVisible = querySnapshot.docs[querySnapshot.docs?.length - 1]
          }, (error) => {
            console.error(error)
          })
        }, 500)
    },
  },
})
</script>
