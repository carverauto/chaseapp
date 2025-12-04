<template>
  <div :class="vidCols">
    <div v-if="newWindow">
      <div class="bg-white px-4 py-5 sm:px-6">
        <div class="flex space-x-3">
          <div class="flex-shrink-0">
            <img
              class="h-10 w-10 rounded-full"
              src="/map.png"
              alt="special map marker"
            >
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium text-gray-900">
              {{ newWindowData.name }}
            </p>
            <p class="text-sm text-gray-500 overflow-y-scroll">
              {{ newWindowData.desc }}
            </p>
            <div class="text-sm text-gray-500">
              <div
                v-for="url in newWindowData.urls"
                :key="url"
              >
                <p
                  class="hover:underline"
                  @click="openYoutube(url)"
                >
                  {{ url }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="h-72 w-auto max-w-3xl lg:rounded-lg lg:shadow bg-white lg:max-w-4xl lg:mb-4">
      <MglMap
        :access-token="accessToken"
        :map-style="mapStyle"
        glyphs="https://fonts.openmaptiles.org/{fontstack}/{range}.pbf"
        :center="myCoords"
        :zoom="zoom"
        @load="onMapLoaded"
      >
        <span v-if="birds && birds.length > 0">
          <MglMarker
            v-for="(l, key) in birds"
            :key="key"
            :coordinates="[l.lon,l.lat]"
            :class="myClass"
          >
            <div
              slot="marker"
              class="cursor-pointer absolute"
            >
              <div
                class="inline-block h-8 w-8 rounded-full"
                :class="l.type"
                :style="`transform: rotate(${l.track}deg)`"
              />
              <div
                v-if="chase.live"
                class="animate-ping h-9 w-9 rounded-full border-0 bg-red-500"
              />
            </div>

            <MglPopup
              :close-button="false"
              :offset="20"
            >
              <div class="flex h-auto w-auto">
                <a
                  target="_blank"
                  class="outline-none"
                  :href="'https://registry.faa.gov/AircraftInquiry/Search/NNumberResult?nNumberTxt=' + l.tailno"
                >{{ l.tailno }}</a>
                <span class="inline-block relative">
                  <div v-if="l.tailno === 'N29HD'">
                    <img
                      class="h-14 w-14 rounded-full"
                      :src="n29hdImageURL"
                      alt="n29hd network image"
                    >
                  </div>
                  <div v-else-if="l.imageUrl">
                    <img
                      class="h-14 w-14 rounded-full"
                      :src="l.imageUrl"
                      alt="network or agency logo"
                    >
                  </div>
                  <div v-if="!l.imageUrl">
                    <span class="inline-block h-14 w-14 rounded-full overflow-hidden bg-gray-100">
                      <svg
                        class="h-full w-full text-gray-300"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path d="M24 20.993V24H0v-2.996A14.977 14.977 0 0112.004 15c4.904 0 9.26 2.354 11.996 5.993zM16.002 8.999a4 4 0 11-8 0 4 4 0 018 0z" />
                      </svg>
                    </span>
                  </div>
                  <span class="absolute top-0 right-0 block h-4 w-4 rounded-full ring-2 ring-white bg-green-400" />
                </span>
                <div v-if="nickname && l.tailno === 'N29HD'">
                  <span class="absolute bottom-2 right-0 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                    {{ nickname }}
                  </span>
                </div>
                <div v-else>
                  <span class="absolute bottom-2 right-0 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                    {{ l.group }}
                  </span>
                </div>
              </div>
            </MglPopup>
          </MglMarker>
        </span>

        <weather-list-stations />
        <rocket-list-launches />
        <ships-list-ships :ships="ships[0]" />
        <markers-list-markers
          @openInfoWin="openInfoWin"
          @openYoutube="openYoutube"
        />

        <MglNavigationControl position="top-right" />
        <MglGeolocateControl ref="geolocateControl" />
        <MglFullscreenControl position="top-right" />
      </MglMap>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { onValue, ref } from 'firebase/database'
import { ADSB, AisBoat } from '@/types'
import { rtdb } from '~/plugins/firebase'
import {GeoJSON} from "geojson"
import * as turf from "@turf/turf"


// eslint-disable-next-line @typescript-eslint/no-var-requires
const MapboxTraffic = require('@mapbox/mapbox-gl-traffic')

export default Vue.extend({
  name: 'Map',
  props: {
    chase: {
      type: Object,
      required: true,
    },
  },
  data () {
    return {
      // lottieSettings: { path: '/json/lf30_editor_9uk17mde.json', class: 'h-32 w-32' },
      quake: true,
      youtube: '',
      vidCols: 'lg:grid-cols-1',
      myClass: 'animate-ping',
      nickname: '',
      database: {},
      spriteClass: 'sprite015',
      birds: undefined as undefined|ADSB[],
      ships: [] as []|AisBoat[],
      markers: [],
      markerCoordinates: [-118.2607073, 34.0201613] as number[],
      zoom: 3,
      heliSprite: 'sprite3045',
      planeSprite: 'plane3045',
      ac130Sprite: 'ac1303045',
      center: this.$store.state.coords as number[],
      n29hdImageURL: undefined as undefined|string,
      mapStyle: 'mapbox://styles/mapbox/dark-v10?optimize=true',
      accessToken: 'pk.eyJ1IjoibWZyZWVtYW40NTEiLCJhIjoiY2tyaWRyYnNlMXJleTJwbTRjYTAzYWhjaCJ9.hhwesJS258kLg3-XdLmPqg',
      newWindow: false,
      newWindowData: {},
      activeInfoWin: false,
      video: false,
    }
  },
  computed: {
    /*
    openInfoWin: {
      get () {
        console.log('get')
        // @ts-ignore
        return this.newWindow
      },
      set (newValue) {
        console.log('set')
        // @ts-ignore
        this.infoWin = newValue
      },
    },
     */
    myCoords () {
      // eslint-disable-next-line vue/no-side-effects-in-computed-properties
      return this.$store.state.mapCoords
    },
  },
  watch: {
  },
  mounted () {
    const convDate = new Date().getUTCHours()
    this.n29hdImageURL = this.setN29HDimageURL(convDate)
    this.getWeather()
    this.getAirships()
    this.getBoats()
    if (this.$route.query.lat && this.$route.query.lon) {
      this.center = [this.$route.query.lat, this.$route.query.lon]
      this.zoom = this.$route.query.zoom ? this.$route.query.zoom : 10
    }
  },
  methods: {
    createBookmark (val: string) {
      if (process.browser)
        return this.browser.bookmarks.create(val)
    },
    openYoutube (val: string) {
      // Grab text
      const matches = /(?:http:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=)?([^& \n]+)/g.exec(val)
      if (!matches)
        return true

      this.video = true
      this.youtube = 'https://www.youtube.com/embed/' + matches[1]
      this.newWindow = false
    },
    openInfoWin (val: any) {
      this.newWindow = true
      this.newWindowData = val
    },
    updateCoords (coords: number[]) {
      this.center = [coords[0], coords[1]]
      this.zoom = this.$route.query.zoom ? this.$route.query.zoom : 10
      this.newWindow = false
      this.video = false
    },
    getWeather () {
      // fetch radar rasters for geoserver WMS service and add to map
    },
    getTFRs (): GeoJSON {
      const tfrRef = ref(rtdb, 'tfr/activeTFRs')
      let tfrData: GeoJSON
      onValue(tfrRef, (snapshot) => {
        // const ships: AisBoat[] = []
        tfrData = snapshot.val() as GeoJSON
        // need to build the features array and stuff it full of crap
        if (tfrData) 
          this.tfrCallback(tfrData)
         else
          console.log('No data returned from server/missing')
      }, (error) => {
        console.error(error)
      })
    },
    getBoats () {
      const adsbRef = ref(rtdb, 'ships')
      onValue(adsbRef, (snapshot) => {
        const ships: AisBoat[] = []
        const shipData = snapshot.val()
        // need to build the features array and stuff it full of crap
        if (shipData)
          for (const [key, value] of Object.entries(shipData)) {
            if (key) {
              const item = value as AisBoat
              ships.push(item)
              this.ships = ships
            }
          }
        else
          console.log('No data returned from server/missing')
      }, (error) => {
        console.error(error)
      })
    },
    getAirships () {
      const adsbRef = ref(rtdb, 'adsb')
      onValue(adsbRef, (snapshot) => {
        const birds: ADSB[] = []
        // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
        const birdData = snapshot.val()
        // need to build the features array and stuff it full of crap
        if (birdData)
          for (const [key, value] of Object.entries(birdData)) {
            if (key) {
              const item = value as ADSB
              birds.push(item)
              this.birds = birds
            }
          }
        else
          console.log('No data returned from server/missing')
      }, (error) => {
        console.error(error)
      })
    },
    between (x: number, min: number, max: number): boolean {
      return x >= min && x <= max
    },
    // TODO: need to clean this up
    setN29HDimageURL (hour: number) {
      // Des is in the chopper for CBS
      if (hour < 1) {
        this.nickname = 'Desmond'
        return '/networks/cbsla.jpg'
      }
      if (this.between(hour, 1, 6)) {
        this.nickname = 'Desmond'
        return '/networks/cbsla.jpg'
      }
      if (this.between(hour, 23, 24)) {
        this.nickname = 'Desmond'
        return '/networks/cbsla.jpg'
      }
      // Stu for FOX
      if (this.between(hour, 14, 22)) {
        this.nickname = 'Stu'
        return '/networks/foxla.jpg'
      }
    },
    convertTZ (date: Date, tzString: string) {
      return new Date((date).toLocaleString('en-US', { timeZone: tzString }))
    },
    tfrCallback: function (data: any) {
      // iterate through the features and add them to the map
      if (data) 
        for (const feature of data.features) {
          const tfr = feature as GeoJSON.Feature<GeoJSON.Polygon>
          // let center = turf.point([-74.50, 40]);
          let center = turf.point([tfr.geometry.coordinates[0], tfr.geometry.coordinates[1]])

          let radius = tfr.properties.radiusArc
          let options = {
            steps: 80,
            units: 'kilometers',
          };

          let circle = turf.circle(center, radius, options);

          this.map.addLayer({
            id: "circle-fill",
            type: "fill",
            source: {
              type: "geojson",
              data: circle,
            },
            paint: {
              "fill-color": "pink",
              "fill-opacity": 0.5,
            },
            layout: {},
          });

          this.map.addLayer({
            'id': 'tfr-labels',
            'type': 'symbol',
            'source': {
              type: 'geojson',
              data: tfr,
            },
            'layout': {
              'text-field': 'TFR',
              'text-font': ['Open Sans Bold', 'Arial Unicode MS Bold'],
              'text-size': 12,
            },
            'paint': {
              'text-color': 'rgba(0,0,0,0.5)',
            },
          });

        }
       else 
        console.log('No TFRs')
      
    },
    async onMapLoaded (event: any) {
      let tfrGeoJSON = this.getTFRs()


      /*
      this.map.addSource('tfr', {
        type: 'geojson',
        data: tfrGeoJSON,
      })

      this.map.addLayer({
        "id": "circle500",
        "type": "circle",
        "source": "tfr",
        "paint": {
          "circle-radius": {
            stops: [
              [5, 1],
              [15, 1024],
            ],
            base: 2,
          },
          "circle-color": "red",
          "circle-opacity": 0.6,
        },
      });

       */


      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-assignment
      this.map = event.map
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      this.map.on('click', (e: { lngLat: { lng: number; lat: number } }) => {
        this.$store.commit('SET_COORDS', [e.lngLat.lng, e.lngLat.lat])
        this.newWindow = false
      })
      this.map.addSource('radar', {
        type: 'raster',
        tiles: [
        'https://mesonet.agron.iastate.edu/cache/tile.py/1.0.0/nexrad-n0q-900913/{z}/{x}/{y}.png',
        ],
        tileSize: 256,
      })
      this.map.addSource('warnings', {
        type: 'raster',
        tiles: [
        'https://opengeo.ncep.noaa.gov/geoserver/wwa/warnings/ows?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&FORMAT=image%2Fpng&TRANSPARENT=true&TILED=true&LAYERS=warnings&WIDTH=256&HEIGHT=256&SRS=EPSG%3A3857&BBOX={bbox-epsg-3857}',
        ],
        tileSize: 256,
      })
      this.map.addLayer({
        id: 'warnings',
        type: 'raster',
        source: 'warnings',
        minzoom: 0,
        maxzoom: 22,
      })
      this.map.addLayer({
        id: 'radar',
        type: 'raster',
        source: 'radar',
        minzoom: 0,
        maxzoom: 24,
      })
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      this.map.addControl(new MapboxTraffic({ showTraffic: false, showTrafficButton: true }))
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      this.map.attributionControl = false
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-return
      await this.$nextTick().then(() => this.map.resize())
      // this.map.resize()
    },
  },
})
</script>

<style src="mapbox-gl/dist/mapbox-gl.css" />

<style>
.halloween {
  position: absolute;
  font-family: sans-serif;
  margin-top: 5px;
  margin-left: 5px;
  padding: 5px;
  width: 10%;
}

.quake-info {
  position: absolute;
  font-family: sans-serif;
  margin-top: 5px;
  margin-left: 5px;
  padding: 5px;
  width: 30%;
  border: 2px solid black;
  font-size: 14px;
  color: #222;
  background-color: #fff;
  border-radius: 3px;
}

.marker {
  width: 0px;
  height: 0px;
}

.mapboxgl-popup {
  max-width: 400px;
  font: 12px/20px 'Helvetica Neue', Arial, Helvetica, sans-serif;
}
.mapboxgl-ctrl-traffic {
  background-image: url('data:image/svg+xml;charset=utf8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22%23333333%22%20preserveAspectRatio%3D%22xMidYMid%20meet%22%20viewBox%3D%22-2%20-2%2022%2022%22%3E%0D%0A%3Cpath%20d%3D%22M13%2C4.1L12%2C3H6L5%2C4.1l-2%2C9.8L4%2C15h10l1-1.1L13%2C4.1z%20M10%2C13H8v-3h2V13z%20M10%2C8H8V5h2V8z%22%2F%3E%0D%0A%3C%2Fsvg%3E');
}

.mapboxgl-ctrl-map {
  background-image: url('data:image/svg+xml;charset=utf8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22%23333333%22%20viewBox%3D%22-10%20-10%2060%2060%22%20preserveAspectRatio%3D%22xMidYMid%20meet%22%3E%3Cg%3E%3Cpath%20d%3D%22m25%2031.640000000000004v-19.766666666666673l-10-3.511666666666663v19.766666666666666z%20m9.140000000000008-26.640000000000004q0.8599999999999923%200%200.8599999999999923%200.8600000000000003v25.156666666666666q0%200.625-0.625%200.783333333333335l-9.375%203.1999999999999993-10-3.5133333333333354-8.906666666666668%203.4383333333333326-0.2333333333333334%200.07833333333333314q-0.8616666666666664%200-0.8616666666666664-0.8599999999999994v-25.156666666666663q0-0.625%200.6233333333333331-0.7833333333333332l9.378333333333334-3.198333333333334%2010%203.5133333333333336%208.905000000000001-3.4383333333333344z%22%3E%3C%2Fpath%3E%3C%2Fg%3E%3C%2Fsvg%3E');
}

.mapboxgl-ctrl-attrib-inner a:last-of-type { display: none; }

.mapbox-marker {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  cursor: pointer;
}

</style>
