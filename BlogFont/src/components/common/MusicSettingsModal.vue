<template>
  <div v-if="music.settingsOpen" class="ms-overlay" @click.self="music.closeSettings()">
    <div class="ms-modal" role="dialog" aria-label="音乐设置">
      <div class="ms-header">
        <h3>音乐设置</h3>
        <button type="button" class="ms-close" @click="music.closeSettings()">×</button>
      </div>
      <div class="ms-body">
        <label class="ms-row">
          <span>总开关</span>
          <select :value="music.prefs.enabled ? 'on' : 'off'" @change="onEnabled">
            <option value="on">开</option>
            <option value="off">关</option>
          </select>
        </label>
        <label class="ms-row">
          <span>默认页（首页/列表等）</span>
          <select :value="music.prefs.defaultPages ? 'on' : 'off'" @change="onDefault">
            <option value="on">开</option>
            <option value="off">关</option>
          </select>
        </label>
        <label class="ms-row">
          <span>文章/漫游详情</span>
          <select :value="music.prefs.detailPages ? 'on' : 'off'" @change="onDetail">
            <option value="on">开</option>
            <option value="off">关</option>
          </select>
        </label>
        <label class="ms-row">
          <span>播放方式</span>
          <select :value="music.prefs.autoplay ? 'auto' : 'click'" @change="onAutoplay">
            <option value="auto">进入网页自动播放</option>
            <option value="click">点击后才播放</option>
          </select>
        </label>

        <div class="ms-block">
          <label class="ms-row">
            <span>我的全局歌单</span>
            <select :value="music.prefs.customEnabled ? 'on' : 'off'" @change="onCustomEnabled">
              <option value="off">用站点默认</option>
              <option value="on">用我的</option>
            </select>
          </label>
          <template v-if="music.prefs.customEnabled">
            <label class="ms-row">
              <span>类型</span>
              <select :value="music.prefs.customMode" @change="onCustomMode">
                <option value="playlist">网易云歌单</option>
                <option value="song">网易云单曲</option>
              </select>
            </label>
            <label class="ms-row ms-vol">
              <span>{{ music.prefs.customMode === 'song' ? '歌曲 ID / 链接' : '歌单 ID / 链接' }}</span>
              <input
                type="text"
                class="ms-input"
                :value="music.prefs.customId"
                :placeholder="music.prefs.customMode === 'song' ? 'song id' : 'playlist id'"
                @change="onCustomId"
              />
            </label>
          </template>
        </div>

        <label class="ms-row ms-vol">
          <span>音量 {{ Math.round(music.prefs.volume * 100) }}%</span>
          <input
            type="range"
            min="0"
            max="100"
            :value="Math.round(music.prefs.volume * 100)"
            @input="onVolume"
          />
        </label>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useMusicStore } from "@/stores/music";
import { extractNeteaseId } from "@/utils/music.js";

const music = useMusicStore();

function onEnabled(e) {
  music.setPref({ enabled: e.target.value === "on" });
}
function onDefault(e) {
  music.setPref({ defaultPages: e.target.value === "on" });
}
function onDetail(e) {
  music.setPref({ detailPages: e.target.value === "on" });
}
function onAutoplay(e) {
  music.setPref({ autoplay: e.target.value === "auto" });
}
function onCustomEnabled(e) {
  music.setPref({ customEnabled: e.target.value === "on" });
}
function onCustomMode(e) {
  music.setPref({ customMode: e.target.value === "song" ? "song" : "playlist" });
}
function onCustomId(e) {
  const raw = e.target.value || "";
  const id = extractNeteaseId(raw) || raw.trim();
  music.setPref({ customId: id });
}
function onVolume(e) {
  music.setPref({ volume: Number(e.target.value) / 100 });
}
</script>

<style scoped>
.ms-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
.ms-modal {
  width: min(420px, 100%);
  background: var(--float-panel-bg, var(--panel-bg, #fff));
  color: var(--text-primary, #111);
  border-radius: 12px;
  border: 1px solid var(--panel-border, #eee);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
}
.ms-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid var(--panel-border, #eee);
}
.ms-header h3 {
  margin: 0;
  font-size: 18px;
}
.ms-close {
  border: none;
  background: transparent;
  font-size: 24px;
  cursor: pointer;
  color: var(--text-primary);
  line-height: 1;
}
.ms-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.ms-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 15px;
}
.ms-row select {
  min-width: 160px;
  max-width: 220px;
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid var(--input-border, #ccc);
  background: var(--input-bg, #fff);
  color: var(--input-text, #111);
}
.ms-vol {
  flex-direction: column;
  align-items: stretch;
}
.ms-vol input[type="range"] {
  width: 100%;
}
.ms-block {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-top: 4px;
  border-top: 1px dashed var(--panel-border, #ddd);
}
.ms-input {
  width: 100%;
  box-sizing: border-box;
  padding: 8px 10px;
  border-radius: 6px;
  border: 1px solid var(--input-border, #ccc);
  background: var(--input-bg, #fff);
  color: var(--input-text, #111);
}
.ms-tip {
  margin: 0;
  font-size: 12px;
  color: var(--text-muted, #888);
  line-height: 1.4;
}
</style>
