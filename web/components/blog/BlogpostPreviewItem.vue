<template>
  <NuxtLink
    :to="localePath({ name: 'blog-slug', params: { slug: post.slug } })"
  >
    <div class="flex-1 flex-col rounded-lg shadow-lg overflow-hidden">
      <div class="flex-shrink-0">
        <img class="h-48 w-full object-cover" alt="post image url" :src="post.imgUrl">
      </div>
      <div class="flex-1 bg-white p-6 flex flex-col justify-between">
        <div class="flex-1">
          <p class="text-sm font-medium text-pink-800">
            <a href="#" class="hover:underline">
              Article
            </a>
          </p>
          <a href="#" class="block mt-2">
            <p class="text-xl font-semibold text-gray-900">
              {{ post.title }}
            </p>
            <p class="mt-3 text-base text-gray-500">
              {{ post.description }}
            </p>
          </a>
        </div>
        <div class="mt-6 flex items-center">
          <div class="flex-shrink-0">
            <a href="#">
              <span class="sr-only">
                <BlogpostAuthor
                  v-for="(author, index) in post.authors"
                  :key="index"
                  :author="author"
                />
              </span>
            </a>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-900">
              <a href="#" class="hover:underline">
                <BlogpostAuthor
                  v-for="(author, index) in post.authors"
                  :key="index"
                  :author="author"
                />
              </a>
            </p>
            <div class="flex space-x-1 text-sm text-gray-500">
              <time datetime="formatDateByLocale(post.date)">
                <format-date :value="$dayjs(post.date)" />
              </time>
              <span aria-hidden="true">
                &middot;
              </span>
              <span>
                {{ post.readingTime.text }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </NuxtLink>
</template>

<script>
export default {
  props: {
    post: {
      type: Object,
      required: true,
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
