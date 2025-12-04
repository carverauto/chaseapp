// middleware/cacheControl.ts
// @ts-ignore
/*
import { defineNuxtServerMiddleware } from '@nuxtjs/composition-api'
const staticPage = 'max-age=15778476, s-maxage=15778476'
// I may have some other caching scemes here in the future
const pathMap: { [key: string]: string } = {
    '/': staticPage,
    '/about': staticPage,
    '/privacy': staticPage,
    '/tos': staticPage
}
export default defineNuxtServerMiddleware((req: any , res: any, next: any) => {
    const {
        originalUrl,
        headers: { host }
    } = req
    if (!originalUrl || !host || host.split(':')[0] === 'localhost') return next()
    const cache = pathMap[originalUrl]
    if (cache) res.setHeader('Cache-Control', `public, ${cache}`)
    return next()
})
 */