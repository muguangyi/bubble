import Vue from 'vue'
import App from './app.vue'
import BootstrapVue from 'bootstrap-vue'
import axios from "axios";
import 'bootswatch/dist/lumen/bootstrap.css'

Vue.use(BootstrapVue)
Vue.prototype.$axios = axios
Vue.prototype.HOME = '/api/v1/'
Vue.config.productionTip = false

new Vue({
    render: h => h(App),
}).$mount('#app')
