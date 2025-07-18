<template>
  <div class="wanderland">
    <!-- 背景图片 -->
    <img :src="getImageUrl(article.Img)" alt="背景图" class="bg-img" />
    
    <!-- 左侧分类与文章列表 -->
    <div class="slide">
      <div class="category-container">
        <!-- 分类加载错误提示 -->
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
      <!-- 加载中遮罩 -->
      <div v-if="isLoadingArticle" class="loading-overlay">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <!-- 文章错误提示 -->
      <div v-else-if="articleError" class="error-message article-error">
        {{ articleError }}
        <button @click="retryLoadArticle" class="retry-btn">重试</button>
      </div>

      <!-- 文章内容展示 -->
      <div v-else-if="article && article.Title" class="article-content">
        <!-- 文章标题区 -->
        <div class="white-box title-section">
          <h1 class="article-main-title">{{ article.Title }}</h1>
          <div class="article-meta-main">
            <span>作者: {{ getAuthorName(article.UserId) }}</span>
            <span>发布时间: {{ TimeFormat(article.CreatedAt) }}</span>
            <span>浏览量: {{ article.ViewCount || 0 }}</span>
            <span>评论数: {{ article.CommentCount || 0 }}</span>
          </div>
          <div class="content-box" v-html="renderedContent"></div>
        </div>

        <!-- 评论区 -->
        <div class="white-box comment-section" v-if="articleLoaded">
          <h3 class="comment-title">评论区</h3>
          <Comments ref="commentsRef" :article-id="articleId" :key="'comments-' + articleId" />
        </div>
      </div>

      <!-- 文章未找到 -->
      <div v-else class="empty-state">
        <p>未找到指定文章</p>
        <button @click="goToDefaultArticle" class="default-btn">
          查看默认文章
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import Comments from "../Article/Comments.vue";
import { useUserStore } from "@/stores/user";
import { ref, onMounted, computed, watch, nextTick, onUnmounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { getCurrentInstance } from "vue";
import { marked } from "marked";
import hljs from "highlight.js";
import "highlight.js/styles/github.css";
import { useArticleStore } from "@/stores/article";
import { decodeArticleId } from '@/utils/utils.js';

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
const isLoadingCategories = ref(true);
const isLoadingArticle = ref(false);
const categoryError = ref(null);
const articleError = ref(null);
const articleLoaded = ref(false);

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
    console.log("从路由参数获取文章ID:", articleId.value);
  } catch (e) {
    console.error("解码文章ID失败:", e);
    articleId.value = 95; // 回退默认值
  }
};

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
const fetchArticleDetail = async (id) => {
  try {
    console.log(`开始获取文章详情，ID: ${id}`);
    const response = await fetch(`${URL}/path-to-article/${id}`);
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP错误: ${response.status}, 响应: ${errorText}`);
    }

    const data = await response.json();
    console.log("原始文章数据:", data);

    if (!data || !data.data) {
      throw new Error("文章数据为空");
    }

    // 确保Content是字符串
    if (typeof data.data.Content !== 'string') {
      console.warn('文章内容非字符串类型，转为JSON字符串');
      data.data.Content = JSON.stringify(data.data.Content || {});
    }
    
    article.value = data.data;
    
    // 获取文章作者信息
    if (article.value.UserId) {
      await fetchUserInfo(article.value.UserId);
    }
  } catch (err) {
    console.error("获取文章失败:", err);
    throw err;
  }
};

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

// 加载并显示文章
const loadAndShowArticle = async (id) => {
  console.log(`开始加载文章 ID: ${id}, 当前文章ID: ${articleId.value}`);

  // 防止重复加载相同文章
  if (!id || id === articleId.value) {
    console.log(`跳过加载，ID相同或无效`);
    return;
  }

  // 设置加载状态和清除旧数据
  isLoadingArticle.value = true;
  articleError.value = null;
  article.value = {};
  articleLoaded.value = false;
  console.log(`清除旧数据，设置加载状态: ${isLoadingArticle.value}`);

  try {
    // 更新文章ID状态
    articleId.value = id;
    articleStore.setArticleId(id);
    localStorage.setItem("articleId", id);
    console.log(`更新文章ID: ${id}`);

    // 更新URL
    const encodedId = window.btoa(id.toString());
    router.replace({
      name: route.name,
      params: { ...route.params, articleId: encodedId },
    });
    console.log(`更新路由参数: ${encodedId}`);

    // 加载文章详情
    await fetchArticleDetail(id);
    
    // 强制DOM更新后再执行清理和高亮
    await nextTick();
    
    // 彻底清理DOM中的无效内容
    cleanupDOM();
    
    // 应用代码高亮
    hljs.highlightAll();
    bindCopyButtons();
    
    // 标记文章加载完成
    articleLoaded.value = true;
    
    // 延迟加载评论组件，避免干扰主内容
    setTimeout(() => {
      if (commentsRef.value) {
        console.log(`开始加载评论`);
        commentsRef.value.fetchComments();
      }
    }, 300);

    // 滚动到顶部
    window.scrollTo({ top: 0, behavior: "smooth" });
  } catch (err) {
    console.error("加载文章失败:", err);
    articleError.value = "加载文章失败，请稍后重试";
  } finally {
    isLoadingArticle.value = false;
    console.log(`加载完成，状态: ${isLoadingArticle.value}`);
  }
};

// 安全渲染Markdown内容
const renderMarkdownSafely = (content) => {
  if (!content) return '<div class="empty-content">暂无文章内容</div>';
  
  try {
    // 深度转换嵌套对象为字符串
    const safeContent = typeof content === 'string'
      ? content
      : JSON.stringify(content, (key, value) => {
          // 处理特殊值
          if (value === null || value === undefined) return '';
          return value;
        }, 2);
    
    // 替换特殊字符串
    const cleanContent = safeContent
      .replace(/\[object Object\]/g, '')
      .replace(/null/g, '')
      .replace(/undefined/g, '');
    
    // 渲染Markdown
    return marked.parse(cleanContent);
  } catch (error) {
    console.error('Markdown解析失败:', error);
    return `<div class="error">解析内容出错: ${error.message}</div>`;
  }
};

// 渲染文章内容
const renderedContent = computed(() => {
  if (isLoadingArticle.value) {
    return `<div class="loading">加载中...</div>`;
  }

  if (articleError.value) {
    return `<div class="error">${articleError.value}</div>`;
  }

  return renderMarkdownSafely(article.value?.Content);
});

// Markdown渲染配置
const renderer = {
  heading(text, level) {
    // 确保text是字符串
    const safeText = typeof text === 'string' 
      ? text 
      : JSON.stringify(text, null, 2);
    
    // 生成安全的slug
    const slug = safeText
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-')
      .trim()
      .slice(0, 50); // 限制长度
    
    return `<h${level} id="${slug}">${safeText}</h${level}>`;
  },
  
  code(code, language) {
    return `<div class="code-block">
              <pre><code class="language-${language}">${hljs.highlightAuto(code).value}</code></pre>
              <button class="copy-button">复制</button>
            </div>`;
  },
  
  paragraph(text) {
    return `<p>${text}</p>`;
  },
  
  image(href, title, text) {
    return `<img src="${href}" alt="${text}" title="${title || text}" class="article-image">`;
  },
  
  blockquote(quote) {
    return `<blockquote>${quote}</blockquote>`;
  },
  
  html(html) {
    return html;
  },
  
  br() {
    return "<br>";
  }
};

// 配置marked
marked.setOptions({
  gfm: true,
  breaks: true,
  renderer: renderer,
});

// 获取作者名称
const getAuthorName = (userId) => {
  return users.value[userId]?.Username || "未知作者";
};

// 获取图片URL
const getImageUrl = (imgName) => {
  if (!imgName) return `${proxy.$imageBaseUrl}default-article.jpg`;
  return `${proxy.$imageBaseUrl}${imgName}`;
};

// 时间格式化
const TimeFormat = (time) => {
  if (!time) return "";
  const date = new Date(time);
  return date.toLocaleString('zh-CN');
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

// 增强DOM清理
const cleanupDOM = () => {
  // 清理主内容区域
  const contentElement = document.querySelector('.content-box');
  if (!contentElement) return;
  
  // 清理空标签
  contentElement.querySelectorAll('*').forEach(el => {
    if (el.textContent.trim() === '' && !el.querySelector('img, code, pre')) {
      el.remove();
    }
  });
  
  // 清理包含无效内容的标签
  contentElement.querySelectorAll('p, h1, h2, h3').forEach(el => {
    if (el.textContent === 'null' || el.textContent === 'undefined') {
      el.remove();
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

  // 加载分类
  await fetchCategories();
  console.log(`分类加载完成`);

  // 加载文章
  if (articleId.value) {
    loadAndShowArticle(articleId.value);
  }
});

// 清理工作
onUnmounted(() => {
  // 移除所有事件监听器
  document.querySelectorAll('.copy-button').forEach(btn => {
    btn.removeEventListener('click', null);
  });
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
  border-radius: 4px;
  margin: 15px 0;
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
</style>
