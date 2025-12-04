<template>
  <div>
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <div v-if="loading.email" class="mt-8">
        <SignupStage :stage="two" />
      </div>
      <div v-else>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Sign in to your account
        </h2>
        <p class="mt-2 text-center text-sm text-black-800 max-w">
          or
          <span class="font-medium hover:text-black-500">
            Sign-in to start a free account
          </span>
        </p>
        <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
          <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
            <div class="space-y-6">
              <div>
                <label for="email" class="block text-sm font-medium text-gray-700">
                  Email address
                </label>
                <div class="mt-1">
                  <input
                    id="email"
                    v-model="email"
                    name="email"
                    type="email"
                    autocomplete="email"
                    required
                    class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  >
                </div>
              </div>

              <div>
                <button class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500" @click="sendEmailLink">
                  <IconSpinner
                    v-if="loading.email"
                    class="w-5 h-5"
                    primary="text-red-500"
                    secondary="text-blue-300"
                  />
                  Sign in
                </button>
              </div>
            </div>

            <div class="mt-6">
              <div class="relative">
                <div class="absolute inset-0 flex items-center">
                  <div class="w-full border-t border-gray-300" />
                </div>
                <div class="relative flex justify-center text-sm">
                  <span class="px-2 bg-white text-gray-500">
                    Or continue with
                  </span>
                </div>
              </div>

              <div class="mt-6 grid grid-cols-3 gap-3">
                <div>
                  <PushButton class="w-full" @click="login('google')">
                    <IconSpinner
                      v-if="loading.google"
                      class="w-5 h-5"
                      primary="text-red-500"
                      secondary="text-blue-300"
                    />
                    <icon-google v-else class="w-5 h-5" alt="Sign-In with Google" />
                  </PushButton>
                </div>
                <div>
                  <PushButton class="w-full" @click="login('facebook')">
                    <IconSpinner
                      v-if="loading.facebook"
                      class="w-5 h-5"
                      primary="text-black"
                      secondary="text-gray-600"
                    />
                    <svg v-show="!loading.facebook" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                      <path fill-rule="evenodd" d="M20 10c0-5.523-4.477-10-10-10S0 4.477 0 10c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V10h2.54V7.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V10h2.773l-.443 2.89h-2.33v6.988C16.343 19.128 20 14.991 20 10z" clip-rule="evenodd" />
                    </svg>
                  </PushButton>
                </div>
                <div>
                  <PushButton class="w-full" @click="login('twitter')">
                    <IconSpinner
                      v-if="loading.twitter"
                      class="w-5 h-5"
                      primary="text-orange-500"
                      secondary="text-orange-300"
                    />
                    <svg v-show="!loading.twitter" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                      <path d="M6.29 18.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0020 3.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.073 4.073 0 01.8 7.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 010 16.407a11.616 11.616 0 006.29 1.84" />
                    </svg>
                  </PushButton>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

import Vue from 'vue'
import {
  getAuth,
  sendSignInLinkToEmail,
  signInWithPopup,
  FacebookAuthProvider,
  GoogleAuthProvider,
  TwitterAuthProvider,
  onAuthStateChanged,
} from 'firebase/auth'
import { doc, setDoc } from 'firebase/firestore'
import { db } from '~/plugins/firebase'

// import { analytics } from "~/plugins/firebaseAnalytics.client"
// import { logEvent } from "firebase/analytics"

const auth = getAuth()

export default Vue.extend({
  data () {
    return {
      email: '',
      two: 2,
      password: undefined as undefined|any,
      loading: {
        google: false,
        facebook: false,
        twitter: false,
        email: false,
      },
    }
  },
  methods: {
    sendEmailLink () {
      if (this.email) {
        this.loading.email = true
        const actionCodeSettings = {
          url: 'https://chaseapp.tv/signup/',
          handleCodeInApp: true,
          // When multiple custom dynamic link domains are defined, specify which
          // one to use.
          // dynamicLinkDomain: "m.chaseapp.tv"
        }
        sendSignInLinkToEmail(auth, this.email, actionCodeSettings)
          .then(() => {
            // Verification email sent.
            window.localStorage.setItem('emailForSignIn', this.email)
            onAuthStateChanged(auth, (user) => {
              if (user) {
                console.log(user)
                this.save(this.email, user.uid)
                // logEvent(analytics, 'login')
                this.$router.push('/')
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
      } else
        // @ts-ignore
        this.$toast.show({ type: 'danger', title: 'ChaseApp - Missing e-mail', message: 'Must supply an email address', timeout: 3 })
    },
    async save (email: string, uid: string) {
      if (uid)
        await setDoc(doc(db, 'users', uid), {
          email,
          uid,
          lastUpdated: Date.now(),
        }, { merge: true })
      else
        console.log('Must supply UID')
    },
    login (provider: string) {
      // @ts-ignore
      this.loading[provider] = true

      // firebaseProvider, not facebook
      let fbProvider
      switch (provider) {
        case 'google':
          fbProvider = new GoogleAuthProvider()
          break
        case 'facebook':
          fbProvider = new FacebookAuthProvider()
          break
        case 'twitter':
          fbProvider = new TwitterAuthProvider()
      }

      if (fbProvider)
        signInWithPopup(auth, fbProvider)
          .then((result) => {
            const user = result.user
            if (user) {
              if (!user.email)
              // lets get their email before we continue
                this.$router.push(`/signup/${user.uid}`)
              else {
                // console.log(user)
                this.save(user.email, user.uid)
                // logEvent(analytics, 'login')
                this.$router.push('/')
              }
              // @ts-ignore
              this.loading[provider] = false
            }
          }).catch((error) => {
            const errorCode = error.code
            const errorMessage = error.message
            console.error(error)
            // @ts-ignore
            this.$toast.show({
              type: 'danger',
              title: 'Login failed',
              message: `${provider} login failed - <b>${errorCode}/${errorMessage}</b>`,
              timeout: 0,
            })
            // @ts-ignore
            this.loading[provider] = false
            this.$router.push('/signin')
          })
      else
        console.log('Must supply auth provider')
    },
  },
})
</script>
