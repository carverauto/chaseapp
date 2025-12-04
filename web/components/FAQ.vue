<template>
  <div class="bg-gray-50">
    <div class="max-w-7xl mx-auto py-12 px-4 sm:py-16 sm:px-6 lg:px-8">
      <div class="max-w-3xl mx-auto divide-y-2 divide-gray-200">
        <h2 class="text-center text-3xl font-extrabold text-gray-900 sm:text-4xl">
          Frequently asked questions
        </h2>
        <dl class="mt-6 space-y-6 divide-y divide-gray-200">
          <div class="pt-6">
            <dt class="text-lg">
              <!-- Expand/collapse question button -->
              <button class="text-left w-full flex justify-between items-start text-gray-400">
              <span class="font-medium text-gray-900">
                What is ChaseApp?
              </span>
                <span class="ml-6 h-7 flex items-center">
                <!--
                  Expand/collapse icon, toggle lib based on question open state.

                  Heroicon name: outline/chevron-down

                  Open: "-rotate-180", Closed: "rotate-0"
                -->
                <svg class="-rotate-180 h-6 w-6 transform" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </span>
              </button>
            </dt>
            <dd class="flex mt-2 pr-12 justify-start">
              <p class="text-base text-gray-500">
                <!-- <img class="mr-2 mb-1 mt-1 inline-block h-24 w-auto rounded-md float-left" src="/police.gif" /> -->
                <Donut class="mr-4 mb-1 mt-1 inline-block h-20 w-auto rounded-md float-left" />
                ChaseApp is a social platform for delivering real-time chase notifications and premium content to chase enthusiasts and people that are just generally interested in live events and news.
              </p>
            </dd>
          </div>

          <div class="pt-6">
            <dt class="text-lg">
              <!-- Expand/collapse question button -->
              <button class="text-left w-full flex justify-between items-start text-gray-400">
              <span class="font-medium text-gray-900">
                How does it work?
              </span>
                <span class="ml-6 h-7 flex items-center">
                <!--
                  Expand/collapse icon, toggle lib based on question open state.

                  Heroicon name: outline/chevron-down

                  Open: "-rotate-180", Closed: "rotate-0"
                -->
                <svg class="-rotate-180 h-6 w-6 transform" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </span>
              </button>
            </dt>
            <dd class="flex mt-2 pr-12 lg:items-center">
              <p class="text-base text-gray-500">
                <img class="ml-1 mb-1 mt-1 inline-block h-24 w-auto rounded-md float-right" src="/hackers.gif" />
                Artificial Intelligence and Machine Learning algorithms combine with a human operations and coordination team to sift through data, finding new chases and live streams.
              </p>
            </dd>
          </div>

          <div class="pt-6">
            <dt class="text-lg">
              <!-- Expand/collapse question button -->
              <button class="text-left w-full flex justify-between items-start text-gray-400">
              <span class="font-medium text-gray-900">
                How is this better than Twitter?
              </span>
                <span class="ml-6 h-7 flex items-center">
                <!--
                  Expand/collapse icon, toggle lib based on question open state.

                  Heroicon name: outline/chevron-down

                  Open: "-rotate-180", Closed: "rotate-0"
                -->
                <svg class="-rotate-180 h-6 w-6 transform" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </span>
              </button>
            </dt>
            <dd class="flex mt-2 pr-12 items-center">
              <p class="text-base text-gray-500">
                <ChaseWheels w1="fill:#FF0000" w2="fill:#5cda3f" class="mr-4 mt-1 mb-1 h-16 w-auto float-left"/>
                We look at multiple sources to find you new chases, adding useful content, voting, chat, and other elements like the ChaseWheels &trade; tracker, ChaseStats &trade;, and Birds of a Feather (BoF) &trade; helicopter tracking system.
              </p>
            </dd>
          </div>
        </dl>
      </div>
    </div>
  </div>
</template>

<script>
export default {
name: "FAQ"
}
</script>

<style scoped>

</style>
