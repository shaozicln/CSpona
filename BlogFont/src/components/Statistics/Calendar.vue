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
                  'active': month.value === displayedMonth && pickerYear === displayedYear,
                  'disabled': month.value > currentMonth && pickerYear === currentYear || pickerYear > currentYear
                }"
                :disabled="month.value > currentMonth && pickerYear === currentYear || pickerYear > currentYear"
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
          v-for="post in visiblePosts" 
          :key="post.id" 
          class="timeline-post"
          :style="getPostStyle(post)"
        >
          <div class="post-content">
            <div class="post-time">{{ formatPostTime(post.start) }} - {{ formatPostTime(post.end) }}</div>
            <div class="post-title">{{ post.title }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="calendar-footer">
      <div class="time-indicator">{{ displayedYear }}年{{ displayedMonth }}月</div>
      <div class="hint">右上角可切换不同时间</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// 模拟数据 - 包含精确的开始和结束时间
const posts = ref([
  { 
    id: 1, 
    title: 'AI技术：裸工作？', 
    category: 'AI',
    start: '2024-05-03T09:00:00',
    end: '2024-05-05T18:00:00'
  },
  { 
    id: 2, 
    title: 'vue /veud ? : 哪个更划算？', 
    category: 'Vue',
    start: '2024-05-12T10:30:00',
    end: '2024-05-14T15:00:00'
  },
  { 
    id: 3, 
    title: 'wpload 界面样式总结', 
    category: 'CSS',
    start: '2024-05-18T14:00:00',
    end: '2024-05-20T11:00:00'
  },
  { 
    id: 4, 
    title: 'React Hooks深入解析', 
    category: 'React',
    start: '2024-05-25T13:00:00',
    end: '2024-06-02T17:00:00' // 跨月的示例
  },
  { 
    id: 5, 
    title: 'TypeScript最佳实践', 
    category: 'TypeScript',
    start: '2024-05-28T09:00:00',
    end: '2024-05-30T18:00:00'
  }
])

// 当前日期信息
const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth() + 1)
const currentDate = computed(() => {
  return `今天 ${currentMonth.value}/${now.getDate()} ${['日','一','二','三','四','五','六'][now.getDay()]}`
})

const currentWeather = ref('晴')

// 显示的年份和月份
const displayedYear = ref(currentYear.value)
const displayedMonth = ref(currentMonth.value)

// 计算当月天数
const daysInMonth = computed(() => {
  const year = displayedYear.value
  const month = displayedMonth.value
  return new Date(year, month, 0).getDate()
})

// 获取当月第一天的日期对象
const firstDayOfMonth = computed(() => {
  return new Date(displayedYear.value, displayedMonth.value - 1, 1)
})

// 获取当月最后一天的日期对象
const lastDayOfMonth = computed(() => {
  return new Date(displayedYear.value, displayedMonth.value, 0)
})

// 筛选出在当前月有活动的文章
const visiblePosts = computed(() => {
  return posts.value.filter(post => {
    const postEnd = new Date(post.end)
    const postStart = new Date(post.start)
    return (postStart <= lastDayOfMonth.value && postEnd >= firstDayOfMonth.value)
  })
})

// 月份选择器相关
const showMonthPicker = ref(false)
const pickerYear = ref(displayedYear.value)
const months = [
  { label: '1月', value: 1 },
  { label: '2月', value: 2 },
  { label: '3月', value: 3 },
  { label: '4月', value: 4 },
  { label: '5月', value: 5 },
  { label: '6月', value: 6 },
  { label: '7月', value: 7 },
  { label: '8月', value: 8 },
  { label: '9月', value: 9 },
  { label: '10月', value: 10 },
  { label: '11月', value: 11 },
  { label: '12月', value: 12 },
]

// 格式化文章时间显示
const formatPostTime = (dateString) => {
  const date = new Date(dateString)
  return `${date.getMonth()+1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

// 计算文章卡片样式
const getPostStyle = (post) => {
  const postStart = new Date(post.start)
  const postEnd = new Date(post.end)
  
  // 计算卡片开始位置
  const startDay = Math.max(1, postStart.getDate())
  const startPercent = postStart.getDate() === startDay ? 
    (postStart.getHours() * 60 + postStart.getMinutes()) / 1440 * 100 : 0
  
  // 计算卡片结束位置
  const endDay = Math.min(daysInMonth.value, postEnd.getDate())
  const endPercent = postEnd.getDate() === endDay ? 
    (postEnd.getHours() * 60 + postEnd.getMinutes()) / 1440 * 100 : 100
  
  // 计算卡片宽度和位置
  const left = `${((startDay - 1) + startPercent / 100) / daysInMonth.value * 100}%`
  const width = `${((endDay - startDay) + (endPercent - startPercent) / 100) / daysInMonth.value * 100}%`
  
  return {
    left,
    width,
    backgroundColor: getPostColor(post)
  }
}

// 根据文章类别获取颜色
const getPostColor = (post) => {
  const colors = {
    'AI': '#4ecdc4',
    'Vue': '#45b7d1',
    'CSS': '#88d8b0',
    'React': '#ff6b6b',
    'TypeScript': '#ffa502',
    'default': '#a7dbd8'
  }
  return colors[post.category] || colors.default
}

// 切换月份选择器
function toggleMonthPicker() {
  showMonthPicker.value = !showMonthPicker.value
  if (showMonthPicker.value) {
    pickerYear.value = displayedYear.value
  }
}

// 选择年份
function prevYear() {
  if (pickerYear.value > 2024) {
    pickerYear.value--
  }
}

function nextYear() {
  if (pickerYear.value < currentYear.value) {
    pickerYear.value++
  }
}

// 选择月份
function selectMonth(month) {
  displayedYear.value = pickerYear.value
  displayedMonth.value = month
  showMonthPicker.value = false
}

// 点击外部关闭月份选择器
onMounted(() => {
  document.addEventListener('click', (e) => {
    if (!e.target.closest('.month-selector')) {
      showMonthPicker.value = false
    }
  })
})
</script>

<style scoped>
.calendar-container {
  font-family: 'Helvetica Neue', Arial, sans-serif;
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
  margin-top: 20px;
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
}

.timeline-post {
  position: absolute;
  height: 40px;
  border-radius: 4px;
  padding: 5px;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
}

.timeline-post:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.post-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.post-time {
  font-size: 0.7rem;
  opacity: 0.9;
  margin-bottom: 2px;
}

.post-title {
  font-size: 0.8rem;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
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