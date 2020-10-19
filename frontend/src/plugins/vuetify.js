import Vue from "vue";
import Vuetify from "vuetify/lib";

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: "#673ab7",
        secondary: "#03a9f4",
        accent: "#9c27b0",
        error: "#f44336",
        warning: "#ffc107",
        info:      "#cddc39",
        success:   "#8bc34a",
      }
    }
  }
});
