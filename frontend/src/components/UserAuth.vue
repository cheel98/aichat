<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-header">
        <img src="/logo.png" alt="Logo" class="logo">
        <h1>欢迎使用 AI 聊天</h1>
      </div>
      
      <!-- 标签切换 -->
      <div class="tabs">
        <div 
          class="tab" 
          :class="{ active: activeTab === 'login' }" 
          @click="activeTab = 'login'"
        >
          登录
        </div>
        <div 
          class="tab" 
          :class="{ active: activeTab === 'register' }" 
          @click="activeTab = 'register'"
        >
          注册
        </div>
      </div>
      
      <!-- 登录表单 -->
      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label for="account">账号</label>
          <div class="input-container">
            <i class="icon-user"></i>
            <input 
              type="text" 
              id="account" 
              v-model="loginForm.account" 
              placeholder="邮箱 / 手机号"
              required
            />
          </div>
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <div class="input-container">
            <i class="icon-lock"></i>
            <input 
              type="password" 
              id="password" 
              v-model="loginForm.password" 
              placeholder="请输入密码"
              required
            />
          </div>
        </div>
        
        <div v-if="loginError" class="error-message">{{ loginError }}</div>
        
        <div class="form-actions">
          <button type="submit" class="btn-primary" :disabled="loginLoading">
            {{ loginLoading ? '登录中...' : '登录' }}
          </button>
          
          <div class="additional-links">
            <a href="#" @click.prevent="activeTab = 'register'">没有账号？立即注册</a>
            <a href="#" @click.prevent="forgotPassword">忘记密码？</a>
          </div>
        </div>
      </form>
      
      <!-- 注册表单 -->
      <form v-else @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <div class="input-container">
            <i class="icon-user"></i>
            <input 
              type="text" 
              id="username" 
              v-model="registerForm.username" 
              placeholder="请设置用户名"
              required
              minlength="3"
              maxlength="50"
            />
          </div>
        </div>
        
        <div class="form-group">
          <label for="contact">联系方式</label>
          <div class="input-container">
            <i class="icon-contact"></i>
            <input 
              type="text" 
              id="contact" 
              v-model="registerForm.contact" 
              placeholder="邮箱或手机号"
              required
              @input="detectContactType"
            />
          </div>
        </div>
        
        <div class="form-group">
          <label for="reg-password">密码</label>
          <div class="input-container">
            <i class="icon-lock"></i>
            <input 
              type="password" 
              id="reg-password" 
              v-model="registerForm.password" 
              placeholder="请设置密码"
              required
              minlength="6"
            />
          </div>
        </div>
        
        <div class="form-group">
          <label for="confirm-password">确认密码</label>
          <div class="input-container">
            <i class="icon-lock"></i>
            <input 
              type="password" 
              id="confirm-password" 
              v-model="registerForm.confirmPassword" 
              placeholder="请再次输入密码"
              required
              minlength="6"
            />
          </div>
          <div v-if="passwordError" class="field-error">
            {{ passwordError }}
          </div>
        </div>
        
        <div v-if="registerError" class="error-message">{{ registerError }}</div>
        
        <div class="form-actions">
          <button type="submit" class="btn-primary" :disabled="registerLoading || !passwordsMatch">
            {{ registerLoading ? '注册中...' : '注册' }}
          </button>
          
          <div class="additional-links">
            <a href="#" @click.prevent="activeTab = 'login'">已有账号？立即登录</a>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import userStore from '../store/userStore';

export default {
  name: 'UserAuth',
  emits: ['auth-success'],
  
  data() {
    return {
      activeTab: 'login',
      
      // 登录表单
      loginForm: {
        account: '',
        password: '',
        login_type: 1 // 默认邮箱登录
      },
      loginLoading: false,
      loginError: '',
      
      // 注册表单
      registerForm: {
        username: '',
        email: '', // 邮箱
        phone: '', // 手机
        contact: '', // 邮箱或手机
        password: '',
        confirmPassword: '',
        login_type: 1 // 默认邮箱注册
      },
      registerLoading: false,
      registerError: '',
      passwordError: '',
      passwordsMatch: true
    };
  },
  
  computed: {
    // 检查两次输入的密码是否一致
    passwordsMatch() {
      // 如果两个输入框都有值，才进行比较
      if (this.registerForm.password && this.registerForm.confirmPassword) {
        return this.registerForm.password === this.registerForm.confirmPassword;
      }
      // 如果有一个输入框为空，则不显示错误
      return true;
    },
    
    // 密码错误信息
    passwordError() {
      return !this.passwordsMatch ? '两次输入的密码不一致' : '';
    }
  },
  
  methods: {
    // 检测联系方式类型（邮箱或手机号）
    detectContactType() {
      // 邮箱正则
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      // 手机号正则（简单验证）
      const phoneRegex = /^1[3-9]\d{9}$/;
      
      if (emailRegex.test(this.registerForm.contact)) {
        this.registerForm.loginType = 1; // 邮箱类型
      } else if (phoneRegex.test(this.registerForm.contact)) {
        this.registerForm.loginType = 2; // 手机号类型
      }
      // 如果都不匹配，保持默认值
    },
    
    // 处理登录
    async handleLogin() {
      this.loginLoading = true;
      this.loginError = '';
      
      // 自动检测登录类型
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      this.loginForm.login_type = emailRegex.test(this.loginForm.account) ? 1 : 2;
      
      try {
        await userStore.login(this.loginForm);
        this.$emit('auth-success');
      } catch (error) {
        this.loginError = error.response?.data?.error || '登录失败，请检查账号和密码';
      } finally {
        this.loginLoading = false;
      }
    },
    
    // 处理注册
    async handleRegister() {
      this.registerLoading = true;
      this.registerError = '';
      
      // 检测联系方式类型
      this.detectContactType();
      
      // 检查密码确认
      if (!this.passwordsMatch) {
        this.registerError = '两次输入的密码不一致';
        this.registerLoading = false;
        return;
      }
      
      // 准备数据
      const registerData = {
        username: this.registerForm.username,
        password: this.registerForm.password,
        login_type: this.registerForm.loginType
      };
      
      // 根据登录类型设置email或phone字段
      if (this.registerForm.loginType === 1) {
        registerData.email = this.registerForm.contact;
      } else {
        registerData.phone = this.registerForm.contact;
      }
      
      try {
        console.log(registerData);
        await userStore.register(registerData);
        
        // 注册成功后提示并切换到登录页
        this.activeTab = 'login';
        this.loginForm.account = this.registerForm.contact;
        
        // 清空注册表单
        this.registerForm = {
          username: '',
          contact: '',
          password: '',
          confirmPassword: '',
          loginType: 1
        };
        
        // 显示成功消息
        this.loginError = '注册成功，请登录';
      } catch (error) {
        this.registerError = error.response?.data?.error || '注册失败，请稍后再试';
      } finally {
        this.registerLoading = false;
      }
    },
    
    // 忘记密码
    forgotPassword() {
      alert('密码找回功能正在开发中，请联系管理员');
    }
  }
};
</script>

<style scoped>
.auth-page {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.auth-container {
  width: 100%;
  max-width: 420px;
  background: var(--card-bg);
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  padding: 30px;
}

.auth-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo {
  width: 60px;
  height: 60px;
  margin-bottom: 15px;
}

.auth-header h1 {
  font-size: 24px;
  color: var(--text-color);
  font-weight: 600;
}

/* 标签切换 */
.tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 12px 0;
  font-size: 16px;
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.tab.active {
  color: var(--primary-color, #409eff);
  font-weight: 500;
}

.tab.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: var(--primary-color, #409eff);
}

/* 表单样式 */
.auth-form {
  margin-top: 20px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  color: var(--text-color);
  font-weight: 500;
  font-size: 14px;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.input-container i {
  position: absolute;
  left: 12px;
  color: var(--text-secondary);
}

.input-container input {
  padding-left: 40px;
  height: 46px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background-color: var(--input-bg);
  color: var(--text-color);
  width: 100%;
  font-size: 15px;
}

.input-container input:focus {
  border-color: var(--primary-color, #409eff);
  outline: none;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.error-message {
  margin-bottom: 15px;
  color: #f56c6c;
  font-size: 14px;
}

.field-error {
  margin-top: 5px;
  color: #f56c6c;
  font-size: 12px;
}

.form-actions {
  margin-top: 25px;
}

.btn-primary {
  width: 100%;
  height: 46px;
  background-color: var(--primary-color, #409eff);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary:hover {
  background-color: var(--secondary-color, #66b1ff);
}

.btn-primary:disabled {
  background-color: rgba(64, 158, 255, 0.5);
  cursor: not-allowed;
}

.additional-links {
  margin-top: 15px;
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.additional-links a {
  color: var(--primary-color, #409eff);
  text-decoration: none;
}

.additional-links a:hover {
  text-decoration: underline;
}

/* 图标样式 */
[class^="icon-"] {
  display: inline-block;
  width: 20px;
  height: 20px;
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain;
}

.icon-user {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23909399' d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/%3E%3C/svg%3E");
}

.icon-lock {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23909399' d='M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z'/%3E%3C/svg%3E");
}

.icon-contact {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23909399' d='M22 6c0-1.1-.9-2-2-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6zm-2 0l-8 5-8-5h16zm0 12H4V8l8 5 8-5v10z'/%3E%3C/svg%3E");
}
</style> 