<template>
    <div id="Content" style="height:100vh;">
        <div class="sidebar">
            <!-- 分类列表 -->
            <div class="category" v-for="(category, index) in categories" :key="category.Id"
                @click="scrollToCategory(index)">
                <div class="category-content">
                    <img :src="'../Public/Pictures/' + category.Img" alt="图片丢失了!">
                    <span class="category-name">{{ category.Name }} ({{ getArticleCounts()[category.Name] }})</span>
                </div>
            </div>
        </div>
        <div class="articles">
            <!-- 文章分类 -->
            <div class="article-category" v-for="(category, index) in categories" :key="'category-' + category.Id">
                <h2 :id="'category-' + index">{{ category.Name }}</h2>
                    <!-- 如果当前分类下有文章，则渲染文章列表 -->
                <div class="articles-list" v-if="category.Articles.length > 0">
                    <article class="article" v-for="(article, articleIndex) in category.Articles"
                        :key="'article-' + articleIndex">
                        <div class="article-content">
                            <img @click="getArticleContent(article.Id)" :src="'../Public/Pictures/' + article.Img"
                                alt="文章图片丢失了!">
                            <span class="article-name">{{ article.Title }}</span>
                        </div>
                    </article>
                </div>
                <!-- 如果当前分类下没有文章，则显示“敬请期待” -->
                <div v-else class="no-articles"><h3>————敬请期待————</h3></div>
            </div>
        </div>能增删改查文章、用户权限分级等功能
    </div>
</template>
//用js和vue分别做过能增删改查文章、用户权限分级等功能的博客的前端，用go/gin/gorm做的后端


开发过完整的本地博客，用js和vue分别做过的博客的前端，用go/gin/gorm做的后端

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from 'vue';
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import { useRouter } from 'vue-router';
const router = useRouter();

import { ref, onMounted } from 'vue';

const categories = ref([]);

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

<!-- <style>
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
    /* 添加相对定位 */
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

.articles-list {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
}

.article {
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
    /* Add this */
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
</style> -->
<style>
#Content {
    display: flex;
    height: 100vh;
    /* 设置整个页面的高度 */
}

.sidebar {
    width: 25%;
    /* 或者其他您认为合适的宽度 */
    height: 100vh;
    /* 设置侧边栏的高度 */
    overflow-y: auto;
    /* 允许侧边栏内容滚动 */
    padding: 20px;
    box-sizing: border-box;
    background-color: #f9f9f9;
    /* 可选，根据您的设计调整 */
}

.articles {
    width: 75%;
    /* 或者其他您认为合适的宽度 */
    height: 100vh;
    /* 设置文章列表的高度 */
    overflow-y: auto;
    /* 允许文章列表内容滚动 */
    padding: 20px;
    box-sizing: border-box;
    background-color: #ffffff;
    /* 可选，根据您的设计调整 */
}

.category {
    cursor: pointer;
    margin-bottom: 10px;
    position: relative;
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

.category-content:hover {
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

.articles-list {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
}

.article {
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

.no-articles {
    text-align: center;
    padding: 20px;
    font-size: 1.2em;
    color: #1d1d1d;
}
</style>