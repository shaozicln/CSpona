<template>
  <div class="site-music">
    <audio
      ref="audioEl"
      :src="current?.url || ''"
      preload="metadata"
      playsinline
      @timeupdate="onTimeUpdate"
      @ended="onEnded"
      @play="onPlay"
      @pause="onPause"
      @canplay="onCanPlay"
    />

    <div class="sm-stack">
      <Transition name="sm-pop">
        <div v-if="listOpen" class="sm-list-card" role="dialog" aria-label="播放列表">
          <div class="sm-list-head">
            <h3>播放列表 · {{ playlist.length }}</h3>
            <button type="button" class="sm-list-close" @click="listOpen = false">×</button>
          </div>
          <p class="sm-list-tip">点击切歌，拖动手柄排序</p>
          <ul class="sm-list-ul">
            <li
              v-for="(track, i) in playlist"
              :key="track.url + '-' + i"
              class="sm-list-item"
              :class="{ active: i === index, dragging: dragFrom === i }"
              @dragover.prevent="onDragOver(i)"
              @drop.prevent="onDrop(i)"
              @click="selectTrack(i)"
            >
              <span
                class="sm-drag"
                title="拖拽排序"
                draggable="true"
                @click.stop
                @dragstart="onDragStart(i, $event)"
                @dragend="onDragEnd"
              >⋮⋮</span>
              <div class="sm-track-meta">
                <div class="sm-track-name">{{ track.name || "未命名" }}</div>
                <div class="sm-track-artist">{{ track.artist || " " }}</div>
              </div>
            </li>
          </ul>
          <p v-if="!playlist.length" class="sm-list-empty">暂无曲目</p>
        </div>
      </Transition>

      <div class="sm-row">
        <button type="button" class="sm-toggle" :title="expanded ? '收起' : '展开'" @click="expanded = !expanded">
          ♪
        </button>
        <div v-if="expanded" class="sm-panel">
          <div class="sm-title" :title="current?.name">{{ current?.name || "未播放" }}</div>
          <div class="sm-sub">{{ statusText }}</div>
          <div class="sm-controls">
            <button type="button" @click="prev" :disabled="!playlist.length">‹</button>
            <button type="button" class="sm-play" @click="toggle">
              {{ playing ? "暂停" : "播放" }}
            </button>
            <button type="button" @click="next" :disabled="!playlist.length">›</button>
            <button type="button" class="sm-list-btn" title="播放列表" @click="listOpen = !listOpen">
              列表
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useMusicStore } from "@/stores/music";
import { storeToRefs } from "pinia";

const music = useMusicStore();
const { playlist, volume, shouldPlay, status, errorMsg, prefs, autoplayTick } =
  storeToRefs(music);
const route = useRoute();

const audioEl = ref(null);
const index = ref(0);
const playing = ref(false);
const expanded = ref(true);
const listOpen = ref(false);
const userUnlocked = ref(false);
const lastPlaylistSig = ref("");
const autoplayBlocked = ref(false);
const dragFrom = ref(-1);
const suppressPlaylistReset = ref(false);
/** 仅预拉「下一首」元数据，避免整单抢带宽 */
let nextPrefetchEl = null;
let nextPrefetchUrl = "";
let nextPrefetchTimer = 0;
let lastProgressSyncAt = 0;

const current = computed(() => playlist.value[index.value] || null);
const statusText = computed(() => {
  if (!prefs.value.enabled) return "已关闭";
  if (status.value === "loading") return "加载中…";
  if (autoplayBlocked.value) return "浏览器拦截自动播放，请点播放";
  if (errorMsg.value && !playlist.value.length) return errorMsg.value;
  if (!playlist.value.length) return "暂无曲目";
  if (errorMsg.value) return errorMsg.value;
  return current.value?.artist || "";
});

onMounted(async () => {
  await music.loadSiteSettings();
  if (!isDetailPath(route.path)) {
    music.setRouteContext({ detail: false });
  }
  window.addEventListener("pointerdown", unlockOnce, { once: true, capture: true });
  await nextTick();
  if (prefs.value.autoplay) {
    tryAutoplay();
  }
});

onUnmounted(() => {
  window.removeEventListener("pointerdown", unlockOnce, { capture: true });
  clearNextPrefetch();
  music.flushProgress();
});

function clearNextPrefetch() {
  if (nextPrefetchTimer) {
    clearTimeout(nextPrefetchTimer);
    nextPrefetchTimer = 0;
  }
  if (nextPrefetchEl) {
    try {
      nextPrefetchEl.removeAttribute("src");
      nextPrefetchEl.load();
    } catch {
      /* ignore */
    }
    nextPrefetchEl = null;
  }
  nextPrefetchUrl = "";
}

/** 当前曲能播之后，空闲时只预拉下一首 metadata（不抢当前播放） */
function scheduleNextPrefetch() {
  if (nextPrefetchTimer) clearTimeout(nextPrefetchTimer);
  nextPrefetchTimer = window.setTimeout(() => {
    nextPrefetchTimer = 0;
    const list = playlist.value;
    if (!list.length || list.length < 2) return;
    const nextUrl = list[(index.value + 1) % list.length]?.url;
    if (!nextUrl || nextUrl === current.value?.url) return;
    if (nextUrl === nextPrefetchUrl && nextPrefetchEl) return;

    clearNextPrefetch();
    const a = new Audio();
    a.preload = "metadata";
    a.muted = true;
    a.src = nextUrl;
    try {
      a.load();
    } catch {
      return;
    }
    nextPrefetchEl = a;
    nextPrefetchUrl = nextUrl;
  }, 2500);
}

function unlockOnce() {
  userUnlocked.value = true;
  autoplayBlocked.value = false;
  if (prefs.value.autoplay && playlist.value.length && shouldPlay.value) {
    playCurrent();
  }
}

watch(
  () => route.fullPath,
  () => {
    music.flushProgress();
    if (!isDetailPath(route.path)) {
      music.setRouteContext({ detail: false, musicType: "", musicRef: "", articleId: "" });
    }
  }
);

function isDetailPath(path) {
  return /^\/article\//.test(path) || /^\/wanderland\//.test(path);
}

function playlistSignature(list) {
  return (list || []).map((t) => t.url).join("|");
}

function syncLiveProgress(force = false) {
  const now = Date.now();
  if (!force && now - lastProgressSyncAt < 1000) return;
  lastProgressSyncAt = now;
  const a = audioEl.value;
  music.setLiveProgress({
    index: index.value,
    currentTime: a?.currentTime || 0,
    playing: playing.value,
  });
}

function onTimeUpdate() {
  syncLiveProgress(false);
}
function onPlay() {
  playing.value = true;
  syncLiveProgress(true);
  scheduleNextPrefetch();
}
function onPause() {
  playing.value = false;
  syncLiveProgress(true);
}
function onCanPlay() {
  if (playing.value) scheduleNextPrefetch();
}

async function applyResumeHint(hint) {
  if (!hint) return false;
  suppressPlaylistReset.value = true;
  index.value = Math.max(0, Math.min(hint.index || 0, playlist.value.length - 1));
  await nextTick();
  lastPlaylistSig.value = playlistSignature(playlist.value);
  const seek =
    typeof hint.currentTime === "number" && Number.isFinite(hint.currentTime)
      ? Math.max(0, hint.currentTime)
      : 0;

  // 刷新前在播或开了自动播放 → 尝试续播；否则只定位进度（点播放接着听）
  const tryPlay = !!hint.wasPlaying || prefs.value.autoplay;

  if (tryPlay) {
    await playCurrent(seek);
  } else {
    const a = audioEl.value;
    if (a) {
      a.pause();
      const applySeek = () => {
        try {
          if (seek > 0) a.currentTime = seek;
        } catch {
          /* ignore */
        }
      };
      if (a.readyState >= 1) applySeek();
      else a.addEventListener("loadedmetadata", applySeek, { once: true });
    }
    playing.value = false;
  }
  await nextTick();
  syncLiveProgress(true);
  suppressPlaylistReset.value = false;
  return true;
}

watch(playlist, async (list, oldList) => {
  if (suppressPlaylistReset.value) return;
  clearNextPrefetch();

  const hint = music.takeResumeHint();
  if (hint) {
    await applyResumeHint(hint);
    return;
  }

  const sig = playlistSignature(list);
  const oldSig = playlistSignature(oldList || []);
  if (
    sig &&
    oldSig &&
    [...sig.split("|")].sort().join("|") === [...oldSig.split("|")].sort().join("|") &&
    sig !== oldSig
  ) {
    const curUrl = current.value?.url;
    if (curUrl) {
      const ni = list.findIndex((t) => t.url === curUrl);
      if (ni >= 0) index.value = ni;
    }
    lastPlaylistSig.value = sig;
    syncLiveProgress(true);
    return;
  }
  if (sig && sig === lastPlaylistSig.value) {
    if (prefs.value.autoplay) tryAutoplay();
    return;
  }
  lastPlaylistSig.value = sig;

  index.value = 0;
  const a = audioEl.value;
  const wasPlaying = playing.value;
  if (a) {
    a.pause();
    a.load();
  }
  playing.value = false;
  autoplayBlocked.value = false;

  await nextTick();
  const canAuto =
    list.length &&
    shouldPlay.value &&
    (prefs.value.autoplay || userUnlocked.value || wasPlaying);

  if (canAuto) {
    await tryAutoplay();
  }
  syncLiveProgress(true);
});

watch(autoplayTick, () => {
  if (prefs.value.autoplay) tryAutoplay();
});

watch(
  () => prefs.value.autoplay,
  (on) => {
    if (on) tryAutoplay();
  }
);

watch(volume, (v) => {
  if (audioEl.value) audioEl.value.volume = v;
});

watch(
  () => prefs.value.enabled,
  (on) => {
    if (!on && audioEl.value) {
      audioEl.value.pause();
      playing.value = false;
    }
  }
);

async function tryAutoplay() {
  if (!prefs.value.autoplay && !userUnlocked.value) return;
  if (!playlist.value.length || !shouldPlay.value) return;
  await playCurrent();
}

async function playCurrent(seekTo) {
  const a = audioEl.value;
  if (!a || !current.value?.url) return;
  a.volume = volume.value;
  const wantSeek =
    typeof seekTo === "number" && Number.isFinite(seekTo) && seekTo >= 0;
  try {
    if (wantSeek) {
      const apply = () => {
        try {
          a.currentTime = seekTo;
        } catch {
          /* ignore */
        }
      };
      if (a.readyState >= 1) apply();
      else a.addEventListener("loadedmetadata", apply, { once: true });
    }
    await a.play();
    playing.value = true;
    userUnlocked.value = true;
    autoplayBlocked.value = false;
    syncLiveProgress(true);
    scheduleNextPrefetch();
  } catch {
    playing.value = false;
    // 即使自动播放失败，也尽量落到续播位置
    if (wantSeek) {
      const apply = () => {
        try {
          a.currentTime = seekTo;
        } catch {
          /* ignore */
        }
      };
      if (a.readyState >= 1) apply();
      else a.addEventListener("loadedmetadata", apply, { once: true });
    }
    if (prefs.value.autoplay) {
      autoplayBlocked.value = true;
    }
    syncLiveProgress(true);
  }
}

function toggle() {
  const a = audioEl.value;
  if (!a || !current.value?.url) return;
  if (playing.value) {
    a.pause();
  } else {
    playCurrent();
  }
}

function next() {
  if (!playlist.value.length) return;
  index.value = (index.value + 1) % playlist.value.length;
  playCurrent(0);
}

function prev() {
  if (!playlist.value.length) return;
  index.value = (index.value - 1 + playlist.value.length) % playlist.value.length;
  playCurrent(0);
}

function onEnded() {
  if (playlist.value.length > 1) next();
  else playCurrent(0);
}

function selectTrack(i) {
  if (i < 0 || i >= playlist.value.length) return;
  const same = i === index.value;
  index.value = i;
  if (same && audioEl.value) {
    try {
      audioEl.value.currentTime = 0;
    } catch {
      /* ignore */
    }
  }
  playCurrent(0);
}

function onDragStart(i, e) {
  dragFrom.value = i;
  try {
    e.dataTransfer.setData("text/plain", String(i));
    e.dataTransfer.effectAllowed = "move";
  } catch {
    /* ignore */
  }
}

function onDragOver(i) {
  if (dragFrom.value < 0 || dragFrom.value === i) return;
}

function onDragEnd() {
  dragFrom.value = -1;
}

function onDrop(to) {
  const from = dragFrom.value;
  dragFrom.value = -1;
  if (from < 0 || from === to) return;
  const curUrl = current.value?.url;
  music.reorderPlaylist(from, to);
  if (curUrl) {
    const ni = playlist.value.findIndex((t) => t.url === curUrl);
    if (ni >= 0) index.value = ni;
  }
  syncLiveProgress(true);
  music.flushProgress();
}
</script>

<style scoped>
.site-music {
  position: fixed;
  right: 16px;
  bottom: 16px;
  z-index: 1500;
}
.sm-stack {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}
.sm-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}
.sm-toggle {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  border: 1px solid var(--panel-border, #ddd);
  background: var(--panel-bg, rgba(255, 255, 255, 0.92));
  color: var(--text-primary, #222);
  cursor: pointer;
  font-size: 18px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}
.sm-panel {
  min-width: 200px;
  max-width: 260px;
  padding: 10px 12px;
  border-radius: 12px;
  background: var(--panel-bg, rgba(255, 255, 255, 0.94));
  border: 1px solid var(--panel-border, #ddd);
  color: var(--text-primary, #222);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.14);
}
.sm-title {
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sm-sub {
  font-size: 11px;
  color: var(--text-muted, #888);
  margin: 2px 0 8px;
  min-height: 14px;
}
.sm-controls {
  display: flex;
  gap: 6px;
  align-items: center;
}
.sm-controls button {
  border: 1px solid var(--input-border, #ccc);
  background: var(--input-bg, #fff);
  color: var(--text-primary);
  border-radius: 6px;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 12px;
}
.sm-play {
  flex: 1;
}
.sm-list-btn {
  flex-shrink: 0;
}

.sm-list-card {
  width: min(280px, calc(100vw - 32px));
  max-height: min(42vh, 360px);
  background: var(--panel-bg, rgba(255, 255, 255, 0.96));
  border: 1px solid var(--panel-border, #ddd);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.16);
  color: var(--text-primary, #222);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.sm-pop-enter-active,
.sm-pop-leave-active {
  transition: opacity 0.18s ease, transform 0.2s ease;
}
.sm-pop-enter-from,
.sm-pop-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.96);
}
.sm-list-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px 4px;
}
.sm-list-head h3 {
  margin: 0;
  font-size: 14px;
}
.sm-list-close {
  border: none;
  background: transparent;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  color: var(--text-primary);
}
.sm-list-tip {
  margin: 0 12px 6px;
  font-size: 11px;
  color: var(--text-muted, #888);
}
.sm-list-ul {
  list-style: none;
  margin: 0;
  padding: 0 6px 10px;
  overflow-y: auto;
  flex: 1;
}
.sm-list-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
}
.sm-list-item:hover {
  background: var(--panel-bg-soft, rgba(0, 0, 0, 0.04));
}
.sm-list-item.active {
  background: var(--toc-active-bg, rgba(55, 105, 170, 0.12));
}
.sm-list-item.dragging {
  opacity: 0.55;
}
.sm-drag {
  color: var(--text-muted, #999);
  cursor: grab;
  user-select: none;
  padding: 4px 2px;
  letter-spacing: -2px;
  touch-action: none;
}
.sm-drag:active {
  cursor: grabbing;
}
.sm-track-meta {
  flex: 1;
  min-width: 0;
}
.sm-track-name {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sm-track-artist {
  font-size: 11px;
  color: var(--text-muted, #888);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sm-list-empty {
  text-align: center;
  color: var(--text-muted, #888);
  padding: 16px;
  font-size: 13px;
}
</style>
