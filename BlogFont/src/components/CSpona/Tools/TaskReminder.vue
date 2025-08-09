<template>
    <div class="page-container">
        <!-- 顶部导航栏 -->
        <header class="header">
            <div class="header-container">
                <div class="header-left">
                    <i class="fa fa-bell header-icon"></i>
                    <h1 class="header-title">任务提醒助手</h1>
                </div>
                <div class="header-right md:flex">
                    <button @click="showSettingsModal = true" class="header-button">
                        <i class="fa fa-cog header-button-icon"></i>
                        <span class="header-button-text">设置</span>
                    </button>
                    <button @click="showAboutModal = true" class="header-button">
                        <i class="fa fa-info-circle header-button-icon"></i>
                        <span class="header-button-text">关于</span>
                    </button>
                </div>
            </div>
            <!-- 移动端菜单 -->
            <div class="mobile-menu" v-if="mobileMenuOpen">
                <div class="mobile-menu-container">
                    <button @click="handleMobileSettingClick" class="mobile-menu-button">
                        <i class="fa fa-cog mobile-menu-button-icon"></i>
                        <span class="mobile-menu-button-text">设置</span>
                    </button>
                    <button @click="handleMobileAboutClick" class="mobile-menu-button">
                        <i class="fa fa-info-circle mobile-menu-button-icon"></i>
                        <span class="mobile-menu-button-text">关于</span>
                    </button>
                </div>
            </div>
        </header>

        <!-- 主内容区 -->
        <main class="main-content">
            <!-- 添加任务卡片 -->
            <section class="card add-reminder-card">
                <h2 class="card-title">
                    <i class="fa fa-plus-circle card-icon"></i>添加新提醒
                </h2>
                <form @submit.prevent="handleAddReminder" class="form">
                    <div class="form-group">
                        <label for="taskTitle" class="form-label">提醒标题</label>
                        <input type="text" id="taskTitle" v-model="newReminder.title" placeholder="例如：取快递"
                            class="form-input" required />
                    </div>
                    <div class="form-group">
                        <label for="taskContent" class="form-label">提醒内容</label>
                        <textarea id="taskContent" v-model="newReminder.content" placeholder="例如：明天下午3点05到小区门口取快递"
                            class="form-input form-textarea" required></textarea>
                    </div>
                    <div class="form-group">
                        <label for="reminderTime" class="form-label">提醒时间</label>
                        <input type="datetime-local" id="reminderTime" v-model="newReminder.time" class="form-input"
                            required />
                    </div>
                    <button type="submit" class="form-button">
                        <i class="fa fa-save form-button-icon"></i>保存提醒
                    </button>
                </form>
            </section>

            <!-- 任务列表 -->
            <section class="task-list-section">
                <div class="task-list-header">
                    <h2 class="task-list-title">
                        <i class="fa fa-list-alt task-list-icon"></i>我的提醒
                    </h2>
                    <div class="task-list-actions">
                        <button @click="deleteAllReminders" class="delete-all-button" :disabled="reminders.length === 0"
                            :class="{ disabled: reminders.length === 0 }">
                            <i class="fa fa-trash delete-all-icon"></i>清空所有
                        </button>
                    </div>
                </div>

                <!-- 任务列表容器 -->
                <div id="remindersList" class="task-list">
                    <template v-if="reminders.length === 0">
                        <div class="empty-task">
                            <i class="fa fa-calendar-o empty-task-icon"></i>
                            <p class="empty-task-text">
                                暂无提醒任务，添加你的第一个提醒吧！
                            </p>
                        </div>
                    </template>

                    <template v-else>
                        <div v-for="(reminder, index) in sortedReminders" :key="reminder.Id" class="task-item"
                            :class="{ 'past': isPast(reminder.RemindAt) }">
                            <div class="task-item-content">
                                <div class="task-item-text">
                                    <h3 class="task-title" :class="{ 'completed': isPast(reminder.RemindAt) }">
                                        {{ reminder.Title }}
                                    </h3>
                                    <p class="task-content" :class="{ 'completed': isPast(reminder.RemindAt) }">
                                        {{ reminder.Content }}
                                    </p>
                                    <p class="task-time">
                                        <i class="fa fa-clock-o task-time-icon"></i>
                                        {{ formatDateTime(reminder.RemindAt) }}
                                        <span v-if="isPast(reminder.RemindAt)" class="expired-tag">已过期</span>
                                    </p>
                                </div>
                                <div class="task-item-actions">
                                    <button class="edit-task-button" @click="prepareEditReminder(reminder)">
                                        <i class="fa fa-edit"></i>
                                    </button>
                                    <button class="delete-task-button" @click="deleteReminder(reminder.Id)">
                                        <!-- 修正为 Id -->
                                        <i class="fa fa-trash"></i>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </template>
                </div>
            </section>
        </main>

        <!-- 编辑提醒模态框 -->
        <div id="editModal" class="modal" v-if="editingReminder" @click.self="editingReminder = null">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">编辑提醒</h3>
                    <button @click="editingReminder = null" class="modal-close-button" aria-label="关闭">
                        <i class="fa fa-times"></i>
                    </button>
                </div>
                <form @submit.prevent="saveEditReminder" class="modal-form">
                    <div class="modal-form-group">
                        <label for="editTitle" class="modal-form-label">提醒标题</label>
                        <input type="text" id="editTitle" v-model="editingReminder.title" class="modal-form-input"
                            required />
                    </div>
                    <div class="modal-form-group">
                        <label for="editContent" class="modal-form-label">提醒内容</label>
                        <textarea id="editContent" v-model="editingReminder.content"
                            class="modal-form-input modal-form-textarea" required></textarea>
                    </div>
                    <div class="modal-form-group">
                        <label for="editTime" class="modal-form-label">提醒时间</label>
                        <input type="datetime-local" id="editTime" v-model="editingReminder.time"
                            class="modal-form-input" required />
                    </div>
                    <div class="modal-form-buttons">
                        <button type="submit" class="modal-form-button">保存修改</button>
                    </div>
                </form>
            </div>
        </div>

        <!-- 邮件设置模态框 -->
        <div id="settingsModal" class="modal" v-if="showSettingsModal" @click.self="showSettingsModal = false">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">邮件通知设置</h3>
                    <button @click="showSettingsModal = false" class="modal-close-button" aria-label="关闭">
                        <i class="fa fa-times"></i>
                    </button>
                </div>
                <form @submit.prevent="saveEmailSettings" class="modal-form">
                    <div class="modal-form-group">
                        <label for="Email" class="modal-form-label">邮箱地址</label>
                        <input type="email" id="Email" v-model="emailSettings.Email" placeholder="your@email.com"
                            class="modal-form-input" required />
                    </div>
                    <div class="modal-form-group">
                        <label for="smtpPassword" class="modal-form-label">
                            SMTP授权码
                            <span class="modal-form-note">(不是邮箱密码)</span>
                        </label>
                        <input type="password" id="smtpPassword" v-model="emailSettings.smtpPassword"
                            placeholder="请输入邮箱SMTP授权码" class="modal-form-input" required />
                        <p class="modal-form-help">
                            如何获取授权码？<a href="https://service.mail.qq.com/cgi-bin/help?subtype=1&&id=28&&no=1001256"
                                target="_blank" class="modal-form-help-link">点击查看教程</a>
                        </p>
                    </div>
                    <button type="submit" class="modal-form-button">保存设置</button>
                </form>
            </div>
        </div>

        <!-- 关于模态框 -->
        <div id="aboutModal" class="modal" v-if="showAboutModal" @click.self="showAboutModal = false">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">关于提醒助手</h3>
                    <button @click="showAboutModal = false" class="modal-close-button" aria-label="关闭">
                        <i class="fa fa-times"></i>
                    </button>
                </div>
                <div class="modal-body">
                    <p class="modal-body-text">
                        提醒助手是一个简单实用的工具，可以帮助你管理日常任务和重要事项。
                    </p>
                    <p class="modal-body-text">
                        设置提醒后，系统会在指定时间通过邮箱发送通知，确保你不会错过重要事情。
                    </p>
                    <div class="modal-features">
                        <h4 class="modal-features-title">功能特点：</h4>
                        <ul class="modal-features-list">
                            <li>添加自定义提醒任务（含标题和详细内容）</li>
                            <li>设置具体提醒时间</li>
                            <li>邮箱通知提醒</li>
                            <li>管理所有提醒任务（查看、编辑、删除）</li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>

        <!-- 通知提示 -->
        <div id="notification" class="notification" :class="{
            show: notification.show,
            success: notification.type === 'success',
            error: notification.type === 'error',
            warning: notification.type === 'warning',
            reminder: notification.type === 'reminder',
            info: notification.type === 'info',
        }">
            <i class="fa" :class="notificationIcon"></i>
            <span class="notification-text">{{ notification.message }}</span>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from "vue";
import { getCurrentInstance } from "vue";

// 获取全局实例和配置
const instance = getCurrentInstance();
const { proxy } = instance || {};
const URL = instance?.appContext.config.globalProperties.URL;

// 响应式数据
const reminders = ref([]);
const emailSettings = ref({
    Email: "",
    smtpPassword: "",
});
const newReminder = ref({
    title: "", // 新增title字段
    content: "", // 原description改为content
    time: "",
});
const showSettingsModal = ref(false);
const showAboutModal = ref(false);
const mobileMenuOpen = ref(false);
const notification = ref({
    show: false,
    message: "",
    type: "info", // 'success', 'error', 'warning', 'reminder', 'info'
});
const editingReminder = ref(null); // 用于编辑的提醒对象

// 初始化
onMounted(() => {
    fetchEmailSettings();
    // 设置默认提醒时间为当前时间加1小时
    setDefaultReminderTime();
    // 启动提醒检查器
    startReminderChecker();
    // 请求通知权限
    if (
        Notification.permission !== "granted" &&
        Notification.permission !== "denied"
    ) {
        Notification.requestPermission();
    }
    // 从后端加载提醒
    fetchReminders();
});

const fetchEmailSettings = async () => {
    if (!userId.value) {
        showNotification("未获取到用户信息，无法加载邮箱设置", "error");
        return;
    }

    try {
        const response = await fetch(`${URL}/users/${userId.value}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (!response.ok) {
            throw new Error(`获取邮箱设置失败: ${response.status}`);
        }

        const result = await response.json();
        emailSettings.value = {
            Email: result.data.Email || "",
            // 密码/授权码只在编辑时使用，显示时可做隐藏处理
            smtpPassword: result.data.EmailSmtp || "",
        };
    } catch (error) {
        console.error("获取邮箱设置失败:", error);
        // 不显示错误通知，避免干扰用户，首次使用时可能无设置
    }
};



// 计算属性：按时间排序的提醒列表
const sortedReminders = computed(() => {
    return [...reminders.value].sort((a, b) => {
        // 先判断两个提醒是否过期
        const aIsPast = isPast(a.RemindAt);
        const bIsPast = isPast(b.RemindAt);

        // 规则1：未过期的排在前面，已过期的排在后面
        if (!aIsPast && bIsPast) return -1; // a未过期，b已过期 → a在前
        if (aIsPast && !bIsPast) return 1;  // a已过期，b未过期 → b在前

        // 规则2：同状态（都未过期或都已过期）则按时间升序排列
        return new Date(a.RemindAt) - new Date(b.RemindAt);
    });
});

// 计算属性：通知图标
const notificationIcon = computed(() => {
    switch (notification.value.type) {
        case "success":
            return "fa-check-circle";
        case "error":
            return "fa-exclamation-circle";
        case "warning":
            return "fa-exclamation-triangle";
        case "reminder":
            return "fa-bell";
        default:
            return "fa-info-circle";
    }
});

// 获取当前用户ID
const userId = computed(() => localStorage.getItem("userId"));

// 设置默认提醒时间
const setDefaultReminderTime = () => {
    const defaultTime = new Date();
    defaultTime.setHours(defaultTime.getHours() + 1);
    defaultTime.setMinutes(0);
    newReminder.value.time = defaultTime.toISOString().slice(0, 16);
};

// 提醒管理方法 - 从后端获取提醒
const fetchReminders = async () => {
    if (!userId.value) {
        showNotification("未获取到用户信息", "error");
        return;
    }

    try {
        const response = await fetch(`${URL}/reminder/${userId.value}`);
        if (!response.ok) {
            throw new Error(`获取提醒失败: ${response.status}`);
        }
        const result = await response.json();
        reminders.value = result.data || [];
      reminders.value.forEach(reminder => {
      reminder.isTriggered = false; // 每个提醒默认未触发
    });
    } catch (error) {
        console.error("获取提醒失败:", error);
        showNotification("加载提醒失败", "error");
    }
};

// 添加提醒
const handleAddReminder = async () => {
    if (
        !newReminder.value.title ||
        !newReminder.value.content ||
        !newReminder.value.time
    ) {
        showNotification("请填写完整的提醒信息（标题、内容和时间）", "warning");
        return;
    }

    if (!userId.value) {
        showNotification("未获取到用户信息", "error");
        return;
    }

    try {
        // 转换时间格式为后端要求的"2006-01-02 15:04:05"
        const formattedTime = new Date(newReminder.value.time)
            .toISOString()
            .slice(0, 19)
            .replace("T", " ");

        const formData = new FormData();
        formData.append("user_id", userId.value);
        formData.append("title", newReminder.value.title);
        formData.append("content", newReminder.value.content);
        formData.append("remind_at", formattedTime);

        const response = await fetch(`${URL}/reminder`, {
            method: "POST",
            body: formData,
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "创建提醒失败");
        }

        const result = await response.json();
        reminders.value.push({
      ...result.data, // 后端返回的新提醒对象
      isTriggered: false // 确保新提醒未触发
    });
        showNotification("提醒已添加成功！", "success");
        
    fetchReminders();

        // 重置表单
        newReminder.value.title = "";
        newReminder.value.content = "";
        setDefaultReminderTime();
    } catch (error) {
        console.error("添加提醒失败:", error);
        showNotification(error.message, "error");
    }
};

// 准备编辑提醒
const prepareEditReminder = (reminder) => {
    editingReminder.value = {
        ...reminder,
        time: new Date(reminder.RemindAt).toISOString().slice(0, 16)
    };
    // 自动填充标题和内容（直接从 reminder 中获取）
    editingReminder.value.title = reminder.Title;  // 填入原标题
    editingReminder.value.content = reminder.Content;  // 填入原内容
    try {
        const date = new Date(reminder.RemindAt);
        editingReminder.value.time = date.toISOString().slice(0, 16);
    } catch (error) {
        console.error("编辑时时间解析失败:", error);
        const defaultTime = new Date();
        defaultTime.setHours(defaultTime.getHours() + 1);
        editingReminder.value.time = defaultTime.toISOString().slice(0, 16);
    }
};

// 保存编辑后的提醒
const saveEditReminder = async () => {
    if (!editingReminder.value || !editingReminder.value.Id) {
        // 修正为 Id
        showNotification("无效的提醒数据", "error");
        return;
    }

    try {
        // 转换时间格式为后端要求的格式
        const formattedTime = new Date(editingReminder.value.time)
            .toISOString()
            .slice(0, 19)
            .replace("T", " ");
        const formData = new FormData();
        formData.append("title", editingReminder.value.title);
        formData.append("content", editingReminder.value.content);
        formData.append("remind_at", formattedTime);
        // 请求路径使用 Id 而非 id
        const response = await fetch(
            `${URL}/reminder/${editingReminder.value.Id}`,
            {
                method: "PUT",
                body: formData,
            }
        );

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "更新提醒失败");
        }

        const result = await response.json();
        // 更新本地列表
        const index = reminders.value.findIndex(
            (r) => r.id === editingReminder.value.id
        );
        if (index !== -1) {
            reminders.value[index] = result.data;
        }

        showNotification("提醒已更新成功！", "success");
        editingReminder.value = null; // 关闭编辑状态
        
    fetchReminders();
    } catch (error) {
        console.error("更新提醒失败:", error);
        showNotification(error.message, "error");
    }
};

// 删除提醒
const deleteReminder = async (id) => {
    if (!confirm("确定要删除这个提醒吗？")) {
        return;
    }

    try {
        const response = await fetch(`${URL}/reminder/${id}`, { method: "DELETE" });
        if (!response.ok) throw new Error("删除失败");
        // 过滤时使用 Id 而非 id
        reminders.value = reminders.value.filter((r) => r.Id !== id);
        showNotification("提醒已删除", "info");
        
    fetchReminders();
    } catch (error) {
        console.error("删除提醒失败:", error);
        showNotification(error.message, "error");
    }
};

// 删除所有提醒
const deleteAllReminders = async () => {
    if (reminders.value.length === 0) {
        showNotification("没有可删除的提醒", "info");
        return;
    }

    if (!confirm("确定要删除所有提醒吗？此操作不可恢复。")) {
        return;
    }

    try {
        // 循环删除所有提醒
        for (const reminder of reminders.value) {
            await fetch(`${URL}/reminder/${reminder.Id}`, { method: "DELETE" }); // 修正为 Id
        }
        reminders.value = [];
        showNotification("所有提醒已删除", "info");
        
    fetchReminders();
    } catch (error) {
        console.error("删除所有提醒失败:", error);
        showNotification("删除提醒失败", "error");
    }
};

const encrypt = (str) => btoa(str); // Base64编码
const decrypt = (str) => atob(str); // Base64解码

watch(showSettingsModal, (newValue) => {
    if (newValue) {
        fetchEmailSettings();
    }
});
// 邮箱设置保存
const saveEmailSettings = async () => {
    if (!emailSettings.value.Email || !emailSettings.value.smtpPassword) {
        showNotification("请填写完整的邮箱设置", "warning");
        return;
    }

    try {
        const formData = new FormData();
        formData.append("email", emailSettings.value.Email);
        formData.append("email_smtp", emailSettings.value.smtpPassword);

        const response = await fetch(`${URL}/users/${userId.value}`, {
            method: "PUT",
            body: formData,
        });

        if (!response.ok) {
            throw new Error(`保存失败: ${response.status}`);
        }

        // 保存成功后重新从后端获取最新设置
        await fetchEmailSettings();
        showSettingsModal.value = false;
        showNotification("邮箱设置已保存", "success");
    } catch (error) {
        console.error("保存邮箱设置失败:", error);
        showNotification("保存邮箱设置失败，请重试", "error");
    }
};
// 提醒检查
const checkReminders = () => {
    const now = new Date();
    const currentTime = now.getTime();
    reminders.value.forEach((reminder) => {
        const reminderTime = new Date(reminder.RemindAt).getTime(); // 修正为 RemindAt
        if (
            Math.abs(reminderTime - currentTime) < 60000 && 
            !reminder.isTriggered
        ) {
            triggerReminder(reminder);
            // 标记为已触发（更新引用）
            reminder.isTriggered = true;
        }
    });
};

const startReminderChecker = () => {
    // 每分钟检查一次
    setInterval(checkReminders, 60000);
    // 立即检查一次
    checkReminders();
};

const triggerReminder = (reminder) => {
    // 显示浏览器通知
    if (Notification.permission === "granted") {
        try {
            new Notification("提醒助手", {
                body: `${reminder.Title}: ${reminder.Content}`,
                icon: "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCA1MTIgNTEyIj48cGF0aCBmaWxsPSIjM0I4MkY2IiBkPSJNMzg0IDM2OGgtMjU2Yy0xNy43IDAtMzIgMTQuMy0zMiAzMnY0OGMxNy43IDAgMzIgMTQuMyAzMiAzMnYxNmMwIDE3LjcgMTQuMyAzMiAzMiAzMmgxOTJjMTcuNyAwIDMyLTE0LjMgMzItMzJ2LTE2YzAgLTE3LjctMTQuMy0zMi0zMi0zMlY0MDAgYzAtMTcuNy0xNC4zLTMyLTMyLTMyem0tMjU2LTMyMHYxMjhjMCAxNy43IDE0LjMgMzIgMzIgMzJoMjU2YzE3LjcgMCAzMi0xNC4zIDMyLTMyVjQ4Yy0xNy43IDAtMzIgMTQuMy0zMiAzMnptMjU2IDI1NnoiLz48L3N2Zz4=",
            });
        } catch (error) {
            console.error("显示通知失败:", error);
        }
    }

    // 显示页面内通知
    showNotification(`提醒: ${reminder.Title} - ${reminder.Content}`, "reminder");

    // 发送邮件通知
    if (emailSettings.value.Email && emailSettings.value.smtpPassword) {
        sendEmailNotification(reminder);
    } else {
        setTimeout(() => {
            showNotification("请先设置邮箱以接收邮件提醒", "warning");
        }, 3000);
    }
};

const sendEmailNotification = async (reminder) => {
    try {
        if (!emailSettings.value.Email ||!emailSettings.value.smtpPassword) {
            throw new Error("邮箱设置不完整");
        }
        await fetch(`${URL}/reminder/${reminder.Id}/email`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: emailSettings.value.Email,
                smtpPassword: emailSettings.value.smtpPassword
            }),
        });
        showNotification("邮件提醒已发送", "success");
    } catch (error) {
        console.error("发送邮件失败:", error);
        showNotification("发送邮件提醒失败", "error");
    }
};


// 工具方法
const showNotification = (message, type = "info") => {
    notification.value = {
        show: true,
        message,
        type,
    };

    setTimeout(() => {
        notification.value.show = false;
    }, 3000);
};

const formatDateTime = (dateString) => {
    try {
        // 先判断是否为空或无效值
        if (!dateString || dateString === "Invalid Date") {
            return "未设置时间";
        }

        const date = new Date(dateString);
        // 检查日期是否有效
        if (isNaN(date.getTime())) {
            return "时间格式错误";
        }

        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        const hours = String(date.getHours()).padStart(2, "0");
        const minutes = String(date.getMinutes()).padStart(2, "0");

        return `${year}-${month}-${day} ${hours}:${minutes}`;
    } catch (error) {
        console.error("时间格式化失败:", error, "原始值:", dateString);
        return "解析失败";
    }
};

// 移动端菜单处理
const handleMobileSettingClick = () => {
    mobileMenuOpen.value = false;
    showSettingsModal.value = true;
};

const handleMobileAboutClick = () => {
    mobileMenuOpen.value = false;
    showAboutModal.value = true;
};

// 在脚本末尾添加 isPast 函数
const isPast = (dateString) => {
    try {
        if (!dateString) return false;
        const reminderTime = new Date(dateString).getTime();
        const now = new Date().getTime();
        return reminderTime < now; // 提醒时间早于当前时间 → 已过期
    } catch (error) {
        console.error('过期判断失败:', error);
        return false;
    }
};

// 显式暴露函数（确保模板可访问）
defineExpose({
    isPast,
    formatDateTime
});
</script>

<style scoped>
/* 全局样式 */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

.page-container {
    background-color: #f3f4f6;
    min-height: 100vh;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica,
        Arial, sans-serif;
    color: #1f2937;
    padding-bottom: 60px;
}

/* 头部样式 */
.header {
    background-color: #ffffff;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 50;
    transition: all 0.3s ease;
}

.header-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    height: 60px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 10px;
}

.header-icon {
    color: #3b82f6;
    font-size: 24px;
}

.header-title {
    font-size: 18px;
    font-weight: bold;
    color: #3b82f6;
}

.header-right {
    display: flex;
    align-items: center;
    gap: 16px;
}

.header-button {
    display: flex;
    align-items: center;
    gap: 5px;
    padding: 6px 12px;
    border-radius: 6px;
    background: none;
    border: none;
    cursor: pointer;
    transition: all 0.3s ease;
    color: #1f2937;
}

.header-button:hover {
    background-color: #e5e7eb;
}

.header-button-icon {
    font-size: 16px;
}

.header-button-text {
    font-size: 14px;
}

.header-menu-button {
    background: none;
    border: none;
    font-size: 20px;
    cursor: pointer;
    color: #1f2937;
}

/* 移动端菜单 */
.mobile-menu {
    background-color: #ffffff;
    border-top: 1px solid #e5e7eb;
}

.mobile-menu-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 10px 20px;
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.mobile-menu-button {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 15px;
    border-radius: 6px;
    background: none;
    border: none;
    cursor: pointer;
    transition: all 0.3s ease;
    text-align: left;
    color: #1f2937;
}

.mobile-menu-button:hover {
    background-color: #e5e7eb;
}

.mobile-menu-button-icon {
    font-size: 18px;
}

.mobile-menu-button-text {
    font-size: 16px;
}

/* 主内容区 */
.main-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 80px 20px 20px;
}

/* 卡片样式 */
.card {
    background-color: #ffffff;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 32px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
}

.card:hover {
    box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
}

.add-reminder-card {
    max-width: 768px;
    margin-left: auto;
    margin-right: auto;
}

.card-title {
    font-size: 20px;
    font-weight: bold;
    margin-bottom: 16px;
    display: flex;
    align-items: center;
    color: #1f2937;
}

.card-icon {
    color: #3b82f6;
    margin-right: 8px;
    font-size: 22px;
}

/* 表单样式 */
.form {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.form-label {
    font-size: 14px;
    font-weight: 500;
    color: #6b7280;
}

.form-input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 16px;
    transition: all 0.3s ease;
}

.form-input:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.form-button {
    width: 100%;
    background-color: #3b82f6;
    color: #ffffff;
    border: none;
    border-radius: 8px;
    padding: 10px 16px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.form-button:hover {
    background-color: #2563eb;
}

.form-button-icon {
    font-size: 18px;
}

/* 任务列表样式 */
.task-list-section {
    max-width: 768px;
    margin-left: auto;
    margin-right: auto;
}

.task-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
}

.task-list-title {
    font-size: 20px;
    font-weight: bold;
    display: flex;
    align-items: center;
    color: #1f2937;
}

.task-list-icon {
    color: #3b82f6;
    margin-right: 8px;
    font-size: 22px;
}

.task-list-actions {
    display: flex;
    gap: 8px;
}

.delete-all-button {
    display: flex;
    align-items: center;
    gap: 5px;
    background: none;
    border: none;
    color: #ef4444;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    padding: 4px 8px;
    border-radius: 4px;
}

.delete-all-button:hover {
    background-color: rgba(239, 68, 68, 0.1);
}

.delete-all-button.disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.delete-all-icon {
    font-size: 14px;
}

.task-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

/* 任务项样式 */
.task-item {
    background-color: #ffffff;
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
}

.task-item:hover {
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
}

.task-item.past {
    opacity: 0.7;
}

.task-item-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
}

.task-item-text {
    flex: 1;
    margin-right: 16px;
}

.task-description {
    font-size: 16px;
    font-weight: 500;
    color: #1f2937;
}

.task-description.completed {
    color: #6b7280;
    text-decoration: line-through;
}

.task-time {
    font-size: 14px;
    color: #6b7280;
    margin-top: 4px;
    display: flex;
    align-items: center;
    gap: 4px;
}

.task-time-icon {
    font-size: 14px;
}

.expired-tag {
    font-size: 12px;
    background-color: #e5e7eb;
    color: #6b7280;
    padding: 2px 6px;
    border-radius: 4px;
    margin-left: 8px;
}

.delete-task-button {
    background: none;
    border: none;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.3s ease;
    padding: 4px;
    border-radius: 4px;
}

.delete-task-button:hover {
    color: #ef4444;
    background-color: rgba(239, 68, 68, 0.1);
}

/* 空任务样式 */
.empty-task {
    background-color: #ffffff;
    border-radius: 8px;
    padding: 48px 20px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    text-align: center;
    color: #6b7280;
}

.empty-task-icon {
    font-size: 64px;
    opacity: 0.3;
    margin-bottom: 12px;
}

.empty-task-text {
    font-size: 16px;
}

/* 页脚样式 */
.footer {
    background-color: #1f2937;
    color: #ffffff;
    padding: 24px 0;
    margin-top: auto;
}

.footer-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    text-align: center;
}

.footer-text {
    font-size: 14px;
    opacity: 0.8;
}

/* 模态框样式 */
.modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
}

.modal-content {
    background-color: #ffffff;
    border-radius: 12px;
    width: 100%;
    max-width: 500px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    overflow: hidden;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #e5e7eb;
}

.modal-title {
    font-size: 18px;
    font-weight: bold;
    color: #1f2937;
}

.modal-close-button {
    background: none;
    border: none;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.3s ease;
    padding: 4px;
    border-radius: 4px;
    font-size: 18px;
}

.modal-close-button:hover {
    color: #1f2937;
    background-color: #e5e7eb;
}

.modal-body {
    padding: 20px;
    color: #1f2937;
}

.modal-body-text {
    margin-bottom: 16px;
    line-height: 1.5;
}

.modal-features {
    margin-top: 20px;
}

.modal-features-title {
    font-weight: 600;
    margin-bottom: 8px;
}

.modal-features-list {
    padding-left: 20px;
    line-height: 1.5;
}

.modal-features-list li {
    margin-bottom: 6px;
}

.modal-form {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.modal-form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.modal-form-label {
    font-size: 14px;
    font-weight: 500;
    color: #6b7280;
}

.modal-form-note {
    font-size: 12px;
    color: #6b7280;
    margin-left: 4px;
}

.modal-form-input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 16px;
    transition: all 0.3s ease;
}

.modal-form-input:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.modal-form-help {
    font-size: 12px;
    color: #6b7280;
    margin-top: 4px;
}

.modal-form-help-link {
    color: #3b82f6;
    text-decoration: none;
}

.modal-form-help-link:hover {
    text-decoration: underline;
}

.modal-form-button {
    width: 100%;
    background-color: #3b82f6;
    color: #ffffff;
    border: none;
    border-radius: 8px;
    padding: 10px 16px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
}

.modal-form-button:hover {
    background-color: #2563eb;
}

/* 通知样式 */
.notification {
    position: fixed;
    bottom: 20px;
    right: 20px;
    padding: 12px 16px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    display: flex;
    align-items: center;
    gap: 8px;
    color: #ffffff;
    z-index: 1000;
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.3s ease;
    max-width: 300px;
}

.notification.show {
    opacity: 1;
    transform: translateY(0);
}

.notification i {
    font-size: 18px;
}

.notification-text {
    font-size: 14px;
    line-height: 1.4;
}

.notification.success {
    background-color: #10b981;
}

.notification.error {
    background-color: #ef4444;
}

.notification.warning {
    background-color: #f59e0b;
}

.notification.reminder {
    background-color: #3b82f6;
}

.notification.info {
    background-color: #1f2937;
}

/* 响应式样式 */
@media (max-width: 768px) {
    .header-title {
        font-size: 16px;
    }

    .card {
        padding: 16px;
    }

    .card-title,
    .task-list-title {
        font-size: 18px;
    }

    .form-input,
    .modal-form-input {
        font-size: 14px;
        padding: 8px 10px;
    }

    .form-button,
    .modal-form-button {
        font-size: 14px;
        padding: 8px 12px;
    }

    .task-description {
        font-size: 14px;
    }

    .task-time {
        font-size: 12px;
    }
}

/* 辅助类 */
.hidden {
    display: none !important;
}

/* 辅助类 - 确保大屏显示按钮容器，小屏隐藏 */
.md:flex {
    display: flex !important;
}

/* 小屏（≤768px）时：隐藏按钮容器，显示三条杠 */
@media (max-width: 768px) {
    .md:hidden {
        display: block !important;
        /* 三条杠按钮显示 */
    }

    .md:flex {
        display: none !important;
        /* 按钮容器隐藏 */
    }
}
</style>
