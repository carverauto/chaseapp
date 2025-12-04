<template>
  <SignupStage stage="3" />
</template>

<script lang="ts">

import Vue from 'vue'
import {
  getAuth,
  isSignInWithEmailLink,
  signInWithEmailLink,
} from 'firebase/auth'
const auth = getAuth()

// import { analytics } from "~/plugins/firebaseAnalytics.client";
// import { logEvent } from "firebase/analytics"

export default Vue.extend({
  data () {
    return {
      email: undefined as undefined|any,
    }
  },
  mounted () {
    this.completeSignIn()
  },
  methods: {
    completeSignIn () {
    // Confirm the link is a sign-in with email link.
      if (isSignInWithEmailLink(auth, window.location.href)) {
      // Additional state parameters can also be passed via URL.
      // This can be used to continue the user's intended action before triggering
      // the sign-in operation.
      // Get the email if available. This should be available if the user completes
      // the flow on the same device where they started it.
        this.email = window.localStorage.getItem('emailForSignIn') as string
        if (!this.email)
        // User opened the link on a different device. To prevent session fixation
        // attacks, ask the user to provide the associated email again. For example:
          this.email = window.prompt('Please provide your email for confirmation')

        // The client SDK will parse the code from the link for you.
        signInWithEmailLink(auth, this.email, window.location.href)
          .then(() => {
          // Clear email from storage.
            window.localStorage.removeItem('emailForSignIn')
            // You can access the new user via result.user
            // Additional user info profile not available via:
            // result.additionalUserInfo.profile == null
            // You can check if the user is new or existing:
            // result.additionalUserInfo.isNewUser
            // logEvent(analytics, 'sign_up')
            this.$router.push('/')
          })
          .catch((e) => {
          // Some error occurred, you can inspect the code: error.code
          // Common errors could be invalid email and invalid or expired OTPs.
          // @ts-ignore
            this.$toast.show({
              type: 'error',
              title: 'ChaseApp - Email Link Authentication',
              message: `Problem with signInWithEmailLink <b>${e}</b>`,
              timeout: 0,
            })
          })
      }
    },
  },
// validate ({ params }) { // Must be a number return params.id !== undefined }
})
</script>
