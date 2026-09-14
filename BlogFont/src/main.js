import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import { createApp } from "vue";
import { createPinia } from 'pinia'
import App from "./App.vue";
import router from "./router";
import "./styles/theme-sync.css";

import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { IMAGE_BASE } from "@/utils/image.js";
import { ensureVisitorId } from "@/utils/visitor.js";

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

ensureVisitorId();

const app = createApp(App).use(router).use(pinia).use(ElementPlus);

// 同源 /api（开发走 Vite 代理，生产走 Nginx 反代）
app.config.globalProperties.URL = "/api";
app.config.globalProperties.URL2 = "";

// 图片一律走 IMAGE_BASE（开发默认 https://cspona.top/Pictures/）
app.config.globalProperties.$imageBaseUrl = IMAGE_BASE;

app.mount("#app");
