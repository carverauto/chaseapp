<template>
  <fieldset class="space-y-8 divide-y divide-gray-200">
    <div class="space-y-8 divide-y divide-gray-200">
      <div>
        <div>
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Profile
          </h3>
          <p class="mt-1 text-sm text-gray-500">
            This information will be displayed publicly so be careful what you share.
          </p>
        </div>

        <div class="mt-6 grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
          <div class="sm:col-span-4">
            <label for="username" class="block text-sm font-medium text-gray-700">
              Username
            </label>
            <div class="mt-1 flex rounded-md shadow-sm">
              <input
                id="username"
                v-model="userName"
                type="text"
                name="username"
                autocomplete="username"
                class="flex-1 focus:ring-blue-500 focus:border-blue-500 block w-full min-w-0 rounded-none rounded-r-md sm:text-sm border-gray-300"
              >
            </div>
          </div>

          <!--
          <div class="sm:col-span-6">
            <label for="photo" class="block text-sm font-medium text-gray-700">
              Photo
            </label>
            <div class="mt-1 flex items-center">
              <span class="h-12 w-12 rounded-full overflow-hidden bg-gray-100">
                <div v-if="photoURL">
                  <img :src="photoURL">
                </div>
                <div v-else>
                  <svg class="h-full w-full text-gray-300" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M24 20.993V24H0v-2.996A14.977 14.977 0 0112.004 15c4.904 0 9.26 2.354 11.996 5.993zM16.002 8.999a4 4 0 11-8 0 4 4 0 018 0z" />
                  </svg>
                </div>
              </span>
              <button type="button" class="ml-5 bg-white py-2 px-3 border border-gray-300 rounded-md shadow-sm text-sm leading-4 font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                Change
              </button>
            </div>
          </div>
          -->
        </div>
      </div>
    </div>

    <div class="pt-5">
      <div class="flex justify-end">
        <nuxt-link to="/">
          <button type="button" class="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
            Go Back
          </button>
        </nuxt-link>
        <button type="submit" class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500" @click="save">
          Save
        </button>
      </div>
    </div>
  </fieldset>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { doc, onSnapshot, setDoc } from 'firebase/firestore'
import { db } from '~/plugins/firebase'
import { AuthUser } from '@/types'

export default Vue.extend({
  middleware: 'auth',
  data () {
    return {
      userData: {} as any|AuthUser,
      userName: '',
      photoURL: '',
      email: '',
    }
  },
  computed: {
    ...mapGetters({
      isLoggedIn: 'isLoggedIn',
    }),
    authUser () {
      return this.$store.state.authUser as AuthUser
    },
  },
  mounted () {
    // @ts-ignore
    this.getUserInfo()
  },
  methods: {
    cancel () {
      this.$router.push('/profile')
    },
    async save () {
      // @ts-ignore
      if (this.authUser.uid)
        try {
          // @ts-ignore
          await setDoc(doc(db, 'users', this.authUser.uid), {
            // @ts-ignore
            email: this.authUser.email,
            // @ts-ignore
            uid: this.authUser.uid,
            lastUpdated: Date.now(),
            // @ts-ignore
            photoURL: this.authUser.photoURL,
            // @ts-ignore
            userName: this.userName,
          }, { merge: true })
        } catch (e) {
          console.log(`Can't save profile: ${e}`)
          this.$toast.show({
            type: 'danger',
            title: 'ChaseApp - Profile',
            message: `Unable to save profile <b>${e}</b>`,
            timeout: 0,
          })
        }
      this.$toast.show({
        type: 'success',
        title: 'ChaseApp - Profile',
        message: 'Profile Updated',
        timeout: 3,
      })
    },
    getUserInfo (): void {
      // @ts-ignore
      onSnapshot(doc(db, 'users', this.authUser.uid), (doc) => {
        // @ts-ignore
        this.userData = doc.data() as AuthUser
        // @ts-ignore
        const displayName = this.authUser.displayName
        // @ts-ignore
        this.userName = this.userData.userName || displayName
        // @ts-ignore
        this.email = this.userData.email
        // @ts-ignore
        this.photoURL = this.userData.photoURL || this.authUser.photoURL
      })
    },
  },
})
</script>
