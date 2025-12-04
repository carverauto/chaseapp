<template>
  <div class="flow-root">
    <div v-if="myevent.eventType === 'system'">
      <System :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'bof'">
      <Bof :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'twitter'">
      <TwitterEvent :eventProp="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'launch'">
      <Launch :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'weather'">
      <Weather :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'youtube'">
      <LiveStream :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'volcano'">
      <Volcano :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'airshow'">
      <Airshow :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'fire'">
      <Fire :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'halloween'">
      <Halloween :event="myevent" />
    </div>
    <div v-else-if="myevent.eventType === 'thanksgiving'">
      <Thanksgiving :event="myevent" />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  props: {
    myevent: {
      type: Object,
      required: true,
    },
  },
  data () {
    return {
      bof: false,
    }
  },
})
</script>
