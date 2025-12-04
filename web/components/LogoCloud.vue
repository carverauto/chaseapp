<template>
  <!-- This example requires Tailwind CSS v2.0+ -->
  <div class="bg-white">
    <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:py-16 lg:px-8">
      <p class="text-center text-base font-semibold uppercase text-gray-600 tracking-wider">
        Bringing you LIVE video feeds from over 25 different networks and social media platforms
      </p>
      <div class="mt-6 grid grid-cols-2 gap-0.5 md:grid-cols-3 lg:mt-8">
      <div v-for="(network, $index) in Networks" :key="$index" :network="network">
        <div class="col-span-1 flex justify-center py-8 px-8 bg-gray-50">
          <img class="max-h-12" :src="network" alt="Network">
        </div>
      </div>
        </div>
    </div>
  </div>
</template>

<script>
const myNetworks = [
    '/networks/abc7.jpg', '/networks/broadcastify.jpg',
    '/networks/cbsla.jpg', '/networks/facebook.jpg', '/networks/fox4.jpg',
    '/networks/fox10.png', '/networks/foxla.jpg', '/networks/iheart.jpg',
    '/networks/kctv5.jpg', '/networks/kdoc.jpg', '/networks/kfor.jpg',
    '/networks/khou.jpg', '/networks/ktla.jpg', '/networks/nbc6.jpg',
    '/networks/nbcdfw.jpg', '/networks/nbclosangeles.jpg', '/networks/wesh.jpg',
    '/networks/wsvn.jpg',
]

export default {
  name: "LogoCloud",
  data() {
    return {
      Networks: myNetworks,
    }
  },
  mounted () {
  },
  methods: {
  }
}
</script>

<style scoped>

</style>