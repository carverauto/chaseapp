<template>
  <list-chases :id="$route.params.id" />
</template>

<script lang="ts">

import Vue from 'vue'
// eslint-disable-next-line import/named
import { MetaInfo } from 'vue-meta'
import { doc, getDoc } from 'firebase/firestore'
import { db } from '~/plugins/firebase'
import { Chase } from '@/types'
import PageMeta, { PageMetaInfo } from '@/lib/PageMeta'
// import { convertTimestamp } from 'convert-firebase-timestamp'

export default Vue.extend({
  // @ts-ignore
  async asyncData ({ params }) {
    const docRef = doc(db, 'chases', params.id)
    const docSnap = await getDoc(docRef)

    if (docSnap.exists()) {
      const chase = {
        ID: params.id,
        ...(docSnap.data()),
      }
      return { chase }
    }
  },
  data () {
    return {
      chase: undefined as undefined|Chase,
      myChase: undefined as any|Chase,
      image: 'https://chaseapp.tv/icon.png' as string,
      imgRegex: /\.([0-9a-z]+)(?:[?#]|$)/i,
      smImgReplace: '_200x200.webp?',
      lgImgReplace: '_1200x600.webp?',
    }
  },
  head () {
    // @ts-ignore
    const chase = this.chase as Chase
    // @ts-ignore
    const image = this.image as string
    const meta = new PageMeta()

    const newImage = chase?.ImageURL?.replace(this.imgRegex, this.lgImgReplace)
    // TODO: social image share shit goes here
    const info = {
      title: chase?.Name,
      description: chase?.Desc,
      url: `${process.env.WEB_URL}/chase/${this.$route.params.id}`,
      image: newImage || image,
    } as unknown as PageMetaInfo
    return meta.info(info) as MetaInfo
  },
})

</script>
