<template>
  <div
    class="max-w-3xl px-4 py-4 sm:px-6 lg:rounded-lg z-20 lg:shadow dark:bg-gray-800 bg-white lg:max-w-4xl"
    :class="chase.Live ? 'border-2 border-green-600' : ''"
  >
    <div class="flex items-center justify-between">
      <NuxtLink
        :to="`/chase/${chase.ID}`"
        class="focus:outline-none truncate text-md md:text-2xl font-bold text-gray-700 dark:text-white cursor-pointer hover:text-gray-600 dark:hover:text-gray-200 hover:underline"
      >
        <p class="focus:outline-none truncate text-md md:text-2xl font-bold text-gray-700 dark:text-white cursor-pointer hover:text-gray-600 dark:hover:text-gray-200 hover:underline">
          {{ chase.Name }}
        </p>
      </NuxtLink>
      <div class="ml-2 flex-shrink-0 flex">
        <div v-if="chase.Live">
          <p class="inline-flex text-xs leading-5 font-semibold rounded-md bg-green-100 text-green-800">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-md uppercase text-xs font-medium bg-green-100 text-green-800">
              <svg
                id="Capa_1"
                class="-ml-0.5 mr-2 h-4 w-4 text-green-400 animate-pulse"
                fill="currentColor"
                enable-background="new 0 0 512 512"
                height="20"
                viewBox="0 0 512 512"
                width="20"
                xmlns="http://www.w3.org/2000/svg"
              >
                <g><circle
                  cx="256"
                  cy="256"
                  r="63.2"
                /><path d="m141.695 141.695-21.213-21.213c-36.251 36.251-56.215 84.379-56.215 135.518s19.964 99.267 56.216 135.518l21.213-21.213c-30.586-30.585-47.429-71.18-47.429-114.305s16.843-83.72 47.428-114.305z" /><path d="m96.273 96.273-21.212-21.212c-48.404 48.403-75.061 112.662-75.061 180.939s26.657 132.536 75.061 180.939l21.213-21.213c-42.737-42.737-66.274-99.462-66.274-159.726s23.537-116.989 66.273-159.727z" /><path d="m436.939 75.061-21.213 21.213c42.737 42.737 66.274 99.462 66.274 159.726s-23.537 116.989-66.273 159.727l21.213 21.213c48.403-48.404 75.06-112.663 75.06-180.94s-26.657-132.536-75.061-180.939z" /><path d="m391.518 120.482-21.213 21.213c30.585 30.585 47.428 71.18 47.428 114.305s-16.844 83.72-47.429 114.305l21.213 21.213c36.251-36.251 56.216-84.379 56.216-135.518s-19.964-99.267-56.215-135.518z" /><path d="m346.095 165.905-21.213 21.213c18.433 18.433 28.584 42.896 28.584 68.882s-10.151 50.449-28.584 68.882l21.213 21.213c24.1-24.099 37.372-56.096 37.372-90.095s-13.272-65.996-37.372-90.095z" /><path d="m165.905 165.905c-24.1 24.099-37.372 56.095-37.372 90.095s13.272 65.996 37.372 90.095l21.213-21.213c-18.433-18.433-28.584-42.896-28.584-68.882s10.151-50.449 28.584-68.882z" /></g></svg>
              live
            </span>
          </p>
        </div>
        <div v-else>
          <p class="inline-flex text-xs leading-5 font-semibold rounded-md bg-red-100 text-red-800">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-md uppercase text-xs font-medium bg-red-100 text-red-800">
              <svg
                class="-ml-0.5 mr-1.5 h-2 w-2 text-red-400"
                fill="currentColor"
                viewBox="0 0 2 2"
              >
                <rect
                  width="2"
                  height="2"
                  rx="0.5"
                  ry="0.5"
                />
              </svg>
              over
            </span>
          </p>
        </div>
        <div class="flex-shrink-0 self-center flex">
          <div class="relative z-30 inline-block text-left">
            <div>
              <button
                type="button"
                class="-m-2 p-2 rounded-full flex items-center text-gray-400 hover:text-gray-600"
                aria-expanded="false"
                aria-haspopup="true"
                @click="activeInfoWin = true"
              >
                <span class="sr-only">Open options</span>
                <!-- Heroicon name: solid/dots-vertical -->
                <svg
                  class="h-5 w-5"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
                </svg>
              </button>
            </div>

            <transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div
                class="origin-top-right absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none"
                role="menu"
                aria-orientation="vertical"
                aria-labelledby="menu-0-button"
                tabindex="-1"
              >
                <div
                  v-if="activeInfoWin"
                  v-on-clickaway="showInfoWin"
                  class="py-1"
                  role="none"
                >
                  <!-- Active: "bg-gray-100 text-gray-900", Not Active: "text-gray-700" -->
                  <div
                    id="menu-0-item-0"
                    class="cursor-pointer text-gray-700 flex px-4 py-2 text-sm"
                    role="menuitem"
                    tabindex="-1"
                  >
                    <!-- Heroicon name: solid/star -->
                    <svg
                      class="mr-3 h-5 w-5 text-gray-400"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                    <span @click="createBookmark(chase.ID)">Add to favorites</span>
                  </div>
                  <div
                    id="menu-0-item-1"
                    class="cursor-pointer text-gray-700 flex px-4 py-2 text-sm"
                    role="menuitem"
                    tabindex="-1"
                    @click="startShare"
                  >
                    <!-- Heroicon name: solid/share -->
                    <svg
                      class="mr-3 h-5 w-5 text-gray-400"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M15 8a3 3 0 10-2.977-2.63l-4.94 2.47a3 3 0 100 4.319l4.94 2.47a3 3 0 10.895-1.789l-4.94-2.47a3.027 3.027 0 000-.74l4.94-2.47C13.456 7.68 14.19 8 15 8z"
                      />
                    </svg>
                    <span>Share</span>
                    <div v-if="showShare">
                      <Share :chase="chase" />
                    </div>
                  </div>
                  <a
                    id="menu-0-item-2"
                    href="#"
                    class="text-gray-700 flex px-4 py-2 text-sm"
                    role="menuitem"
                    tabindex="-1"
                  >
                    <!-- Heroicon name: solid/flag -->
                    <svg
                      class="mr-3 h-5 w-5 text-gray-400"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M3 6a3 3 0 013-3h10a1 1 0 01.8 1.6L14.25 8l2.55 3.4A1 1 0 0116 13H6a1 1 0 00-1 1v3a1 1 0 11-2 0V6z"
                        clip-rule="evenodd"
                      />
                    </svg>
                    <span>Report content</span>
                  </a>
                </div>
              </div>
            </transition>
          </div>
        </div>
      </div>
    </div>

    <div class="sm:flex">
      <div
        v-if="chase.ImageURL"
        class="mt-2 sm:mb-0 mb-4 flex-shrink-0 sm:mr-4 lg:self-center"
      >
        <div v-if="$route.params.id">
          <img
            v-lazy-load
            class="h-32 w-full sm:w-32"
            :src="lgImage"
            alt="Chase Image or Screengrab"
          >
        </div>
        <div v-else>
          <NuxtLink :to="`/chase/${chase.ID}`">
            <img
              v-lazy-load
              class="h-32 w-full sm:w-32"
              :src="lgImage"
              alt="Chase image or screengrab"
            >
          </NuxtLink>
        </div>
      </div>
      <div class="sm:flex sm:justify-between py-3 md:py-5 mb-4">
        <p class="flex text-sm text-justify md:text-lg text-gray-600 dark:text-gray-300">
          {{ chase.Desc }}
        </p>
      </div>
    </div>

    <ChaseSentiment
      v-if="chase.sentiment"
      v-tooltip="tooltips.msg.sentiment"
      :sentiment="chase.sentiment"
      class="mb-4"
    />

    <div class="items-center justify-between mt-2 grid grid-cols-3">
      <div class="flex items-center text-gray-500 sm:mt-0 lg:ml-0.5">
        <div v-if="chase.Type !== 'Rocket'">
          <div v-if="chase.Wheels">
            <ChaseWheels
              v-tooltip="tooltips.msg.wheels"
              :w1="chase.Wheels.W1"
              :w2="chase.Wheels.W2"
              :w3="chase.Wheels.W3"
              :w4="chase.Wheels.W4"
              class="h-7 w-7"
            />
          </div>
        </div>
        <ChaseVoting
          v-tooltip="tooltips.msg.donuts"
          :chase="chase"
        />
      </div>
      <div class="z-25 flex flex-row mt-1 items-center justify-center w-parent">
        <button
          aria-label="Go to top of page"
          class="group bg-white rounded-md text-gray-500 items-center text-base font-medium hover:text-gray-900 focus:outline-none"
          @click="isOpen = !isOpen"
        >
          <!--
              Heroicon name: solid/chevron-down

              Item active: "text-gray-600", Item inactive: "text-gray-400"
              -->
          <div v-if="!$route.params.id">
            <svg
              :class="{'rotate-180': isOpen}"
              class="z-1 h-7 w-7 transform text-gray-400 group-hover:text-gray-500"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                fill-rule="evenodd"
                d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                clip-rule="evenodd"
              />
            </svg>
          </div>
        </button>
      </div>
      <!-- Heroicon name: solid/calendar -->

      <!-- <div class="items-center justify-between mt-2 grid grid-cols-3">
        <div class="flex items-center text-gray-500 sm:mt-0 lg:ml-0.5">
        -->
      <client-only>
        <div class="flex justify-end">
          <div
            v-if="chase.CreatedAt"
            class="text-xs font-light text-gray-600 dark:text-gray-400"
          >
            <div class="py-0.5">
              <time
                class="flex items-center"
                :datetime="chase.CreatedAt.toDate()"
              >
                <img
                  v-tooltip="tooltips.msg.chaseTime"
                  src="/alarm-clock.svg"
                  aria-hidden="true"
                  height="4"
                  width="4"
                  class="h-4 w-auto mr-1"
                  alt="chase start"
                >
                <span>
                  {{ $dayjs(chase.CreatedAt.toDate()).fromNow() }}
                </span>
              </time>
            </div>
            <div
              v-if="chase.EndedAt"
              v-show="!chase.Live"
            >
              <div class="flex items-center">
                <img
                  v-tooltip="tooltips.msg.chaseDuration"
                  src="/stopwatch.svg"
                  aria-hidden="true"
                  height="4"
                  width="4"
                  class="h-4 w-auto object-right mr-1 ml-0.5"
                  alt="duration"
                >
                <span class="">
                  {{ dateDiff }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </client-only>
    </div>
    <div class="overflow-hidden">
      <!--
      enter-active-class="transform transition-all ease-in duration-1000"
      enter-from-class="-translate-y-full"
      enter-to-class="-translate-y-0"
      leave-active-class="transform transition-all ease-out duration-500"
      leave-from-class="-translate-y-0"
      leave-to-class="-translate-y-full"
      -->
      <transition
        enter-active-class="transition duration-300 ease-in-out transform"
        enter-from-class="-translate-y-full"
        enter-to-class="translate-y-0"
        leave-active-class="transition duration-300 ease-in-out transform"
        leave-from-class="translate-y-0"
        leave-to-class="-translate-y-full"
      >
        <div
          v-show="$route.params.id || isOpen"
          class="h-full z-10 origin-top"
        >
          <div class="md:px-6 pt-4 flex-1 min-w-0 max-w-full min-h-content z-5 origin-top">
            <div v-if="chase.Networks">
              <list-networks
                v-for="(network, $index) in sortedNetworksES6"
                :key="$index"
                :network="network"
              />
              <div class="items-center justify-between grid grid-cols-1">
                <div class="flex pt-2 md:pt-3 items-start justify-between">
                  <span class="flex-col text-left text-sm text-gray-300">
                    <time
                      class="whitespace-nowrap"
                      :datetime="$dayjs(chase.CreatedAt.toDate()).format('MMM Do YYYY')"
                    ><i class="fas fa-calendar-day" /> {{ $dayjs(chase.CreatedAt.toDate()).format('MMM Do YYYY') }}</time>
                    <time
                      class="whitespace-nowrap"
                      :datetime="$dayjs(chase.CreatedAt.toDate()).tz('America/Los_Angeles').format('h:mm A z')"
                    ><i class="far fa-clock" /> {{ $dayjs(chase.CreatedAt.toDate()).tz('America/Los_Angeles').format('h:mm A z') }}</time>
                  </span>
                </div>
              </div>
              <div
                v-if="!$route.params.id"
                class="z-25 flex flex-row mt-1 items-center justify-center w-parent"
              >
                <nuxt-link
                  class="text-xs text-gray-600"
                  :to="`/chase/${chase.ID}`"
                >
                  [ CLICK TO CHAT ⌨️ ]
                </nuxt-link>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </div>

    <div v-if="$route.params.id">
      <div v-if="showChat">
        <div v-if="chase.ID">
          <ShowChat
            class="mt-2"
            :chase-id="chase.ID"
            :chase-name="chase.Name"
          />
        </div>
        <div class="grid grid-cols-1 place-content-center">
          <button
            class="text-xs text-gray-600"
            @click="showChat = false"
          >
            HIDE CHAT
          </button>
        </div>
      </div>
      <div v-else>
        <div class="grid grid-cols-1 place-content-center">
          <button
            class="text-xs text-gray-600"
            @click="showChat = true"
          >
            SHOW CHAT
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

import Vue, { PropType } from 'vue'
import { mapState, mapGetters } from 'vuex'
import { doc, onSnapshot, setDoc } from 'firebase/firestore'
import { AuthUser } from '@/types'
// import { auth } from '~/plugins/firebase'
import { Chase } from '@/types/'
import { db } from '~/plugins/firebase'
// eslint-disable-next-line @typescript-eslint/no-var-requires,@typescript-eslint/no-unsafe-assignment
const clone = require('lodash.clone')


export default Vue.extend({
  props: {
    chase: {
      type: Object as PropType<Chase>,
      required: true,
    },
  },
  asyncData ({ $dayjs, chase }) {
    try {
      return {
        newDate: $dayjs(chase.CreatedAt).fromNow() as Date,
      }
    } catch (error) {
      console.error(error)
    }
  },
  data () {
    return {
      showConfetti: false,
      myConfetti: '',
      imgRegex: /\.([0-9a-z]+)(?:[?#]|$)/i,
      smImgReplace: '_200x200.webp?' as string,
      lgImgReplace: '_1200x600.webp?' as string,
      theme: 'dark',
      isOpen: !!this.chase.Live,
      picture: true,
      duration: '',
      showChat: true as boolean,
      showDemoOptions: false,
      showShare: false,
      updatingData: false,
      currentUserId: undefined as undefined|string,
      activeInfoWin: false,
      favorites: [],
      tooltips: {
        msg: {
          sentiment: "Chase Sentiment",
          wheels: "ChaseWheels Wheel Tracker",
          donuts: "Votes",
          chatters: "Active in Chat",
          chaseTime: "Time of Chase",
          chaseDuration: "Duration of Chase",
        },
      },
    }
  },
  computed: {
    ...mapState({
      authUser: state => state.authUser as AuthUser,
    }),
    ...mapGetters({
      isLoggedIn: 'isLoggedIn',
    }),
    lgImage (): string {
      return this.chase.ImageURL.replace(this.imgRegex, this.lgImgReplace)
    },
    smImage (): string {
      return this.chase.ImageURL.replace(this.imgRegex, this.smImgReplace)
    },
    authUser () {
      return this.$store.state.authUser
    },
    dateDiff () {
      if (this.chase.EndedAt) {
        const x = this.$dayjs(this.chase.CreatedAt.toDate())
        const y = this.$dayjs(this.chase.EndedAt.toDate())
        const mins = y.diff(x, 'minutes', true)
        return this.$dayjs.duration(mins, 'minutes').humanize()
      }
      return false
    },
    /*
    created () {
      // return this.getAuthUser()
    },
     */
    sortedNetworksES6 () {
      if (this.chase.Networks) {
        const sortedNets = clone(this.chase.Networks)
        return sortedNets.sort((a: { Tier: number }, b: { Tier: number }) => {
          return b.Tier - a.Tier
        })
      } else
        return false
    },
  },
  watch: {
    // isOpen: function () {this.$fire.analytics.logEvent('Chase LIVE status changed')}
    /*
    currentUserId () {
      this.showChat = false
      setTimeout(() => (this.showChat = true), 150)
    },
     */
  },
  mounted () {
    if (this.isLoggedIn)
      this.getBookmarks()
  },
  methods: {
    startShare () {
      // eslint-disable-next-line no-return-assign
      return this.showShare = !this.showShare
    },
    getBookmarks (): any {
      onSnapshot(doc(db, 'users', this.authUser.uid), async (doc) => {
        const userData = await doc.data()
        this.favorites = userData.favorites
      })
    },
    async createBookmark (val: string) {
      if (this.isLoggedIn) {
        try {
          this.favorites.push(val)
          await setDoc(doc(db, 'users', this.authUser.uid), {
            favorites: this.favorites,
            lastUpdated: Date.now(),
          }, { merge: true })
        } catch (e) {
          this.$toast.show({
            type: 'danger',
            title: 'ChaseApp - Favorites',
            message: `Unable to update favorites <b>${e}</b>`,
            timeout: 0,
          })
        }
        this.$toast.show({
          type: 'success',
          title: 'ChaseApp - Favorites',
          message: 'Favorites Updated',
          timeout: 3,
        })
      } else
        this.$toast.denied('Not logged in')
    },
    showInfoWin () {
      this.activeInfoWin = !this.activeInfoWin
    },
    async getAuthUser () {
      // this.currentUserId = await this.authUser?.uid || null
    },
  },
})
</script>
