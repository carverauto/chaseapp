<template>
  <time
    :datetime="value.toISOString()"
    :data-tooltip="value.toDate().format('dddd, MMM Do, `YY @ hh:mm:ss a')"
  > {{ formatted(format) }}
  </time>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      required: false,
      default: 'fromNow',
    },
    format: {
      type: [ String, Boolean ],
      required: false,
      default: false,
    },
  },

  methods: {
    formatted () {
      if (this.format)
        return this.value.format(this.format)

      if (this.type === 'long')
        return this.value.format('dddd, MMM Do, `YY @ hh:mm:ss a')

      if (this.type === 'human')
        return this.value.format('ddd, MMM Do, h:mm a')

      if (this.type === 'time')
        return this.value.format('h:mm a')

      if (this.type === 'fromNow')
        return this.value.fromNow()
    },
  },
}
</script>
