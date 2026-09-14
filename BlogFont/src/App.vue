<script setup>
import MyHead from './components/MyHead.vue'
import SiteMusicPlayer from './components/common/SiteMusicPlayer.vue'
import MusicSettingsModal from './components/common/MusicSettingsModal.vue'
import { onMounted, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { resolveImageUrl } from '@/utils/image.js'

const userStore = useUserStore()
const themeStore = useThemeStore()
const bgReady = ref(false)
const bgSrc = resolveImageUrl('background.jpg')

function markBgReady() {
  bgReady.value = true
}

onMounted(() => {
  themeStore.setMode(themeStore.mode)
  userStore.restoreSession()
  const img = document.querySelector('#background .bg-img')
  if (img && img.complete && img.naturalWidth > 0) {
    markBgReady()
  }
})
</script>

<template>
  <!-- z-index:0，不要用 -1：body 不透明底色会把 -1 背景完全盖住 -->
  <div id="background" aria-hidden="true" :class="{ 'is-ready': bgReady }">
    <img
      class="bg-img"
      :src="bgSrc"
      alt=""
      decoding="async"
      fetchpriority="high"
      @load="markBgReady"
      @error="markBgReady"
    />
    <div class="bg-mask"></div>
  </div>
  <div id="app-root">
    <MyHead v-if="!$route.meta.hideMyHead"/>
    <router-view />
    <SiteMusicPlayer />
    <MusicSettingsModal />
  </div>
</template>

<style>
:root {
  --panel-bg: rgba(255, 255, 255, 0.8);
  --panel-bg-soft: rgba(240, 248, 255, 0.7);
  /* 浮层（nav 下拉 / 音乐卡片与设置）：与白天同透明度 */
  --float-panel-bg: rgba(255, 255, 255, 0.8);
  --float-panel-bg-soft: rgba(240, 248, 255, 0.7);
  --panel-border: #ffffff;
  --text-primary: rgba(0, 0, 0, 0.88);
  --text-secondary: #666;
  --text-muted: rgba(0, 0, 0, 0.55);
  --code-bg: #ffffff;
  --code-border: #e5e5e5;
  --code-text: #24292f;
  --toc-active-bg: rgba(55, 105, 170, 0.12);
  --toc-active-bar: rgba(55, 105, 170, 0.9);
  --toc-active-text: rgba(30, 70, 120, 0.95);
  --nav-text: #000000;
  --bg-mask: transparent;
  --copy-icon: rgba(0, 0, 0, 0.55);
  --copy-icon-hover: rgba(0, 0, 0, 0.9);
  --input-bg: rgba(255, 255, 255, 0.92);
  --input-text: #222;
  --input-border: #ddd;
  --input-placeholder: #999;
  --hover-bg: rgba(0, 0, 0, 0.06);
  --card-bg: rgba(255, 255, 255, 0.85);
  --btn-bg: rgb(139, 189, 234);
  --btn-bg-hover: rgb(119, 172, 218);
  --btn-text: #ffffff;
  --page-bg-fallback: #d7e3ef;
}

html[data-theme="dark"] {
  --panel-bg: rgba(0, 0, 0, 0.38);
  --panel-bg-soft: rgba(0, 0, 0, 0.26);
  /* 与白天 0.8 / 0.7 对齐，浮层更实、更易读 */
  --float-panel-bg: rgba(0, 0, 0, 0.8);
  --float-panel-bg-soft: rgba(0, 0, 0, 0.7);
  --panel-border: rgba(255, 255, 255, 0.16);
  --text-primary: rgba(255, 255, 255, 0.92);
  --text-secondary: rgba(255, 255, 255, 0.62);
  --text-muted: rgba(255, 255, 255, 0.5);
  --code-bg: rgba(0, 0, 0, 0.32);
  --code-border: rgba(255, 255, 255, 0.16);
  --code-text: #e6edf3;
  --toc-active-bg: rgba(120, 170, 230, 0.18);
  --toc-active-bar: rgba(140, 190, 255, 0.95);
  --toc-active-text: rgba(210, 230, 255, 0.98);
  --nav-text: rgba(255, 255, 255, 0.92);
  --bg-mask: rgba(40, 40, 40, 0.26);
  --copy-icon: rgba(255, 255, 255, 0.65);
  --copy-icon-hover: rgba(255, 255, 255, 0.95);
  --input-bg: rgba(0, 0, 0, 0.4);
  --input-text: rgba(255, 255, 255, 0.92);
  --input-border: rgba(255, 255, 255, 0.2);
  --input-placeholder: rgba(255, 255, 255, 0.45);
  --hover-bg: rgba(255, 255, 255, 0.1);
  --card-bg: rgba(0, 0, 0, 0.34);
  --btn-bg: rgb(120, 190, 255);
  --btn-bg-hover: rgb(140, 200, 255);
  --btn-text: #ffffff;
  --page-bg-fallback: #2a2f36;
}

* {
  padding: 0px;
  margin: 0px;
  font-family: Cormorant SC, serif;
}

html,
body {
  min-height: 100%;
  /* 透明：否则会盖住 fixed 背景层 */
  background: transparent;
}

#app {
  position: relative;
  min-height: 100vh;
}

#background {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
  background: var(--page-bg-fallback);
}

#background .bg-img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  display: block;
  transition: opacity 0.35s ease;
}

#background.is-ready .bg-img {
  opacity: 0.8;
}

#background .bg-mask {
  position: absolute;
  inset: 0;
  background: var(--bg-mask);
  transition: background 0.25s ease;
}

#app-root {
  position: relative;
  z-index: 1;
  color: var(--text-primary);
  min-height: 100vh;
}

::-webkit-scrollbar {
  width: 2px;
  height: 2px;
}
</style>
