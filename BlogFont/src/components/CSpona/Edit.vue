<template>
  <div>
    <div class="edit">
      <input @change="fileChange" type="file" placeholder="来一张文章底图" />
      <div>
        <input v-model="title" type="text" placeholder="输入标题" />
        <button class="button" @click="submitArticle">提交</button>
      </div>
    </div>

    <div class="bu">
      <div class="button-group">
        <button @click="showWrite = true" :class="{ active: showWrite }">
          Write
        </button>
        <button @click="showWrite = false" :class="{ active: !showWrite }">
          Preview
        </button>
      </div>
      <div class="format-buttons">
        <button @click="insertText('**', '**')" title="Bold"><b>B</b></button>
        <button @click="insertText('*', '*')" title="Italic"><i>I</i></button>
        <button @click="insertText('# ', '')" title="Heading">H</button>
        <button @click="insertText('> ', '')" title="Blockquote">❝</button>
        <button @click="insertText('- ', '')" title="List">•</button>
        <button @click="insertText('1. ', '')" title="Numbered List">1.</button>
        <button @click="insertText('```\n', '\n```')" title="Code Block">
          { }
        </button>
        <button @click="insertText('[', '](url)')" title="Link">🔗</button>
        <button @click="showImageDialog = true" title="Insert Image">🖼️</button>
      </div>
    </div>
    <div v-if="showWrite" class="editor-area">
      <textarea
        v-model="content"
        class="markdown-input"
        placeholder="你想在漫游地(Wanderland)分享些什么？...若遇卡顿请在备份内容后刷新"

        ref="textarea"
      ></textarea>
    </div>
    <div v-else class="preview-area">
      <div class="markdown-preview" v-html="renderedContent"></div>
    </div>

    <div v-if="showImageDialog" class="image-dialog">
      <div class="dialog-content">
        <h3>在文章中插入图片</h3>
        <div class="image-options">
          <div class="upload-option">
            <input
              type="file"
              @change="handleImageUpload"
              accept="image/*"
              class="file-input"
            />
          </div>
          <div class="url-option">
            <!-- <h4>From URL</h4>
            <input
              v-model="imageUrl"
              type="text"
              placeholder="Enter image URL"
            /> -->
          </div>
        </div>
        <div class="dialog-buttons">
          <button @click="insertImage">Insert</button>
          <button @click="showImageDialog = false">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from "vue";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;
const URL2 = instance?.appContext.config.globalProperties.URL2;

const categories = ref([]);
const selectedCategoryId = ref("");
const title = ref("");
const content = ref("");
const categoryId = ref("");
const userId = localStorage.getItem("userId");
const showWrite = ref(true);
const textarea = ref(null);
const showImageDialog = ref(false);
const imageUrl = ref("");
const uploadedImage = ref(null);
const file = ref(null);

// Configure marked
marked.setOptions({
  breaks: true,
  gfm: true,
});

const renderedContent = computed(() => {
  return DOMPurify.sanitize(marked(content.value));
});

function insertText(prefix, suffix) {
  const textareaEl = textarea.value;
  if (!textareaEl) return;

  const start = textareaEl.selectionStart;
  const end = textareaEl.selectionEnd;
  const selectedText = content.value.substring(start, end);

  content.value =
    content.value.substring(0, start) +
    prefix +
    selectedText +
    suffix +
    content.value.substring(end);

  // Set cursor position
  nextTick(() => {
    textareaEl.selectionStart = start + prefix.length;
    textareaEl.selectionEnd = end + prefix.length;
    textareaEl.focus();
  });
}

function handleImageUpload(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    const fileToUpload = files[0];
    const formData = new FormData();
    formData.append("img", fileToUpload); // 确保与后端参数名一致

    // 上传图片到服务器
    fetch(`${URL}/upload-image`, {
      method: "POST",
      body: formData,
    })
      .then((response) => {
        if (!response.ok) throw new Error("上传失败");
        return response.json();
      })
      .then((data) => {
        // 假设后端返回格式为 { imageUrl: "文件名.png" }
        imageUrl.value = `${URL2}${data.imageUrl}`;
        console.log("图片已上传，完整URL:", imageUrl.value);
      })
      .catch((error) => {
        console.error("图片上传失败:", error);
        alert("图片上传失败，请重试");
      });
  }
}

function insertImage() {
  let imageMarkdown = "";

  // 优先使用服务器返回的URL（推荐方式）
  if (imageUrl.value) {
    imageMarkdown = `![Image](${imageUrl.value})`;
  }
  // 仅作为备选方案（实际应该不会执行到这）
  else if (uploadedImage.value) {
    imageMarkdown = `![Image](${uploadedImage.value})`;
  }

  if (imageMarkdown) {
    insertText(imageMarkdown, "");
    showImageDialog.value = false;
    imageUrl.value = ""; // 清空URL
    uploadedImage.value = null; // 清空base64
  }
}
function fileChange(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    file.value = files[0];
    console.log(file.value.name);
  } else {
    console.log("没有选择文件");
  }
}

const fetchCategories = async () => {
  try {
    const response = await fetch(`${URL}/categories`);
    const data = await response.json();
    categories.value = data.data;
    console.log(data.data);
  } catch (error) {
    console.error("Failed to fetch categories:", error);
  }
};

function updateCategoryId() {
  categoryId.value = selectedCategoryId.value;
  console.log(selectedCategoryId.value);
}

const submitArticle = async () => {
  if (!title.value || !content.value || !userId) {
    alert("信息不完善");
    return;
  }

  if (categoryId.value == 2) {
    const additionalContent =
      '<div style="color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>';
    content.value += additionalContent;
  }

  const formData = new FormData();
  formData.append("title", title.value);
  formData.append("content", content.value);
  formData.append("category_id", 1000);
  formData.append("user_id", userId);
  if (file.value) {
    formData.append("img", file.value);
  }

  try {
    const response = await fetch(`${URL}/articles`, {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      throw new Error("Network response was not ok");
    }

    const data = await response.json();
    console.log(data);
    alert("文章上传成功 (^_^) 快去看看吧 ! ");
    window.location.reload();
  } catch (error) {
    console.error("Error submitting article:", error);
    alert("文章上传失败");
  }
};

onMounted(() => {
  fetchCategories();
});
</script>

<style scoped>
.edit{
    display: flex;
    flex-direction: column;
}

.button-group button {
  padding: 5px 10px;
  margin-right: 5px;
  background: none;
  border: 1px solid #d1d5da;
  border-radius: 4px;
  cursor: pointer;
}

.button-group button.active {
  background-color: rgb(139, 189, 234);
  color: white;
  border-color: rgb(139, 189, 234);
}

.markdown-input {
  width: 100%;
  height: 400px;
  padding: 10px;
  border: 1px solid rgb(255, 255, 255);
  border-radius: 5px;
  resize: vertical;
}

.markdown-preview {
  height: 400px;
  padding: 10px;
  background-color: rgba(255, 255, 255, 0.86);
  border: 1px solid #ddd;
  border-radius: 5px;
  overflow: auto;
}

.format-buttons {
  display: flex;
  gap: 7px;
  flex-wrap: wrap;
}

.format-buttons button {
  padding: 5px 10px;
  background: none;
  border: 1px solid #d1d5da;
  border-radius: 4px;
  cursor: pointer;
  min-width: 30px;
}

.format-buttons button:hover {
  background-color: #f6f8fa;
}

/* Image dialog styles */
.image-dialog {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.dialog-content {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
}

.image-options {
  display: flex;
  gap: 20px;
  margin: 20px 0;
}

.upload-option,
.url-option {
  flex: 1;
}

.dialog-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.dialog-buttons button {
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.dialog-buttons button:first-child {
  background-color: rgb(139, 189, 234);
  color: white;
  border: none;
}

.dialog-buttons button:last-child {
  background: none;
  border: 1px solid #d1d5da;
}

/* Preview styles */
.markdown-preview :deep() img {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

.markdown-preview :deep() pre {
  background-color: #f6f8fa;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
}

.markdown-preview :deep() code {
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
  background-color: rgba(27, 31, 35, 0.05);
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-size: 85%;
}

.bu {
  width: 100%;
  display: flex;
  justify-items: space-between;
  align-content: center;
  margin-bottom: 10px;
}

.file-input {
  width: 480px;
}
</style>
