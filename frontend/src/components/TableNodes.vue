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

  props: ['blabla'],

  data: () => ({
    search: '',
    loading: true,
    isDisabled: null,

    isActive: null,
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
        // console.log('nodes: isActive: ' + this.isActive)
        // console.debug('nodes: activeTabName: ' + this.activeTabName)
        // console.debug('nodes: window.location: ' + window.location)
        // console.log('nodes: isDisabled: ' + this.isDisabled)
        // console.debug(window.location.hash)
        // console.log('prop blalba: ' + this.blabla)
        const path = '/v1/table/nodes'
        console.log('nodes: retrieving: ' + path)
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
    time: 150,
    method: 'getNodes',
    autoStart: true
  },
  mounted() {
    console.log('nodes: mounted')
    this.isActive = true
    this.getNodes()
  },
  onUnmounted() {
    console.log('nodes: onUnmounted')
    this.isActive = true
  },
  unmounted() {
    console.log('nodes: unmounted')
    this.isActive = true
  },
  beforeCreate() {
    console.log('nodes: beforeCreate')
  },
  created() {
    console.log('nodes: created')
  },
  beforeMount() {
    console.log('nodes: beforeMount')
  },
  beforeUpdate() {
    console.log('nodes: beforeUpdate')
  },
  updated() {
    console.log('nodes: updated')
  },
  activated() {
    console.log('nodes: activated')
  },
  deactivated() {
    console.log('nodes: deactivated')
  },
  beforeUnmount() {
    console.log('nodes: beforeUnmount')
  },
  errorCaptured() {
    console.log('nodes: errorCaptured')
  },
  renderTracked() {
    console.log('nodes: renderTracked')
  },
  renderTriggered() {
    console.log('nodes: renderTriggered')
  },
  destroyed() {
    this.isActive = false
    //   this.$cron.stop('getNodes')
    //   this.running = false
    console.log('nodes: destroyed')
  }

  // unmounted() {
  //   this.stopCronJob()
  // }
}
</script>
