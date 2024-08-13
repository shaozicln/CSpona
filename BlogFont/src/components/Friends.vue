<template>
    <div id="friends" class="dividing">
        <h2>Shaoの友链</h2>
        <br>
        <div id="web">
            <ul>
                <li v-for="friend in friends" :key="friend.Id">
                    <a :href="friend.Web" target="_blank" class="aaa">
                        <img :src="'../Public/Pictures/' + friend.Img" :alt="friend.Name">
                        <span class="name">
                            {{ friend.Name }}<br>
                            {{ friend.Web }}<br>
                            {{ friend.Introduction }}<br>
                        </span>
                    </a>
                </li>
            </ul>
        </div>
        <br>
        <div id="send">
            <h3>有个的网站做得不错, 想放个友链 ?</h3>
            <router-link :to="{ path: '/feedback' }"><button @click="send()">点我联系作者发送友链</button></router-link>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const friends = ref([])
const sendVisible = ref(false)
const sendQuery = ref('')
const sendInput = ref(null)

onMounted(async () => {
    try {
        const response = await fetch('http://127.0.0.1:8081/friendsWeb', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        const data = await response.json()
        friends.value = data.data
    } catch (error) {
        console.error(error)
    }
})

const send = () => {
    sendVisible.value = true
    setTimeout(() => {
        sendInput.value.focus();
    }, 0);
}

const closeSend = () => {
    sendVisible.value = false
}

const performSend = async () => {
    // perform send logic here
}

</script>

<style scope>
#friends {
    height: 100vh;
    padding: 50px;
    display: flex;
    justify-self: center;
    align-items: center;
    display: grid;
    grid-template-rows: auto auto auto;
}

h2 {
    text-align: center;
    font-size: 30px;
    font-family: cursive;
}

h3 {
    text-align: center;
    font-size: 40px;
    font-family: cursive;
}

#web ul {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-gap: 30px;
    list-style: none;
    padding: 0;
    margin: 0;
}

.aaa {
    width: 100%;
    height: 320px;
    border-radius: 20px;
}

.name {
    position: absolute;
    top: 100%;
    left: 50%;
    transform: translate(-50%, -50%);
    color: #fff;
    font-size: 25px;
    font-weight: bold;
    opacity: 0;
    transition: top 0.5s, opacity 0.5s;
    /* animate top and opacity */
    border-radius: 20px;
}

img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 20px;
}

a {
    display: block;
    position: relative;
}

a:hover .name {
    top: 50%;
    opacity: 1;
    border-radius: 20px;
}

a:hover::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 20px;
    z-index: -1;
}

a:hover img {
    filter: brightness(0.5);
}

.send-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: rgba(0, 0, 0, 0.5);
}

.send-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: transparent;
    z-index: -10;
    pointer-events: none;
}

.send-form-container {
    background-color: rgba(255, 255, 255, 0.8);
    padding: 50px 50px 40px 50px;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    transition: all 0.3s ease-in-out;
}

.send-item {
    display: flex;
    align-items: center;
    margin-right: 20px;
}

.send-item input {
    margin-right: 5px;
    width: 100%;
    height: 40px;
    border: 1px solid #ccc;
    border-radius: 5px;
    font-size: 20px;
}

.send-item button {
    background-color: rgb(139, 189, 234);
    color: #fff;
    border: none;
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease-in-out;
    border-radius: 20px;
}

.send-item button:hover {
    background-color: rgba(139, 190, 234, 0.686);
    transform: translateY(-5px);
    box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}

#send {
    gap: 20px;
    height: 70vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

#send button {
    background-color: rgb(139, 189, 234);
    color: #fff;
    border: none;
    padding: 10px 30px;
    font-size: 20px;
    cursor: pointer;
    transition: all 0.3s ease-in-out;
    border-radius: 20px;
    width:100%;
}

#send button:hover {
    background-color: rgba(139, 190, 234, 0.8);
    transform: translateY(-5px);
    box-shadow: 0 5px 10px rgba(0, 0, 0, 0.2);
}
</style>