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
                            <img :src="'../Public/Pictures/' + application.Img">
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import AuthorBack from './AuthorBack.vue'
import { ref } from 'vue';
import MarkdownEditor from './MarkdownEditor.vue';
import { nextTick } from 'vue';

const title = ref('');
const content = ref('');
const showModal = ref(false);
const advices = ref([]);


const advice = async () => {
    showModal.value = true;
    try {
        const response = await fetch('http://127.0.0.1:8081/advice')
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
        const response = await fetch('http://127.0.0.1:8081/advice')
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
        const response = await fetch('http://127.0.0.1:8081/advice')
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
        const response = await fetch('http://127.0.0.1:8081/advice')
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
        const response = await fetch('http://127.0.0.1:8081/application')
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
        const response = await fetch('http://127.0.0.1:8081/application?id=' + index)
        const data = await response.json()
        console.log(data)
        newFriendsWeb.value = data.data;
        await nextTick(); // wait for the assignment to complete
    } catch (error) {
        console.error(error)
    }
    try {
        const response = await fetch('http://127.0.0.1:8081/friendsWeb', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: newFriendsWeb.value[0].Name,
                web: newFriendsWeb.value[0].Web,
                introduction: newFriendsWeb.value[0].Introduction,
                img: newFriendsWeb.value[0].Img,
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
        const response = await fetch('http://127.0.0.1:8081/application/' + index, {
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
</style>