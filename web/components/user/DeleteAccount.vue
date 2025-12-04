<template>
  <div class="mt-6 bg-white shadow sm:rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <h3 class="text-lg leading-6 font-medium text-gray-900">
        Delete your account
      </h3>
      <div class="mt-2 max-w-xl text-sm text-gray-500">
        <p>
          Once you delete your account, you will lose all data associated with it.
        </p>
      </div>
      <div class="mt-5">
        <button class="bg-red-600 py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500" @click="showDelete = true">
          Delete Account
        </button>
        <ModalBase v-show="showDelete">
          <ModalDelete />
        </ModalBase>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DeleteAccount',
  data () {
    return {
      deleteState: 'active',
      showDelete: false,
    }
  },
  methods: {
    async deleteAccount () {
      this.deleteState = 'loading'

    },
  },
}
</script>

<style scoped>

</style>
