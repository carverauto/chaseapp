// vuex.d.ts
// eslint-disable-next-line @typescript-eslint/no-unused-vars
// import { ComponentCustomProperties } from 'vue'
import { Store } from 'vuex'
import { State } from '~/types/index'

declare module '@vue/runtime-core' {
  // declare your own store states
  /*
  interface State {
    count: number
    authUser: AuthUser
  }
   */

  // provide typings for `this.$store`
  interface ComponentCustomProperties {
    $store: Store<State>
  }
}
