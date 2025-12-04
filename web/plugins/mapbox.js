import Vue from 'vue'
import {
  MglMap,
  MglMarker,
  MglNavigationControl,
  MglGeolocateControl,
  MglPopup,
  MglImageLayer,
  MglGeojsonLayer, MglFullscreenControl,
} from 'vue-mapbox'
import Mapbox from 'mapbox-gl'

Vue.component('MglMap', MglMap)
Vue.component('MglMarker', MglMarker)
Vue.component('MglPopup', MglPopup)
Vue.component('MglImageLayer', MglImageLayer)
Vue.component('MglGeojsonLayer', MglGeojsonLayer)
Vue.component('MglGeolocateControl', MglGeolocateControl)
Vue.component('MglNavigationControl', MglNavigationControl)
Vue.component('MglFullscreenControl', MglFullscreenControl)

Vue.prototype.$mapbox = Mapbox
