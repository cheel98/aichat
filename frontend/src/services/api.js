import axios from 'axios';
import { API_URL } from '../config';

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器，添加认证头
apiClient.interceptors.request.use(
  (config) => {
    // 如果是登录请求，则不添加认证头
    if (config.url.includes('/auth')) {
      return config;
    }
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }else{
      // 如果token不存在，则拦截请求
      return Promise.reject(new Error('未登录'));
    }
    return config;
  },
);
// 响应拦截器，处理认证错误
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
    return Promise.reject(error);
  }
);

// 用户认证相关API
export const authAPI = {
  // 用户注册
  register: (userData) => {
    return apiClient.post('/auth/register', userData);
  },
  
  // 用户登录
  login: (credentials) => {
    return apiClient.post('/auth/login', credentials);
  }
};

// 用户相关API
export const userAPI = {
  // 用户登出
  logout: () => {
    return apiClient.post('/user/logout');
  },
  // 获取用户资料
  getProfile: () => {
    return apiClient.get('/user/profile');
  },
  
  // 更新用户资料
  updateProfile: (profileData) => {
    return apiClient.put('/user/profile', profileData);
  },
  
  // 更新用户密码
  updatePassword: (passwordData) => {
    return apiClient.put('/user/password', passwordData);
  },
  
  // 获取用户设置
  getSettings: () => {
    return apiClient.get('/user/settings');
  },
  
  // 更新用户设置
  updateSettings: (settingsData) => {
    return apiClient.put('/user/settings', settingsData);
  }
};

export default apiClient; 