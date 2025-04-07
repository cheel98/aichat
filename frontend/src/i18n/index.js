import { reactive, readonly } from 'vue';
import zhCN from './zh-CN.js';
import enUS from './en-US.js';

// 支持的语言
export const SUPPORTED_LANGUAGES = {
  'zh-CN': '简体中文',
  'en-US': 'English'
};

// 获取系统语言或存储的语言
const getDefaultLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage && SUPPORTED_LANGUAGES[savedLanguage]) {
    return savedLanguage;
  }
  
  // 尝试匹配浏览器语言
  const browserLang = navigator.language;
  if (browserLang.startsWith('zh')) {
    return 'zh-CN';
  }
  
  // 默认为中文
  return 'zh-CN';
};

// 语言包映射
const messages = {
  'zh-CN': zhCN,
  'en-US': enUS
};

// 创建响应式状态
const state = reactive({
  currentLanguage: getDefaultLanguage(),
  messages: {}
});

// 加载当前语言的消息
const loadMessages = () => {
  state.messages = messages[state.currentLanguage] || messages['zh-CN'];
};

// 初始化加载语言
loadMessages();

// 改变语言的方法
const setLanguage = (lang) => {
  if (SUPPORTED_LANGUAGES[lang]) {
    state.currentLanguage = lang;
    localStorage.setItem('language', lang);
    loadMessages();
    
    // 更改文档语言
    document.querySelector('html').setAttribute('lang', lang);
    
    // 更新网页标题
    document.title = state.messages.app.title;
    
    return true;
  }
  return false;
};

// 翻译方法
const t = (key) => {
  // 支持点表示法，例如 'app.title'
  const keys = key.split('.');
  let value = state.messages;
  
  for (const k of keys) {
    if (!value || typeof value !== 'object') {
      return key; // 如果找不到翻译，返回原key
    }
    value = value[k];
  }
  
  return value || key;
};

// 导出
export default {
  state: readonly(state),
  t,
  setLanguage
};