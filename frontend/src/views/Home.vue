<template>
  <div class="home">
    <img alt="Vue logo" src="/favicon.ico" />
  </div>
</template>

<script>
// @ is an alias to /src
// import HelloWorld from "@/components/HelloWorld.vue";

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

  // components: {
  //   HelloWorld
  // }
};
</script>
