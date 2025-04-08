/**
 * 应用程序配置文件
 */

// API基础URL
export const API_BASE_URL = 'http://localhost:8080';

// API完整URL
export const API_URL = `${API_BASE_URL}/api`;

// 其他配置项可以在此处添加
export const CONFIG = {
  // 应用名称
  appName: 'Raspberry',
  
  // 版本号
  version: '1.0.0',
  
  // 最大消息长度
  maxMessageLength: 2000,
};

export default {
  API_BASE_URL,
  API_URL,
  CONFIG
}; 