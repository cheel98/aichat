<template>
  <div class="user-settings">
    <h2>用户设置</h2>
    
    <div v-if="loading" class="loading">
      加载中...
    </div>
    
    <div v-else-if="error" class="error-message">
      {{ error }}
    </div>
    
    <form v-else @submit.prevent="saveSettings" class="settings-form">
      <div class="form-group">
        <label for="theme">主题</label>
        <select id="theme" v-model="form.theme">
          <option value="light">浅色</option>
          <option value="dark">深色</option>
        </select>
      </div>
      
      <div class="form-group">
        <label for="language">语言</label>
        <select id="language" v-model="form.language">
          <option value="zh-CN">简体中文</option>
          <option value="en-US">English</option>
        </select>
      </div>
      
      <div class="form-group checkbox-group">
        <label class="checkbox-label">
          <input 
            type="checkbox" 
            v-model="form.notificationEnabled" 
          />
          启用通知
        </label>
      </div>
      
      <div v-if="successMessage" class="success-message">
        {{ successMessage }}
      </div>
      
      <div class="form-actions">
        <button 
          type="submit" 
          class="btn-primary" 
          :disabled="submitLoading"
        >
          {{ submitLoading ? '保存中...' : '保存设置' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import userStore from '../store/userStore';

export default {
  name: 'UserSettings',
  
  data() {
    return {
      form: {
        theme: 'light',
        language: 'zh-CN',
        notificationEnabled: false
      },
      loading: true,
      submitLoading: false,
      error: '',
      successMessage: ''
    };
  },
  
  created() {
    this.fetchUserSettings();
  },
  
  methods: {
    // 获取用户设置
    async fetchUserSettings() {
      this.loading = true;
      this.error = '';
      
      try {
        const settings = await userStore.getUserSettings();
        this.form.theme = settings.theme || 'light';
        this.form.language = settings.language || 'zh-CN';
        this.form.notificationEnabled = settings.notification_enabled === 1;
        
        // 应用主题
        this.applyTheme(this.form.theme);
      } catch (error) {
        this.error = error.response?.data?.error || '获取用户设置失败';
      } finally {
        this.loading = false;
      }
    },
    
    // 保存设置
    async saveSettings() {
      this.submitLoading = true;
      this.error = '';
      this.successMessage = '';
      
      try {
        await userStore.updateSettings({
          theme: this.form.theme,
          language: this.form.language,
          notification_enabled: this.form.notificationEnabled
        });
        
        // 应用主题
        this.applyTheme(this.form.theme);
        
        this.successMessage = '设置保存成功';
        setTimeout(() => {
          this.successMessage = '';
        }, 3000);
      } catch (error) {
        this.error = error.response?.data?.error || '保存设置失败';
      } finally {
        this.submitLoading = false;
      }
    },
    
    // 应用主题
    applyTheme(theme) {
      document.documentElement.setAttribute('data-theme', theme);
    }
  }
};
</script>

<style scoped>
.user-settings {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

h2 {
  text-align: center;
  margin-bottom: 30px;
  color: var(--text-color);
}

.loading, .error-message {
  text-align: center;
  padding: 20px;
}

.loading {
  color: var(--text-color);
}

.error-message {
  color: #f56c6c;
}

.success-message {
  color: #67c23a;
  margin-bottom: 15px;
  padding: 10px;
  background-color: rgba(103, 194, 58, 0.1);
  border-radius: 4px;
  text-align: center;
}

.settings-form {
  background-color: var(--card-bg);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 10px;
  color: var(--text-color);
  font-weight: 500;
}

select {
  width: 100%;
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--input-bg);
  color: var(--text-color);
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23606266' d='M6 8.825L1.175 4 0 5.175 6 11.175 12 5.175 10.825 4 6 8.825z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 10px center;
  padding-right: 30px;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #409eff;
}

.form-actions {
  margin-top: 30px;
}

.btn-primary {
  width: 100%;
  background: #409eff;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 500;
}

.btn-primary:hover {
  background: #66b1ff;
}

.btn-primary:disabled {
  background: #a0cfff;
  cursor: not-allowed;
}
</style> 