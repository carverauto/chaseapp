<template>
  <div class="relative">
    <!-- Item active: "text-gray-900", Item inactive: "text-gray-500" -->
    <!-- <button @click="showNotifyList" v-on-clickaway="showNotifyMenu" class="flex-shrink-0 bg-gray-800 p-1 rounded-full text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white"> -->
    <!-- <button type="button" class="text-gray-500 group bg-gray-800 rounded-md inline-flex items-center text-base font-medium hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500" aria-expanded="false"> -->
    <button @click="showNotifyMenu" class="flex-shrink-0 bg-gray-800 p-1 rounded-full text-gray-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white" aria-expanded="false">
      <span class="sr-only">View notifications</span>
      <!-- Heroicon name: outline/bell -->
      <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
      </svg>
    </button>

    <div v-if="activeNotifyMenu">
      <transition
          enter-active-class="transition ease-out duration-200"
          enter-from-class="opacity-0 translate-y-1"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition ease-in duration-150"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 translate-y-1"
      >
        <div class="absolute z-10 left-1/2 transform -translate-x-1/2 mt-3 px-2 w-screen max-w-xs sm:px-0">
          <div class="rounded-lg shadow-lg ring-1 ring-black ring-opacity-5 overflow-hidden">
            <div class="relative grid gap-6 bg-white px-5 py-6 sm:gap-8 sm:p-8">
              <a href="#" class="-m-3 p-3 block rounded-md hover:bg-gray-50 transition ease-in-out duration-150">
                <p class="text-base font-medium text-gray-900">
                  Notification
                </p>
                <p class="mt-1 text-sm text-gray-500">
                  Received: 1 day ago
                </p>
              </a>

              <a href="#" class="-m-3 p-3 block rounded-md hover:bg-gray-50 transition ease-in-out duration-150">
                <p class="text-base font-medium text-gray-900">
                  Notification
                </p>
                <p class="mt-1 text-sm text-gray-500">
                  Received: 3 days ago
                </p>
              </a>

            </div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script>
export default {
  name: "NotifyList",
  data () {
    return {
      activeNotifyMenu: this.showNotify
    }
  },
  props: {
    showNotify: {
      type: Boolean,
      required: true,
    },
  },
  methods: {
    showNotifyMenu() {
      this.activeNotifyMenu = !this.activeNotifyMenu
    },
    notifyMenuOff() {
      this.activeNotifyMenu = false
    },
  }
}
</script>

<style scoped>

</style>
