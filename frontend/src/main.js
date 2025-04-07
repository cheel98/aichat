import { createApp } from 'vue'
import App from './App.vue'
import './assets/main.css'
import i18n from './i18n'
import router from './router'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap-icons/font/bootstrap-icons.css'
import 'bootstrap'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { ElMessage } from 'element-plus'
// 设置HTML的语言属性
document.querySelector('html').setAttribute('lang', i18n.state.currentLanguage)
// 设置网页标题
document.title = i18n.t('app.title')

const app = createApp(App)
app.use(ElementPlus)
// 添加全局属性，在任何组件中都可以通过this.$t访问
app.config.globalProperties.$t = i18n.t
app.config.globalProperties.$i18n = i18n
// 添加全局路由实例，确保在组件中可以访问
app.config.globalProperties.$router = router

// 全局注册 $message 方法
app.config.globalProperties.$message = ElMessage

app.use(router)
app.mount('#app') 