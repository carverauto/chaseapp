<template>
  <div class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <!--
        Background overlay, show/hide based on modal state.

        Entering: "ease-out duration-300"
          From: "opacity-0"
          To: "opacity-100"
        Leaving: "ease-in duration-200"
          From: "opacity-100"
          To: "opacity-0"
      -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" />

      <!-- This element is to trick the browser into centering the modal contents. -->
      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <!--
        Modal panel, show/hide based on modal state.

        Entering: "ease-out duration-300"
          From: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          To: "opacity-100 translate-y-0 sm:scale-100"
        Leaving: "ease-in duration-200"
          From: "opacity-100 translate-y-0 sm:scale-100"
          To: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
      -->
      <div class="inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6">
        <div>
          <div class="mx-auto flex items-center justify-center h-24 w-24 rounded-full bg-red-100">
            <Lottie ref="sharing" :lottie-settings="lottieSettings" />
          </div>
          <div class="mt-3 text-center sm:mt-5">
            <h3 id="modal-title" class="text-lg leading-6 font-medium text-gray-900">
              Sharing is Caring..
            </h3>
            <div class="mt-2">
              <p class="text-sm text-gray-500">
                Give the gift of ChaseApp today.
              </p>
            </div>
          </div>
        </div>
        <div id="shareNetworks" class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense">
          <div id="shareTwitter">
            <ShareNetwork :title="chase.Name" :url="chaseURL" network="twitter" hashtags="#pursuit">
              <button type="button" aria-labelledby="shareNetworks" aria-label="Share on Twitter" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-2 sm:text-sm" @click="updateAnalytics(share)">
                <IconClient icon="openmoji:twitter" icon-class="h-6 w-6" />
                <div class="ml-2">
                  Twitter
                </div>
              </button>
            </ShareNetwork>
          </div>
         <div id="shareFacebook">
           <ShareNetwork :title="chase.Name" :url="chaseURL" network="facebook" hashtags="#Pursuit">
             <button type="button" aria-label="Share on Facebook" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-2 sm:text-sm" @click="updateAnalytics(share)">
               <IconClient icon="openmoji:facebook" icon-class="h-6 w-6" />
               <div class="ml-2">
                 Facebook
               </div>
             </button>
           </ShareNetwork>
         </div>
          <div id="shareTelegram">
            <ShareNetwork :title="chase.Name" :url="chaseURL" network="telegram">
              <button type="button" aria-label="Share on Telegram" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-1 sm:text-sm" @click="updateAnalytics(share)">
                <IconClient icon="logos:telegram" icon-class="h-6 w-6" />
                <div class="ml-2">
                  Telegram
                </div>
              </button>
            </ShareNetwork>
          </div>
          <div>
            <ShareNetwork :title="chase.Name" :url="chaseURL" network="reddit">
              <button type="button" aria-label="Share on Reddit" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-1 sm:text-sm" @click="updateAnalytics(share)">
                <IconClient icon="logos:reddit-icon" icon-class="h-6 w-6" />
                <div class="ml-2">
                  Reddit
                </div>
              </button>
            </ShareNetwork>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

// eslint-disable-next-line import/named
import Vue, { PropType } from 'vue'
import { Chase } from '~/types'

// import { analytics } from "~/plugins/firebaseAnalytics.client"
// import { logEvent } from "firebase/analytics"

export default Vue.extend({
  props: {
    chase: {
      type: Object as PropType<Chase>,
      required: true,
    },
  },
  data () {
    return {
      lottieSettings: { path: '/json/christmas-gifts.json' },
    }
  },
  computed: {
    chaseURL (): string {
      return 'https://chaseapp.tv/chase/' + this.chase.ID
    },
  },
  methods: {
    updateAnalytics(event: string) {
      // logEvent(analytics, event)
    }
  }
})
</script>
