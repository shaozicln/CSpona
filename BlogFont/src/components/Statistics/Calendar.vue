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
              <button
                v-for="month in months"
                :key="month.value"
                @click="selectMonth(month.value)"
                :class="{
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
                }"
                :disabled="
                  (month.value > currentMonth && pickerYear === currentYear) ||
                  pickerYear > currentYear ||
                  !hasPostsInMonth[month.value]
                "
              >
                {{ month.label }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="timeline-container">
      <div class="timeline-scale">
        <div v-for="day in daysInMonth" :key="day" class="scale-day">
          <div class="scale-number">{{ day }}</div>
          <div class="scale-line"></div>
        </div>
      </div>

      <div class="timeline-content">
        <div
          v-for="(post, index) in sortedVisiblePosts"
          :key="post.id"
          class="timeline-post"
          :style="{
            ...getPostStyle(post),
            top: `${index * 10 + 5}px`,
            'border-radius': getPostBorderRadius(post),
          }"
          @click="goToArticle(post.id)"
          @mouseenter="hoveredPost = post"
          @mouseleave="hoveredPost = null"
        >
          <div class="post-content">
            <div class="post-time">
              {{ formatPostTime(post.start) }} - {{ formatPostTime(post.end) }}
            </div>
            <div class="post-title">{{ post.title }}</div>
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
      category.Articles.map((article) => ({
        id: article.Id,
        title: article.Title,
        category: category.Name,
        start: article.CreatedAt,
        end: article.UpdatedAt,
      }))
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
  return `今天 ${currentMonth.value}/${currentDay.value} ${
    ["日", "一", "二", "三", "四", "五", "六"][now.getDay()]
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
const visiblePosts = computed(() => {
  return posts.value.filter((post) => {
    const postEnd = new Date(post.end);
    const postStart = new Date(post.start);

    // 确保文章年份不超过当前年份
    if (postStart.getFullYear() > currentYear.value) {
      return false;
    }

    // 如果是当前年份，确保月份不超过当前月份
    if (
      postStart.getFullYear() === currentYear.value &&
      postStart.getMonth() + 1 > currentMonth.value
    ) {
      return false;
    }

    // 检查文章时间范围是否与当前显示月份有交集
    return (
      postStart <= lastDayOfMonth.value && postEnd >= firstDayOfMonth.value
    );
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
  const postEnd = new Date(post.end);
  const firstDay = new Date(displayedYear.value, displayedMonth.value - 1, 1);
  const lastDay = new Date(displayedYear.value, displayedMonth.value, 0);

  // 处理跨月开始时间
  let startDay = 1;
  let startPercent = 0;
  if (postStart >= firstDay) {
    startDay = Math.max(1, postStart.getDate());
    startPercent =
      postStart.getDate() === startDay
        ? ((postStart.getHours() * 60 + postStart.getMinutes()) / 1440) * 100
        : 0;
  }

  // 处理跨月结束时间
  let endDay = daysInMonth.value;
  let endPercent = 100;
  if (postEnd <= lastDay) {
    endDay = Math.min(daysInMonth.value, postEnd.getDate());
    endPercent =
      postEnd.getDate() === endDay
        ? ((postEnd.getHours() * 60 + postEnd.getMinutes()) / 1440) * 100
        : 100;
  }

  // 计算卡片宽度和位置
  const left = `${
    ((startDay - 1 + startPercent / 100) / daysInMonth.value) * 100
  }%`;
  const width = `${
    ((endDay - startDay + (endPercent - startPercent) / 100) /
      daysInMonth.value) *
    100
  }%`;

  return {
    left,
    width,
    backgroundColor: getPostColor(post),
  };
};

// 切换月份选择器
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

// 计算文章卡片边缘圆角
const getPostBorderRadius = (post) => {
  const postStart = new Date(post.start);
  const postEnd = new Date(post.end);
  const firstDay = new Date(displayedYear.value, displayedMonth.value - 1, 1);
  const lastDay = new Date(displayedYear.value, displayedMonth.value, 0);

  // 判断是否跨月
  const startsBeforeMonth = postStart < firstDay;
  const endsAfterMonth = postEnd > lastDay;

  if (startsBeforeMonth && endsAfterMonth) {
    return "0";
  } else if (startsBeforeMonth) {
    return "0 4px 4px 0";
  } else if (endsAfterMonth) {
    return "4px 0 0 4px";
  }
  return "4px";
};

// 无文章按钮变灰
const hasPostsInMonth = computed(() => {
  const postsByMonth = {};
  
  // 初始化所有月份为false
  for (let i = 1; i <= 12; i++) {
    postsByMonth[i] = false;
  }
  
  // 检查当前选择年份的文章
  posts.value.forEach(post => {
    const postStart = new Date(post.start);
    const postEnd = new Date(post.end);
    const postStartYear = postStart.getFullYear();
    const postEndYear = postEnd.getFullYear();
    
    // 只有当文章时间范围与当前选择年份有交集时才处理
    if (postStartYear <= pickerYear.value && postEndYear >= pickerYear.value) {
      // 确定在当前选择年份内的实际开始和结束月份
      const startMonth = postStartYear < pickerYear.value ? 1 : postStart.getMonth() + 1;
      const endMonth = postEndYear > pickerYear.value ? 12 : postEnd.getMonth() + 1;
      
      // 标记这些月份为有文章
      for (let month = startMonth; month <= endMonth; month++) {
        postsByMonth[month] = true;
      }
    }
  });
  
  return postsByMonth;
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
  background: linear-gradient(135deg, #f0f9ff 0%, #e6f7ff 100%);
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
}

.timeline-post {
  position: relative;
  height: 58px;
  padding: 8px;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
  z-index: 1;
  box-sizing: border-box;
}

.timeline-post:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  z-index: 2;
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
  font-size: 1.2rem;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

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
</style>
