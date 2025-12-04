<template>
  <fieldset class="space-y-8 divide-y divide-gray-200">
    <div class="space-y-8 divide-y divide-gray-200">
      <div>
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Profile
          </h3>
          <p class="mt-1 text-sm text-gray-500">
            Please update your e-mail address to continue.
          </p>
        </div>

        <div class="mt-6 grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
          <div class="sm:col-span-4">
            <label for="email" class="block text-sm font-medium text-gray-700">
              E-mail
            </label>
            <div class="mt-1 flex rounded-md shadow-sm">
              <input
                id="email"
                v-model="email"
                required
                type="text"
                name="email"
                autocomplete="email"
                class="flex-1 focus:ring-blue-500 focus:border-blue-500 block w-full min-w-0 rounded-none rounded-r-md sm:text-sm border-gray-300"
              >
            </div>
          </div>
        </div>
      </div>
    </div>

    <div>
      <button class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500" @click="save">
        <IconSpinner
          v-if="loading"
          class="w-5 h-5"
          primary="text-red-500"
          secondary="text-blue-300"
        />
        Update Profile
      </button>
    </div>
  </fieldset>
</template>

<script lang="ts">
import Vue from 'vue'
import {
  getAuth,
  sendSignInLinkToEmail,
  isSignInWithEmailLink,
  signInWithEmailLink,
  onAuthStateChanged,
} from 'firebase/auth'
import { mapGetters } from 'vuex'
import { doc, setDoc } from 'firebase/firestore'
import { AuthUser } from '~/types'
import { db } from '~/plugins/firebase'
const auth = getAuth()

import { analytics } from "~/plugins/firebaseAnalytics.client"
import { logEvent } from "firebase/analytics"

export default Vue.extend({
  data () {
    return {
      uid: this.$route.params.id,
      loading: false,
      email: '' as '' | string,
    }
  },
  computed: {
    ...mapGetters({
      isLoggedIn: 'isLoggedIn',
    }),
    authUser (): AuthUser {
      return this.$store.state.authUser as AuthUser
    },
  },
  mounted () {
    this.completeSignIn()
  },
  methods: {
    completeSignIn () {
      // Confirm the link is a sign-in with email link.
      if (isSignInWithEmailLink(auth, window.location.href)) {
        let email = window.localStorage.getItem('emailForSignIn')
        if (!email)
        // User opened the link on a different device. To prevent session fixation
        // attacks, ask the user to provide the associated email again. For example:
          email = window.prompt('Please provide your email for confirmation') as string
        // The client SDK will parse the code from the link for you.
        signInWithEmailLink(auth, email, window.location.href).then(() => {
          // Clear email from storage.
          window.localStorage.removeItem('emailForSignIn')
          // You can access the new user via result.user
          // Additional user info profile not available via:
          // result.additionalUserInfo.profile == null
          // You can check if the user is new or existing:
          // result.additionalUserInfo.isNewUser
          logEvent(analytics, 'sign_up')
          this.$router.push('/')
        })
          .catch((e) => {
            // Some error occurred, you can inspect the code: error.code
            // Common errors could be invalid email and invalid or expired OTPs.
            this.$toast.show({
              type: 'error',
              title: 'ChaseApp - Email Link Authentication',
              message: `Problem with signInWithEmailLink <b>${e}</b>`,
              timeout: 0,
            })
          })
      }
    },
    cancel () {
      this.$router.push(`/signup/${this.uid}`)
    },
    async save () {
      const myEmail = this.email
      if (this.authUser.uid)
        try {
          await setDoc(doc(db, 'users', this.authUser.uid), {
            email: myEmail,
            uid: this.authUser.uid,
            lastUpdated: Date.now(),
          }, { merge: true })

          // @ts-ignore
          this.loading = true
          const actionCodeSettings = {
            url: 'https://chaseapp.tv/signup/',
            handleCodeInApp: true,
            // When multiple custom dynamic link domains are defined, specify which
            // one to use.
            // dynamicLinkDomain: "m.chaseapp.tv"
          }
          sendSignInLinkToEmail(auth, myEmail, actionCodeSettings)
            .then(() => {
              window.localStorage.setItem('emailForSignIn', myEmail)
              onAuthStateChanged(auth, (user) => {
                if (user) {
                  this.save()
                  this.$router.push('/networks')
                }
              })
            })
            .catch((e) => {
              // Error occurred. Inspect error.code.
              const errorCode = e.code
              const errorMessage = e.message
              console.log(`errorCode: ${errorCode} and errorMessage: ${errorMessage}`)
              // this.$toast.show({ type: 'danger', title: 'ChaseApp - Email Link Authentication', message: `E-mail link not sent - <b>${e}</b>`, timeout: 0, })
            })

          return this.$toast.show({
            type: 'success',
            title: 'ChaseApp - Profile',
            message: 'Profile Updated',
            timeout: 3,
          })
        } catch (e) {
          console.log(`Can't save profile: ${e}`)
          this.$toast.show({
            type: 'danger',
            title: 'ChaseApp - Profile',
            message: `Unable to save profile <b>${e}</b>`,
            timeout: 3,
          })
        }
    },
  },
})
</script>
