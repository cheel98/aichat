<template>
  <div class="theme-toggle">
    <button @click="toggleTheme" class="theme-button" :title="theme === 'dark' ? '切换到亮色主题' : '切换到暗色主题'">
      <svg v-if="theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="5"></circle>
        <line x1="12" y1="1" x2="12" y2="3"></line>
        <line x1="12" y1="21" x2="12" y2="23"></line>
        <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
        <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
        <line x1="1" y1="12" x2="3" y2="12"></line>
        <line x1="21" y1="12" x2="23" y2="12"></line>
        <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
        <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
      </svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
      </svg>
    </button>
  </div>
</template>

<script>
export default {
  name: 'ThemeToggle',
  
  data() {
    return {
      theme: 'dark'
    }
  },
  
  methods: {
    toggleTheme() {
      this.theme = this.theme === 'dark' ? 'light' : 'dark'
      document.documentElement.setAttribute('data-theme', this.theme)
      localStorage.setItem('theme', this.theme)
    },
    
    initTheme() {
      const savedTheme = localStorage.getItem('theme')
      if (savedTheme) {
        this.theme = savedTheme
        document.documentElement.setAttribute('data-theme', this.theme)
      } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
        // 如果用户系统偏好亮色主题
        this.theme = 'light'
        document.documentElement.setAttribute('data-theme', 'light')
      }
    }
  },
  
  mounted() {
    this.initTheme()
  }
}
</script>

<style scoped>
.theme-toggle {
  position: relative;
}

.theme-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background-color: transparent;
  border: none;
  border-radius: 50%;
  color: var(--text-color);
  cursor: pointer;
  transition: background-color 0.2s;
}

.theme-button:hover {
  background-color: rgba(128, 128, 128, 0.2);
}
</style> 