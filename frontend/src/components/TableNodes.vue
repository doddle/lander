<template>
  <v-card class="mx-auto">
    <v-card-title>
      Nodes
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-icon="mdi-magnify"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="nodes"
      :search="search"
      :dense="true"
      :items-per-page="50"
      :multi-sort="true"
      sort-by="seconds"
    ></v-data-table>
  </v-card>
</template>
<script>
export default {
  name: "TableNodes",

  data: () => ({
    search: "",
    nodes: [],
    headers: [
      {
        text: "Name",
        align: "start",
        // sortable: true,
        value: "name"
      },
      { text: "ready", value: "ready" },
      { text: "age", value: "age" },
      { text: "age(s)", value: "seconds" }
    ]
  }),

  methods: {
    async getNodes() {
      try {
        const resp = await fetch("/v1/table/nodes");
        const data = await resp.json();
        console.log("retrieving v1/table/nodes");
        this.nodes = data.nodes;
        this.headers = data.headers;
      } catch (error) {
        console.error(error);
      }
    }
  },
  cron: {
    time: 15000,
    method: "getNodes",
    autoStart: true
  },
  beforeMount() {
    this.getNodes();
  }
};
</script>
