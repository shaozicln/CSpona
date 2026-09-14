<template>
  <div class="home-container">
    <div id="TypingTextContainer">
      <div id="TypingText">
        <span>{{ displayText }}</span>
        <span class="blink">|</span>
      </div>
      <div class="welcome-container">
        <p>----------Welcome----------</p>
      </div>
    </div>
    <div class="myself">
      <div class="me">
        <Author />
      </div>
      <div class="markdown-content">
        <div v-if="aboutError" class="about-error">
          <p>{{ aboutError }}</p>
          <button type="button" class="about-retry" @click="loadAboutMe">重试</button>
        </div>
        <div v-else-if="!aboutReady" class="about-skeleton" aria-busy="true" aria-live="polite">
          <div class="sk sk-title"></div>
          <div class="sk sk-line"></div>
          <div class="sk sk-line"></div>
          <div class="sk sk-line short"></div>
          <div class="sk sk-gap"></div>
          <div class="sk sk-line"></div>
          <div class="sk sk-line short"></div>
          <div class="sk sk-line"></div>
          <div class="sk sk-line short"></div>
          <div class="sk sk-line"></div>
        </div>
        <div v-else class="about-body" v-html="aboutHtml"></div>
      </div>
    </div>
    <div class="contact-info" style="display: flex; align-items: center; gap: 4px;flex-direction:column;justify-content:center">
      <span>----------欢迎来到CSpona!----------</span>
      <span style="display: flex; align-items: center; gap: 4px"
        >备案号：<a
          href="https://beian.miit.gov.cn/"
          target="_blank"
          rel="noopener noreferrer"
          style="color:rgb( 79, 79, 79);"
        >
          黑ICP备2025038200号-1
        </a>
      </span>
    </div>
  </div>
</template>

<script setup>
import Author from "@/components/Author/Author.vue";
import { apiJson } from "@/utils/api";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { ref, onMounted, onUnmounted, computed } from "vue";

marked.setOptions({ gfm: true, breaks: true });

const aboutMd = ref("");
const aboutReady = ref(false);
const aboutError = ref("");

const aboutHtml = computed(() => {
  if (!aboutMd.value) return "<p>暂无内容</p>";
  return DOMPurify.sanitize(marked.parse(aboutMd.value));
});

async function loadAboutMe() {
  aboutError.value = "";
  aboutReady.value = false;
  try {
    const { res, data } = await apiJson("/about-me");
    if (!res.ok) {
      aboutError.value = data?.message || "加载失败";
      return;
    }
    aboutMd.value = data?.data?.content || "";
  } catch (e) {
    console.error(e);
    aboutError.value = "加载失败";
  } finally {
    aboutReady.value = true;
  }
}

const fullTexts = ref([
  "欢迎你的到来, 新朋友 ..",
  "感受Coding带来的创造的乐趣吧 !",
]);

const displayText = ref("");
const currentTextIndex = ref(0);
const currentProgress = ref(0);
const speed = 4;
let intervalId;

const currentText = computed(() => {
  return fullTexts.value[currentTextIndex.value].substring(
    0,
    currentProgress.value
  );
});

onMounted(() => {
  loadAboutMe();
  intervalId = setInterval(() => {
    if (
      currentProgress.value < fullTexts.value[currentTextIndex.value].length
    ) {
      currentProgress.value += 1;
    } else {
      currentProgress.value = 0;
      currentTextIndex.value =
        (currentTextIndex.value + 1) % fullTexts.value.length;
    }
    displayText.value = currentText.value;
  }, 1000 / speed);
});

onUnmounted(() => {
  clearInterval(intervalId);
});
</script>

<style scoped>
.myself {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.home-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.welcome-container {
  position: absolute;
  top: 62%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 30px;
  font-family: Cormorant SC, serif;
  text-align: center;
}

#TypingTextContainer {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  height: 60vh;
  width: 100vw;
  font-size: 30px;
  font-family: cursive;
  position: relative;
}

#TypingText {
  position: relative;
}

.blink {
  animation: blink 0.75s infinite;
}

@keyframes blink {
  0% {
    opacity: 1;
  }

  50% {
    opacity: 0;
  }

  100% {
    opacity: 1;
  }
}

.me {
  margin: 25px;
}

.markdown-content {
  margin: 25px;
  font-size: 18px;
  font-family: "楷体";
  line-height: 1.6;
  color: var(--text-primary);
  padding: 20px 20px 40px 20px;
  background-color: var(--panel-bg);
  border-radius: 8px;
  max-width: 800px;
  min-width: 280px;
}

.about-error {
  color: var(--text-secondary);
  padding: 12px 0;
  text-align: center;
}

.about-retry {
  margin-top: 10px;
  padding: 6px 14px;
  border: none;
  border-radius: 6px;
  background: var(--btn-bg);
  color: var(--btn-text);
  cursor: pointer;
}

.about-skeleton {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 280px;
  padding: 4px 0;
}

.sk {
  border-radius: 8px;
  background: linear-gradient(
    90deg,
    var(--hover-bg) 25%,
    rgba(255, 255, 255, 0.18) 50%,
    var(--hover-bg) 75%
  );
  background-size: 200% 100%;
  animation: home-shimmer 1.2s ease-in-out infinite;
}

.sk-title {
  height: 32px;
  width: 36%;
  margin-bottom: 6px;
}

.sk-line {
  height: 16px;
  width: 100%;
}

.sk-line.short {
  width: 72%;
}

.sk-gap {
  height: 10px;
  width: 40%;
  opacity: 0;
}

@keyframes home-shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin-top: 20px;
  margin-bottom: 10px;
}

.markdown-content :deep(h1:first-child) {
  margin-top: 0;
}

.markdown-content :deep(p) {
  margin-bottom: 15px;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
  padding-left: 1.2em;
  list-style-position: outside;
}

.markdown-content :deep(li) {
  font-size: 16px;
}

.markdown-content :deep(a) {
  color: var(--text-primary);
  text-decoration: underline;
}

.contact-info {
  width: 100vw;
  margin-top: 100px;
  margin-bottom: 0px;
  font-size: 20px;
  color: var(--text-secondary);
  text-align: center;
  background-color: var(--panel-bg);
  padding: 40px;
}

.contact-info span {
  margin: 5px 0;
}
</style>
