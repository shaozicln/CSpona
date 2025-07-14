<template>
  <div class="friends">
    <h2 class="h2222">Shaoの友链</h2>
    <div class="bookshelf-container">
      <div class="bookshelf">
        <div
          v-for="(friend, index) in friends"
          :key="friend.Id"
          class="book"
          :style="getBookPosition(index)"
          @click="openBook(friend)"
        >
          <div
            class="book-spine"
            :style="{
              backgroundColor: colors[index % colors.length],
            }"
          >
            <span class="book-title">{{ friend.Name || "敬请期待" }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 全屏背景层 (rgb(0,0,0.9) + 流星星河效果) -->
    <div class="galaxy-background" v-if="showModal">
      <!-- 流星星河效果 -->
      <div @click="closeModal" class="back-btn">
        <el-icon>
          <Back />
        </el-icon>
      </div>
      <div class="stars"></div>
      <div class="shooting-stars"></div>
    </div>

    <!-- 书籍内容层 (独立于背景层，确保交互性) -->
    <div
      class="book-container"
      v-if="showModal"
      :class="{ 'book-enter': showModal }"
    >
      <ul class="friend-details">
        <li>
          <img
            :src="getImageUrl(currentFriend.Img)"
            class="friendImg2"
            alt=""
          />
        </li>
        <li><img :src="getImageUrl('img.png')" alt="" class="friendImg2" /></li>
        <li><img :src="getImageUrl('img.png')" alt="" class="friendImg2" /></li>
        <li><img :src="getImageUrl('img.png')" alt="" class="friendImg2" /></li>
        <li>
          <div class="front-content">
            <img
              :src="getImageUrl(currentFriend.Avatar)"
              alt=""
              class="friendImg2"
            />
          </div>
          <div class="back-content">
            <img
              :src="getImageUrl(currentFriend.Avatar)"
              alt="作者头像"
              class="author-avatar"
            />
            <p class="author-message">"{{ currentFriend.Description }}"</p>
          </div>
        </li>
        <li>
          <img
            :src="getImageUrl(currentFriend.Background)"
            class="friendImg2"
            alt=""
          />
          <h2>{{ currentFriend.Name }}</h2>
          <p>{{ currentFriend.Introduction }}</p>
          <a
            :href="currentFriend.Web || 'https://cspona.top/'"
            target="_blank"
            class="aaa"
          >
            <button>点击跳转</button></a
          >
        </li>
      </ul>
    </div>

    <div id="send">
      <h2>有个的网站做得不错, 想放个友链 ?</h2>
      <router-link :to="{ path: '/feedback' }">
        <button @click="send()">点我联系作者发送友链</button>
      </router-link>
    </div>
  </div>
</template>

<style scope>
.friends {
  min-height: 100vh;
  padding: 50px;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}

h3 {
  text-align: center;
  font-size: 30px;
  font-family: cursive;
}

h2 {
  text-align: center;
  font-size: 50px;
  font-family: cursive;
}

.bookshelf-container {
  width: 100%;
  display: flex;
  justify-content: center;
  margin: 25px 0;
}

.bookshelf {
  display: flex;
  align-items: flex-end;
  height: 400px;
  perspective: 1000px;
  position: relative;
  padding: 20px 0;
}

.book {
  position: relative;
  height: 300px;
  width: 50px;
  margin: 0 -12px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); /* 更平滑的缓动函数 */
  transform-style: preserve-3d;
  cursor: pointer;
  z-index: 1;
  transform-origin: bottom; /* 从底部弹起 */
}

.book:hover {
  transform: translateY(-20px) scaleY(1.05); /* 增加轻微缩放效果 */
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3); /* 更强的阴影效果 */
  z-index: 10; /* 确保悬停的书在最上层 */
}

.book-spine {
  /* 保持原有样式 */
  transition: all 0.3s ease; /* 添加过渡效果 */
}

.book:hover .book-spine {
  box-shadow: 5px 15px 25px rgba(0, 0, 0, 0.4); /* 悬停时更强的阴影 */
}

.book-spine {
  position: absolute;
  width: 100%;
  height: 100%;
  /* 删除或注释掉这行 background 属性 */
  /* background: linear-gradient(90deg, #8a4b08, #a05a2c); */
  transform: rotateY(0deg) translateZ(25px);
  backface-visibility: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.3);
  border-radius: 2px 0 0 2px;
}

.book-title {
  color: white;
  font-weight: bold;
  writing-mode: vertical-rl;
  text-align: center;
  padding: 10px 0;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
  font-family: "Microsoft YaHei", sans-serif;
  font-size: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

/* 全屏背景层 */
.galaxy-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  z-index: 999;
  transition: opacity 0.5s ease;
}

/* 流星星河效果 */
.stars {
  background-image: url("https://picsum.photos/id/1015/1200/800");
  background-size: cover;
  background-position: center;
  /* 叠加发光效果 */
  box-shadow: inset 0 0 200px rgba(255, 215, 0, 0.3);
}

.shooting-stars {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

/* 流星效果 */
.shooting-stars::before {
  content: "";
  position: absolute;
  width: 4px;
  height: 4px;
  background: linear-gradient(90deg, transparent, white);
  border-radius: 50%;
  box-shadow: 0 0 10px 2px white;
  top: 10%;
  left: 100%;
  transform: rotate(45deg);
  animation: shooting-star 3s linear infinite;
}

.shooting-stars::after {
  content: "";
  position: absolute;
  width: 3px;
  height: 3px;
  background: linear-gradient(90deg, transparent, white);
  border-radius: 50%;
  box-shadow: 0 0 8px 1px white;
  top: 30%;
  left: 100%;
  transform: rotate(45deg);
  animation: shooting-star 4s linear 1s infinite;
}

@keyframes stars {
  from {
    background-position: 0 0;
  }

  to {
    background-position: -200px -200px;
  }
}

@keyframes shooting-star {
  0% {
    left: 100%;
    top: 10%;
    opacity: 0;
  }

  10% {
    opacity: 1;
  }

  80% {
    opacity: 1;
  }

  100% {
    left: -20%;
    top: 80%;
    opacity: 0;
  }
}

.back-icon-svg {
  width: 100px;
  height: 100px;
  color: white;
}

.back-btn {
  position: fixed;
  top: 10%;
  left: 15%;
  z-index: 1010;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  border-radius: 5px;
  padding: 8px 15px;
  font-size: 30px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background-color: rgba(0, 0, 0, 0.9);
  transform: scale(1.05);
}

.book-container {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) scale(0.9);
  transition: transform 0.5s ease 0.2s;
  z-index: 1005;
  width: 60%;
  display: flex;
  justify-content: center;
  pointer-events: auto;
}

.book-container.book-enter {
  transform: translate(-50%, -50%) scale(1);
}

/* 原有的ul和li样式 */
.friend-details {
  list-style: none;
  width: 291px;
  height: 420px;
  transform-style: preserve-3d;
  perspective: 900px;
  transition: 1.5s;
  position: relative;
}

.friend-details li {
  position: absolute;
  transform-origin: left;
  width: 291px;
  height: 420px;
}

.friend-details li:nth-child(1) {
  transform: rotateY(-25deg);
  transition: 2s;
}

.friend-details li:nth-child(2) {
  transform: rotateY(-23deg);
  transition: 1.7s;
}

.friend-details li:nth-child(3) {
  transform: rotateY(-21deg);
  transition: 1.4s;
}

.friend-details li:nth-child(4) {
  transform: rotateY(-19deg);
  transition: 1.1s;
}

.friend-details li:nth-child(5) {
  transform: rotateY(-17deg);
  transition: 0.8s;
}

.friend-details li:nth-child(6) {
  transform: rotateY(-15deg);
  transition: 0.5s;
}

.friend-details:hover li:nth-child(1) {
  transform: rotateY(-180deg);
  transition: 1.4s;
}

.friend-details:hover li:nth-child(2) {
  transform: rotateY(-180deg);
  opacity: 0.6;
  transition: 2s;
}

.friend-details:hover li:nth-child(3) {
  transform: rotateY(-180deg);
  opacity: 0.6;
  transition: 2.6s;
}

.friend-details:hover li:nth-child(4) {
  transform: rotateY(-180deg);
  opacity: 0.6;
  transition: 3.2s;
}

.friend-details:hover li:nth-child(5) {
  transform: rotateY(-180deg);
  transition: 3.8s;
}

.friend-details:hover {
  transform: translateX(150px);
}

.friendImg2 {
  width: 100%;
  height: 100%;
  object-fit: cover;
  box-shadow: 1px 4px 5px rgba(0, 0, 0, 0.2);
}

.friendImg {
  width: 120px;
  height: 120px;
  object-fit: cover;
  box-shadow: 1px 4px 5px rgba(0, 0, 0, 0.2);
}

/* 第五张纸的特殊样式 - 正面 */
.friend-details li:nth-child(5) .front-content {
  position: absolute;
  width: 100%;
  height: 100%;
}

/* 第五张纸的特殊样式 - 背面 */
.friend-details li:nth-child(5) .back-content {
  position: absolute;
  width: 100%;
  height: 100%;
  transform: rotateY(180deg);
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: white;
}

/* 圆形作者图片 */
.author-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid rgb(112, 184, 227);
  margin-bottom: 20px;
}

/* 作者留言 */
.author-message {
  font-style: italic;
  text-align: center;
  max-width: 80%;
  line-height: 1.5;
}

/* 最后一个li的特殊样式 */
.friend-details li:last-child {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  color: white;
  background: rgba(0, 0, 0, 0.9);
}

.friend-details li:last-child img {
  position: absolute;
  top: 0;
  left: 0;
  opacity: 0.5;
  z-index: -1;
}

.friend-details li:last-child h2 {
  font-size: 24px;
  margin-bottom: 10px;
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.8);
}

.friend-details li:last-child p {
  font-size: 16px;
  margin-bottom: 20px;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.8);
}

.friend-details li:last-child button {
  padding: 10px 20px;
  background: rgb(112, 184, 227);
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.friend-details li:last-child button:hover {
  background: rgb(112, 184, 227);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
}

#send {
  margin-top: 80px;
  gap: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

#send button {
  background-color: rgb(139, 189, 234);
  color: #fff;
  border: none;
  padding: 10px 30px;
  font-size: 20px;
  cursor: pointer;
  transition: all 0.3s ease-in-out;
  border-radius: 20px;
  width: 100%;
}

#send button:hover {
  background-color: rgba(139, 190, 234, 0.8);
  transform: translateY(-5px);
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}
@media (max-width: 768px) {
  .bookshelf {
    height: 250px;
  }
  .book {
    height: 200px;
    width: 40px;
  }
  .book-title {
    font-size: 14px;
  }
}
</style>

<script setup>
import { Back } from "@element-plus/icons-vue";
// 获取全局URL属性
import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

import { ref, onMounted } from "vue";

const { proxy } = getCurrentInstance();
const getImageUrl = (imgName) => {
  return `${proxy.$imageBaseUrl}${imgName}`;
};

const friends = ref([]);
const sendVisible = ref(false);
const sendQuery = ref("");
const sendInput = ref(null);
const showModal = ref(false);
const currentFriend = ref({});

const colors = [
  "rgb(178,188,198)",
  "rgb(190,188,148)",
  "rgb(129,174,175)",
  "rgb(197,219,214)",
  "rgb(162,204,199)",
  "rgb(174,197,181)",
  "rgb(236,196,150)",
  "rgb(227,152,123)",
  "rgb(228,156,141)",
  "rgb(214,111,112)",
];

onMounted(async () => {
  try {
    const response = await fetch(`${URL}/friendsWeb`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    const data = await response.json();
    const realFriends = data.data || [];
    
    // 固定显示10本书
    const totalBooks = 10;
    
    // 创建最终数组，先全部填充默认值
    const allFriends = Array.from({ length: totalBooks }, (_, i) => ({
      Id: `default-${i}`,
      Name: "敬请期待",
      Web: "#",
      Img: "",
      Introduction: "",
    }));
    
    // 明确指定要优先占用的位置（第5和第6本，索引4和5）
    const priorityPositions = [4, 5]; // 第5和第6个位置
    
    // 先填充优先位置
    let realIndex = 0;
    for (const pos of priorityPositions) {
      if (realIndex < realFriends.length && pos < allFriends.length) {
        allFriends[pos] = realFriends[realIndex++];
      }
    }
    
    // 剩余友链按顺序填充其他位置
    for (let i = 0; i < allFriends.length && realIndex < realFriends.length; i++) {
      if (!priorityPositions.includes(i)) {
        allFriends[i] = realFriends[realIndex++];
      }
    }
    
    friends.value = allFriends;
    
  } catch (error) {
    console.error("获取友链数据失败：", error);
    friends.value = Array.from({ length: 10 }, (_, i) => ({
      Id: `error-default-${i}`,
      Name: "敬请期待",
      Web: "#",
      Img: "",
      Introduction: "",
    }));
  }
});

// 位置计算函数（保持从中间展开效果）
const getBookPosition = (index) => {
  const center = 4.5; // 书架中心在第4和第5本之间
  const position = index - center;
  
  return {
    transform: `translateX(${position * 35}px)`,
    zIndex: 10 - Math.abs(position), // 中心层级最高
  };
};

const openBook = (friend) => {
  if (!friend.Name || friend.Name.trim() === "敬请期待") {
    alert("敬请期待~");
    return;
  }
  currentFriend.value = friend;
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
};

const send = () => {
  sendVisible.value = true;
  setTimeout(() => {
    sendInput.value.focus();
  }, 0);
};

const performSend = async () => {
  // perform send logic here
};
</script>
