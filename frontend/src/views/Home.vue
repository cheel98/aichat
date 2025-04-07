<template>
  <div class="app-container">
    <!-- 头部区域 -->
    <header class="header">
      <div class="logo-container">
        <img src="../assets/logo.svg" alt="DeepSeek Logo" class="logo" />
        <h1>{{ $t('app.title') }}</h1>
      </div>
      <div class="header-right">
        <button v-if="!isSidebarOpen" class="toggle-history-btn" @click="toggleSidebar">
          <i class="bi bi-clock-history"></i>
        </button>
        <ThemeToggle />
        <LanguageToggle />
        <UserMenu @navigate="handleNavigation" />
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
          class="sidebar"
          ref="chatHistory"
        />
      </transition>
      
      <!-- 主要内容区域 - 聊天界面 -->
      <div 
        class="content-area"
        :class="{ 'with-sidebar': isSidebarOpen }"
      >
        <ChatInterface 
          :conversation-id="currentConversationId"
          @conversation-updated="handleConversationUpdate"
          ref="chatInterface"
        />
      </div>
    </div>
    
    <!-- 模态框 - 用户设置 -->
    <div v-if="showSettingsModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal-container">
        <div class="modal-header">
          <h2>{{ $t('settings.title') }}</h2>
          <button class="modal-close-btn" @click="closeModals">×</button>
        </div>
        <div class="modal-body">
          <UserSettings />
        </div>
      </div>
    </div>
    <!-- 模态框 - 用户资料 -->
    <div v-if="showProfileModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal-container">
        <div class="modal-header">
          <h2>{{ $t('profile.title') }}</h2>
          <button class="modal-close-btn" @click="closeModals">×</button>
        </div>
        <div class="modal-body">
          <UserProfile />
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
    
    closeModals() {
      this.showSettingsModal = false;
      this.showProfileModal = false;
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
}

.logo-container {
  display: flex;
  align-items: center;
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
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.toggle-history-btn {
  background: transparent;
  border: none;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--text-color);
  transition: background-color 0.2s;
}

.toggle-history-btn:hover {
  background-color: rgba(0, 0, 0, 0.06);
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
}

.sidebar {
  width: 250px;
  border-right: 1px solid var(--border-color);
  background-color: var(--chat-bg);
  flex-shrink: 0;
  overflow-y: auto;
}

.content-area {
  flex: 1;
  overflow: hidden;
  transition: margin-left 0.3s ease;
}

.content-area.with-sidebar {
  margin-left: 0;
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

/* 侧边栏过渡动画 */
.slide-enter-active, .slide-leave-active {
  transition: transform 0.3s ease;
}
.slide-enter-from, .slide-leave-to {
  transform: translateX(-100%);
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
</style> 