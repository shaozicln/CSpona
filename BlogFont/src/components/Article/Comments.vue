<template>
  <div class="comment-container">
    <div class="comment-section">
      <!-- 评论输入框 -->
      <div class="comment-input">
        <el-input v-model="newComment" type="textarea" :rows="4" placeholder="写下你的评论..."
          class="comment-textarea" />
        <el-button @click="submitComment" type="primary" class="comment-submit-btn">
          发表评论
        </el-button>
      </div>

      <!-- 评论列表 -->
      <div class="comment-list">
        <div v-for="comment in comments" :key="comment.Id" class="comment-item">
          <!-- 主评论 -->
          <div class="comment-main">
            <img :src="getImageUrl(comment.User.Avatar || 'boli.jpg')" class="avatar" />
            <div class="comment-content">
              <!-- 主评论内容 -->
              <div class="comment-header">
                <span class="username">{{ comment.User.Username }}</span>
                <span class="time">{{ formatTime(comment.CreatedAt) }}</span>
                <span v-if="comment.IsPinned === 1" class="pinned-tag">置顶</span>

                <!-- 主评论操作菜单 -->
                <div v-if="hasCommentActions(comment)" class="comment-actions">
                  <button class="menu-btn" @click.stop="toggleMenu(comment.Id)">⋮</button>
                  <div v-if="menuStates[comment.Id]" class="action-menu">
                    <button v-if="isAdmin" @click="pinComment(comment.Id, true)">
                      {{ comment.IsPinned === 1 ? "取消置顶" : "置顶" }}
                    </button>
                    <button v-if="isAdmin || comment.User.Id == currentUserId" @click="deleteComment(comment.Id)">
                      删除
                    </button>
                  </div>
                </div>
              </div>
              <div class="comment-text">{{ comment.Idea }}</div>
              <div class="comment-footer">
                <span class="reply-btn" @click="toggleReply(comment.Id, false)">回复</span>
              </div>

              <!-- 主评论回复框 -->
              <div v-if="activeReplies.parent === comment.Id" class="reply-input parent-reply-input">
                <el-input v-model="replyContents[comment.Id]" type="textarea" :rows="3" placeholder="写下你的回复..."
                  class="reply-textarea" />
                <div class="reply-buttons">
                  <el-button @click="submitReply(comment.Id)" type="primary" class="reply-submit-btn">提交回复</el-button>
                  <el-button @click="cancelReply(false)" type="text" class="reply-cancel-btn">取消</el-button>
                </div>
              </div>
            </div>
          </div>

          <!-- 回复列表 -->
          <div v-if="comment.Replies && comment.Replies.length > 0" class="reply-list">
            <div v-for="(reply, index) in sortedReplies(comment)" :key="reply.Id" class="reply-item">
              <img :src="getImageUrl(reply.User.Avatar || 'boli.jpg')" class="avatar" />
              <div class="comment-content">
                <div class="comment-header">
                  <span class="username">
                    {{ reply.User.Username }}
                    <span v-if="reply.ReplyId > 0" class="reply-target">
                      回复 @{{ getReplyUsername(reply.ReplyId) || '已删除用户' }}
                    </span>
                  </span>
                  <span class="time">{{ formatTime(reply.CreatedAt) }}</span>
                  <span v-if="reply.IsPinned === 1" class="pinned-tag">置顶</span>

                  <!-- 子评论操作菜单 -->
                  <div v-if="hasReplyActions(reply, comment)" class="comment-actions">
                    <button class="menu-btn" @click.stop="toggleMenu(reply.Id)">⋮</button>
                    <div v-if="menuStates[reply.Id]" class="action-menu">
                      <button v-if="isAdmin || comment.User.Id == currentUserId"
                        @click="pinComment(reply.Id, false, comment.Id)">
                        {{ reply.IsPinned === 1 ? "取消置顶" : "置顶" }}
                      </button>
                      <button v-if="isAdmin || reply.User.Id == currentUserId || comment.User.Id == currentUserId"
                        @click="deleteComment(reply.Id)">
                        删除
                      </button>
                    </div>
                  </div>
                </div>
                <div class="comment-text">{{ reply.Idea }}</div>
                <div class="comment-footer">
                  <span class="reply-btn" @click="toggleReply(reply.Id, true)">回复</span>
                </div>

                <!-- 子评论回复框（关键调整：嵌套在子评论content内） -->
                <div v-if="activeReplies.child === reply.Id" class="reply-input child-reply-input">
                  <el-input v-model="replyContents[reply.Id]" type="textarea" :rows="3" placeholder="写下你的回复..."
                    class="reply-textarea" />
                  <div class="reply-buttons">
                    <el-button @click="submitReply(reply.Id)" type="primary" class="reply-submit-btn">提交回复</el-button>
                    <el-button @click="cancelReply(true)" type="text" class="reply-cancel-btn">取消</el-button>
                  </div>
                </div>
              </div>
            </div>

            <!-- 显示更多/收起按钮 -->
            <div v-if="comment.Replies.length > 2" class="show-more-replies">
              <button @click="toggleShowAllReplies(comment)">
                {{ replyExpansionState[comment.Id] ? "收起回复" : `显示更多回复(${comment.Replies.length - 2}条)` }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, shallowRef, watch } from "vue";
import { ElInput, ElButton } from "element-plus";
import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;
const { proxy } = getCurrentInstance();
const getImageUrl = (imgName) => {
  return `${proxy.$imageBaseUrl}${imgName}`;
};

import { useArticleStore } from '@/stores/article';
const articleStore = useArticleStore();
const articleId =ref(articleStore.articleId);


const currentUserId = parseInt(localStorage.getItem("userId")) || 0;

const comments = ref([]);
const newComment = ref("");
const replyContents = ref({});

const userQx = localStorage.getItem("userQx");
const isAdmin = userQx === "A"; // 管理员权限

const usernameCache = shallowRef({});
const menuStates = ref({}); // 菜单状态
const activeReplies = ref({ parent: null, child: null }); // 回复框状态

// 权限检查
const hasCommentActions = computed(() => (comment) => {
  return isAdmin || comment.User.Id == currentUserId;
});
const hasReplyActions = computed(() => (reply, comment) => {
  return isAdmin || reply.User.Id == currentUserId || comment.User.Id == currentUserId;
});

// 重置评论相关状态（切换文章时调用）
const resetCommentState = () => {
  comments.value = [];
  newComment.value = "";
  replyContents.value = {};
  menuStates.value = {};
  activeReplies.value = { parent: null, child: null };
  usernameCache.value = {};
  replyExpansionState.value = {}; // 切换文章时重置展开状态
};

// 菜单控制
const toggleMenu = (commentId) => {
  menuStates.value = { ...menuStates.value, [commentId]: !menuStates.value[commentId] };
};
const closeMenu = () => { menuStates.value = {}; };

// 置顶功能
const pinComment = async (commentId, isParent, parentCommentId = null) => {
  try {
    let target;
    if (isParent) {
      target = comments.value.find((c) => c.Id === commentId);
    } else {
      const parentComment = comments.value.find((c) => c.Id === parentCommentId);
      target = parentComment?.Replies.find((r) => r.Id === commentId);
    }
    if (!target) return;

    const newStatus = target.IsPinned === 1 ? 0 : 1;
    const response = await fetch(`${URL}/comments/${commentId}/pin`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ is_pinned: newStatus }),
    });
    if (response.ok) {
      target.IsPinned = newStatus;
      if (isParent) {
        comments.value = [...comments.value].sort((a, b) => {
          if (a.IsPinned === 1 && b.IsPinned !== 1) return -1;
          if (a.IsPinned !== 1 && b.IsPinned === 1) return 1;
          return new Date(b.CreatedAt) - new Date(a.CreatedAt);
        });
      }
    }
      fetchComments();
  } catch (error) {
    console.error("置顶失败:", error);
  }
};

// 获取评论
const fetchComments = async () => {
  if (!articleId.value) return; // 文章ID为空时不请求
  try {
    const response = await fetch(`${URL}/comments/${articleId.value}`);
    if (!response.ok) throw new Error("获取评论失败");
    const data = await response.json();
    comments.value = (data.data || []).sort((a, b) => {
        // 置顶评论优先
      if (a.IsPinned === 1 && b.IsPinned !== 1) return -1;
      if (a.IsPinned !== 1 && b.IsPinned === 1) return 1;

        // 同为置顶评论，按 PinnedAt 降序（后置顶在前）
      if (a.IsPinned === 1 && b.IsPinned === 1) {
        return new Date(b.PinnedAt) - new Date(a.PinnedAt);
      }
        // 同为非置顶评论，按创建时间降序
      return new Date(b.CreatedAt) - new Date(a.CreatedAt);
    
      
    });  
  } catch (error) {
    console.error("获取评论失败:", error);
  }
};

// 回复用户名查找
const getReplyUsername = (replyId) => {
  if (usernameCache.value[replyId]) return usernameCache.value[replyId];
  const findUsername = (commentsList) => {
    for (const item of commentsList) {
      if (item.Id === replyId) {
        usernameCache.value[replyId] = item.User?.Username || null;
        return item.User?.Username;
      }
      if (item.Replies && item.Replies.length) {
        const username = findUsername(item.Replies);
        if (username) return username;
      }
    }
    return null;
  };
  return findUsername(comments.value);
};

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim()) {
    alert("评论内容不能为空哦 (＞﹏＜)");
    return;
  }
  try {
    const response = await fetch(`${URL}/path-to-article/${articleId.value}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userId: parseInt(currentUserId), idea: newComment.value }),
    });
    if (!response.ok) throw new Error((await response.json()).error || `HTTP error! status: ${response.status}`);
    alert("评论发布成功 (^_^) 快去看看吧 !");
    newComment.value = "";
    resetUsernameCache();
    fetchComments();
  } catch (err) {
    console.error("评论提交错误:", err);
    alert(`评论发布失败啦Σ(っ °Д °;)っ\n原因: ${err.message}`);
  }
};

// 提交回复
const submitReply = async (parentId) => {
  const content = replyContents.value[parentId];
  if (!content || !content.trim()) return;
  try {
    const response = await fetch(`${URL}/path-to-article/${articleId.value}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userId: parseInt(currentUserId), idea: content, parentId: parentId }),
    });
    if (!response.ok) throw new Error((await response.json()).error || `HTTP error! status: ${response.status}`);
    replyContents.value[parentId] = "";
    activeReplies.value.parent = null;
    activeReplies.value.child = null;
    resetUsernameCache();
    fetchComments();
  } catch (error) {
    console.error("提交回复失败:", error);
    alert(`回复失败啦Σ(っ °Д °;)っ\n原因: ${error.message}`);
  }
};

// 删除评论
const deleteComment = async (commentId) => {
  if (!confirm("确定要删除这条评论吗？")) return;
  try {
    const response = await fetch(`${URL}/path-to-article/${commentId}`, {
      method: "DELETE",
    });
    if (!response.ok) throw new Error((await response.json()).error || `HTTP error! status: ${response.status}`);
    fetchComments();
  } catch (error) {
    console.error("删除评论时出错:", error);
    alert(`删除失败啦Σ(っ °Д °;)っ\n原因: ${error.message}`);
  }
};

// 切换回复框显示
const toggleReply = (commentId, isChild = false) => {
  const key = isChild ? "child" : "parent";
  if (activeReplies.value[key] === commentId) {
    activeReplies.value[key] = null;
  } else {
    activeReplies.value[key] = commentId;
    replyContents.value[commentId] = "";
    activeReplies.value[isChild ? "parent" : "child"] = null;
  }
};

// 取消回复
const cancelReply = (isChild = false) => {
  activeReplies.value[isChild ? "child" : "parent"] = null;
};


const replyExpansionState = ref({}); 

// 切换显示所有回复
const toggleShowAllReplies = (comment) => {
    // comment.showAllReplies = !comment.showAllReplies;
  // 操作保存的状态
  replyExpansionState.value[comment.Id] = !replyExpansionState.value[comment.Id];
};

// 格式化时间
const formatTime = (timeStr) => {
  const date = new Date(timeStr);
  return date.toLocaleString();
};

// 子评论排序
const sortedReplies = (comment) => {
  if (!comment.Replies) return [];
  const sorted = [...comment.Replies].sort((a, b) => {
    // 置顶回复优先
    if (a.IsPinned === 1 && b.IsPinned !== 1) return -1;
    if (a.IsPinned !== 1 && b.IsPinned === 1) return 1;

     // 同为置顶回复，按 PinnedAt 降序（后置顶在前）
    if (a.IsPinned === 1 && b.IsPinned === 1) {
      return new Date(b.PinnedAt) - new Date(a.PinnedAt); // 关键修改：降序排列
    }

    // 同为非置顶回复，按创建时间降序
    return new Date(b.CreatedAt) - new Date(a.CreatedAt);
  });
 const isExpanded = replyExpansionState.value[comment.Id] || false;
  return isExpanded ? sorted : sorted.slice(0, 2);
};

// 重置用户名缓存
const resetUsernameCache = () => {
  usernameCache.value = {};
};

// 生命周期钩子
onMounted(() => {
  document.addEventListener("click", closeMenu);
  fetchComments(); // 初始加载
});

onUnmounted(() => {
  document.removeEventListener("click", closeMenu);
});

// 监听Pinia中的articleId变化（替代原来的localStorage监听）
watch(
  () => articleStore.articleId,
  (newId) => {
    if (newId) {
      articleId.value = newId; // 同步到本地ref
      resetCommentState();
      fetchComments(); // 触发评论刷新
    }
  },
  { immediate: true } // 初始加载时立即执行
);

</script>

<style scoped>
.comment-container {
  width: 70%;
  margin: 0 15%;
}

.comment-section {
  margin-top: 30px;
  padding: 20px;
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 8px;
}

.comment-input {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 20px;
}

.comment-textarea {
  flex: 1;
}

.comment-submit-btn {
  min-width: 100px;
  background-color: rgb(139, 189, 234) !important;
  border-color: rgb(139, 189, 234) !important;
  color: white !important;
}

.comment-submit-btn:hover {
  background-color: rgb(119, 172, 218) !important;
  border-color: rgb(119, 172, 218) !important;
}

.comment-list {
  margin-top: 20px;
}

.comment-item {
  margin-bottom: 20px;
}

.comment-main {
  display: flex;
  gap: 15px;
}

.reply-item {
  display: flex;
  gap: 15px;
  margin-left: 40px;
  margin-top: 15px;
  padding-left: 15px;
  border-left: 2px solid #eee;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.username {
  font-weight: bold;
}

.time {
  color: #999;
  font-size: 0.9em;
}

.comment-text {
  margin-bottom: 8px;
  line-height: 1.5;
}

.comment-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.reply-btn {
  color: rgb(20, 87, 111);
  font-size: 0.9em;
  cursor: pointer;
}

.reply-btn:hover {
  text-decoration: underline;
}

.delete-btn {
  padding: 4px 8px;
  background-color: #ff4d4f;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.8em;
  cursor: pointer;
}

.delete-btn:hover {
  background-color: #ff7875;
}

/* 回复框样式调整 */
.parent-reply-input {
  margin-top: 10px;
  margin-left: 55px;
}

.child-reply-input {
  margin-top: 10px;
  margin-left: 0; /* 子评论回复框不再强制缩进 */
}

.reply-input textarea {
  width: 100%;
  height: 60px;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: vertical;
}

.reply-buttons {
  margin-top: 8px;
  display: flex;
  gap: 8px;
}

.reply-buttons button {
  padding: 6px 12px;
  font-size: 0.9em;
}

.reply-submit-btn {
  background-color: rgb(139, 189, 234) !important;
  border-color: rgb(139, 189, 234) !important;
  color: white !important;
}

.reply-submit-btn:hover {
  background-color: rgb(119, 172, 218) !important;
  border-color: rgb(119, 172, 218) !important;
}

.reply-cancel-btn {
  color: #666 !important;
}

.reply-cancel-btn:hover {
  color: #333 !important;
}

.show-more-replies {
  margin: 10px 0 10px 55px;
}

.show-more-replies button {
  background: none;
  border: none;
  color: rgb(139, 189, 234);
  cursor: pointer;
  font-size: 0.9em;
  padding: 5px 0;
}

.show-more-replies button:hover {
  text-decoration: underline;
}

/* 置顶标签样式 */
.pinned-tag {
  display: inline-block;
  padding: 2px 6px;
  background-color: #f0f0f0;
  color: #666;
  border-radius: 4px;
  font-size: 0.8em;
  margin-left: 8px;
}

/* 三点菜单按钮 */
.menu-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #666;
  padding: 0 5px;
  font-size: 1.2em;
  line-height: 1;
}

.menu-btn:hover {
  color: rgb(139, 189, 234);
}

/* 操作菜单 */
.action-menu {
  position: absolute;
  right: 0;
  top: 100%;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 10;
  min-width: 100px;
}

.action-menu button {
  display: block;
  width: 100%;
  padding: 8px 12px;
  text-align: left;
  background: none;
  border: none;
  cursor: pointer;
}

.action-menu button:hover {
  background-color: #f5f5f5;
  color: rgb(139, 189, 234);
}

/* 评论操作容器 */
.comment-actions {
  position: relative;
  margin-left: auto;
}

.reply-target {
  color: #666;
  font-size: 0.9em;
  margin-left: 4px;
}

.reply-target:hover {
  color: rgb(139, 189, 234);
}
</style>