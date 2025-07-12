<template>
  <div class="home-container">
    <div id="TypingTextContainer">
      <div id="TypingText">
        <span>{{ displayText }}</span>
        <span class="blink">|</span>
      </div>
      <div class="welcome-container">
        <p>----------Welcome----------</p>
      </div>
    </div>
    <div class="myself">
      <div class="me">
        <Author />
      </div>
      <div class="markdown-content">
        <h1>关于我</h1>
        <p>
          你好，我是一个热爱编程的开发者。专注于Web开发，熟悉 Vue.js 和 Go
          语言。喜欢探索新技术并将其应用于实际项目中。
        </p>
        <div class="profile-info">
          <div class="profile-item">
            <i class="fas fa-microscope"></i>
            <span>🔭 I'm 长柄木勺, a sophomore student</span>
          </div>
          <div class="profile-item">
            <i class="fas fa-seedling"></i>
            <span
              >🌱 I am learning computer technology related content and basic
              algorithmic knowledge</span
            >
          </div>
          <div class="profile-item">
            <i class="fas fa-smile"></i>
            <span>😄 Tech stack: Javascript Vue Go Gin Mysql</span>
          </div>
          <div class="profile-item">
            <i class="fas fa-comment"></i>
            <span
              >💬 I've done: two front-end blogs, some web mini-games and fun
              little features</span
            >
          </div>
          <div class="profile-item">
            <i class="fas fa-lightbulb"></i>
            <span
              >🤔 Current mini-goal: algorithmic fundamentals and
              general-purpose technical implementations</span
            >
          </div>
          <div class="profile-item">
            <i class="fas fa-bolt"></i>
            <span
              >⚡ Hobbies: gaming, watching anime and travelling, love to see
              different stories and landscapes</span
            >
          </div>
          <div class="profile-item">
            <i class="fas fa-envelope"></i>
            <span>📫 How to reach me: changbingmushao@qq.com </span>
          </div>
        </div>
      </div>
    </div>
    <div class="contact-info"style="display: flex; align-items: center; gap: 4px;flex-direction:column;justify-content:center" >
      <span>----------欢迎来到CSpona!----------</span>
      <span style="display: flex; align-items: center; gap: 4px"
        >备案号：<a
          href="https://beian.miit.gov.cn/"
          target="_blank"
          rel="noopener noreferrer"
          style="color:rgb( 79, 79, 79);"
        >
          黑ICP备2025038200号-1
        </a>
      </span>
    </div>
  </div>
</template>

<script setup>
// 获取全局URL属性
import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;
import Author from "@/components/Author/Author.vue";

import { ref, onMounted, onUnmounted, computed } from "vue";
import { library } from "@fortawesome/fontawesome-svg-core";
import { faGithub } from "@fortawesome/free-brands-svg-icons";
import {
  faMicroscope,
  faSeedling,
  faSmile,
  faComment,
  faLightbulb,
  faBolt,
  faEnvelope,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

library.add(faGithub, faEnvelope);

// 注册 FontAwesomeIcon 组件
import { defineComponent } from "vue";
defineComponent({
  components: {
    FontAwesomeIcon,
  },
});

const fullTexts = ref([
  "你好, 新朋友 ..",
  "上决浮云, 下绝地纪; 截云断岳, 剑出山倾 ..",
  "你想活出怎样的人生 ?",
  "感受Coding带来的创造的乐趣吧 !",
]);

// 显示
const displayText = ref("");
// 索引
const currentTextIndex = ref(0);
// 进度
const currentProgress = ref(0);
// 速度-----每秒显示的字符数
const speed = 4;
// 定时器ID
let intervalId;

// 计算当前应该显示的文本
const currentText = computed(() => {
  return fullTexts.value[currentTextIndex.value].substring(
    0,
    currentProgress.value
  );
});

onMounted(() => {
  // 启动逐字显示
  intervalId = setInterval(() => {
    if (
      currentProgress.value < fullTexts.value[currentTextIndex.value].length
    ) {
      currentProgress.value += 1; // 增加显示的字符数
    } else {
      // 当前语句显示完毕，准备显示下一条语句
      currentProgress.value = 0;
      currentTextIndex.value =
        (currentTextIndex.value + 1) % fullTexts.value.length;
    }
    displayText.value = currentText.value;
  }, 1000 / speed); // 根据速度设置定时器的时间间隔
});

onUnmounted(() => {
  clearInterval(intervalId); // 清除定时器
});
</script>

<style scoped>
.myself {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.home-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.welcome-container {
  position: absolute;
  top: 62%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 30px;
  font-family: Cormorant SC, serif;
  text-align: center;
  /* 确保文本居中 */
}

#TypingTextContainer {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  height: 60vh;
  width: 100vw;
  font-size: 30px;
  font-family: cursive;
  /* 保持字体一致 */
  position: relative;
  /* 确保子元素的绝对定位相对于此容器 */
}

#TypingText {
  position: relative;
  /* 为动画元素定位 */
}

.blink {
  animation: blink 0.75s infinite;
  /* 调整闪烁速度 */
}

@keyframes blink {
  0% {
    opacity: 1;
  }

  50% {
    opacity: 0;
  }

  100% {
    opacity: 1;
  }
}

.me {
  margin: 25px;
}

/* 美化 markdownContent 的样式 */
.markdown-content {
  margin: 25px;
  font-size: 18px;
  font-family: "楷体";
  line-height: 1.6;
  color: #333;
  padding: 20px 20px 40px 20px;
  background-color: #fff;
  border-radius: 8px;
  /*box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);*/
  max-width: 800px;
}

.markdown-content h1,
.markdown-content h2,
.markdown-content h3,
.markdown-content h4,
.markdown-content h5,
.markdown-content h6 {
  margin-top: 20px;
  margin-bottom: 10px;
  /*color: #333;*/
}

.markdown-content p {
  margin-bottom: 15px;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
}

.profile-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  /*color: #555;*/
}

.profile-item i {
  /*color: #007BFF;*/
  font-size: 20px;
}

.github-stats {
  margin-top: 20px;
  text-align: center;
}

.github-stats img {
  border-radius: 8px;
  /*box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);*/
}

.contact-info {
  width: 100vw;
  /* display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center; */
  margin-top: 100px;
  margin-bottom: 0px;
  font-size: 20px;
  color: #494949;
  /* 灰色字 */
  text-align: center;
  background-color: rgba(255, 255, 255, 0.7);
  /* 灰色背景 */
  padding: 40px;
  /* border-radius: 16px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); */
}

.contact-info span {
  margin: 5px 0;
}

/* .contact-info a {
    color: #00000053;
    text-decoration: none;
}

.contact-info a:hover {
    text-decoration: underline;
}

.contact-info i {
    margin-right: 5px;
    color: #ad0e0e;
} */
</style>
