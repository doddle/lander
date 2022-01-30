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
    loading: true, // used to indicate if data is being retreived
    isActive: null,
    deployments: [
      { name: 'foo', ready: true, ns: 'ns1', age: '12m' },
      { name: 'bar', ready: false, ns: 'ns2', age: '1h' }
    ],
    headers: [
      {
        text: 'Name',
        align: 'start',
        value: 'name'
      },
      { text: 'ready', value: 'true' },
      { text: 'namespace', value: 'ns' },
      {
        text: 'age',
        value: 'age'
      }
    ]
  }),

  methods: {
    async getDeployments() {
      try {
        this.loading = true
        const path = '/v1/table/deployments'
        console.debug('retrieving: ' + path)
        // const resp = await fetch(path)
        // const data = await resp.json()
        // this.deployments = data.deployments
        // this.headers = data.headers
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
