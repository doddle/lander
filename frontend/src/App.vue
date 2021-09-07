<template>
  <v-app>
    <v-app-bar app color="blue-grey lighten-1" class="align-centre">
      <v-row class="justify-space-between">
        <v-col cols=1>
          <div class="d-flex">
            <a :href="'/favicon-' + host + '.png'">
              <img
                :src="`favicon-${host}.ico`"
                alt="identicon"
                class="shrink mr-2"
                contain
                transition="scale-transition"
                width="40"
              />
            </a>
          </div>
        </v-col>

        <v-col cols=3 class="d-flex justify-space-around"> 
          <v-menu offset-y>
            <template v-slot:activator="{ on, attrs }">
              <v-btn color="primary" dark v-bind="attrs" v-on="on">
                {{ settings.cluster }}
              </v-btn>
            </template>
            <v-list>
              <v-list-item v-for="(item, index) in settings.clusters" :key="index">
                <v-list-item-title>
                  <a :href="'https://' + item">{{ item }}</a>
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </v-col>
        
        <v-col cols=1 style="text-align:right">
          <v-btn href="https://github.com/digtux/lander" target="_blank" text>
            <span class="mr-2"></span>
            <v-icon>mdi-open-in-new</v-icon>
          </v-btn>
        </v-col>

      </v-row>
    </v-app-bar>

    <v-main>
      <!--      <div class="d-flex flex-column mb-6">-->
      <v-container>
        <v-row no-gutters align-content="center" justify="center">
          <v-col sm="3" md="3" key="1">
            <OverviewPieDeployments />
          </v-col>
          <v-col sm="3" md="3" key="2">
            <OverviewPieStatefulSets />
          </v-col>
          <v-col sm="3" md="3" key="3">
            <OverviewPieNodes />
          </v-col>
        </v-row>
      </v-container>
      <v-container fluid>
        <v-toolbar :color="settings.colorscheme" tabs>
          <template>
            <!--          <template v-slot:extension>-->
            <v-tabs v-model="tab" centered slider-color="black">
              <!--        <v-tabs-slider></v-tabs-slider>-->
              <v-tabs-slider color="grey"></v-tabs-slider>
              <v-tab href="#tab-1">
                Links
              </v-tab>
              <v-tab href="#tab-2">
                Nodes
              </v-tab>
            </v-tabs>
          </template>
        </v-toolbar>

        <v-tabs-items v-model="tab">
          <v-tab-item key="1" value="tab-1">
            <ClusterLinks></ClusterLinks>
          </v-tab-item>
          <v-tab-item key="2" value="tab-2">
            <TableNodes></TableNodes>
          </v-tab-item>
        </v-tabs-items>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
// import Home from "./components/Home";
import ClusterLinks from "./components/ClusterLinks";
import OverviewPieDeployments from "./components/OverviewPieDeployments";
import OverviewPieNodes from "./components/OverviewPieNodes";
import OverviewPieStatefulSets from "./components/OverviewPieStatefulSets";
import TableNodes from "./components/TableNodes";

export default {
  name: "App",
  data: function() {
    const hostLocation = location.host;
    const hostName = hostLocation.split(":")[0];
    return {
      tab: null,
      host: hostName,
      settings: {
        colorscheme: "blue lighten-5",
        cluster: "unknown",
        clusters: ["cluster1.acmecorp.org"]
      }
    };
  },
  title() {
    return `${this.host}`;
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
    // Home,
    OverviewPieDeployments,
    OverviewPieStatefulSets,
    OverviewPieNodes,
    TableNodes,
    ClusterLinks
  }
};
</script>
