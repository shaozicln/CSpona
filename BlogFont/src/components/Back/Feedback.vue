<template>
  <div class="container">
    <div class="left">
      <Author />
    </div>
    <div class="right">
      <div class="panel top">
        <h3>反馈作者</h3>
        <div class="form-block">
          <select v-model="adviceForm.type">
            <option value="" disabled>请选择类型</option>
            <option v-for="type in types" :key="type" :value="type">
              {{ type }}
            </option>
          </select>
          <template v-if="adviceForm.type === '背景音乐の推荐'">
            <input v-model="musicTitle" type="text" placeholder="歌名" />
            <select v-model="musicKind">
              <option value="netease">网易云链接</option>
              <option value="file">上传音乐文件</option>
            </select>
            <input
              v-if="musicKind === 'netease'"
              v-model="musicRef"
              type="text"
              placeholder="网易云歌曲/歌单分享链接"
            />
            <div v-else class="file-input-group">
              <label>音乐文件:</label>
              <input type="file" accept="audio/*,.mp3,.m4a,.flac,.ogg,.wav,.aac" @change="musicFileChange" />
              <span v-if="musicFileName" class="file-name">{{ musicFileName }}</span>
            </div>
            <textarea
              v-model="adviceForm.content"
              placeholder="想对作者说的话（可选）"
            ></textarea>
          </template>
          <textarea
            v-else
            v-model="adviceForm.content"
            placeholder="给点具体建议?"
          ></textarea>
        </div>
        <div class="button-container">
          <button @click="createAdvice" class="button" type="button">提交</button>
        </div>
      </div>

      <div class="panel bottom">
        <h3>友链申请</h3>
        <div class="form-block">
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
            <span v-if="coverFileName" class="file-name">{{ coverFileName }}</span>
          </div>
          <div class="file-input-group">
            <label>网站背景:</label>
            <input @change="backgroundChange" type="file" accept="image/*" />
            <span v-if="backgroundFileName" class="file-name">{{ backgroundFileName }}</span>
          </div>
        </div>
        <div class="button-container">
          <button @click="createApplication" class="button" type="button">提交</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getCurrentInstance, ref } from "vue";
import Author from "../Author/Author.vue";

const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

const types = ref(["技术栈の文章", "对网站の建议", "背景音乐の推荐"]);

const adviceForm = ref({
  type: "",
  content: "",
});
const musicTitle = ref("");
const musicKind = ref("netease");
const musicRef = ref("");
const musicFile = ref(null);
const musicFileName = ref("");

const websiteName = ref("");
const websiteUrl = ref("");
const websiteDescription = ref("");
const friendDescription = ref("");

const coverFileName = ref("");
const backgroundFileName = ref("");
const coverFile = ref(null);
const backgroundFile = ref(null);
const MAX_FILE_SIZE = 1 * 1024 * 1024;
const MAX_MUSIC_SIZE = 15 * 1024 * 1024;

function musicFileChange(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    const file = files[0];
    if (file.size > MAX_MUSIC_SIZE) {
      alert("音频不能超过 15MB");
      event.target.value = "";
      musicFile.value = null;
      musicFileName.value = "";
      return;
    }
    musicFile.value = file;
    musicFileName.value = file.name;
  } else {
    musicFile.value = null;
    musicFileName.value = "";
  }
}

function formatFileSize(bytes) {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}

function coverChange(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    const file = files[0];
    if (file.size > MAX_FILE_SIZE) {
      alert("封面图片大小不能超过1MB，请重新选择");
      event.target.value = "";
      coverFile.value = null;
      coverFileName.value = "";
      return;
    }
    if (!file.type.startsWith("image/")) {
      alert("请选择有效的图片文件");
      event.target.value = "";
      coverFile.value = null;
      coverFileName.value = "";
      return;
    }
    coverFile.value = file;
    coverFileName.value = file.name;
    console.log(`封面图片: ${file.name}, 大小: ${formatFileSize(file.size)}`);
  } else {
    coverFile.value = null;
    coverFileName.value = "";
  }
}

function backgroundChange(event) {
  const files = event.target.files;
  if (files && files.length > 0) {
    const file = files[0];
    if (file.size > MAX_FILE_SIZE) {
      alert("背景图片大小不能超过1MB，请重新选择");
      event.target.value = "";
      backgroundFile.value = null;
      backgroundFileName.value = "";
      return;
    }
    if (!file.type.startsWith("image/")) {
      alert("请选择有效的图片文件");
      event.target.value = "";
      backgroundFile.value = null;
      backgroundFileName.value = "";
      return;
    }
    backgroundFile.value = file;
    backgroundFileName.value = file.name;
    console.log(`背景图片: ${file.name}, 大小: ${formatFileSize(file.size)}`);
  } else {
    backgroundFile.value = null;
    backgroundFileName.value = "";
  }
}

const usernameWeb = localStorage.getItem("username");
const emailWeb = localStorage.getItem("email");
const avatar = localStorage.getItem("avatar");

const createAdvice = async () => {
  try {
    if (!adviceForm.value.type) {
      alert("请选择类型");
      return;
    }

    if (adviceForm.value.type === "背景音乐の推荐") {
      if (!musicTitle.value.trim()) {
        alert("请填写歌名");
        return;
      }
      if (musicKind.value === "netease" && !musicRef.value.trim()) {
        alert("请填写网易云链接");
        return;
      }
      if (musicKind.value === "file" && !musicFile.value) {
        alert("请上传音乐文件");
        return;
      }

      const formData = new FormData();
      formData.append("username", usernameWeb || "");
      formData.append("email", emailWeb || "");
      formData.append("type", adviceForm.value.type);
      formData.append("content", adviceForm.value.content || `推荐：${musicTitle.value}`);
      formData.append("music_title", musicTitle.value.trim());
      formData.append("music_kind", musicKind.value);
      if (musicKind.value === "netease") {
        formData.append("music_ref", musicRef.value.trim());
      } else {
        formData.append("music_file", musicFile.value);
      }
      const response = await fetch(`${URL}/advice`, {
        method: "POST",
        body: formData,
      });
      if (!response.ok) {
        const err = await response.json().catch(() => ({}));
        alert(err.error || "提交失败");
        return;
      }
      alert("作者收到啦(^_^) 感谢推荐! ");
      adviceForm.value = { type: "", content: "" };
      musicTitle.value = "";
      musicKind.value = "netease";
      musicRef.value = "";
      musicFile.value = null;
      musicFileName.value = "";
      return;
    }

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
    adviceForm.value = { type: "", content: "" };
  } catch (error) {
    console.error(error);
  }
};

const createApplication = async () => {
  try {
    if (!coverFile.value) {
      alert("请上传封面图片");
      return;
    }
    if (!backgroundFile.value) {
      alert("请上传背景图片");
      return;
    }
    if (coverFile.value.size > MAX_FILE_SIZE) {
      alert("封面图片大小不能超过1MB");
      return;
    }
    if (backgroundFile.value.size > MAX_FILE_SIZE) {
      alert("背景图片大小不能超过1MB");
      return;
    }

    const formData = new FormData();
    formData.append("username", usernameWeb);
    formData.append("email", emailWeb);
    formData.append("name", websiteName.value);
    formData.append("web", websiteUrl.value);
    formData.append("introduction", websiteDescription.value);
    formData.append("img", coverFile.value);
    formData.append("avatar", avatar);
    formData.append("background", backgroundFile.value);
    formData.append("description", friendDescription.value);
    const response = await fetch(`${URL}/application`, {
      method: "POST",
      body: formData,
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    console.log(data);
    alert("作者收到啦(^_^) 感谢友链! 请给作者一点时间 ~ ");
    websiteName.value = "";
    websiteUrl.value = "";
    websiteDescription.value = "";
    friendDescription.value = "";
    coverFile.value = null;
    backgroundFile.value = null;
    coverFileName.value = "";
    backgroundFileName.value = "";
  } catch (error) {
    console.error(error);
    alert("图片上传失败");
  }
};
</script>

<style scoped>
.container {
  display: flex;
  padding: 20px 40px 40px;
  align-items: flex-start;
  justify-content: center;
  gap: 40px;
  min-height: calc(100vh - 120px);
  box-sizing: border-box;
  color: var(--text-primary);
}

.left {
  width: 35%;
  max-width: 360px;
  padding: 20px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.right {
  width: 60%;
  max-width: 720px;
  padding: 5px 0;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.panel {
  background-color: var(--panel-bg);
  border: 1px solid var(--panel-border);
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  padding: 24px 28px;
  box-sizing: border-box;
}

.panel h3 {
  margin: 0 0 16px;
  font-size: 28px;
  text-align: center;
  color: var(--text-primary);
  font-family: Cormorant SC, serif;
}

.form-block {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

input,
select,
textarea {
  width: 100%;
  box-sizing: border-box;
  padding: 10px 12px;
  margin: 0;
  border: 1px solid var(--input-border);
  border-radius: 5px;
  background-color: var(--input-bg);
  color: var(--input-text);
  font-size: 20px;
  font-family: Cormorant SC, serif;
}

input::placeholder,
textarea::placeholder {
  color: var(--input-placeholder);
}

select {
  height: 48px;
  cursor: pointer;
}

textarea {
  min-height: 120px;
  resize: vertical;
  font-size: 22px;
  line-height: 1.5;
}

.file-input-group {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  margin-top: 4px;
  color: var(--text-secondary);
  font-size: 18px;
}

.file-input-group label {
  min-width: 88px;
  color: var(--text-primary);
}

.file-input-group input[type="file"] {
  width: auto;
  flex: 1;
  min-width: 180px;
  padding: 6px;
  font-size: 16px;
  background: transparent;
  border: 1px dashed var(--input-border);
  color: var(--text-secondary);
}

.file-name {
  width: 100%;
  font-size: 15px;
  color: var(--text-muted);
  padding-left: 98px;
  box-sizing: border-box;
}

.button-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.button {
  width: 100px;
  height: 40px;
  background-color: var(--btn-bg);
  color: var(--btn-text);
  padding: 8px 10px;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease-in-out;
  font-size: 20px;
  font-family: Cormorant SC, serif;
}

.button:hover {
  background-color: var(--btn-bg-hover);
  transform: translateY(-5px);
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}

@media (max-width: 900px) {
  .container {
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }

  .left,
  .right {
    width: 100%;
    max-width: 640px;
  }
}
</style>
