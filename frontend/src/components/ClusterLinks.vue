<template>
  <!--  <div class="d-flex flex-column mb-6">-->
  <!--    &lt;!&ndash;    <v-container>&ndash;&gt;-->
  <!--    <v-card class="mx-auto">-->
  <v-list>
    <v-list-item v-for="item in stacks" :key="item.title" :href="item.address">
      <v-list-item-icon>
        <v-icon v-if="item.oauth2proxy" color="green">
          mdi-account-lock
        </v-icon>
      </v-list-item-icon>

      <v-list-item-content>
        <v-list-item-title v-text="item.name"></v-list-item-title>
        <v-list-item-subtitle v-html="item.address + ' - ' + item.description "></v-list-item-subtitle>
      </v-list-item-content>

      <v-list-item-avatar>
        <v-img :src="/assets/ + item.icon"></v-img>
      </v-list-item-avatar>
    </v-list-item>
  </v-list>
  <!--    </v-card>-->
  <!--  </div>-->
</template>

<script>
export default {
  name: "ClusterLinks",

  data: function() {
    return {
      stacks: [
        {
          address: "https://mycluster.com/alerts",
          class: "nginx",
          description: "something",
          icon: "prometheus.png",
          oauth2proxy: true
        }
      ]
    };
  },

  methods: {
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
  }
};
</script>
