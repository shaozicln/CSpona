import { createRouter, createWebHistory } from "vue-router";
import Home from "./components/Home/Home.vue"; //首页
import MessageBoard from "./components/Friends/MessageBoard.vue"; //留言/友链
import Sort from "./components/Statistics/Sort.vue";
import Articles from "./components/Article/Articles.vue";
import Login from "./components/CSpona/Login.vue";
import Register from "./components/CSpona/Register.vue";
import Back from "./components/Back/Back.vue";
import Search from "./components/Search/Search.vue";
import Feedback from "./components/Back/Feedback.vue";
import Author from "./components/Author/Author.vue";
import AuthorBack from "./components/Author/AuthorBack.vue";
import Advice from "./components/Manage/Advice.vue";
import Application from "./components/Manage/Application.vue";
import ArticlesManage from "./components/Manage/ArticlesManage.vue";
import ArticleContent from "./components/Article/ArticleContent.vue";
import Slide from "./components/Slide.vue";
import User from "./components/CSpona/User.vue";
import Wanderland from "./components/Friends/Wanderland.vue";
import Log from "./components/CSpona/Log.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/message-board",
    name: "MessageBoard",
    component: MessageBoard,
  },
  {
    path: "/sort",
    name: "Sort",
    component: Sort,
  },
  {
    path: "/articles",
    name: "Articles",
    component: Articles,
  },
  {
    path: "/wanderland/:articleId",
    name: "Wanderland",
    component: Wanderland,
  },
  {
    path: "/log",
    name: "Log",
    component: Log,
  },
  {
    path: "/user",
    name: "User",
    component: User,
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/register",
    name: "Register",
    component: Register,
  },
  {
    path: "/search",
    name: "Search",
    component: Search,
  },
  {
    path: "/back",
    name: "Back",
    component: Back,
  },
  {
    path: "/feedback",
    name: "Feedback",
    component: Feedback,
  },
  {
    path: "/author",
    name: "Author",
    component: Author,
  },
  {
    path: "/author-back",
    name: "AuthorBack",
    component: AuthorBack,
  },
  {
    path: "/advice",
    name: "Advice",
    component: Advice,
  },
  {
    path: "/application",
    name: "Application",
    component: Application,
  },
  {
    path: "/articles-manage",
    name: "ArticlesManage",
    component: ArticlesManage,
  },
  {
    path: "/articleContent/:articleId",
    name: "ArticleContent",
    component: ArticleContent,
    meta: {
      hideMyHead: true,
    },
  },
  {
    path: "/slide",
    name: "Slide",
    component: Slide,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  if (sessionStorage.getItem('shouldRefresh') === 'true') {
    sessionStorage.removeItem('shouldRefresh'); // 清除标记
    location.reload(); // 刷新页面
    // 注意：这里不调用 next()，因为刷新后会自动加载目标路由
    return;
  }
  next();
});

export default router;
