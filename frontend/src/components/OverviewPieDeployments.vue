<template>
  <div>
    <v-card flat>
      <v-card-title>Deployments</v-card-title>
      <v-card-subtitle>total: {{ total }}</v-card-subtitle>
      <apexcharts type="pie" :options="chartOptions" :series="series" />
    </v-card>
  </div>
</template>

<script>
import VueApexCharts from 'vue-apexcharts'

export default {
  name: 'Apex',
  components: {
    apexcharts: VueApexCharts,
  },
  data: function() {
    return {
      chartOptions: {
        colors: [],
        chart: {
          id: 'pie-deployments',
          dropShadow: {
            effect: false,
          },
        },
        legend: {
          position: 'bottom',
        },
        labels: [],
      },
      series: [],
      total: 0,
    }
  },
  methods: {
    async getPieDeploy() {
      try {
        const path = '/v1/deploymentPieChart'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        const data = await resp.json()
        this.colors = data.colors
        this.chartOptions = data.chartOptions
        this.series = data.series
        this.total = data.total
      } catch (error) {
        console.error(error)
      }
    },
  },
  cron: {
    time: 10000,
    method: 'getPieDeploy',
    autoStart: true,
  },
  beforeMount() {
    this.getPieDeploy()
  },
}
</script>
