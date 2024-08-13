<script setup>
import MyHead from './components/MyHead.vue'
import TypingText from './components/Home.vue'
import MessageBoard from './components/MessageBoard.vue'
</script>

<template>
  <!-- <Head id="h"></Head> -->
  <img id="background" src="../public/Pictures/动漫背景2.jpg">
  <div id="app">
    <MyHead />
    <router-view />
  </div>

</template>



<style>
* {
  padding: 0px;
  margin: 0px;
  font-family: cursive;
}

#background {
  width: 100%;
  height: 100%;
  top: 0;
  opacity: 0.80;
  z-index: -1;
  position: fixed;
}
</style>