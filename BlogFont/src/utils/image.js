/**
 * 图片公网前缀（文章等历史图：开发期也可走线上）。
 * 可用 .env 的 VITE_IMAGE_BASE 覆盖。
 */
export const IMAGE_BASE = (
  import.meta.env.VITE_IMAGE_BASE || "https://cspona.top/Pictures/"
).replace(/\/?$/, "/");

/**
 * 统一图片地址：文件名或 /Pictures/xxx → IMAGE_BASE + 文件名
 */
export function resolveImageUrl(path, fallback = "boli.jpg") {
  let s = (path ?? "").toString().trim();
  if (!s) s = fallback;

  if (/^https?:\/\//i.test(s)) {
    const idx = s.indexOf("/Pictures/");
    if (idx >= 0) {
      return IMAGE_BASE + s.slice(idx + "/Pictures/".length).replace(/^\/+/, "");
    }
    return s;
  }

  s = s.replace(/\\/g, "/").replace(/^\/+/, "");
  if (/^pictures\//i.test(s)) s = s.slice("pictures/".length);
  return IMAGE_BASE + s;
}

/**
 * 同源 /Pictures/（新上传的友链图等）：
 * - 本地：Vite 直接读 public/Pictures，或走后端静态
 * - 线上：nginx / 后端静态同一目录
 * 不要用线上 CDN 拼本地刚上传、尚未同步到公网的文件。
 */
export function resolveLocalPicturesUrl(path, fallback = "boli.jpg") {
  let s = (path ?? "").toString().trim();
  if (!s) s = fallback;
  if (/^https?:\/\//i.test(s)) return s;
  s = s.replace(/\\/g, "/").replace(/^\/+/, "");
  if (/^pictures\//i.test(s)) s = s.slice("pictures/".length);
  const encoded = s
    .split("/")
    .filter(Boolean)
    .map((p) => encodeURIComponent(p))
    .join("/");
  return `/Pictures/${encoded}`;
}

/**
 * 与后端 utils.BuildMarkdownTOC / heading 渲染共用的 slug 规则
 */
export function slugifyHeading(text) {
  const slug = text
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/\s+/g, "-")
    .replace(/[^\w\u4e00-\u9fa5-]/g, "")
    .replace(/-+/g, "-")
    .replace(/^-+|-+$/g, "")
    .toLowerCase();
  return slug || "section";
}

export function headingIdFromCounter(text, counter) {
  const slug = slugifyHeading(text);
  const n = (counter[slug] || 0) + 1;
  counter[slug] = n;
  return n > 1 ? `${slug}-${n}` : slug;
}
