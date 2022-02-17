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
      sort-by="age"
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
      { text: 'ready', value: 'ready' },
      { text: 'progressing', value: 'progressing' },
      { text: 'replicas (desired)', value: 'replicas' },
      { text: 'replicas (available)', value: 'replicas_available' }
    ]
  }),

  methods: {
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
    }
  },

  cron: {
    time: 10000,
    method: 'getDeployments'
  },

  mounted() {
    this.getDeployments()
  }
}
</script>
