// @ts-ignore
window.twttr = (function (d, s, id) {
  const fjs = d.getElementsByTagName(s)[0]
  // @ts-ignore
  const t = window.twttr || {}
  if (d.getElementById(id)) return t
  const js = d.createElement(s)
  js.id = id
  // @ts-ignore
  js.src = 'https://platform.twitter.com/widgets.js'
  fjs.parentNode?.insertBefore(js, fjs)

  t._e = []
  t.ready = function (f) {
    t._e.push(f)
  }

  return t
}(document, 'script', 'twitter-wjs'))
