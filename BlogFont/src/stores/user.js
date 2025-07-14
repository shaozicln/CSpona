import { defineStore } from "pinia";
import { ref } from "vue";
import { useRouter } from "vue-router";
import { getCurrentInstance } from "vue";

export const useUserStore = defineStore("user", () => {
  const instance = getCurrentInstance();
  const URL = instance?.appContext.config.globalProperties.URL;

  const attention = ref("");
  const username = ref("");
  const password = ref("");
  const email = ref("");
  const avatar = ref("");
  const newUsername = ref("");
  const newEmail = ref("");
  const newPassword = ref("");
  const newPasswordAgain = ref("");

  const router = useRouter();
  const showForgetPasswordDialog = ref(false);

  function checkpassword() {
    if (newPassword.value === newPasswordAgain.value) {
      attention.value = "<font color='green'>ok密码一致</font>";
    } else {
      attention.value = "<font color='red'>密码不一致</font>";
    }
  }

  const check = async () => {
    if (username.value.trim() === "") {
      alert("请输入名字");
      return;
    }
    if (email.value.trim() === "") {
      alert("请输入邮箱");
      return;
    }
    if (password.value.trim() === "") {
      alert("请输入密码");
      return;
    }

    try {
      const response = await fetch(`${URL}/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: username.value,
          password: password.value,
          email: email.value,
        }),
      });

      if (!response.ok) {
        throw new Error("登录失败，状态码：" + response.status);
      }

      const data = await response.json();

      if (data.message === "登录成功") {
        alert("登录成功，点击确认进入");
        const userId = data.id;
        const username = data.username;
        const email = data.email;
        const avatar = data.avatar;

        localStorage.setItem("userId", userId);
        localStorage.setItem("username", username);
        localStorage.setItem("email", email);
        localStorage.setItem("avatar", avatar);

        await router.push("/");
        window.location.reload();
      } else if (data.message === "密码错误") {
        alert("密码错误，请重试");
      } else if (data.message === "喵喵喵？注册了吗就来登录？") {
        alert("喵喵喵？注册了吗就来登录？");
      }
    } catch (error) {
      console.error("Error:", error);
      alert("登录失败, 请重试");
    }
  };

  const forgetPassword = () => {
    showForgetPasswordDialog.value = true;
    document.getElementById("con")?.classList?.remove("centered");
    document.getElementById("con")?.classList?.add("not-centered");
  };

  const resetPassword = async () => {
    try {
      if (newPassword.value !== newPasswordAgain.value) {
        alert("两次输入不一致");
        return;
      }

      const response = await fetch(`${URL}/users/${newUsername.value}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          password: newPasswordAgain.value,
        }),
      });

      const data = await response.json();

      if (data.error) {
        alert("重置密码失败");
      } else {
        alert("重置成功! 不要再忘记了");
        showForgetPasswordDialog.value = false;
        document.getElementById("con")?.classList?.remove("not-centered");
        document.getElementById("con")?.classList?.add("centered");
      }
    } catch (error) {
      console.error("Error:", error);
      alert("重置密码失败");
    }
  };

  return {
    attention,
    username,
    password,
    avatar,
    email,
    newUsername,
    newPassword,
    newEmail,
    newPasswordAgain,
    showForgetPasswordDialog,
    checkpassword,
    check,
    forgetPassword,
    resetPassword,
  };
});
