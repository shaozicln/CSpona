<template>
  <div class="container">
    <div class="page-bg" aria-hidden="true">
      <img v-if="coverUrl" :src="coverUrl" alt="" class="bg-img" />
      <div class="bg-mask"></div>
    </div>
    <div class="grid-container">
      <div class="left-side">
        <transition name="fade">
          <div class="author-container">
            <AuthorBack />
          </div>
        </transition>
        <transition name="toc-transition">
          <div class="toc-box">
            <PageLoading
              v-if="showDetailLoading"
              variant="toc"
              message="索引加载中…"
            />
            <template v-else-if="!loadError">
              <div id="toc-title">内容索引</div>
              <ul>
                <li
                  v-for="item in toc"
                  :key="item.id"
                  class="toc-item"
                  :class="[
                    `toc-level-${Math.min(item.level || 1, 4)}`,
                    { active: activeSlug === item.id },
                  ]"
                  @click.prevent="scrollToContent(item.id)"
                >
                  <a :href="'#' + item.id">{{ item.text }}</a>
                </li>
              </ul>
            </template>
          </div>
        </transition>
      </div>
      <div class="right-side">
        <div class="white-box" ref="scrollBox">
          <PageLoading
            v-if="showDetailLoading"
            variant="detail"
            message="文章加载中…"
          />
          <PageLoading
            v-else-if="loadError"
            :error="loadError"
            @retry="retryLoad"
          />
          <template v-else>
            <div class="title-box">
              <div style="display: flex; justify-content: space-between; align-items: center;">
                <h1>{{ article.Title }}</h1>
                <button v-if="userQx === 'A'" class="edit-btn" @click="showEditModal = true">
                  <i class="fas fa-edit"></i>
                </button>
              </div>
              <p>
                浏览量：{{ article.ViewCount }} | 创建时间：{{
                  TimeFormat(article.CreatedAt)
                }}
                | 评论数：{{ article.CommentCount }}
              </p>
            </div>
            <div class="content-box" v-html="renderedContent"></div>
            <div class="comment-box" v-if="articleId">
              <Comments
                :article-id="articleId"
                :initial-comments="pageComments"
                :key="'comments-' + articleId"
              />
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
  <div v-if="showEditModal" class="modal-overlay" @click="showEditModal = false">
    <div class="modal-content" @click.stop>
      <ArticlePut :article-id="articleId" :article-data="article" @close="showEditModal = false" />
    </div>
  </div>
</template>

<script setup>
import Comments from "./Comments.vue";
import { library } from "@fortawesome/fontawesome-svg-core";
import { faEdit } from "@fortawesome/free-solid-svg-icons";
library.add(faEdit);

import {
  getCurrentInstance,
  ref,
  computed,
  watch,
  nextTick,
} from "vue";
import AuthorBack from "../Author/AuthorBack.vue";
import { marked } from "marked";
import hljs from "highlight.js";
import "highlight.js/styles/github.css";
import { useRoute } from "vue-router";
import { decodeArticleId } from "@/utils/utils.js";
import { resolveImageUrl, headingIdFromCounter } from "@/utils/image.js";
import ArticlePut from "./ArticlePut.vue";
import { useArticleStore } from "@/stores/article";
import PageLoading from "@/components/common/PageLoading.vue";
import { useDelayedLoading } from "@/composables/useDelayedLoading.js";

const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;
const route = useRoute();
const articleStore = useArticleStore();

const userQx = ref(localStorage.getItem("userQx") || "");
const showEditModal = ref(false);
const articleId = ref("");
const article = ref({});
const toc = ref([]);
const pageComments = ref([]);
const scrollBox = ref(null);
const activeSlug = ref("");
const loadError = ref("");
let headingCounter = {};
let loadSeq = 0;

const { pending: detailPending, visible: detailVisible, start: startDetailLoading, stop: stopDetailLoading } =
  useDelayedLoading({ delayMs: 0, minShowMs: 320 });
const showDetailLoading = computed(() => detailPending.value || detailVisible.value);

const coverUrl = computed(() => {
  // 加载中 / 无封面时不请求默认图，避免骨架闪过后蹦出 boli.jpg
  if (showDetailLoading.value) return "";
  const a = article.value;
  if (!a) return "";
  if (a.ImgUrl) return resolveImageUrl(a.ImgUrl, "");
  if (a.Img) return resolveImageUrl(a.Img, "");
  return "";
});

// marked v15：renderer 收到 token 对象
const renderer = {
  heading({ tokens, depth, text }) {
    const plain = text || "";
    const id = headingIdFromCounter(plain, headingCounter);
    const inner = this.parser.parseInline(tokens);
    return `<h${depth} id="${id}">${inner}</h${depth}>\n`;
  },
  code({ text, lang }) {
    const language = (lang || "").trim();
    let highlighted = "";
    try {
      if (language && hljs.getLanguage(language)) {
        highlighted = hljs.highlight(text, {
          language,
          ignoreIllegals: true,
        }).value;
      } else {
        // 不用 highlightAuto（很慢），无语言时只做转义
        highlighted = text
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;");
      }
    } catch {
      highlighted = text
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;");
    }
    const icon = `<svg class="copy-icon" viewBox="0 0 24 24" width="18" height="18" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15V5a2 2 0 0 1 2-2h10"/></svg>`;
    const check = `<svg class="copy-icon copy-icon-check" viewBox="0 0 24 24" width="18" height="18" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6L9 17l-5-5"/></svg>`;
    return `<div class="code-block"><button type="button" class="copy-button" title="复制" aria-label="复制代码">${icon}${check}</button><pre><code class="language-${language}">${highlighted}</code></pre></div>\n`;
  },
};

marked.setOptions({ gfm: true, breaks: true });
marked.use({ renderer });

const renderedContent = computed(() => {
  if (!article.value?.Content) return "";
  headingCounter = {};
  return marked.parse(article.value.Content);
});

async function loadPageArticle(id) {
  if (!id) return;
  const seq = ++loadSeq;
  loadError.value = "";
  startDetailLoading();
  try {
    const response = await fetch(`${URL}/page/article/${id}`);
    const data = await response.json();
    if (seq !== loadSeq) return;
    if (!data.data) {
      console.error("加载文章失败:", data.msg);
      loadError.value = data.msg || "加载文章失败，请稍后重试";
      article.value = {};
      toc.value = [];
      pageComments.value = [];
      return;
    }
    const payload = data.data;
    article.value = payload.article || {};
    toc.value = (payload.toc || []).map((t) => ({
      id: t.id || t.ID || t.Id,
      text: t.text || t.Text,
      level: t.level || t.Level,
    }));
    pageComments.value = payload.comments || [];
    activeSlug.value = toc.value[0]?.id || "";
    await nextTick();
    if (seq !== loadSeq) return;
    bindCopyButtons();
    const hash = decodeURIComponent(window.location.hash.replace(/^#/, ""));
    if (hash) scrollToContent(hash, false);
  } catch (error) {
    if (seq !== loadSeq) return;
    console.error("Failed to fetch article page:", error);
    loadError.value = "加载文章失败，请稍后重试";
    article.value = {};
    toc.value = [];
    pageComments.value = [];
  } finally {
    if (seq === loadSeq) await stopDetailLoading();
  }
}

function retryLoad() {
  if (articleId.value) loadPageArticle(articleId.value);
}

watch(
  () => route.params.articleId,
  async (newId) => {
    if (!newId) return;
    const decodedId = decodeArticleId(newId);
    articleId.value = decodedId;
    localStorage.setItem("articleId", decodedId);
    articleStore.setArticleId(decodedId);
    await loadPageArticle(decodedId);
  },
  { immediate: true }
);

function TimeFormat(time) {
  if (!time) return "";
  const date = new Date(time);
  return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()} - ${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
}

function bindCopyButtons() {
  document.querySelectorAll(".copy-button").forEach((button) => {
    button.onclick = async () => {
      const codeBlock = button.parentElement?.querySelector("code");
      const codeText = codeBlock?.innerText || "";
      try {
        if (navigator.clipboard?.writeText) {
          await navigator.clipboard.writeText(codeText);
        } else {
          const textarea = document.createElement("textarea");
          textarea.value = codeText;
          document.body.appendChild(textarea);
          textarea.select();
          document.execCommand("copy");
          document.body.removeChild(textarea);
        }
        button.classList.add("copied");
        button.title = "已复制";
        setTimeout(() => {
          button.classList.remove("copied");
          button.title = "复制";
        }, 1600);
      } catch (e) {
        console.error("复制失败", e);
      }
    };
  });
}

function scrollToContent(slug, updateHash = true) {
  const box = scrollBox.value;
  const element = document.getElementById(slug);
  if (!box || !element) {
    console.error(`Element not found for slug: ${slug}`);
    return;
  }
  const boxRect = box.getBoundingClientRect();
  const elRect = element.getBoundingClientRect();
  const top = elRect.top - boxRect.top + box.scrollTop - 16;
  box.scrollTo({ top, behavior: "smooth" });
  activeSlug.value = slug;
  if (updateHash) {
    history.replaceState(null, "", `#${slug}`);
  }
}

function handleScroll() {
  const box = scrollBox.value;
  if (!box || !toc.value.length) return;
  const boxTop = box.getBoundingClientRect().top;
  const probe = boxTop + 48;
  let current = toc.value[0]?.id || "";
  for (const item of toc.value) {
    const el = document.getElementById(item.id);
    if (!el) continue;
    if (el.getBoundingClientRect().top <= probe) {
      current = item.id;
    } else {
      break;
    }
  }
  activeSlug.value = current;
}

watch(scrollBox, (box, _, onCleanup) => {
  if (!box) return;
  box.addEventListener("scroll", handleScroll, { passive: true });
  onCleanup(() => box.removeEventListener("scroll", handleScroll));
});
</script>

<style scoped>
.container {
  position: relative;
  height: calc(100vh - 4.75rem);
  overflow: hidden;
}

.page-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
  background: transparent;
}

.page-bg .bg-img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.page-bg .bg-mask {
  position: absolute;
  inset: 0;
  background: var(--bg-mask);
  transition: background 0.25s ease;
}

.grid-container {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 1fr 3fr;
  grid-gap: 20px;
  padding: 20px;
  height: calc(100% - 0px);
  box-sizing: border-box;
}

.left-side {
  grid-column: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 0;
  overflow: auto;
}

.right-side {
  grid-column: 2;
  min-height: 0;
  height: 100%;
}

.white-box {
  background-color: var(--panel-bg);
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
  overflow-y: auto;
  color: var(--text-primary);
}

.title-box {
  margin-bottom: 20px;
}

.title-box h1 {
  font-size: 40px;
  margin-bottom: 10px;
  color: var(--text-primary);
}

.title-box p {
  font-size: 15px;
  color: var(--text-secondary);
}

.content-box {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 20px;
  overflow-x: hidden;
  line-height: 1.8;
  font-size: 16px;
  color: var(--text-primary);
}

.content-box :deep(h1),
.content-box :deep(h2),
.content-box :deep(h3),
.content-box :deep(h4) {
  scroll-margin-top: 16px;
}

.content-box :deep(p) {
  margin-bottom: 20px !important;
  margin-top: 0 !important;
}

.content-box :deep(img) {
  max-width: 100% !important;
  height: auto !important;
  display: block !important;
}

.comment-box {
  margin-top: 32px;
  margin-bottom: 40px;
  padding-bottom: 24px;
}

/* 内容索引：层级左对齐 + 聚焦竖线高亮 */
.toc-box {
  margin: 20px 0;
  background-color: var(--panel-bg-soft);
  padding: 16px 12px;
  height: calc(100vh - 100px);
  overflow-y: auto;
  border-radius: 20px;
  border: solid 2px var(--panel-border);
  font-size: 16px;
  width: 300px;
  text-align: left;
  box-sizing: border-box;
}

.toc-box ul {
  padding: 0;
  margin: 0;
  list-style: none;
}

.toc-box .toc-item {
  margin: 2px 0;
  border-radius: 6px;
  text-align: left;
}

.toc-box .toc-item a {
  display: block;
  padding-top: 7px;
  padding-bottom: 7px;
  padding-right: 10px;
  padding-left: 8px;
  text-decoration: none;
  color: var(--text-muted);
  line-height: 1.35;
  border-left: 3px solid transparent;
  border-radius: 6px;
  transition: background 0.15s ease, color 0.15s ease, border-color 0.15s ease;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toc-box .toc-item:hover a {
  color: var(--text-primary);
  background: var(--toc-active-bg);
  text-decoration: none;
}

/* 一级顶格；二三四级递增缩进（提高优先级压过通用 a） */
.toc-box .toc-item.toc-level-1 > a {
  padding-left: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}
.toc-box .toc-item.toc-level-2 > a {
  padding-left: 22px;
  font-size: 15px;
}
.toc-box .toc-item.toc-level-3 > a {
  padding-left: 36px;
  font-size: 14px;
}
.toc-box .toc-item.toc-level-4 > a {
  padding-left: 50px;
  font-size: 13px;
}

#toc-title {
  font-size: 22px;
  text-align: center;
  margin-bottom: 12px;
  color: var(--text-primary);
}

.toc-box .toc-item.active > a {
  color: var(--toc-active-text);
  font-weight: 600;
  background: var(--toc-active-bg);
  border-left-color: var(--toc-active-bar);
}

.edit-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 20px;
  color: var(--text-primary);
}

.edit-btn i {
  color: inherit;
}

.edit-btn:hover {
  color: var(--text-secondary);
}

pre {
  background-color: var(--code-bg);
  padding: 16px;
  border-radius: 4px;
  overflow-x: auto;
  font-family: "Consolas", "Monaco", "Courier New", monospace;
  font-size: 14px;
  line-height: 1.5;
  margin: 0px;
  border: 1px solid var(--code-border);
}

code {
  font-family: "Consolas", "Monaco", "Courier New", monospace;
  font-size: 14px;
  color: #d63384;
}

pre code {
  color: inherit;
  background-color: transparent;
  padding: 0;
  margin: 0px;
  line-height: 1.5;
}

/* v-html 注入的节点必须用 :deep，否则 position 不生效会跑到代码块外面 */
.content-box :deep(.code-block) {
  position: relative;
  margin: 12px 0;
  border-radius: 6px;
  border: 1px solid var(--code-border);
  background-color: var(--code-bg);
  overflow: hidden;
  color: var(--code-text);
}

.content-box :deep(.code-block pre),
.content-box :deep(.code-block code) {
  margin: 0;
  color: var(--code-text);
  background: transparent;
}

.content-box :deep(.code-block pre) {
  padding: 18px 16px 16px;
  border: none;
  border-radius: 0;
  overflow-x: auto;
}

.content-box :deep(.copy-button) {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 3;
  width: 30px;
  height: 30px;
  padding: 0;
  margin: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: var(--copy-icon);
  cursor: pointer;
  line-height: 0;
  opacity: 0.85;
  transition: opacity 0.15s ease, color 0.15s ease;
}

.content-box :deep(.code-block:hover .copy-button),
.content-box :deep(.copy-button:hover),
.content-box :deep(.copy-button.copied) {
  opacity: 1;
}

.content-box :deep(.copy-button:hover) {
  background: transparent;
  color: var(--copy-icon-hover);
}

.content-box :deep(.copy-button .copy-icon) {
  display: block;
  pointer-events: none;
}

.content-box :deep(.copy-button .copy-icon-check) {
  display: none;
}

.content-box :deep(.copy-button.copied) {
  background: transparent;
  color: #1a7f37;
}

.content-box :deep(.copy-button.copied .copy-icon:not(.copy-icon-check)) {
  display: none;
}

.content-box :deep(.copy-button.copied .copy-icon-check) {
  display: block;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  max-width: 90vw;
  max-height: 90vh;
  overflow: auto;
  background: var(--panel-bg);
  color: var(--text-primary);
  border: 1px solid var(--panel-border);
  border-radius: 12px;
  padding: 16px;
}

.fade-enter-active,
.fade-leave-active,
.toc-transition-enter-active,
.toc-transition-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to,
.toc-transition-enter-from,
.toc-transition-leave-to {
  opacity: 0;
}
</style>
