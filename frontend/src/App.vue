<template>
  <v-app>
    <v-app-bar app color="blue-grey lighten-1">
      <div class="d-flex align-center">
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

      <v-spacer></v-spacer>
      <ClusterLinks />

      <v-spacer></v-spacer>

      <v-btn href="https://github.com/digtux/lander" target="_blank" text>
        <span class="mr-2"></span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>
    <v-main>
      <br />
      <Home />
      <v-divider> </v-divider>
      <v-container class="mb-6">
        <v-row no-gutters align-content="center">
          <v-col sm="4" md="3" key="1">
            <OverviewPieDeployments />
          </v-col>
          <v-col sm="4" md="3" key="2">
            <OverviewPieStatefulSets />
          </v-col>
          <v-col sm="4" md="3" key="3">
            <OverviewPieNodes />
          </v-col>
          <!--
          <v-col sm="4" md="3" key="3">
            <OverviewPieDeployments2 />
          </v-col>
          <v-col sm="4" md="3" key="4">
            <OverviewPieStatefulSets2 />
          </v-col>
          -->
        </v-row>
      </v-container>
      <v-divider> </v-divider>
      <TableNodes />
    </v-main>
  </v-app>
</template>

<script>
import Home from "./components/Home";
import ClusterLinks from "./components/CluserLinks";
import OverviewPieDeployments from "./components/OverviewPieDeployments";
import OverviewPieNodes from "./components/OverviewPieNodes";
import OverviewPieStatefulSets from "./components/OverviewPieStatefulSets";
import TableNodes from "./components/TableNodes";

export default {
  name: "App",
  data() {
    const hostLocation = location.host;
    const hostName = hostLocation.split(":")[0];
    return {
      host: hostName
    };
  },
  title() {
    return `${this.host}`;
  },
  components: {
    Home,
    OverviewPieDeployments,
    OverviewPieStatefulSets,
    OverviewPieNodes,
    TableNodes,
    ClusterLinks
  }
};
</script>
