<template>
  <div class="background">
    <div id="con" :class="['con', {'centered': !userStore.showForgetPasswordDialog, 'not-centered': userStore.showForgetPasswordDialog}]">
      <div id="con-signin">
        <form class="form" id="form-in">
          <h2 id="tit">Login</h2>
          <input v-model="userStore.username" placeholder="请输入用户名" type="text" name="username" id="username" /><br />
          <input v-model="userStore.email" placeholder="请输入邮箱" type="text" name="email" id="email" /><br />
          <input v-model="userStore.password" placeholder="请输入密码" type="password" name="password" id="password" />
        </form>
        <a href="#" @click.prevent="userStore.forgetPassword()" id="forget">忘记密码?</a>
        <br />
        <div>
          <button id="btn登录" @click.prevent="userStore.check()" class="button">登录</button>
          <br />
          <div id="con-register">
            <br />
            <router-link :to="{ path: '/register' }" class="register-link">没有账号? 去注册</router-link>
          </div>
        </div>
      </div>
      
      <div id="con-dialog" v-if="userStore.showForgetPasswordDialog">
        <div class="dialog-content">
          <form id="reset">
            <input v-model="userStore.newUsername" placeholder="请输入用户名" type="text" /><br />
            <input v-model="userStore.newEmail" placeholder="请输入邮箱" type="text" name="email" id="email" /><br />
            <input v-model="userStore.newPassword" placeholder="请输入新密码" type="password" /><br />
            <input v-model="userStore.newPasswordAgain" placeholder="请再次输入新密码" type="password" @keyup="userStore.checkpassword()" /><br />
            <span id="attention" v-html="userStore.attention"></span>
            <button @click.prevent="userStore.resetPassword()" class="button">重置密码</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '@/stores/user';
const userStore = useUserStore();

import { onBeforeRouteLeave } from 'vue-router';
onBeforeRouteLeave((to, from) => {
  if (to.name === 'Articles') { // 仅当跳转到 Articles 路由时设置刷新标记
    sessionStorage.setItem('refreshAfterEnter', 'Articles');
  }
});
</script>

<style scoped>
h2 {
  font-size: 60px;
}

.centered {
  display: flex;
  justify-content: center;
  align-items: center;
}

.background {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.con {
  /* width: 600px; /* adjust the width to fit both login and dialog */
  padding-top: 25px;
  padding-bottom: 25px;
  width: 25vw;
  border: 1px solid #ddd;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: rgba(255, 255, 255, 0.5);
  transition: all 0.5s ease-in-out;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: all 0.5s ease-in-out;
}

.con.not-centered {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: row;
  width: 60vw;
  padding: 25px;
  gap: 20px;
}

.register-link {
  color: #000000;
  /* black color */
  text-decoration: none;
  font-size: 17px;
}

.register-link:hover {
  color: #00000099;
  /* gray color on hover */
}

dialog {
  position: relative;
  /* changed from absolute to relative */
  width: 80%;
  background-color: rgba(255, 255, 255, 0.5);
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  transition: all 0.5s ease-in-out;
  z-index: 1000 !important;
  display: flex !important;
  display: flex;
  align-items: center;
  justify-content: center;
}

#con-signin {
  width: 40%;
  /* adjust the width to fit the login form */
  display: inline-block;
  vertical-align: top;
}

#con-dialog {
  flex: 1;
  margin-left: 20px;
  transition: all 0.5s ease-in-out;
  transition-delay: 0.5s;
}

#con-register {
  width: 100%;
  /* adjust the width to 40% */
  display: inline-block;
  /* make the container inline-block */
  vertical-align: top;
  /* vertically align the container */
}

#tit {
  text-align: center;
  margin-bottom: 20px;
}

#forget {
  font-size: 15px;
}

.con input,
select {
  width: 100%;
  height: 40px;
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  font-size: 20px;
}

.button {
  width: 100px;
  height: 40px;
  background-color: rgb(139, 189, 234);
  /* 修改按钮颜色 */
  color: #fff;
  padding: 10px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  position: relative;
  transition: all 0.3s ease-in-out;
  font-size: 20px;
}

.button:hover {
  background-color: rgb(139, 189, 234);
  transform: translateY(-5px);
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.5s ease-in-out;
}

.slide-enter,
.slide-leave-to {
  right: -300px;
  /* adjust this value to match the width of the dialog */
}

.slide-enter-to,
.slide-leave {
  right: 0;
}

.slide-out {
  transform: translateX(300px);
  opacity: 0;
}

#reset {
  background-color: rgba(255, 255, 255, 0.493);
}

a {
  text-decoration: none;
  border: none;
  box-shadow: none;
  outline: none;
}

a:focus {
  outline: none;
  border: none;
  box-shadow: none;
}

/* Add floating effect to the "Reset" button */
#btn-reset {
  position: relative;
  transition: all 0.3s ease-in-out;
}

#btn-reset:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}
</style>
