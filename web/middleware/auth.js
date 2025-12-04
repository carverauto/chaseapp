export default function ({ store, redirect }) {
  if (!store.getters.isLoggedIn) {
    console.log('Not logged in, redirecting')
    return redirect('/signin')
  }
}
