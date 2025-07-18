import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import { createApp } from "vue";
import { createPinia } from 'pinia'
import App from "./App.vue";
import router from "./router";

import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const app = createApp(App).use(router).use(pinia).use(ElementPlus);

// 重构前的后端，仅有一个main.go，Back文件，已弃用
// app.config.globalProperties.URL = 'http://127.0.0.1:8081';

// 重构后的后端，支持图片文件上传，go文件分类管理，BlogBack
// app.config.globalProperties.URL = 'http://127.0.0.1:3000/api';

const isProduction = process.env.NODE_ENV === "production";
app.config.globalProperties.URL = isProduction
  ? "/api"
  : "http://127.0.0.1:3000/api";

app.config.globalProperties.URL2 = isProduction
  ? "https://cspona.top"
  : "http://127.0.0.1:3000";

// 图片
app.config.globalProperties.$imageBaseUrl =
  process.env.NODE_ENV === "production" ? "/Pictures/" : "../Public/Pictures/";

app.mount("#app");
