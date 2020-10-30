<template>
  <div class="d-flex flex-column mb-6">
    <v-container>
      <v-card class="mx-auto">
        <v-toolbar v-bind:color="fuck" >
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

  data() {
    var hostLocation = location.host;
    var hostName = hostLocation.split(":")[0];
    var fuck = "yello"

    // var dict = {
    //   "localhost": "blue",
    //   "prod": "red",
    //   "staging": "green"
    // }

    // fuck = dict.localhost
    // if (hostName.contains("localhost")) {
    //   fuck = "green"
    // }




    switch (hostName) {
      case "localhost":
        fuck = "red"
        break;
    
      default:
        fuck = "blue"
        break;
    }

    return {
      stacks: [],
      fuck
    };


  // const href = window.location.href; const findTerm = (term) => { if (href.includes(term)){ return href; } }; switch (href) { case findTerm('google'): searchWithGoogle(); break; case findTerm('yahoo'): searchWithYahoo(); break; default: console.log('No search engine found'); }; 


  },

  // x() {
  //   return `${this.fuck}`;
  // },

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
