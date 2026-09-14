import { defineStore } from "pinia";
import { ref } from "vue";
import { useRouter } from "vue-router";
import { apiFetch, apiJson } from "@/utils/api.js";

function persistProfile(data) {
  if (!data) return;
  if (data.id != null) localStorage.setItem("userId", String(data.id));
  if (data.username) localStorage.setItem("username", data.username);
  if (data.email) localStorage.setItem("email", data.email);
  if (data.avatar) localStorage.setItem("avatar", data.avatar);
  if (data.qx) localStorage.setItem("userQx", data.qx);
}

function clearProfile() {
  localStorage.removeItem("userId");
  localStorage.removeItem("username");
  localStorage.removeItem("email");
  localStorage.removeItem("avatar");
  localStorage.removeItem("userQx");
}

export const useUserStore = defineStore("user", () => {
  const attention = ref("");
  const username = ref("");
  const password = ref("");
  const email = ref("");
  const avatar = ref("");
  const userQx = ref("");
  const newUsername = ref("");
  const newEmail = ref("");
  const newPassword = ref("");
  const newPasswordAgain = ref("");
  const isLoggedIn = ref(!!localStorage.getItem("userId"));

  const router = useRouter();
  const showForgetPasswordDialog = ref(false);

  function checkpassword() {
    if (newPassword.value === newPasswordAgain.value) {
      attention.value = "<font color='green'>ok密码一致</font>";
    } else {
      attention.value = "<font color='red'>密码不一致</font>";
    }
  }

  /** 启动时用 Cookie 恢复会话（不碰 token） */
  const restoreSession = async () => {
    try {
      const { res, data } = await apiJson("/auth/me");
      if (!res.ok || !data?.data) {
        isLoggedIn.value = false;
        return false;
      }
      persistProfile(data.data);
      username.value = data.data.username || "";
      email.value = data.data.email || "";
      avatar.value = data.data.avatar || "";
      userQx.value = data.data.qx || "";
      isLoggedIn.value = true;
      return true;
    } catch {
      isLoggedIn.value = false;
      return false;
    }
  };

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
      const { res, data } = await apiJson("/login", {
        method: "POST",
        body: JSON.stringify({
          username: username.value,
          password: password.value,
          email: email.value,
        }),
      });

      if (!res.ok) {
        throw new Error("登录失败，状态码：" + res.status);
      }

      if (data.message === "登录成功" || data.msg === "登录成功") {
        alert("登录成功，点击确认进入");
        persistProfile({
          id: data.id,
          username: data.username,
          email: data.email,
          avatar: data.avatar,
          qx: data.qx,
        });
        isLoggedIn.value = true;
        await router.push("/");
        window.location.reload();
      } else if (data.message === "密码错误" || data.msg === "密码错误") {
        alert("密码错误，请重试");
      } else if (
        data.message === "喵喵喵？注册了吗就来登录？" ||
        data.msg === "喵喵喵？注册了吗就来登录？"
      ) {
        alert("喵喵喵？注册了吗就来登录？");
      } else {
        alert(data.msg || data.message || "登录失败");
      }
    } catch (error) {
      console.error("Error:", error);
      alert("登录失败, 请重试");
    }
  };

  const logout = async () => {
    try {
      await apiFetch("/logout", { method: "POST" });
    } catch {
      /* ignore */
    }
    clearProfile();
    isLoggedIn.value = false;
    await router.push("/login");
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

      const { res, data } = await apiJson("/password-reset", {
        method: "POST",
        body: JSON.stringify({
          username: newUsername.value,
          email: newEmail.value,
          password: newPasswordAgain.value,
        }),
      });

      if (!res.ok || data?.error) {
        alert(data?.error || "重置密码失败");
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
    userQx,
    isLoggedIn,
    newUsername,
    newPassword,
    newEmail,
    newPasswordAgain,
    showForgetPasswordDialog,
    checkpassword,
    check,
    logout,
    restoreSession,
    forgetPassword,
    resetPassword,
  };
});
