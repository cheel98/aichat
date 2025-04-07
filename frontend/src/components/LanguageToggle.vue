<template>
  <div class="language-toggle">
    <button class="language-btn" @click="toggleDropdown">
      {{ currentLanguageLabel }}
      <i class="bi bi-caret-down-fill"></i>
    </button>
    
    <div v-if="showDropdown" class="language-dropdown">
      <div 
        v-for="(name, code) in SUPPORTED_LANGUAGES" 
        :key="code"
        class="language-item"
        :class="{ active: code === currentLanguage }"
        @click="changeLanguage(code)"
      >
        {{ name }}
      </div>
    </div>
  </div>
</template>

<script>
import i18n, { SUPPORTED_LANGUAGES } from '../i18n';

export default {
  name: 'LanguageToggle',
  
  data() {
    return {
      showDropdown: false,
      SUPPORTED_LANGUAGES
    };
  },
  
  computed: {
    currentLanguage() {
      return i18n.state.currentLanguage;
    },
    
    currentLanguageLabel() {
      return SUPPORTED_LANGUAGES[this.currentLanguage];
    }
  },
  
  mounted() {
    // 添加点击外部关闭下拉菜单
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
      const dropdown = this.$el.querySelector('.language-dropdown');
      const button = this.$el.querySelector('.language-btn');
      
      if (dropdown && !dropdown.contains(event.target) && !button.contains(event.target)) {
        this.showDropdown = false;
      }
    },
    
    changeLanguage(code) {
      i18n.setLanguage(code);
      this.showDropdown = false;
    }
  }
};
</script>

<style scoped>
.language-toggle {
  position: relative;
}

.language-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 6px 10px;
  font-size: 14px;
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s;
}

.language-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.language-btn .bi {
  font-size: 12px;
}

.language-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 4px;
  background: var(--chat-bg);
  min-width: 120px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.language-item {
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.3s;
  color: var(--text-color);
}

.language-item:hover {
  background: rgba(64, 158, 255, 0.1);
}

.language-item.active {
  color: var(--primary-color);
  background-color: rgba(64, 158, 255, 0.08);
}
</style> 