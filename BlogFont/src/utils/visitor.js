const VISITOR_KEY = "cspona_vid";

function randomId() {
  if (typeof crypto !== "undefined" && crypto.randomUUID) {
    return crypto.randomUUID();
  }
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

/** 匿名访客稳定 ID（localStorage，约等于一人一号） */
export function ensureVisitorId() {
  try {
    let id = localStorage.getItem(VISITOR_KEY);
    if (!id || !/^[a-f0-9-]{8,64}$/i.test(id)) {
      id = randomId();
      localStorage.setItem(VISITOR_KEY, id);
    }
    // 同步到 cookie，便于后端无 Header 时也能识别
    document.cookie = `${VISITOR_KEY}=${encodeURIComponent(id)};path=/;max-age=${365 * 24 * 3600};SameSite=Lax`;
    return id;
  } catch {
    return randomId();
  }
}

export function visitorHeaders() {
  return { "X-Visitor-Id": ensureVisitorId() };
}
