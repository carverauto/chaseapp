import { GetterTree } from 'vuex'
import { State } from '~/types'

export const state = () => ({
  mapCoords: undefined,
})

export const mutations = {
  add (state: State, coords: number[]) {
    state.mapCoords = coords
  },
  remove (state: State) {
    state.mapCoords = [-98.2856656, 36.2612751]
  },
}

export const getters = <GetterTree<State, any>>{
  getCoords: state => state.mapCoords,
}