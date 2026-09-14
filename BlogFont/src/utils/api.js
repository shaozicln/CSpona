import { getCurrentInstance } from "vue";
import router from "@/router";

function getApiBase() {
  const instance = getCurrentInstance();
  const fromApp = instance?.appContext?.config?.globalProperties?.URL;
  if (fromApp) return fromApp;
  return "/api";
}

/**
 * 带 Cookie 的 API 请求。浏览器不碰 token，只自动带 HttpOnly Cookie。
 */
export async function apiFetch(path, options = {}) {
  const base = getApiBase();
  const url = path.startsWith("http")
    ? path
    : `${base}${path.startsWith("/") ? path : `/${path}`}`;
  const headers = { ...(options.headers || {}) };
  if (options.body && !(options.body instanceof FormData) && !headers["Content-Type"]) {
    headers["Content-Type"] = "application/json";
  }
  return fetch(url, {
    ...options,
    headers,
    credentials: "include",
  });
}

export async function apiJson(path, options = {}) {
  const res = await apiFetch(path, options);
  let data = null;
  try {
    data = await res.json();
  } catch {
    data = null;
  }
  return { res, data };
}

/** 401 时弹窗，确认则跳转登录页。返回 true 表示已按未登录处理。 */
export function promptLoginIfUnauthorized(res, data) {
  if (!res || res.status !== 401) return false;
  const msg = data?.msg || data?.message || "请先登录后再操作";
  const go = window.confirm(`${msg}\n\n是否前往登录页？`);
  if (go) {
    const redirect = encodeURIComponent(
      window.location.pathname + window.location.search + window.location.hash
    );
    router.push({ path: "/login", query: { redirect } }).catch(() => {
      window.location.href = `/login?redirect=${redirect}`;
    });
  }
  return true;
}
