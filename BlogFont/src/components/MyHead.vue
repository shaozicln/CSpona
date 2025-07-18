<template>
  <div class="nav-bar">
    <!-- 外层容器：用于整体鼠标离开检测 -->
    <div class="nav-container" @mouseleave="closeAllDropdowns">
      <ul>
        <li 
          v-for="(item, index) in navItems" 
          :key="index" 
          :class="{ active: activeIndex === index }"
          class="nav-item"
        >
          <!-- 主链接容器 -->
          <div class="main-item" @click.stop>
            <router-link :to="item.to" class="main-link">
              {{ item.text }}
            </router-link>
            
            <!-- 箭头图标 -->
            <template v-if="item.children && item.children.length">
              <span 
                class="arrow" 
                   @mouseenter.stop="toggleDropdown(index)"
                :class="{ 'arrow-open': item.isOpen }"
              >
                ▼
              </span>
            </template>
          </div>
          
          <!-- 下拉菜单：增加hover范围 -->
          <div class="dropdown-wrapper" v-if="item.isOpen">
            <ul class="dropdown-menu">
              <li 
                v-for="(child, i) in item.children" 
                :key="i"
                class="dropdown-item"
              >
                <router-link 
                  :to="child.to" 
                  @click.stop="handleChildClick(index)"
                   @mouseenter.stop="keepDropdownOpen(index)"
                >
                  {{ child.text }}
                </router-link>
              </li>
            </ul>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { getCurrentInstance, ref } from 'vue';
import { useRoute } from 'vue-router';

// 全局配置获取
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;
const route = useRoute();

// 激活状态索引
const activeIndex = ref(0);

// 用户类型判断
const username = localStorage.getItem("username");
const email = localStorage.getItem("email");
let userType = "";
if (username === "长柄木勺" && email === "changbingmushao@qq.com") {
  userType = "管理";
} else {
  userType = "反馈";
}

// 导航项配置
const navItems = ref([
  { 
    // to: '/log', 
    to: '/', 
    text: 'CSpona:)', 
    children: [
      { text: '个人中心', to: '/user' },
      { text: '登录/注册', to: '/login' }
    ],
    isOpen: false 
  },
  { to: '/', text: '首页', children: [], isOpen: false },
  { 
    to: '/message-board', 
    text: '留言/友链', 
    children: [
      { text: '漫游地', to: '/wanderland/95' }
    ],
    isOpen: false 
  },
  { to: '/articles', text: '文章', children: [], isOpen: false },
  { to: '/sort', text: '数据', children: [], isOpen: false },
  { 
    to: userType === "管理" ? '/back' : '/Feedback', 
    text: userType, 
    children: [], 
    isOpen: false 
  },
  { to: '/search', text: '搜索一下', children: [], isOpen: false }
]);

// 切换下拉菜单展开/收起
const toggleDropdown = (index) => {
  navItems.value[index].isOpen = !navItems.value[index].isOpen;
  // 收起其他下拉菜单
  navItems.value.forEach((item, i) => {
    if (i !== index) item.isOpen = false;
  });
  activeIndex.value = index;
};

// 子项点击处理
const handleChildClick = (parentIndex) => {
  activeIndex.value = parentIndex;
  navItems.value[parentIndex].isOpen = false;
};

// 关闭所有下拉菜单（全局离开时触发）
const closeAllDropdowns = () => {
  navItems.value.forEach(item => {
    item.isOpen = false;
  });
};


</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-size: 26px;
}

.nav-bar {
  color: #000000;
  padding: 1em 50px;
  text-align: center;
  transition: all 0.3s;
  padding-right: 50px;
}

.nav-bar:hover {
    background-color: rgba(255, 255, 255, 0);
    color: #000000;
    padding: 1em;
    text-align: center;
}

/* 外层容器：扩大鼠标检测范围 */
.nav-container {
  position: relative;
}

.nav-container > ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  justify-content: space-between;
}

/* 导航项容器：增加上下内边距，扩大hover范围 */
.nav-item {
  position: relative;
  margin-right: 20px;
  cursor: pointer;
  
}

.nav-item:last-child {
  margin-right: 0;
}

/* 主链接容器 */
.main-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

/* 主链接样式 */
.main-link {
  color: #000000;
  text-decoration: none;
  transition: color 0.2s;
}

.main-link:hover {
  color: #333;
}

/* 箭头样式与旋转动画 */
.arrow {
  display: inline-block;
  transition: transform 0.3s ease; /* 平滑旋转动画 */
  transform-origin: center; /* 关键：以自身中心为旋转点 */
  font-size: 18px;
  cursor: pointer;
}

/* 展开时顺时针旋转180度（原地旋转） */
.arrow-open {
  transform: rotate(180deg); /* 以中心为轴旋转 */
}

/* 下拉菜单外层容器 */
.dropdown-wrapper {
  position: absolute;
  top: 100%;
  left: 0;
  padding: 10px 0; /* 增加上下外间距，扩大检测范围 */
  width: 100%;
}

/* 下拉菜单样式：增加宽度和内边距 */
.dropdown-menu {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  min-width: 200px; /* 加宽菜单，更容易点击 */
  background-color: rgb(255, 255, 255, 0.9);
  border: 1px solid rgb(255, 255, 255);
  border-radius: 6px;
  padding: 10px 0;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  z-index: 100;
  list-style: none;
}

/* 下拉子项：增加高度和内边距，扩大点击区域 */
.dropdown-item {
  padding: 10px 25px; /* 增加上下内边距，点击区域更大 */
  text-align: left;
  transition: background-color 0.2s;
  display:flex;

  justify-content: center;
  align-items: center;
}

.dropdown-item:hover {
  background-color: #f5f5f5;
}

.dropdown-item a {
  color: #333;
  text-decoration: none;
  font-size: 22px;
  font-weight: 400;
  display: block; /* 整行可点击 */
  width: 100%;
}

/* 激活状态样式 */
.nav-item.active .main-link {
  color: #000;
}

/* 悬停效果 */
.nav-bar:hover {
  background-color: rgba(255, 255, 255, 0);
}


.nav-bar li {
    margin-right: 20px;
}

.nav-bar a {
    color: #000000;
    text-decoration: none;
}

.nav-bar a:hover {
    color: #000000bb;
}

.nav-bar:hover a {
    color: #000000;
}

.nav-bar li:last-child {
    margin-right: 0;
}
</style>