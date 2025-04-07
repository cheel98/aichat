import userStore from '../store/userStore';

// 不需要验证的路由
const publicRoutes = ['/login', '/register'];

// 路由守卫
export function setupRouterGuards(router) {
  router.beforeEach(async (to, from, next) => {
    // 检查是否是公开路由
    if (publicRoutes.includes(to.path)) {
      next();
      return;
    }

    // 检查是否已登录
    const isAuthenticated = userStore.isAuthenticated;
    
    // if (!isAuthenticated) {
    //   // 未登录，重定向到登录页
    //   next({
    //     path: '/login',
    //     query: { redirect: to.fullPath }
    //   });
    //   return;
    // }

    // 已登录，允许访问
    next();
  });
} 