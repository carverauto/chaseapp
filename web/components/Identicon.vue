<template>
  <!--
  <circle cx={s / 2} cy={s / 2} r={s / 2} fill="#eee"/>
  <circle cx={c} cy={c - r} r={z} fill={colors[i++]}/>
  <circle cx={c} cy={c - ro2} r={z} fill={colors[i++]}/>
  <circle cx={c - rroot3o4} cy={c - r3o4} r={z} fill={colors[i++]}/>
  <circle cx={c - rroot3o2} cy={c - ro2} r={z} fill={colors[i++]}/>
  <circle cx={c - rroot3o4} cy={c - ro4} r={z} fill={colors[i++]}/>
  <circle cx={c - rroot3o2} cy={c} r={z} fill={colors[i++]}/>
  <circle cx={c - rroot3o2} cy={c + ro2} r={z} fill={colors[i++]}/>
  <circle cx={c - rroot3o4} cy={c + ro4} r={z} fill={colors[i++]}/>
  <circle cx={c - rroot3o4} cy={c + r3o4} r={z} fill={colors[i++]}/>
  <circle cx={c} cy={c + r} r={z} fill={colors[i++]}/>
  <circle cx={c} cy={c + ro2} r={z} fill={colors[i++]}/>
  <circle cx={c + rroot3o4} cy={c + r3o4} r={z} fill={colors[i++]}/>
  <circle cx={c + rroot3o2} cy={c + ro2} r={z} fill={colors[i++]}/>
  <circle cx={c + rroot3o4} cy={c + ro4} r={z} fill={colors[i++]}/>
  <circle cx={c + rroot3o2} cy={c} r={z} fill={colors[i++]}/>
  <circle cx={c + rroot3o2} cy={c - ro2} r={z} fill={colors[i++]}/>
  <circle cx={c + rroot3o4} cy={c - ro4} r={z} fill={colors[i++]}/>
  <circle cx={c + rroot3o4} cy={c - r3o4} r={z} fill={colors[i++]}/>
  <circle cx={c} cy={c} r={z} fill={colors[i++]}/>
</svg>)
-->
  <svg viewBox="0 0 64 64">
    <circle cx="32" cy="32" fill="#eee" r="32" />
    <circle cx="32" cy="8" :fill="colors[0]" r="5" />
    <circle cx="32" cy="20" :fill="colors[1]" r="5" />
    <circle cx="21.607695154586736" cy="14" :fill="colors[2]" r="5" />
    <circle cx="11.215390309173472" cy="20" :fill="colors[3]" r="5" />
    <circle cx="21.607695154586736" cy="26" :fill="colors[4]" r="5" />
    <circle cx="11.215390309173472" cy="32" :fill="colors[5]" r="5" />
    <circle cx="11.215390309173472" cy="44" :fill="colors[6]" r="5" />
    <circle cx="21.607695154586736" cy="38" :fill="colors[7]" r="5" />
    <circle cx="21.607695154586736" cy="50" :fill="colors[8]" r="5" />
    <circle cx="32" cy="56" :fill="colors[9]" r="5" />
    <circle cx="32" cy="44" :fill="colors[10]" r="5" />
    <circle cx="42.392304845413264" cy="50" :fill="colors[11]" r="5" />
    <circle cx="52.78460969082653" cy="44" :fill="colors[12]" r="5" />
    <circle cx="42.392304845413264" cy="38" :fill="colors[13]" r="5" />
    <circle cx="52.78460969082653" cy="32" :fill="colors[14]" r="5" />
    <circle cx="52.78460969082653" cy="20" :fill="colors[15]" r="5" />
    <circle cx="42.392304845413264" cy="26" :fill="colors[16]" r="5" />
    <circle cx="42.392304845413264" cy="14" :fill="colors[17]" r="5" />
    <circle cx="32" cy="32" :fill="colors[18]" r="5" />
  </svg>
</template>

<script>
const { blake2b } = require('blakejs')
export default {
  props: {
    value: {
      type: String,
      required: true,
    },
  },
  computed: {
    colors () {
      const schema = {
        target: { freq: 1, colors: [ 0, 28, 0, 0, 28, 0, 0, 28, 0, 0, 28, 0, 0, 28, 0, 0, 28, 0, 1 ] },
        cube: { freq: 20, colors: [ 0, 1, 3, 2, 4, 3, 0, 1, 3, 2, 4, 3, 0, 1, 3, 2, 4, 3, 5 ] },
        quazar: { freq: 16, colors: [ 1, 2, 3, 1, 2, 4, 5, 5, 4, 1, 2, 3, 1, 2, 4, 5, 5, 4, 0 ] },
        flower: { freq: 32, colors: [ 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 3 ] },
        cyclic: { freq: 32, colors: [ 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6 ] },
        vmirror: { freq: 128, colors: [ 0, 1, 2, 3, 4, 5, 3, 4, 2, 0, 1, 6, 7, 8, 9, 7, 8, 6, 10 ] },
        hmirror: { freq: 128, colors: [ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 8, 6, 7, 5, 3, 4, 2, 11 ] },
      }

      const total = Object.keys(schema).map(k => schema[k].freq).reduce((a, b) => a + b)
      const findScheme = (d) => {
        let cum = 0
        const ks = Object.keys(schema)
        for (const i in ks) {
          const n = schema[ks[i]].freq
          cum += n
          if (d < cum)
            return schema[ks[i]]
        }
        throw new Error('Impossible')
      }
      let nid = this.value
      // console.log(nid)

      const zero = blake2b(new Uint8Array([ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ]))
      nid = Array.from(blake2b(nid)).map((x, i) => (x + 256 - zero[i]) % 256)

      const sat = (Math.floor(nid[29] * 70 / 256 + 26) % 80) + 30
      const d = Math.floor((nid[30] + nid[31] * 256) % total)
      const scheme = findScheme(d)
      const palette = Array.from(nid).map((x, i) => {
        const b = (x + i % 28 * 58) % 256
        if (b === 0)
          return '#444'

        if (b === 255)
          return 'transparent'

        const h = Math.floor(b % 64 * 360 / 64)
        const l = [ 53, 15, 35, 75 ][Math.floor(b / 64)]
        return `hsl(${h}, ${sat}%, ${l}%)`
      })

      const rot = (nid[28] % 6) * 3

      return scheme.colors.map((_, i) => palette[scheme.colors[i < 18 ? (i + rot) % 18 : 18]])
    },
  },
}
</script>