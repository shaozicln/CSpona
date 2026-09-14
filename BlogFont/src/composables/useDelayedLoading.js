import { ref, onUnmounted } from "vue";

/**
 * 短延迟防闪：请求很快完成则不展示；一旦展示则至少停留一小段时间。
 * delayMs <= 0 时立即展示（避免 setTimeout 被更快的请求抢先取消）。
 */
export function useDelayedLoading({ delayMs = 150, minShowMs = 220 } = {}) {
  const pending = ref(false);
  const visible = ref(false);

  let delayTimer = null;
  let shownAt = 0;

  const clearDelay = () => {
    if (delayTimer != null) {
      clearTimeout(delayTimer);
      delayTimer = null;
    }
  };

  const reveal = () => {
    if (!pending.value || visible.value) return;
    visible.value = true;
    shownAt = Date.now();
  };

  const start = () => {
    pending.value = true;
    clearDelay();
    if (delayMs <= 0) {
      reveal();
    } else {
      delayTimer = setTimeout(reveal, delayMs);
    }
  };

  const stop = async () => {
    pending.value = false;
    clearDelay();
    if (visible.value) {
      const left = minShowMs - (Date.now() - shownAt);
      if (left > 0) {
        await new Promise((r) => setTimeout(r, left));
      }
    }
    visible.value = false;
  };

  onUnmounted(() => {
    clearDelay();
    pending.value = false;
    visible.value = false;
  });

  return { pending, visible, start, stop };
}
