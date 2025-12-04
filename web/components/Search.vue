<template>
  <client-only>
    <ais-instant-search
      v-on-clickaway="closeAutocomplete"
      :search-client="searchClient"
      index-name="chases"
    >
      <ais-autocomplete>
        <div slot-scope="{ indices, refine }" class="relative">
          <div>
            <input
              v-model="keywords"
              type="search"
              placeholder="Search"
              class="border border-transparent focus:bg-white focus:border-gray-300 placeholder-gray-600 rounded-lg bg-gray-200 py-2 pr-4 pl-10 block w-full appearance-none leading-normal"
              @input="refine($event.currentTarget.value)"
            >
            <div
              v-if="keywords && !hideAutocomplete"
              class="absolute left-0 mt-2 py-2 w-full z-50 bg-white rounded-lg shadow-xl"
            >
              <a
                v-for="hit in indices[0].hits"
                :key="hit.objectID"
                class="block px-4 py-2 text-gray-800 hover:bg-gray-200 hover:text-gray-800 cursor-pointer"
                @click="goToChase(hit.objectID)"
              >
                <ais-highlight attribute="Name" :hit="hit" />
              </a>
              <ais-powered-by class="px-4 py-2" />
            </div>
          </div>
          <div class="pointer-events-none absolute inset-y-0 left-0 pl-4 flex items-center">
            <svg class="fill-current pointer-events-none text-gray-600 w-4 h-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
              <path d="M12.9 14.32a8 8 0 1 1 1.41-1.41l5.35 5.33-1.42 1.42-5.33-5.34zM8 14A6 6 0 1 0 8 2a6 6 0 0 0 0 12z" />
            </svg>
          </div>
        </div>
      </ais-autocomplete>
    </ais-instant-search>
  </client-only>
</template>
<script>

import algoliasearch from 'algoliasearch/lite'

export default {
  data () {
    return {
      searchClient: {},
      keywords: '',
      hideAutocomplete: false,
      isAlgoliaLoaded: false,
    }
  },
  watch: {
    keywords (value) {
      if (value)
        this.hideAutocomplete = false
    },
  },
  mounted () {
    // eslint-disable-next-line nuxt/no-env-in-hooks
    if (process.client && algoliasearch)
      this.searchClient = algoliasearch(
        process.env.ALGOLIA_APPLICATION_ID,
        process.env.ALGOLIA_SEARCH_API_KEY)
  },
  methods: {
    goToChase (chaseId) {
      this.keywords = ''
      this.$router.push({
        name: 'chase-id',
        params: {
          id: chaseId,
        },
      })
    },
    closeAutocomplete () {
      this.hideAutocomplete = true
    },
  },
}
</script>
