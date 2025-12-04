// import firebase from 'firebase'
// import UserCredential = firebase.auth.UserCredential

// @ts-ignore
// const auth = this.$fire.auth

export default {
  data () {
    return {
      passphrase: 'ChaseChat',
    }
  },
  methods: {
    /*
    signInWithGoogleAuthentication () {
      const provider = new firebase.auth.GoogleAuthProvider()
      return new Promise((resolve, reject) => {
        auth
          .signInWithPopup(provider)
          .then(function (result: UserCredential) {
            resolve(result.user)
          })
          .catch(function (error: any) {
            reject(error)
          })
      })
    },
    saveUserToLocalStorage (user: firebase.auth.UserCredential) {
      if (process.client) {
        // @ts-ignore
        const encryptWithAES = CryptoJS.AES.encrypt(JSON.stringify(user), this.passphrase).toString()
        localStorage.setItem('enus', encryptWithAES)
      }
    },
    saveUserToStore (user: firebase.auth.UserCredential) {
      // @ts-ignore
      this.$store.commit('user.ts/add', user)
    },
    decryptUser (): any {
      if (process.client) {
        const encryptedUser = localStorage.getItem('enus')
        if (encryptedUser) {
          // @ts-ignore
          const bytes = CryptoJS.AES.decrypt(encryptedUser, this.passphrase)
          const decryptWithAES = bytes.toString(CryptoJS.enc.Utf8) // TODO: see if this is broken or not..
          return JSON.parse(decryptWithAES)
        }
      }
    },
     */
  },
}
