<template>
  <span v-if="markers">
    <MglMarker
      v-for="(marker, key) in markers"
      :key="key"
      :coordinates="[marker.location._long,marker.location._lat]"
      :class="myClass"
    >
      <div slot="marker" class="absolute cursor-pointer">
        <div v-if="marker.type === 'thanksgiving'">
          <div @click="$emit('openInfoWin', marker)">
            <Lottie :lottie-settings="thanksgivingLottie" />
          </div>
        </div>
        <div v-if="marker.type === 'halloween'">
          <Lottie :lottie-settings="halloweenLottie" />
        </div>
        <div v-if="marker.type === 'volcano'">
          <img class="inline-block h-6 w-6 rounded-full" src="/volcano.png" alt="Active Volcano" @click="$emit('openInfoWin', marker)">
        </div>
        <div v-if="marker.type === 'fire'">
          <img class="inline-block h-8 w-8 rounded-full" src="/fire.svg" alt="Active Fire">
        </div>
        <div v-if="marker.type === 'crime'">
          <img class="inline-block h-6 w-6 rounded-full" src="/crime.svg" alt="Criminal Investigation">
        </div>
        <div v-if="marker.type === 'sad'">
          <img class="inline-block h-6 w-6 rounded-full" src="/sad.png" alt="Disturbance/Event">
        </div>
        <div v-if="marker.type === 'patrol'">
          <img class="inline-block h-6 w-6 rounded-full" src="/patrol.png" alt="Live on Patrol">
        </div>
        <div v-if="marker.type === 'airshow'">
          <img class="inline-block h-6 w-6 rounded-full" src="/blueAngels.jpg" alt="Airshow">
        </div>
        <div v-if="marker.type === 'liveatc'">
          <img class="inline-block h-6 w-6 rounded-full" src="/trans.png" alt="LiveATC">
        </div>
        <div v-if="marker.type === 'map'">
          <img class="inline-block h-6 w-6 rounded-full" src="/marker.png" alt="Map">
        </div>
        <span v-if="marker.live" class="absolute animate-ping top-0 right-0 block h-2.5 w-2.5 rounded-full ring-2 bg-green-400" />
      </div>
    </MglMarker>
  </span>
</template>

<script lang="ts">

import Vue from 'vue'
import { collection, onSnapshot, query } from 'firebase/firestore'
import { Marker } from '~/types'
import { db } from '~/plugins/firebase'

export default Vue.extend({
  data () {
    return {
      markers: undefined as any|Marker[],
      live: false,
      myClass: 'animate-ping',
      thanksgivingLottie: {
        path: '/json/56305-cool-turkey.json',
        class: 'inline-block h-12 w-12 rounded-full',
      },
      halloweenLottie: {
        path: '/json/9976-halloween-witch-and-broom.json',
        class: 'inline-block h-12 w-12 rounded-full',
      },
    }
  },
  mounted () {
    const q = query(collection(db, 'markers'))
    onSnapshot(q, (querySnapshot) => {
      this.markers = []
      this.markers.push(...querySnapshot.docs.map(doc => ({ id: doc.id, ...doc.data() })))
    }, (error) => {
      console.error(error)
    })
  },
})
</script>
