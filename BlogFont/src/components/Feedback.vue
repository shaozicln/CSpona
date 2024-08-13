<template>
    <div class="container">
        <div class="left">
            <Author />
        </div>
        <div class="right">
            <div class="top">
                <h3>反馈作者</h3>
                <div>
                    <select v-model="adviceForm.type">
                        <option value="" disabled selected>请选择类型</option>
                        <option v-for="type in types" :key="type" :value="type">{{ type }}</option>
                    </select>
                    <textarea v-model="adviceForm.content" placeholder="给点具体建议?"></textarea>
                    <br>
                </div>
                <div class="button-container">
                    <button @click="createAdvice" class="button">提交</button>
                </div>
            </div>
            <div class="bottom">
                <h3>友链申请</h3><br>
                <div>
                    <input v-model="applicationForm.websiteName" type="text" placeholder="网站名称">
                    <input v-model="applicationForm.websiteUrl" type="text" placeholder="请输入网址">
                    <input v-model="applicationForm.websiteDescription" type="text" placeholder="来一句网站介绍">
                    <input v-model="applicationForm.image" type="text" placeholder="来一张图片展示">
                    <br>
                    <button @click="createApplication" class="button">提交</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import Author from './Author.vue'

import { ref } from 'vue'

const types = ref([
    "技术栈の文章",
    "对网站の建议"
])

const adviceForm = ref({
    type: '',
    content: ''
})

const applicationForm = ref({
    websiteName: '',
    websiteUrl: '',
    websiteDescription: '',
    image: ''
})

const usernameWeb = localStorage.getItem("username")
const emailWeb = localStorage.getItem("email")

const createAdvice = async () => {
    try {
        const response = await fetch('http://127.0.0.1:8081/advice', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: usernameWeb,
                email: emailWeb,
                type: adviceForm.value.type,
                content: adviceForm.value.content,
            }),
        })
        const data = await response.json()
        console.log(data)
        alert("作者收到啦(^_^) 感谢建议! ")
        window.location.reload()
    } catch (error) {
        console.error(error)
    }
}
const createApplication = async () => {
    try {
        const response = await fetch('http://127.0.0.1:8081/application', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: usernameWeb,
                email: emailWeb,
                name: applicationForm.value.websiteName,
                web: applicationForm.value.websiteUrl,
                introduction: applicationForm.value.websiteDescription,
                img: applicationForm.value.image
            }),
        })
        const data = await response.json()
        console.log(data)
        alert("作者收到啦(^_^) 感谢友链! 请给作者一点时间 ~ ")
        window.location.reload()
    } catch (error) {
        console.error(error)
    }
}
</script>



<style scope>
.container {
    display: flex;
    padding: 20px;
    align-items: center;
    justify-content: center;
}

.left {
    width: 35%;
    padding: 20px;
    display: flex;
    flex-direction: column;
}

.right {
    width: 60%;
    padding: 5px;
    margin-left: 100px;
    font-size: 2px;
    flex-wrap: wrap;
}

.top {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    margin-bottom: 20px;
    font-size: 20px;
}

.bottom {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    margin-bottom: 20px;
    font-size: 20px;
}

input {
    width: 100%;
    height: 30px;
    width: 650px;
    padding: 10px;
    margin-bottom: 10px;
    border: 1px solid #ccc;
    background-color: rgba(255, 255, 255, 0.9);
    font-size: 20px;
}

select {
    width: 100%;
    height: 52px;
    width: 650px;
    margin-bottom: 10px;
    /* Add a margin-bottom to the select element */
    border: 1px solid #ccc;
    background-color: rgba(255, 255, 255, 0.9);
    font-size: 20px;
}

textarea {
    width: 100%;
    height: 120px;
    width: 650px;
    margin-bottom: 10px;
    /* Add a margin-bottom to the select element */
    border: 1px solid #ccc;
    background-color: rgba(255, 255, 255, 0.9);
    font-size: 25px;
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
</style>