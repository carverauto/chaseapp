<template>
  <main class="campaign-builder">
    <section class="campaign-builder__header">
      <p class="campaign-builder__eyebrow">
        Tools
      </p>
      <div class="campaign-builder__title-row">
        <div>
          <h1>Campaign Builder</h1>
          <p>
            Build a ready-to-review GetChase growth campaign from an event,
            region, budget, and audience goal.
          </p>
        </div>
        <button
          type="button"
          class="campaign-builder__primary"
          @click="copyBrief"
        >
          Copy brief
        </button>
      </div>
    </section>

    <section class="campaign-builder__grid">
      <form
        class="campaign-builder__panel"
        @submit.prevent
      >
        <div class="campaign-builder__section-heading">
          <h2>Campaign inputs</h2>
          <span>Draft only</span>
        </div>

        <label>
          Ad account
          <select
            v-model="form.adAccount"
            :disabled="!connectedAdAccounts.length"
          >
            <option
              v-for="adAccount in metaAdAccounts"
              :key="adAccount.id"
              :value="adAccount.label"
            >
              {{ adAccount.label }}
            </option>
          </select>
        </label>

        <div class="campaign-builder__actions">
          <button
            class="campaign-builder__secondary"
            :disabled="loadingAdAccounts"
            type="button"
            @click="loadMetaAdAccounts"
          >
            {{ loadingAdAccounts ? 'Loading accounts' : 'Connect agency Business Manager' }}
          </button>
          <p>{{ metaConnectionStatus }}</p>
        </div>

        <label>
          Campaign goal
          <select v-model="form.goal">
            <option
              v-for="goal in goals"
              :key="goal"
              :value="goal"
            >
              {{ goal }}
            </option>
          </select>
        </label>

        <label>
          Tracking event
          <select v-model="form.trackingEvent">
            <option
              v-for="trackingEvent in trackingEvents"
              :key="trackingEvent"
              :value="trackingEvent"
            >
              {{ trackingEvent }}
            </option>
          </select>
        </label>

        <label>
          Target region
          <input
            v-model="form.region"
            type="text"
            placeholder="Los Angeles, CA"
          >
        </label>

        <div class="campaign-builder__split">
          <label>
            Budget
            <input
              v-model.number="form.budget"
              min="50"
              step="50"
              type="number"
            >
          </label>

          <label>
            Flight length
            <select v-model="form.flightLength">
              <option
                v-for="flightLength in flightLengths"
                :key="flightLength"
                :value="flightLength"
              >
                {{ flightLength }}
              </option>
            </select>
          </label>
        </div>

        <label>
          Attribution window
          <select v-model="form.attributionWindow">
            <option
              v-for="attributionWindow in attributionWindows"
              :key="attributionWindow"
              :value="attributionWindow"
            >
              {{ attributionWindow }}
            </option>
          </select>
        </label>

        <div class="campaign-builder__section-heading campaign-builder__section-heading--sub">
          <h2>Audience setup</h2>
          <span>AI-assisted</span>
        </div>

        <div class="campaign-builder__split">
          <label>
            Minimum age
            <input
              v-model.number="form.ageMin"
              min="13"
              max="65"
              type="number"
            >
          </label>

          <label>
            Maximum age
            <input
              v-model.number="form.ageMax"
              min="13"
              max="65"
              type="number"
            >
          </label>
        </div>

        <label>
          Gender
          <select v-model="form.gender">
            <option
              v-for="gender in genders"
              :key="gender"
              :value="gender"
            >
              {{ gender }}
            </option>
          </select>
        </label>

        <label>
          Custom audience
          <select v-model="form.customAudience">
            <option
              v-for="customAudience in customAudiences"
              :key="customAudience"
              :value="customAudience"
            >
              {{ customAudience }}
            </option>
          </select>
        </label>

        <label>
          Audience instructions for AI
          <textarea
            v-model="form.audiencePrompt"
            rows="5"
            placeholder="Describe who should see this ad. The AI should translate this into interests, behaviors, and demographic suggestions."
          />
        </label>

        <label>
          Tone
          <select v-model="form.tone">
            <option
              v-for="tone in tones"
              :key="tone"
              :value="tone"
            >
              {{ tone }}
            </option>
          </select>
        </label>

        <label>
          Landing destination
          <input
            v-model="form.destination"
            type="text"
            placeholder="https://chaseapp.tv"
          >
        </label>

        <label>
          Extra context
          <textarea
            v-model="form.context"
            rows="4"
            placeholder="Add promo details, exclusions, or regions to avoid."
          />
        </label>
      </form>

      <section class="campaign-builder__panel campaign-builder__preview">
        <div class="campaign-builder__section-heading">
          <h2>Generated campaign</h2>
          <span>Human approval required</span>
        </div>

        <article>
          <h3>Strategy</h3>
          <p>{{ strategy }}</p>
        </article>

        <article>
          <h3>Meta campaign structure</h3>
          <dl>
            <div>
              <dt>Ad account</dt>
              <dd>{{ form.adAccount }}</dd>
            </div>
            <div>
              <dt>Objective</dt>
              <dd>{{ metaObjective }}</dd>
            </div>
            <div>
              <dt>Tracking event</dt>
              <dd>{{ form.trackingEvent }}</dd>
            </div>
            <div>
              <dt>Attribution</dt>
              <dd>{{ form.attributionWindow }}</dd>
            </div>
            <div>
              <dt>Budget</dt>
              <dd>${{ form.budget }} over {{ form.flightLength }}</dd>
            </div>
            <div>
              <dt>Audience</dt>
              <dd>{{ audienceSummary }}</dd>
            </div>
            <div>
              <dt>Destination</dt>
              <dd>{{ form.destination }}</dd>
            </div>
          </dl>
        </article>

        <article>
          <h3>AI audience recommendation</h3>
          <p>{{ aiAudienceRecommendation }}</p>
        </article>

        <article>
          <h3>Ad copy</h3>
          <ul>
            <li
              v-for="copy in adCopy"
              :key="copy"
            >
              {{ copy }}
            </li>
          </ul>
        </article>

        <article>
          <h3>Creative concepts</h3>
          <ul>
            <li
              v-for="concept in creativeConcepts"
              :key="concept"
            >
              {{ concept }}
            </li>
          </ul>
        </article>

        <article>
          <h3>Companion push copy</h3>
          <p>{{ pushCopy }}</p>
        </article>

        <article>
          <h3>Approval checklist</h3>
          <ul>
            <li>Review event claims and regional targeting.</li>
            <li>Confirm creative follows Meta policy before upload.</li>
            <li>Send approved assets to Meta Bulk Uploader.</li>
            <li>Attach UTM tracking before launch.</li>
          </ul>
        </article>
      </section>
    </section>
  </main>
</template>

<script>
export default {
  name: 'CampaignBuilder',
  data () {
    return {
      connectedAdAccounts: [],
      loadingAdAccounts: false,
      metaConnectionStatus: 'Agency Business Manager not connected',
      goals: [
        'App installs',
        'Subscriptions',
        'Event awareness',
        'User reactivation',
        'Local market growth',
      ],
      trackingEvents: [
        'App install',
        'Subscribe',
        'Start trial',
        'View content',
        'Open live chase',
        'Enable notifications',
        'Purchase',
      ],
      attributionWindows: [
        '1-day click',
        '7-day click',
        '1-day click or view',
        '7-day click or 1-day view',
        '28-day click or 1-day view',
      ],
      genders: [
        'All',
        'Men',
        'Women',
      ],
      customAudiences: [
        'None',
        'All app users',
        'Free users',
        'Lapsed subscribers',
        'Trial started but not subscribed',
        'Website visitors',
        'High-intent lookalike seed',
      ],
      flightLengths: [
        '3 days',
        '7 days',
        '14 days',
        '30 days',
      ],
      tones: [
        'Urgent',
        'Local',
        'Community',
        'News-style',
        'Direct response',
      ],
      form: {
        adAccount: '',
        goal: 'App installs',
        trackingEvent: 'App install',
        region: 'Los Angeles, CA',
        ageMin: 18,
        ageMax: 54,
        gender: 'All',
        customAudience: 'None',
        audiencePrompt: 'Find people likely to follow live breaking events, police activity, local news, scanner pages, traffic alerts, and real-time public safety updates.',
        budget: 500,
        flightLength: '7 days',
        attributionWindow: '7-day click or 1-day view',
        tone: 'Urgent',
        destination: 'https://chaseapp.tv',
        context: '',
      },
    }
  },
  head () {
    return {
      title: 'Campaign Builder',
    }
  },
  mounted () {
    this.loadMetaAdAccounts()
  },
  computed: {
    metaAdAccounts () {
      return this.connectedAdAccounts.length
        ? this.connectedAdAccounts
        : [{ id: 'not-connected', label: 'Connect agency Business Manager to load ad accounts' }]
    },
    metaObjective () {
      const objectives = {
        'App installs': 'App promotion',
        Subscriptions: 'Sales or conversions',
        'Event awareness': 'Awareness',
        'User reactivation': 'Traffic or engagement',
        'Local market growth': 'Awareness with regional reach',
      }

      return objectives[this.form.goal] || 'Traffic'
    },
    audienceSummary () {
      return `${this.form.gender}, ages ${this.form.ageMin}-${this.form.ageMax}, ${this.form.customAudience.toLowerCase()}, ${this.form.region}`
    },
    aiAudienceRecommendation () {
      return `Use the audience instructions to suggest interests, behaviors, and demographic filters for ${this.form.region}. Prioritize users interested in local news, breaking news, police activity, traffic alerts, live streams, emergency alerts, and real-time community updates. Seed or exclude custom audiences using "${this.form.customAudience}" as the starting point.`
    },
    strategy () {
      const adAccount = this.form.adAccount || 'the selected Meta ad account'
      return `${this.form.goal} campaign in ${this.form.region} optimized for ${this.form.trackingEvent.toLowerCase()} on ${adAccount}. Use a ${this.form.tone.toLowerCase()} message, target ${this.audienceSummary.toLowerCase()}, route traffic to ${this.form.destination}, and keep publishing gated behind final approval.`
    },
    adCopy () {
      const region = this.form.region
      const trackingEvent = this.form.trackingEvent.toLowerCase()

      return [
        `Live GetChase alerts for ${region}. Know what is happening as it unfolds.`,
        `GetChase tracks fast-moving events in ${region} so you do not miss the moments people talk about.`,
        `Follow real-time updates, maps, and community context from GetChase.`,
        `Designed to drive ${trackingEvent}: fast alerts, live maps, and the context people need.`,
        `When breaking events move fast in ${region}, GetChase helps you keep up.`,
      ]
    },
    creativeConcepts () {
      return [
        `Map-first static showing an active ${this.form.region} alert with a clear app install call to action.`,
        `Split creative with live event feed on one side and phone notification preview on the other.`,
        `Short vertical video concept: alert arrives, map opens, viewer sees the event timeline.`,
        `Community proof creative showing active viewers following the same live update.`,
      ]
    },
    pushCopy () {
      return `New ${this.form.trackingEvent.toLowerCase()} campaign ready for ${this.form.region}. Review copy, creative, budget, attribution, and targeting before export.`
    },
    brief () {
      return [
        'Campaign Builder brief',
        `Ad account: ${this.form.adAccount}`,
        `Goal: ${this.form.goal}`,
        `Tracking event: ${this.form.trackingEvent}`,
        `Attribution window: ${this.form.attributionWindow}`,
        `Region: ${this.form.region}`,
        `Audience: ${this.audienceSummary}`,
        `Audience instructions for AI: ${this.form.audiencePrompt}`,
        `AI recommendation: ${this.aiAudienceRecommendation}`,
        `Budget: $${this.form.budget} over ${this.form.flightLength}`,
        `Tone: ${this.form.tone}`,
        `Destination: ${this.form.destination}`,
        `Strategy: ${this.strategy}`,
        'Ad copy:',
        ...this.adCopy.map(copy => `- ${copy}`),
        'Creative concepts:',
        ...this.creativeConcepts.map(concept => `- ${concept}`),
        `Extra context: ${this.form.context || 'None'}`,
      ].join('\n')
    },
  },
  methods: {
    normalizeAdAccounts (accounts) {
      const seen = new Map()

      accounts.forEach((account) => {
        if (!account || !account.id) {
          return
        }

        seen.set(account.id, {
          id: account.id,
          label: `${account.name || 'Unnamed account'} - ${account.id}${account.currency ? ` - ${account.currency}` : ''}`,
        })
      })

      return Array.from(seen.values()).sort((a, b) => a.label.localeCompare(b.label))
    },
    async loadMetaAdAccounts () {
      this.loadingAdAccounts = true
      this.metaConnectionStatus = 'Loading agency ad accounts'

      try {
        const payload = await this.$axios.$get('/api/v1/meta/ad-accounts')
        this.connectedAdAccounts = this.normalizeAdAccounts(payload.data || [])

        if (this.connectedAdAccounts.length) {
          this.form.adAccount = this.connectedAdAccounts[0].label
          this.metaConnectionStatus = `${this.connectedAdAccounts.length} ad account${this.connectedAdAccounts.length === 1 ? '' : 's'} loaded`
        } else {
          this.metaConnectionStatus = 'No ad accounts found'
        }
      } catch (error) {
        this.metaConnectionStatus = error.message || 'Meta connection failed'
      } finally {
        this.loadingAdAccounts = false
      }
    },
    async copyBrief () {
      if (!navigator.clipboard) {
        return
      }

      await navigator.clipboard.writeText(this.brief)
    },
  },
}
</script>

<style scoped>
.campaign-builder {
  width: min(1180px, calc(100vw - 24px));
  margin: 24px auto 56px;
  color: #111827;
}

.campaign-builder__header {
  margin-bottom: 20px;
}

.campaign-builder__eyebrow {
  margin: 0 0 6px;
  color: #4f46e5;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0;
  text-transform: uppercase;
}

.campaign-builder__title-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 18px;
}

.campaign-builder h1 {
  margin: 0;
  color: #0f172a;
  font-size: 34px;
  line-height: 1.1;
  font-weight: 800;
  letter-spacing: 0;
}

.campaign-builder__title-row p {
  max-width: 680px;
  margin: 10px 0 0;
  color: #475569;
  font-size: 16px;
  line-height: 1.6;
}

.campaign-builder__primary {
  min-height: 44px;
  padding: 0 18px;
  border: 0;
  border-radius: 6px;
  background: #2563eb;
  color: white;
  font-weight: 700;
  cursor: pointer;
}

.campaign-builder__secondary {
  min-height: 40px;
  padding: 0 14px;
  border: 1px solid #2563eb;
  border-radius: 6px;
  background: #eff6ff;
  color: #1d4ed8;
  font-weight: 700;
  cursor: pointer;
}

.campaign-builder__secondary:disabled {
  cursor: wait;
  opacity: 0.7;
}

.campaign-builder__actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: -2px 0 18px;
}

.campaign-builder__actions p {
  margin: 0;
  color: #64748b;
  font-size: 13px;
}

.campaign-builder__grid {
  display: grid;
  grid-template-columns: minmax(320px, 420px) minmax(0, 1fr);
  gap: 18px;
  align-items: start;
}

.campaign-builder__panel {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.08);
  padding: 20px;
}

.campaign-builder__section-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.campaign-builder__section-heading--sub {
  margin-top: 22px;
  padding-top: 18px;
  border-top: 1px solid #e2e8f0;
}

.campaign-builder__section-heading h2 {
  margin: 0;
  color: #111827;
  font-size: 18px;
  line-height: 1.3;
  font-weight: 800;
}

.campaign-builder__section-heading span {
  color: #64748b;
  font-size: 13px;
  white-space: nowrap;
}

.campaign-builder label {
  display: grid;
  gap: 7px;
  margin-bottom: 14px;
  color: #334155;
  font-size: 14px;
  font-weight: 700;
}

.campaign-builder input,
.campaign-builder select,
.campaign-builder textarea {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  background: #ffffff;
  color: #0f172a;
  font-size: 15px;
  line-height: 1.4;
}

.campaign-builder input,
.campaign-builder select {
  min-height: 42px;
  padding: 0 12px;
}

.campaign-builder textarea {
  resize: vertical;
  padding: 10px 12px;
}

.campaign-builder__split {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.campaign-builder__preview {
  display: grid;
  gap: 16px;
}

.campaign-builder article {
  border-top: 1px solid #e2e8f0;
  padding-top: 16px;
}

.campaign-builder article:first-of-type {
  border-top: 0;
  padding-top: 0;
}

.campaign-builder h3 {
  margin: 0 0 8px;
  color: #0f172a;
  font-size: 15px;
  font-weight: 800;
}

.campaign-builder p,
.campaign-builder li,
.campaign-builder dd {
  color: #475569;
  font-size: 14px;
  line-height: 1.6;
}

.campaign-builder ul {
  display: grid;
  gap: 8px;
  margin: 0;
  padding-left: 18px;
}

.campaign-builder dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin: 0;
}

.campaign-builder dt {
  color: #64748b;
  font-size: 12px;
  font-weight: 800;
  text-transform: uppercase;
}

.campaign-builder dd {
  margin: 3px 0 0;
}

@media (max-width: 900px) {
  .campaign-builder__title-row,
  .campaign-builder__grid {
    display: grid;
    grid-template-columns: 1fr;
  }

  .campaign-builder__primary {
    width: 100%;
  }
}

@media (max-width: 560px) {
  .campaign-builder {
    width: min(100vw - 16px, 1180px);
    margin-top: 16px;
  }

  .campaign-builder h1 {
    font-size: 28px;
  }

  .campaign-builder__panel {
    padding: 16px;
  }

  .campaign-builder__split,
  .campaign-builder dl {
    grid-template-columns: 1fr;
  }
}
</style>
