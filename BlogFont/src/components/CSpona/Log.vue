<template>
</template>
<script setup>  
import {onBeforeRouteLeave} from 'vue-router'
onBeforeRouteLeave((to, from) => {
  if (to.name === 'Articles') { // 仅当跳转到 Articles 路由时设置刷新标记
    sessionStorage.setItem('refreshAfterEnter', 'Articles');
  }
});
</script>
<style scoped> </style>