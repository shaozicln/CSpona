<template>
    <div class="background">
        <div id="con" class="con centered">
            <div id="con-signin">
                <form action="#" class="form" id="form-in">
                    <h2 id="tit">Login</h2>
                    <input v-model="username" placeholder="请输入用户名" type="text" name="username" id="username"><br>
                    <input v-model="email" placeholder="请输入邮箱" type="text" name="email" id="email"><br>
                    <input v-model="password" placeholder="请输入密码" type="password" name="password" id="password">
                </form>
                <a href="#" @click="forgetPassword()" id="forget">忘记密码?</a>
                <br>
                <div>
                    <button id="btn登录" @click="check()" class="button">登录</button>
                    <br>
                    <div id="con-register">
                        <br>
                        <router-link :to="{ path: '/register' }" class="register-link">没有账号? 去注册</router-link>
                    </div>
                </div>
            </div>
            <div id="con-dialog" v-if="showForgetPasswordDialog" ref="forgetPasswordDialog" title="重置密码"
                :style="{ right: forgetPasswordDialogRight }">
                <transition name="slide">
                    <dialog v-if="showForgetPasswordDialog" ref="forgetPasswordDialog" title="重置密码"
                        :style="{ right: forgetPasswordDialogRight }" id="dialog">

                        <form id="reset">
                            <input v-model="newUsername" placeholder="请输入用户名" type="text"><br>
                            <input v-model="newEmail" placeholder="请输入邮箱" type="text" name="email" id="email"><br>
                            <input v-model="newPassword" placeholder="请输入新密码" type="password"><br>
                            <input v-model="newPasswordAgain" placeholder="请再次输入新密码" type="password"
                                @keyup="checkpassword()"><br>
                            <span id="attention" v-html="attention"></span>
                            <button @click="resetPassword()" class="button">重置密码</button>
                        </form>
                    </dialog>
                </transition>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const attention = ref('')
function checkpassword() {
    if (newPassword.value === newPasswordAgain.value) {
        attention.value = "<font color='green'>ok密码一致</font>"
    } else {
        attention.value = "<font color='red'>密码不一致</font>"
    }
}

const username = ref('');
const password = ref('');
const email = ref('');
const newUsername = ref('');
const newEmail = ref('');
const newPassword = ref('');
const newPasswordAgain = ref('');
const router = useRouter();

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
        fetch("http://127.0.0.1:8081/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: username.value,
                password: password.value,
                email: email.value,
            }),
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error("登录失败，状态码：" + response.status);
                }
                return response.json();
            })
            .then((data) => {
                if (data.message === "登录成功") {
                    alert("登录成功，点击确认进入");
                    var userId = data.id;
                    // var userQx = data.qx;
                    var username = data.username;
                    var email = data.email;
                    var password = data.password;
                    // var avatar = data.avatar;
                    // localStorage.setItem("avatar", avatar);
                    localStorage.setItem("userId", userId);
                    localStorage.setItem("username", username);
                    localStorage.setItem("email", email);
                    localStorage.setItem("password", password);
                    // localStorage.setItem("userQx", userQx);
                    router.push("/Back");
                    // window.location.reload()
                } else if (data.message === "密码错误") {
                    alert("密码错误，请重试");
                } else if (data.message === "喵喵喵？注册了吗就来登录？") {
                    alert("喵喵喵？注册了吗就来登录？");
                }
                console.log(data);
            })
    } catch (error) {
        console.error("Error:", error);
        alert("登录失败, 请重试");
    };
}

const forgetPasswordDialog = ref(null);
const forgetPasswordDialogRight = ref('0px');
const showForgetPasswordDialog = ref(false);

const forgetPassword = async () => {
    showForgetPasswordDialog.value = true;
    document.getElementById('con').classList.remove('centered');
    document.getElementById('con').classList.add('not-centered');
    forgetPasswordDialogRight.value = '320px';
}
const resetPassword = async () => {
    try {
        if (newPassword.value != newPasswordAgain.value) {
            alert("两次输入不一致");
            return;
        }
        const username = newUsername.value;
        fetch(`http://127.0.0.1:8081/users/${username}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                password: newPasswordAgain.value
            })
        })
            .then(respon => respon.json())
            .then(data => {
                if (data.error) {
                    alert("重置密码失败");
                } else {
                    alert("重置成功! 不要再忘记了");
                    forgetPasswordDialog.value.hide();
                    showForgetPasswordDialog.value = false; // <--- Set to false to hide the dialog
                    forgetPasswordDialogRight.value = '-300px';
                    document.getElementById('con').classList.remove('not-centered');
                    document.getElementById('con').classList.add('centered');
                    document.getElementById('con-dialog').classList.add('slide-out');
                }
            })
            .catch(error => console.error(error));
    } catch (error) {
        console.error("Error:", error);
        alert("重置密码失败");
    }
}

</script>

<style scope>
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
