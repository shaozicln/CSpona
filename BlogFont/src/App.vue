<script setup>
import MyHead from './components/MyHead.vue'
</script>

<template>
  <!-- <Head id="h"></Head> -->
  <img id="background" src="../public/Pictures/background.jpg">
  <div id="app">
    <MyHead />
    <router-view />
  </div>

</template>



<style>
* {
  padding: 0px;
  margin: 0px;
  /* font-family: cursive; */
  /* font-family:Marcellus SC, serif; */
  font-family: Cormorant SC, serif;

  /*background-color: #f5f5f5;*/
}

#background {
  width: 100%;
  height: 100%;
  top: 0;
  opacity: 0.80;
  z-index: -1;
  position: fixed;
}

::-webkit-scrollbar {
  width: 2px;
  /* 垂直滚动条的宽度 */
  height: 2px;
  /* 水平滚动条的高度 */
}
</style>