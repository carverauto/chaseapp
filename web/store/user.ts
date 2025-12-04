// eslint-disable-next-line import/named
import { GetterTree } from 'vuex'
import { State, AuthUser } from '@/types'

export const state = () => ({
  user: undefined,
})

export const mutations = {
  add (state: State, { uid, displayName, photoURL, email }:AuthUser) {
    state.user = {
      uid,
      displayName,
      photoURL,
      email,
    }
  },
  remove (state: State) {
    // @ts-ignore
    state.user = undefined
  },
}

export const getters: GetterTree<State, State> = {
  getUser: state => state.user,
}
