<template>
  <div>
    <v-card flat>
      <v-card-title>StatefulSets</v-card-title>
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
          id: 'pie-statefulsets',
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
    async getPieStatefulSets() {
      try {
        const path = '/v1/statefulSetPieChart'
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
    method: 'getPieStatefulSets',
    autoStart: true,
  },
  beforeMount() {
    this.getPieStatefulSets()
  },
}
</script>
