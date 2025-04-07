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
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
);
// 响应拦截器，统一处理错误
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // 处理网络错误
     // 忽略错误，将response原样返回给上层处理
     if (error.response) {
       // 如果服务器返回了响应，直接返回响应对象
       return error.response;
     } 
     return  
  }
);

// 响应拦截器，处理认证错误
// apiClient.interceptors.response.use(
//   (response) => {
//     return response;
//   },
//   (error) => {
//     if (error.response && error.response.status === 401) {
//       localStorage.removeItem('token');
//       localStorage.removeItem('user');
//       window.location.href = '/login';
//     }
//     return Promise.reject(error);
//   }
// );

// 用户认证相关API
export const authAPI = {
  // 用户注册
  register: (userData) => {
    return apiClient.post('/auth/register', userData);
  },
  
  // 用户登录
  login: (credentials) => {
    return apiClient.post('/auth/login', credentials);
  },
  
  // 用户登出
  logout: () => {
    return apiClient.post('/auth/logout');
  }
};

// 用户相关API
export const userAPI = {
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