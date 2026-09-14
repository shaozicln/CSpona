<template>
  <div
    class="page-loading"
    :class="[`is-${variant}`, { 'is-error': !!error }]"
    role="status"
    :aria-busy="!error"
    :aria-live="error ? 'assertive' : 'polite'"
  >
    <div v-if="error" class="pl-error">
      <p class="pl-error-text">{{ error }}</p>
      <button
        v-if="retryable"
        type="button"
        class="pl-retry"
        @click="$emit('retry')"
      >
        重试
      </button>
    </div>

    <template v-else>
      <div class="pl-head">
        <div class="pl-spinner" aria-hidden="true"></div>
        <p class="pl-msg">{{ message }}</p>
      </div>

      <!-- 详情骨架：标题 + 多行正文 -->
      <div v-if="variant === 'detail'" class="pl-skeleton pl-detail">
        <div class="bone bone-title"></div>
        <div class="bone bone-meta"></div>
        <div class="bone bone-line"></div>
        <div class="bone bone-line"></div>
        <div class="bone bone-line short"></div>
        <div class="bone bone-line"></div>
        <div class="bone bone-line short"></div>
      </div>

      <!-- 文章列表骨架：分类标题 + 卡片格 -->
      <div v-else-if="variant === 'list'" class="pl-skeleton pl-list">
        <div class="bone bone-section"></div>
        <div class="pl-grid">
          <div v-for="n in 4" :key="'g' + n" class="pl-card-bone">
            <div class="bone bone-thumb"></div>
            <div class="bone bone-line"></div>
          </div>
        </div>
      </div>

      <!-- 漫游地左侧卡片骨架 -->
      <div v-else-if="variant === 'cards'" class="pl-skeleton pl-cards">
        <div v-for="n in 3" :key="'c' + n" class="pl-side-card">
          <div class="bone bone-thumb tall"></div>
          <div class="bone bone-line"></div>
          <div class="bone bone-line short"></div>
        </div>
      </div>

      <!-- TOC 侧栏骨架 -->
      <div v-else-if="variant === 'toc'" class="pl-skeleton pl-toc">
        <div class="bone bone-section"></div>
        <div class="bone bone-line"></div>
        <div class="bone bone-line indent"></div>
        <div class="bone bone-line indent"></div>
        <div class="bone bone-line"></div>
        <div class="bone bone-line indent"></div>
      </div>
    </template>
  </div>
</template>

<script setup>
defineProps({
  /** detail | list | cards | toc | spinner */
  variant: { type: String, default: "spinner" },
  message: { type: String, default: "加载中…" },
  error: { type: String, default: "" },
  retryable: { type: Boolean, default: true },
});

defineEmits(["retry"]);
</script>

<style scoped>
.page-loading {
  width: 100%;
  box-sizing: border-box;
  padding: 24px 20px;
  color: var(--text-secondary);
}

.pl-head {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 20px;
}

.pl-spinner {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid var(--panel-border);
  border-top-color: rgb(139, 189, 234);
  animation: pl-spin 0.75s linear infinite;
  flex-shrink: 0;
}

.pl-msg {
  margin: 0;
  font-size: 15px;
  color: var(--text-secondary);
}

.pl-error {
  text-align: center;
  padding: 28px 12px;
}

.pl-error-text {
  margin: 0 0 14px;
  color: var(--text-primary);
  font-size: 15px;
}

.pl-retry {
  padding: 8px 18px;
  border: none;
  border-radius: 6px;
  background: rgb(139, 189, 234);
  color: #fff;
  cursor: pointer;
  font-size: 14px;
}

.pl-retry:hover {
  background: rgb(119, 172, 218);
}

.bone {
  border-radius: 8px;
  background: linear-gradient(
    90deg,
    var(--hover-bg) 25%,
    rgba(255, 255, 255, 0.12) 50%,
    var(--hover-bg) 75%
  );
  background-size: 200% 100%;
  animation: pl-shimmer 1.2s ease-in-out infinite;
}

.bone-title {
  height: 28px;
  width: 55%;
  margin-bottom: 12px;
}

.bone-meta {
  height: 14px;
  width: 40%;
  margin-bottom: 22px;
}

.bone-line {
  height: 14px;
  width: 100%;
  margin-bottom: 12px;
}

.bone-line.short {
  width: 68%;
}

.bone-line.indent {
  width: 82%;
  margin-left: 16px;
}

.bone-section {
  height: 22px;
  width: 30%;
  margin: 8px auto 18px;
}

.bone-thumb {
  height: 120px;
  width: 100%;
  margin-bottom: 10px;
  border-radius: 10px;
}

.bone-thumb.tall {
  height: 140px;
}

.pl-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 12px;
}

.pl-card-bone {
  width: calc(50% - 8px);
  min-width: 140px;
}

.pl-cards {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.pl-side-card {
  border-radius: 8px;
  overflow: hidden;
}

.pl-toc .bone-section {
  margin-left: 0;
  width: 50%;
}

.is-spinner .pl-skeleton {
  display: none;
}

@keyframes pl-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes pl-shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
