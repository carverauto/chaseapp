<template>
  <div>
    <div class="max-w-3xl lg:min-w-3xl px-4 py-4 sm:px-6 lg:max-w-4xl">
      <div class="flex items-center justify-between">
        <p class="skeleton w-48">
          &nbsp;
        </p>
        <div class="ml-2 flex-shrink-0 flex">
          <p class="px-2 inline-flex leading-5 rounded-full mr-4">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-md uppercase skeleton w-12">
                &nbsp;
                &nbsp;
            </span>
          </p>
        </div>
      </div>
      <div class="mt-2 sm:flex sm:justify-between">
        <div class="sm:flex">
          <p class="flex items-center skeleton w-80">
            &nbsp;
          </p>
          <p class="flex items-center skeleton w-80">
            &nbsp;
          </p>
        </div>
      </div>
      <div class="flex items-center justify-between">
        <div class="mt-4 flex items-center sm:mt-0 lg:ml-0.5 sm:ml-6">
          &nbsp;
          <div class="skeleton w-24 mt-3">
              &nbsp;
          </div>
        </div>
        <div class="flex mt-4">
          <div class="text-right skeleton w-14 mt-3 mr-6">
            &nbsp;
            &nbsp;
          </div>
        </div>
      </div>
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
