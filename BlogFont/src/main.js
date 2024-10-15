// import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

const app = createApp(App).use(router)

// 重构前的后端，仅有一个main.go，Back文件，已弃用
// app.config.globalProperties.URL = 'http://127.0.0.1:8081';

// 重构后的后端，文件分类管理，BlogBack
app.config.globalProperties.URL = 'http://127.0.0.1:3000/api';

app.mount('#app')
