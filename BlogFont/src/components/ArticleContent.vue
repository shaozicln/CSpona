<template>
    <div class="container">
        <img :src="'../Public/Pictures/' + article.Img" alt="" class="bg-img" />
        <div class="grid-container">
            <div class="left-side">
                <Author />
            </div>
            <div class="right-side">
                <div class="white-box">
                    <div class="title-box">
                        <h1>{{ article.Title }}</h1>
                        <p>
                            浏览量：{{ article.ViewCount }} |
                            创建时间：{{ TimeFormat(article.CreatedAt) }} |
                            评论数：{{ article.CommentCount }}
                        </p>
                    </div>
                    <div class="content-box" v-html=renderedContent></div>
                    <div class="comment-box" style="height: 50vh;"></div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import Author from './Author.vue'
import { marked } from 'marked';

import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';

const articleId = +localStorage.getItem('articleId');

const article = ref('')

const route = useRoute();

onMounted(
    async () => {
        try {
            const response = await fetch('http://127.0.0.1:8081/path-to-article/' + articleId);
            const data = await response.json();
            article.value = data.data;
            console.log(data.data);

        } catch (error) {
            console.error('Failed to fetch categories:', error);
        }
    })
const renderedContent = computed(() => {
    if (!article.value) return '';
    return marked(article.value.Content);
});
function TimeFormat(time) {
    const date = new Date(time); // <--- Pass the `time` parameter here
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hour = date.getHours();
    const minute = date.getMinutes();
    const second = date.getSeconds()
    return `${year}/${month}/${day} - ${hour}:${minute}:${second}`;
}
</script>

<style scoped>
.container {
    position: relative;
    height: 100vh;
    overflow: hidden;
}

.bg-img {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100vh;
    object-fit: cover;
    z-index: -1;
}

.grid-container {
    display: grid;
    grid-template-columns: 1fr 3fr;
    grid-gap: 20px;
    padding: 20px;
}

.left-side {
    grid-column: 1;
}

.right-side {
    grid-column: 2;
}

.white-box {
    background-color: rgba(255, 255, 255, 0.8);
    padding: 20px;
    height: 95vh;
    overflow-y: auto;
}

.title-box {
    margin-bottom: 20px;
}

.title-box h1 {
    font-size: 40px;
    margin-bottom: 10px;
}

.title-box p {
    font-size: 15px;
    color: #666;
}

.content-box {
    margin-bottom: 20px;
}

.comment-box {
    margin-bottom: 20px;
}
</style>