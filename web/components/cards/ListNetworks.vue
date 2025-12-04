<template>
  <div class="flex w-full min-w-0 pb-4 md:pb-3 lg:pb-2">
    <div class="w-6 md:w-10 mr-2 md:mr-4 flex-shrink-0">
      <div v-if="network.Logo">
        <img class="border border-gray-200 h-6 w-6 md:h-10 md:w-10 rounded-md" :alt="network.Name" :src="checkAvatars(network.Logo)">
      </div>
      <div v-else>
        <img v-if="checkAvatars(network.URL)" class="border border-gray-200 h-6 w-6 md:h-10 md:w-10 rounded-md" :alt="network.Name" :src="checkAvatars(network.URL)">
        <Placeholder v-else primary="#FFFFFF" secondary="#000000" class="border border-gray-200 h-6 w-6 md:h-10 md:w-10 rounded-md fill-current text-gray-300" />
      </div>
    </div>
    <div class="align-middle min-w-0 flex flex-col">
      <div class="align-middle flex font-bold text-gray-600 leading-none">
        <span class="py-0.5 inline-flex">{{ network.Name }}</span>
        <span v-if="network.Tier === 1" class="ml-1 items-center px-2 py-0.5 rounded text-xs font-bold bg-blue-100 text-blue-800">
          Primary
        </span>
        <span v-else-if="network.Tier === 2" class="ml-1 inline-flex items-center px-2 py-0.5 rounded text-xs font-bold bg-yellow-100 text-yellow-800">
          Sponsored
        </span>
        <span v-else-if="network.Tier < 0" class="ml-1 inline-flex items-center px-2 py-0.5 rounded text-xs font-bold bg-red-100 text-red-800">
          Dead
        </span>
      </div>
      <div class="flex align-text-bottom items-center min-w-1 max-w-full">
        <a class="truncate text-orange-600 dark:text-orange-400 hover:underline" :href="network.URL" rel="noopener" target="_blank">
          {{ network.URL }}
        </a>
        <span class="fas fa-external-link-alt leading-0 ml-2 text-orange-600 dark:text-orange-400 text-xs mr-1" />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ListNetworks',
  middleware: 'auth',
  props: {
    network: {
      type: Object,
      required: true,
    },
  },
  data: () => ({
    avatar_urls: {
      'foxla.com': '/networks/foxla.jpg',
      'losangeles.cbslocal.com': '/networks/cbsla.jpg',
      'nbclosangeles.com': '/networks/nbclosangeles.jpg',
      'abc7.com': '/networks/abc7.jpg',
      'ktla.com': '/networks/ktla.jpg',
      'broadcastify.com': '/networks/broadcastify.jpg',
      'audio12.broadcastify.com': '/networks/broadcastify.jpg',
      'facebook.com': '/networks/facebook.jpg',
      'm.facebook.com': '/networks/facebook.jpg',
      'wsvn.com': '/networks/wsvn.jpg',
      'twitter.com': '/networks/twitter.svg',
      'nbcmiami.com': '/networks/nbc6.jpg',
      'iheart.com': '/networks/iheart.jpg',
      'kdoc.tv': '/networks/kdoc.jpg',
      'youtube.com': '/networks/youtube.jpg',
      'youtu.be': '/networks/youtube.jpg',
      'pscp.tv': '/networks/periscope.png',
      'fox10phoenix.com': '/networks/fox10.png',
      'kfor.com': '/networks/kfor.jpg',
      'kctv5.com': '/networks/kctv5.jpg',
      'fox4news.com': '/networks/fox4.jpg',
      'nbcdfw.com': '/networks/nbcdfw.jpg',
      'khou.com': '/networks/khou.jpg',
      'wesh.com': '/networks/wesh.jpg',
      'cbs8.com': '/networks/cbs8.png',
      'abc13.com': '/networks/abc13.png',
      'wsoctv.com': '/networks/wsoc.png',
      'news9.com': '/networks/news9.jpg',
    },
  }),
  methods: {
    checkAvatars (toCheck) {
      try {
        return this.avatar_urls[new URL(toCheck).hostname.replace('www.', '')]
      } catch {
        return this.avatar_urls[toCheck]
      }
    },
  },
}
</script>

<style scoped>

</style>
