import { IMAGE_BASE } from "@/utils/image.js";

/** 从分享链接或纯数字提取网易云 ID */
export function extractNeteaseId(input) {
  const s = (input ?? "").toString().trim();
  if (!s) return "";
  if (/^\d+$/.test(s)) return s;
  const m =
    s.match(/[?&#]id=(\d+)/i) ||
    s.match(/\/song\/(\d+)/i) ||
    s.match(/\/playlist\/(\d+)/i) ||
    s.match(/\/(\d+)\s*$/);
  return m ? m[1] : "";
}

/** 本地 music/xxx → 可播放 URL */
export function resolveMusicUrl(path) {
  if (!path) return "";
  let s = path.toString().trim();
  if (/^https?:\/\//i.test(s)) return s;
  s = s.replace(/\\/g, "/").replace(/^\/+/, "");
  if (/^pictures\//i.test(s)) s = s.slice("pictures/".length);
  return IMAGE_BASE + s;
}

const METING_CACHE_KEY = "cspona_music_meting_cache";
const METING_CACHE_TTL_MS = 6 * 60 * 60 * 1000; // 6h：列表元数据，加速二次进入

function readMetingCache(type, tid) {
  try {
    const all = JSON.parse(localStorage.getItem(METING_CACHE_KEY) || "{}");
    const hit = all[`${type}:${tid}`];
    if (!hit || !Array.isArray(hit.tracks)) return null;
    if (Date.now() - (hit.ts || 0) > METING_CACHE_TTL_MS) return null;
    return hit.tracks;
  } catch {
    return null;
  }
}

function writeMetingCache(type, tid, tracks) {
  try {
    const all = JSON.parse(localStorage.getItem(METING_CACHE_KEY) || "{}");
    all[`${type}:${tid}`] = { ts: Date.now(), tracks };
    // 控制体积：最多保留 20 个歌单缓存
    const keys = Object.keys(all);
    if (keys.length > 20) {
      keys
        .sort((a, b) => (all[a].ts || 0) - (all[b].ts || 0))
        .slice(0, keys.length - 20)
        .forEach((k) => delete all[k]);
    }
    localStorage.setItem(METING_CACHE_KEY, JSON.stringify(all));
  } catch {
    /* ignore quota */
  }
}

export async function fetchMetingTracks(apiBase, type, id) {
  const tid = extractNeteaseId(id);
  if (!tid) return [];

  const cached = readMetingCache(type, tid);
  if (cached?.length) {
    // 延后刷新，避免与首播抢带宽
    const defer =
      typeof window !== "undefined" && window.requestIdleCallback
        ? (fn) => window.requestIdleCallback(fn, { timeout: 20000 })
        : (fn) => setTimeout(fn, 8000);
    defer(() => {
      refreshMetingTracks(apiBase, type, tid).catch(() => {});
    });
    return cached.map((t) => ({ ...t }));
  }

  return refreshMetingTracks(apiBase, type, tid);
}

async function refreshMetingTracks(apiBase, type, tid) {
  const url = `${apiBase}/music/meting?type=${encodeURIComponent(type)}&id=${encodeURIComponent(tid)}`;
  const res = await fetch(url, { credentials: "include" });
  if (!res.ok) throw new Error("网易云解析失败");
  const data = await res.json();
  const list = Array.isArray(data) ? data : data ? [data] : [];
  const tracks = list
    .map((t) => ({
      name: t.name || t.title || "未知曲目",
      artist: t.artist || t.author || "",
      url: rewritePlayUrl(t.url || "", apiBase),
      cover: t.pic || t.cover || "",
    }))
    .filter((t) => t.url);
  if (tracks.length) writeMetingCache(type, tid, tracks);
  return tracks;
}

/** 把第三方 meting 的 type=url 改走本站代理，避免直链被拦 */
function rewritePlayUrl(playUrl, apiBase) {
  if (!playUrl) return "";
  try {
    const u = new URL(playUrl, typeof window !== "undefined" ? window.location.origin : "http://localhost");
    const typ = u.searchParams.get("type");
    const id = u.searchParams.get("id");
    if (typ === "url" && id && /^\d+$/.test(id)) {
      return `${apiBase}/music/meting?type=url&id=${encodeURIComponent(id)}`;
    }
  } catch {
    /* ignore */
  }
  return playUrl;
}
