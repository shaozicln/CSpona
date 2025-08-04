<template>
  <div class="calendar-container">
    <div class="calendar-header">
      <h2>博客创作月历</h2>
      <div class="header-controls">
        <div class="current-date">
          <span>{{ currentDate }}</span>
          <span>{{ currentWeather }}</span>
        </div>
        <div class="month-selector">
          <button @click="toggleMonthPicker" class="month-button">
            {{ displayedYear }}年{{ displayedMonth }}月 ▼
          </button>
          <div v-if="showMonthPicker" class="month-picker">
            <div class="year-selector">
              <button @click="prevYear" class="nav-button">◀</button>
              <span>{{ pickerYear }}年</span>
              <button @click="nextYear" class="nav-button">▶</button>
            </div>
            <div class="month-grid">
              <button v-for="month in months" :key="month.value" @click="selectMonth(month.value)" :class="{
                'month-button': true,
                active:
                  month.value === displayedMonth &&
                  pickerYear === displayedYear,
                disabled:
                  (month.value > currentMonth &&
                    pickerYear === currentYear) ||
                  pickerYear > currentYear ||
                  !hasPostsInMonth[month.value],
                'no-posts':
                  !hasPostsInMonth[month.value] &&
                  !(
                    month.value > currentMonth && pickerYear === currentYear
                  ) &&
                  !(pickerYear > currentYear),
              }" :disabled="(month.value > currentMonth && pickerYear === currentYear) ||
                  pickerYear > currentYear ||
                  !hasPostsInMonth[month.value]
                  ">
                {{ month.label }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="timeline-container">
      <button @click="prevMonth" class="timeline-nav-button left" 
    :disabled="!hasPrevValidMonth">
        <</button>

          <div class="timeline-scale">
            <div v-for="day in daysInMonth" :key="day" class="scale-day">
              <div class="scale-number">{{ day }}</div>
              <div class="scale-line"></div>
            </div>
          </div>

          <div class="timeline-content">
            <div v-for="(post, index) in sortedVisiblePosts" :key="post.id" class="timeline-post" :style="{
              ...getPostStyle(post),
              top: `${index * 75 + 10}px`,
            }" @click="goToArticle(post.id)">
              <div class="post-content">
                <!-- 只保留创建时间，去掉结束时间 -->
                <div class="post-time">{{ formatPostTime(post.start) }}</div>
                <!-- 为标题添加动态滚动类 -->
                <div class="post-title" :class="{ 'scroll-title': shouldScrollTitle(post) }">
                  {{ post.title }}
                </div>
              </div>
            </div>
          </div>
    </div>

    <div class="calendar-footer">
      <div class="time-indicator">
        {{ displayedYear }}年{{ displayedMonth }}月
      </div>
      <div class="hint">右上角可切换不同时间</div>
    </div>

    <button @click="nextMonth" class="timeline-nav-button right" 
    :disabled="!hasNextValidMonth">></button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
// 获取全局URL属性
import { getCurrentInstance } from "vue";
const instance = getCurrentInstance();
const URL = instance?.appContext.config.globalProperties.URL;

// 模拟数据 - 包含精确的开始和结束时间
const posts = ref([]);

// 直接点击卡片跳转
const router = useRouter();
const goToArticle = (articleId) => {
  const encodedId = window.btoa(articleId.toString());
  router.push(`/articleContent/${encodedId}`);
};

// 从API获取文章数据
const fetchArticles = async () => {
  try {
    const response = await fetch(`${URL}/categories-with-articles`);
    const { data } = await response.json();

    // 提取所有文章并转换格式
    const allArticles = data.flatMap((category) =>
      category.Articles.filter((article) => article.CategoryId !== 1000).map(
        (article) => ({
          id: article.Id,
          title: article.Title,
          category: category.Name,
          categoryId: article.CategoryId,
          start: article.CreatedAt,
          end: article.UpdatedAt,
        })
      )
    );

    posts.value = allArticles;
    console.log("文章数据加载完成:", posts.value);
  } catch (error) {
    console.error("获取文章数据失败:", error);
  }
};

// 根据文章类别获取颜色
const getPostColor = (post) => {
  const categoryColors = {
    "Coding-Study": "#aed66b",
    "Project-Practice": "#a7dbd8",
    Plan: "#e9d2c8",
    "N-A-C-G": "#efcbc0",
    Travel: "#b3cbe2",
    测试: "#8fabda",
    default: "#e4e17d",
  };

  const categoryName = post.category || "";
  const matchedColor = Object.entries(categoryColors).find(
    ([name]) => name.toLowerCase() === categoryName.toLowerCase()
  );

  return matchedColor ? matchedColor[1] : categoryColors.default;
};

// 当前日期信息
const now = new Date();
const currentYear = ref(now.getFullYear());
const currentMonth = ref(now.getMonth() + 1);
const currentDay = ref(now.getDate());
const currentDate = computed(() => {
  return `今天 ${currentMonth.value}/${currentDay.value} ${["日", "一", "二", "三", "四", "五", "六"][now.getDay()]
    }`;
});

const currentWeather = ref("晴");

// 显示的年份和月份
const displayedYear = ref(currentYear.value);
const displayedMonth = ref(currentMonth.value);

// 计算当月天数
const daysInMonth = computed(() => {
  const year = displayedYear.value;
  const month = displayedMonth.value;
  return new Date(year, month, 0).getDate();
});

// 获取当月第一天的日期对象
const firstDayOfMonth = computed(() => {
  return new Date(displayedYear.value, displayedMonth.value - 1, 1);
});

// 获取当月最后一天的日期对象
const lastDayOfMonth = computed(() => {
  return new Date(displayedYear.value, displayedMonth.value, 0);
});

// 筛选出在当前月有活动的文章
// 修改 visiblePosts 计算属性
const visiblePosts = computed(() => {
  return posts.value.filter((post) => {
    const postStart = new Date(post.start);
    const postYear = postStart.getFullYear();
    const postMonth = postStart.getMonth() + 1; // 月份从0开始，+1转为实际月份

    // 核心：仅保留 年份和月份完全匹配当前显示年月 的文章
    const isSameYear = postYear === displayedYear.value;
    const isSameMonth = postMonth === displayedMonth.value;

    // 额外校验：不显示未来时间的文章（年份>当前年 或 同年份月份>当前月）
    const isFuture =
      postYear > currentYear.value ||
      (postYear === currentYear.value && postMonth > currentMonth.value);

    return isSameYear && isSameMonth && !isFuture;
  });
});

// 月份选择器相关
const showMonthPicker = ref(false);
const pickerYear = ref(displayedYear.value);
const months = [
  { label: "1月", value: 1 },
  { label: "2月", value: 2 },
  { label: "3月", value: 3 },
  { label: "4月", value: 4 },
  { label: "5月", value: 5 },
  { label: "6月", value: 6 },
  { label: "7月", value: 7 },
  { label: "8月", value: 8 },
  { label: "9月", value: 9 },
  { label: "10月", value: 10 },
  { label: "11月", value: 11 },
  { label: "12月", value: 12 },
];

// 格式化文章时间显示
const formatPostTime = (dateString) => {
  const date = new Date(dateString);
  return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()}`;
};

const getPostStyle = (post) => {
  const postStart = new Date(post.start);
  const postYear = postStart.getFullYear();
  const postMonth = postStart.getMonth() + 1;
  const startDay = postStart.getDate();

  if (postYear !== displayedYear.value || postMonth !== displayedMonth.value) {
    return { display: 'none' };
  }

  const left = `${((startDay - 1) / daysInMonth.value) * 100}%`;

  return {
    left,
    width: 'calc(100% - 40px)',  // 预留左侧颜色条空间
    '--post-category-color': getPostColor(post),
  };

};// 切换月份选择器
function toggleMonthPicker() {
  showMonthPicker.value = !showMonthPicker.value;
  if (showMonthPicker.value) {
    pickerYear.value = displayedYear.value;
  }
}

// 选择年份
function prevYear() {
  if (pickerYear.value > 2023) {
    // 设置最小年份
    pickerYear.value--;
  }
}

function nextYear() {
  if (pickerYear.value < currentYear.value) {
    pickerYear.value++;
  }
}

// 选择月份
function selectMonth(month) {
  // 确保选择的日期不超过当前日期
  if (pickerYear.value === currentYear.value && month > currentMonth.value) {
    return;
  }

  displayedYear.value = pickerYear.value;
  displayedMonth.value = month;
  showMonthPicker.value = false;
}

// 按开始时间排序的可见文章
const sortedVisiblePosts = computed(() => {
  return [...visiblePosts.value].sort((a, b) => {
    return new Date(a.start) - new Date(b.start);
  });
});

// 添加标题滚动条件计算属性
const shouldScrollTitle = computed(() => (post) => {
  const postStart = new Date(post.start);
  const startDay = postStart.getDate();
  // 条件：当月下旬（日期>20）或标题长度>15字
  return startDay > 20 || post.title.length > 15;
});

// 无文章按钮变灰
const hasPostsInMonth = computed(() => {
  const postsByMonth = {};

  // 初始化所有月份为false
  for (let i = 1; i <= 12; i++) {
    postsByMonth[i] = false;
  }

  // 仅检查当前选择年份的文章
  posts.value.forEach((post) => {
    const postStart = new Date(post.start);
    const postYear = postStart.getFullYear();
    const postMonth = postStart.getMonth() + 1; // 月份从0开始

    // 如果文章创建年份与当前选择年份相同，则标记该月份有文章
    if (postYear === pickerYear.value) {
      postsByMonth[postMonth] = true;
    }
  });

  return postsByMonth;
});


// 获取指定年份中所有有文章的月份（升序排列）
const getValidMonthsInYear = (year) => {
  const validMonths = [];
  posts.value.forEach(post => {
    const postStart = new Date(post.start);
    if (postStart.getFullYear() === year) {
      const month = postStart.getMonth() + 1;
      if (!validMonths.includes(month)) {
        validMonths.push(month);
      }
    }
  });
  return validMonths.sort((a, b) => a - b); // 按月份升序排列
};
// 所有有效月份（有文章+不晚于当前时间），按时间升序
const validMonths = computed(() => {
  const valid = new Set();
  
 valid.add(`${currentYear.value}-${currentMonth.value}`);

  posts.value.forEach(post => {
    const postStart = new Date(post.start);
    const year = postStart.getFullYear();
    const month = postStart.getMonth() + 1;
    
    // 只保留不晚于当前时间的月份
    const isNotFuture = 
      year < currentYear.value || 
      (year === currentYear.value && month <= currentMonth.value);
    
    if (isNotFuture) {
      valid.add(`${year}-${month}`);
    }
  });
  
  return Array.from(valid)
    .map(str => {
      const [y, m] = str.split('-').map(Number);
      return { year: y, month: m };
    })
    .sort((a, b) => a.year - b.year || a.month - b.month);
});

// 上一个有效月份
const prevMonth = () => {
  const current = { year: displayedYear.value, month: displayedMonth.value };
  const validList = validMonths.value;
  const currentIndex = validList.findIndex(
    item => item.year === current.year && item.month === current.month
  );
  
  if (currentIndex > 0) {
    const prev = validList[currentIndex - 1];
    displayedYear.value = prev.year;
    displayedMonth.value = prev.month;
  }
};

// 下一个有效月份
const nextMonth = () => {
  const current = { year: displayedYear.value, month: displayedMonth.value };
  const validList = validMonths.value;
  const currentIndex = validList.findIndex(
    item => item.year === current.year && item.month === current.month
  );
  
  if (currentIndex !== -1 && currentIndex < validList.length - 1) {
    const next = validList[currentIndex + 1];
    displayedYear.value = next.year;
    displayedMonth.value = next.month;
  }
};

// 按钮禁用状态
const hasPrevValidMonth = computed(() => {
  const current = { year: displayedYear.value, month: displayedMonth.value };
  const currentIndex = validMonths.value.findIndex(
    item => item.year === current.year && item.month === current.month
  );
  return currentIndex > 0;
});

const hasNextValidMonth = computed(() => {
  const current = { year: displayedYear.value, month: displayedMonth.value };
  const currentIndex = validMonths.value.findIndex(
    item => item.year === current.year && item.month === current.month
  );
  return currentIndex !== -1 && currentIndex < validMonths.value.length - 1;
});
// 点击外部关闭月份选择器
onMounted(() => {
  document.addEventListener("click", (e) => {
    if (!e.target.closest(".month-selector")) {
      showMonthPicker.value = false;
    }
  });
  fetchArticles();
});
</script>

<style scoped>
.calendar-container {
  font-family: "Helvetica Neue", Arial, sans-serif;
  /*background: linear-gradient(135deg, #f0f9ff 0%, #e6f7ff 100%);*/
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 15px rgba(0, 120, 150, 0.1);
  color: #2d3748;
  position: relative;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid rgba(69, 183, 209, 0.3);
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 20px;
}

.calendar-header h2 {
  margin: 0;
  color: #2b6cb0;
  font-size: 1.5rem;
  font-weight: 600;
}

.current-date {
  display: flex;
  gap: 10px;
  color: #4a5568;
  font-size: 0.95rem;
}

.month-selector {
  position: relative;
}

.month-button {
  background: #45b7d1;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
}

.month-button:hover {
  background: #3aa5c0;
}

.month-button.active {
  background: #2b6cb0;
  font-weight: bold;
}

.month-button.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #a0aec0;
}

.month-button.disabled.no-posts {
  background-color: #f8f9fa;
  color: #adb5bd;
  border: 1px dashed #dee2e6;
}

.month-picker {
  position: absolute;
  top: 100%;
  right: 0;
  background: white;
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 10;
  width: 180px;
}

.year-selector {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-weight: bold;
  color: #2b6cb0;
}

.nav-button {
  background: none;
  border: none;
  color: #2b6cb0;
  font-size: 1rem;
  cursor: pointer;
  padding: 0 8px;
}

.nav-button:hover {
  color: #1a4e8a;
}

.month-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.timeline-container {
  position: relative;
  height: 300px;
}

.timeline-scale {
  display: flex;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 30px;
  border-bottom: 1px solid #cbd5e0;
}

.scale-day {
  flex: 1;
  position: relative;
  text-align: center;
}

.scale-number {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  font-size: 0.8rem;
  color: #4a5568;
}

.scale-line {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 1px;
  height: 5px;
  background-color: #cbd5e0;
}

.timeline-content {
  position: absolute;
  top: 40px;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: auto;
  max-height: calc(100% - 40px);
  gap: 10px;
}

.timeline-post {
  position: absolute;
  height: 60px;
  padding: 8px 12px;
  color: #000;
  border-radius: 4px;
  transition: all 0.2s;
  box-sizing: border-box;
  border-left: 4px solid transparent;
}

/* 通过样式注入左侧颜色条（替代原背景色） */
.timeline-post::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  border-radius: 2px 0 0 2px;
  background-color: var(--post-category-color);
}

.post-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.post-time {
  font-size: 1rem;
  opacity: 0.9;
  margin-bottom: 2px;
}

.post-title {
  font-size: 1.1rem;
  font-weight: 500;
  white-space: nowrap;
  /* 强制不换行 */
  overflow: hidden;
  /* 隐藏溢出内容 */
  text-overflow: ellipsis;
  /* 默认省略号 */
  margin-top: 4px;
}

/* .scroll-title {
  text-overflow: clip; 
  animation: scroll-left 10s linear infinite; 
}


@keyframes scroll-left {
  0% {
    transform: translateX(0);
  }
  10% {
    transform: translateX(0); 
  }
  90% {
    transform: translateX(calc(-100% + 200px)); 
  }
  100% {
    transform: translateX(calc(-100% + 200px)); 
  }
}

.scroll-title:hover {
  animation-play-state: paused;
} */

.year-indicator {
  font-size: 0.6rem;
  opacity: 0.8;
}

.calendar-footer {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid rgba(69, 183, 209, 0.3);
  font-size: 0.85rem;
  color: #4a5568;
  display: flex;
  justify-content: space-between;
}

.hint {
  color: #4299e1;
  cursor: pointer;
}

/* 时间轴切换按钮基础样式 */
.timeline-nav-button {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  border: none;
  border-radius: 50%;
  background-color: #45b7d1;
  color: white;
  font-size: 1rem;
  cursor: pointer;
  z-index: 5;
  /* 确保在卡片上方 */
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 左侧按钮位置 */
.timeline-nav-button.left {
  
  width: 30px;
  height: 30px;
  left: -15px;
  /* 向左超出容器一点 */
}

/* 右侧按钮位置 */
.timeline-nav-button.right {
  width: 40px;
  height: 40px;
  right: -1px;
  /* 向右超出容器一点 */
}

/* 按钮交互效果 */
.timeline-nav-button:hover {
  background-color: #2b6cb0;
  transform: translateY(-50%) scale(1.1);
}

/* 禁用状态（无可用月份时） */
.timeline-nav-button:disabled {
  background-color: #cbd5e0;
  cursor: not-allowed;
  transform: translateY(-50%);
}

.timeline-nav-button span {
  display: inline-block;
  width: 100%;
  text-align: center;
}

/* 调整时间轴容器样式，避免按钮被遮挡 */
.timeline-container {
  position: relative;
  height: 300px;
  padding: 0 20px;
  /* 左右预留空间，避免按钮遮挡内容 */
  box-sizing: border-box;
}
</style>
