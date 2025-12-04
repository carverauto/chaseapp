<template>
  <article>
    <header class="flex items-left justify-between flex-col mt-12">
      <div class="flex flex-1 flex-col mb-8">
        <h1 class="text-4xl font-semibold mb-4 leading-tight">
          {{ post.title }}
        </h1>
        <div class="text-sm flex justify-between flex-col sm:flex-row">
          <div>
            <BlogpostAuthor
              v-for="(author, index) in post.authors"
              :key="index"
              :author="author"
            />
          </div>
          <div class="mt-1">
            {{ formatDateByLocale(post.date) }}
            <span class="text-xs mx-1">&bullet;</span>
            {{ post.readingTime.text }}
          </div>
        </div>
      </div>
      <AppImage :src="post.imgUrl" ratio="16:9" sizes="80vh" class="rounded" />
    </header>
    <div class="mt-12">
      <nav>
        <ul>
          <h1 class="text-3xl font-semibold mb-4 leading-tight">
            Table of Contents
          </h1>
          <li v-for="link of post.toc" :key="link.id">
            <NuxtLink :class="{ 'py-2': link.depth === 2, 'ml-2 pb-2': link.depth === 3 }" :to="`#${link.id}`">{{ link.text }}</NuxtLink>
          </li>
        </ul>
      </nav>
      <nuxt-content :document="post" />
    </div>
    <div
      v-if="hasTags"
      class="border-t border-light-border dark:border-dark-border my-10"
    >
      <div
        class="flex flex-row flex-wrap justify-start my-10"
      >
        <span
          v-for="(tag, id) in post.tags"
          :key="id"
          class="inline-flex items-center px-2.5 py-0.5 rounded-md text-sm font-medium bg-blue-100 text-blue-800"
        >
          <svg class="-ml-0.5 mr-1.5 h-2 w-2 text-pink-400" fill="currentColor" viewBox="0 0 8 8">
            <circle cx="4" cy="4" r="3" />
          </svg>
          {{ tag }}
        </span>
      </div>
    </div>
  </article>
</template>

<script>
export default {
  name: 'BlogpostItem',
  props: {
    post: {
      type: Object,
      required: true,
    },
  },
  computed: {
    hasTags () {
      return this.post.tags
    },
  },
  methods: {
    formatDateByLocale (d) {
      const currentLocale = this.$i18n.locale || 'en'
      const options = { year: 'numeric', month: 'long', day: 'numeric' }
      return new Date(d).toLocaleDateString(currentLocale, options)
    },
  },
}
</script>
