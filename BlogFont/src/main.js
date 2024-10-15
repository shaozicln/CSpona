// import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const app = createApp(App).use(router)

app.config.globalProperties.URL = 'http://127.0.0.1:8081';

app.mount('#app')
