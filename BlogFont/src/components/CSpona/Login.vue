<template>
  <div class="background">
    <div id="con" :class="['con', { 'with-reset': userStore.showForgetPasswordDialog }]">
      <div id="con-signin">
        <form class="form" id="form-in" @submit.prevent="userStore.check()">
          <h2 id="tit">Login</h2>
          <input v-model="userStore.username" placeholder="请输入用户名" type="text" name="username" id="username" autocomplete="username" />
          <input v-model="userStore.email" placeholder="请输入邮箱" type="text" name="email" id="email" autocomplete="email" />
          <input v-model="userStore.password" placeholder="请输入密码" type="password" name="password" id="password" autocomplete="current-password" />
          <div class="form-meta">
            <button
              type="button"
              class="link-btn"
              @click="userStore.showForgetPasswordDialog = !userStore.showForgetPasswordDialog"
            >
              {{ userStore.showForgetPasswordDialog ? '收起重置' : '忘记密码？' }}
            </button>
          </div>
          <button type="submit" id="btn登录" class="button primary-btn">登录</button>
          <router-link :to="{ path: '/register' }" class="register-link">没有账号？去注册</router-link>
        </form>
      </div>

      <div id="con-dialog" v-if="userStore.showForgetPasswordDialog">
        <div class="dialog-content">
          <h3 class="reset-title">重置密码</h3>
          <form id="reset" @submit.prevent="userStore.resetPassword()">
            <input v-model="userStore.newUsername" placeholder="请输入用户名" type="text" autocomplete="username" />
            <input v-model="userStore.newEmail" placeholder="请输入邮箱" type="text" name="email" autocomplete="email" />
            <input v-model="userStore.newPassword" placeholder="请输入新密码" type="password" autocomplete="new-password" />
            <input
              v-model="userStore.newPasswordAgain"
              placeholder="请再次输入新密码"
              type="password"
              autocomplete="new-password"
              @keyup="userStore.checkpassword()"
            />
            <span id="attention" v-html="userStore.attention"></span>
            <button type="submit" class="button primary-btn">重置密码</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
const userStore = useUserStore();
</script>

<style scoped>
h2 {
  font-size: 60px;
}

.background {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 24px;
  box-sizing: border-box;
}

.con {
  width: min(420px, 92vw);
  padding: 28px 28px 32px;
  border: 1px solid var(--panel-border, #ddd);
  border-radius: 12px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: var(--panel-bg);
  color: var(--text-primary);
  transition: width 0.35s ease, padding 0.35s ease;
  display: flex;
  flex-direction: row;
  align-items: stretch;
  gap: 0;
  box-sizing: border-box;
}

.con.with-reset {
  width: min(860px, 96vw);
  gap: 28px;
}

#con-signin {
  flex: 1;
  min-width: 0;
}

#con-dialog {
  flex: 1;
  min-width: 0;
  padding-left: 28px;
  border-left: 1px solid var(--panel-border, rgba(0, 0, 0, 0.12));
  box-sizing: border-box;
}

#tit {
  text-align: center;
  margin-bottom: 20px;
}

.reset-title {
  text-align: center;
  font-size: 28px;
  margin: 0 0 18px;
  font-weight: 600;
}

.form,
#reset {
  display: flex;
  flex-direction: column;
}

.form-meta {
  display: flex;
  justify-content: flex-end;
  margin: -8px 0 16px;
}

.link-btn {
  border: none;
  background: transparent;
  padding: 0;
  color: var(--text-secondary);
  font-size: 15px;
  cursor: pointer;
  line-height: 1.4;
}

.link-btn:hover {
  color: var(--text-primary);
  text-decoration: underline;
}

.register-link {
  display: block;
  margin-top: 16px;
  text-align: center;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 16px;
}

.register-link:hover {
  color: var(--text-primary);
}

.con input {
  width: 100%;
  height: 42px;
  margin-bottom: 14px;
  padding: 10px 12px;
  border: 1px solid var(--input-border);
  border-radius: 8px;
  font-size: 18px;
  background-color: var(--input-bg);
  color: var(--input-text);
  box-sizing: border-box;
}

.primary-btn {
  width: 100%;
  height: 44px;
  margin-top: 4px;
  background-color: var(--btn-bg);
  color: var(--btn-text);
  padding: 10px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
  font-size: 20px;
}

.primary-btn:hover {
  background-color: var(--btn-bg-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.15);
}

#attention {
  min-height: 1.2em;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-secondary);
}

@media (max-width: 720px) {
  .con.with-reset {
    flex-direction: column;
    width: min(420px, 92vw);
  }

  #con-dialog {
    padding-left: 0;
    padding-top: 20px;
    border-left: none;
    border-top: 1px solid var(--panel-border, rgba(0, 0, 0, 0.12));
  }
}
</style>
