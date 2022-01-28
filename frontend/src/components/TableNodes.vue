<template>
  <v-card class="mx-auto">
    <v-card-title>
      Nodes
      <!--      <button v-if="running" v-on:click="stopCronJob">Notes-</button>-->
      <!--      <button v-if="!running" v-on:click="getNodes">Notes+</button>-->
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
    // isActive: false,
    // running: false,
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
        const path = '/v1/table/nodes'
        console.log('retrieving: ' + path)
        // console.log(this.isActive())
        this.loading = true
        const resp = await fetch(path)
        const data = await resp.json()
        this.nodes = data.nodes
        this.headers = data.headers
        this.loading = false
        // this.running = true
      } catch (error) {
        console.error(error)
      }
    }
    // stopCronJob() {
    //   this.$cron.stop('getNodes')
    //   this.running = false
    // }
  },
  cron: {
    time: 1500,
    method: 'getNodes',
    autoStart: true
  },
  mounted() {
    this.getNodes()
  }

  // unmounted() {
  //   this.stopCronJob()
  // }
}
</script>
