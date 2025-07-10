<template>
    <div id="Content" style="height:100vh;">
        <div class="sidebar">
            <!-- 分类列表 -->
            <div class="category" v-for="(category, index) in categories" :key="category.Id"
                @click="scrollToCategory(index)">
                <div class="category-content">
                    <img :src="getImageUrl(category.Img)" alt="图片丢失了!">
                    <span class="category-name">{{ category.Name }} ({{ getArticleCounts()[category.Name] }})</span>
                </div>
            </div>
        </div>
        <div class="articles">
            <!-- 文章分类 -->
            <div class="article-category" v-for="(category, index) in categories" :key="'category-' + category.Id">
                <h3 :id="'category-' + index">{{ category.Name }}</h3>
                <!-- 如果当前分类下有文章，则渲染文章列表 -->
                <div class="articles-list" v-if="category.Articles.length > 0">
                    <article class="article" v-for="(article, articleIndex) in category.Articles"
                        :key="'article-' + articleIndex">
                        <div class="article-content">
                            <img @click="getArticleContent(article.Id)" :src="getImageUrl(article.Img)"
                                alt="文章图片丢失了!">
                            <span class="article-name">{{ article.Title }}</span>
                        </div>
                    </article>
                </div>
                <!-- 如果当前分类下没有文章，则显示“敬请期待” -->
                <div v-else class="no-articles">
                    <h3>————敬请期待————</h3>
                </div>
            </div>
        </div>
    </div>
</template>

<!-- 增加算法类和实习经验类 -->

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from 'vue';
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import { useRouter } from 'vue-router';
const router = useRouter();

import { ref, onMounted } from 'vue';

const categories = ref([]);


const { proxy } = getCurrentInstance()
const getImageUrl = (imgName) => {
  return `${proxy.$imageBaseUrl}${imgName}`
}
// 获取分类和文章数据
const fetchCategories = async () => {
    try {
        const response = await fetch(`${URL}/categories-with-articles`);
        const data = await response.json();
        categories.value = data.data;

        // 文章数量统计
        const articleCounts = {};
        data.data.forEach(category => {
            articleCounts[category.Name] = category.Articles.length;
        });
        localStorage.setItem('articleCounts', JSON.stringify(articleCounts));

        console.log(data.data);

    } catch (error) {
        console.error('Failed to fetch categories:', error);
    }
};

const getArticleCounts = () => {
    const articleCounts = localStorage.getItem('articleCounts');
    return articleCounts ? JSON.parse(articleCounts) : {};
};

const scrollToCategory = (index) => {
    const categoryElement = document.getElementById(`category-${index}`);
    if (categoryElement) {
        categoryElement.scrollIntoView({ behavior: 'smooth', block: 'start' });
    } else {
        console.warn(`Category element with id 'category-${index}' not found.`);
    }
};

// 跳转到文章页面
const getArticleContent = (articleId) => {
    localStorage.setItem("articleId", articleId);
    router.push({ path: '/articleContent' });
};

onMounted(() => {
    fetchCategories();
});

</script>

<style>
#Content {
    display: flex;
    height: 100vh;
}

.sidebar {
    width: 25%;
    height: 100vh;
    overflow-y: auto;
    padding: 20px;
    box-sizing: border-box;
}

.articles {
    width: 75%;
    height: 100vh;
    overflow-y: auto;
    padding: 20px;
    box-sizing: border-box;
}


.category {
    cursor: pointer;
    margin-bottom: 10px;
    position: relative;
    /* animation: slide-fade-in both;
    /* 使用浏览器的视图时间线,允许动画根据视口的变化进行同步 */
    /*animation-timeline: view();*/
    /* 定义动画时间范围 */
    /*animation-range: contain 0% contain 50%;*/
}

.category img {
    width: 100%;
    height: 18vh;
    display: block;
    margin-bottom: 5px;
}

.category-content {
    position: relative;
    transition: all 0.3s ease-in-out;
}

.category-content {
    transform: translateY(-5px);
    box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}


.category-name {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    color: white;
    text-align: center;
    padding: 5px 0;
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
        transform: scale(.8) translateY(25px);
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
    border: 1px solid black;
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
    color: #1d1d1d;
}
</style>
