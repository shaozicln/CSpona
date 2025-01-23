<template>
  <div class="left-side">
    <input v-model="categoryId" type="number" placeholder="分类Id">
    <input @change="fileChange" type="file" placeholder="来一张文章底图">
    <input v-model="title" type="text" placeholder="输入标题" />
    <button class="button" @click="submitArticle">提交</button>
    <div class="markdown-editor">
      <textarea v-model="content" class="markdown-input" />
      <div class="markdown-preview" v-html="renderedContent"></div>
    </div>
  </div>
</template>

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from 'vue';
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import { ref, computed } from 'vue';
import { marked } from 'marked';
const title = ref('')
const content = ref('')
const categoryId = ref('')
const userId = ref(1)

const renderedContent = computed(() => {
  return marked(content.value);
});

//图片文件名字赋值
const file = ref(null);
function fileChange(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    file.value = files[0];
    console.log(file.value.name);
  } else {
    console.log('没有选择文件');
  }
}

const submitArticle = async () => {
  if (!title.value || !content.value || !categoryId.value || !userId.value) {
    alert("信息不完善");
    return;
  }
  if (categoryId.value == 2) {
    const additionalContent = '<div style="color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>';
    content.value += additionalContent; // 在文章内容末尾追加HTML
  }
  try {
    const formData = new FormData();
    formData.append('title', title.value)
    formData.append('content', content.value)
    formData.append('category_id', categoryId.value)
    formData.append('user_id', userId.value)
    formData.append('img', file.value)

    const response = await fetch(`${URL}/articles`, {
      method: 'POST',
      body: formData,
    })
    // 检查响应状态码
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json()
    console.log(data)
    console.log("标记标记");
    alert("文章上传成功 (^_^) 快去看看吧 ! ")
    window.location.reload()
  } catch (error) {
    console.error(error)
  }
}

</script>

<style scoped>
.left-side {
  flex: 1;
  padding: 20px;
  max-height: 100vh;
  border-right: 3px solid #fff;
  height: 100%;
}

.left-side,
input {
  margin-right: 50px;
}

.markdown-editor {
  display: flex;
  flex-direction: row;
  gap: 20px;
}

.markdown-input {
  width: 50%;
  height: 90vh;
  padding: 10px;
  border: 1px solid rgb(255, 255, 255);
  border-radius: 5px;
  resize: vertical;
}

.markdown-preview {
  height: 90vh;
  width: 50%;
  padding: 10px;
  background-color: rgba(255, 255, 255, 0.86);
  border: 1px solid #ddd;
  border-radius: 5px;
  overflow: auto;
}
</style>