<template>
    <div id="TypingTextContainer">
        <div id="TypingText">
            <span>{{ displayText }}</span>
            <span class="blink">|</span>
        </div>
        <div class="welcome-container">
            <h3>----------Welcome----------</h3>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';


const fullTexts = ref([
    "你好, 新朋友 ..",
    "上决浮云, 下绝地纪; 截云断岳, 剑出山倾 ..",
    "你想活出怎样的人生 ?",
    "感受Coding带来的创造的乐趣吧 !",
]);

//显示
const displayText = ref("");
//索引
const currentTextIndex = ref(0);
//进度
const currentProgress = ref(0);
//速度-----每秒显示的字符数
const speed = 4;
// 定时器ID
let intervalId;

// 计算当前应该显示的文本
const currentText = computed(() => {
    return fullTexts.value[currentTextIndex.value].substring(0, currentProgress.value);
});

onMounted(() => {
    // 启动逐字显示
    intervalId = setInterval(() => {
        if (currentProgress.value < fullTexts.value[currentTextIndex.value].length) {
            currentProgress.value += 1; // 增加显示的字符数
        } else {
            // 当前语句显示完毕，准备显示下一条语句
            currentProgress.value = 0;
            currentTextIndex.value = (currentTextIndex.value + 1) % fullTexts.value.length;
        }
        displayText.value = currentText.value;
    }, 1000 / speed); // 根据速度设置定时器的时间间隔
});

onUnmounted(() => {
    clearInterval(intervalId); // 清除定时器
});
</script>

<style scoped>
.welcome-container {
    position: absolute;
    top: 72%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 20px;
    font-family: "楷体";
}


#TypingTextContainer {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    font-size: 40px;
}

#TypingText {
    position: relative;
    /* 为动画元素定位 */
    top: auto;
    /* 移除top属性，因为我们使用Flexbox居中 */
    animation: none;
    /* 确保没有其他动画影响 */
}

.blink {
    animation: blink 0.25s infinite;
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
</style>