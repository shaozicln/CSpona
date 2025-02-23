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
                                :class="{ 'active': activeSlug === item.slug }" @click="scrollToContent(item.slug)">
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



<!-- 改正：
 一-整体变成无上栏大页面，左上角新增一个侧边栏图标，点击侧边栏，左侧left-side变为侧边栏Slide组件，鼠标移到对应题目上面会变大并显示出相应图片
 二点击图标打开，侧边栏，然后左侧是目录和文章，点击可跳转，侧边栏布局如下：
    |--------|
    |我的头像|    
    |--------|

 《返回图标，会返回到Article》 《书本图标，会返回到现有content页面》 《首页图标，会返回到Home》
    - 目录1
    --- 文章1
    - 目录2
    --- 文章2
    --- 文章3
    - 目录3
    --- 文章4
    --- 文章5
    --- 文章6
    - 目录4
    --- 文章7
    --- 文章8
    --- 文章9
侧边栏背景为温迪
-->

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


const articleId = +localStorage.getItem('articleId');
const article = ref('');
const route = useRoute();
const toc = ref([]);
const isAuthorVisible = ref(true); // 控制 Author 组件的显示和隐藏
const isAvatarVisible = ref(false);
const authorAvatar = ref('35.jpg');
let headingCounter = {};

const renderer = {
    heading(text, level) {
        const slug = text
            .normalize("NFD") // 将字符串转换为兼容形式，处理中文字符
            .replace(/[\u0300-\u036f]/g, "") // 去掉重音符号
            .replace(/\s+/g, '-') // 替换空格为连字符
            .replace(/[^\w\u4e00-\u9fa5-]/g, '') // 只保留字母、数字、连字符和汉字
            .replace(/-+/g, '-') // 替换连续的连字符为单个连字符
            .trim() // 去除首尾空白字符
            .toLowerCase(); // 转换为小写

        let counter = (headingCounter[slug] || 0) + 1;
        headingCounter[slug] = counter;

        const id = counter > 1 ? `${slug}-${counter}` : slug;

        return `<h${level} id="${id}">${text}</h${level}>`;
    },
    code(code, language) {
        console.log("Code block detected.");
        const copyButton = `<button class="copy-button">复制</button>`;
        return `<div class="code-block"><pre><code class="language-${language}">${hljs.highlightAuto(code).value}</code>${copyButton}</pre></div>`;
    }
};

// 确保在解析 Markdown 之前应用自定义渲染器
marked.setOptions({
    gfm: true,
    breaks: true,
});

marked.use(renderer);

const renderedContent = computed(() => {
    if (!article.value) return '';
    toc.value = []; // 清空 TOC 数组
    headingCounter = {}; // 重置标题计数器

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
            const slugBase = plainText
                .normalize("NFD") // 将字符串转换为兼容形式，处理中文字符
                .replace(/[\u0300-\u036f]/g, "") // 去掉重音符号
                .replace(/[^\w\u4e00-\u9fa5-]/g, '-') // 只保留字母、数字、连字符和汉字
                .replace(/-+/g, '-') // 替换连续的连字符为单个连字符
                .toLowerCase(); // 转换为小写
            const slug = `${slugBase}-${headingCounter[slugBase] || 0}`; // 添加计数器
            headingCounter[slugBase] = (headingCounter[slugBase] || 0) + 1; // 更新计数器

            toc.value.push({ slug, text: plainText, level }); // 添加到 TOC 数组
            // console.log(`Generated slug for "${text}": ${slug}`);
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
        console.log('Fetched article data:', data.data);
        console.log('Article content:', article.value.Content);
        // 打印 marked 渲染结果
        console.log('Rendered content:', marked(article.value.Content));
    } catch (error) {
        console.error('Failed to fetch categories:', error);
    }
    // 在内容渲染后调用 highlight.js
    nextTick(() => {
        hljs.highlightAll();
        bindCopyButtons();
    });
});

function bindCopyButtons() {
    const copyButtons = document.querySelectorAll('.copy-button');
    copyButtons.forEach(button => {
        button.addEventListener('click', () => {
            const codeBlock = button.parentElement.querySelector('code');
            const codeText = codeBlock.innerText;

            const textarea = document.createElement('textarea');
            textarea.value = codeText;
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand('copy');
            document.body.removeChild(textarea);

            button.textContent = '已复制';
            setTimeout(() => {
                button.textContent = '复制';
            }, 2000);
        });
    });
}

function toggleAuthor() {
    isAuthorVisible.value = true;
    isAvatarVisible.value = false;
}

const activeSlug = ref(''); // 用于存储当前激活的标题 slug

const handleScroll = () => {
    const scrollPosition = window.scrollY;
    const tocItems = toc.value;

    for (let i = tocItems.length - 1; i >= 0; i--) {
        const element = document.getElementById(tocItems[i].slug);
        if (element) {
            const elementTop = element.offsetTop; // 使用 offsetTop 获取元素顶部位置
            const elementBottom = element.offsetTop + element.offsetHeight; // 获取元素底部位置

            if (scrollPosition + window.innerHeight / 2 >= elementTop && scrollPosition < elementBottom) {
                activeSlug.value = tocItems[i].slug;
                return;
            }
        }
    }

    activeSlug.value = '';
};

onMounted(() => {
    window.addEventListener('scroll', handleScroll);
    nextTick(() => {
        hljs.highlightAll();
        bindCopyButtons();
    });
});

// 滚动到相应内容
function scrollToContent(slug) {
    console.log('Scrolling to:', slug); // 添加日志
    nextTick(() => { // 确保 DOM 更新完成
        const element = document.getElementById(slug);
        if (element) {
            element.scrollIntoView({ behavior: 'smooth' });
        } else {
            console.error(`Element not found for slug: ${slug}`); // 添加错误日志
        }
    });
}

// 复制代码功能
function copyCode(button) {
    const codeBlock = button.parentElement.querySelector('code');
    const codeText = codeBlock.innerText;

    const textarea = document.createElement('textarea');
    textarea.value = codeText;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);

    button.textContent = '已复制';
    setTimeout(() => {
        button.textContent = '复制';
    }, 2000);
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

.toc-box li.active a {
    color: rgba(0, 0, 0, 0.2);
    /* 高亮颜色 */
    font-weight: bold;
}

.code-block {
    position: relative;
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

.copy-button {
    position: absolute;
    top: 10px;
    right: 10px;
    background-color: rgba(0, 0, 0, 0.2);
    color: rgb(0, 0, 0, 0.8);
    border: none;
    border-radius: 4px;
    padding: 5px 10px;
    cursor: pointer;
    font-size: 12px;
}

.copy-button:hover {
    background-color: rgba(0, 0, 0, 0.4);
}
</style>
