<template>
  <div class="border overflow-y-scroll h-48 max-w-3xl px-4 py-4 sm:px-6 lg:rounded-lg z-20 lg:shadow dark:bg-gray-800 bg-white lg:max-w-4xl lg:mb-4">
    <div class="flow-root">
      <ul class="-mb-8">
        <li v-for="index in 4" :key="index">
          <div class="relative pb-8">
            <span class="absolute top-4 left-4 -ml-px h-full bg-gray-200" aria-hidden="true" />
            <div class="flex items-center justify-between skeleton">
              &nbsp;
              <div class="hidden lg:block ml-3">
                <p class="text-sm font-sans font-medium text-gray-800 group-hover:text-gray-900 skeleton">&nbsp;</p>
                <p class="text-xs font-sans font-medium text-gray-500 group-hover:text-gray-700 skeleton">&nbsp;</p>
              </div>
              <div class="min-w-0 flex-1 pt-1.5 justify-between space-x-4">
                <div>
                  <p class="font-medium font-sans text-gray-900 sm:truncate ml-2 skeleton">&nbsp;</p>
                  <div class="absolute top-0 right-0 skeleton">
                    &nbsp;
                  </div>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<style>
.skeleton {
  @apply text-opacity-0;
  background-image: linear-gradient(100deg, #edf2f7 0%, #f4f7fa 20%, #edf2f7 40%);
  background-position: 50%;
  background-size: 200%;
  animation: skeleton 1.5s infinite linear;
}
@keyframes skeleton {
  0% {
    background-position: 50%;
  }
  50%, 100% {
    background-position: -90%
  }
}
</style>
