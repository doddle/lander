<template>
  <div>
    <v-card flat>
      <v-card-title>Nodes</v-card-title>
      <v-card-subtitle>total: {{ total }}</v-card-subtitle>
      <apexcharts type="pie" :options="chartOptions" :series="series" />
    </v-card>
  </div>
</template>

<script>
import VueApexCharts from 'vue-apexcharts'

export default {
  name: 'OverviewPieNodes',
  components: {
    apexcharts: VueApexCharts
  },
  data: function() {
    return {
      chartOptions: {
        stroke: {
          width: 0
        },
        theme: {
          palette: 'pallet3'
        },
        colors: [],
        chart: {
          id: 'pie-nodes',
          dropShadow: {
            effect: false
          }
        },
        legend: {
          position: 'bottom'
        },
        labels: []
      },
      series: [],
      total: 0
    }
  },
  methods: {
    async getPieNodes() {
      try {
        const path = '/v1/pie/nodes'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        const data = await resp.json()
        this.chartOptions = data.chartOptions
        this.series = data.series
        this.total = data.total
      } catch (error) {
        console.error(error)
      }
    }
  },
  cron: {
    time: 10000,
    method: 'getPieNodes',
    autoStart: true
  },
  beforeMount() {
    this.getPieNodes()
  }
}
</script>
