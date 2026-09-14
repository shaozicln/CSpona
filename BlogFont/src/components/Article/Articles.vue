<template>
  <div id="Content" style="height: 100vh">
    <div class="sidebar">
      <!-- 顶部圆形图片及随机句子区域 -->
      <div class="sidebar-header">
        <img class="header-avatar" :src="getImageUrl('Venti-2.webp')" alt="头像" />
        <p class="daily-quote">{{ currentQuote }}</p>
      </div>

      <!-- 使用计算属性 filteredCategories 渲染列表 -->
      <div class="categories-list">
        <PageLoading
          v-if="showListLoading"
          variant="cards"
          message="分类加载中…"
        />
        <PageLoading
          v-else-if="loadError"
          error="无法加载分类，请稍后重试"
          @retry="fetchCategories"
        />
        <template v-else>
          <div class="category-item" v-for="(category, index) in filteredCategories" :key="category.Id"
            @click="scrollToCategory(index)">
            <img class="category-icon" :src="getImageUrl(category.Img, category.ImgUrl)" alt="分类图标" />
            <span class="category-name">
              {{ category.Name }} ({{ getArticleCounts()[category.Name] || 0 }})
            </span>
          </div>
        </template>
      </div>
    </div>
    <div class="articles">
      <PageLoading
        v-if="showListLoading"
        variant="list"
        message="文章加载中…"
      />
      <PageLoading
        v-else-if="loadError"
        error="无法加载文章列表，请稍后重试"
        @retry="fetchCategories"
      />
      <template v-else>
        <!-- 文章分类区域（保持不变） -->
        <div class="article-category" v-for="(category, index) in categories" :key="'category-' + category.Id">
          <h3 v-if="category && category.Id !== 1000" :id="'category-' + index">
            {{ category.Name }}
          </h3>
          <div class="articles-list" v-if="category.Articles?.length > 0">
            <article
              class="article"
              v-for="(article, articleIndex) in category.Articles"
              :key="'article-' + article.Id + '-' + articleIndex"
            >
              <div class="article-content">
                <img @click="getArticleContent(article.Id)" :src="getImageUrl(article.Img, article.ImgUrl)" alt="文章图片丢失了!" />
                <span class="article-name">{{ article.Title }}</span>
              </div>
            </article>
          </div>
          <div v-else-if="category && category.Id !== 1000" class="no-articles">
            <h3>————敬请期待————</h3>
          </div>
        </div>
      </template>
    </div>
  </div>

</template>

<!--TODO： 增加算法类和实习经验类 -->

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import { useRouter, useRoute } from "vue-router";
const route = useRoute();
const router = useRouter();

import { ref, onMounted, computed } from "vue";
import { resolveImageUrl } from "@/utils/image.js";
import PageLoading from "@/components/common/PageLoading.vue";
import { useDelayedLoading } from "@/composables/useDelayedLoading.js";

const { pending: listPending, visible: listVisible, start: startListLoading, stop: stopListLoading } =
  useDelayedLoading({ delayMs: 0, minShowMs: 320 });
const loadError = ref("");
const showListLoading = computed(() => listPending.value || listVisible.value);

// 每日随机句子数组
const dailyQuotes = [
  "莫待花谢空折枝",
  "过我嶙峋，拥我九春",
  "蝴蝶飞来的夜晚绝非偶然",
  "错过本身就是惊艳的组成部分",
  "我喻我以长青，我拥此春待亭亭",
  "Everything wins.",
  "Die Luft der Freiheit weht.",
  "放轻松，就当漫游地球",
   "快乐与悲伤都是太沉重的字眼，我宁愿说今天跳了支不错的舞",
  "人生如逆旅，我亦是行人",
  "且将新火试新茶，诗酒趁年华",
  "一点浩然气，千里快哉风",
  "自歌自舞自开怀且喜无拘无碍",
];
const getImageUrl = (imgName, resolved) => {
  return resolveImageUrl(resolved || imgName);
};

const categories = ref([]);

// 获取分类和文章数据（BFF 聚合）
const fetchCategories = async () => {
  loadError.value = "";
  startListLoading();
  try {
    const response = await fetch(`${URL}/page/home`);
    const data = await response.json();
    const list = data?.data?.categories || data?.data || [];

    list.forEach((category) => {
      if (category.Articles && category.Articles.length > 0) {
        category.Articles.sort((a, b) => {
          return new Date(b.CreatedAt) - new Date(a.CreatedAt);
        });
      }
    });

    categories.value = list;

    const articleCounts = {};
    list.forEach((category) => {
      articleCounts[category.Name] = (category.Articles || []).length;
    });
    localStorage.setItem("articleCounts", JSON.stringify(articleCounts));
  } catch (error) {
    console.error("Failed to fetch categories:", error);
    loadError.value = "load_failed";
    categories.value = [];
  } finally {
    await stopListLoading();
  }
};
// 过滤掉 Id === 1000 的分类
const filteredCategories = computed(() => {
  return categories.value.filter(
    (category) => category && category.Id !== 1000
  );
});
const getArticleCounts = () => {
  const articleCounts = localStorage.getItem("articleCounts");
  return articleCounts ? JSON.parse(articleCounts) : {};
};

const scrollToCategory = (index) => {
  const categoryElement = document.getElementById(`category-${index}`);
  if (categoryElement) {
    categoryElement.scrollIntoView({ behavior: "smooth", block: "start" });
  } else {
    console.warn(`Category element with id 'category-${index}' not found.`);
  }
};

// 跳转到文章页面
const getArticleContent = (articleId) => {
  const encodedId = window.btoa(articleId.toString());
  router.push(`/articleContent/${encodedId}`);
};

const currentQuote = ref("");
onMounted(() => {
  fetchCategories();
  // 根据日期计算随机句子索引（确保每天一句）
  const day = new Date().getDate();
  const randomIndex = day % dailyQuotes.length;
  currentQuote.value = dailyQuotes[randomIndex];
});
</script>

<style scoped>
#Content {
  display: flex;
  height: 100vh;
  flex-direction: row;
}

h2 {
  text-align: center;
  font-size: 30px;
  font-family: cursive;
  color: var(--text-primary);
}

h3 {
  text-align: center;
  font-size: 40px;
  font-family: cursive;
  color: var(--text-primary);
}

.sidebar {
  width: 27%;
  height: 100vh;
  overflow-y: auto;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 30px;
  /* 顶部区域与分类列表间距 */
}

/* 顶部头像和句子区域 */
.sidebar-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  padding: 20px;
  border-bottom: 1px solid var(--panel-border);
  background-color: var(--panel-bg);
  border-radius: 20px;
}

.header-avatar {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  /* 圆形图片 */
  object-fit: cover;
  border: 3px solid #f0f0f0;
}

.daily-quote {
  font-family: Cormorant SC, serif;
  font-size: 20px;
  color: var(--text-primary);
  text-align: center;
  line-height: 1.5;
  max-width: 90%;
  /* background-color:rgb(255, 255, 255, 0.5);
  border: 1px solid white; */
  border-radius: 20px;
  transition: all 0.3s ease;
}

.articles {
  width: 75%;
  height: 100vh;
  overflow-y: auto;
  padding: 20px;
  box-sizing: border-box;
}

/* 分类列表容器 */
.categories-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 单个分类项样式 */
.category-item {
  display:flex;
  justify-content: center;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 15px;
  background-color: var(--panel-bg);
  border: 1px solid var(--panel-border);
  border-radius: 15px;
  transition: all 0.3s ease;
  color: var(--text-primary);
}

.category-item:hover {
  background-color: var(--hover-bg);
  transform: translateX(5px);
}

.category-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  /* 圆形图标 */
  object-fit: cover;
}

.category-name {
  font-size: 16px;
  color: var(--text-primary);
  white-space: nowrap;
  /* 防止文字换行 */
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-category {
  margin-bottom: 40px;
}

.article-category h2 {
  margin-top: 0;
}

.article-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

@keyframes slide-fade-in {
  from {
    opacity: 0;
    box-shadow: none;
    transform: scale(0.8) translateY(25px);
  }
}

.articles-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

.article {
  animation: slide-fade-in both;
  /* 使用浏览器的视图时间线,允许动画根据视口的变化进行同步 */
  animation-timeline: view();
  /* 定义动画时间范围 */
  animation-range: contain 0% contain 50%;
  width: 45%;
  margin: 10px;
}

.article-content {
  position: relative;
  border: 1px solid var(--panel-border);
  padding: 10px;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease-in-out;
}

.article-content:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}

.article-content img {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 10px;
  margin-bottom: auto;
}

.article-name {
  position: absolute;
  bottom: 10px;
  left: 10px;
  background-color: rgba(0, 0, 0, 0.5);
  color: white;
  padding: 5px 10px;
  border-radius: 10px;
}

/* #Aimg {
    width: 300px;
    height: 300px;
} */

/* 确保滚动条在视高固定的情况下不会导致内容区域溢出 */
.Content,
.articles {
  overflow: auto;
  mask-image: linear-gradient(to bottom, white, black 20px, black 100%);
}

.no-articles {
  text-align: center;
  padding: 20px;
  font-size: 1.2em;
  color: var(--text-primary);
}
</style>
