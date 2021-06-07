<template>
  <div class="d-flex flex-column mb-6">
    <v-container>
      <v-card class="mx-auto">
        <v-toolbar v-bind:color="settings.colorscheme">
          <!-- <v-app-bar-nav-icon></v-app-bar-nav-icon> -->

          <v-toolbar-title>Cluster Links</v-toolbar-title>

          <v-spacer></v-spacer>
        </v-toolbar>
        <v-list>
          <v-list-item
            v-for="item in stacks"
            :key="item.title"
            :href="item.address"
          >
            <v-list-item-icon>
              <v-icon v-if="item.oauth2proxy" color="green">
                mdi-account-lock
              </v-icon>
            </v-list-item-icon>

            <v-list-item-content>
              <v-list-item-title v-text="item.address"></v-list-item-title>

              <v-list-item-subtitle
                v-html="item.description"
              ></v-list-item-subtitle>
            </v-list-item-content>

            <v-list-item-avatar>
              <v-img :src="item.icon"></v-img>
            </v-list-item-avatar>
          </v-list-item>
        </v-list>
      </v-card>
    </v-container>
  </div>
</template>

<script>
export default {
  name: "Home",

  data: function() {
    return {
      stacks: [],
      settings: { colorscheme: "blue lighten-5", cluster: "unknown" }
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
    },
    async getEndpoints() {
      try {
        const resp = await fetch("/v1/endpoints");
        const data = await resp.json();
        console.log("retrieving v1/endpoints");
        this.stacks = data;
      } catch (error) {
        console.error(error);
      }
    }
  },
  cron: {
    time: 10000,
    method: "getEndpoints",
    autoStart: true
  },
  beforeMount() {
    this.getEndpoints();
    this.getSettings();
  }
};
</script>
