<template>
  <div class="d-flex flex-column mb-6">
    <v-container>
      <v-card class="mx-auto">
        <v-toolbar color="indigo" dark>
          <!-- <v-app-bar-nav-icon></v-app-bar-nav-icon> -->

          <v-toolbar-title>Cluster Links</v-toolbar-title>

          <v-spacer></v-spacer>

          <!-- <v-btn icon>
            <v-icon>mdi-magnify</v-icon>
          </v-btn> -->

          <!-- <v-btn icon>
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn> -->
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

  data() {
    return {
      stacks: []
    };
  },

  methods: {
    async getStacks() {
      try {
        const resp = await fetch("/v1/endpoints");
        const data = await resp.json();
        console.log("retrieving endpoints");
        this.stacks = data;
      } catch (error) {
        console.error(error);
      }
    }
  },
  cron: {
    time: 5000,
    method: "getStacks",
    autoStart: true
  },
  beforeMount() {
    this.getStacks();
  }
};
</script>
