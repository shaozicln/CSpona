<template>
  <div class="not-found-container">
    <h1 class="four-zero-four">404</h1>
    <p class="page-not-found">这个页面找不到了</p>
    <div class="tips">
      <p>可能跟温迪一起出游了吧...</p>
      <p>可能被作者吃掉了也说不定...</p>
      <p>可能和可莉一起去炸鱼了也没准呢...</p>
      <p>可能随着无名客登上列车航行宇宙了大概...</p>
      <p>可能不小心迷路到零号空洞了...</p>
      <p>可能去了很远的地方...</p>
    </div>
    <router-link to="/" class="back-home">返回首页 »</router-link>
  </div>
</template>

<script setup>
// 这里如果有需要可以引入 Vue Router 的 useRouter 等进行更多逻辑处理，当前简单示例暂时不需要额外逻辑
</script>

<style scoped>
.not-found-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #1f2937; /* 深色背景，可根据实际需求调整 */
  color: #fff; /* 文字颜色为白色 */
  text-align: center;
  padding: 0 20px;
}

.four-zero-four {
  font-size: 80px;
  font-weight: bold;
  margin-bottom: 20px;
  /* 让字体更圆润，可根据喜好调整字体，这里用通用的圆润处理方式 */
  font-family: "Helvetica Neue", Helvetica, Arial, sans-serif; 
}

.page-not-found {
  font-size: 24px;
  margin-bottom: 30px;
}

.tips {
  margin-bottom: 40px;
  line-height: 1.6;
}

.back-home {
  color: #3B82F6;
  text-decoration: none;
  font-size: 16px;
  border: 1px solid #3B82F6;
  padding: 10px 20px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.back-home:hover {
  background-color: #3B82F6;
  color: #fff;
}
</style>