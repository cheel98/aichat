<template>
  <div class="user-settings">
    <div v-if="loading" class="loading">
      <div class="loading-spinner"></div>
    </div>
    
    <div v-else-if="error" class="error-message">
      {{ error }}
    </div>
    
    <div v-else class="settings-content">
      <!-- 左侧导航标签 -->
      <div class="settings-nav">
        <div
          v-for="(item, index) in navItems"
          :key="index"
          :class="['nav-item', { active: activeTab === item.key }]"
          @click="activeTab = item.key"
        >
          {{ $t(item.label) }}
        </div>
      </div>
      
      <!-- 右侧内容区域 -->
      <div class="settings-body">
        <!-- 提示词设置 -->
        <div v-if="activeTab === 'prompt'" class="tab-content">
          <div class="setting-group-vertical">
            <div class="setting-label">{{ $t('settings.promptText') }}</div>
            <div class="setting-control">
              <textarea
                v-model="form.prompt"
                :placeholder="$t('settings.promptPlaceholder')"
                class="textarea-control"
                rows="6"
              ></textarea>
            </div>
          </div>
        </div>
        
        <!-- 规则设置 -->
        <div v-if="activeTab === 'rules'" class="tab-content">
          <div class="setting-group-vertical">
            <div class="setting-label">{{ $t('settings.rulesText') }}</div>
            <div class="setting-control">
              <textarea
                v-model="form.rules"
                :placeholder="$t('settings.rulesPlaceholder')"
                class="textarea-control"
                rows="6"
              ></textarea>
            </div>
          </div>
        </div>
        
        <!-- 通知设置 -->
        <div v-if="activeTab === 'notifications'" class="tab-content">
          <div class="setting-group">
            <div class="setting-label">{{ $t('settings.enableNotifications') }}</div>
            <div class="setting-control">
              <div class="switch-control">
                <input 
                  type="checkbox" 
                  id="notification-toggle" 
                  v-model="form.notificationEnabled"
                  class="toggle-input"
                />
                <label for="notification-toggle" class="toggle-label"></label>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 保存按钮 -->
        <div class="settings-action">
          <button 
            class="btn-primary save-btn" 
            :disabled="submitLoading"
            @click="saveSettings"
          >
            {{ submitLoading ? $t('settings.savingSettings') : $t('settings.saveSettings') }}
          </button>
        </div>
      </div>
    </div>
    
    <div v-if="successMessage" class="success-toast">
      {{ successMessage }}
    </div>
  </div>
</template>

<script>
import userStore from '../store/userStore';
import i18n from '../i18n';

export default {
  name: 'UserSettings',
  
  data() {
    return {
      form: {
        notificationEnabled: false,
        prompt: '',
        rules: ''
      },
      loading: true,
      submitLoading: false,
      error: '',
      successMessage: '',
      activeTab: 'prompt',
      navItems: [
        { key: 'prompt', label: 'settings.promptSettings' },
        { key: 'rules', label: 'settings.rulesSettings' },
        { key: 'notifications', label: 'settings.notifications' }
      ]
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
        this.form.notificationEnabled = settings.notification_enabled === 1;
        this.form.prompt = settings.prompt || '';
        this.form.rules = settings.rules || '';
      } catch (error) {
        this.error = error.response?.data?.error || this.$t('common.error');
      } finally {
        this.loading = false;
      }
    },
    
    // 保存设置
    async saveSettings() {
      if (this.submitLoading) return;
      
      this.submitLoading = true;
      this.error = '';
      
      try {
        await userStore.updateSettings({
          notification_enabled: this.form.notificationEnabled,
          prompt: this.form.prompt,
          rules: this.form.rules
        });
        
        this.showSuccessMessage(this.$t('settings.settingsSaved'));
      } catch (error) {
        this.error = error.response?.data?.error || this.$t('common.error');
      } finally {
        this.submitLoading = false;
      }
    },
    
    // 显示成功消息
    showSuccessMessage(message) {
      this.$message.success(message);
      setTimeout(() => {
          this.$emit('close');
        }, 500);
    }
  }
};
</script>

<style scoped>
.user-settings {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.loading-spinner {
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top: 3px solid var(--primary-color);
  border-radius: 50%;
  width: 24px;
  height: 24px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  padding: 16px;
  color: #f56c6c;
  text-align: center;
}

.settings-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.settings-nav {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}

.nav-item {
  padding: 12px 20px;
  cursor: pointer;
  font-weight: 500;
  color: var(--text-secondary);
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
}

.nav-item:hover {
  color: var(--text-color);
}

.nav-item.active {
  color: var(--primary-color);
  border-bottom-color: var(--primary-color);
}

.settings-body {
  flex-grow: 1;
  padding: 0 20px;
  overflow-y: auto;
}

.tab-content {
  padding-bottom: 20px;
}

.setting-group {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid var(--border-color);
}

.setting-group-vertical {
  padding: 16px 0;
  border-bottom: 1px solid var(--border-color);
}

.setting-label {
  font-weight: 500;
  margin-bottom: 8px;
}

.setting-control {
  margin-top: 8px;
}

.select-control {
  min-width: 120px;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
  background-color: var(--chat-bg);
  color: var(--text-color);
  font-size: 14px;
}

.textarea-control {
  width: 100%;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
  background-color: var(--chat-bg);
  color: var(--text-color);
  font-size: 14px;
  resize: vertical;
  font-family: var(--font-family);
}

.textarea-control:focus {
  outline: none;
  border-color: var(--primary-color);
}

.switch-control {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 22px;
}

.toggle-input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-label {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 22px;
}

.toggle-label:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

.toggle-input:checked + .toggle-label {
  background-color: var(--primary-color);
}

.toggle-input:checked + .toggle-label:before {
  transform: translateX(22px);
}

.divider {
  height: 1px;
  background-color: var(--border-color);
  margin: 16px 0;
}

.settings-action {
  display: flex;
  justify-content: flex-end;
  padding: 16px 0;
}

.save-btn {
  min-width: 120px;
  padding: 8px 16px;
  border-radius: 4px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  font-size: 14px;
  cursor: pointer;
}

.save-btn:hover {
  background-color: var(--secondary-color);
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.success-toast {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background-color: #67c23a;
  color: white;
  padding: 10px 20px;
  border-radius: 4px;
  font-size: 14px;
  z-index: 100;
  animation: fade-in-out 3s ease-in-out;
}

@keyframes fade-in-out {
  0%, 100% { opacity: 0; }
  10%, 90% { opacity: 1; }
}
</style> 