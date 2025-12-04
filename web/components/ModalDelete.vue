<template>
  <div>
    <div class="rounded-md bg-yellow-50 p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <!-- Heroicon name: solid/exclamation -->
          <svg class="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-yellow-800">
            Warning: you are about to delete your account
          </h3>
          <div class="mt-2 text-sm text-yellow-700">
            <p>
              This is an irreversible procedure, are you sure you want to continue?
            </p>
            <div class="flex inline-flex m-4">
              <button class="m-1 bg-red-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" @click="deleteAccount">
                Confirm Delete
              </button>
              <nuxt-link to="/">
                <button class="m-2 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                  Cancel
                </button>
              </nuxt-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {mapGetters, mapState} from "vuex";
import {AuthUser} from "~/types";

export default {
  middleware: 'auth',
  data () {
    return {
      loading: {
        delete: false,
      },
    }
  },
  mounted () {
    console.log(this.authUser)
  },
  computed: {
    ...mapState({
      // @ts-ignore
      authUser: state => state.authUser as AuthUser,
    }),
    ...mapGetters({
      isLoggedIn: 'isLoggedIn',
    }),
    getAvatar () {
      return this.authUser?.photoURL || '/user.png'
    },
  },
  methods: {
    async deleteAccount () {
      this.deleteState = 'loading'
      const user = this.$fire.auth.currentUser
      // TODO: delete tokens for user
      await this.$axios.$post('/DeleteUser', { id: this.authUser.uid }).then((res) => {
        const user = this.$fire.auth.currentUser
        user.delete().then(function () {
          this.$toast.show({
            type: 'success',
            title: 'User deleted',
            message: 'User deleted',
            timeout: 0,
          })
        })
        this.$router.push('/')
      }).catch((error) => {
        console.error(error)
      })
      user.delete().catch((error) => {
        console.error(error)
      })
      await this.$router.push('/signin')
    },
  },
}
</script>

<style scoped>

</style>
