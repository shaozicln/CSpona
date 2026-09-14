import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";
import { apiJson } from "@/utils/api";
import {
  extractNeteaseId,
  fetchMetingTracks,
  resolveMusicUrl,
} from "@/utils/music.js";

const STORAGE_KEY = "cspona_music_prefs";
const ORDER_KEY = "cspona_music_order";
const PROGRESS_KEY = "cspona_music_progress";

function readPrefs() {
  try {
    const raw = JSON.parse(localStorage.getItem(STORAGE_KEY) || "{}");
    return {
      enabled: raw.enabled !== false,
      defaultPages: raw.defaultPages !== false,
      detailPages: raw.detailPages !== false,
      autoplay: raw.autoplay === true,
      volume:
        typeof raw.volume === "number"
          ? Math.min(1, Math.max(0, raw.volume))
          : 0.45,
      customEnabled: raw.customEnabled === true,
      customMode: raw.customMode === "song" ? "song" : "playlist",
      customId: typeof raw.customId === "string" ? raw.customId : "",
    };
  } catch {
    return {
      enabled: true,
      defaultPages: true,
      detailPages: true,
      autoplay: false,
      volume: 0.45,
      customEnabled: false,
      customMode: "playlist",
      customId: "",
    };
  }
}

function getApiBase() {
  return "/api";
}

function hasArticleMusic(art) {
  return (
    !!art &&
    (art.type === "netease" || art.type === "playlist" || art.type === "file") &&
    !!(art.ref || "").trim()
  );
}

function cloneTracks(list) {
  return (list || []).map((t) => ({ ...t }));
}

function readJson(key, fallback) {
  try {
    return JSON.parse(localStorage.getItem(key) || "null") ?? fallback;
  } catch {
    return fallback;
  }
}

function writeJson(key, value) {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch {
    /* quota */
  }
}

/** 按上次拖拽顺序重排；新增曲目追加末尾 */
function applySavedOrder(tracks, listKey) {
  if (!listKey || !tracks?.length) return tracks || [];
  const map = readJson(ORDER_KEY, {});
  const order = map[listKey];
  if (!Array.isArray(order) || !order.length) return tracks;

  const byUrl = new Map(tracks.map((t) => [t.url, t]));
  const sorted = [];
  for (const u of order) {
    if (byUrl.has(u)) {
      sorted.push(byUrl.get(u));
      byUrl.delete(u);
    }
  }
  for (const t of byUrl.values()) sorted.push(t);
  return sorted;
}

function saveOrder(listKey, urls) {
  if (!listKey || !urls?.length) return;
  const map = readJson(ORDER_KEY, {});
  map[listKey] = urls;
  const keys = Object.keys(map);
  if (keys.length > 40) {
    keys.slice(0, keys.length - 40).forEach((k) => delete map[k]);
  }
  writeJson(ORDER_KEY, map);
}

function readProgress(listKey) {
  if (!listKey) return null;
  const all = readJson(PROGRESS_KEY, {});
  const p = all[listKey];
  if (!p || typeof p !== "object") return null;
  // 超过 30 天丢弃
  if (Date.now() - (p.updatedAt || 0) > 30 * 24 * 60 * 60 * 1000) return null;
  return p;
}

function saveProgress(listKey, partial) {
  if (!listKey) return;
  const all = readJson(PROGRESS_KEY, {});
  all[listKey] = {
    ...(all[listKey] || {}),
    ...partial,
    listKey,
    updatedAt: Date.now(),
  };
  const keys = Object.keys(all);
  if (keys.length > 40) {
    keys
      .sort((a, b) => (all[a].updatedAt || 0) - (all[b].updatedAt || 0))
      .slice(0, keys.length - 40)
      .forEach((k) => delete all[k]);
  }
  writeJson(PROGRESS_KEY, all);
}

function progressToHint(p, tracks, preferPlay) {
  if (!p || !tracks?.length) return null;
  let index = 0;
  if (p.url) {
    const i = tracks.findIndex((t) => t.url === p.url);
    if (i >= 0) index = i;
    else if (typeof p.index === "number") {
      index = Math.min(Math.max(0, p.index), tracks.length - 1);
    }
  } else if (typeof p.index === "number") {
    index = Math.min(Math.max(0, p.index), tracks.length - 1);
  }
  const currentTime =
    typeof p.currentTime === "number" && Number.isFinite(p.currentTime)
      ? Math.max(0, p.currentTime)
      : 0;
  // 同歌单从头附近且第一首：无需续播
  if (index === 0 && currentTime < 0.8) return null;
  return {
    index,
    currentTime,
    // 刷新前在播 / 开了自动播放 → 尽量接着播
    wasPlaying: preferPlay || !!p.playing,
  };
}

export const useMusicStore = defineStore("music", () => {
  const prefs = ref(readPrefs());
  const settingsOpen = ref(false);
  const siteSettings = ref({ mode: "netease", playlistId: "", tracks: [] });
  const articleMusic = ref({ type: "", ref: "", articleId: "" });
  const isDetailRoute = ref(false);
  /** 当前列表来源：default | article */
  const activeSource = ref("default");
  /** 当前列表稳定键（用于顺序/进度持久化） */
  const listKey = ref("");
  /** 切入文章配乐前的快照（含进度） */
  const previousSnapshot = ref(null);
  /** 播放器实时进度，切入文章前写入快照 */
  const liveProgress = ref({ index: 0, currentTime: 0, playing: false });
  /** 还原后给播放器的续播指令 */
  const resumeHint = ref(null);
  const playlist = ref([]);
  const status = ref("idle");
  const errorMsg = ref("");
  const loadToken = ref(0);
  /** 递增以通知播放器尝试自动播放 */
  const autoplayTick = ref(0);

  let persistTimer = null;

  const volume = computed(() => prefs.value.volume);
  const shouldPlay = computed(() => {
    if (!prefs.value.enabled) return false;
    if (isDetailRoute.value) {
      if (!prefs.value.detailPages || !hasArticleMusic(articleMusic.value)) {
        return prefs.value.defaultPages;
      }
      return prefs.value.detailPages;
    }
    return prefs.value.defaultPages;
  });

  watch(
    prefs,
    (v) => {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(v));
    },
    { deep: true }
  );

  function openSettings() {
    settingsOpen.value = true;
  }
  function closeSettings() {
    settingsOpen.value = false;
  }
  function setPref(partial) {
    prefs.value = { ...prefs.value, ...partial };
  }

  function requestAutoplay() {
    autoplayTick.value += 1;
  }

  function defaultListKey() {
    const p = prefs.value;
    if (p.customEnabled) {
      const id = extractNeteaseId(p.customId) || (p.customId || "").trim();
      if (id) return `custom:${p.customMode}:${id}`;
    }
    const site = siteSettings.value;
    if (site.mode === "files") {
      const sig = (site.tracks || [])
        .map((t) => t.url || "")
        .filter(Boolean)
        .join("|");
      return `files:${sig || "empty"}`;
    }
    const id = extractNeteaseId(site.playlistId) || (site.playlistId || "").trim();
    if (site.mode === "song" && id) return `site:song:${id}`;
    if (id) return `site:playlist:${id}`;
    return "site:empty";
  }

  function articleListKey(art) {
    const id = (art.articleId || "").toString().trim();
    if (id) return `article:${id}`;
    return `article:${art.type}:${(art.ref || "").trim()}`;
  }

  function assignPlaylist(tracks, key, source, { resume = true } = {}) {
    // 换歌单前先把旧进度写入 localStorage
    if (listKey.value && playlist.value.length) {
      flushProgress();
    }
    const ordered = applySavedOrder(tracks, key);
    listKey.value = key;
    activeSource.value = source;
    if (resume) {
      const hint = progressToHint(
        readProgress(key),
        ordered,
        prefs.value.autoplay === true
      );
      if (hint) resumeHint.value = hint;
    }
    playlist.value = ordered;
  }

  /** 手动重排当前播放列表；返回新的当前曲下标 */
  function reorderPlaylist(fromIndex, toIndex) {
    const list = playlist.value.slice();
    if (
      fromIndex < 0 ||
      toIndex < 0 ||
      fromIndex >= list.length ||
      toIndex >= list.length ||
      fromIndex === toIndex
    ) {
      return -1;
    }
    const [item] = list.splice(fromIndex, 1);
    list.splice(toIndex, 0, item);
    playlist.value = list;
    if (listKey.value) {
      saveOrder(
        listKey.value,
        list.map((t) => t.url)
      );
    }
    return toIndex;
  }

  function moveTrack(fromIndex, delta) {
    return reorderPlaylist(fromIndex, fromIndex + delta);
  }

  function setLiveProgress(partial) {
    liveProgress.value = { ...liveProgress.value, ...partial };
    schedulePersistProgress();
  }

  function schedulePersistProgress() {
    if (!listKey.value) return;
    if (persistTimer) return;
    persistTimer = setTimeout(() => {
      persistTimer = null;
      flushProgress();
    }, 1000);
  }

  function flushProgress() {
    if (!listKey.value || !playlist.value.length) return;
    const idx = Math.min(
      Math.max(0, liveProgress.value.index || 0),
      playlist.value.length - 1
    );
    const track = playlist.value[idx];
    saveProgress(listKey.value, {
      index: idx,
      url: track?.url || "",
      currentTime: liveProgress.value.currentTime || 0,
      playing: !!liveProgress.value.playing,
    });
  }

  function takeResumeHint() {
    const h = resumeHint.value;
    resumeHint.value = null;
    return h;
  }

  function snapshotBeforeArticle() {
    if (!playlist.value.length) return;
    flushProgress();
    previousSnapshot.value = {
      tracks: cloneTracks(playlist.value),
      listKey: listKey.value,
      index: liveProgress.value.index || 0,
      currentTime: liveProgress.value.currentTime || 0,
      wasPlaying: !!liveProgress.value.playing,
      savedAt: Date.now(),
    };
  }

  function restorePreviousOrKeep() {
    const snap = previousSnapshot.value;
    if (snap?.tracks?.length) {
      previousSnapshot.value = null;
      status.value = "ready";
      errorMsg.value = "";
      const key = snap.listKey || defaultListKey();
      // 内存快照优先；若 localStorage 更新则合并
      const saved = readProgress(key);
      let index = Math.min(snap.index || 0, snap.tracks.length - 1);
      let currentTime = snap.currentTime || 0;
      let wasPlaying = !!snap.wasPlaying;
      if (saved && (saved.updatedAt || 0) >= (snap.savedAt || 0)) {
        const fromLs = progressToHint(saved, snap.tracks, prefs.value.autoplay);
        if (fromLs) {
          index = fromLs.index;
          currentTime = fromLs.currentTime;
          wasPlaying = fromLs.wasPlaying;
        }
      }
      resumeHint.value = { index, currentTime, wasPlaying };
      listKey.value = key;
      playlist.value = cloneTracks(snap.tracks);
      activeSource.value = "default";
      return true;
    }
    return false;
  }

  async function loadSiteSettings() {
    try {
      const { res, data } = await apiJson("/music/settings");
      if (res.ok && data?.data) {
        siteSettings.value = {
          mode: data.data.mode || "netease",
          playlistId: data.data.playlistId || "",
          tracks: Array.isArray(data.data.tracks) ? data.data.tracks : [],
        };
      }
    } catch (e) {
      console.warn("load music settings failed", e);
    }
  }

  function setRouteContext({
    detail,
    musicType = "",
    musicRef = "",
    articleId = "",
  }) {
    flushProgress();
    isDetailRoute.value = !!detail;
    articleMusic.value = {
      type: (musicType || "").trim(),
      ref: (musicRef || "").trim(),
      articleId: (articleId ?? "").toString().trim(),
    };
  }

  function clearArticleMusic() {
    articleMusic.value = { type: "", ref: "", articleId: "" };
  }

  async function buildDefaultTracks() {
    const p = prefs.value;
    if (p.customEnabled) {
      const id = extractNeteaseId(p.customId) || (p.customId || "").trim();
      if (id) {
        const typ = p.customMode === "song" ? "song" : "playlist";
        return fetchMetingTracks(getApiBase(), typ, id);
      }
    }

    const site = siteSettings.value;
    if (site.mode === "files" && site.tracks?.length) {
      return site.tracks
        .map((t) => ({
          name: t.title || "曲目",
          artist: "",
          url: resolveMusicUrl(t.url),
          cover: "",
        }))
        .filter((t) => t.url);
    }
    if (site.mode === "song" && site.playlistId) {
      return fetchMetingTracks(
        getApiBase(),
        "song",
        extractNeteaseId(site.playlistId)
      );
    }
    if (site.playlistId) {
      return fetchMetingTracks(
        getApiBase(),
        "playlist",
        extractNeteaseId(site.playlistId)
      );
    }
    return [];
  }

  async function buildArticleTracks(art) {
    if (art.type === "file") {
      const url = resolveMusicUrl(art.ref);
      return url
        ? [{ name: "文章配乐", artist: "", url, cover: "" }]
        : [];
    }
    const id = extractNeteaseId(art.ref);
    if (!id) return [];
    const typ = art.type === "playlist" ? "playlist" : "song";
    return fetchMetingTracks(getApiBase(), typ, id);
  }

  /**
   * - 关「文章/漫游详情」音乐：永不加载文章配乐
   * - 有配乐且开关开：切文章曲（先快照当前列表）
   * - 离开有配乐文章：还原快照（继续之前的音乐），没有快照再拉默认
   * - 无配乐详情：不打断
   */
  async function rebuildPlaylist() {
    const token = ++loadToken.value;
    const art = articleMusic.value;
    const allowArticle =
      prefs.value.detailPages &&
      isDetailRoute.value &&
      hasArticleMusic(art);

    if (!prefs.value.enabled) {
      flushProgress();
      status.value = "idle";
      return;
    }

    // 正在播文章曲，但现在不允许文章配乐（关开关 / 离开文章 / 无配乐）→ 恢复之前
    if (activeSource.value === "article" && !allowArticle) {
      flushProgress();
      if (restorePreviousOrKeep()) {
        return;
      }
      // 无快照则尝试默认
      if (!prefs.value.defaultPages) {
        status.value = "idle";
        return;
      }
      status.value = "loading";
      try {
        const next = await buildDefaultTracks();
        if (token !== loadToken.value) return;
        const cleaned = (next || []).filter((t) => t.url);
        assignPlaylist(cleaned, defaultListKey(), "default");
        status.value = cleaned.length ? "ready" : "idle";
        if (!cleaned.length) errorMsg.value = "尚未配置默认可播曲目";
        if (prefs.value.autoplay && cleaned.length) requestAutoplay();
      } catch (e) {
        if (token !== loadToken.value) return;
        console.error(e);
        status.value = "error";
        errorMsg.value = "音乐加载失败";
      }
      return;
    }

    // 详情无配乐或关了文章音乐：保持当前，不加载文章曲
    if (isDetailRoute.value && !allowArticle) {
      return;
    }

    // 非详情：确保默认歌单
    if (!isDetailRoute.value) {
      if (!prefs.value.defaultPages) {
        status.value = "idle";
        return;
      }
      const key = defaultListKey();
      if (
        activeSource.value === "default" &&
        playlist.value.length &&
        listKey.value === key
      ) {
        return; // 同一歌单已在播，不重复拉
      }
      status.value = "loading";
      errorMsg.value = "";
      try {
        const next = await buildDefaultTracks();
        if (token !== loadToken.value) return;
        const cleaned = (next || []).filter((t) => t.url);
        if (cleaned.length) {
          assignPlaylist(cleaned, key, "default");
          status.value = "ready";
          if (prefs.value.autoplay) requestAutoplay();
        } else {
          status.value = playlist.value.length ? "ready" : "idle";
          errorMsg.value = "尚未配置默认可播曲目";
        }
      } catch (e) {
        if (token !== loadToken.value) return;
        console.error(e);
        status.value = playlist.value.length ? "ready" : "error";
        errorMsg.value = "音乐加载失败";
      }
      return;
    }

    // 详情 + 允许文章配乐
    const nextKey = articleListKey(art);
    if (
      activeSource.value === "article" &&
      listKey.value === nextKey &&
      playlist.value.length
    ) {
      return; // 同一篇文章配乐，保持当前进度
    }

    status.value = "loading";
    errorMsg.value = "";
    try {
      // 从默认歌单切入文章：快照；文章→文章只落盘旧进度
      if (activeSource.value !== "article" && playlist.value.length) {
        snapshotBeforeArticle();
      } else {
        flushProgress();
      }
      const next = await buildArticleTracks(art);
      if (token !== loadToken.value) return;
      const cleaned = (next || []).filter((t) => t.url);
      if (cleaned.length) {
        assignPlaylist(cleaned, nextKey, "article");
        status.value = "ready";
        errorMsg.value = "";
        if (prefs.value.autoplay) requestAutoplay();
      } else {
        // 文章曲失败：不打断，必要时还原
        if (activeSource.value !== "article") {
          restorePreviousOrKeep();
        }
        status.value = playlist.value.length ? "ready" : "idle";
        errorMsg.value = "当前配乐暂不可播";
      }
    } catch (e) {
      if (token !== loadToken.value) return;
      console.error(e);
      if (activeSource.value !== "article") {
        restorePreviousOrKeep();
      }
      status.value = playlist.value.length ? "ready" : "error";
      errorMsg.value = "当前配乐暂不可播";
    }
  }

  watch(
    [
      isDetailRoute,
      articleMusic,
      siteSettings,
      () => prefs.value.enabled,
      () => prefs.value.defaultPages,
      () => prefs.value.detailPages,
      () => prefs.value.customEnabled,
      () => prefs.value.customMode,
      () => prefs.value.customId,
    ],
    () => {
      rebuildPlaylist();
    },
    { deep: true }
  );

  watch(
    () => prefs.value.autoplay,
    (on) => {
      if (on) requestAutoplay();
    }
  );

  if (typeof window !== "undefined") {
    window.addEventListener("pagehide", () => flushProgress());
    document.addEventListener("visibilitychange", () => {
      if (document.visibilityState === "hidden") flushProgress();
    });
  }

  return {
    prefs,
    settingsOpen,
    siteSettings,
    articleMusic,
    isDetailRoute,
    activeSource,
    listKey,
    playlist,
    status,
    errorMsg,
    volume,
    shouldPlay,
    autoplayTick,
    resumeHint,
    openSettings,
    closeSettings,
    setPref,
    setLiveProgress,
    flushProgress,
    takeResumeHint,
    loadSiteSettings,
    setRouteContext,
    clearArticleMusic,
    rebuildPlaylist,
    requestAutoplay,
    reorderPlaylist,
    moveTrack,
  };
});
