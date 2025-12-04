<template>
  <div v-if="leo.length > 0 || media.length > 0" class="flex items-center justify-between">
    <p class="mr-1 text-sm text-cool-gray-700">Cluster Alert</p>
    <div v-if="leo.length > 0" class="flex -space-x-2 relative z-0 overflow-hidden">
      <ShowLEO v-for="(leoData, $index) in leo" :key="$index" :index="$index" :leo="leoData" />
    </div>
    <div v-if="media.length > 0" class="ml-1 flex -space-x-2 relative z-0 overflow-hidden">
      <ShowMedia v-for="(mediaData, $index) in media" :key="$index" :index="$index" :media="mediaData" />
    </div>
    <div class="ml-3">
      <p v-if="leo" class="lg:text-sm md:text-xs text-xs font-medium text-gray-700 group-hover:text-gray-900">
        Police: {{ leo.length }}
      </p>
      <p v-if="media" class="text-xs font-medium text-gray-500 group-hover:text-gray-700">
        Media: {{ media.length }}
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { ref, onValue } from 'firebase/database'
import { rtdb } from '~/plugins/firebase'
import { Cluster, GeojsonFeature } from '~/types'

export default Vue.extend({
  data () {
    return {
      clusters: {} as any|Cluster,
      media: [] as any|GeojsonFeature[],
      leo: [] as any|GeojsonFeature[],
    }
  },
  mounted () {
    this.getClusters()
  },
  methods: {
    getClusters () {
      const bofRef = ref(rtdb, 'bof')
      onValue(bofRef, (snapshot) => {
        this.leo = []
        this.media = []
        const data = snapshot.val() as GeojsonFeature[]
        data.forEach((bof) => {
          if (bof.properties.dbscan !== 'noise')
            switch (bof.properties.group) {
              case 'leo':
                this.leo.push(bof)
                break
              case 'media':
                this.media.push(bof)
            }
        })
        this.media.sort((a: { properties: { cluster: number } }, b: { properties: { cluster: number } }) => a.properties.cluster - b.properties.cluster)
        this.leo.sort((a: { properties: { cluster: number } }, b: { properties: { cluster: number } }) => a.properties.cluster - b.properties.cluster)
      }, (error) => {
        console.error(error)
      })
    },
  },
})
</script>
