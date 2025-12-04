import { State } from '@/types'

export default {
  getCoords: (state: State) => {
    try {
      return state.mapCoords !== null
    } catch {
      return false
    }
  },
  getStreamToken: (state: State) => {
    try {
      return state.streamToken !== null
    } catch {
      return false
    }
  },
  getFirehoseData: (state: State) => {
    try {
      return state.firehoseData !== null
    } catch {
      return false
    }
  },
  getMessagingToken: (state: State) => {
    try {
      return state.messagingToken !== null
    } catch {
      return false
    }
  },
  isLoggedIn: (state: State) => {
    try {
      return state.authUser.uid !== null
    } catch {
      return false
    }
  },
}
