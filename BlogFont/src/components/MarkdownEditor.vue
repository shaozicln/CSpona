<template>
  <div class="left-side">
    <input v-model="categoryId" type="number" placeholder="分类Id">
    <input v-model="img" type="text" placeholder="图片名称">
    <input v-model="title" type="text" placeholder="输入标题" />
    <button class="button" @click="submitArticle">提交</button>
    <div class="markdown-editor">
      <textarea v-model="content" class="markdown-input" />
      <div class="markdown-preview" v-html="renderedContent"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { marked } from 'marked';
const title = ref('')
const content = ref('')
const categoryId = ref('')
const userId = ref(1)
const img = ref('')

const renderedContent = computed(() => {
  return marked(content.value);
});

const submitArticle = async () => {
  if (!title.value || !content.value || !categoryId.value || !userId.value || !img.value) {
    alert("信息不完善");
    return;
  }
  try {
    const response = await fetch('http://127.0.0.1:8081/articles', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        title: title.value,
        content: content.value,
        categoryId: categoryId.value,
        userId: userId.value,
        img: img.value
      }),
    })
    const data = await response.json()
    console.log(data)
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