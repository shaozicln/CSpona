<template>
  <div class="wanderland">
    <!-- 背景图片 -->
    <img :src="getImageUrl(article.Img)" alt="背景图" class="bg-img" />
    <!-- 左侧分类与文章列表 -->
    <div class="slide">
      <div class="category-container">
        <!-- 错误提示 -->
        <div v-if="categoryError" class="error-message">
          {{ categoryError }}
        </div>

        <!-- 分类列表 -->
        <div v-else>
          <div v-for="category in categories" :key="category.Id" class="category-item">
            <!-- 特殊处理ID为1000的分类 -->
            <div v-if="category.Id === 1000" class="special-articles">
              <div v-for="articleItem in category.Articles" :key="articleItem.Id" class="article-card"
                @click="loadAndShowArticle(articleItem.Id)" :class="{ 'active-article': articleItem.Id === articleId }">
                <img :src="getImageUrl(articleItem.Img)" alt="文章缩略图" class="article-img" loading="lazy" />
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
        </div>
      </div>
    </div>

    <!-- 右侧文章内容区 -->
    <div class="content">
      <!-- 文章错误提示 -->
      <div v-if="articleError && !isLoadingArticle" class="error-message article-error">
        {{ articleError }}
        <button @click="retryLoadArticle" class="retry-btn">重试</button>
      </div>

      <!-- 文章内容展示 -->
      <div v-else-if="article && article.Title" class="article-content">
        <!-- 文章标题区 -->
        <div class="white-box title-section">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            
          <h1 class="article-main-title">{{ article.Title }}</h1>
            <!-- 使用Font Awesome的编辑图标 -->
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
          <div class="content-box" v-html="renderMarkdown(article.Content)"></div>
        </div>

        <!-- 评论区 -->
        <div class="white-box comment-section" v-if="!isLoadingArticle">
          <h3 class="comment-title" style="text-align: center; font-size: 40px; font-family: cursive">
            评论区
          </h3>
          <!-- 使用更唯一的key强制评论组件在文章ID变化时重新渲染 -->
          <Comments ref="commentsRef" :article-id="articleId" :key="'comments-' + articleId" />
        </div>
      </div>

      <!-- 文章未找到 -->
      <div v-else-if="!isLoadingArticle" class="empty-state">
        <p>未找到指定文章</p>
        <button @click="goToDefaultArticle" class="default-btn">
          查看默认文章
        </button>
      </div>
    </div>
  </div>

  <div v-if="isLoadingArticle" class="debug">
    加载状态: {{ isLoadingArticle }} | 文章ID: {{ articleId }}
  </div>

   <!-- 编辑弹窗遮罩层 -->
  <div v-if="showEditModal" class="modal-overlay" @click="showEditModal = false">
    <!-- 编辑弹窗内容 - 阻止事件冒泡 -->
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
import { useRouter, useRoute, onBeforeRouteLeave } from "vue-router";
import { getCurrentInstance } from "vue";
import { marked } from "marked";
import hljs from "highlight.js";
import "highlight.js/styles/github.css";
import { useArticleStore } from "@/stores/article";
import { decodeArticleId } from "@/utils/utils.js";

import { library } from '@fortawesome/fontawesome-svg-core';
import { faEdit } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
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
const categories = ref([]);
const users = ref({}); // 缓存用户信息
const toc = ref([]);
const headingCounter = {};
const activeSlug = ref("");
const isLoadingCategories = ref(true); // 分类加载状态
const isLoadingArticle = ref(false);
const categoryError = ref(null);
const articleError = ref(null);
const isArticleOwner = ref(false);

// 获取全局实例与URL
const instance = getCurrentInstance();
const { proxy } = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

// 初始化文章ID
const initArticleId = () => {
  if (!route.params.articleId) {
    articleId.value = 95; // 默认值
    console.log("使用默认文章ID:", articleId.value);
    return;
  }

  try {
    articleId.value = decodeArticleId(route.params.articleId);
    isLoadingArticle.value = false;
    console.log("从路由参数获取文章ID:", articleId.value);
  } catch (e) {
    console.error("解码文章ID失败:", e);
    articleId.value = 95; // 回退默认值
  }
};

const renderMarkdown = (raw) => {
  if (typeof raw !== "string") {
    raw = String(raw || "");
  }
  try {
    const htmlContent = marked(raw);
    console.log(htmlContent);

    return htmlContent;
  } catch (e) {
    console.error("Markdown 渲染失败", e);
    return `<p>渲染失败</p>`;
  }
};

const fetchArticleDetail = async (id) => {
  const res = await fetch(`${URL}/path-to-article/${id}`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  const { data } = await res.json();

  article.value = {
    ...data,
    Content: data.Content ?? data.content ?? "",
    Title: data.Title ?? data.title ?? "未命名文章",
  };

  // 判断文章作者ID与当前用户ID是否相等
  isArticleOwner.value = String(article.value.UserId) === String(userId.value); // 转为字符串避免类型问题

  if (data.UserId) await fetchUserInfo(data.UserId);
};

// 加载并显示文章
const loadAndShowArticle = async (id) => {
  isLoadingArticle.value = true;
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
    console.log("文章加载完成", article.value);
    nextTick(() => {
      if (commentsRef.value) {
        commentsRef.value.fetchComments();
      }
    });
    window.scrollTo({ top: 0, behavior: "smooth" });
    nextTick(() => {
      hljs.highlightAll();
      bindCopyButtons();
    });
  } catch (err) {
    console.error("加载文章失败:", err);
    articleError.value = "加载文章失败，请稍后重试";
  } finally {
    isLoadingArticle.value = false;
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
      loadAndShowArticle(95); // 回退到默认文章
    }
  },
  { immediate: true }
);

// 配置 marked
marked.setOptions({
  gfm: true,
  breaks: true,
  renderer: new marked.Renderer(),
});

// 使用自定义渲染器来处理代码块
const renderer = new marked.Renderer();
renderer.code = (code, language) => {
  const validLanguage =
    !language || !hljs.getLanguage(language) ? "plaintext" : language;
  return `<pre class="hljs"><code class="hljs language-${validLanguage}">${hljs.highlightAuto(code).value
    }</code></pre>`;
};
// 如果是本地文件路径，添加前缀
renderer.image = function () {
  // 参数解析
  let href, title, text;

  if (arguments.length >= 3) {
    [href, title, text] = arguments;
  } else if (arguments[0] && typeof arguments[0] === "object") {
    const token = arguments[0];
    href = token.href;
    title = token.title;
    text = token.text;
  } else {
    console.error("无法解析图片参数:", arguments);
    href = "";
  }

  // 确保 href 是字符串
  if (typeof href !== "string") {
    href = String(href);
  }

  // 编码 URL 并创建图片标签
  return `<img src="${encodeURI(href)}" 
               alt="${(text || "").replace(/"/g, "&quot;")}" 
               title="${(title || "").replace(/"/g, "&quot;")}"
               class="markdown-image"
               style="
                 max-width: 100%;
                 height: auto;
                 display: block;
                 margin: 15px auto;
                 border-radius: 4px;
                 background: #f8f8f8;
                 border: 1px solid #eee;
                 padding: 4px;
                 box-sizing: border-box;
               ">`;
};

marked.setOptions({
  renderer: renderer,
});
// 获取图片URL
const getImageUrl = (imgName) => {
  if (!imgName) return `${proxy.$imageBaseUrl}default-article.jpg`;
  return `${proxy.$imageBaseUrl}${imgName}`;
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

// 获取分类数据
const fetchCategories = async () => {
  try {
    isLoadingCategories.value = true;
    const response = await fetch(`${URL}/categories-with-articles`);
    if (!response.ok) throw new Error(`HTTP错误: ${response.status}`);

    const data = await response.json();
    if (!data || !data.data) throw new Error("数据格式不正确");

    categories.value = data.data || [];

    // 预加载用户信息
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
    isLoadingCategories.value = false;
  }
};

// 获取文章详情

// 获取用户信息
const fetchUserInfo = async (userId) => {
  try {
    if (users.value[userId]) return;

    const response = await fetch(`${URL}/users/${userId}`);
    if (!response.ok) throw new Error(`获取用户 ${userId} 信息失败`);

    const data = await response.json();
    users.value[userId] = data.data;
  } catch (err) {
    console.error("获取用户信息失败:", err);
    // 使用默认用户信息
    users.value[userId] = {
      Username: "未知作者",
      Avatar: "default-avatar.jpg",
    };
  }
};

// 获取作者名称
const getAuthorName = (userId) => {
  return users.value[userId]?.Username || "未知作者";
};

// 复制按钮绑定
const bindCopyButtons = () => {
  const copyButtons = document.querySelectorAll(".copy-button");
  copyButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const codeBlock = button.parentElement.querySelector("code");
      const codeText = codeBlock.innerText;

      navigator.clipboard
        .writeText(codeText)
        .then(() => {
          button.textContent = "已复制";
          button.classList.add("copied");
          setTimeout(() => {
            button.textContent = "复制";
            button.classList.remove("copied");
          }, 2000);
        })
        .catch((err) => {
          console.error("复制失败:", err);
          button.textContent = "复制失败";
          setTimeout(() => {
            button.textContent = "复制";
          }, 2000);
        });
    });
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
    loadArticle(articleId.value);
    loadAndShowArticle(articleId.value);
  } else {
    goToDefaultArticle();
  }
};

// 跳转到默认文章
const goToDefaultArticle = () => {
  const encodedId = window.btoa(95);
  router.push(`/article/${encodedId}`);
};

// 初始化
onMounted(async () => {
  console.log(`组件挂载开始`);

  // 初始化文章ID
  initArticleId();
  console.log(`初始化文章ID: ${articleId.value}`);

  // 只加载分类，文章加载由路由监听处理
  await fetchCategories();
  console.log(`分类加载完成`);

  const defaultId = 95;
  await loadAndShowArticle(defaultId);

  // 绑定事件监听
  window.addEventListener("scroll", handleScroll);
  console.log(`滚动监听已添加`);
});

// 清理工作
onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
});

onBeforeRouteLeave((to, from) => {
  if (to.name === "Articles") {
    // 仅当跳转到 Articles 路由时设置刷新标记
    sessionStorage.setItem("refreshAfterEnter", "Articles");
  }
});
</script>

<style scoped>
/* 基础布局 */
.wanderland {
  display: flex;
  width: 100%;
  min-height: 100vh;
  gap: 10px;
}

/* 左侧区域 */
.slide {
  width: 25%;
  padding: 20px;
  box-sizing: border-box;
  overflow-y: auto;
  max-height: 100vh;
  /* 细边框分隔 */
  background-color: rgba(255, 255, 255, 0.95);
  height: 100vh;
}

/* 右侧内容区 */
.content {
  width: 75%;
  height: 100vh;
  background-color: rgba(255, 255, 255, 0.8);
  position: relative;
  box-sizing: border-box;
  overflow-y: auto;
}

/* 背景图 */
.bg-img {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  object-fit: cover;
  z-index: -1;
  /* filter: brightness(0.8); */
}

.article-count {
  color: #666;
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
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
}

.active-article {
  border: 2px solid #42b983;
  box-shadow: 0 2px 12px rgba(66, 185, 131, 0.3);
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
  color: #666;
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
  color: #333;
  line-height: 1.3;
}

.article-meta-main {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  font-size: 14px;
  color: #666;
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
  color: #2c3e50;
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
  border-top: 2px solid #0000002a;
  padding: 10px 0 0 0;
}

.comment-title {
  font-size: 20px;
  margin-bottom: 20px;
  padding-bottom: 10px;
}

/* 代码块样式 */
.code-block {
  position: relative;
  margin: 15px 0;
  border-radius: 4px;
  overflow: hidden;
}

pre {
  background-color: #f8f9fa;
  padding: 15px;
  overflow-x: auto;
}

code {
  font-family: "Consolas", "Monaco", monospace;
  font-size: 14px;
}

.copy-button {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: rgba(255, 255, 255, 0.8);
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-button:hover {
  background-color: #fff;
  border-color: #42b983;
}

/* 加载状态样式 */
.loading {
  padding: 20px;
  text-align: center;
  color: #666;
  background-color: rgba(255, 255, 255, 0.9);
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
  background-color: rgba(255, 255, 255, 0.9);
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
  color: #666;
  background-color: rgba(255, 255, 255, 0.8);
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
  opacity: 0.7;
  transition: opacity 0.3s;
}

.edit-btn:hover {
  opacity: 1;
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
  background-color: white;
  border-radius: 8px;
  padding: 20px;
}
</style>
