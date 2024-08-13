<template>
    <div class="info-box">
        <img :src="'../Public/Pictures/' + authorAvatar" alt="Author Avatar" id="avatar">
        <ul>
            <li>网站作者: {{ author }}</li>
            <li>邮箱: {{ email }}</li>
            <li>gitHub: {{ gitHub }}</li>
            <li>性别: {{ sex }}</li>
            <li>大学: {{ university }}</li>
            <li>专业: {{ major }}</li>
            <li>现居地: {{ address }}</li>
        </ul>
    </div>
</template>

<script setup>
import { ref } from 'vue'

const authorAvatar = ref('35.jpg')
const author = ref('长柄木勺')
const email = ref('changbingmushao@qq.com')
const gitHub = ref('https://github.com/shaozicln')
const sex = ref('女')
const university = ref('东北农业大学')
const major = ref('计算机科学与技术')
const address = ref('中国-哈尔滨')
</script>

<style>
.info-box {
    background-color: rgba(255, 255, 255, 0.5);
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    text-align: center;
}

#avatar {
    width: 200px;
    height: 200px;
    border-radius: 50%;
    margin-bottom: 20px;
    margin: 0 auto;
}

.info-box ul {
    list-style: none;
    padding: 0;
    margin: 0;
    font-size: 25px;
}
</style>