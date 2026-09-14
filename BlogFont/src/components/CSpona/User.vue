<template>
    <!-- 背景图片 -->
    <img :src="userAvatar" alt="背景图" class="bg-img" />
    <div class="user-profile-container">
        <!-- 左侧用户信息展示区（保持不变） -->
        <div class="user-info-section">
            <div class="profile-card">
                <!-- 加载状态 -->
                <div v-if="loading" class="loading-state">
                    <div class="spinner"></div>
                    <p>加载用户信息中...</p>
                </div>

                <!-- 错误状态 -->
                <div v-if="error" class="error-state">
                    <i class="fas fa-exclamation-circle"></i>
                    <p>{{ error }}</p>
                    <button @click="fetchUserInfo" class="retry-btn">重试</button>
                </div>

                <!-- 用户信息（加载成功后显示） -->
                <div v-else class="user-info-content">
                    <!-- 圆形头像 -->
                    <div class="avatar-container">
                        <img :src="userAvatar" alt="用户头像" class="avatar" @error="handleAvatarError" />
                    </div>

                    <h2 class="username">{{ userInfo.Username || "未设置" }}</h2>
                    <p class="user-role" :class="roleClass">{{ roleText }}</p>

                    <div class="info-grid">
                        <div class="info-item">
                            <span class="info-label">用户ID</span>
                            <span class="info-value">{{ userInfo.Id || "N/A" }}</span>
                        </div>

                        <div class="info-item">
                            <span class="info-label">邮箱</span>
                            <span class="info-value">{{ userInfo.Email || "未设置" }}</span>
                        </div>

                        <div class="info-item">
                            <span class="info-label">注册时间</span>
                            <span class="info-value">{{ formattedDate || "N/A" }}</span>
                        </div>

                        <div class="info-item">
                            <span class="info-label">文章数量</span>
                            <span class="info-value">{{ userInfo.ArticleCount || 0 }}</span>
                        </div>

                        <div class="info-item full-width">
                            <span class="info-label">个人简介</span>
                            <p class="info-value description">
                                {{ userInfo.Description || "暂无简介" }}
                            </p>
                        </div>
                    </div>

                    <button class="edit-btn" @click="showEditDialog = true">
                        <i class="fas fa-edit"></i> 修改资料
                    </button>
                </div>
            </div>
        </div>

        <!-- 右侧编辑区域（优化后） -->
        <div class="edit-section">
            <div class="edit-card">
                <h2 class="edit-title">在漫游地(Wanderland)发布文章</h2>
                <div class="edit-content">
                    <Edit />
                </div>
            </div>
        </div>

        <!-- 修改资料弹窗（优化后：文件上传+布局留白） -->
        <div class="modal-overlay" v-if="showEditDialog" @click="closeModal">
            <div class="modal" @click.stop>
                <div class="modal-header">
                    <h3>修改个人资料</h3>
                    <button class="close-btn" @click="showEditDialog = false">
                        <i class="fas fa-times"></i>
                    </button>
                </div>

                <div class="modal-body">
                    <form @submit.prevent="handleUpdateProfile">
                        <!-- 用户名 -->
                        <div class="form-group">
                            <label for="username">用户名</label>
                            <input type="text" id="username" v-model="updatedUser.username" required
                                class="form-control" />
                        </div>

                        <!-- 邮箱 -->
                        <div class="form-group">
                            <label for="email">邮箱</label>
                            <input type="email" id="email" v-model="updatedUser.email" required class="form-control" />
                        </div>

                        <!-- 个人简介 -->
                        <div class="form-group">
                            <label for="description">个人简介</label>
                            <textarea id="description" v-model="updatedUser.description" rows="3" class="form-control"
                                placeholder="请输入个人简介..." style="resize: vertical"></textarea>
                        </div>

                        <!-- 头像上传（文件选择样式） -->
                        <div class="form-group avatar-upload-group">
                            <label>头像设置</label>
                            <div class="avatar-upload-container">
                                <!-- 头像预览 -->
                                <div class="avatar-preview">
                                    <img :src="avatarPreviewUrl" alt="预览头像" class="preview-img" />
                                </div>

                                <!-- 文件上传按钮 -->
                                <div class="file-upload-wrapper">
                                    <input type="file" id="avatar" class="file-upload-input" accept="image/*"
                                        @change="handleAvatarUpload" />
                                    <label for="avatar" class="file-upload-btn">
                                        <i class="fas fa-upload"></i> 选择头像图片
                                    </label>
                                    <p class="file-info" v-if="selectedFileName">
                                        已选择: {{ selectedFileName }}
                                    </p>
                                    <p class="file-hint">支持JPG、PNG格式，建议尺寸200x200px</p>
                                </div>
                            </div>
                        </div>

                        <!-- 操作按钮 -->
                        <div class="form-actions">
                            <button type="button" class="cancel-btn" @click="showEditDialog = false">
                                取消
                            </button>
                            <button type="submit" class="save-btn" :disabled="updateLoading">
                                <span v-if="updateLoading"><i class="fas fa-spinner fa-spin"></i> 保存中...</span>
                                <span v-else>保存修改</span>
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { useUserStore } from "@/stores/user";
import Edit from "./Edit.vue";
import { getCurrentInstance } from "vue";

// 获取全局实例和配置
const instance = getCurrentInstance();
const { proxy } = instance || {};
const URL = instance?.appContext.config.globalProperties.URL;

// 初始化用户存储
const userStore = useUserStore();

// 状态管理
const loading = ref(true);
const error = ref("");
const updateLoading = ref(false);
const showEditDialog = ref(false);
const userInfo = reactive({});

// 头像上传相关状态
const avatarPreviewUrl = ref("");
const selectedFileName = ref("");
const selectedFile = ref(null);

// 从localStorage获取用户ID
const userId = computed(() => localStorage.getItem("userId"));

// 用于编辑的用户信息副本
const updatedUser = reactive({
    username: "",
    email: "",
    description: "",
    avatar: "",
});

// 图片URL处理函数
const getImageUrl = (imgName) => {
    if (!imgName) return `${proxy?.$imageBaseUrl}boli.jpg`;
    return `${proxy?.$imageBaseUrl}${imgName}`;
};

// 头像URL计算
const userAvatar = computed(() => getImageUrl(userInfo.Avatar));

// 处理头像加载失败
const handleAvatarError = (e) => {
    e.target.src = `${proxy?.$imageBaseUrl}boli.jpg`;
};

// 处理头像上传
const handleAvatarUpload = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const maxSize = 1 * 1024 * 1024;
    if (file.size > maxSize) {
        alert("头像大小不要超过1MB");
        e.target.value = "";
        return;
    }

    // 验证文件类型
    if (!file.type.startsWith("image/")) {
        alert("请选择图片文件");
        e.target.value = "";
        return;
    }

    // 保存文件信息
    selectedFile.value = file;
    selectedFileName.value = file.name;

    // 预览头像
    const reader = new FileReader();
    reader.onload = (event) => {
        avatarPreviewUrl.value = event.target.result;
    };
    reader.readAsDataURL(file);
};

// 格式化日期
const formattedDate = computed(() => {
    if (!userInfo.CreatedAt) return "";
    const date = new Date(userInfo.CreatedAt);
    return date.toLocaleString("zh-CN", {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
});

// 角色文本和样式
const roleText = computed(() => {
    switch (userInfo.RoleQx) {
        case "A":
            return "管理员";
        case "B":
            return "高级用户";
        case "C":
            return "普通用户";
        default:
            return "未知角色";
    }
});

const roleClass = computed(() => {
    switch (userInfo.RoleQx) {
        case "A":
            return "role-admin";
        case "B":
            return "role-premium";
        default:
            return "role-default";
    }
});

// 获取用户信息
const fetchUserInfo = async () => {
    try {
        loading.value = true;
        error.value = "";

        if (!userId.value) {
            throw new Error("未获取到用户ID，请重新登录");
        }

        // 调用API获取用户详情
        const response = await fetch(`${URL}/users/${userId.value}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (!response.ok) {
            throw new Error(`获取用户信息失败，状态码：${response.status}`);
        }

        const result = await response.json();
        Object.assign(userInfo, result.data || result);

        // 初始化编辑表单
        updatedUser.username = userInfo.Username || "";
        updatedUser.email = userInfo.Email || "";
        updatedUser.description = userInfo.Description || "";
        updatedUser.avatar = userInfo.Avatar || "";

        // 初始化头像预览
        avatarPreviewUrl.value = getImageUrl(userInfo.Avatar);
        selectedFileName.value = userInfo.Avatar || "";
    } catch (err) {
        error.value = err.message || "获取用户信息失败，请稍后重试";
        console.error("获取用户信息错误:", err);
    } finally {
        loading.value = false;
    }
};

// 保存用户资料修改
const handleUpdateProfile = async () => {
    try {
        updateLoading.value = true;

        // 准备更新数据（如果有新头像则使用新文件名）
        const formData = new FormData();
        // 添加文本字段
        formData.append("username", updatedUser.username);
        formData.append("email", updatedUser.email);
        formData.append("description", updatedUser.description);
        if (selectedFile.value) {
            formData.append("avatar", selectedFile.value);
        }

        // 调用API更新用户信息
        const response = await fetch(`${URL}/users/${userId.value}`, {
            method: "PUT",
            body: formData,
        });

        const data = await response.json();

        if (!response.ok || data.error) {
            throw new Error(data.error || `更新失败，状态码：${response.status}`);
        }

        // 更新本地用户信息
        Object.assign(userInfo, {
            Username: updatedUser.username,
            Email: updatedUser.email,
            Avatar: selectedFileName.value,
            Description: updatedUser.description,
        });

        // 更新Pinia store
        userStore.username = updatedUser.username;
        userStore.email = updatedUser.email;
        userStore.avatar = selectedFileName.value;

        // 更新localStorage
        localStorage.setItem("username", updatedUser.username);
        localStorage.setItem("email", updatedUser.email);
        localStorage.setItem("avatar", selectedFileName.value);

        // 关闭弹窗并提示成功
        showEditDialog.value = false;
        alert("资料更新成功！");
    } catch (err) {
        console.error("更新用户信息错误:", err);
        alert(err.message || "更新失败，请稍后重试");
    } finally {
        updateLoading.value = false;
    }
};

// 关闭模态框
const closeModal = (e) => {
    if (e.target === e.currentTarget) {
        showEditDialog.value = false;
    }
};

// 页面加载时获取用户信息
onMounted(() => {
    if (userId.value) {
        fetchUserInfo();
    } else {
        loading.value = false;
        error.value = "未检测到登录状态，请先登录";
    }
});
</script>

<style scoped>
.user-profile-container {
  position: relative;
  z-index: 1;
  display: flex;
  min-height: calc(100vh - 4.75rem);
  width: 100%;
  gap: 0;
  box-sizing: border-box;
}

.user-info-section {
  width: 38%;
  padding: 2.5rem 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
  box-sizing: border-box;
}

.profile-card {
  width: 100%;
  max-width: 520px;
  background: var(--panel-bg);
  color: var(--text-primary);
  border: 1px solid var(--panel-border);
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  padding: 2.75rem 2.5rem;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.edit-section {
  width: 62%;
  padding: 2.5rem 2.5rem 3rem;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  background-color: var(--panel-bg-soft);
  box-sizing: border-box;
}

.edit-card {
  width: 100%;
  max-width: 860px;
  background: var(--panel-bg);
  color: var(--text-primary);
  border: 1px solid var(--panel-border);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  padding: 2.75rem 2.5rem;
  margin-top: 1rem;
}

.edit-title {
  margin-top: 0;
  color: var(--text-primary);
  font-size: 1.75rem;
  padding-bottom: 1.1rem;
  border-bottom: 1px solid var(--panel-border);
}

.edit-content {
  margin-top: 1.75rem;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.55);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal {
  background-color: var(--panel-bg);
  color: var(--text-primary);
  border: 1px solid var(--panel-border);
  border-radius: 16px;
  width: 100%;
  max-width: 640px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  animation: slideIn 0.3s ease;
  overflow: hidden;
}

.modal-header {
  padding: 1.5rem 2rem;
  border-bottom: 1px solid var(--panel-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.35rem;
  color: var(--text-primary);
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: var(--text-secondary);
  transition: color 0.2s ease, background-color 0.2s ease;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #f87171;
  background-color: rgba(248, 113, 113, 0.12);
}

.modal-body {
  padding: 2rem;
}

.form-group {
  margin-bottom: 1.9rem;
  position: relative;
}

.form-group label {
  display: block;
  margin-bottom: 0.8rem;
  font-weight: 500;
  color: var(--text-primary);
  font-size: 1rem;
}

.form-control {
  width: 100%;
  box-sizing: border-box;
  padding: 0.95rem 1.1rem;
  border: 1px solid var(--input-border);
  border-radius: 8px;
  font-size: 1.05rem;
  transition: all 0.2s ease;
  background-color: var(--input-bg);
  color: var(--input-text);
  min-height: 48px;
}

.form-control:focus {
  outline: none;
  border-color: var(--btn-bg);
  box-shadow: 0 0 0 3px rgba(120, 190, 255, 0.2);
}

.avatar-upload-group {
  margin-top: 1rem;
}

.avatar-upload-container {
  display: flex;
  gap: 1.5rem;
  align-items: center;
  flex-wrap: wrap;
}

.avatar-preview {
  width: 128px;
  height: 128px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px dashed var(--input-border);
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--panel-bg-soft);
}

.preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-upload-wrapper {
  flex: 1;
  min-width: 250px;
}

.file-upload-input {
  display: none;
}

.file-upload-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.85rem 1.6rem;
  background-color: var(--btn-bg);
  color: var(--btn-text);
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-upload-btn:hover {
  background-color: var(--btn-bg-hover);
  transform: translateY(-1px);
}

.file-info {
  margin: 0.75rem 0 0.25rem;
  font-size: 0.95rem;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.file-hint {
  margin: 0.25rem 0 0;
  font-size: 0.85rem;
  color: var(--text-muted);
  line-height: 1.5;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2.5rem;
}

.cancel-btn {
  padding: 0.85rem 1.6rem;
  background-color: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--input-border);
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  background-color: var(--hover-bg);
  color: var(--text-primary);
}

.save-btn {
  padding: 0.85rem 1.6rem;
  background-color: var(--btn-bg);
  color: var(--btn-text);
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.save-btn:hover:not(:disabled) {
  background-color: var(--btn-bg-hover);
  transform: translateY(-1px);
}

.save-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

@media (max-width: 900px) {
  .user-profile-container {
    flex-direction: column;
  }

  .user-info-section,
  .edit-section {
    width: 100%;
  }

  .edit-card {
    padding: 2rem 1.5rem;
    margin-top: 0;
  }

  .avatar-upload-container {
    flex-direction: column;
    align-items: flex-start;
  }
}

.bg-img {
  position: fixed;
  inset: 0;
  width: 100%;
  height: 100vh;
  object-fit: cover;
  z-index: 0;
  opacity: 0.85;
}

.loading-state {
  text-align: center;
  padding: 2.5rem 0;
  color: var(--text-secondary);
}

.spinner {
  width: 40px;
  height: 40px;
  margin: 0 auto 1rem;
  border: 4px solid var(--panel-border);
  border-top: 4px solid var(--btn-bg);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.error-state {
  text-align: center;
  padding: 2rem 0;
  color: #f87171;
}

.error-state i {
  font-size: 2rem;
  margin-bottom: 1rem;
}

.retry-btn {
  margin-top: 1rem;
  padding: 0.65rem 1.2rem;
  background-color: var(--btn-bg);
  color: var(--btn-text);
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.retry-btn:hover {
  background-color: var(--btn-bg-hover);
}

.avatar-container {
  display: flex;
  justify-content: center;
  margin-bottom: 1.75rem;
}

.avatar {
  width: 148px;
  height: 148px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid var(--panel-border);
  transition: transform 0.3s ease;
}

.avatar:hover {
  transform: scale(1.04);
}

.username {
  text-align: center;
  margin: 0 0 0.6rem 0;
  color: var(--text-primary);
  font-size: 2rem;
}

.user-role {
  text-align: center;
  padding: 0.35rem 1.1rem;
  border-radius: 20px;
  font-weight: 500;
  margin: 0.5rem auto 1.75rem;
  display: inline-block;
}

.role-admin {
  background-color: rgba(120, 190, 255, 0.2);
  color: rgb(120, 190, 255);
}

.role-premium {
  background-color: rgba(251, 191, 36, 0.2);
  color: #fbbf24;
}

.role-default {
  background-color: rgba(74, 222, 128, 0.18);
  color: #4ade80;
}

html:not([data-theme="dark"]) .role-admin {
  background-color: #e0f2fe;
  color: #0284c7;
}

html:not([data-theme="dark"]) .role-premium {
  background-color: #fef3c7;
  color: #d97706;
}

html:not([data-theme="dark"]) .role-default {
  background-color: #dcfce7;
  color: #166534;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
  margin-bottom: 2rem;
}

.info-item {
  background-color: var(--panel-bg-soft);
  border: 1px solid var(--panel-border);
  padding: 1rem 1.15rem;
  border-radius: 10px;
  transition: transform 0.2s ease;
}

.info-item:hover {
  transform: translateY(-2px);
}

.full-width {
  grid-column: span 2;
}

.info-label {
  display: block;
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-bottom: 0.4rem;
  font-weight: 500;
}

.info-value {
  color: var(--text-primary);
  font-weight: 500;
  font-size: 1.05rem;
  word-break: break-word;
}

.description {
  margin: 0;
  line-height: 1.7;
}

.edit-btn {
  width: 100%;
  padding: 0.95rem;
  background-color: var(--btn-bg);
  color: var(--btn-text);
  border: none;
  border-radius: 8px;
  font-size: 1.05rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.edit-btn:hover {
  background-color: var(--btn-bg-hover);
  transform: translateY(-2px);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
