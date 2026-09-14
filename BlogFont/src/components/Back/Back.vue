<template>
    <div class="container">
        <div class="left-side">
            <markdown-editor v-model="content" />
        </div>
        <span class="vertical-bar"></span>
        <div class="right-side">
            <div id="author">
                <AuthorBack />
            </div>
            <br>
            <div class="button-group">
                <button class="button" @click="advice">建议查看</button>
                <button class="button" @click="application">友链申请</button>
                <button class="button" @click="article">文章管理</button>
                <button class="button" @click="message">留言管理</button>
                <button class="button" @click="openAboutMe">关于我</button>
            </div>
        </div>
        <div class="modal" v-if="showModal" @click.self="hideModal">
            <div class="modal-content">
                <div id="button-group">
                    <button class="button" @click="sortByDefault">默认顺序</button>
                    <button class="button" @click="sortByType">类型顺序</button>
                    <button class="button" @click="sortByUsername">用户顺序</button>
                </div>
                <div class="advice-list">
                    <div v-for="(advice, index) in advices" :key="advice.id">
                        <h1>{{ index + 1 }}</h1>
                        <p>U：{{ advice.Username }}</p>
                        <p>T：{{ advice.Type }}</p>
                        <p>C：{{ advice.Content }}</p>
                    </div>
                </div>
            </div>
        </div>
        <div class="modal" v-if="showModal2" @click.self="hideModal2">
            <div class="model-content">
                <div class="application-list">
                    <div v-for="(application, index) in applications" :key="application.id">
                        <div class="p">
                            <div>
                                <h1>{{ index + 1 }}</h1>
                                <div id="button-group">
                                    <button class="button" @click="applicationPost(application.Id)">同意申请</button>
                                    <button class="button" @click="applicationDelete(application.Id)">删除申请</button>
                                </div>
                                <p>U：{{ application.Username }}</p>
                                <p>E：{{ application.Email }}</p>
                                <p>N：{{ application.Name }}</p>
                                <p>W：{{ application.Web }}</p>
                                <p>I：{{ application.Introduction }}</p>
                            </div>
                            <img :src="getImageUrl(application.Img)">
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="modal" v-if="showAboutModal" @click.self="hideAboutModal">
            <div class="about-modal">
                <div class="about-modal-header">
                    <h2>编辑「关于我」</h2>
                    <div class="about-actions">
                        <button class="button" type="button" @click="aboutPreview = !aboutPreview">
                            {{ aboutPreview ? '编辑' : '预览' }}
                        </button>
                        <button class="button" type="button" :disabled="aboutSaving" @click="saveAboutMe">
                            {{ aboutSaving ? '保存中…' : '保存' }}
                        </button>
                    </div>
                </div>
                <p class="about-hint">内容保存在服务器 Markdown 文件，保存后首页刷新即可看到。</p>
                <textarea
                    v-if="!aboutPreview"
                    v-model="aboutContent"
                    class="about-editor"
                    placeholder="用 Markdown 写关于我…"
                ></textarea>
                <div v-else class="about-preview" v-html="aboutPreviewHtml"></div>
            </div>
        </div>
    </div>
</template>

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from 'vue';
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import AuthorBack from '../Author/AuthorBack.vue'
import { ref, computed } from 'vue';
import MarkdownEditor from '../Manage/MarkdownEditor.vue';
import { nextTick } from 'vue';
import { apiJson, promptLoginIfUnauthorized } from '@/utils/api';
import { marked } from 'marked';
import DOMPurify from 'dompurify';

marked.setOptions({ gfm: true, breaks: true });

const { proxy } = getCurrentInstance()
const getImageUrl = (imgName) => {
  return `${proxy.$imageBaseUrl}${imgName}`
}

const title = ref('');
const content = ref('');
const showModal = ref(false);
const advices = ref([]);
const showAboutModal = ref(false);
const aboutContent = ref('');
const aboutPreview = ref(false);
const aboutSaving = ref(false);

const aboutPreviewHtml = computed(() => {
  return DOMPurify.sanitize(marked.parse(aboutContent.value || ''));
});

const openAboutMe = async () => {
    showAboutModal.value = true;
    aboutPreview.value = false;
    try {
        const { res, data } = await apiJson('/about-me');
        if (!res.ok) {
            alert(data?.message || '加载失败');
            return;
        }
        aboutContent.value = data?.data?.content || '';
    } catch (error) {
        console.error(error);
        alert('加载失败');
    }
};

const hideAboutModal = () => {
    showAboutModal.value = false;
};

const saveAboutMe = async () => {
    aboutSaving.value = true;
    try {
        const { res, data } = await apiJson('/about-me', {
            method: 'PUT',
            body: JSON.stringify({ content: aboutContent.value }),
        });
        if (promptLoginIfUnauthorized(res, data)) return;
        if (!res.ok) {
            alert(data?.message || '保存失败');
            return;
        }
        alert('关于我已保存');
        showAboutModal.value = false;
    } catch (error) {
        console.error(error);
        alert('保存失败');
    } finally {
        aboutSaving.value = false;
    }
};

const advice = async () => {
    showModal.value = true;
    try {
        const response = await fetch(`${URL}/advice`)
        const data = await response.json()
        console.log(data)
        advices.value = data.data;
    } catch (error) {
        console.error(error)
    }
}

const sortByDefault = async () => {
    showModal.value = true;
    try {
        const response = await fetch(`${URL}/advice`)
        const data = await response.json()
        console.log(data)
        advices.value = data.data;
    } catch (error) {
        console.error(error)
    }
}

const sortByType = async () => {
    showModal.value = true;
    try {
        const response = await fetch(`${ URL } /advice`)
        const data = await response.json()
        console.log(data)
        advices.value = data.data.sort((a, b) => {
            if (a.Type < b.Type) return -1;
            if (a.Type > b.Type) return 1;
            return 0;
        });
    } catch (error) {
        console.error(error)
    }
}

const sortByUsername = async () => {
    showModal.value = true;
    try {
        const response = await fetch(`${URL}/advice`)
        const data = await response.json()
        console.log(data)
        advices.value = data.data.sort((a, b) => {
            if (a.Username < b.Username) return -1;
            if (a.Username > b.Username) return 1;
            return 0;
        });
    } catch (error) {
        console.error(error)
    }
}

function hideModal() {
    showModal.value = false;
}


const showModal2 = ref(false);
const applications = ref([]);


const application = async () => {
    showModal2.value = true;
    try {
        const response = await fetch(`${URL}/application`)
        const data = await response.json()
        console.log(data)
        applications.value = data.data;
    } catch (error) {
        console.error(error)
    }
}
function hideModal2() {
    showModal2.value = false;
}

const applicationPost = async (index) => {
    const newFriendsWeb = ref('');
    try {
        const response = await fetch(`${URL}/application?id=` + index)
        const data = await response.json()
        console.log(data)
        newFriendsWeb.value = data.data;
        await nextTick(); // wait for the assignment to complete
    } catch (error) {
        console.error(error)
    }
    try {
        const response = await fetch(`${URL}/friendsWeb`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: newFriendsWeb.value[0].Name,
                web: newFriendsWeb.value[0].Web,
                introduction: newFriendsWeb.value[0].Introduction,
                img: newFriendsWeb.value[0].Img,
                avatar: newFriendsWeb.value[0].Avatar,
                background: newFriendsWeb.value[0].Background,
                description: newFriendsWeb.value[0].Description,
            }),
        })
        const data = await response.json()
        console.log(data)
        alert("友链添加成功 (^_^) !")
    } catch (error) {
        console.error(error)
    }
}
const applicationDelete = async (index) => {
    try {
        const response = await fetch(`${URL}/application/` + index, {
            method: 'DELETE',
        })
        const data = await response.json()
        console.log(data)
        alert("友链申请删除成功 (-w-) !")
    } catch (error) {
        console.error(error)
    }
}


function article() {
    // button 1 click handler
}

function message() {

}

import { onBeforeRouteLeave } from 'vue-router';
onBeforeRouteLeave((to, from) => {
  if (to.name === 'Articles') { // 仅当跳转到 Articles 路由时设置刷新标记
    sessionStorage.setItem('refreshAfterEnter', 'Articles');
  }
});

</script>

<style scoped>
markdown-editor {
    height: 100vh;
    /* 添加滚动条 */
}

.container {
    display: flex;
    flex-direction: row;
    height: 100vh;
}

.left-side {
    flex: 1;
    padding: 20px;
    max-height: 100vh;
    height: 100%;
}

.left-side,
input {
    margin-right: 50px;
}

.right-side {
    width: 200px;
    padding: 12px;
    height: 100vh;
    width: 15%;
    display: flex;
    flex-direction: column;
    justify-content: center;
}


.button-group {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 50px;
    gap: 20px;
}

.button {
    width: 100px;
    height: 40px;
    background-color: rgb(139, 189, 234);
    /* 修改按钮颜色 */
    color: #fff;
    padding: 10px;
    border: none;
    border-radius: 20px;
    cursor: pointer;
    position: relative;
    transition: all 0.3s ease-in-out;
    font-size: 20px;
}

.button:hover {
    background-color: rgb(139, 189, 234);
    transform: translateY(-5px);
    box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}

.modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
}

.modal-content {
    background-color: #fff;
    opacity: 0.8;
    padding: 25px;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    width: 40vw;
    height: 80vh;
    overflow-y: auto;
}

#button-group {
    display: flex;
    flex-direction: row;
    align-items: center;
    padding: 10px;
    gap: 10px;
}

.modal-button {
    width: 100px;
    height: 40px;
    background-color: rgb(139, 189, 234);
    color: #fff;
    padding: 10px;
    border: none;
    border-radius: 20px;
    cursor: pointer;
    position: relative;
    transition: all 0.3s ease-in-out;
    font-size: 20px;
}

.advice-list {
    padding: 20px;
}

.application-list {
    padding: 40px;
    overflow-y: auto;
    /* 添加滚动条 */
    background-color: rgba(255, 255, 255, 0.837);
    /* 设置白色背景 */
    width: 40vw;
    height: 80vh;
    border-radius: 20px;
}

.application-list div {
    margin-bottom: 20px;
}

.application-list h1 {
    font-size: 28px;
    margin-bottom: 10px;
}

h1 {
    font-size: 28px;
}

.application-list p {
    font-size: 25px;
    margin-bottom: 10px;
    display: flex;
    flex-direction: column;
    flex: 1;
    /* add this to make the paragraph take up the remaining space */
}

.application-list img {
    width: 200px;
    height: 200px;
    border-radius: 10px;
    margin: 10px;
    margin-left: 0;
    /* remove the margin left */
}

.p {
    display: flex;
    align-content: center;
    justify-content: space-between;
}

.about-modal {
    background-color: var(--panel-bg, #fff);
    color: var(--text-primary, #222);
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    width: min(720px, 92vw);
    height: min(80vh, 720px);
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.about-modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
}

.about-modal-header h2 {
    font-size: 22px;
    margin: 0;
}

.about-actions {
    display: flex;
    gap: 10px;
}

.about-hint {
    font-size: 14px;
    color: var(--text-secondary, #666);
    margin: 0;
}

.about-editor,
.about-preview {
    flex: 1;
    min-height: 0;
    width: 100%;
    box-sizing: border-box;
    border: 1px solid var(--input-border, #ddd);
    border-radius: 8px;
    padding: 14px;
    font-size: 16px;
    line-height: 1.6;
    background: var(--input-bg, #fff);
    color: var(--input-text, #222);
    overflow: auto;
    resize: none;
    font-family: Consolas, Monaco, monospace;
}

.about-preview {
    font-family: "楷体", serif;
}
</style>