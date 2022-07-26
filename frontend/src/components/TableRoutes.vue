<template>
  <v-card class="mx-auto">
    <v-card-title>
      <v-spacer></v-spacer>
      <v-text-field
        append-icon="mdi-magnify"
        hide-details
        label="Filter routes"
        single-line
        v-model="searchProp"
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :dense="true"
      :headers="headers"
      :items-per-page="50"
      :items="routes"
      :loading="loading"
      :search="searchProp"
      sort-by="age"
    >
    </v-data-table>
  </v-card>
</template>
<script>
export default {
  name: 'TableRoutes',

  data: () => ({
    searchProp: '',
    loading: true, // used to indicate if data is being retrieved
    isActive: null,
    routes: [
      {
        hostname: 'foo.acme.org',
        path: '/',
        ns: 'ns1',
        svc: 'svc1',
        class: 'private',
      },
      {
        hostname: 'bar.acme.org',
        path: '/',
        ns: 'ns2',
        svc: 'svc2',
        class: 'public',
      },
    ],
    headers: [
      {
        text: 'HostName',
        align: 'start',
        value: 'hostname',
      },
      {
        text: 'IngressClass',
        value: 'class',
      },
      { text: 'path', value: 'path' },
      { text: 'namespace', value: 'ns' },
      {
        text: 'service',
        value: 'svc',
      },
    ],
  }),

  methods: {
    async getRoutes() {
      try {
        this.loading = true
        const path = '/v1/routes'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        this.routes = await resp.json()
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
    },
  },

  cron: {
    time: 10000,
    method: 'getRoutes',
  },

  mounted() {
    this.getRoutes()
  },
}
</script>
