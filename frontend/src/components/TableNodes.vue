<template>
  <v-card class="mx-auto">
    <v-card-title>
      <v-spacer></v-spacer>
      <v-text-field
        append-icon="mdi-magnify"
        hide-details
        label="Filter nodes"
        single-line
        v-model="searchProp"
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :dense="true"
      :headers="headers"
      :items-per-page="25"
      :items="nodes"
      :loading="loading"
      :search="searchProp"
      sort-by="age"
    >
      <!-- mod seconds into something human readable yet still searchable -->
      <template v-slot:item.age="{ item }">
        {{ convertSeconds(item.age) }}
      </template>

      <!-- highlight unready node states -->
      <template v-slot:item.ready="{ item }">
        <v-chip :color="markTrueGood(item.ready)" dark>
          {{ item.ready }}
        </v-chip>
      </template>
      <!-- highlight unschedulable node states -->
      <template v-slot:item.schedulable="{ item }">
        <v-chip :color="markTrueGood(item.schedulable)" dark>
          {{ item.schedulable }}
        </v-chip>
      </template>
    </v-data-table>
  </v-card>
</template>
<script>
export default {
  name: 'TableNodes',

  data: () => ({
    searchProp: '',
    loading: true, // used to indicate if data is being retreived
    isActive: null,
    nodes: [],
    headers: [
      {
        text: 'Name',
        align: 'start',
        value: 'name'
      },
      { text: 'ready', value: 'true' },
      // {
      //   text: 'age',
      //   value: 'age'
      //   // sortable: true
      // },
      {
        sortable: true,
        text: 'age',
        value: 'age'
      }
    ]
  }),

  methods: {
    async getNodes() {
      try {
        this.loading = true
        const path = '/v1/table/nodes'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        const data = await resp.json()
        this.nodes = data.nodes
        this.headers = data.headers
        this.loading = false
      } catch (error) {
        console.error(error)
      }
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
    markTrueGood(inputString) {
      console.log(inputString)
      if (inputString === true) {
        return 'green'
      } else {
        return 'red'
      }
    }
  },
  cron: {
    time: 5000,
    method: 'getNodes'
  },
  mounted() {
    this.getNodes()
  }
}
</script>
