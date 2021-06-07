<template>
  <div>
    <apexcharts
      width="500"
      height="350"
      type="pie"
      :options="chartOptions"
      :series="series"
    />
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
          id: "deployments"
        }
      },
      series: []
    };
  },
  methods: {
    async getPieDeploy() {
      try {
        const resp = await fetch("/v1/pie/deployments");
        const data = await resp.json();
        console.log("retrieving v1/pie/deployments");
        this.series = data;
      } catch (error) {
        console.error(error);
      }
    }
  },
  cron: {
    time: 10000,
    method: "getPieDeploy",
    autoStart: true
  },
  beforeMount() {
    this.getPieDeploy();
  }
};
</script>
