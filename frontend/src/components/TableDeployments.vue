<template>
  <v-card class="mx-auto">
    <v-card-title>
      <v-spacer></v-spacer>
      <v-text-field
        append-icon="mdi-magnify"
        hide-details
        label="Filter deployments"
        single-line
        v-model="searchProp"
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :dense="true"
      :headers="headers"
      :items-per-page="50"
      :items="deployments"
      :loading="loading"
      :search="searchProp"
      :sort-by="['lastChangedTimestamp']"
      :sort-desc="[true]"
    >
      <!-- highlight and set node state colours for the "ready" states -->
      <template v-slot:item.ready="{ item }">
        <v-tooltip bottom>
          <template v-slot:activator="{ on, attrs }">
            <v-chip
              :color="markTrueGood(item.ready)"
              dark
              v-bind="attrs"
              v-on="on"
            >
              {{ item.ready }}
            </v-chip>
          </template>
          <span>
            How many pods are "Ready" (Passing readiness probes)
          </span>
        </v-tooltip>
      </template>

      <template v-slot:item.progressing="{ item }">
        <v-tooltip bottom>
          <template v-slot:activator="{ on, attrs }">
            <v-chip
              :color="markTrueGood(item.progressing)"
              dark
              v-bind="attrs"
              v-on="on"
            >
              {{ item.progressing }}
            </v-chip>
          </template>
          <span>Progressing pods (k8s hasn't given up on them)</span>
        </v-tooltip>
      </template>

      <template v-slot:item.replicas="{ item }">
        <v-tooltip bottom>
          <template v-slot:activator="{ on, attrs }">
            <v-chip
              :color="looksOK(item.replicas, item.replicas_available)"
              dark
              v-bind:active="attrs"
              v-on="on"
            >
              {{ item.replicas }}
            </v-chip>
          </template>
          <span>desired replicas</span>
        </v-tooltip>
      </template>

      <template v-slot:item.replicas_available="{ item }">
        <v-tooltip bottom>
          <template v-slot:activator="{ on, attrs }">
            <v-chip
              :color="looksOK(item.replicas_available, item.replicas)"
              dark
              v-bind:active="attrs"
              v-on="on"
            >
              {{ item.replicas_available }}
            </v-chip>
          </template>
          <span>available pods (running)</span>
        </v-tooltip>
      </template>

      <template v-slot:item.lastChangedTimestamp="{ item }">
        {{ howManySecondsAgoFriendly(item.lastChangedTimestamp) }}
      </template>
    </v-data-table>
  </v-card>
</template>
<script>
export default {
  name: 'TableDeployments',

  data: () => ({
    searchProp: '',
    loading: true, // used to indicate if data is being retrieved
    isActive: null,
    deployments: [],
    headers: [
      { text: 'namespace', value: 'ns', align: 'start' },
      { text: 'name', value: 'name' },
      { text: 'version', value: 'tag' },
      { text: 'ready', value: 'ready' },
      { text: 'progressing', value: 'progressing' },
      { text: 'replicas (desired)', value: 'replicas' },
      { text: 'replicas (available)', value: 'replicas_available' },
      { text: 'last changed', value: 'lastChangedTimestamp' },
    ],
  }),

  methods: {
    // quick and dirty, just show "how long ago" it was till something changed
    // the input timestamp is a UTC timestamp string, if the users
    // browser's timestamp is far off in the past bad things could happen here however
    // so maybe fall back to just showing the timestamp?
    howManySecondsAgoFriendly(timestamp) {
      const seconds = this.howManySecondsAgo(timestamp)
      return this.convertSeconds(seconds)
    },

    howManySecondsAgo(timestamp) {
      const secondsNow = new Date().getTime()
      const secondsThen = new Date(timestamp).getTime()

      return Math.floor((secondsNow - secondsThen) / 1000)
    },

    convertSeconds(inputSeconds) {
      const seconds = inputSeconds.toFixed(1)
      const minutes = (inputSeconds / 60).toFixed(1)
      const hours = (inputSeconds / (60 * 60)).toFixed(1)
      const days = (inputSeconds / (60 * 60 * 24)).toFixed(1)
      if (seconds < 60) {
        return seconds + 's'
      } else if (minutes < 60) {
        return minutes + 'm'
      } else if (hours < 24) {
        return hours + 'h'
      } else {
        return days + 'd'
      }
    },

    async getDeployments() {
      try {
        this.loading = true
        const path = '/v1/table/deployments'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        this.deployments = await resp.json()
        this.loading = false
      } catch (error) {
        console.error(error)
      }
    },

    // looks OK?
    looksOK(current, desired) {
      if (current !== desired) {
        return 'red'
      } else {
        return 'green'
      }
    },

    // returns green or red based on if the input is input is true or false
    markTrueGood(inputString) {
      if (inputString === true) {
        return 'green'
      } else {
        return 'red'
      }
    },
  },

  cron: {
    time: 10000,
    method: 'getDeployments',
  },

  mounted() {
    this.getDeployments()
  },
}
</script>
