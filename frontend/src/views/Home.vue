<template>
  <div class="home">
    <img alt="favicon" src="/favicon.ico" />
  </div>
</template>

<script>
export default {
  name: 'Home',

  data() {
    return {
      stacks: [],
    }
  },
  methods: {
    async getStacks() {
      try {
        const resp = await fetch('/v1/endpoints')
        const data = await resp.json()
        console.log('retrieving endpoints')
        this.stacks = data
      } catch (error) {
        console.error(error)
      }
    },
  },
  cron: {
    time: 15000,
    method: 'getStacks',
    autoStart: true,
  },
  mounted() {
    this.getStacks()
  },
}
</script>
