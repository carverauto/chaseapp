<template>
  <div
    v-if="preload"
    class="loading-page"
  >
    <p>Loading..</p>
  </div>
  <div v-else>
    <!-- <nav class="bg-gradient-to-l from-yellow-300 to-blue-500 z-30 sticky top-0"> -->
    <nav class="bg-gradient-to-l from-red-600 via-indigo-700 to-blue-500 z-30 sticky top-0">
      <div class="max-w-7xl mx-auto px-2 sm:px-4 lg:px-8">
        <div class="relative flex items-center justify-between h-16">
          <div class="flex items-center px-2 lg:px-0">
            <div class="flex-shrink-0">
              <NuxtLink to="/">
                <span class="inline-block relative">
                  <img
                    class="rounded-md block lg:hidden h-8 w-auto"
                    src="/icon.png"
                    alt="ChaseApp"
                  >
                  <brand-logo
                    primary="#FFFFFF"
                    secondary="#FFFFFF"
                    class="hidden lg:block mt-1.5 h-10 w-auto"
                  />
                </span>
              </NuxtLink>
            </div>
          </div>
          <div class="flex-1 flex justify-center px-2 lg:ml-6 lg:justify-end">
            <div class="max-w-lg w-full lg:max-w-xs">
              <div class="relative">
                <Search />
              </div>
            </div>
          </div>
          <div class="flex lg:hidden">
            <!-- Mobile menu button -->
            <button
              type="button"
              class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
              aria-controls="mobile-menu"
              aria-expanded="false"
              @click="showSettingsMenu"
            >
              <span class="sr-only">Open main menu</span>
              <!-- Icon when menu is closed. -->
              <!--
                Heroicon name: outline/menu

                Menu open: "hidden", Menu closed: "block"
              -->
              <div v-if="!activeSettingsMenu">
                <svg
                  class="block h-6 w-6"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  aria-hidden="true"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                  />
                </svg>
              </div>
              <div v-else>
                <svg
                  class="hidden h-6 w-6"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  aria-hidden="true"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </div>
            </button>
          </div>

          <div class="hidden lg:block lg:ml-4">
            <div class="flex items-center">
              <!-- Settings dropdown -->
              <client-only>
                <div class="ml-4 relative flex-shrink-0">
                  <div>
                    <FlyoutMenu @userSubscribed="showNotifications = true" v-if="!cookieSet" />
                    <button
                      id="user-menu"
                      class="bg-gray-800 rounded-full flex text-sm text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white"
                      aria-haspopup="true"
                      @click="showSettingsMenu"
                    >
                      <span class="sr-only">Open user menu</span>
                      <div class="flex-shrink-0">
                        <img
                          class="h-10 w-10 rounded-full"
                          :src="getAvatar"
                          alt="Avatar"
                        >
                      </div>
                    </button>
                  </div>
                  <!--
                    Settings dropdown panel, show/hide based on dropdown state.

                    Entering: "transition ease-out duration-100"
                      From: "transform opacity-0 scale-95"
                      To: "transform opacity-100 scale-100"
                    Leaving: "transition ease-in duration-75"
                      From: "transform opacity-100 scale-100"
                      To: "transform opacity-0 scale-95"
                  -->
                  <div
                    v-if="activeSettingsMenu"
                    v-on-clickaway="showSettingsMenu"
                    class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5"
                    role="menu"
                    aria-orientation="vertical"
                    aria-labelledby="user-menu"
                    @click="activeSettingsMenu = false"
                  >
                    <nuxt-link
                      to="/"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      role="menuitem"
                    >
                      Home
                    </nuxt-link>
                    <nuxt-link
                      to="/about"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      role="menuitem"
                    >
                      About
                    </nuxt-link>
                    <nuxt-link
                      to="/settings"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      role="menuitem"
                    >
                      Settings
                    </nuxt-link>
                    <nuxt-link
                      to="/profile"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      role="menuitem"
                    >
                      Profile
                    </nuxt-link>
                    <div class="relative">
                      <div
                        class="absolute inset-0 flex items-center"
                        aria-hidden="true"
                      >
                        <div class="w-full border-t border-gray-300" />
                      </div>
                    </div>
                    <nuxt-link
                      to="/privacy"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      role="menuitem"
                    >
                      Privacy Policy
                    </nuxt-link>
                    <nuxt-link
                      to="/tos"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      role="menuitem"
                    >
                      Terms of Service
                    </nuxt-link>
                    <div class="relative">
                      <div
                        class="absolute inset-0 flex items-center"
                        aria-hidden="true"
                      >
                        <div class="w-full border-t border-gray-300" />
                      </div>
                    </div>
                    <div v-if="isLoggedIn">
                      <a
                        href="#"
                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        role="menuitem"
                        @click="logout"
                      >Sign out</a>
                    </div>
                    <div v-else>
                      <a
                        href="/signin"
                        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        role="menuitem"
                      >Sign In</a>
                    </div>
                  </div>
                </div>
                <div class="ml-2">
                  <client-only>
                    <button
                      class="ml-auto flex-shrink-0 bg-gray-800 p-1 rounded-full text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white"
                      @click="toggleNotify()"
                    >
                      <span class="sr-only">View notifications</span>
                      <!-- Heroicon name: outline/bell -->
                      <div v-if="showNotifications">
                        <svg
                          class="h-6 w-6"
                          xmlns="http://www.w3.org/2000/svg"
                          fill="white"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                          aria-hidden="true"
                          style="color: rgb(239, 243, 244);"
                        >
                          <path d="M23.61.15c-.375-.184-.822-.03-1.006.34L19.74 6.266l-1.703-1.81c-.283-.303-.758-.316-1.06-.032-.302.284-.316.76-.032 1.06l2.443 2.596c.143.15.34.235.546.235.036 0 .073-.003.11-.008.243-.036.452-.19.562-.41l3.342-6.74c.184-.372.032-.822-.34-1.006zm-4.592 16.475c-.083-.064-2.044-1.625-2.01-5.76.022-2.433-.78-4.596-2.256-6.09-1.324-1.34-3.116-2.083-5.046-2.092h-.013c-1.93.01-3.722.75-5.046 2.092C3.172 6.27 2.37 8.433 2.39 10.867 2.426 15 .467 16.56.39 16.62c-.26.193-.367.53-.266.838.102.308.39.515.712.515h4.08c.088 2.57 2.193 4.64 4.785 4.64s4.698-2.07 4.785-4.64h4.082c.32 0 .604-.206.707-.51s-.002-.643-.256-.838zM9.7 20.513c-1.434 0-2.6-1.127-2.684-2.54h5.368c-.085 1.413-1.25 2.54-2.684 2.54z" />
                        </svg>
                      </div>
                      <div v-else>
                        <IconClient
                          icon="ic:twotone-notification-add"
                          icon-class="h-6 w-6"
                        />
                      </div>
                    </button>
                  </client-only>
                </div>
                <div
                  v-if="!isLoggedIn"
                  class="ml-4 relative flex-shrink-0"
                >
                  <div>
                    <nuxt-link to="/signin">
                      <button
                        type="button"
                        class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                      >
                        Sign In
                      </button>
                    </nuxt-link>
                  </div>
                </div>
              </client-only>
            </div>
          </div>
        </div>
      </div>

      <!-- Mobile menu, show/hide based on menu state. -->
      <div v-show="activeSettingsMenu">
        <div
          id="mobile-menu"
          class="lg:hidden bg-gray-200"
        >
          <div class="pt-4 pb-3 border-t border-gray-700">
            <div class="flex items-center px-5">
              <div
                v-show="isLoggedIn"
                v-if="getAvatar"
                class="flex-shrink-0"
              >
                <img
                  class="h-10 w-10 rounded-full"
                  :src="getAvatar"
                  alt="Avatar"
                >
              </div>
              <client-only>
                <div
                  v-if="isLoggedIn"
                  class="ml-3"
                >
                  <div v-if="authUser.email">
                    <div class="text-sm font-medium text-gray-800">{{ authUser.email }}</div>
                  </div>
                </div>
              </client-only>
            </div>
            <div
              class="mt-3 px-2 space-y-1"
              @click="activeSettingsMenu = false"
            >
              <client-only>
                <div>
                  <button
                    class="ml-auto flex-shrink-0 bg-gray-800 p-1 rounded-full text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white"
                    @click="toggleNotify()"
                  >
                    <span class="sr-only">View notifications</span>
                    <!-- Heroicon name: outline/bell -->
                    <div v-if="showNotifications">
                      <svg
                        class="h-6 w-6"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="white"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        aria-hidden="true"
                        style="color: rgb(239, 243, 244);"
                      >
                        <path d="M23.61.15c-.375-.184-.822-.03-1.006.34L19.74 6.266l-1.703-1.81c-.283-.303-.758-.316-1.06-.032-.302.284-.316.76-.032 1.06l2.443 2.596c.143.15.34.235.546.235.036 0 .073-.003.11-.008.243-.036.452-.19.562-.41l3.342-6.74c.184-.372.032-.822-.34-1.006zm-4.592 16.475c-.083-.064-2.044-1.625-2.01-5.76.022-2.433-.78-4.596-2.256-6.09-1.324-1.34-3.116-2.083-5.046-2.092h-.013c-1.93.01-3.722.75-5.046 2.092C3.172 6.27 2.37 8.433 2.39 10.867 2.426 15 .467 16.56.39 16.62c-.26.193-.367.53-.266.838.102.308.39.515.712.515h4.08c.088 2.57 2.193 4.64 4.785 4.64s4.698-2.07 4.785-4.64h4.082c.32 0 .604-.206.707-.51s-.002-.643-.256-.838zM9.7 20.513c-1.434 0-2.6-1.127-2.684-2.54h5.368c-.085 1.413-1.25 2.54-2.684 2.54z" />
                      </svg>
                    </div>
                    <div v-else>
                      <IconClient
                        icon="ic:twotone-notification-add"
                        icon-class="h-6 w-6"
                      />
                    </div>
                  </button>
                  <nuxt-link
                    to="/"
                    class="block px-3 py-2 rounded-md text-base font-medium text-gray-800 hover:text-white hover:bg-gray-700"
                  >
                    Home
                  </nuxt-link>
                  <nuxt-link
                    to="/about"
                    class="block px-3 py-2 rounded-md text-base font-medium text-gray-800 hover:text-white hover:bg-gray-700"
                  >
                    About
                  </nuxt-link>
                  <nuxt-link
                    to="/settings"
                    class="block px-3 py-2 rounded-md text-base font-medium text-gray-800 hover:text-white hover:bg-gray-700"
                  >
                    Settings
                  </nuxt-link>
                  <nuxt-link
                    to="/profile"
                    class="block px-3 py-2 rounded-md text-base font-medium text-gray-800 hover:text-white hover:bg-gray-700"
                  >
                    Profile
                  </nuxt-link>
                  <div v-if="isLoggedIn">
                    <a
                      href="#"
                      class="block px-3 py-2 rounded-md text-base font-medium text-gray-800 hover:text-white hover:bg-gray-700"
                      @click="logout"
                    >Sign out</a>
                  </div>
                  <div v-else>
                    <nuxt-link
                      to="/signin"
                      class="block px-3 py-2 rounded-md text-base font-medium text-gray-800 hover:text-white hover:bg-gray-700"
                    >
                      Sign In
                    </nuxt-link>
                  </div>
                </div>
              </client-only>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <div class="sm:flex items-center justify-center sm:space-x-5">
      <div class="lg:mt-4 lg:mb-5 flex-shrink-0 items-center justify-center lg:space-y-4">
        <Nuxt />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
declare global {
  interface Window {
    registration:any;
  }
}

import Vue from 'vue'
// import { signOut, getAuth, onAuthStateChanged } from 'firebase/auth'
import { signOut } from 'firebase/auth'
import { mapState, mapGetters } from 'vuex'
import { AuthUser } from '@/types'
import { auth, db, firebaseApp } from '~/plugins/firebase'

const topic = "chases-notifications"

export default Vue.extend({
  data () {
    return {
      notificationBtn: { style: '' },
      activeSettingsMenu: false,
      showNotifications: false,
      clicked: false,
      loading: {
        logout: false,
      },
      preload: true,
      cookieSet: false,
    }
  },
  head () {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-return
    return {
      script: [
        {
          hid: 'pusherBeams' as string,
          type: 'text/javascript' as string,
          src: "https://js.pusher.com/beams/2.0.0-beta.1/push-notifications-cdn.js" as string,
          defer: true as boolean,
          callback: () => { this.notifyStartup() },
        },
      ],
    } as any
  },
  computed: {
    ...mapState({
      authUser: state => state.authUser as AuthUser,
    }),
    ...mapGetters({
      isLoggedIn: 'isLoggedIn',
    }),
    getAvatar () {
      return this.authUser?.photoURL || '/user.png'
    },
  },
  mounted () {
    this.preload = false
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    this.cookieSet = this.$cookies.get('notify') as boolean
    this.showNotifications = this.cookieSet
  },
  methods: {
    notifyStartup () {
      console.log('notifyStartup')
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
      if (this.showNotifications) {
        console.log('in showNotifications')
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-assignment,no-undef,@typescript-eslint/no-unsafe-call
        const beamsClient = new PusherPushNotifications.Client({
          // TODO: change this to take a process.env and make it discern from a production/staging key
          instanceId: '4430414d-cce4-4722-9586-f32db3d7d433',
        });
        console.log('Starting beamsClient')
        // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
        beamsClient.start()
            // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-return
            .then(() => beamsClient.addDeviceInterest(topic))
            .then(() => console.log(`Registered to ${topic}`))
            .catch(console.error);

        console.log('Waiting for new notifications')
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,no-undef
        PusherPushNotifications.onNotificationReceived = ({pushEvent, payload}: any) => {
          // logEvent(analytics, "notification_foreground")
          // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
          pushEvent.waitUntil(
              // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
              self.registration.showNotification(payload.notification.title, {
                // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
                body: payload.notification.body,
                // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
                icon: payload.notification.icon,
                // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
                data: payload.data,
              }),
          )
        }
      }
    },
    checkiOS () {
      return [
        'iPad Simulator',
        'iPhone Simulator',
        'iPod Simulator',
        'iPad',
        'iPhone',
        'iPod',
      ].includes(navigator.platform) ||
          // iPad on iOS 13 detection
          (navigator.userAgent.includes('Mac') && 'ontouchend' in document)
    },
    isSafari () {
      return (navigator.vendor.match(/apple/i) || '').length > 0
    },
    toggleNotify () {
      if (this.showNotifications) {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
        this.$toast.success('Disabling notifications')
        this.showNotifications = !this.showNotifications
        // TODO: make this do the beamsClient.removeDeviceInterest
        // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
        this.$cookies.set('notify', 'false', {
          path: '/',
          maxAge: 60 * 60 * 24 * 7,
        })
      } else {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        this.$toast.success('Enabling notifications')
        if (navigator.maxTouchPoints && navigator.userAgent.includes('Safari') && !navigator.userAgent.includes('Chrome')) {
          // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
          this.$toast.denied('Push notifications not supported in iOS')
          this.showNotifications = false
          return
        }
        console.log('Toggled')
        this.notifyStartup()
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        this.$cookies.set('notify', 'true', {
          path: '/',
          maxAge: 60 * 60 * 24 * 7,
        })
        this.showNotifications = !this.showNotifications
      }
    },
    logout (): void {
      signOut(auth).then(() => {
        this.activeSettingsMenu = false
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        this.$router.push('/signin')
      }).catch((e) => {
        console.error(e)
      })
    },
    showSettingsMenu () {
      this.activeSettingsMenu = !this.activeSettingsMenu
    },
  },
})
</script>

<style>
  .loading-page {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(255, 255, 255, 0.8);
    text-align: center;
    padding-top: 200px;
    font-size: 30px;
    font-family: sans-serif;
  }
  html,body {
    @apply bg-gray-400
  }
</style>
