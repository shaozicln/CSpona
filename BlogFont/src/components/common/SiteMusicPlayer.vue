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
      @error="onAudioError"
    />

    <div class="sm-stack">
      <Transition name="sm-pop">
        <div v-if="listOpen" class="sm-list-card" role="dialog" aria-label="播放列表">
          <div class="sm-list-head">
            <h3>播放列表 · {{ playlist.length }}</h3>
            <div class="sm-list-actions">
              <button
                type="button"
                class="sm-mode-btn"
                :title="playModeTitle"
                @click="cyclePlayMode"
              >
                {{ playModeLabel }}
              </button>
              <button type="button" class="sm-list-close" @click="listOpen = false">×</button>
            </div>
          </div>
          <p class="sm-list-tip">{{ playModeTip }} · 点击切歌，拖动手柄排序</p>
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
const trackUnavailable = ref(false);
const dragFrom = ref(-1);
const suppressPlaylistReset = ref(false);
/** 仅预拉「下一首」元数据，避免整单抢带宽 */
let nextPrefetchEl = null;
let nextPrefetchUrl = "";
let nextPrefetchTimer = 0;
let lastProgressSyncAt = 0;
/** 连续跳过不可播（VIP/失效）曲目，防止死循环 */
let skipFailCount = 0;
let playGen = 0;
let loadWatchTimer = 0;

const current = computed(() => playlist.value[index.value] || null);
const playMode = computed(() => prefs.value.playMode || "loop");
const playModeLabel = computed(() => {
  if (playMode.value === "one") return "单曲";
  if (playMode.value === "shuffle") return "随机";
  return "循环";
});
const playModeTitle = computed(() => {
  if (playMode.value === "one") return "单曲循环（点击切换）";
  if (playMode.value === "shuffle") return "随机播放（点击切换）";
  return "列表循环（点击切换）";
});
const playModeTip = computed(() => {
  if (playMode.value === "one") return "单曲循环";
  if (playMode.value === "shuffle") return "随机播放";
  return "列表循环";
});
const statusText = computed(() => {
  if (!prefs.value.enabled) return "已关闭";
  if (status.value === "loading") return "加载中…";
  if (trackUnavailable.value) return "当前曲暂不可播，已跳过…";
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
  // 任意交互解锁后，若开了自动播放则补播
  window.addEventListener("pointerdown", unlockOnce, { once: true, capture: true });
  window.addEventListener("keydown", unlockOnce, { once: true, capture: true });
  await nextTick();
  if (prefs.value.autoplay) {
    tryAutoplay();
  }
});

onUnmounted(() => {
  window.removeEventListener("pointerdown", unlockOnce, { capture: true });
  window.removeEventListener("keydown", unlockOnce, { capture: true });
  if (loadWatchTimer) clearTimeout(loadWatchTimer);
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
  const a = audioEl.value;
  if (a) {
    try {
      a.muted = false;
      a.volume = volume.value;
    } catch {
      /* ignore */
    }
  }
  if (prefs.value.autoplay && playlist.value.length && shouldPlay.value) {
    if (!playing.value || a?.paused) {
      playCurrent(undefined, { forceSound: true });
    }
  }
}

function isAutoplayPolicyError(err) {
  const name = err?.name || "";
  const msg = String(err?.message || "");
  return (
    name === "NotAllowedError" ||
    /user didn't interact|autoplay|not allowed/i.test(msg)
  );
}

function clearLoadWatch() {
  if (loadWatchTimer) {
    clearTimeout(loadWatchTimer);
    loadWatchTimer = 0;
  }
}

/** 开播后尽量取消静音（延迟多拍，提高首次进页有声成功率） */
function forceUnmute(a) {
  if (!a) return;
  const apply = () => {
    try {
      a.muted = false;
      a.defaultMuted = false;
      a.removeAttribute("muted");
      a.volume = volume.value;
    } catch {
      /* ignore */
    }
  };
  apply();
  setTimeout(apply, 0);
  setTimeout(apply, 120);
  setTimeout(apply, 400);
}

/** 不可播（VIP/失效/加载失败）→ 自动下一首 */
async function skipUnplayable() {
  const n = playlist.value.length;
  if (n <= 0) return;
  skipFailCount += 1;
  trackUnavailable.value = true;
  autoplayBlocked.value = false;
  if (skipFailCount >= n) {
    skipFailCount = 0;
    playing.value = false;
    trackUnavailable.value = false;
    // 整圈都不可播
    syncLiveProgress(true);
    return;
  }
  index.value = (index.value + 1) % n;
  // 跳过坏曲时不受单曲循环限制
  await nextTick();
  await playCurrent(0, {
    fromSkip: true,
    preferMutedAutoplay: prefs.value.autoplay && !userUnlocked.value,
  });
}

function onAudioError() {
  // 空 src 的 error 忽略
  if (!current.value?.url) return;
  skipUnplayable();
}

function onCanPlay() {
  clearLoadWatch();
  trackUnavailable.value = false;
  if (playing.value) scheduleNextPrefetch();
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
  autoplayBlocked.value = false;
  // 自动播放若以静音启动，这里立刻尝试开声
  if (prefs.value.autoplay || userUnlocked.value) {
    forceUnmute(audioEl.value);
  }
  syncLiveProgress(true);
  scheduleNextPrefetch();
}
function onPause() {
  playing.value = false;
  syncLiveProgress(true);
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
  // 开了自动播放：优先静音开播再开声（浏览器对静音 autoplay 几乎总是放行）
  await playCurrent(undefined, { preferMutedAutoplay: prefs.value.autoplay && !userUnlocked.value });
}

/**
 * @param {number} [seekTo]
 * @param {{ fromSkip?: boolean, forceSound?: boolean, preferMutedAutoplay?: boolean }} [opts]
 */
async function playCurrent(seekTo, opts = {}) {
  const a = audioEl.value;
  if (!a || !current.value?.url) return;
  const gen = ++playGen;
  a.volume = volume.value;
  const wantSeek =
    typeof seekTo === "number" && Number.isFinite(seekTo) && seekTo >= 0;
  if (!opts.fromSkip) {
    skipFailCount = 0;
    trackUnavailable.value = false;
  }

  const startMuted =
    !!opts.preferMutedAutoplay &&
    !opts.forceSound &&
    !userUnlocked.value &&
    !!prefs.value.autoplay;

  // 关键属性再 play：静音自动播放策略才稳定
  if (startMuted) {
    a.defaultMuted = true;
    a.muted = true;
    a.setAttribute("muted", "");
  } else {
    a.defaultMuted = false;
    a.muted = false;
    a.removeAttribute("muted");
  }

  clearLoadWatch();
  // 仅在真正加载失败时跳过；已在播放中不要因缓冲慢误跳
  loadWatchTimer = window.setTimeout(() => {
    loadWatchTimer = 0;
    if (gen !== playGen) return;
    if (playing.value || !a.paused) return;
    if (a.error || a.readyState === 0) skipUnplayable();
  }, 12000);

  const applySeek = () => {
    try {
      if (wantSeek) a.currentTime = seekTo;
    } catch {
      /* ignore */
    }
  };

  const markPlaying = () => {
    if (gen !== playGen) return;
    playing.value = true;
    autoplayBlocked.value = false;
    trackUnavailable.value = false;
    skipFailCount = 0;
    clearLoadWatch();
    forceUnmute(a);
    if (!a.muted) userUnlocked.value = true;
    syncLiveProgress(true);
    scheduleNextPrefetch();
  };

  try {
    if (wantSeek) {
      if (a.readyState >= 1) applySeek();
      else a.addEventListener("loadedmetadata", applySeek, { once: true });
    }
    // 等一帧，确保 muted 属性已生效
    await nextTick();
    await a.play();
    markPlaying();
  } catch (err) {
    if (gen !== playGen) return;

    if (!isAutoplayPolicyError(err)) {
      playing.value = false;
      clearLoadWatch();
      await skipUnplayable();
      return;
    }

    // 有声被拦 → 强制静音再开，再开声
    if (prefs.value.autoplay || opts.preferMutedAutoplay) {
      try {
        a.defaultMuted = true;
        a.muted = true;
        a.setAttribute("muted", "");
        await nextTick();
        await a.play();
        markPlaying();
        return;
      } catch {
        try {
          a.muted = false;
          a.removeAttribute("muted");
        } catch {
          /* ignore */
        }
      }
    }

    playing.value = false;
    if (wantSeek) {
      if (a.readyState >= 1) applySeek();
      else a.addEventListener("loadedmetadata", applySeek, { once: true });
    }
    if (prefs.value.autoplay) {
      autoplayBlocked.value = true;
    }
    clearLoadWatch();
    syncLiveProgress(true);
  }
}

function toggle() {
  const a = audioEl.value;
  if (!a || !current.value?.url) return;
  if (playing.value) {
    a.pause();
  } else {
    autoplayBlocked.value = false;
    userUnlocked.value = true;
    playCurrent(undefined, { forceSound: true });
  }
}

function cyclePlayMode() {
  const order = ["loop", "one", "shuffle"];
  const cur = playMode.value;
  const nextMode = order[(order.indexOf(cur) + 1) % order.length];
  music.setPref({ playMode: nextMode });
}

/** 按播放模式取下一曲下标；dir: 1 下一首 / -1 上一首 */
function pickAdjacentIndex(fromIndex, dir = 1) {
  const n = playlist.value.length;
  if (n <= 0) return 0;
  if (n === 1) return 0;
  const mode = playMode.value;
  if (mode === "shuffle") {
    let r = Math.floor(Math.random() * n);
    let guard = 0;
    while (r === fromIndex && guard++ < 12) {
      r = Math.floor(Math.random() * n);
    }
    return r;
  }
  return (fromIndex + dir + n) % n;
}

function next() {
  if (!playlist.value.length) return;
  skipFailCount = 0;
  trackUnavailable.value = false;
  index.value = pickAdjacentIndex(index.value, 1);
  playCurrent(0);
}

function prev() {
  if (!playlist.value.length) return;
  skipFailCount = 0;
  trackUnavailable.value = false;
  index.value = pickAdjacentIndex(index.value, -1);
  playCurrent(0);
}

function onEnded() {
  skipFailCount = 0;
  trackUnavailable.value = false;
  if (playMode.value === "one") {
    playCurrent(0);
    return;
  }
  if (playlist.value.length > 1) {
    index.value = pickAdjacentIndex(index.value, 1);
    playCurrent(0);
  } else {
    playCurrent(0);
  }
}

function selectTrack(i) {
  if (i < 0 || i >= playlist.value.length) return;
  skipFailCount = 0;
  trackUnavailable.value = false;
  autoplayBlocked.value = false;
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
  background: var(--float-panel-bg, var(--panel-bg, rgba(255, 255, 255, 0.92)));
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
  background: var(--float-panel-bg, var(--panel-bg, rgba(255, 255, 255, 0.94)));
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
  background: var(--float-panel-bg, var(--panel-bg, rgba(255, 255, 255, 0.96)));
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
  gap: 8px;
}
.sm-list-head h3 {
  margin: 0;
  font-size: 14px;
  min-width: 0;
  flex: 1;
}
.sm-list-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}
.sm-mode-btn {
  border: 1px solid var(--input-border, #ccc);
  background: var(--input-bg, #fff);
  color: var(--text-primary);
  border-radius: 6px;
  padding: 2px 8px;
  font-size: 12px;
  cursor: pointer;
  line-height: 1.4;
}
.sm-mode-btn:hover {
  background: var(--hover-bg, rgba(0, 0, 0, 0.06));
}
.sm-list-close {
  border: none;
  background: transparent;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  color: var(--text-primary);
  padding: 0 2px;
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
  background: var(--float-panel-bg-soft, var(--panel-bg-soft, rgba(0, 0, 0, 0.04)));
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
