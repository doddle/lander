<template>
  <!--  <div class="d-flex flex-column mb-6">-->
  <!--    <v-container>-->
  <!--  <v-card class="mx-auto">-->

  <v-toolbar :color="settings.colorscheme" tabs flat>
    <v-app-bar-nav-icon></v-app-bar-nav-icon>
    <v-toolbar-title>insight</v-toolbar-title>
    <v-spacer></v-spacer>
    <!--    <v-btn icon> <v-icon>mdi-magnify</v-icon> </v-btn> <v-btn icon> <v-icon>mdi-dots-vertical</v-icon> </v-btn>-->
    <template v-slot:extension>
      <v-tabs v-model="tab" align-with-title>
        <!--        <v-tabs-slider></v-tabs-slider>-->
        <v-tabs-slider color="grey"></v-tabs-slider>
        <v-tab href="#tab-1">
          One
        </v-tab>
        <v-tab href="#tab-2">
          Two
        </v-tab>
      </v-tabs>

      <v-tabs-items v-model="tab">
        <v-tab-item key="1" value="tab-1">
          <v-card-text>1111</v-card-text>
        </v-tab-item>
        <v-tab-item key="2" value="tab-2">
          <v-card-text>222222222222</v-card-text>
        </v-tab-item>
      </v-tabs-items>
    </template>
  </v-toolbar>
</template>

<script>
// import ClusterLinks from "./components/ClusterLinks";
export default {
  name: "Home",

  data: function() {
    return {
      tab: null,
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
  },
  components: {
    // ClusterLinks
  }
};
</script>
