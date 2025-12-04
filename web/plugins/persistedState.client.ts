import createPersistedState from 'vuex-persistedstate'
// import { Store } from 'vuex'

// @ts-ignore
export default ({ store }) => {
  createPersistedState()(store)
}
