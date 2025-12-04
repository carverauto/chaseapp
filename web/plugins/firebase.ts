import { initializeApp, getApps } from 'firebase/app'
import { getFirestore } from 'firebase/firestore'
import { getDatabase } from 'firebase/database'
import { FirebaseApp, FirebaseOptions } from '@firebase/app-types'
import { getAuth } from 'firebase/auth'

const firebaseConfig = {
  apiKey: "AIzaSyDZVvCuh81AYFsNqNhdI5GUzwQC91na580",
  authDomain: "chaseapp-8459b.firebaseapp.com",
  databaseURL: "https://chaseapp-8459b.firebaseio.com",
  projectId: "chaseapp-8459b",
  storageBucket: "chaseapp-8459b.appspot.com",
  messagingSenderId: "1020122644146",
  appId: "1:1020122644146:web:68f163a80a77facbcc13ab",
  measurementId: "G-BYC6KDR1PM"
} as FirebaseOptions

const apps = getApps()
// eslint-disable-next-line import/no-mutable-exports
let firebaseApp: FirebaseApp

if (!apps.length)
  firebaseApp = initializeApp(firebaseConfig) as FirebaseApp
else
  firebaseApp = apps[0] as FirebaseApp

const auth = getAuth(firebaseApp)
const db = getFirestore(firebaseApp)
const rtdb = getDatabase(firebaseApp)

export { db, rtdb, auth, firebaseApp }
