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
                <a
                  v-if="child.action === 'logout'"
                  href="#"
                  class="dropdown-link"
                  @click.prevent.stop="onLogout(index)"
                  @mouseenter.stop="keepDropdownOpen(index)"
                >
                  {{ child.text }}
                </a>
                <a
                  v-else-if="child.action === 'theme'"
                  href="#"
                  class="dropdown-link theme-switch-link"
                  @click.prevent.stop="onToggleTheme(index)"
                  @mouseenter.stop="keepDropdownOpen(index)"
                >
                  {{ child.text }}
                </a>
                <a
                  v-else-if="child.action === 'music'"
                  href="#"
                  class="dropdown-link"
                  @click.prevent.stop="onOpenMusic(index)"
                  @mouseenter.stop="keepDropdownOpen(index)"
                >
                  {{ child.text }}
                </a>
                <router-link 
                  v-else
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
import { getCurrentInstance, ref, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { useThemeStore } from '@/stores/theme';
import { useMusicStore } from '@/stores/music';

const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;
const route = useRoute();
const userStore = useUserStore();
const themeStore = useThemeStore();
const musicStore = useMusicStore();

const activeIndex = ref(0);

const username = localStorage.getItem("username");
const email = localStorage.getItem("email");
let userType = "";
if (username === "长柄木勺" && email === "changbingmushao@qq.com") {
  userType = "管理";
} else {
  userType = "反馈";
}

const accountChildren = computed(() => {
  const items = [
    { text: '个人中心', to: '/user' },
    { text: '登录/注册', to: '/login' },
    { text: '工具箱', to: '/tools' },
    { text: themeStore.switchLabel, to: '#theme', action: 'theme' },
    { text: '音乐设置', to: '#music', action: 'music' },
  ];
  if (userStore.isLoggedIn || localStorage.getItem('userId')) {
    items.push({ text: '退出登录', to: '#logout', action: 'logout' });
  }
  return items;
});

const navItems = ref([
  { 
    to: '/', 
    text: 'CSpona:)', 
    children: accountChildren.value,
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

watch(accountChildren, (kids) => {
  navItems.value[0].children = kids;
}, { immediate: true });

const toggleDropdown = (index) => {
  navItems.value[index].isOpen = !navItems.value[index].isOpen;
  navItems.value.forEach((item, i) => {
    if (i !== index) item.isOpen = false;
  });
  activeIndex.value = index;
};

const handleChildClick = (parentIndex) => {
  activeIndex.value = parentIndex;
  navItems.value[parentIndex].isOpen = false;
};

const keepDropdownOpen = (index) => {
  navItems.value[index].isOpen = true;
};

const closeAllDropdowns = () => {
  navItems.value.forEach(item => {
    item.isOpen = false;
  });
};

const onLogout = async (parentIndex) => {
  navItems.value[parentIndex].isOpen = false;
  await userStore.logout();
};

const onToggleTheme = (parentIndex) => {
  themeStore.toggle();
  navItems.value[parentIndex].isOpen = false;
};

const onOpenMusic = (parentIndex) => {
  musicStore.openSettings();
  navItems.value[parentIndex].isOpen = false;
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
  position: relative;
  z-index: 100;
  color: var(--nav-text);
  padding: 1em 50px;
  text-align: center;
  transition: all 0.3s;
  padding-right: 50px;
}

.nav-bar:hover {
    background-color: rgba(255, 255, 255, 0);
    color: var(--nav-text);
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
  color: var(--nav-text);
  text-decoration: none;
  transition: color 0.2s;
}

.main-link:hover {
  color: var(--text-muted);
}

/* 箭头样式与旋转动画 */
.arrow {
  display: inline-block;
  transition: transform 0.3s ease;
  transform-origin: center;
  font-size: 18px;
  cursor: pointer;
  color: var(--nav-text);
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
  min-width: 200px;
  background-color: var(--panel-bg);
  border: 1px solid var(--panel-border);
  border-radius: 6px;
  padding: 10px 0;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  z-index: 100;
  list-style: none;
}

/* 下拉子项：增加高度和内边距，扩大点击区域 */
.dropdown-item {
  padding: 10px 25px;
  text-align: left;
  transition: background-color 0.2s;
  display:flex;
  justify-content: center;
  align-items: center;
}

.dropdown-item:hover {
  background-color: var(--toc-active-bg);
}

.dropdown-item a,
.dropdown-item .dropdown-link {
  color: var(--text-primary);
  text-decoration: none;
  font-size: 22px;
  font-weight: 400;
  display: block;
  width: 100%;
  cursor: pointer;
}

.theme-switch-link {
  color: var(--text-primary);
}

.theme-switch-link:hover {
  color: var(--text-secondary);
}

/* 激活状态样式 */
.nav-item.active .main-link {
  color: var(--nav-text);
}

/* 悬停效果 */
.nav-bar:hover {
  background-color: transparent;
}

.nav-bar li {
    margin-right: 20px;
}

.nav-bar a {
    color: var(--nav-text);
    text-decoration: none;
}

.nav-bar a:hover {
    color: var(--text-muted);
}

.nav-bar:hover a {
    color: var(--nav-text);
}

.nav-bar li:last-child {
    margin-right: 0;
}
</style>