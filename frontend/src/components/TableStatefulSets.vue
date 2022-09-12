<template>
  <v-card class="mx-auto">
    <v-card-title>
      <v-spacer></v-spacer>
      <v-text-field
        append-icon="mdi-magnify"
        hide-details
        label="Filter objects"
        single-line
        v-model="searchProp"
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :dense="true"
      :headers="headers"
      :items-per-page="50"
      :items="objects"
      :loading="loading"
      :search="searchProp"
      sort-by="age"
    >
      <!-- highlight and set node state colours for the "ready" states -->
      <template v-slot:item.ready="{ item }">
        <v-chip :color="markTrueGood(item.ready)" dark>
          {{ item.ready }}
        </v-chip>
      </template>
      <template v-slot:item.progressing="{ item }">
        <v-chip :color="markTrueGood(item.progressing)" dark>
          {{ item.progressing }}
        </v-chip>
      </template>
      <template v-slot:item.replicas="{ item }">
        <v-chip :color="looksOK(item.replicas, item.replicas_updated)" dark>
          {{ item.replicas }}
        </v-chip>
      </template>
    </v-data-table>
  </v-card>
</template>
<script>
export default {
  name: 'TableStatefulSets',

  data: () => ({
    searchProp: '',
    loading: true, // used to indicate if data is being retrieved
    isActive: null,
    objects: [],
    headers: [
      { text: 'namespace', value: 'ns', align: 'start' },
      { text: 'name', value: 'name' },
      { text: 'ready', value: 'ready' },
      { text: 'progressing', value: 'progressing' },
      { text: 'replicas (desired)', value: 'replicas' },
      { text: 'replicas (updated)', value: 'replicas_updated' },
      { text: 'replicas (current)', value: 'replicas_current' },
    ],
  }),

  methods: {
    async getStatefulSets() {
      try {
        this.loading = true
        const path = '/v1/statefulSetTable'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        this.objects = await resp.json()
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
    method: 'getStatefulSets',
  },

  mounted() {
    this.getStatefulSets()
  },
}
</script>
