<template>
  <teleport to="body">
    <transition name="fade">
      <div v-if="show" class="mask" @click.self="hide">
        <div class="info-box">
          <img :src="getImageUrl(authorAvatar)" alt="Author Avatar" id="avatar">
          <ul>
            <li>网站作者: {{ author }}</li>
            <li><font-awesome-icon :icon="['fas', 'envelope']"></font-awesome-icon> 邮箱：<a
                href="mailto:changbingmushao@qq.com">{{ email }}</a></li>
            <li><font-awesome-icon :icon="['fab', 'github']"></font-awesome-icon> gitHub：<a
                href="https://github.com/shaozicln" target="_blank">{{ gitHub }}</a> </li>
          </ul>
        </div>
      </div>
    </transition>
  </teleport>
  <img :src="getImageUrl(authorAvatar)" alt="Author Avatar" id="avatar" @click="show = true">
</template>
  
  <script setup>
  // 获取全局URL属性
  import { getCurrentInstance } from 'vue';
  const instance = getCurrentInstance();
  const URL = instance?.appContext.config.globalProperties.URL;

  
const { proxy } = getCurrentInstance()
const getImageUrl = (imgName) => {
  return `${proxy.$imageBaseUrl}${imgName}`
}

  import { library } from '@fortawesome/fontawesome-svg-core'
  import { faGithub } from '@fortawesome/free-brands-svg-icons'
  import { faMicroscope, faSeedling, faSmile, faComment, faLightbulb, faBolt, faEnvelope } from '@fortawesome/free-solid-svg-icons'
  import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

  library.add(faGithub, faEnvelope)

  // 注册 FontAwesomeIcon 组件
  import { defineComponent } from 'vue';
  defineComponent({
    components: {
      FontAwesomeIcon
    }
  });

  import { ref } from 'vue'
  
  const authorAvatar = ref('Venti-2.jpg')
  const author = ref('长柄木勺')
  const email = ref('changbingmushao@qq.com')
  const gitHub = ref('shaozicln')
  const show = ref(false)
  
  function hide() {
    show.value = false
  }
  </script>
  
  <style scoped>
  .mask {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .info-box {
    background-color: var(--panel-bg);
    color: var(--text-primary);
    padding: 20px;
    border: 1px solid var(--panel-border);
    border-radius: 10px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    text-align: center;
    z-index: 100;
    width: 400px;
  }

  #avatar {
    width: 200px;
    height: 200px;
    border-radius: 50%;
    margin-bottom: 20px;
    margin: 0 auto;
  }


  .info-box ul {
    list-style: none;
    padding: 0;
    margin: 0;
    font-size: 20px;
  }

  .info-box ul li {
    list-style: none;
    padding: 0;
    margin: 0;
    font-size: 20px;
    color: var(--text-primary);
  }

  .info-box a,
  .info-box a:link,
  .info-box a:visited {
    color: var(--text-primary);
    text-decoration: none;
  }

  .info-box a:hover {
    color: var(--text-primary);
    opacity: 0.8;
  }

  .author-info {
    list-style-type: none;
    padding: 0;
  }


  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.5s;
  }

  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }
</style>