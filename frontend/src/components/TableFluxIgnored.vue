<template>
  <v-card class="mx-auto">
    <v-card-title>
      <v-spacer></v-spacer>
      <v-text-field
        append-icon="mdi-magnify"
        hide-details
        label="Filter Flux Ignored"
        single-line
        v-model="searchProp"
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :dense="true"
      :headers="headers"
      :items-per-page="50"
      :items="fluxIgnored"
      :loading="loading"
      :search="searchProp"
      sort-by="age"
    >
    </v-data-table>
  </v-card>
</template>
<script>
export default {
  name: 'TableFluxIgnored',

  data: () => ({
    searchProp: '',
    loading: true, // used to indicate if data is being retrieved
    isActive: null,
    fluxIgnored: [
      { ns: 'default', name: 'exampleResource', kind: 'deployment' }
    ],
    headers: [
      {
        text: 'Namespace',
        align: 'start',
        value: 'ns'
      },
      {
        text: 'Name',
        value: 'name'
      },
      { text: 'Kind', value: 'kind' }
    ]
  }),

  methods: {
    async getFluxIgnored() {
      try {
        this.loading = true
        const path = '/v1/table/fluxIgnored'
        console.log('retrieving: ' + path)
        // const resp = await fetch(path)
        // this.fluxIgnored = await resp.json()
        // this.headers = data.headers
        this.loading = false
      } catch (error) {
        console.error(error)
      }
    }
  },

  cron: {
    time: 10000,
    method: 'getFluxIgnored'
  },

  mounted() {
    this.getFluxIgnored()
  }
}
</script>
