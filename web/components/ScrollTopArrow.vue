<template>
  <ScrollTop>
    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 bg-opacity-75 text-gray-800">
      <img src="/chevron-up.svg" alt="auto scroll up" class="h-8 w-15" />
    </span>
  </ScrollTop>
</template>

<script>
import ScrollTop from './ScrollTop'
export default {
  components: {
    ScrollTop,
  },
}
</script>
