import { onAuthStateChanged, getAuth } from 'firebase/auth'

export default ({ store }: any) => {
  const auth = getAuth()

  onAuthStateChanged(auth, (user) => {
    if (user) {
      const photoUrl = user.providerData[0].photoURL

      store.commit('SET_AUTH_USER', {
        email: user.email,
        uid: user.uid,
        photoURL: photoUrl,
        emailVerified: user.emailVerified,
        displayName: user.displayName,
      })
    }
    else
      store.commit('RESET_STORE')
  })
}
