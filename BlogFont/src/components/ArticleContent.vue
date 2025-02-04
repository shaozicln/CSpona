<template>
    <div class="container">
        <img :src="'../Public/Pictures/' + article.Img" alt="" class="bg-img" />
        <div class="grid-container">
            <div class="left-side">
                <transition name="fade">
                    <!--  v-if="isAuthorVisible" -->
                    <div class="author-container">
                        <AuthorBack />
                    </div>
                </transition>
                <!-- <transition name="avatar">
                    <div class="avatar-container"  v-if="isAvatarVisible" @click="toggleAuthor">
                        <img :src="'../Public/Pictures/' + authorAvatar" alt="Author Avatar" class="avatar" />
                    </div>
                </transition> -->
                <transition name="toc-transition">
                    <div class="toc-box">
                        <div id="toc-title">内容索引</div>
                        <ul>
                            <li v-for="item in toc" :key="item.slug"
                                :style="{ marginLeft: (item.level - 1) * 20 + 'px' }"
                                @click="scrollToContent(item.slug)">
                                <a :href="'#' + item.slug">{{ item.text }}</a>
                            </li>
                        </ul>
                    </div>
                </transition>
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
                    <div class="content-box" v-html="renderedContent"></div>
                    <div class="comment-box" style="height: 50vh;"></div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from 'vue';
import AuthorBack from './AuthorBack.vue'
import { marked } from 'marked';
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css'; // 使用 GitHub 风格的代码高亮样式
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { nextTick } from 'vue';

const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

marked.setOptions({
    gfm: true, // 启用 GitHub Flavored Markdown
    breaks: true, // 将换行符转换为 <br>
});

const articleId = +localStorage.getItem('articleId');
const article = ref('');
const route = useRoute();
const toc = ref([]);
const isAuthorVisible = ref(true); // 控制 Author 组件的显示和隐藏
const isAvatarVisible = ref(false);
const authorAvatar = ref('35.jpg');

const renderer = {
    heading(text, level) {
        // 使用原始标题作为 slug，并进行处理
        const slug = text.replace(/[^\w\u4e00-\u9fa5-]+/g, '-'); // 只保留字母、数字、连字符和汉字
        toc.value.push({ slug, text, level });
        console.log(`Generated slug for "${text}": ${slug}`); // 添加日志
        return `<h${level} id="${slug}">${text}</h${level}>`;
    },
    code(code, language) {
        console.log(`Generated slug for "${text}": ${slug}`); // 添加日志
        return `<pre><code class="language-${language}">${code}</code></pre>`;
    }
};

// 确保在解析 Markdown 之前应用自定义渲染器
marked.use(renderer);

const renderedContent = computed(() => {
    if (!article.value) return '';
    toc.value = [];

    // 使用 marked 解析 Markdown 内容
    const content = marked(article.value.Content);

    // 提取 <h1> 到 <h6> 标题
    const parser = new DOMParser();
    const doc = parser.parseFromString(content, 'text/html');
    const headings = doc.querySelectorAll('h1, h2, h3, h4, h5, h6');

    headings.forEach(heading => {
        const level = parseInt(heading.tagName.replace('H', ''), 10); // 获取标题级别（1-6）
        const text = heading.textContent; // 获取标题文本

        // 使用 DOMParser 提取纯文本
        const parser = new DOMParser();
        const doc = parser.parseFromString(text, 'text/html');
        const plainText = doc.body.textContent || '';

        // 忽略包含网址的标题
        if (!isUrl(text)) {
            const slug = text.replace(/[^\w\u4e00-\u9fa5-]+/g, '-'); // 只保留字母、数字、连字符和汉字
            toc.value.push({ slug, text, level }); // 添加到 TOC 数组
        }
    });
    console.log(content); // 打印生成的 HTML 内容
    console.log('TOC:', toc.value); // 打印 TOC 数组
    return content; // 返回完整的渲染内容
});

// 判断字符串是否为网址
function isUrl(text) {
    const urlPattern = /https?:\/\/[^\s]+/; // 匹配 http:// 或 https:// 开头的字符串
    return urlPattern.test(text);
}

function TimeFormat(time) {
    const date = new Date(time);
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hour = date.getHours();
    const minute = date.getMinutes();
    const second = date.getSeconds();
    return `${year}/${month}/${day} - ${hour}:${minute}:${second}`;
}

onMounted(async () => {
    try {
        const response = await fetch(`${URL}/path-to-article/` + articleId);
        const data = await response.json();
        article.value = data.data;
        console.log(data.data);
        // 打印 marked 渲染结果
        console.log('Rendered content:', marked(article.value.Content));
    } catch (error) {
        console.error('Failed to fetch categories:', error);
    }
    // 在内容渲染后调用 highlight.js
    nextTick(() => {
        hljs.highlightAll();
    });
});

function toggleAuthor() {
    isAuthorVisible.value = true;
    isAvatarVisible.value = false;
}

// 监听滚动事件
let lastScrollY = window.scrollY;
const handleScroll = () => {
    if (window.scrollY > lastScrollY) {
        // 向下滑动，隐藏 Author 组件
        isAuthorVisible.value = false;
        isAvatarVisible.value = true;
    }
    lastScrollY = window.scrollY;
};

onMounted(() => {
    window.addEventListener('scroll', handleScroll);
});

// 滚动到相应内容
function scrollToContent(slug) {
    console.log('Scrolling to:', slug); // 添加日志
    const element = document.getElementById(slug);
    if (element) {
        element.scrollIntoView({ behavior: 'smooth' });
    } else {
        console.log('Element not found:', slug); // 添加日志
    }
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
    display: flex;
    flex-direction: column;
    align-items: center;
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
    line-height: 1.8;
    overflow-y: auto;
}

.comment-box {
    margin-bottom: 20px;
}

.toc-box {
    margin: 20px;
    background-color: rgba(240, 248, 255, 0.7);
    padding: 20px;
    height: auto;
    overflow-y: auto;
    border-radius: 20px;
    border: solid 2px white;
    font-size: 20px;
    width: 300px;
}

.toc-box h3 {
    font-size: 20px;
    margin-bottom: 10px;
}

.toc-box ul {
    list-style-type: none;
    padding: 0;
}

.toc-box li {
    margin-bottom: 5px;
}

.toc-box a {
    text-decoration: none;
    color: rgba(0, 0, 0, 0.7);
}

.toc-box a:hover {
    text-decoration: underline;
}

#toc-title {
    font-size: 25px;
    text-align: center;
}

pre {
    background-color: rgba(0, 0, 0, 0.3);
    padding: 16px;
    border-radius: 4px;
    overflow-x: auto;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 14px;
    line-height: 1.5;
    margin: 0px;
    border: 1px solid #ddd;
}

code {
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 14px;
    color: #d63384;
}

pre code {
    color: inherit;
    background-color: rgba(0, 0, 0, 0.3);
    padding: 0;
    margin: 0px;
    line-height: 1.5;
}

.avatar-container {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 20px;
    cursor: pointer;
}

.avatar {
    object-fit: cover;
    width: 200px;
    height: 200px;
    border-radius: 50%;
    margin-bottom: 20px;
    margin: 0 auto;
}

/* .fade-enter-active,
.fade-leave-active {
    transition: transform opacity 1s ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.toc-transition-enter-active,
.toc-transition-leave-active {
    transition: transform 1s ease-in-out, opacity 1s ease-in-out;
}

.toc-transition-enter-from,
.toc-transition-leave-to {
    transform: translateY(60px);
    opacity: 0;
}

.avatar-enter-active {
    transition: transform 1s ease-in-out, opacity 1s ease-in-out;
}

.avatar-enter-from {
    transform: translateY(60px);
    width: 130px;
    height: 130px;
    opacity: 0;
} */
</style>
