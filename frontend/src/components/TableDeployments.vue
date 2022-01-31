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
    deployments: [
      { name: 'foo', ns: 'ns1', created: '12m' },
      { name: 'bar', ns: 'ns2', created: '1h' }
    ],
    headers: [
      { text: 'namespace', value: 'ns', align: 'start' },
      { text: 'name', value: 'name' },
      { text: 'ready', value: 'ready' },
      { text: 'progressing', value: 'progressing' },
      { text: 'replicas (desired)', value: 'replicas' },
      { text: 'replicas (available)', value: 'replicas_available' }
      // {
      //   text: 'created',
      //   value: 'created'
      // }
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
