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
            <option v-for="type in types" :key="type" :value="type">
              {{ type }}
            </option>
          </select>
          <textarea
            v-model="adviceForm.content"
            placeholder="给点具体建议?"
          ></textarea>
          <br />
        </div>
        <div class="button-container">
          <button @click="createAdvice" class="button">提交</button>
        </div>
      </div>
      <div class="bottom">
        <h3>友链申请</h3>
        <br />
        <div>
          <input v-model="websiteName" type="text" placeholder="网站名称" />
          <input v-model="websiteUrl" type="text" placeholder="请输入网址" />
          <input
            v-model="websiteDescription"
            type="text"
            placeholder="来一句网站介绍"
          />
          <input
            v-model="friendDescription"
            type="text"
            placeholder="对来看该友链的人想说什么"
          />
          <div class="file-input-group">
            <label>封面展示:</label>
            <input @change="coverChange" type="file" accept="image/*" />
           
          </div>
          <div class="file-input-group">
            <label>网站背景:</label>
            <input @change="backgroundChange" type="file" accept="image/*" />
           
          </div>
          <br />
          <button @click="createApplication" class="button">提交</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

const { proxy } = getCurrentInstance();
const getImageUrl = (imgName) => {
  return `${proxy.$imageBaseUrl}${imgName}`;
};

import Author from "../Author/Author.vue";

import { ref } from "vue";

const types = ref(["技术栈の文章", "对网站の建议"]);

const adviceForm = ref({
  type: "",
  content: "",
});

const websiteName = ref();
const websiteUrl = ref();
const websiteDescription = ref();
const friendDescription = ref();

//图片文件名字赋值
const coverFile = ref(null);
const backgroundFile = ref(null);
function coverChange(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    coverFile.value = files[0];
  }
}

function backgroundChange(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    backgroundFile.value = files[0];
  }
}

const usernameWeb = localStorage.getItem("username");
const emailWeb = localStorage.getItem("email");
const avatar = localStorage.getItem("avatar");

const createAdvice = async () => {
  try {
    const response = await fetch(`${URL}/advice`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: usernameWeb,
        email: emailWeb,
        type: adviceForm.value.type,
        content: adviceForm.value.content,
      }),
    });
    const data = await response.json();
    console.log(data);
    alert("作者收到啦(^_^) 感谢建议! ");
    window.location.reload();
  } catch (error) {
    console.error(error);
  }
};
const createApplication = async () => {
  try {
    const formData = new FormData();
    formData.append("username", usernameWeb);
    formData.append("email", emailWeb);
    formData.append("name", websiteName.value);
    formData.append("web", websiteUrl.value);
    formData.append("introduction", websiteDescription.value);
    formData.append("img", coverFile.value); // 封面图片
    formData.append("avatar", avatar); // 头像图片
    formData.append("background", backgroundFile.value); // 背景图片
    formData.append("description", friendDescription.value);
    const response = await fetch(`${URL}/application`, {
      method: "POST",
      body: formData,
    });
    // 检查响应状态码
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    console.log(data);
    alert("作者收到啦(^_^) 感谢友链! 请给作者一点时间 ~ ");
    window.location.reload();
  } catch (error) {
    console.error(error);
    alert("图片上传失败");
  }
};

import { onBeforeRouteLeave } from 'vue-router';
onBeforeRouteLeave((to, from) => {
  if (to.name === 'Articles') { // 仅当跳转到 Articles 路由时设置刷新标记
    sessionStorage.setItem('refreshAfterEnter', 'Articles');
  }
});
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
