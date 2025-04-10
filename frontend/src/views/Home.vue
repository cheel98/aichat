<template>
  <div class="app-container" :class="{ 'collapsed-sidebar': !isSidebarOpen }">
    <!-- 顶部浮动控制栏，只在侧边栏折叠时显示 -->
    <div v-if="!isSidebarOpen" class="floating-controls">
      <div class="left-controls">
        <button 
          class="expand-btn"
          @click="toggleSidebar"
          aria-label="展开侧边栏"
        >
          <i class="bi bi-layout-sidebar"></i>
        </button>
      </div>
    </div>
    
    <!-- 头部区域 -->
    <header v-if="isSidebarOpen" class="header">
      <div class="logo-container">
        <img src="../assets/logo.svg" alt="DeepSeek Logo" class="logo" />
        <h1>{{ $t('app.title') }}</h1>
        <button class="toggle-sidebar-btn" @click="toggleSidebar" title="收起侧边栏">
          <i class="bi bi-layout-sidebar-inset"></i>
        </button>
      </div>
      <div class="header-right">
        <ThemeToggle />
        <LanguageToggle />
        <UserMenu @navigate="handleNavigation" @logout="handleLogout" @auth-success="handleAuthSuccess" />
      </div>
    </header>
    
    <!-- 主区域：侧边栏与内容区 -->
    <div class="main-content">
      <!-- 历史对话侧边栏 -->
      <transition name="slide">
        <ChatHistory 
          v-if="isSidebarOpen" 
          @select-conversation="loadConversation"
          @new-chat="startNewChat"
          @toggle-sidebar="toggleSidebar"
          @conversation-updated="handleConversationUpdate"
          class="sidebar"
          ref="chatHistory"
        />
      </transition>
      
      <!-- 主要内容区域 - 聊天界面 -->
      <div 
        class="content-area"
        :class="{ 'with-sidebar': isSidebarOpen, 'full-width': !isSidebarOpen }"
      >
        <ChatInterface 
          :conversation-id="currentConversationId"
          @conversation-updated="handleConversationUpdate"
          ref="chatInterface"
        />
      </div>
    </div>
    
    <!-- 模态框 - 用户设置 -->
    <div v-if="showSettingsModal"  class="modal-overlay" @click.self="closeModals">
      <div class="modal-container">
        <div class="modal-header">
          <h2>{{ $t('settings.title') }}</h2>
          <button class="modal-close-btn" @click="closeModals">×</button>
        </div>
        <div class="modal-body">
          <UserSettings @close="this.showSettingsModal = false" />
        </div>
      </div>
    </div>
    <!-- 模态框 - 用户资料 -->
    <div v-if="showProfileModal"  class="modal-overlay" @click.self="closeModals">
      <div class="modal-container">
        <div class="modal-header">
          <h2>{{ $t('profile.title') }}</h2>
          <button class="modal-close-btn" @click="closeModals">×</button>
        </div>
        <div class="modal-body">
          <UserProfile @close="this.showProfileModal = false" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ChatInterface from '../components/ChatInterface.vue';
import ChatHistory from '../components/ChatHistory.vue';
import ThemeToggle from '../components/ThemeToggle.vue';
import UserMenu from '../components/UserMenu.vue';
import UserProfile from '../components/UserProfile.vue';
import UserSettings from '../components/UserSettings.vue';
import LanguageToggle from '../components/LanguageToggle.vue';

export default {
  name: 'HomeView',
  components: {
    ChatInterface,
    ChatHistory,
    ThemeToggle,
    UserMenu,
    UserProfile,
    UserSettings,
    LanguageToggle
  },
  
  data() {
    return {
      isSidebarOpen: true,
      currentConversationId: null,
      showSettingsModal: false,
      showProfileModal: false
    };
  },

  methods: {
    closeProfileModal() {
      this.showProfileModal = false;
    },
    closeSettingsModal() {
      this.showSettingsModal = false;
    },
    closeModals() {
      this.showSettingsModal = false;
      this.showProfileModal = false;
    },
    handleLogout() {
      this.$refs.chatHistory.logout();
      this.$refs.chatInterface.logout();
      this.currentConversationId = null;
    },
    handleAuthSuccess(){
      console.log("handleAuthSuccess")
      this.$refs.chatHistory.fetchConversations();
    },
    handleNavigation(page) {
      switch (page) {
        case 'profile':
          this.showProfileModal = true;
          break;
        case 'settings':
          this.showSettingsModal = true;
          break;
        case 'logout':
          this.$router.push('/login');
          break;
      }
    },
    
    
    toggleSidebar() {
      this.isSidebarOpen = !this.isSidebarOpen;
    },
    
    loadConversation(conversationId) {
      this.currentConversationId = conversationId;
      // 通知ChatInterface加载特定的对话
      if (this.$refs.chatInterface) {
        this.$refs.chatInterface.loadConversation(conversationId);
      }
      
      // 更新URL以包含会话ID
      if (this.$router) {
        this.$router.replace({
          path: '/',
          query: { id: conversationId }
        });
      }
    },
    
    startNewChat() {
      this.currentConversationId = null;
      if (this.$refs.chatInterface) {
        this.$refs.chatInterface.clearMessages();
      }
      
      // 移除URL中的会话ID
      if (this.$router) {
        this.$router.replace({ path: '/' });
      }
    },
    
    handleConversationUpdate() {
      // 当对话更新时通知ChatHistory组件刷新列表
      if (this.$refs.chatHistory) {
        this.$refs.chatHistory.fetchConversations();
      }
      
      // 如果当前有活跃对话，更新页面标题
      if (this.currentConversationId) {
        this.updatePageTitle(this.currentConversationId);
      }
    },
    
    // 更新页面标题的方法
    async updatePageTitle(conversationId) {
      try {
        if (this.$refs.chatHistory) {
          const conversation = await this.$refs.chatHistory.getConversationById(conversationId);
          if (conversation && conversation.title) {
            document.title = conversation.title;
          }
        }
      } catch (error) {
        console.error('更新页面标题失败', error);
      }
    },
    
    // 从URL中获取会话ID
    getSessionIdFromUrl() {
      if (this.$route && this.$route.query.id) {
        const sessionId = this.$route.query.id;
        this.loadConversation(sessionId);
      }
    }
  },
  
  mounted() {
    // 组件挂载时从URL获取会话ID
    this.getSessionIdFromUrl();
  },
  
  watch: {
    // 监听URL变化
    '$route.query.id': function(newId) {
      if (newId && newId !== this.currentConversationId) {
        this.loadConversation(newId);
      } else if (!newId && this.currentConversationId) {
        this.startNewChat();
      }
    },
    
    // 监听对话ID变化，更新页面标题
    currentConversationId: {
      handler: async function(newId) {
        if (newId) {
          try {
            // 假设可以通过 this.$refs.chatHistory 获取会话标题
            // 或者通过 chatInterface 获取
            const conversation = await this.$refs.chatHistory.getConversationById(newId);
            if (conversation && conversation.title) {
              document.title = conversation.title || this.$t('app.title');
            } else {
              document.title = this.$t('app.title');
            }
          } catch (error) {
            console.error('获取会话标题失败', error);
            document.title = this.$t('app.title');
          }
        } else {
          document.title = this.$t('app.title');
        }
      },
      immediate: true
    }
  }
};
</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--background-color);
  color: var(--text-color);
  transition: all 0.3s ease;
}

.header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--chat-bg);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.header.collapsed {
  height: 50px;
  padding: 0 10px;
  justify-content: flex-end;
}

.collapsed-sidebar .header {
  box-shadow: none;
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar {
  width: 280px;
  border-right: 1px solid var(--border-color);
  background-color: var(--chat-bg);
  flex-shrink: 0;
  overflow-y: auto;
  z-index: 10;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.content-area {
  flex: 1;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  width: 100%;
  transform-origin: left center;
}

.content-area.with-sidebar {
  margin-left: 0;
  width: calc(100% - 280px);
  transform: translateX(0);
}

.content-area.full-width {
  width: 100%;
  padding-top: 10px;
  transform: translateX(0);
}

.logo-container {
  display: flex;
  align-items: center;
  position: relative; /* 为绝对定位的子元素提供参考 */
}

.logo {
  height: 32px;
  width: 32px;
  margin-right: 10px;
}

.logo-container h1 {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
  margin-right: 10px; /* 为侧边栏折叠按钮留出空间 */
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.toggle-sidebar-btn {
  background: transparent;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text-color);
  transition: all 0.2s;
  font-size: 16px;
  margin-left: 5px; /* 与标题保持一定距离 */
}

.toggle-sidebar-btn:hover {
  background-color: rgba(0, 0, 0, 0.06);
  color: var(--primary-color);
}

/* 侧边栏过渡动画 */
.slide-enter-active, .slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-leave-from, .slide-enter-to {
  transform: translateX(0);
  opacity: 1;
}

.slide-enter-from, .slide-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

/* 全屏模式下的样式调整 */
.collapsed-sidebar .content-area {
  height: 100vh;
  padding-top: 0;
  width: 100%;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateX(0);
}

.collapsed-sidebar .main-content {
  margin-left: 0;
}

/* 聊天界面过渡动画 */
.chat-interface-enter-active,
.chat-interface-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.chat-interface-enter-from,
.chat-interface-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.chat-interface-enter-to,
.chat-interface-leave-from {
  opacity: 1;
  transform: translateY(0);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  width: 500px;
  max-width: 90%;
  max-height: 85vh;
  background-color: var(--chat-bg);
  border-radius: 10px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 20px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.modal-close-btn {
  background: transparent;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text-color);
  font-size: 22px;
  transition: background-color 0.2s;
}

.modal-close-btn:hover {
  background-color: rgba(0, 0, 0, 0.06);
}

.modal-body {
  padding: 0;
  overflow-y: auto;
  max-height: calc(85vh - 60px);
}

/* 打开侧边栏按钮 */
.open-sidebar-btn {
  position: absolute;
  left: 20px;
  top: 20px;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: var(--card-bg);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 5;
  transition: all 0.2s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.open-sidebar-btn:hover {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

/* 折叠后顶部栏样式 */
.collapsed-sidebar .header {
  height: 0;
  min-height: 0;
  padding: 0;
  overflow: hidden;
  border-bottom: none;
  opacity: 0;
}

.collapsed-sidebar .content-area {
  height: 100vh;
  padding-top: 0;
  width: 100%;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@media (max-width: 768px) {
  .sidebar {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    z-index: 100;
    box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
  }
  
  .content-area.with-sidebar {
    margin-left: 0;
  }
  
  .modal-container {
    width: 95%;
    max-width: none;
  }
}

/* 浮动控制栏样式 */
.floating-controls {
  position: fixed;
  top: 15px;
  left: 15px;
  z-index: 100;
  display: flex;
  align-items: center;
}

.left-controls {
  display: flex;
  align-items: center;
}

.expand-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: var(--primary-color);
  color: white;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.expand-btn:hover {
  background-color: var(--secondary-color, #66b1ff);
}

.floating-user-menu {
  background-color: var(--card-bg);
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
</style> 