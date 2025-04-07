import { createRouter, createWebHistory } from 'vue-router';
import { setupRouterGuards } from './guards';
import UserAuth from '../components/UserAuth.vue';
import Home from '../views/Home.vue';

// 定义路由
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: UserAuth,
    meta: {
      title: '登录'
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: UserAuth,
    meta: {
      title: '注册'
    }
  },
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      requiresAuth: true,
      title: '主页'
    },
    props: route => ({ 
      conversationId: route.query.id 
    })
  },
  // 捕获所有未匹配的路由并重定向到首页
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
});

// 设置路由守卫
// setupRouterGuards(router);

// 标题更新
router.afterEach((to) => {
  document.title = to.meta.title || '智能对话';
});

export default router; 