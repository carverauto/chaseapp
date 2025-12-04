<template>
  <div v-if="$route.params.id">
    <div v-if="chatSupported">
      <span class="text-left lg:text-lg sm:text-3xl mb-4 mr-2 text-gray-900 sm:text-2xl">
        Live chat
      </span>
      <span v-if="watcher_count">
        <span class="inline-flex items-center px-3 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
          <svg
            v-tooltip="tooltips.msg.watcherCount"
            class="w-4 h-4 mr-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
          /><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
          /></svg>
          {{ watcher_count }}
        </span>
      </span>
      <client-only>
        <div
          vue-chat-scroll
          class="messages container mx-auto overflow-y-auto mb-1 mt-2 lg:h-96 h-72"
        >
          <div v-if="isLoggedIn">
            <div
              v-if="chats.length > 0"
              class="flex-col items-center justify-center font-sans lg:space-y-4 divide-y divide-gray-300 lg:divide-y-0"
            >
              <span v-if="userData">
                <show-message
                  v-for="(myMessage, $index) in chats"
                  :key="$index"
                  class="message"
                  :message="myMessage"
                  :user-data="userData"
                />
              </span>
            </div>
            <div v-else>
              <p class="mt-4">Loading..</p>
            </div>
          </div>
          <div v-else>
            <p class="mt-4 text-xs text-gray-600">You must login to chat</p>
          </div>
        </div>

        <div
          class="w-full flex justify-between"
          style="bottom: 0px;"
        >
          <img
            v-if="photoURL"
            class="mt-4 inline-block h-8 w-8 rounded-full"
            :src="photoURL"
            alt="avatar"
          >
          <textarea
            id="message"
            v-model="message"
            class="flex-grow m-2 py-2 px-4 mr-1 rounded-full border border-gray-300 bg-gray-200 resize-none"
            rows="1"
            placeholder="Say something..."
            style="outline: none;"
            @keyup.enter="onSubmit"
          />

          <button
            class="m-2"
            style="outline: none;"
            @click="onSubmit"
          >
            <svg
              class="svg-inline--fa text-green-400 fa-paper-plane fa-w-16 w-12 h-12 py-2 mr-2"
              aria-hidden="true"
              focusable="false"
              data-prefix="fas"
              data-icon="paper-plane"
              role="img"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 512 512"
            >
              <path
                fill="currentColor"
                d="M476 3.2L12.5 270.6c-18.1 10.4-15.8 35.6 2.2 43.2L121 358.4l287.3-253.2c5.5-4.9 13.3 2.6 8.6 8.3L176 407v80.5c0 23.6 28.5 32.9 42.5 15.8L282 426l124.6 52.2c14.2 6 30.4-2.9 33-18.2l72-432C515 7.8 493.3-6.8 476 3.2z"
              />
            </svg>
          </button>
        </div>
      </client-only>
    </div>
    <div v-else>
      <div v-if="isLoggedIn">
        <p>Chat temporarily disabled</p>
      </div>
      <div v-else>
        <p>Sign-In to chat</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import VueChatScroll from 'vue-chat-scroll'
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { doc, getDoc } from 'firebase/firestore'
import { db } from '~/plugins/firebase'
import { StreamChat, TokenOrProvider} from 'stream-chat'

const chatClient = StreamChat.getInstance(process.env.GETSTREAM_API_KEY)

// eslint-disable-next-line @typescript-eslint/no-unsafe-argument
Vue.use(VueChatScroll)

export default {
  props: {
    chaseId: {
      type: String,
      required: true,
    },
    chaseName: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      chatSupported: false,
      uid: '',
      nickname: '',
      message: '',
      displayName: '',
      messagesLoaded: false,
      userData: '',
      userName: '',
      email: '',
      photoURL: 'https://chaseapp.tv/icon.png',
      data: {type: '', nickname: '', message: ''},
      chats: [],
      errors: [],
      offStatus: false,
      chState: {},
      channel: {},
      watcher_count: '',
      tooltips: {
        msg: {
          watcherCount: "Connected Users",
        },
      },
    }
  },
  computed: {
    ...mapGetters({
      isLoggedIn: 'isLoggedIn',
    }),
    authUser() {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-return,@typescript-eslint/no-unsafe-member-access
      return this.$store.state.authUser
    },
    screenHeight() {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-return,@typescript-eslint/no-unsafe-member-access,@typescript-eslint/restrict-plus-operands
      return this.isDevice ? window.innerHeight + 'px' : 'calc(100vh - 80px)'
    },
  },
  async created() {
    await this.startStreamChat()
  },
  destroyed() {
    this.resetMessages()
  },
  methods: {
    async addReaction(n) {
      console.log('Got event from child')
      // console.log(n)
      if (this.channel)
        await this.channel.sendReaction(n.messageId, {
          // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
          type: n.reaction,
        })
      // console.log(result)

    },
    async startStreamChat() {
      if (this.isLoggedIn) {
        await this.getUserInfo()
        // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
        const uid = this.authUser.uid
        this.nickname = this.userName
        this.message = ''
        if (uid)
            // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
          await this.$axios.post('https://us-central1-chaseapp-8459b.cloudfunctions.net/GetStreamToken',
              {
                // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
                user_id: uid,
              },
          ).then(async (res: { data: { data: { message: TokenOrProvider } } }) => {
            try {
              await chatClient.connectUser({
                    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
                    id: this.authUser.uid,
                    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
                    name: this.userName ? this.userName : this.authUser.name,
                    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
                    image: this.authUser.photoURL ? this.authUser.photoURL : this.photoURL,
                  },
                  res.data.data.message,
              )
              // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
              this.chatSupported = true
            } catch (e) {
              console.log(e)
              // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
              this.chatSupported = false
              return
            }

            // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment,@typescript-eslint/no-unsafe-member-access
            this.channel = chatClient.channel('livestream', this.chaseId, {
              // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
              name: this.chaseName ? this.chaseName : "NA",
            })

            if (this.channel) {
              // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
              this.chState = await this.channel.watch({ presence: true })
              // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-assignment
              this.watcher_count = this.chState.watcher_count
              // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access,@typescript-eslint/no-unsafe-assignment
              this.chats = this.chState.messages

              this.channel.on('message.new', (event: { message: any }) => {
                // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
                this.chats.push(event.message)
              })

              this.channel.on('message.updated', (myEvent: any) => {
                console.log(`Message Updated: ${myEvent}`)
              })

              this.channel.on('reaction.new'), (myEvent: any) => {
                console.log('We got a new reaction')
                for (const i of this.chats) 
                    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
                  if (i.id == myEvent.message.id)
                    console.log('We got a new reaction to ${myEvent.message.id}')
                  
                
              }
              this.channel.on('reaction.updated', myEvent => {
                console.log('We got a reaction to update')
                console.log(myEvent)
                for (const i of this.chats)
                  if (i.id === myEvent.message.id)
                    console.log('We got a reaction to our reaction')
                  
                
                /*
                this.chats.forEach((value => {
                  if (value.id === event.message.id) {
                    console.log(value)
                    value.reaction_counts = event.message.reaction_counts
                  }
                }))
                 */
              })

              this.channel.on('user.watching.start', () => {
                console.log('User watching')
                this.watcher_count++
              })

              this.channel.on('user.watching.stop', () => {
                console.log('User leaving')
                this.watcher_count--
              })
            }
          })

      }
    },
    resetMessages () {
      this.chats = []
      this.messagesLoaded = false
    },
    async getUserInfo () {
      const docRef = doc(db, 'users', this.authUser.uid)
      const docSnap = await getDoc(docRef)
      if (docSnap.exists()) {
        this.userData = docSnap.data()
        this.uid = this.userData.uid
        this.displayName = this.authUser.displayName
        this.userData.displayName = this.displayName || null
        this.userName = this.userData.userName || this.displayName
        this.email = this.userData.email
        this.photoURL = this.userData.photoURL || this.authUser.photoURL
      } else
        console.log('No docSnap')
    },
    async onSubmit (evt) {
      evt.preventDefault()
      // if (process.client)
      //  window.addEventListener('touchstart', ontouchstart(evt), { passive: true })

      if (this.isLoggedIn) {
        // Don't allow empty messages to be sent
        if (/\S/.test(this.message)) {
          const response = await this.channel.sendMessage({
            text: this.message,
          })
        }

        this.message = ''
      } else
        this.$toast.show('You must login to chat')
    },
  },
}
</script>
