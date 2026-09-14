import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";

const STORAGE_KEY = "cspona_theme";

function readStored() {
  const v = localStorage.getItem(STORAGE_KEY);
  return v === "dark" ? "dark" : "light";
}

function applyDomTheme(mode) {
  document.documentElement.setAttribute("data-theme", mode);
}

export const useThemeStore = defineStore("theme", () => {
  const mode = ref(readStored());
  applyDomTheme(mode.value);

  const isDark = computed(() => mode.value === "dark");
  /** 下拉显示「另一个」模式的名字 */
  const switchLabel = computed(() =>
    mode.value === "dark" ? "日光模式" : "夜色模式"
  );

  function setMode(next) {
    mode.value = next === "dark" ? "dark" : "light";
  }

  function toggle() {
    setMode(mode.value === "dark" ? "light" : "dark");
  }

  watch(
    mode,
    (v) => {
      localStorage.setItem(STORAGE_KEY, v);
      applyDomTheme(v);
    },
    { immediate: true }
  );

  return { mode, isDark, switchLabel, setMode, toggle };
});
