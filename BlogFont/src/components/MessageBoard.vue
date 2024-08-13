<template>
  <div class="message-board">
    <div class="send-box">
      <input type="text" v-model="message" placeholder="你好呀,放个留言 ?" />
      <button @click="sendMessage"><i class="iconfont icon-send-line"></i></button>
    </div>
    <div class="message-list">
      <ul>
        <li v-for="(msg, index) in displayMessages" :key="index" :style="{ top: `${msg.top}px`, left: `${msg.left}px` }"
          :class="{ 'sending': sendingIndex === index }">
          {{ msg.text }}
        </li>
      </ul>
    </div>
  </div>
  <Friends />
</template>

<script setup>
import Friends from './Friends.vue'
import { ref, onMounted, onUnmounted } from 'vue'

const message = ref('')
const messages = ref([])
const displayMessages = ref([])
const sendingIndex = ref(-1)
const refreshInterval = ref(null)


//随机位置
const randomPosition = (type) => {
  if (type === 'top') {
    return randomTop()
  } else if (type === 'left') {
    return randomLeft()
  }
}
const randomTop = () => {
  const inputHeight = 100
  const inputTop = window.innerHeight / 2 - inputHeight / 2
  let randomTop = Math.random() * (window.innerHeight - inputHeight)
  while (randomTop > inputTop && randomTop < inputTop + inputHeight) {
    randomTop = Math.random() * (window.innerHeight - inputHeight)
  }
  // 检查重叠
  for (let i = 0; i < displayMessages.value.length; i++) {
    const existingMessage = displayMessages.value[i]
    const distance = Math.sqrt(Math.pow(randomTop - existingMessage.top, 2) + Math.pow(randomLeft - existingMessage.left, 2))
    if (distance < 100) {
      randomTop = Math.random() * (window.innerHeight - inputHeight)
      break
    }
  }
  return randomTop
}
const randomLeft = () => {
  const inputWidth = 300
  const inputLeft = window.innerWidth / 2 - inputWidth / 2
  let randomLeft = Math.random() * (window.innerWidth - inputWidth)
  while (randomLeft > inputLeft && randomLeft < inputLeft + inputWidth) {
    randomLeft = Math.random() * (window.innerWidth - inputWidth)
  }
  // 检查重叠
  for (let i = 0; i < displayMessages.value.length; i++) {
    const existingMessage = displayMessages.value[i]
    const distance = Math.sqrt(Math.pow(randomTop - existingMessage.top, 2) + Math.pow(randomLeft - existingMessage.left, 2))
    if (distance < 100) {
      randomLeft = Math.random() * (window.innerWidth - inputWidth)
      break
    }
  }
  return randomLeft
}
// 获取留言列表
const getMessages = async () => {
  try {
    const response = await fetch('http://127.0.0.1:8081/commentBoard', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    const data = await response.json()
    messages.value = data.data.map(item => item.Content)

    // 随机选取 15 条留言
    const randomMessages = messages.value.sort(() => Math.random() - 0.5).slice(0, 13)

    displayMessages.value = randomMessages.map((msg, index) => {
      return {
        text: msg,
        top: randomPosition('top'),
        left: randomPosition('left')
      }
    })
  } catch (error) {
    console.error(error)
  }
}

// 发送留言
const sendMessage = async () => {
  console.log(message.value.trim());
  console.log(typeof message.value.trim());
  const userid = localStorage.getItem("userId");
  if (userid !== null && userid !== "") {
    const parsedUserId = parseInt(userid, 10);
    if (isNaN(parsedUserId)) {
      console.error("UserId is not a valid integer");
    } else {
      // use parsedUserId
    }
  } else {
    console.error("UserId is not stored in localStorage");
  }

  if (!message.value.trim()) return; // 忽略空格或空留言
  try {
    const response = await fetch('http://127.0.0.1:8081/commentBoard', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        UserId: parseInt(userid),
        Content: message.value.trim(),
      }),
    });
    console.log(typeof (parseInt(userid)));
    console.log(parseInt(userid));
    const data = await response.json()
    // 添加到 messages 数组
    messages.value.push(data.data.Content)
    // 添加到 displayMessages 数组
    displayMessages.value.push({
      text: message.value,
      top: randomPosition('top'),
      left: randomPosition('left')
    })
    // 设置发送状态
    sendingIndex.value = displayMessages.value.length - 1;
    // 清空输入框
    message.value = '';
    // 1秒后取消发送状态
    setTimeout(() => {
      sendingIndex.value = -1;
    }, 1000);
  } catch (error) {
    console.error('发送留言失败:', error)
  }
}

//组件挂载卸载对于定时器的设置和清除
onMounted(() => {
  getMessages()
  refreshInterval.value = setInterval(() => {
    getMessages()
  }, 5000)
})
onUnmounted(() => {
  clearInterval(refreshInterval.value)
})
</script>

<style scoped>
.message-board {
  position: relative;
  height: 100vh;
  overflow: hidden;
}

.send-box {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 5px;
  background-color: #f0f0f02d;
  border: 1px solid #ccc;
  border-radius: 20px;
  z-index: 1;
  width: 300px;
  height: 40px;
  display: flex;
  flex-direction: row;
}

.send-box input {
  width: 80%;
  height: 20px;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
}

.send-box button {
  width: 20%;
  height: 30px;
  background-color: #38cde400;
  color: #fff;
  border: none;
  border-radius: 20px;
  cursor: pointer;
}

.icon-send-line {
  font-size: 30px !important;
  color: black !important;
}

.message-list {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.message-list ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.message-list li {
  position: absolute;
  padding: 10px;
  font-size: 20px;
  background-color: #ffffff5f;
  border: 1px solid #ccc;
  border-radius: 5px;
  transition: all 1s;
  font-family: cursive;
}

.message-list li.sending {
  animation: fadeIn 1s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.5);
  }

  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>