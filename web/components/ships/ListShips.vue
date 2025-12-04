<template>
  <span v-if="ships">
    <MglMarker
      v-for="(ship, key) in ships"
      :key="key"
      :coordinates="[ship.longitude,ship.latitude]"
      :class="myClass"
    >
      <div slot="marker" class="absolute cursor-pointer">
        <div
          class="inline-block h-8 w-8 rounded-full boat"
          :style="`transform: rotate(${ship.heading}deg)`"
        />
      </div>
      <MglPopup :close-button="false" :offset="20">
        <div class="flex">
          <div class="mr-4 flex-shrink-0 self-center">
            <img class="h-8 w-8" height="8" width="8" src="/boat.png" alt="boat">
          </div>
          <div>
            <h4 class="text-lg font-bold">{{ ship.name }}</h4>
            <p class="mt-1">
              MMSI: {{ ship.mmsi }}
            </p>
            <p v-if="ship.dest">
              DEST: {{ ship.dest }}
            </p>
          </div>
        </div>
      </MglPopup>
    </MglMarker>
  </span>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    ships: {
      type: Array,
      required: false,
    },
  },
  data () {
    return {
      live: false,
      myClass: 'animate-ping',
    }
  },
})
</script>
