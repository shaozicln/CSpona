<template>
    <div>
        <button class="button" @click="advice">建议查看</button>
    </div>
    <div class="modal" v-if="showModal" @click.self="hideModal">
        <div class="modal-content">
            <div id="button-group">
                <button class="button" @click="sortByDefault">默认顺序</button>
                <button class="button" @click="sortByType">类型顺序</button>
                <button class="button" @click="sortByUsername">用户顺序</button>
            </div>
            <div class="advice-list">
                <div v-for="(advice, index) in advices" :key="advice.id">
                    <h1>{{ index + 1 }}</h1>
                    <p>U：{{ advice.Username }}</p>
                    <p>P：{{ advice.Type }}</p>
                    <p>C：{{ advice.Content }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>

import { ref } from 'vue';

const showModal = ref(false);
const advices = ref([]);

const advice = async () => {
    showModal.value = true;
    try {
        const response = await fetch('http://127.0.0.1:8081/advice')
        const data = await response.json()
        console.log(data)
        advices.value = data.data;
    } catch (error) {
        console.error(error)
    }
}

const sortByDefault = async () => {
    showModal.value = true;
    try {
        const response = await fetch('http://127.0.0.1:8081/advice')
        const data = await response.json()
        console.log(data)
        advices.value = data.data;
    } catch (error) {
        console.error(error)
    }
}

const sortByType = async () => {
    showModal.value = true;
    try {
        const response = await fetch('http://127.0.0.1:8081/advice')
        const data = await response.json()
        console.log(data)
        advices.value = data.data.sort((a, b) => {
            if (a.Type < b.Type) return -1;
            if (a.Type > b.Type) return 1;
            return 0;
        });
    } catch (error) {
        console.error(error)
    }
}

const sortByUsername = async () => {
    showModal.value = true;
    try {
        const response = await fetch('http://127.0.0.1:8081/advice')
        const data = await response.json()
        console.log(data)
        advices.value = data.data.sort((a, b) => {
            if (a.Username < b.Username) return -1;
            if (a.Username > b.Username) return 1;
            return 0;
        });
    } catch (error) {
        console.error(error)
    }
}

function hideModal() {
    showModal.value = false;
}
</script>

<style scope>
.button {
    width: 100px;
    height: 40px;
    background-color: rgb(139, 189, 234);
    /* 修改按钮颜色 */
    color: #fff;
    padding: 10px;
    border: none;
    border-radius: 20px;
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

.modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
}

.modal-content {
    background-color: #fff;
    opacity: 0.8;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    width: 40vw;
    height: 80vh;
    overflow-y: auto;
}

#button-group {
    display: flex;
    flex-direction: row;
    align-items: center;
    padding: 10px;
    gap: 10px;
}

.modal-button {
    width: 100px;
    height: 40px;
    background-color: rgb(139, 189, 234);
    color: #fff;
    padding: 10px;
    border: none;
    border-radius: 20px;
    cursor: pointer;
    position: relative;
    transition: all 0.3s ease-in-out;
    font-size: 20px;
}

.advice-list {
    padding: 20px;
}
</style>