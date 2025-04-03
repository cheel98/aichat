import { createApp } from 'vue'
import App from './App.vue'
import './assets/main.css'

// 在挂载应用之前，根据本地存储中的主题设置应用主题
const theme = localStorage.getItem('theme') || 'light'
document.documentElement.setAttribute('data-theme', theme)

createApp(App).mount('#app') 