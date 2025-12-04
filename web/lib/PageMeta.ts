
// eslint-disable-next-line import/named
import { MetaInfo } from 'vue-meta'

export interface PageMetaInfo {
    title: string
    description: string
    url: string
    image: string
    updated_time: string
}

export default class {
  info (meta: PageMetaInfo): MetaInfo {
    const structuredData = {
      '@type': 'Article' as String,
      // datePublished: this.chase.CreatedAt.seconds, // TODO: fix datePublished in PageMetaInfo
      headline: meta.title as String,
      image: meta.image as String,
      updated_time: meta.updated_time as String,
    } as Object

    return {
      title: meta.title,
      meta: [
        { hid: 'description', name: 'description', content: meta.description },

        // Schema.org
        { hid: 'itemprop:name', itemprop: 'name', content: meta.title },
        { hid: 'itemprop:description', itemprop: 'description', content: meta.description },
        { hid: 'itemprop:image', itemprop: 'image', content: meta.image },

        // facebook
        { hid: 'og:type', property: 'og:type', content: 'website' },
        { hid: 'og:site_name', property: 'og:site_name', content: 'ChaseApp' },
        { hid: 'og:url', property: 'og:url', content: meta.url },
        { hid: 'og:image', property: 'og:image', content: meta.image },
        { hid: 'og:title', property: 'og:title', content: meta.title },
        { hid: 'og:description', property: 'og:description', content: meta.description },
        { hid: 'og:image:width', property: 'og:image:width', content: '600' },
        { hid: 'og:image:height', property: 'og:image:height', content: '600' },

        // twitter
        { hid: 'twitter:card', name: 'twitter:card', content: 'summary_large_image' },
        { hid: 'twitter:image', name: 'twitter:image', content: meta.image },
        { hid: 'twitter:title', name: 'twitter:title', content: meta.title },
        { hid: 'twitter:description', name: 'twitter:description', content: meta.description },
      ],
      link: [
        { rel: 'canonical', href: meta.url },
      ],
      script: [
        { type: 'application/ld+json', json: structuredData as any },
      ],
    }
  }
}
