<script>
import { Pie } from "vue-chartjs";

export default {
  extends: Pie,

  data() {
    return {
      chartData: {
        labels: ["Error", "Running", "Unknown"],
        datasets: [
          {
            borderWidth: 0,
            borderColor: [
              "rgba(255,99,132)",
              "rgb(65,194,93,1)",
              "rgb(235,190,54)"
            ],
            backgroundColor: [
              "rgba(255, 99, 132, 0.2)",
              "rgba(86, 255, 123, 0.2)",
              "rgba(235,190,54,0.2)"
            ],
            data: [5, 343, 0]
          }
        ]
      },
      chartOptions: {
        legend: {
          display: true
        },
        responsive: false,
        maintainAspectRatio: false
      }
    };
  },
  mounted() {
    this.renderChart(this.chartData, this.options);
  }
};
</script>
