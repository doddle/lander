<template>

  <!-- <v-container>
    <p v-for="post in stacks" :key="post.id">
      {{ post.address }}
    </p>
  </v-container> -->


      <v-data-table
        :headers="headers"
        :items="stacks"
        :search="search"
        class="elevation-3"
        :rows-per-page-items="[100, 200, 300, 400]"

      >
      </v-data-table>

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
    time: 1000,
    method: "getStacks",
    autoStart: true
  },
  mounted() {
    this.GetStacks();
  }
};
</script>
