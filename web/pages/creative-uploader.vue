<template>
  <main class="creative-uploader">
    <section class="creative-uploader__header">
      <p class="creative-uploader__eyebrow">
        Tools
      </p>
      <div class="creative-uploader__title-row">
        <div>
          <h1>Creative Uploader</h1>
          <p>
            Upload or import creatives, split them into broad testing ad sets,
            add copy fields, and prepare them for Meta Ads Manager without manual setup.
          </p>
        </div>
        <button
          class="creative-uploader__secondary"
          type="button"
          @click="copyPlan"
        >
          Copy setup summary
        </button>
      </div>
    </section>

    <section class="creative-uploader__grid">
      <form
        class="creative-uploader__panel"
        @submit.prevent
      >
        <div class="creative-uploader__section-heading">
          <h2>Upload setup</h2>
          <span>Creative testing</span>
        </div>

        <label>
          Ad account
          <select
            v-model="form.adAccount"
            :disabled="!connectedAdAccounts.length"
          >
            <option
              v-for="account in metaAdAccounts"
              :key="account.id"
              :value="account.label"
            >
              {{ account.label }}
            </option>
          </select>
        </label>

        <div class="creative-uploader__actions">
          <button
            class="creative-uploader__secondary"
            :disabled="loadingAdAccounts"
            type="button"
            @click="loadMetaAdAccounts"
          >
            {{ loadingAdAccounts ? 'Loading accounts' : 'Connect agency Business Manager' }}
          </button>
          <p>{{ metaConnectionStatus }}</p>
        </div>

        <label>
          Creative source
          <select v-model="form.creativeSource">
            <option>Upload files</option>
            <option>Use GetChase Creative Generator</option>
            <option>Use selected creative library assets</option>
          </select>
        </label>

        <label>
          Creative files
          <input
            multiple
            type="file"
            @change="handleCreativeFiles"
          >
        </label>

        <div class="creative-uploader__split">
          <label>
            Number of ad sets
            <input
              v-model.number="form.adSetCount"
              max="20"
              min="1"
              type="number"
            >
          </label>
          <label>
            Launch mode
            <select v-model="form.launchMode">
              <option>Upload paused</option>
              <option>Push active immediately</option>
              <option>Schedule launch</option>
            </select>
          </label>
        </div>

        <label>
          Schedule time
          <input
            v-model="form.scheduleTime"
            type="datetime-local"
          >
        </label>

        <label>
          Ad set name pattern
          <input
            v-model="form.adSetNamePattern"
            type="text"
          >
        </label>

        <label>
          Creative assignment
          <textarea
            v-model="form.creativeAssignment"
            rows="5"
          />
        </label>

        <label>
          Creative name
          <input
            v-model="form.creativeName"
            type="text"
          >
        </label>

        <label>
          Primary text
          <textarea
            v-model="form.creativeText"
            rows="4"
          />
        </label>

        <label>
          Headline
          <input
            v-model="form.creativeHeadline"
            type="text"
          >
        </label>

        <label>
          Description
          <input
            v-model="form.creativeDescription"
            type="text"
          >
        </label>

        <label>
          Call to action
          <select v-model="form.creativeCta">
            <option>Install Now</option>
            <option>Learn More</option>
            <option>Sign Up</option>
            <option>Subscribe</option>
          </select>
        </label>
      </form>

      <section class="creative-uploader__panel creative-uploader__preview">
        <div class="creative-uploader__section-heading">
          <h2>Upload plan</h2>
          <span>Ready for API build</span>
        </div>

        <article>
          <h3>Creative testing structure</h3>
          <p>{{ uploadSummary }}</p>
        </article>

        <article>
          <h3>Broad ad set settings</h3>
          <ul>
            <li>Audience: broad, using campaign-level region, age, and gender.</li>
            <li>Placements: Advantage+ placements unless overridden later.</li>
            <li>Optimization: imported from Campaign Builder or selected manually.</li>
            <li>Budget: inherited from campaign or split evenly by ad set.</li>
          </ul>
        </article>

        <article>
          <h3>Ad set and creative mapping</h3>
          <ul>
            <li
              v-for="assignment in creativeAssignments"
              :key="assignment"
            >
              {{ assignment }}
            </li>
          </ul>
        </article>

        <article>
          <h3>Ad copy payload</h3>
          <dl>
            <div>
              <dt>Name</dt>
              <dd>{{ form.creativeName }}</dd>
            </div>
            <div>
              <dt>Headline</dt>
              <dd>{{ form.creativeHeadline }}</dd>
            </div>
            <div>
              <dt>Description</dt>
              <dd>{{ form.creativeDescription }}</dd>
            </div>
            <div>
              <dt>CTA</dt>
              <dd>{{ form.creativeCta }}</dd>
            </div>
          </dl>
        </article>

        <article>
          <h3>Rules Builder</h3>
          <label>
            Rule condition
            <select v-model="primaryRule.condition">
              <option
                v-for="option in ruleConditionOptions"
                :key="option"
              >
                {{ option }}
              </option>
            </select>
          </label>
          <div class="creative-uploader__split">
            <label>
              Operator
              <select v-model="primaryRule.operator">
                <option
                  v-for="operator in ruleOperators"
                  :key="operator"
                  :value="operator"
                >
                  {{ operator }}
                </option>
              </select>
            </label>
            <label>
              Value
              <input
                v-model="primaryRule.value"
                step="0.01"
                type="number"
              >
            </label>
          </div>

          <div
            v-for="(rule, index) in additionalRules"
            :key="index"
          >
            <label>
              Additional condition {{ index + 1 }}
              <select v-model="rule.condition">
                <option
                  v-for="option in ruleConditionOptions"
                  :key="option"
                >
                  {{ option }}
                </option>
              </select>
            </label>
            <div class="creative-uploader__split">
              <label>
                Operator
                <select v-model="rule.operator">
                  <option
                    v-for="operator in ruleOperators"
                    :key="operator"
                    :value="operator"
                  >
                    {{ operator }}
                  </option>
                </select>
              </label>
              <label>
                Value
                <input
                  v-model="rule.value"
                  step="0.01"
                  type="number"
                >
              </label>
            </div>
          </div>

          <div class="creative-uploader__actions">
            <button
              class="creative-uploader__text-button"
              type="button"
              @click="addCondition"
            >
              Add condition
            </button>
          </div>

          <label>
            Rule action
            <select v-model="ruleAction">
              <option>Notify only</option>
              <option>Pause losing ad</option>
              <option>Pause losing ad set</option>
              <option>Increase budget on winner</option>
            </select>
          </label>
          <p>{{ rulesSummary }}</p>
        </article>

        <article>
          <h3>Meta API sequence</h3>
          <ul>
            <li>Upload image or video asset to the selected ad account.</li>
            <li>Create broad ad sets using the shared testing settings.</li>
            <li>Create ad creatives with text, headline, description, CTA, and asset reference.</li>
            <li>Create ads in the selected status: active, scheduled, or paused.</li>
            <li>Attach rules after upload if automation is enabled.</li>
          </ul>
        </article>
      </section>
    </section>
  </main>
</template>

<script>
export default {
  name: 'CreativeUploader',
  data () {
    return {
      connectedAdAccounts: [],
      loadingAdAccounts: false,
      metaConnectionStatus: 'Agency Business Manager not connected',
      ruleConditionOptions: [
        'Spend reaches test threshold',
        'CPA is above target',
        'CTR is below target',
        'Creative has no conversions',
      ],
      ruleOperators: ['>=', '>', '<=', '<', '='],
      primaryRule: {
        condition: 'Spend reaches test threshold',
        operator: '>=',
        value: '50',
      },
      additionalRules: [],
      ruleAction: 'Notify only',
      form: {
        adAccount: '',
        creativeSource: 'Upload files',
        adSetCount: 3,
        launchMode: 'Upload paused',
        scheduleTime: '',
        adSetNamePattern: 'Broad Creative Test {{number}}',
        creativeAssignment: 'Ad set 1: Creative A, Creative B\nAd set 2: Creative C\nAd set 3: Creative D',
        creativeName: 'GetChase live alerts creative test',
        creativeText: 'Get real-time alerts when fast-moving live events happen near you.',
        creativeHeadline: 'Live alerts before everyone is talking',
        creativeDescription: 'Maps, updates, and community context in one app.',
        creativeCta: 'Install Now',
      },
    }
  },
  head () {
    return {
      title: 'Creative Uploader',
    }
  },
  computed: {
    metaAdAccounts () {
      return this.connectedAdAccounts.length
        ? this.connectedAdAccounts
        : [{ id: 'not-connected', label: 'Connect agency Business Manager to load ad accounts' }]
    },
    adSets () {
      const count = Math.max(1, Number(this.form.adSetCount || 1))
      return Array.from({ length: count }, (_, index) => {
        return this.form.adSetNamePattern.replace('{{number}}', String(index + 1))
      })
    },
    creativeAssignments () {
      const assignments = this.form.creativeAssignment
        .split('\n')
        .map(line => line.trim())
        .filter(Boolean)

      return assignments.length
        ? assignments
        : this.adSets.map(adSet => `${adSet}: assign selected creatives`)
    },
    launchDetail () {
      return this.form.launchMode === 'Schedule launch' && this.form.scheduleTime
        ? `${this.form.launchMode} at ${this.form.scheduleTime}`
        : this.form.launchMode
    },
    uploadSummary () {
      const adSetCount = this.adSets.length
      return `${this.form.creativeSource} will create ${adSetCount} broad ad set${adSetCount === 1 ? '' : 's'}, then upload ads as "${this.launchDetail}".`
    },
    rules () {
      return [this.primaryRule, ...this.additionalRules]
    },
    rulesSummary () {
      return `When ${this.rules.map(this.formatRuleCondition).map(rule => `"${rule}"`).join(' AND ')} happen, Creative Uploader should "${this.ruleAction.toLowerCase()}".`
    },
    setupSummary () {
      return [
        'Creative Uploader plan',
        `Creative source: ${this.form.creativeSource}`,
        `Ad sets: ${this.adSets.join(', ')}`,
        `Launch mode: ${this.launchDetail}`,
        `Creative name: ${this.form.creativeName}`,
        `Primary text: ${this.form.creativeText}`,
        `Headline: ${this.form.creativeHeadline}`,
        `Description: ${this.form.creativeDescription}`,
        `CTA: ${this.form.creativeCta}`,
        `Assignments: ${this.creativeAssignments.join(' | ')}`,
        `Rule: ${this.rules.map(this.formatRuleCondition).join(' AND ')} -> ${this.ruleAction}`,
      ].join('\n')
    },
  },
  mounted () {
    this.loadMetaAdAccounts()
  },
  methods: {
    formatRuleCondition (rule) {
      return `${rule.condition} ${rule.operator} ${rule.value}`
    },
    addCondition () {
      this.additionalRules.push({
        condition: this.ruleConditionOptions[Math.min(this.additionalRules.length + 1, this.ruleConditionOptions.length - 1)],
        operator: '>',
        value: this.additionalRules.length === 0 ? '12' : '0',
      })
    },
    handleCreativeFiles (event) {
      const files = Array.from(event.target.files || []).map(file => file.name)
      const adSetCount = Math.max(1, Number(this.form.adSetCount || 1))

      if (files.length) {
        this.form.creativeAssignment = files.map((file, index) => {
          return `Ad set ${(index % adSetCount) + 1}: ${file}`
        }).join('\n')
      }
    },
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
    async copyPlan () {
      if (!navigator.clipboard) {
        return
      }

      await navigator.clipboard.writeText(this.setupSummary)
    },
  },
}
</script>

<style scoped>
.creative-uploader {
  width: min(1180px, calc(100vw - 24px));
  margin: 24px auto 56px;
  color: #111827;
}

.creative-uploader__header {
  margin-bottom: 20px;
}

.creative-uploader__eyebrow {
  margin: 0 0 6px;
  color: #4f46e5;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0;
  text-transform: uppercase;
}

.creative-uploader__title-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 18px;
}

.creative-uploader h1 {
  margin: 0;
  color: #0f172a;
  font-size: 34px;
  line-height: 1.1;
  font-weight: 800;
  letter-spacing: 0;
}

.creative-uploader__title-row p {
  max-width: 720px;
  margin: 10px 0 0;
  color: #475569;
  font-size: 16px;
  line-height: 1.6;
}

.creative-uploader__secondary,
.creative-uploader__text-button {
  min-height: 40px;
  padding: 0 14px;
  border: 1px solid #2563eb;
  border-radius: 6px;
  background: #eff6ff;
  color: #1d4ed8;
  font-weight: 700;
  cursor: pointer;
}

.creative-uploader__secondary:disabled {
  cursor: wait;
  opacity: 0.7;
}

.creative-uploader__actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: -2px 0 18px;
}

.creative-uploader__actions p {
  margin: 0;
  color: #64748b;
  font-size: 13px;
}

.creative-uploader__grid {
  display: grid;
  grid-template-columns: minmax(320px, 420px) minmax(0, 1fr);
  gap: 18px;
  align-items: start;
}

.creative-uploader__panel {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.08);
  padding: 20px;
}

.creative-uploader__section-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.creative-uploader__section-heading h2 {
  margin: 0;
  color: #111827;
  font-size: 18px;
  line-height: 1.3;
  font-weight: 800;
}

.creative-uploader__section-heading span {
  color: #64748b;
  font-size: 13px;
  white-space: nowrap;
}

.creative-uploader label {
  display: grid;
  gap: 7px;
  margin-bottom: 14px;
  color: #334155;
  font-size: 14px;
  font-weight: 700;
}

.creative-uploader input,
.creative-uploader select,
.creative-uploader textarea {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  background: #ffffff;
  color: #0f172a;
  font-size: 15px;
  line-height: 1.4;
}

.creative-uploader input,
.creative-uploader select {
  min-height: 42px;
  padding: 0 12px;
}

.creative-uploader textarea {
  resize: vertical;
  padding: 10px 12px;
}

.creative-uploader__split {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.creative-uploader__preview {
  display: grid;
  gap: 16px;
}

.creative-uploader article {
  border-top: 1px solid #e2e8f0;
  padding-top: 16px;
}

.creative-uploader article:first-of-type {
  border-top: 0;
  padding-top: 0;
}

.creative-uploader h3 {
  margin: 0 0 8px;
  color: #0f172a;
  font-size: 15px;
  font-weight: 800;
}

.creative-uploader p,
.creative-uploader li,
.creative-uploader dd {
  color: #475569;
  font-size: 14px;
  line-height: 1.6;
}

.creative-uploader ul {
  display: grid;
  gap: 8px;
  margin: 0;
  padding-left: 18px;
}

.creative-uploader dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin: 0;
}

.creative-uploader dt {
  color: #64748b;
  font-size: 12px;
  font-weight: 800;
  text-transform: uppercase;
}

.creative-uploader dd {
  margin: 3px 0 0;
}

@media (max-width: 900px) {
  .creative-uploader__title-row,
  .creative-uploader__grid {
    display: grid;
    grid-template-columns: 1fr;
  }

  .creative-uploader__secondary {
    width: 100%;
  }
}

@media (max-width: 560px) {
  .creative-uploader {
    width: min(100vw - 16px, 1180px);
    margin-top: 16px;
  }

  .creative-uploader h1 {
    font-size: 28px;
  }

  .creative-uploader__panel {
    padding: 16px;
  }

  .creative-uploader__split,
  .creative-uploader dl {
    grid-template-columns: 1fr;
  }
}
</style>
