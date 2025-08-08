<template>
    <div class="page-container">
        <!-- MyHead组件 - 仅在满足条件时显示 -->
        <!-- <transition name="slide"> -->
            <div class="nav-wrapper">
                <MyHead />
            </div>
        <!-- </transition> -->
        <!-- 轮播主容器 -->
        <section class="con" @click="handleContainerClick">
            <div class="slider">
                <!-- 轮播项通过v-for渲染 -->
                <div class="slide" :class="{ active: currentIndex === index }" v-for="(slide, index) in slides"
                    :key="index">
                    <img :src="getImageUrl(slide.imgUrl)" :alt="slide.title">
                    <div class="left">
                        <div id="layer">
                            <h1 class="layerTitle">{{ slide.layerTitle }}</h1>
                        </div>
                        <div id="left-text">
                            <h3>{{ slide.subtitle }}</h3>
                            <p v-html="slide.description"></p>
                            <br>
                            <a>
                                <div id="left-btn">
                                    <h4>enter</h4>
                                </div>
                            </a>
                        </div>
                    </div>
                    <div class="right">
                        <h1 class="layerTitle">{{ slide.rightTitle }}</h1>
                        <h3>{{ slide.rightSubtitle }}</h3>
                    </div>
                </div>
            </div>

            <!-- 左右切换按钮 -->
            <div class="navi">
                <span class="prev" @click.stop="handlePrev">
                    <ArrowLeft class="arrow-icon" />
                </span>
                <span class="next" @click.stop="handleNext">
                    <ArrowRight class="arrow-icon" />
                </span>
            </div>
        </section>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { Back } from "@element-plus/icons-vue";
import { ArrowLeft, ArrowRight } from "@element-plus/icons-vue";
import MyHead from '@/components/MyHead.vue';

import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const { proxy } = instance || {};
const URL = instance?.appContext.config.globalProperties.URL;

const getImageUrl = (imgName) => {
    if (!imgName) return `${proxy?.$imageBaseUrl}boli.jpg`;
    return `${proxy?.$imageBaseUrl}${imgName}`;
};

// 路由实例
const router = useRouter();

// 控制导航栏显示的变量
const isNavVisible = ref(false);
// 导航栏是否被鼠标悬停
const isNavHovered = ref(false);

// 轮播数据
const slides = [
    {
        imgUrl: 'Venti-7.jpg',
        title: '透龙山风景区',
        layerTitle: 'TaskRe ',
        subtitle: 'Tools - 01 基于邮箱SMTP的任务提醒',
        description: '一个简单的网页，让用户可以添加自定义的提醒任务。<br>比如:“明天下午3点取快递”、“这周五下午5点前交实验报告”、“下周三四六级报名截止”等。<br>通过邮箱的SMTP功能通知用户，',
        rightTitle: ' minder',
        rightSubtitle: 'Reminder Tool',
        route: '/task-reminder'
    },
    {
        imgUrl: 'venti-2.jpg',
        title: '日月峡国家森林公园',
        layerTitle: 'RiYue ',
        subtitle: '02日月峡国家森林公园',
        description: '日月峡国家森林公园位于小兴安岭南麓，呼兰河上游，黑龙江省铁力林业局马永顺（林业英雄）林场内；行政区域横跨伊春市铁力市和绥化市庆安县，占地面积29708公顷。',
        rightTitle: ' Xia',
        rightSubtitle: 'Forest Scenery',
        route: '/celebration'
    },
    {
        imgUrl: '../图片库/桃山悬羊峰.webp',
        title: '桃山悬羊峰',
        layerTitle: 'XuanYa ',
        subtitle: '03桃山悬羊峰',
        description: '悬羊峰景区在桃山林业局东南37千米处，白河林场施业区内。悬羊峰始建于1984年，2009年被国土资源部批准为国家级地质公园。2011年被国家旅游局批准为AAAA级景区。',
        rightTitle: ' ngFeng',
        rightSubtitle: 'Cliff Scape',
        route: '/tourism/xuanyang'
    },
    {
        imgUrl: '../图片库/日月峡滑雪场.webp',
        title: '日月峡滑雪场',
        layerTitle: 'Hua ',
        subtitle: '04日月峡滑雪场',
        description: '日月峡滑雪场按照3S级标准建设，现已晋升为4S级，拥有三条高、中、初级和一条空中技巧雪道，总长度5000多米，还有滑雪圈、雪地摩托及雪地自行车场地。',
        rightTitle: ' Xue',
        rightSubtitle: 'Ski Resort',
        route: '/tourism/skiresort'
    }
];

// 当前轮播索引
const currentIndex = ref(0);

// 监听鼠标移动事件
const handleMouseMove = (e) => {
    // 当鼠标在屏幕上方100px范围内且没有悬停在导航上时
    if (e.clientY < 100 && !isNavHovered.value) {
        isNavVisible.value = true;
    }
    // 当鼠标离开上方区域且没有悬停在导航上时隐藏
    else if (e.clientY >= 100 && !isNavHovered.value) {
        isNavVisible.value = false;
    }
};

// 鼠标进入导航区域
const handleNavMouseEnter = () => {
    isNavHovered.value = true;
    isNavVisible.value = true; // 保持显示
};

// 鼠标离开导航区域
const handleNavMouseLeave = () => {
    isNavHovered.value = false;
    // 检查鼠标是否还在上方区域
    isNavVisible.value = window.event?.clientY < 100;
};

// 切换到下一张
const handleNext = () => {
    currentIndex.value = (currentIndex.value + 1) % slides.length;
};

// 切换到上一张
const handlePrev = () => {
    currentIndex.value = (currentIndex.value - 1 + slides.length) % slides.length;
};

// 点击容器（非按钮区域）跳转路由
const handleContainerClick = () => {
    router.push(slides[currentIndex.value].route);
};

// 挂载时添加鼠标监听
onMounted(() => {
    window.addEventListener('mousemove', handleMouseMove);
});

// 卸载时移除监听
onUnmounted(() => {
    window.removeEventListener('mousemove', handleMouseMove);
});
</script>

<style scoped>
/* 保留原CSS样式并添加导航相关样式 */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

/* 导航容器样式 */
.nav-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 999;
    /* 确保悬浮在当前组件上方 */
    transition: opacity 0.3s ease;
}

/* 轮播相关样式 */
.con {
    position: relative;
    width: 100%;
    height: 100vh;
    overflow: hidden;
}

.slider .slide {
    position: absolute;
    width: 100%;
    height: 100%
}

.slide img {
    position: absolute;
    width: 100%;
    height: 100%;
    object-fit: cover;
    pointer-events: none;
    opacity: 0;
    transition: .5s ease;
}

.slide.active img {
    opacity: 1;
}

.slide .left {
    position: absolute;
    top: 0;
    left: 0;
    width: 50%;
    height: 100%;
    transform: translateX(-100%);
    transition: 0;
}

.slide.active .left {
    transform: translateX(0);
    z-index: 1;
    transition: 1s ease;
}

.left #layer {
    position: absolute;
    width: 100%;
    height: 100vh;
    background: rgba(255, 255, 255, .1);
    display: flex;
    justify-content: flex-end;
    align-items: center;
    backdrop-filter: blur(20px);
    -webkit-mask: linear-gradient(#000 0 0), linear-gradient(#000 0 0);
    -webkit-mask-clip: text, padding-box;
    -webkit-mask-composite: xor;
}

#layer h1 {
    font-size: 200px;
    text-shadow: 0 0 10px rgba(255, 255, 255, 0.8);
}

.left #left-text {
    position: absolute;
    bottom: 10%;
    left: 10%;
    color: #fff;
}

#left-text h3 {
    font-size: 30px;
}

#left-text p {
    font-size: 20px;
    margin: 10px 10px 15px 0px;
}

#left-btn {
    display: inline-block;
    padding: 13px 28px;
    background: #fff;
    border: 2px solid #fff;
    border-radius: 6px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    text-decoration: none;
    font-weight: 600;
    color: #444;
    font-size: 16px;
    transition: all .3s linear;
    cursor: pointer;
}

#left-text #left-btn:hover {
    background: transparent;
    color: #fff;
}

.slide .right {
    position: absolute;
    top: 0;
    right: 0;
    width: 50%;
    height: 100%;
    display: flex;
    align-items: center;
    transform: translateX(100%);
    transition: 0s;
}

.slide.active .right {
    transform: translateX(0);
    transition: 1s ease;
}

.right h1 {
    font-size: 200px;
    color: #fff;
    text-shadow:
        0 1px 0 #ccc,
        0 2px 0 #c9c9c9,
        0 3px 0 #bbb,
        0 4px 0 #b9b9b9,
        0 5px 0 #aaa,
        0 6px 1px rgba(0, 0, 0, .1),
        0 0px 5px rgba(0, 0, 0, .1),
        0 1px 3px rgba(0, 0, 0, .3),
        0 3px 5px rgba(0, 0, 0, .2),
        0 5px 10px rgba(0, 0, 0, .25),
        0 10px 10px rgba(0, 0, 0, .2),
        0 20px 20px rgba(0, 0, 0, .15);
}

.right h3 {
    position: absolute;
    font-size: 80px;
    color: #fff;
    text-shadow: 0 0 10px rgba(0, 0, 0, .5);
    transform: translateY(200%);
    margin-left: 10px;
}

.navi {
    position: absolute;
    bottom: 8%;
    right: 5%;
    z-index: 100;
}

.navi span {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    width: 50px;
    height: 50px;
    background: #fff;
    border: 2px solid #fff;
    border-radius: 6px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
    cursor: pointer;
    margin-left: 25px;
    transition: .3s;
}

.navi span:nth-child(1) {
    background: transparent;
}

.navi span:nth-child(1):hover {
    background: #fff;
}

.arrow-icon {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.6);
}

/* 过渡动画核心样式 */
.slide-enter-active,
.slide-leave-active {
    transition: transform 0.5s ease;
    /* 0.5秒平滑过渡 */
}

/* 进入开始状态 - 上方隐藏 */
.slide-enter-from {
    transform: translateY(-100%);
}

/* 离开结束状态 - 上方隐藏 */
.slide-leave-to {
    transform: translateY(-100%);
}

/* 导航容器基础样式 */
.nav-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 999;
    transform: translateY(0);
    /* 默认位置 */
}
</style>