<template>
  <div class="mt-6 m-2 mb-5">
    <form class="space-y-8 divide-y divide-gray-200">
      <div class="space-y-8 divide-y divide-gray-200">
        <div class="pt-0">
          <div>
            <h3 class="text-lg leading-6 font-medium text-gray-900">
              Notifications
            </h3>
            <p class="mt-1 text-sm text-gray-500 mb-3">
              Get real-time notifications, never miss another live chase,
              tweet, or stream for an important or viral event
            </p>
          </div>

          <div>
            <div v-if="isNotSupported">
              <p>
                We're sorry, but push notifications are not supported on iOS at this time.
                This is a limitation imposed by Apple. Please stay tuned as our App
                will be in the Stores soon.
              </p>
            </div>
            <div v-else>
              <div v-if="beamsLoaded">
                <fieldset class="space-y-5">
                  <legend class="sr-only">Notifications</legend>
                  <div class="relative flex items-start">
                    <div class="flex items-center h-5">
                      <input id="chases" v-model="showNotifications" aria-describedby="chases-description" @click="addDeviceInterest('chases-notifications')" name="chases" type="checkbox" class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded">
                    </div>
                    <div class="ml-3 text-sm">
                      <label for="chases" class="font-medium text-gray-700">Chases</label>
                      <p id="chases-description" class="text-gray-500">Get notified when a live chase occurs</p>
                    </div>
                  </div>
                  <div class="relative flex items-start">
                    <div class="flex items-center h-5">
                      <input id="firehose" v-model="firehoseNotifications" aria-describedby="firehose-description" @click="addDeviceInterest('firehose-notifications')" nname="firehose" type="checkbox" class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded">
                    </div>
                    <div class="ml-3 text-sm">
                      <label for="firehose" class="font-medium text-gray-700">Firehose</label>
                      <p id="firehose-description" class="text-gray-500">Get notified when we add something to the Firehose</p>
                    </div>
                  </div>
                </fieldset>
              </div>
            </div>
          </div>
        </div>

        <div>
          <div>
            <h3 class="text-lg leading-6 font-medium text-gray-900 mt-2">
              Account
            </h3>
            <p class="mt-1 text-sm text-gray-500 mb-2">
              Account settings
            </p>
          </div>

          <div class="mt-6">
            <!-- <UserSubscription /> -->
            <UserDeleteAccount />
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script lang="ts">
import { signOut } from 'firebase/auth'
import Vue from 'vue'
import { mapGetters, mapState } from 'vuex'
import { auth } from '~/plugins/firebase'
import { AuthUser } from '~/types'

const topic = "chases-notifications"

export default Vue.extend({
  name: 'Profile',
  middleware: 'auth',
  data () {
    return {
      activeSettingsMenu: false,
      showNotifications: false,
      firehoseNotifications: false,
      beamsLoaded: false,
      loading: {
        logout: false,
      },
      cookieSet: false,
      isNotSupported: false,
      subscriptions: [],
    }
  },
  head () {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-return
    return {
      script: [
        {
          hid: 'pusherBeams',
          type: 'text/javascript',
          src: "https://js.pusher.com/beams/2.0.0-beta.1/push-notifications-cdn.js",
          defer: false,
          callback: () => {
            this.loadBeams()
          },
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
    // See if a cookie is set for notifications
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
    this.cookieSet = this.$cookies.get('notify')
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    this.showNotifications = this.cookieSet
    if (navigator.maxTouchPoints && navigator.userAgent.includes('Safari') && !navigator.userAgent.includes('Chrome'))
      this.isNotSupported = true
  },
  methods: {
    setNotify(interest: string) {
      console.log(`Interested in ${interest}`)
    },
    isSafari () {
      return (navigator.vendor.match(/apple/i) || '').length > 0
    },
    logout (): void {
      signOut(auth).then(async () => {
        this.activeSettingsMenu = false
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        await this.$router.push('/signin')
      }).catch((e) => {
        console.error(e)
      })
    },
    loadBeams () {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access,no-undef
      const beamsClient = new PusherPushNotifications.Client({
        // TODO: change this to take a process.env and make it discern from a production/staging key
        instanceId: '4430414d-cce4-4722-9586-f32db3d7d433',
      });
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      beamsClient.start()
          // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-return
          .then(() =>
              // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-return
              beamsClient.getDeviceInterests().then((interest: string) => {
                if (interest.includes("firehose-notifications"))
                  this.firehoseNotifications = true
              }))
          .catch(console.error);
      this.beamsLoaded = true
    },
    showSettingsMenu () {
      this.activeSettingsMenu = !this.activeSettingsMenu
    },
    addDeviceInterest (interest: string) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access,no-undef
      const beamsClient = new PusherPushNotifications.Client({
        // TODO: change this to take a process.env and make it discern from a production/staging key
        instanceId: '4430414d-cce4-4722-9586-f32db3d7d433',
      });
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      beamsClient.start()
          // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-return
          .then(() => beamsClient.addDeviceInterest(interest))
          .then(() => console.log(`Registered to ${interest}`))
          .catch(console.error);
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      this.$toast.show(`Subscribed to ${interest}`)
    },
    removeDeviceInterest(interest: string) {
      console.log('removeDeviceFromInterest')
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access,no-undef
      const beamsClient = new PusherPushNotifications.Client({
        // TODO: change this to take a process.env and make it discern from a production/staging key
        instanceId: '4430414d-cce4-4722-9586-f32db3d7d433',
      });
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      beamsClient.start()
          // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-return
          .then(() => beamsClient.removeDeviceInterest(interest))
          .then(() => console.log(`Registered to ${interest}`))
          .catch(console.error);
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      this.$toast.show(`Subscribed to ${interest}`)
    },
  },
})
</script>
