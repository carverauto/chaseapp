/*
 * <format-date-diff
 *  class="font-medium"
 *  :start="$dayjs()"
 *  :end="$dayjs($auth.user.created_at).add(8, 'days')"
 *  unit="days"
 * />
*/

<template>
  <span>{{ diff }}</span>
</template>

<script>
export default {
  props: {
    start: {
      type: Object,
      required: true,
    },
    end: {
      type: Object,
      required: true,
    },
    unit: {
      type: String,
      required: false,
      default: 'seconds',
    },
  },
  computed: {
    diff () {
      return this.end.diff(this.start, this.unit)
    },
  },
}
</script>
