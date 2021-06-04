<template>
  <div class="example">
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
    apexcharts: VueApexCharts,
  },
  data: function() {
    return {
      chartOptions: {
        labels: ["one", "two", "three"],
        chart: {
          id: 'basic-bar',
        }
      },
      series: [30, 40, 45]
    };
  }
};
</script>
