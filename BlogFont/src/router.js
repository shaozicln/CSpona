import { createRouter, createWebHistory } from 'vue-router'
import Home from './components/Home.vue'//首页
import MessageBoard from './components/MessageBoard.vue'//留言/友链
import Sort from './components/Sort.vue'
import Articles from './components/Articles.vue'
import Login from './components/Login.vue'
import Register from './components/Register.vue'
import Back from './components/Back.vue'
import Search from './components/Search.vue'
import Feedback from './components/Feedback.vue'
import Author from './components/Author.vue'
import AuthorBack from './components/AuthorBack.vue'
import Advice from './components/Advice.vue'
import Application from './components/Application.vue'
import ArticlesManage from './components/ArticlesManage.vue'
import ArticleContent from './components/ArticleContent.vue'

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
        path: '/articleContent',
        name: 'ArticleContent',
        component: ArticleContent,
        meta: {
            hideMyHead: true
        }
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router