<template>
  <div class="relative bg-gray-50 pt-16 pb-20 px-4 sm:px-6 lg:pt-24 lg:pb-28 lg:px-8">
    <div class="absolute inset-0">
      <div class="bg-white h-1/3 sm:h-2/3" />
    </div>
    <div class="relative max-w-7xl mx-auto">
      <div class="text-center">
        <h2 class="text-3xl tracking-tight font-extrabold text-gray-900 sm:text-4xl">
          <i18n
            path="blog.title"
            tag="h1"
            class="text-3xl xl:text-4xl text-light-onSurfacePrimary dark:text-dark-onSurfacePrimary font-medium leading-normal mb-6 lg:pt-4"
          >
            {{ $t('blog.title') }}
            <template #nuxt>
              <AppTitle />
            </template>
          </i18n>
        </h2>
        <p class="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
          <!--blog description i18n -->
          <i18n
            path="blog.description"
            tag="h3"
            class="xl:text-lg light:text-light-onSurfaceSecondary dark:text-dark-onSurfaceSecondary font-medium leading-relaxed mb-6"
          >
            <template #nuxtTeam>
              <NuxtLink class="text-nuxt-green underline" to="/company">
                {{ $t('blog.chaseapp_team') }}
              </NuxtLink>
            </template>
          </i18n>
        </p>
      </div>
      <div class="mt-12 max-w-lg mx-auto grid gap-5 lg:grid-cols-3 lg:max-w-none">
        <BlogpostPreviewItem
          v-for="(post, index) in posts"
          :key="index"
          :post="post"
        />
      </div>
    </div>
  </div>
</template>

<script>

export default {
  async asyncData ({ $content, app }) {
    let posts = await $content(app.i18n.defaultLocale, 'blog')
      .sortBy('date', 'desc')
      .fetch()

    if (app.i18n.defaultLocale !== app.i18n.locale)
      try {
        const newPosts = await $content(app.i18n.locale, 'blog')
          .sortBy('date', 'desc')
          .fetch()
        console.log(`NewP: ${newPosts}`)

        posts = posts.map((post) => {
          const newPost = newPosts.find(newPost => newPost.slug === post.slug)
          console.log(newPost)

          return newPost || post
        })
      } catch (err) {}

    return {
      posts,
    }
  },
  head () {
    const title = this.$i18n.t('blog.meta.title')
    const description = this.$i18n.t('blog.meta.description')

    return {
      title,
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: description,
        },
        // Open Graph
        {
          hid: 'og:title',
          property: 'og:title',
          content: title,
        },
        {
          hid: 'og:description',
          property: 'og:description',
          content: description,
        },
        // // Twitter Card
        {
          hid: 'twitter:title',
          name: 'twitter:title',
          content: title,
        },
        {
          hid: 'twitter:description',
          name: 'twitter:description',
          content: description,
        },
      ],
    }
  },
}
</script>

<style></style>
