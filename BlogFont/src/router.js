import { createRouter, createWebHistory } from 'vue-router'
import Home from './components/Home/Home.vue'//首页
import MessageBoard from './components/Friends/MessageBoard.vue'//留言/友链
import Sort from './components/Statistics/Sort.vue'
import Articles from './components/Article/Articles.vue'
import Login from './components/Login/Login.vue'
import Register from './components/Login/Register.vue'
import Back from './components/Back/Back.vue'
import Search from './components/Search/Search.vue'
import Feedback from './components/Back/Feedback.vue'
import Author from './components/Author/Author.vue'
import AuthorBack from './components/Author/AuthorBack.vue'
import Advice from './components/Manage/Advice.vue'
import Application from './components/Manage/Application.vue'
import ArticlesManage from './components/Manage/ArticlesManage.vue'
import ArticleContent from './components/Article/ArticleContent.vue'
import Slide from './components/Slide.vue'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/message-board',
        name: 'MessageBoard',
        component: MessageBoard
    },
    {
        path: '/sort',
        name: 'Sort',
        component: Sort
    },
    {
        path: '/articles',
        name: 'Articles',
        component: Articles
    },
    {
        path: '/login',
        name: 'Login',
        component: Login
    },
    {
        path: '/register',
        name: 'Register',
        component: Register
    },
    {
        path: '/search',
        name: 'Search',
        component: Search
    },
    {
        path: '/back',
        name: 'Back',
        component: Back
    },
    {
        path: '/feedback',
        name: 'Feedback',
        component: Feedback
    },
    {
        path: '/author',
        name: 'Author',
        component: Author
    },
    {
        path: '/author-back',
        name: 'AuthorBack',
        component: AuthorBack
    },
    {
        path: '/advice',
        name: 'Advice',
        component: Advice
    },
    {
        path: '/application',
        name: 'Application',
        component: Application
    },
    {
        path: '/articles-manage',
        name: 'ArticlesManage',
        component: ArticlesManage
    },
    {
        path: '/articleContent/:articleId',
        name: 'ArticleContent',
        component: ArticleContent,
        meta: {
            hideMyHead: true
        }
    },
    {
        path: '/slide',
        name: 'Slide',
        component: Slide
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router