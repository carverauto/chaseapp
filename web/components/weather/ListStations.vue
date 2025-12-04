<template>
  <span v-if="weatherStations">
    <MglMarker
      v-for="(w, key) in weatherStations"
      :key="key"
      :coordinates="[w.lng,w.lat]"
      :class="myClass"
    >
      <div slot="marker" class="absolute cursor-pointer">
        <img class="inline-block h-4 w-4 rounded-full" src="/noaa-dot.png" alt="NOAA Weather Station">
        <div v-if="emergency" class="animate-ping h-9 w-9 rounded-full border-0 bg-red-500" />
      </div>
      <MglPopup :close-button="false" :offset="20">
        <div class="flex">
          <div class="rounded-md bg-yellow-50 p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <!-- Heroicon name: solid/exclamation -->
                <svg class="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-yellow-800">
                  NOAA: Storm Surge
                </h3>
                <div class="mt-2 text-sm text-yellow-700">
                  <p>
                    <a target="_blank" :href="'https://tidesandcurrents.noaa.gov/stationhome.html?id=' + w.id ">{{ w.name }}</a>
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </MglPopup>
    </MglMarker>
  </span>
</template>

<script lang="ts">

import Vue from 'vue'
import { WeatherStation } from '~/types'
import { collection, onSnapshot, query } from 'firebase/firestore'
import { db } from '~/plugins/firebase'

export default Vue.extend({
  data() {
    return {
      weatherStations: undefined as any | WeatherStation[],
      emergency: false,
      myClass: 'animate-ping',
    }
  },
  mounted() {
    const q = query(collection(db, 'weather'))
    onSnapshot(q, (querySnapshot) => {
      this.weatherStations = []
      querySnapshot.forEach((doc) => {
        if (doc.data().stations[0].stormsurge) {
          this.weatherStations.push(doc.data().stations[0])
        }
      })
    })
    /*
    const q = query(collection(db, 'weather'))
    onSnapshot(q, (querySnapshot) => {
      this.weatherStations = []
      // this.weatherStations.push(...querySnapshot.docs.map(doc => ({id: doc.id, ...doc.data().stations[0]})))
    }, (error) => {
      console.error(error)
    })
     */
    /*
  this.$fire.firestore.collection('weather')
    .onSnapshot((querySnapshot) => {
      this.weatherStations = []
      querySnapshot.forEach((doc) => {
        // only show sensors that are reporting stormsurge = true
        if (doc.data().stations)
          if (doc.data().stations[0].stormsurge)
            this.weatherStations.push(doc.data().stations[0])
      })
      // console.log(this.weatherStations[0])
    }, (error) => {
      console.error('Problem with firebase/weatherstations: ', error)
    })
  },
     */
  },
})
</script>
