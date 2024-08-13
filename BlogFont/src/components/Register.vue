<template>
    <div class="background">
        <div id="con">
            <form action="#" class="form" id="form-up">
                <h2 id="tit">注册</h2>
                <select name="file" id="file" type="text" placeholder="请选择头像" v-model="avatarPath">
                    <option value="" disabled selected>请选择头像</option>
                    <option v-for="path in avatarPaths" :key="path" :value="path">{{ path }}</option>
                </select><br>
                <input v-model="username" placeholder="请输入用户名" type="text" name="username" id="username"><br>
                <input v-model="newEmail" placeholder="请输入邮箱" type="password" name="email" id="email"><br>
                <input v-model="password" placeholder="请输入密码" type="password" name="password" id="password"><br>
            </form>
            <button id="bs" @click="createUser"class="button">注册</button>
            <div>
                <br>
                <router-link :to="{ path: '/login' }" class="register-link">已有帐号? 去登录</router-link>
            </div>
        </div>

    </div>
</template>

<script setup>
import { ref } from 'vue';

const username = ref('');
const email = ref('');
const password = ref('');
const roleQx = ref('B');
const avatarPath = ref('')
const avatarPaths = ref([
    '31.jpg',
    '35.jpg'
])

const createUser = async () => {
    if (username.value.trim() === '' || password.value.trim() === '' || avatarPath.value.trim() === '') {
        alert("信息不完善,请重新填写");
        return 0;
    }
    const url = "http://127.0.0.1:8081/users";
    fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            Username: username.value,
            Email: email.value,
            Password: password.value,
            RoleQx: roleQx.value,
            Avatar: avatarPath.value,
        }),
    })
        .then((response) => {
            if (response.ok) {
                alert("注册成功");
                return response.json();
            }
            throw new Error("创建文章失败，状态码：" + response.status);
        })
        .catch((error) => console.error("Error:", error));
};
</script>

<style scoped>
.background {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
}

#con {
    padding-top: 50px;
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
    flex-direction: column;
}

#tit {
    text-align: center;
    margin-bottom: 20px;
}

#form-up {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.form input,
select {
    width: 100%;
    height: 40px;
    margin-bottom: 20px;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
    font-size: 20px;
}

select {
    width: 100%;
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

.register-link {
  color: #000000; /* black color */
  text-decoration: none;
  font-size:17px;
}

.register-link:hover {
  color: #00000099; /* gray color on hover */
}

/* 添加浮现效果 */
#bs {
    position: relative;
    transition: all 0.3s ease-in-out;
}

#bs:hover {
    transform: translateY(-5px);
    box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}
</style>