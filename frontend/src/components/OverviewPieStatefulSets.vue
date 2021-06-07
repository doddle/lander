<template>
  <div>
    <v-card flat="true">
      <v-card-title>StatefulSets</v-card-title>
      <v-card-subtitle>total: {{ series[0] + series[1] }}</v-card-subtitle>
      <apexcharts type="pie" :options="chartOptions" :series="series" />
    </v-card>
  </div>
</template>

<script>
import VueApexCharts from "vue-apexcharts";

export default {
  name: "Apex",
  components: {
    apexcharts: VueApexCharts
  },
  data: function() {
    return {
      chartOptions: {
        labels: ["bad", "good"],
        chart: {
          id: "statefulsets"
        }
      },
      series: []
    };
  },
  methods: {
    async getPieStatefulSets() {
      try {
        const resp = await fetch("/v1/pie/statefulsets");
        const data = await resp.json();
        console.log("retrieving v1/pie/statefulsets");
        this.series = data;
      } catch (error) {
        console.error(error);
      }
    }
  },
  cron: {
    time: 10000,
    method: "getPieStatefulSets",
    autoStart: true
  },
  beforeMount() {
    this.getPieStatefulSets();
  }
};
</script>
