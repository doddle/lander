<template>
  <v-app>
    <v-app-bar app color="blue-grey lighten-1">
      <div class="d-flex align-center">
        <v-img
          alt="identicon"
          class="shrink mr-2"
          contain
          src="/favicon-k8s.ico"
          transition="scale-transition"
          width="40"
        />

        <!-- <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-name-dark.png"
          width="100"
        /> -->
      </div>

      <v-spacer></v-spacer>

      <span class="text-center">{{ host }}</span>

      <v-spacer></v-spacer>

      <v-btn href="https://github.com/starkers/lander" target="_blank" text>
        <span class="mr-2"></span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>
    <v-main>
      <br />
      <Home />
    </v-main>
  </v-app>
</template>

<script>
import Home from "./components/Home";

export default {
  name: "App",
  data() {
    var hostLocation = location.host;
    var hostName = hostLocation.split(":")[0];
    return {
      host: hostName
    };
  },
  title() {
    return `${this.host}`;
  },
  components: {
    Home
  }
};
</script>
