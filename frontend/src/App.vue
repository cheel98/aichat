<template>
  <div class="app-container">
    <header class="header">
      <div class="logo-container">
        <img src="./assets/logo.svg" alt="DeepSeek Logo" class="logo" />
        <h1>DeepSeek AI 聊天</h1>
      </div>
      <div class="header-right">
        <ThemeToggle />
        <UserMenu @navigate="handleNavigation" />
      </div>
    </header>
    
    <!-- 主要内容区域 -->
    <component 
      :is="currentComponent" 
      v-if="currentComponent !== 'ChatInterface'"
      @back="showChatInterface"
    />
    <ChatInterface v-else />
  </div>
</template>

<script>
import ChatInterface from './components/ChatInterface.vue'
import ThemeToggle from './components/ThemeToggle.vue'
import UserMenu from './components/UserMenu.vue'
import UserProfile from './components/UserProfile.vue'
import UserSettings from './components/UserSettings.vue'

export default {
  name: 'App',
  components: {
    ChatInterface,
    ThemeToggle,
    UserMenu,
    UserProfile,
    UserSettings
  },
  
  data() {
    return {
      currentComponent: 'ChatInterface'
    };
  },
  
  methods: {
    handleNavigation(page) {
      switch (page) {
        case 'profile':
          this.currentComponent = 'UserProfile';
          break;
        case 'settings':
          this.currentComponent = 'UserSettings';
          break;
        default:
          this.currentComponent = 'ChatInterface';
      }
    },
    
    showChatInterface() {
      this.currentComponent = 'ChatInterface';
    }
  }
}
</script>

<style>
:root {
  --text-color: #2c3e50;
  --bg-color: #f5f7fa;
  --card-bg: #ffffff;
  --border-color: #dcdfe6;
  --input-bg: #ffffff;
}

[data-theme="dark"] {
  --text-color: #e5eaf3;
  --bg-color: #1a1a1a;
  --card-bg: #252525;
  --border-color: #444444;
  --input-bg: #333333;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html, body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  color: var(--text-color);
  background-color: var(--bg-color);
  transition: background-color 0.3s, color 0.3s;
}

#app {
  height: 100%;
}

.app-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 10px 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  height: 32px;
  width: auto;
}

.header h1 {
  color: var(--text-color);
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}
</style> 