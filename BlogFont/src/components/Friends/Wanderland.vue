<template>
  <div class="wanderland">
    <!-- 背景图片 + 夜色灰色遮罩 -->
    <div class="page-bg" aria-hidden="true">
      <img v-if="coverUrl" :src="coverUrl" alt="背景图" class="bg-img" />
      <div class="bg-mask"></div>
    </div>
    <!-- 左侧分类与文章列表 -->
    <div class="slide">
      <div class="category-container">
        <PageLoading
          v-if="showCatLoading"
          variant="cards"
          message="漫游列表加载中…"
        />
        <PageLoading
          v-else-if="categoryError"
          :error="categoryError"
          @retry="fetchCategories"
        />
        <template v-else>
          <div v-for="category in categories" :key="category.Id" class="category-item">
            <div v-if="category.Id === 1000" class="special-articles">
              <div v-for="articleItem in category.Articles" :key="articleItem.Id" class="article-card"
                @click="loadAndShowArticle(articleItem.Id)" :class="{ 'active-article': articleItem.Id === articleId }">
                <img :src="getImageUrl(articleItem.Img, articleItem.ImgUrl)" alt="文章缩略图" class="article-img" loading="lazy" />
                <div class="article-info">
                  <h4 class="article-title">{{ articleItem.Title }}</h4>
                  <p class="article-meta">
                    <span class="author">
                      👤 {{ getAuthorName(articleItem.UserId) }}
                    </span>
                    <span class="date">
                      📅 {{ TimeFormat(articleItem.CreatedAt) }}
                    </span>
                  </p>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 右侧文章内容区 -->
    <div class="content">
      <PageLoading
        v-if="showArtLoading"
        variant="detail"
        message="文章加载中…"
      />
      <PageLoading
        v-else-if="articleError"
        :error="articleError"
        @retry="retryLoadArticle"
      />
      <div v-else-if="article && article.Title" class="article-content">
        <div class="white-box title-section">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            
          <h1 class="article-main-title">{{ article.Title }}</h1>
            <button v-if="isArticleOwner === true" class="edit-btn" @click="showEditModal = true">
              <i class="fas fa-edit"></i>
            </button>
          </div>
          <div class="article-meta-main">
            <span>作者: {{ getAuthorName(article.UserId) }}</span>
            <span>发布时间: {{ TimeFormat(article.CreatedAt) }}</span>
            <span>浏览量: {{ article.ViewCount || 0 }}</span>
            <span>评论数: {{ article.CommentCount || 0 }}</span>
          </div>
          <div class="content-box" v-html="renderedContent"></div>
        </div>

        <div class="white-box comment-section">
          <h3 class="comment-title" style="text-align: center; font-size: 40px; font-family: cursive">
            评论区
          </h3>
          <Comments ref="commentsRef" :article-id="articleId" :initial-comments="pageComments" :key="'comments-' + articleId" />
        </div>
      </div>

      <div v-else class="empty-state">
        <p>未找到指定文章</p>
        <button @click="goToDefaultArticle" class="default-btn">
          查看默认文章
        </button>
      </div>
    </div>
  </div>

   <!-- 编辑弹窗遮罩层 -->
  <div v-if="showEditModal" class="modal-overlay" @click="showEditModal = false">
    <div class="modal-content" @click.stop>
      <WanderlandPut :article-id="articleId" :article-data="article" @close="showEditModal = false" />
    </div>
  </div>
</template>

<script setup>
import Comments from "../Article/Comments.vue";
import { useUserStore } from "@/stores/user";
import {
  ref,
  onMounted,
  computed,
  watch,
  nextTick,
  onUnmounted,
  toRefs,
} from "vue";
import { useRouter, useRoute } from "vue-router";
import { getCurrentInstance } from "vue";
import { marked } from "marked";
import hljs from "highlight.js";
import "highlight.js/styles/github.css";
import { useArticleStore } from "@/stores/article";
import { decodeArticleId } from "@/utils/utils.js";
import { resolveImageUrl, headingIdFromCounter } from "@/utils/image.js";
import { apiFetch } from "@/utils/api.js";
import PageLoading from "@/components/common/PageLoading.vue";
import { useDelayedLoading } from "@/composables/useDelayedLoading.js";

import { library } from '@fortawesome/fontawesome-svg-core';
import { faEdit } from '@fortawesome/free-solid-svg-icons';
import WanderlandPut from "./WanderlandPut.vue";
// 注册图标 
library.add(faEdit);

// 文章更新
const userId = ref(localStorage.getItem('userId') || '');
const showEditModal = ref(false); //控制编辑弹窗显示

// 路由与状态管理
const userStore = useUserStore();
const router = useRouter();
const route = useRoute();
const articleStore = useArticleStore();
const commentsRef = ref(null);

// 响应式变量
const articleId = ref("");
const article = ref({});
const pageComments = ref([]);
const categories = ref([]);
const users = ref({}); // 缓存用户信息
const toc = ref([]);
let headingCounter = {};
const activeSlug = ref("");
const categoryError = ref(null);
const articleError = ref(null);
const isArticleOwner = ref(false);
const PINNED_WANDERLAND_ID = 95;

const {
  pending: catPending,
  visible: catVisible,
  start: startCatLoading,
  stop: stopCatLoading,
} = useDelayedLoading({ delayMs: 0, minShowMs: 280 });
const {
  pending: artPending,
  visible: artVisible,
  start: startArtLoading,
  stop: stopArtLoading,
} = useDelayedLoading({ delayMs: 0, minShowMs: 280 });
const showCatLoading = computed(() => catPending.value || catVisible.value);
const showArtLoading = computed(() => artPending.value || artVisible.value);
let artLoadSeq = 0;

const coverUrl = computed(() => {
  if (showArtLoading.value) return "";
  const a = article.value;
  if (!a) return "";
  if (a.ImgUrl) return resolveImageUrl(a.ImgUrl, "");
  if (a.Img) return resolveImageUrl(a.Img, "");
  return "";
});

// 获取全局实例与URL
const instance = getCurrentInstance();
const { proxy } = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

const renderMarkdown = (raw) => {
  if (typeof raw !== "string") {
    raw = String(raw || "");
  }
  try {
    headingCounter = {};
    return marked.parse(raw);
  } catch (e) {
    console.error("Markdown 渲染失败", e);
    return `<p>渲染失败</p>`;
  }
};

const renderedContent = computed(() => {
  if (!article.value?.Content) return "";
  return renderMarkdown(article.value.Content);
});

const fetchArticleDetail = async (id) => {
  const res = await apiFetch(`/page/article/${id}`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  const body = await res.json();
  const payload = body.data || {};
  const data = payload.article || {};

  article.value = {
    ...data,
    Content: data.Content ?? data.content ?? "",
    Title: data.Title ?? data.title ?? "未命名文章",
    ImgUrl: data.ImgUrl,
  };
  pageComments.value = payload.comments || [];
  if (payload.author) {
    users.value[payload.author.Id] = {
      Username: payload.author.Username,
      Avatar: payload.author.Avatar,
      AvatarUrl: payload.author.AvatarUrl,
    };
  }

  isArticleOwner.value = String(article.value.UserId) === String(userId.value);

  if (data.UserId && !users.value[data.UserId]) await fetchUserInfo(data.UserId);
};

// 加载并显示文章
const loadAndShowArticle = async (id) => {
  const seq = ++artLoadSeq;
  startArtLoading();
  articleError.value = null;
  article.value = {};
  try {
    articleId.value = id;
    articleStore.setArticleId(id);
    localStorage.setItem("articleId", id);
    const encodedId = window.btoa(id.toString());
    router.replace({
      name: route.name,
      params: { ...route.params, articleId: encodedId },
    });
    await fetchArticleDetail(id);
    if (seq !== artLoadSeq) return;
    console.log("文章加载完成", article.value);
    window.scrollTo({ top: 0, behavior: "smooth" });
    nextTick(() => {
      if (seq !== artLoadSeq) return;
      bindCopyButtons();
    });
  } catch (err) {
    if (seq !== artLoadSeq) return;
    console.error("加载文章失败:", err);
    articleError.value = "加载文章失败，请稍后重试";
  } finally {
    if (seq === artLoadSeq) await stopArtLoading();
  }
};

watch(
  () => route.params.articleId,
  (newParamId) => {
    console.log(`路由参数变化: ${newParamId}`);

    if (!newParamId) {
      console.log(`无路由参数，跳过`);
      return;
    }

    try {
      const decodedId = decodeArticleId(newParamId);
      console.log(`解码后ID: ${decodedId}, 当前ID: ${articleId.value}`);

      // 只有当ID变化时才加载
      if (decodedId !== articleId.value) {
        console.log(`ID变化，触发加载`);
        loadAndShowArticle(decodedId);
      } else {
        console.log(`ID相同，跳过加载`);
      }
    } catch (e) {
      console.error("路由参数解析错误:", e);
      loadAndShowArticle(PINNED_WANDERLAND_ID); // 回退到置顶公告
    }
  },
  { immediate: true }
);

// marked v15：renderer 收到 token 对象（与文章详情对齐）
const mdRenderer = {
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
  image({ href, title, text }) {
    const src = typeof href === "string" ? href : String(href || "");
    const alt = (text || "").replace(/"/g, "&quot;");
    const t = (title || "").replace(/"/g, "&quot;");
    return `<img src="${encodeURI(src)}" alt="${alt}" title="${t}" class="markdown-image" />`;
  },
};

marked.setOptions({ gfm: true, breaks: true });
marked.use({ renderer: mdRenderer });

// 获取图片URL
const getImageUrl = (imgName, resolved) => {
  return resolveImageUrl(resolved || imgName || "default-article.jpg");
};

// 时间格式化
function TimeFormat(time) {
  if (!time) return "";
  const date = new Date(time);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hour = String(date.getHours()).padStart(2, "0");
  const minute = String(date.getMinutes()).padStart(2, "0");
  return `${year}/${month}/${day} ${hour}:${minute}`;
}

const isPinnedArticle = (item) => {
  if (!item) return false;
  if (Number(item.Id) === PINNED_WANDERLAND_ID) return true;
  const title = item.Title || "";
  return title.includes("正式开放") && title.includes("Wanderland");
};

const pinWanderlandList = (list) => {
  return (list || []).map((cat) => {
    if (cat.Id !== 1000 || !Array.isArray(cat.Articles)) return cat;
    const arts = [...cat.Articles].sort((a, b) => {
      const pa = isPinnedArticle(a);
      const pb = isPinnedArticle(b);
      if (pa !== pb) return pa ? -1 : 1;
      return new Date(b.CreatedAt) - new Date(a.CreatedAt);
    });
    return { ...cat, Articles: arts };
  });
};

// 获取分类数据（BFF）
const fetchCategories = async () => {
  try {
    categoryError.value = null;
    startCatLoading();
    const response = await apiFetch(`/page/home`);
    if (!response.ok) throw new Error(`HTTP错误: ${response.status}`);

    const data = await response.json();
    const list = data?.data?.categories || data?.data || [];
    if (!list) throw new Error("数据格式不正确");

    categories.value = pinWanderlandList(list);

    const targetCategory = categories.value.find((cat) => cat.Id === 1000);
    if (targetCategory && targetCategory.Articles) {
      for (const article of targetCategory.Articles) {
        if (article.UserId && !users.value[article.UserId]) {
          await fetchUserInfo(article.UserId);
        }
      }
    }
  } catch (err) {
    console.error("获取分类失败:", err);
    categoryError.value = "无法加载分类列表，请稍后重试";
  } finally {
    await stopCatLoading();
  }
};

// 获取用户信息
const fetchUserInfo = async (uid) => {
  try {
    if (users.value[uid]) return;

    const response = await apiFetch(`/users/${uid}`);
    if (!response.ok) throw new Error(`获取用户 ${uid} 信息失败`);

    const data = await response.json();
    users.value[uid] = data.data;
  } catch (err) {
    console.error("获取用户信息失败:", err);
    users.value[uid] = {
      Username: "未知作者",
      Avatar: "default-avatar.jpg",
    };
  }
};

// 获取作者名称
const getAuthorName = (userId) => {
  return users.value[userId]?.Username || "未知作者";
};

// 复制按钮绑定（与文章详情一致，避免重复监听）
const bindCopyButtons = () => {
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
      } catch (err) {
        console.error("复制失败:", err);
      }
    };
  });
};

// 滚动处理
const handleScroll = () => {
  const scrollPosition = window.scrollY;
  const tocItems = toc.value;

  for (let i = tocItems.length - 1; i >= 0; i--) {
    const element = document.getElementById(tocItems[i].slug);
    if (element) {
      const elementTop = element.offsetTop;
      const elementBottom = element.offsetTop + element.offsetHeight;

      if (
        scrollPosition + window.innerHeight / 3 >= elementTop &&
        scrollPosition < elementBottom
      ) {
        activeSlug.value = tocItems[i].slug;
        return;
      }
    }
  }

  activeSlug.value = "";
};

// 滚动到指定内容
const scrollToContent = (slug) => {
  nextTick(() => {
    const element = document.getElementById(slug);
    if (element) {
      element.scrollIntoView({ behavior: "smooth" });
    } else {
      console.error(`未找到ID为 ${slug} 的元素`);
    }
  });
};

// 重试加载文章
const retryLoadArticle = () => {
  if (articleId.value) {
    loadAndShowArticle(articleId.value);
  } else {
    goToDefaultArticle();
  }
};

// 跳转到默认置顶文章
const goToDefaultArticle = () => {
  loadAndShowArticle(PINNED_WANDERLAND_ID);
};

// 初始化：路由 watch(immediate) 负责有参数时加载；无参数时打开置顶公告
onMounted(async () => {
  await fetchCategories();
  if (!route.params.articleId) {
    await loadAndShowArticle(articleId.value || PINNED_WANDERLAND_ID);
  }
  window.addEventListener("scroll", handleScroll);
});

// 清理工作
onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
});
</script>

<style scoped>
/* 基础布局：给顶部导航留空，避免 100vh 把 nav 顶走/盖住 */
.wanderland {
  display: flex;
  width: 100%;
  height: calc(100vh - 4.75rem);
  min-height: 0;
  gap: 10px;
}

.page-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
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

/* 左侧区域 */
.slide {
  position: relative;
  z-index: 1;
  width: 25%;
  padding: 20px;
  box-sizing: border-box;
  overflow-y: auto;
  max-height: 100%;
  /* 细边框分隔 */
  background-color: var(--panel-bg);
  color: var(--text-primary);
  height: 100%;
}

/* 右侧内容区 */
.content {
  position: relative;
  z-index: 1;
  width: 75%;
  height: 100%;
  background-color: var(--panel-bg);
  color: var(--text-primary);
  box-sizing: border-box;
  overflow-y: auto;
}

.article-count {
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: normal;
}

/* 文章卡片样式 */
.special-articles {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.article-card {
  border-radius: 8px;
  overflow: hidden;
  background-color: var(--card-bg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: background-color 0.25s ease, box-shadow 0.25s ease;
  cursor: pointer;
  border: 1px solid transparent;
}

.article-card:hover {
  background-color: var(--hover-bg);
}

.active-article {
  background-color: var(--toc-active-bg);
  box-shadow:
    inset 3px 0 0 var(--toc-active-bar),
    0 2px 10px rgba(0, 0, 0, 0.08);
  border-color: transparent;
}

.active-article .article-title {
  color: var(--toc-active-text);
}

.article-img {
  width: 100%;
  height: 140px;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.article-card:hover .article-img {
  transform: scale(1.03);
}

.article-info {
  padding: 12px;
}

.article-title {
  font-size: 16px;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  height: 40px;
  line-height: 1.5;
  padding: 0 0 8px 0;
}

.article-meta {
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

/* 文章内容区样式 */
.article-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.title-section {
  padding: 25px 30px;
}

.article-main-title {
  font-size: 28px;
  margin-bottom: 15px;
  color: var(--text-primary);
  line-height: 1.3;
}

.article-meta-main {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  font-size: 14px;
  color: var(--text-secondary);
  padding-bottom: 10px;
}

.content-section {
  line-height: 1.8;
  font-size: 16px;
  max-height: calc(100vh - 400px);
  overflow-y: auto;
}

.content-box h1,
.content-box h2,
.content-box h3 {
  margin: 20px 0 15px;
  color: var(--text-primary);
}

.content-box p {
  margin-bottom: 15px;
}

.content-box img {
  max-width: 100%;
  width: auto;
  height: auto;
  border-radius: 4px;
  display: block;
  margin: 15px auto;
}

.comment-section {
  margin-top: 20px;
  border-top: 2px solid var(--panel-border);
  padding: 10px 0 0 0;
}

.comment-title {
  font-size: 20px;
  margin-bottom: 20px;
  padding-bottom: 10px;
}

/* 代码块样式（与文章详情对齐） */
.content-box :deep(.code-block),
.code-block {
  position: relative;
  margin: 15px 0;
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

.content-box :deep(.copy-button) {
  position: absolute;
  top: 8px;
  right: 8px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--copy-icon);
  padding: 4px;
  line-height: 0;
}

.content-box :deep(.copy-button:hover) {
  color: var(--copy-icon-hover);
}

.content-box :deep(.copy-icon-check) {
  display: none;
}

.content-box :deep(.copy-button.copied) {
  color: #1a7f37;
}

.content-box :deep(.copy-button.copied .copy-icon:not(.copy-icon-check)) {
  display: none;
}

.content-box :deep(.copy-button.copied .copy-icon-check) {
  display: inline;
}

.content-box :deep(.markdown-image) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 15px auto;
  border-radius: 4px;
  border: 1px solid var(--panel-border);
  box-sizing: border-box;
}

pre {
  background-color: var(--code-bg);
  color: var(--code-text);
  border: 1px solid var(--code-border);
  padding: 15px;
  overflow-x: auto;
}

code {
  font-family: "Consolas", "Monaco", monospace;
  font-size: 14px;
  color: var(--code-text);
}

.copy-button {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: var(--input-bg);
  color: var(--text-primary);
  border: 1px solid var(--input-border);
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-button:hover {
  background-color: var(--hover-bg);
  border-color: #42b983;
}

/* 加载状态样式 */
.loading {
  padding: 20px;
  text-align: center;
  color: var(--text-secondary);
  background-color: var(--card-bg);
  border-radius: 8px;
}

.loading-spinner {
  border-radius: 50%;
  border-top: 3px solid #42b983;
  width: 30px;
  height: 30px;
  animation: spin 1s linear infinite;
  margin: 0 auto 15px;
}

.article-loading {
  min-height: 300px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* 错误提示样式 */
.error-message {
  padding: 20px;
  text-align: center;
  color: #e74c3c;
  background-color: var(--card-bg);
  border-radius: 8px;
}

.article-error {
  min-height: 300px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 15px;
}

.retry-btn {
  padding: 8px 15px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  margin-top: 10px;
}

.retry-btn:hover {
  background-color: #359469;
}

/* 空状态样式 */
.empty-state,
.empty-hint {
  text-align: center;
  padding: 20px;
  color: var(--text-secondary);
  background-color: var(--card-bg);
  border-radius: 8px;
}

.default-btn {
  padding: 8px 15px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
  transition: background-color 0.3s;
}

.default-btn:hover {
  background-color: #2980b9;
}

/* 加载动画 */
@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

/* 响应式适配 */
@media (max-width: 1024px) {
  .wanderland {
    flex-direction: column;
  }

  .slide,
  .content {
    width: 100%;
    max-height: none;
  }

  .slide {
    border-right: none;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  }

  .special-articles {
    flex-direction: row;
    overflow-x: auto;
    padding-bottom: 10px;
  }

  .article-card {
    min-width: 280px;
  }
}

.content-box {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 20px;
  overflow-x: hidden;
  line-height: 1.8;
  /* 增加行高（核心） */
  font-size: 16px;
  /* 优化字体大小 */
  color: var(--text-primary);
}

.content-box p {
  margin-bottom: 20px !important;
  margin-top: 0 !important;
}

.content-box h1,
.content-box h2,
.content-box h3,
.content-box h4 {
  margin-top: 30px !important;
  margin-bottom: 15px !important;
  line-height: 1.5;
}

.content-box ul,
.content-box ol {
  margin-bottom: 20px !important;
  padding-left: 30px !important;
}

.content-box li {
  margin-bottom: 8px !important;
}


.edit-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  opacity: 0.85;
  transition: opacity 0.3s, color 0.2s;
  color: var(--text-primary);
}

.edit-btn i {
  color: inherit;
}

.edit-btn:hover {
  opacity: 1;
  color: var(--text-primary);
}

/* 新增弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  width: 90%;
  max-width: 1000px;
  max-height: 90vh;
  overflow-y: auto;
  background-color: var(--panel-bg);
  color: var(--text-primary);
  border: 1px solid var(--panel-border);
  border-radius: 8px;
  padding: 20px;
}
</style>
