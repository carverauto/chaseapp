import initialState from './state'
import { AuthUser, State } from '~/types'

export default {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  setWrapper (state: State, payload: any) {
    const { object, prop, value } = payload
    object[prop] = value
  },

  RESET_STORE: (state: State) => {
    Object.assign(state, initialState())
  },

  SET_COORDS: (state: State, coords: number[]) => {
    if (coords)
      state.mapCoords = coords
    else
      console.log('Missing coords but we got SET_COORDS called')
  },

  SET_FIREHOSE_DATA: (state: State, firehoseData: object) => {
    if (firehoseData)
      state.firehoseData = <FirehoseData>firehoseData
  },

  SET_MESSAGING_TOKEN: (state: State, token: string) => {
    if (token)
      state.messagingToken = token
  },

  SET_STREAM_TOKEN: (state: State, token: string) => {
    if (token)
      state.streamToken = token
  },

  SET_AUTH_USER: (state: State, authUser: AuthUser, claims: any) => {
    if (authUser) {
      const { uid, email, emailVerified, displayName, photoURL } = authUser

      state.authUser = {
        uid,
        displayName,
        email,
        emailVerified,
        photoURL: photoURL || null,
        isAdmin: claims?.custom_claim,
      } as AuthUser
    }
  },
}
