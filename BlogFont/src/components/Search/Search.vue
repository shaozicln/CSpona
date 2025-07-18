<template>
    <div class="search-container" :style="{ background: searchBackground }" >
        <div class="search-overlay" tabindex="0" @click="removeDarkBackground"/>
        <div :class="['search-form-container', isFocused ? 'search-form-container-focused' : '']" >
            <form @submit.prevent="performSearch">
                <div class="search-item">
                    <input type="text" placeholder="搜索文章..." v-model="searchQuery" ref="searchInput"
                    @focus="restoreDarkBackground" />
                    <button @mousedown.prevent type="submit"><i class="iconfont icon-sousuo"></i></button>
                </div>
            </form>
            <br>
            <div v-if="searchResultVisible">
                <div v-if="titles.length === 0 && contents.length === 0" class="send">
                    <p>什么都没查询到捏 QAQ </p>
                    <p>想看相关内容? 点击导航栏"反馈"! </p>
                    <router-link :to="{ path: '/feedback' }"><button class="send"
                            @click.prevent="removeDarkBackgroundSecond()">点我私信作者</button>
                    </router-link>
                </div>
                <div v-else class="result">
                    <h3>标题搜索结果</h3>
                    <ul v-for="title in titles" :key="title.Id">
                        <li @click="getArticleContent(title.Id)">
                            {{ title.Title }}
                        </li>
                    </ul>
                    <br>
                    <h3>内容搜索结果</h3>
                    <ul v-for="content in contents" :key="content.Id">
                        <li @click="getArticleContent(content.Id)">
                            {{ content.Title }}
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>

// 获取全局URL属性
import { getCurrentInstance } from 'vue';
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import { ref, onMounted } from 'vue'

import { useRouter } from 'vue-router'
const router = useRouter()

const searchInput = ref(null)
const searchQuery = ref('')
const searchBackground = ref('rgba(0, 0, 0, 0.5)')
const searchPosition = ref('relative')
const searchHeight = ref('88.5vh')
const isFocused = ref(false)
const searchResultVisible = ref(false)

onMounted(() => {
    searchInput.value.focus()
})

const removeDarkBackground = () => {
    searchBackground.value = 'transparent'
    searchPosition.value = 'relative'
    searchHeight.value = '88.5vh'
    isFocused.value = false
    searchResultVisible.value = false
}

const removeDarkBackgroundSecond = ( ) => {
    searchBackground.value = 'transparent'
    searchPosition.value = 'relative'
    searchHeight.value = '88.5vh'
    isFocused.value = false
    searchResultVisible.value = false
    if (confirm('即将跳转界面')) {
        router.push({ path: '/feedback' })
    }
}

const restoreDarkBackground = () => {
    searchBackground.value = 'rgba(0, 0, 0, 0.5)'
    searchPosition.value = 'fixed'
    searchHeight.value = '100vh'
    isFocused.value = true
}

const titles = ref([])
const contents = ref([])

const performSearch = async () => {
    if (!isFocused.value) return
    searchResultVisible.value = true
    setTimeout(() => {
        searchInput.value.focus();
    }, 0);

    const titleUrl = `${ URL }/search?title=${encodeURIComponent(searchQuery.value)}`;
    const contentUrl = `${ URL }/search?content=${encodeURIComponent(searchQuery.value)}`;

    try {
        const titleResponse = await fetch(titleUrl);
        const titleData = await titleResponse.json();
        titles.value = titleData.data

        const contentResponse = await fetch(contentUrl);
        const contentData = await contentResponse.json();
        contents.value = contentData.data

        // if (titleData.data.length === 0 && contentData.data.length === 0) {
        //     titles.value = [{ Title: '什么都没查询到捏' }],
        //         contents.value = [{ Title: '什么都没查询到捏' }]
        // }

        console.log(titleData);
        console.log(contentData);
    } catch (error) {
        console.error(error);
    }
}

//获取文章id
const getArticleContent = (articleId) => {
  const encodedId = window.btoa(articleId.toString());
  router.push(`/articleContent/${encodedId}`);
};

import { onBeforeRouteLeave } from 'vue-router';
onBeforeRouteLeave((to, from) => {
  if (to.name === 'Articles') { // 仅当跳转到 Articles 路由时设置刷新标记
    sessionStorage.setItem('refreshAfterEnter', 'Articles');
  }
});
</script>

<style scoped>
.search-container {
    position: v-bind(searchPosition);
    top: 0;
    left: 0;
    width: 100%;
    height: v-bind(searchHeight);
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: v-bind(searchBackground);
    /* dark background */
}

.search-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: transparent;
    z-index: -10;
}

.search-form-container {
    background-color: rgba(255, 255, 255, 0.5);
    /* white background with 0.5 opacity */
    padding: 50px 50px 40px 50px;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    transition: all 0.3s ease-in-out;
}

.search-form-container-focused {
    background-color: rgba(255, 255, 255, 0.8);
    width: 450px;
    height: 600px;
    padding: 50px 50px 40px 50px;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    gap: 0;
    overflow: auto;
}

.search-item {
    display: flex;
    align-items: center;
    margin-right: 20px;
}

.search-item input {
    margin-right: 5px;
    width: 100%;
    height: 40px;
    border: 1px solid #ccc;
    border-radius: 5px;
    font-size: 20px;
}

.search-item button {
    background-color: rgb(139, 189, 234);
    /* blue button */
    color: #fff;
    border: none;
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease-in-out;
    border-radius: 20px;
}

.search-item button:hover {
    background-color: rgba(139, 190, 234, 0.686);
    transform: translateY(-5px);
    box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}

.result {
    overflow: auto;
    max-height: 70vh;
}

.result h3 {
    font-size: 30px;
    color: rgba(0, 0, 0, 0.848)
}

.result li {
    font-size: 20px;
}

.icon-sousuo {
    background-color: rgba(255, 255, 255, 0) !important;
}

.send {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
}

.send p {
    font-size: 20px;
    margin-bottom: 10px;
}

.send i {
    font-size: 24px;
    margin-bottom: 10px;
}

.send button {
    background-color: rgb(139, 189, 234);
    color: #fff;
    border: none;
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease-in-out;
    border-radius: 20px;
}

.send button:hover {
    background-color: rgba(139, 190, 234, 0.686);
    transform: translateY(-5px);
    box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}
</style>