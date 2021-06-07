<template>
  <div class="text-center">
    <v-menu offset-y>
      <template v-slot:activator="{ on, attrs }">
        <v-btn color="primary" v-bind="attrs" v-on="on">
          {{ settings.cluster }}
        </v-btn>
      </template>
      <v-list>
        <v-list-item-group color="indigo">
          <v-list-item v-for="(item, index) in settings.clusters" :key="index">
            <v-list-item-title>
              <a :href="'https://' + item">
                {{ item }}
              </a>
            </v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
export default {
  name: "ClusterLinks",

  data: function() {
    return {
      stacks: [],
      settings: {
        colorscheme: "blue lighten-5",
        cluster: "unknown",
        clusters: ["cluster1.acmecorp.org"]
      }
    };
  },

  methods: {
    async getSettings() {
      try {
        const resp = await fetch("/v1/settings");
        const data = await resp.json();
        console.log("retrieving settings");
        this.settings = data;
      } catch (error) {
        console.error(error);
      }
    }
  },
  beforeMount() {
    this.getSettings();
  }
};
</script>

<style>

.btn {
  text-transform: none;
}
</style>
