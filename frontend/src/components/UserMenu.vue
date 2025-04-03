<template>
  <div class="user-menu">
    <!-- 未登录状态 -->
    <button 
      v-if="!isAuthenticated" 
      class="login-button" 
      @click="showAuthModal = true"
    >
      <i class="icon-login"></i>
      登录
    </button>
    
    <!-- 已登录状态 -->
    <div v-else class="avatar-container" @click="toggleDropdown">
      <img 
        :src="user?.avatar?.String || '/default-avatar.png'" 
        alt="用户头像" 
        class="avatar"
      />
      <span class="username">{{ user?.username }}</span>
      
      <!-- 下拉菜单 -->
      <div v-if="showDropdown" class="dropdown-menu">
        <div class="dropdown-item" @click="navigateTo('profile')">
          <i class="icon-user"></i>个人资料
        </div>
        <div class="dropdown-item" @click="navigateTo('settings')">
          <i class="icon-settings"></i>设置
        </div>
        <div class="dropdown-divider"></div>
        <div class="dropdown-item" @click="logout">
          <i class="icon-logout"></i>退出登录
        </div>
      </div>
    </div>
    
    <!-- 登录模态框 -->
    <div v-if="showAuthModal" class="auth-modal">
      <div class="modal-content">
        <button class="close-button" @click="showAuthModal = false">×</button>
        <UserAuth @auth-success="onAuthSuccess" />
      </div>
    </div>
  </div>
</template>

<script>
import UserAuth from './UserAuth.vue';
import userStore from '../store/userStore';

export default {
  name: 'UserMenu',
  components: {
    UserAuth
  },
  
  data() {
    return {
      showDropdown: false,
      showAuthModal: false
    };
  },
  
  computed: {
    isAuthenticated() {
      return userStore.state.isAuthenticated;
    },
    
    user() {
      return userStore.state.user;
    }
  },
  
  mounted() {
    // 点击外部区域关闭下拉菜单
    document.addEventListener('click', this.handleOutsideClick);
  },
  
  beforeUnmount() {
    document.removeEventListener('click', this.handleOutsideClick);
  },
  
  methods: {
    toggleDropdown(event) {
      event.stopPropagation();
      this.showDropdown = !this.showDropdown;
    },
    
    handleOutsideClick(event) {
      const dropdown = this.$el.querySelector('.dropdown-menu');
      const avatar = this.$el.querySelector('.avatar-container');
      
      if (dropdown && !dropdown.contains(event.target) && !avatar.contains(event.target)) {
        this.showDropdown = false;
      }
    },
    
    navigateTo(page) {
      this.showDropdown = false;
      // 这里应该使用你的路由系统跳转到对应页面
      // 由于这个简单示例没有配置路由，我们使用事件来通知父组件
      this.$emit('navigate', page);
    },
    
    async logout() {
      this.showDropdown = false;
      
      try {
        await userStore.logout();
      } catch (error) {
        console.error('登出时发生错误:', error);
      }
    },
    
    onAuthSuccess() {
      this.showAuthModal = false;
    }
  }
};
</script>

<style scoped>
.user-menu {
  position: relative;
}

.login-button {
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--primary-color, #409eff);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 6px 12px;
  font-size: 14px;
  transition: all 0.3s;
}

.login-button:hover {
  background: var(--secondary-color, #66b1ff);
}

.icon-login {
  display: inline-block;
  width: 16px;
  height: 16px;
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='white' d='M11 7L9.6 8.4l2.6 2.6H2v2h10.2l-2.6 2.6L11 17l5-5-5-5zm9 12h-8v2h8c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2h-8v2h8v14z'/%3E%3C/svg%3E");
}

.avatar-container {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
}

.avatar-container:hover {
  background: rgba(0, 0, 0, 0.05);
}

.avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  object-fit: cover;
}

.username {
  color: var(--text-color);
  font-size: 14px;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 下拉菜单 */
.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: var(--card-bg);
  min-width: 150px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.dropdown-item {
  padding: 10px 15px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-color);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.dropdown-item:hover {
  background: rgba(64, 158, 255, 0.1);
}

.dropdown-divider {
  height: 1px;
  background: var(--border-color);
  margin: 5px 0;
}

/* 图标样式 */
[class^="icon-"] {
  display: inline-block;
  width: 16px;
  height: 16px;
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain;
}

.icon-user {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23606266' d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/%3E%3C/svg%3E");
}

.icon-settings {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23606266' d='M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58a.49.49 0 0 0 .12-.61l-1.92-3.32a.488.488 0 0 0-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.484.484 0 0 0-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58a.49.49 0 0 0-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z'/%3E%3C/svg%3E");
}

.icon-logout {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23606266' d='M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z'/%3E%3C/svg%3E");
}

/* 登录模态框 */
.auth-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  position: relative;
  background: var(--card-bg);
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.close-button {
  position: absolute;
  top: 10px;
  right: 10px;
  background: none;
  border: none;
  font-size: 24px;
  color: var(--text-color);
  cursor: pointer;
  z-index: 1;
}
</style> 