<template>
  <v-card class="mx-auto">
    <v-card-title>
      Nodes
      <v-spacer></v-spacer>
      <v-text-field
        append-icon="mdi-magnify"
        hide-details
        label="Search"
        single-line
        v-model="search"
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :dense="true"
      :headers="headers"
      :items-per-page="50"
      :items="nodes"
      :loading="loading"
      :multi-sort="true"
      :search="search"
      sort-by="seconds"
    ></v-data-table>
  </v-card>
</template>
<script>
export default {
  name: 'TableNodes',

  data: () => ({
    search: '',
    loading: true,
    nodes: [],
    headers: [
      {
        text: 'Name',
        align: 'start',
        // sortable: true,
        value: 'name'
      },
      { text: 'ready', value: 'ready' },
      { text: 'age', value: 'age' },
      { text: 'age(s)', value: 'seconds' }
    ]
  }),

  methods: {
    async getNodes() {
      try {
        this.loading = true
        const resp = await fetch('/v1/table/nodes')
        const data = await resp.json()
        console.log('retrieving v1/table/nodes')
        this.nodes = data.nodes
        this.headers = data.headers
        this.loading = false
      } catch (error) {
        console.error(error)
      }
    }
  },
  cron: {
    time: 15000,
    method: 'getNodes',
    autoStart: false
  },
  mounted() {
    this.getNodes()
  }
}
</script>
