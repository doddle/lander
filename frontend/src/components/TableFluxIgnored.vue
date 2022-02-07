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
      :items="apiGroups"
      :loading="loading"
      :search="searchProp"
      :sort-by="['Namespace', 'Resource']"
      multi-sort
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
    apiGroups: [],
    headers: [
      {
        text: 'Namespace',
        align: 'start',
        value: 'Namespace'
      },
      {
        text: 'Resource',
        value: 'Resource'
      },
      { text: 'Kind', value: 'Kind' },
      { text: 'APIGroup', value: 'APIGroup' },
      { text: 'APIGroupVersion', value: 'APIGroupVersion' }
    ]
  }),

  methods: {
    async getFluxIgnored() {
      try {
        this.loading = true
        const path = '/v1/table/inventory-flux-ignored'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        this.apiGroups = await resp.json()
        this.loading = false
      } catch (error) {
        console.error(error)
      }
    }
  },

  cron: {
    time: 30000,
    method: 'getFluxIgnored'
  },

  mounted() {
    this.getFluxIgnored()
  }
}
</script>
