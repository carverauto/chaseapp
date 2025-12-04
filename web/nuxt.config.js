// import webpack from 'webpack'
import path from 'path'
import fs from 'fs'
import { defineNuxtConfig } from '@nuxt/bridge'

// eslint-disable-next-line no-undef
require('dotenv').config()
const meta = {
  title: 'ChaseApp',
  description: 'Real-time notifications and chat for live police chases, rocket launches, weather events, disasters, and more',
  url: 'https://chaseapp.tv',
  image: 'https://chaseapp.tv/icon.png',
}

export default defineNuxtConfig ({
  alias: {
    tslib: 'tslib/tslib.es6.js',
  },
  webfontloader: {
    google: {
      families: ['Quicksand'],
    },
  },
  googleAnalytics: {
    id: 'UA-87374124-3',
  },
  ssr: true,
  debug: true,

  purgeCSS: {
    whitelistPatterns: [/mgl-map-wrapper.*/, /mapboxgl.*/],
  },

  /*
  ** Nuxt.js root directory
  ** See https://nuxtjs.org/api/configuration-srcdir/
  */
  /*
  ** Nuxt target
  ** See https://nuxtjs.org/api/configuration-target
  */
  target: 'server',
  server: process.env.NODE_ENV !== 'production'
    ? {
        https: {
          key: fs.readFileSync(path.resolve(__dirname, 'localhost-key.pem')),
          cert: fs.readFileSync(path.resolve(__dirname, 'localhost.pem')),
        },
      }
    : {},
  /*
  ** Headers of the page
  ** See https://nuxtjs.org/api/configuration-head
  */
  head: {
    title: 'ChaseApp',
    htmlAttrs: {
      lang: process.env.NUXT_LOCALE,
      dir: ['en'].includes(process.env.NUXT_LOCALE) ? 'rtl' : 'ltr',
    },
    meta: [
      { charset: 'utf-8' },
      { lang: 'en' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: meta.description },
      { name: 'msapplication-TileColor', content: '#91a1f6' },
      { name: 'theme-color', content: '#362f79' },

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

      // twitter
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:image', content: meta.image },

      { hid: 'twitter:title', name: 'twitter:title', content: meta.title },
      { hid: 'twitter:description', name: 'twitter:description', content: meta.description },

      // mobile
      { name: 'apple-mobile-web-app-capable', content: 'yes' },
      { name: 'mobile-web-app-capable', content: 'yes' },
      { name: 'apple-mobile-mobile-web-app-status-bar-style', content: '#362f79' },

    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/icon.png' },
      { rel: 'preconnect', href: 'https://firebasestorage.googleapis.com/' },
      // { rel: 'preconnect', href: 'https://www.google-analytics.com' },
      // { rel: 'sylesheet', href: 'https://fonts.googleapis.com/css?family=Quicksand&display=swap' },
    ],
    bodyAttrs: {
      class: [
        'font-sans font-medium bg-light-surface bg-dark-surface text-light-onSurfacePrimary dark:text-dark-onSurfacePrimary transition-colors duration-300 ease-linear',
      ],
    },
    script: [
      { src: '//cdnjs.cloudflare.com/ajax/libs/bodymovin/5.7.11/lottie.min.js', async: true },
    ],
  },
  /*
  ** Global CSS
  */
  // css: [ '@/assets/css/main.scss' ],
  css: ['@/assets/css/sprites.css', '@/assets/css/tooltip.css', 'plyr/dist/plyr.css'],

  sitemap: {
    path: '/sitemap.xml',
    hostname: 'https://chaseapp.tv',
    gzip: true,
  },
  /*
  ** Plugins to load before mounting the App
  ** https://nuxtjs.org/guide/plugins
  */
  plugins: [
    { src: '@/plugins/fireAuthInit.client.ts', ssr: false },
    { src: '@/plugins/firebaseAnalytics.client.ts', mode: 'client' },
    { src: '@/plugins/longpress.client.js', mode: 'client' },
    { src: '@/plugins/infiniteloading', ssr: false },
    { src: '@/plugins/persistedState.client.ts', ssr: false },
    { src: '@/plugins/clickaway', ssr: false },
    { src: '@/plugins/vue-instantsearch', mode: 'client' },
    // { src: '@/plugins/twitter', mode: 'client' },
    { src: '@/plugins/mapbox', mode: 'client' },
    { src: '@/plugins/helpers' },
    { src: '@/plugins/firebase' },
    { src: '@/plugins/tooltip' },
    // { src: '@/plugins/vue-plyr', mode: 'client' },
    // { src: '@/plugins/gtag' },
  ],
  /*
  ** Auto import components
  ** See https://nuxtjs.org/api/configuration-components
  */
  components: [
    '@/components',
    '@/components/cards',
    '@/components/brand',
    '@/components/chat',
    '@/components/birds',
    '@/components/firehose',
    '@/components/events',
    '@/components/blog',
    '@components/utils',
  ],
  /*
  ** Nuxt.js dev-modules
  */
  buildModules: [
    // just added back for testing 21-MAR-2022
    // '@nuxt/typescript-build',
    '@nuxtjs/eslint-module',
    'nuxt-windicss',
    '@nuxtjs/pwa',
    'nuxt-storm',
    '@nuxt-modules/compression',
  ],
  compression: {
    algorithm: 'brotliCompress',
  },
  vueI18n: {
    fallbackLocale: 'en',
  },
  fontawesome: {
    icons: {
      solid: true,
      brands: true,
    },
  },
  /*
  ** Nuxt.js modules
  */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/proxy',
    '@nuxtjs/axios',
    '@nuxtjs/dotenv',
    '@nuxtjs/dayjs',
    '@nuxtjs/markdownit',
    '@nuxtjs/google-fonts',
    ['nuxt-lazy-load', { directiveOnly: true }],
    ['nuxt-tailvue', { all: true }],
    '@nuxtjs/google-gtag',
    'vue-social-sharing/nuxt',
    'cookie-universal-nuxt',
  ],
  nuxtPrecompress: {
    enabled: true, // Enable in production
    report: false, // set true to turn one console messages during module init
    test: /\.(mjs|js|css|html|txt|xml|svg)$/, // files to compress on build
    // Serving options
    middleware: {
      // You can disable middleware if you serve static files using nginx...
      enabled: true,
      // Enable if you have .gz or .br files in /static/ folder
      enabledStatic: true,
      // Priority of content-encodings, first matched with request Accept-Encoding will me served
      encodingsPriority: ['br', 'gzip'],
    },
    // build time compression settings
    gzip: {
      // should compress to gzip?
      enabled: true,
      // compression config
      // https://www.npmjs.com/package/compression-webpack-plugin
      filename: '[path].gz[query]', // middleware will look for this filename
      threshold: 10240,
      minRatio: 0.8,
      compressionOptions: { level: 9 },
    },
    brotli: {
      // should compress to brotli?
      enabled: true,
      // compression config
      // https://www.npmjs.com/package/compression-webpack-plugin
      filename: '[path].br[query]', // middleware will look for this filename
      compressionOptions: { level: 11 },
      threshold: 10240,
      minRatio: 0.8,
    },
  },
  markdownit: {
    preset: 'default',
    linkify: true,
    breaks: true,
    use: ['markdown-it-link-attributes'],
  },
  'google-gtag':{
    id: 'G-BYC6KDR1PM',
    config:{
      anonymize_ip: true, // anonymize IP
      send_page_view: false, // might be necessary to avoid duplicated page track on page reload
    },
    debug: true, // enable to track in dev mode
    disableAutoPageTrack: false, // disable if you don't want to track each page route with router.afterEach(...)
    // optional you can add more configuration like [AdWords](https://developers.google.com/adwords-remarketing-tag/#configuring_the_global_site_tag_for_multiple_accounts)
  },
  googleFonts: {
    families: {
      Quicksand: [400, 500, 600, 700],
    },
    display: 'swap',
  },
  pwa: {
    manifest: {
      name: 'ChaseApp',
      description: 'Real-time Police Chase Notifications',
      theme_color: '#c50017',
    },
    workbox: {
      importScripts: ['/firebase-auth-sw.js', '/service-worker.js'],
      // by default the workbox module will not install the service worker in the dev env to avoid conflicts with HMR
      // only set this true for testing and remember to always clear your browser cache in development
      dev: process.env.NODE_ENV === 'development',
    },
  },
  dayjs: {
    locales: ['en'],
    defaultLocale: 'en',
    defaultTimeZone: 'US/Pacific',
    plugins: [
      'utc',
      'duration',
      'relativeTime',
      'timezone',
      'advancedFormat',
    ],
  },
  /*
  ** Axios module configuration
  ** See https://axios.nuxtjs.org/options
  */
  axios: {
    proxy: true,
    prefix: 'https://us-central1-chaseapp-8459b.cloudfunctions.net',
    headers: {
      // accept: 'application/json',
      'Access-Control-Allow-Origin': '*',
      'Content-Type': 'application/json',
      common: {
        'X-ApiKey': process.env.API_KEY,
      },
    },
  },
  proxy: {
    '/video/': {
      target: 'https://nbclim-download.edgesuite.net/',
    },
    //'/Video': '/Video/:path',
    '/ListAirports': '/ListAirports',
    '/GetStreamToken': '/GetStreamToken',
    '/DeleteUser': '/DeleteUser',
    // '/NBCLA': 'http://nbclim-download.edgesuite.net/:path',
  },
  /*
  ** Build configuration
  ** See https://nuxtjs.org/api/configuration-build/
  */
  build: {
    transpile: ['vue-instantsearch', 'instantsearch.js/es', 'dayjs', 'algoliasearch'],
  },

  /*
  ** Runtime Config
  ** See https://nuxtjs.org/guide/runtime-config/
  */
  publicRuntimeConfig: {
    apiUrl: process.env.API_URL,
    nuxtLocale: process.env.NUXT_LOCALE || 'en',
    credentials: false,
  },
})
