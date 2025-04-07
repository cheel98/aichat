import { reactive, readonly } from 'vue';
import { authAPI, userAPI } from '../services/api';

// 初始状态
const state = reactive({
  isAuthenticated: false,
  user: null,
  loading: false,
  error: null,
  settings: null
});

// 从本地存储恢复用户状态
const initializeStore = () => {
  const token = localStorage.getItem('token');
  const user = localStorage.getItem('user');
  
  if (token && user) {
    state.isAuthenticated = true;
    state.user = JSON.parse(user);
  }
};

// 用户操作方法
const actions = {
  // 用户注册
  async register(userData) {
    state.loading = true;
    state.error = null;
    
    try {
      const response = await authAPI.register(userData);
      return response.data;
    } catch (error) {
      state.error = error.response?.data?.error || '注册失败';
      throw error;
    } finally {
      state.loading = false;
    }
  },
  
  // 用户登录
  async login(credentials) {
    state.loading = true;
    state.error = null;
    
    try {
      const response = await authAPI.login(credentials);
      const { token, user } = response.data;
      
      // 保存令牌和用户信息
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      // 更新状态
      state.isAuthenticated = true;
      state.user = user;
      
      return user;
    } catch (error) {
      state.error = error.response?.data?.error || '登录失败';
      console.error('登录失败:', error);
      throw error;
    } finally {
      console.log('登录结束');
      state.loading = false;
    }
  },
  
  // 用户登出
  async logout() {
    state.loading = true;
    
    try {
      // 调用登出API
      await authAPI.logout();
    } catch (error) {
      console.error('登出时发生错误:', error);
    } finally {
      // 无论API调用是否成功，都清除本地存储和状态
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      state.isAuthenticated = false;
      state.user = null;
      state.settings = null;
      state.loading = false;
    }
  },
  
  // 获取用户资料
  async getUserProfile() {
    state.loading = true;
    state.error = null;
    
    try {
      const response = await userAPI.getProfile();
      const user = response.data;
      
      // 更新用户信息
      state.user = user;
      localStorage.setItem('user', JSON.stringify(user));
      
      return user;
    } catch (error) {
      state.error = error.response?.data?.error || '获取用户资料失败';
      throw error;
    } finally {
      state.loading = false;
    }
  },
  
  // 更新用户资料
  async updateProfile(profileData) {
    state.loading = true;
    state.error = null;
    
    try {
      await userAPI.updateProfile(profileData);
      await actions.getUserProfile(); // 刷新用户资料
      return true;
    } catch (error) {
      state.error = error.response?.data?.error || '更新用户资料失败';
      throw error;
    } finally {
      state.loading = false;
    }
  },
  
  // 更新用户密码
  async updatePassword(passwordData) {
    state.loading = true;
    state.error = null;
    
    try {
      await userAPI.updatePassword(passwordData);
      return true;
    } catch (error) {
      state.error = error.response?.data?.error || '更新密码失败';
      throw error;
    } finally {
      state.loading = false;
    }
  },
  
  // 获取用户设置
  async getUserSettings() {
    state.loading = true;
    state.error = null;
    
    try {
      const response = await userAPI.getSettings();
      state.settings = response.data;
      return state.settings;
    } catch (error) {
      state.error = error.response?.data?.error || '获取用户设置失败';
      throw error;
    } finally {
      state.loading = false;
    }
  },
  
  // 更新用户设置
  async updateSettings(settingsData) {
    state.loading = true;
    state.error = null;
    
    try {
      await userAPI.updateSettings(settingsData);
      await actions.getUserSettings(); // 刷新用户设置
      return true;
    } catch (error) {
      state.error = error.response?.data?.error || '更新用户设置失败';
      throw error;
    } finally {
      state.loading = false;
    }
  }
};

// 初始化状态
initializeStore();

// 导出只读状态和操作方法
export default {
  state: readonly(state),
  ...actions,
  // 添加直接访问属性
  get isAuthenticated() {
    return state.isAuthenticated;
  },
  get user() {
    return state.user;
  }
}; 