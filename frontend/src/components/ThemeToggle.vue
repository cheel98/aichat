<template>
  <div class="theme-toggle">
    <button @click="toggleTheme" class="theme-button" :title="theme === 'dark' ? '切换到亮色主题' : '切换到暗色主题'">
      <i v-if="theme === 'dark'" class="bi bi-sun-fill"></i>
      <i v-else class="bi bi-moon-fill"></i>
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