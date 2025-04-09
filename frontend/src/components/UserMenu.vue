<template>
  <div class="user-menu">
    <!-- 未登录状态 -->
    <button 
      v-if="!isAuthenticated" 
      class="login-button" 
      @click="showAuthModal = true"
    >
      <i class="bi bi-box-arrow-in-right"></i>
      {{ $t('app.login') }}
    </button>
    
    <!-- 已登录状态 -->
    <div v-else class="avatar-container" @click="toggleDropdown">
      <img 
        :src="user?.avatar?.String || '/default-avatar.png'" 
        :alt="$t('app.userAvatar')" 
        class="avatar"
      />
      <span class="username">{{ user?.username }}</span>
      
      <!-- 下拉菜单 -->
      <div v-if="showDropdown" class="dropdown-menu">
        <div class="dropdown-item" @click="navigateTo('profile')">
          <i class="bi bi-person"></i>{{ $t('app.profile') }}
        </div>
        <div class="dropdown-item" @click="navigateTo('settings')">
          <i class="bi bi-gear"></i>{{ $t('app.settings') }}
        </div>
        <div class="dropdown-divider"></div>
        <div class="dropdown-item" @click="logout">
          <i class="bi bi-box-arrow-right"></i>{{ $t('app.logout') }}
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
      return userStore.isAuthenticated;
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
        // 通知父组件用户已登出，需要清除聊天记录
        this.$emit('logout');
      } catch (error) {
        console.error(this.$t('app.logoutError'), error);
      }
    },
    
    onAuthSuccess() {
      this.showAuthModal = false;
      
      // 如果在模态框登录成功，跳转到首页
      if (this.$router) {
        this.$router.push('/');
      }
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

/* Bootstrap图标样式 */
.login-button .bi {
  font-size: 16px;
}

.dropdown-item .bi {
  font-size: 16px;
  color: var(--text-secondary);
}

.avatar-container {
  position: relative;
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
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 4px;
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