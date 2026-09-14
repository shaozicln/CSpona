<template>
  <div class="music-fields">
    <label class="mf-label">背景音乐（可选，不选则沿用当前/站点默认）</label>
    <div class="mf-modes">
      <label><input v-model="mode" type="radio" value="" /> 用默认</label>
      <label><input v-model="mode" type="radio" value="netease" /> 网易云单曲</label>
      <label><input v-model="mode" type="radio" value="playlist" /> 网易云歌单</label>
      <label><input v-model="mode" type="radio" value="file" /> 上传音频</label>
    </div>
    <div v-if="mode === 'netease' || mode === 'playlist'" class="mf-row">
      <input
        v-model="neteaseInput"
        type="text"
        :placeholder="mode === 'playlist' ? '歌单 ID 或分享链接' : '歌曲 ID 或分享链接'"
        @change="emitValue"
      />
    </div>
    <div v-if="mode === 'file'" class="mf-row">
      <input type="file" accept="audio/*,.mp3,.m4a,.flac,.ogg,.wav,.aac" @change="onFile" />
      <p v-if="fileRef" class="mf-hint">已选：{{ fileRef }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import { extractNeteaseId } from "@/utils/music.js";
import { apiFetch, promptLoginIfUnauthorized } from "@/utils/api";

const props = defineProps({
  musicType: { type: String, default: "" },
  musicRef: { type: String, default: "" },
});
const emit = defineEmits(["update:musicType", "update:musicRef"]);

function normalizeType(t) {
  if (t === "netease" || t === "playlist" || t === "file") return t;
  return "";
}

const mode = ref(normalizeType(props.musicType));
const neteaseInput = ref(
  props.musicType === "netease" || props.musicType === "playlist" ? props.musicRef : ""
);
const fileRef = ref(props.musicType === "file" ? props.musicRef : "");

watch(
  () => [props.musicType, props.musicRef],
  ([t, r]) => {
    mode.value = normalizeType(t);
    if (t === "netease" || t === "playlist") neteaseInput.value = r || "";
    if (t === "file") fileRef.value = r || "";
  }
);

watch(mode, (m) => {
  if (m === "") {
    emit("update:musicType", "");
    emit("update:musicRef", "");
  } else if (m === "netease" || m === "playlist") {
    const id = extractNeteaseId(neteaseInput.value) || neteaseInput.value.trim();
    emit("update:musicType", m);
    emit("update:musicRef", id);
  } else if (m === "file") {
    emit("update:musicType", "file");
    emit("update:musicRef", fileRef.value);
  }
});

function emitValue() {
  if (mode.value === "netease" || mode.value === "playlist") {
    const id = extractNeteaseId(neteaseInput.value) || neteaseInput.value.trim();
    emit("update:musicType", mode.value);
    emit("update:musicRef", id);
  }
}

async function onFile(e) {
  const file = e.target.files?.[0];
  if (!file) return;
  if (file.size > 15 << 20) {
    alert("音频不能超过 15MB");
    e.target.value = "";
    return;
  }
  const formData = new FormData();
  formData.append("file", file);
  try {
    const res = await apiFetch("/music/upload", { method: "POST", body: formData });
    const data = await res.json().catch(() => ({}));
    if (promptLoginIfUnauthorized(res, data)) return;
    if (!res.ok) {
      alert(data.error || "上传失败");
      return;
    }
    fileRef.value = data.path || "";
    emit("update:musicType", "file");
    emit("update:musicRef", fileRef.value);
  } catch (err) {
    console.error(err);
    alert("上传失败");
  }
}

defineExpose({
  getPayload() {
    if (mode.value === "netease" || mode.value === "playlist") {
      const id = extractNeteaseId(neteaseInput.value) || neteaseInput.value.trim();
      return { music_type: id ? mode.value : "", music_ref: id };
    }
    if (mode.value === "file") {
      return { music_type: fileRef.value ? "file" : "", music_ref: fileRef.value || "" };
    }
    return { music_type: "", music_ref: "" };
  },
});
</script>

<style scoped>
.music-fields {
  margin: 12px 0 16px;
  padding: 12px;
  border: 1px solid var(--panel-border, #ddd);
  border-radius: 8px;
  background: var(--panel-bg-soft, rgba(255, 255, 255, 0.5));
}
.mf-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-secondary, #666);
}
.mf-modes {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-primary);
}
.mf-modes label {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
}
.mf-row input[type="text"] {
  width: 100%;
  box-sizing: border-box;
  padding: 8px 10px;
  border: 1px solid var(--input-border, #ccc);
  border-radius: 6px;
  background: var(--input-bg, #fff);
  color: var(--input-text, #111);
}
.mf-hint {
  margin: 6px 0 0;
  font-size: 13px;
  color: var(--text-muted, #888);
}
</style>
