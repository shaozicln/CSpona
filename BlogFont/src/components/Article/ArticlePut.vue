<template>
  <div class="article-put">
    <div class="modal-header">
      <div class="modal-title">
      <h2>编辑文章</h2>
      <div style="display: flex; gap: 10px; align-items: center;">
        <!-- 删除按钮 -->
        <button class="delete-btn" @click="confirmDelete" title="删除文章">
          <i class="fas fa-trash-alt"></i>
        </button>
      </div>
    </div>
      <button class="close-btn" @click="$emit('close')">×</button>
    </div>

    <div class="form-container">
      <div class="image-upload">
        <label>文章背景图:</label>
        <input type="file" @change="handleBgImageUpload" accept="image/*" class="file-input" />
        <div v-if="bgImageUrl" class="preview">
          <img :src="bgImageUrl" width="200" alt="背景图预览" />
        </div>
      </div>

      <select v-model="selectedCategoryId" @change="updateCategoryId">
        <option value="" disabled>请选择分类</option>
        <option v-for="category in categories" :key="category.Id" :value="category.Id">
          {{ category.Name }}
        </option>
      </select>

      <input v-model="title" type="text" placeholder="输入标题" />

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
          <button @click="insertText('1. ', '')" title="Numbered List">
            1.
          </button>
          <button @click="insertText('```\n', '\n```')" title="Code Block">
            { }
          </button>
          <button @click="insertText('[', '](url)')" title="Link">🔗</button>
          <button @click="showImageDialog = true" title="Insert Image">
            🖼️
          </button>
        </div>
      </div>

      <div v-if="showWrite" class="editor-area">
        <textarea v-model="content" class="markdown-input" placeholder="Type your description here..."
          ref="textarea"></textarea>
      </div>
      <div v-else class="preview-area">
        <div class="markdown-preview" v-html="renderedContent"></div>
      </div>

      <div class="button-group">
        <button class="cancel-btn" @click="$emit('close')">取消</button>
        <button class="submit-btn" @click="submitUpdate">保存修改</button>
      </div>
    </div>

    <!-- 图片插入弹窗（复用原有逻辑） -->
    <div v-if="showImageDialog" class="image-dialog">
      <div class="dialog-content">
        <h3>在文章中插入图片</h3>
        <div class="image-options">
          <div class="upload-option">
            <input type="file" @change="handleImageUpload" accept="image/*" class="file-input" />
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
// 导入Font Awesome相关资源
import { library } from '@fortawesome/fontawesome-svg-core';
import { faTrashAlt, faTimes } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
library.add(faTrashAlt, faTimes);

import { useRouter } from 'vue-router';
const router = useRouter();


// 接收父组件传递的参数
const props = defineProps({
  articleId: {
    type: [String, Number],
    required: true,
  },
  articleData: {
    type: Object,
    required: true,
  },
});

// 定义关闭事件
const emit = defineEmits(["close"]);

// 删除确认方法
const confirmDelete = () => {
  if (confirm('确定要删除这篇文章吗？此操作不可恢复！')) {
    deleteArticle();
  }
};

// 删除文章API调用
const deleteArticle = async () => {
  try {
    const response = await fetch(`${URL}/articles/${props.articleId}`, {
      method: "DELETE",
    });

    if (!response.ok) {
      throw new Error("删除失败");
    }

    alert("文章已成功删除");
    emit('close');
    emit('onDelete');
    router.push('/articles');
  } catch (error) {
    console.error("删除文章失败:", error);
    alert("删除文章失败，请重试");
  }
};

// 获取全局URL配置
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;
const URL2 = instance?.appContext.config.globalProperties.URL2;

// 数据初始化
const categories = ref([]);
const selectedCategoryId = ref("");
const title = ref("");
const content = ref("");
const showWrite = ref(true);
const textarea = ref(null);
const showImageDialog = ref(false);
const imageUrl = ref("");
const uploadedImage = ref(null);
const bgImageFile = ref(null); // 背景图文件对象
const bgImageUrl = ref(""); // 背景图预览URL

// 初始化文章数据
onMounted(() => {
  fetchCategories();
  // 从父组件传入的文章数据初始化表单
  if (props.articleData) {
    title.value = props.articleData.Title;
    content.value = props.articleData.Content;
    selectedCategoryId.value = props.articleData.CategoryId;
    // 加载原有背景图
    if (props.articleData.Img) {
      bgImageUrl.value = `${URL2}/Pictures/${props.articleData.Img}`;
    }
  }
});

// 背景图上传处理
function handleBgImageUpload(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    const file = files[0];
    bgImageFile.value = file;
    // 显示预览
    const reader = new FileReader();
    reader.onload = (e) => {
      bgImageUrl.value = e.target.result;
    };
    reader.readAsDataURL(file);
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await fetch(`${URL}/categories`);
    const data = await response.json();
    categories.value = data.data;
  } catch (error) {
    console.error("Failed to fetch categories:", error);
  }
};

// Markdown渲染配置
marked.setOptions({
  breaks: true,
  gfm: true,
});

const renderedContent = computed(() => {
  return DOMPurify.sanitize(marked(content.value));
});

// 文本插入工具（与原有编辑器逻辑一致）
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

  nextTick(() => {
    textareaEl.selectionStart = start + prefix.length;
    textareaEl.selectionEnd = end + prefix.length;
    textareaEl.focus();
  });
}

// 图片上传逻辑（与原有编辑器逻辑一致）
// 图片上传逻辑（与原有编辑器逻辑一致）
function handleImageUpload(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    const fileToUpload = files[0];
    const formData = new FormData();
    formData.append('img', fileToUpload); // 确保与后端参数名一致
    
    // 上传图片到服务器
    fetch(`${URL}/upload-image`, {
      method: 'POST',
      body: formData,
    })
    .then(response => {
      if (!response.ok) throw new Error('上传失败');
      return response.json();
    })
    .then(data => {
      // 假设后端返回格式为 { imageUrl: "文件名.png" }
      imageUrl.value = `${URL2}${data.imageUrl}`;
      console.log('图片已上传，完整URL:', imageUrl.value);
    })
    .catch(error => {
      console.error('图片上传失败:', error);
      alert('图片上传失败，请重试');
    });
  }
}

// 修复插入图片函数
function insertImage() {
  if (!imageUrl.value) {
    alert("请先上传图片");
    return;
  }
  
  // 生成正确的Markdown图片格式
  const imageMarkdown = `![Image](${imageUrl.value})`;
  
  // 获取textarea元素
  const textareaEl = textarea.value;
  if (!textareaEl) return;
  
  // 在光标位置插入图片标记
  const start = textareaEl.selectionStart;
  const end = textareaEl.selectionEnd;
  
  // 插入图片标记
  content.value = 
    content.value.substring(0, start) + 
    imageMarkdown + 
    content.value.substring(end);
  
  // 重置光标位置
  nextTick(() => {
    textareaEl.focus();
    textareaEl.selectionStart = textareaEl.selectionEnd = start + imageMarkdown.length;
  });
  
  // 关闭弹窗并重置图片URL
  showImageDialog.value = false;
  imageUrl.value = "";
}


// 更新分类ID
function updateCategoryId() {
  // 可以在这里添加额外的验证逻辑
}

// 提交更新
const submitUpdate = async () => {
  if (!title.value || !content.value || !selectedCategoryId.value) {
    alert("信息不完善");
    return;
  }

  // 使用FormData提交数据
  const formData = new FormData();
  formData.append("title", title.value);
  formData.append("content", content.value);
  formData.append("category_id", selectedCategoryId.value);
  formData.append("user_id", props.articleData.UserId);

  // 如果有新背景图，添加到FormData
  if (bgImageFile.value) {
    formData.append("img", bgImageFile.value);
  }

  try {
    const response = await fetch(`${URL}/articles/${props.articleId}`, {
      method: "PUT",
      body: formData, // 不需要设置Content-Type，浏览器会自动处理
    });

    if (!response.ok) {
      throw new Error("更新失败");
    }

    const data = await response.json();
    alert("文章更新成功！");
    emit("close");
    window.location.reload();
  } catch (error) {
    console.error("Error updating article:", error);
    alert("文章更新失败，请重试");
  }
};
</script>

<style scoped>
.image-upload {
  margin: 15px 0;
  padding: 10px;
  border: 1px dashed #ddd;
  border-radius: 4px;
}

.file-input {
  width:480px;
}

.preview {
  margin-top: 10px;
}
.delete-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #e53e3e; /* 红色 */
  transition: all 0.2s ease;
  padding: 4px 8px;
  border-radius: 4px;
}

.delete-btn:hover {
  color: #c53030; /* 深一点的红色 */
  background-color: rgba(229, 62, 62, 0.1);
}

/* 关闭按钮样式 */
.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #666;
  transition: all 0.2s ease;
  padding: 4px 8px;
  border-radius: 4px;
}

.close-btn:hover {
  color: #333;
  background-color: rgba(0, 0, 0, 0.05);
}

.modal-title{
display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 模态框头部样式 */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
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
  background-color: rgb(139,189,234);
  color: white;
  border: none;
}

.dialog-buttons button:last-child {
  background: none;
  border: 1px solid #d1d5da;
}


.button-group {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.cancel-btn {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: none;
  border-radius: 4px;
  cursor: pointer;
}

.button-group .submit-btn {
  padding: 8px 16px;
  background-color: rgb(139, 189, 234);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

/* 复用原有编辑器样式 */
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
.markdown-input,
.markdown-preview {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  resize: vertical; /* 保持 textarea 可垂直resize，和预览区行为呼应 */
  overflow: auto; /* 统一滚动行为 */
  box-sizing: border-box; /* 让 padding、border 算在宽高内，避免超出 */
}

.markdown-input {
  height: 400px; /* 可根据需求调整，和预览区初始高度一致 */
}

.markdown-preview {
  height: 400px;
  background-color: rgba(255, 255, 255, 0.86);
  /* 去掉多余的样式干扰，和 textarea 保持基础盒模型一致 */
}

.format-buttons {
  display: flex;
  gap: 7px;
  flex-wrap: wrap;
}

.format-buttons button {
  background: none;
  border: 1px solid #d1d5da;
  border-radius: 4px;
  cursor: pointer;
  min-width: 30px;
}

/* 图片弹窗样式 */
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
  z-index: 1010;
}

.dialog-content {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
}

.bu {
  width:100%;
  display: flex;
  justify-items: space-between;
  align-content: center;
  margin-bottom:10px;
  gap:10px;
}


.button-group button {
  margin-right: 5px;
  background: none;
  border: 1px solid #d1d5da;
  border-radius: 4px;
  cursor: pointer;
}

.button-group button.active {
  background-color: rgb(139,189,234);
  color: white;
  border-color: rgb(139,189,234);
}   


/* New styles for editor */
.format-buttons {
  display: flex;
  gap: 7px;
  flex-wrap: wrap;
}

.format-buttons button {
  background: none;
  border: 1px solid #d1d5da;
  border-radius: 4px;
  cursor: pointer;
  min-width: 30px;
}

.format-buttons button:hover {
  background-color: #f6f8fa;
}

</style>
