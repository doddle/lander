import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import crono from 'vue-crono'
import vuetify from './plugins/vuetify'
import titleMixin from './mixins/titleMixin'

Vue.mixin(titleMixin)

Vue.use(crono)
Vue.config.productionTip = false

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
