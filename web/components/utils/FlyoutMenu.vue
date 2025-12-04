<template>
  <div v-on-clickaway="hideMenu">
    <div v-if="closeMenu"  class="relative">
      <!-- Item active: "text-gray-900", Item inactive: "text-gray-500" -->
      <!--
        Flyout menu, show/hide based on flyout menu state.

        Entering: "transition ease-out duration-200"
          From: "opacity-0 translate-y-1"
          To: "opacity-100 translate-y-0"
        Leaving: "transition ease-in duration-150"
          From: "opacity-100 translate-y-0"
          To: "opacity-0 translate-y-1"
      -->
      <div class="absolute z-10 left-1/2 transform -translate-x-1/2 mt-3 px-2 w-screen max-w-xs sm:px-0">
        <div class="rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 overflow-hidden">
          <div class="relative grid gap-6 bg-white px-5 py-6 sm:gap-8 sm:p-8">
            <a href="#" class="-m-3 p-3 block rounded-md hover:bg-gray-50 transition ease-in-out duration-150">
              <p class="text-base font-medium text-gray-900">Notifications</p>
              <p class="mt-1 text-sm text-gray-500">Would you like to signup for free notfications?</p>
              <button type="button" @click="enableNotifications" class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white bg-blue-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">Notify Me</button>
              <button type="button" @click="closeMenu = !closeMenu" class="inline-flex items-center mt-2 px-2.5 py-1.5 border border-gray-300 shadow-sm text-xs font-medium rounded text-gray-700 bg-red-500 hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">No Thanks</button>
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

const topic = "chases-notifications"

export default {
  name: "FlyoutMenu",
  data () {
    return {
      closeMenu: true,
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
  methods: {
    hideMenu (): void {
      this.closeMenu = false
    },
    notifyStartup () {
      console.log('notifyStartup')
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
      this.cookieSet = this.$cookies.get('notify') as boolean
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

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,no-undef
      PusherPushNotifications.onNotificationReceived = ({pushEvent, payload}: any) => {
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
    },
    enableNotifications (): void {
      if (navigator.maxTouchPoints && navigator.userAgent.includes('Safari') && !navigator.userAgent.includes('Chrome')) {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,@typescript-eslint/ban-ts-comment
        // @ts-ignore
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        this.$toast.denied('Push notifications not supported in iOS')
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
      this.$emit('userSubscribed', true)
      this.closeMenu = false
      },
    },
}
</script>

<style scoped>

</style>