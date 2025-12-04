export default {
  async updateCoords ({ commit }: any, { coords }: any) {
    await commit('SET_COORDS', coords)
  },

  async updateMessagingToken ({ commit }: any, { token }: any) {
    await commit('SET_MESSAGING_TOKEN', token)
  },

  async updateFirehoseData ({ commit }: any, { firehoseData }: any) {
    await commit('SET_FIREHOSE_DATA', firehoseData)
  },

  async updateStreamToken ({ commit }: any, { token }: any) {
    await commit('SET_STREAM_TOKEN', token)
  },

  async nuxtServerInit ({ dispatch }: any, { res }: any) {
    // INFO -> Nuxt-fire Objects can be accessed in nuxtServerInit action via this.$fire___, ctx.$fire___ and ctx.app.$fire___'

    if (res && res.locals && res.locals.user) {
      const { allClaims: claims, idToken: token, ...authUser } = res.locals.user

      await dispatch('onAuthStateChanged', {
        authUser,
        claims,
        token,
      })
    }
  },

  async onAuthStateChanged ({ commit }: any, { authUser, claims }: any) {
    if (!authUser) {
      // commit('RESET_STORE')
      commit('RESET_STORE')
      return
    }
    if (authUser && authUser.getIdToken)
      try {
        await authUser.getIdToken(true)
      } catch (e) {
        console.error(e)
        return e
      }

    const { uid, email, emailVerified, displayName, photoURL } = authUser
    console.log(authUser)
    commit('SET_AUTH_USER', {
      uid,
      email,
      emailVerified,
      displayName,
      photoURL,
      isAdmin: claims.custom_claim,
    })
  },
}
