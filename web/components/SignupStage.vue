<template>
  <nav aria-label="Progress">
    <ol class="space-y-4 md:flex md:space-y-0 md:space-x-8">
      <li class="md:flex-1">
        <!-- Completed Step -->
        <div v-if='stage >= 1' class="group pl-4 py-2 flex flex-col border-l-4 border-blue-600 hover:border-blue-800 md:pl-0 md:pt-4 md:pb-0 md:border-l-0 md:border-t-4">
          <span class="text-xs text-blue-600 font-semibold tracking-wide uppercase group-hover:text-blue-800">Step 1</span>
          <span class="text-sm font-medium">E-mail link requested</span>
        </div>
      </li>

      <li class="md:flex-1">
        <!-- Current Step -->
        <div v-if='stage > 1' class="pl-4 py-2 flex flex-col border-l-4 border-blue-600 md:pl-0 md:pt-4 md:pb-0 md:border-l-0 md:border-t-4" aria-current="step">
          <span class="text-xs text-blue-600 font-semibold tracking-wide uppercase">Step 2</span>
          <span class="text-sm font-medium">E-mail sent, open and click to continue</span>
        </div>
      </li>

      <li class="md:flex-1">
        <!-- Upcoming Step -->
        <div v-if='stage >= 2' class="group pl-4 py-2 flex flex-col border-l-4 border-gray-200 hover:border-gray-300 md:pl-0 md:pt-4 md:pb-0 md:border-l-0 md:border-t-4">
          <span class="text-xs text-gray-500 font-semibold tracking-wide uppercase group-hover:text-gray-700">Step 3</span>
          <span class="text-sm font-medium">Logging in..</span>
        </div>
      </li>
    </ol>
  </nav>
</template>

<script>

// import { analytics } from "~/plugins/firebaseAnalytics.client"
// import { logEvent } from "firebase/analytics"

export default {
  name: "SignupStage",
  props: {
    stage: {
      type: Number,
      required: false,
      default: 1,
    },
  },
  mounted () {
    // if (this.stage >= 2)
      // logEvent(analytics, 'sign_up')
  },
}
</script>

<style scoped>

</style>
