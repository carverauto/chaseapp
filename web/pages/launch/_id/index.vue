<template>
  <list-launches :id="$route.params.id" />
</template>

<script lang="ts">

import Vue from 'vue'
// eslint-disable-next-line import/named
import { MetaInfo } from 'vue-meta'
// import { NuxtFireInstance } from '@nuxtjs/firebase'
// import { convertTimestamp } from 'convert-firebase-timestamp'
import { Launch } from '@/types'
import PageMeta, { PageMetaInfo } from '@/lib/PageMeta'

export default Vue.extend({
  // @ts-ignore
  // async asyncData ({ params, $fire } /* : { params: Object & { id: string }, $fire: NuxtFireInstance } */) {
  asyncData ({ params } /* : { params: Object & { id: string }, $fire: NuxtFireInstance } */) {
    const launch = {
      id: params.id,
      // ...(await ($fire as NuxtFireInstance).firestore.collection('launches').doc(params.id).get()).data(),
    }
    return { launch }
  },
  data () {
    return {
      launch: undefined as undefined|Launch,
      image: 'https://chaseapp.tv/icon.png' as string,
    }
  },
  head () {
    // @ts-ignore
    const launch = this.launch as Launch
    // @ts-ignore
    const image = this.image as string
    const meta = new PageMeta()

    const info = {
      title: launch?.Name,
      description: launch?.Desc,
      url: `${process.env.WEB_URL}/launch/${this.$route.params.id}`,
      image: launch?.ImageURL ? launch?.ImageURL : image,
    } as unknown as PageMetaInfo
    return meta.info(info) as MetaInfo
  },
})

</script>
