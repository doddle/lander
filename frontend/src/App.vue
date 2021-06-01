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

      <span class="text-center">{{ host }}</span>

      <v-spacer></v-spacer>

      <v-btn href="https://github.com/digtux/lander" target="_blank" text>
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
