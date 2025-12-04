import Vue from "vue"
// Vue.use(VTooltip)

import { VTooltip, VPopover, VClosePopover } from "v-tooltip"

Vue.directive("tooltip", VTooltip);
Vue.directive("close-popover", VClosePopover)
Vue.component("v-popover", VPopover)