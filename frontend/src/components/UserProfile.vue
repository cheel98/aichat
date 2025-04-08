<template>
  <div class="user-profile">
    <div v-if="loading" class="loading">
      <div class="loading-spinner"></div>
    </div>
    
    <div v-else-if="error" class="error-message">
      {{ error }}
    </div>
    
    <div v-else class="profile-content">
      <!-- 基本信息 -->
      <div class="profile-section">
        <div class="profile-row">
          <div class="profile-label">{{ $t('profile.username') }}</div>
          <div class="profile-field">
            <input
              type="text"
              v-model="form.username"
              :placeholder="$t('profile.username')"
              class="text-input"
            />
          </div>
        </div>
        
        <div class="divider"></div>
        
        <div class="profile-row">
          <div class="profile-label">{{ $t('profile.email') }}</div>
          <div class="profile-field">
            <input
              type="email"
              v-model="form.email"
              :placeholder="$t('profile.email')"
              disabled
              class="text-input"
            />
          </div>
        </div>
      </div>
      
      <!-- 修改密码 -->
      <div class="profile-section">
        <div class="section-title">{{ $t('profile.changePassword') }}</div>
        
        <div class="profile-row">
          <div class="profile-label">{{ $t('profile.currentPassword') }}</div>
          <div class="profile-field">
            <input
              type="password"
              v-model="form.currentPassword"
              :placeholder="$t('profile.currentPassword')"
              class="text-input"
            />
          </div>
        </div>
        
        <div class="divider"></div>
        
        <div class="profile-row">
          <div class="profile-label">{{ $t('profile.newPassword') }}</div>
          <div class="profile-field">
            <input
              type="password"
              v-model="form.newPassword"
              :placeholder="$t('profile.newPassword')"
              class="text-input"
            />
          </div>
        </div>
        
        <div class="divider"></div>
        
        <div class="profile-row">
          <div class="profile-label">{{ $t('profile.confirmPassword') }}</div>
          <div class="profile-field">
            <input
              type="password"
              v-model="form.confirmPassword"
              :placeholder="$t('profile.confirmPassword')"
              class="text-input"
            />
          </div>
        </div>
      </div>
      
      <div v-if="passwordError" class="error-panel">
        {{ passwordError }}
      </div>
      
      <!-- 保存按钮 -->
      <div class="profile-actions">
        <button 
          class="btn-primary save-btn" 
          :disabled="submitLoading"
          @click="saveProfile"
        >
          {{ submitLoading ? $t('profile.savingChanges') : $t('profile.saveChanges') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import userStore from '../store/userStore';

export default {
  name: 'UserProfile',
  
  data() {
    return {
      form: {
        username: '',
        email: '',
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      loading: true,
      submitLoading: false,
      error: '',
      passwordError: '',
    };
  },
  
  created() {
    this.fetchUserProfile();
  },
  
  methods: {
    // 获取用户资料
    async fetchUserProfile() {
      this.loading = true;
      this.error = '';
      
      try {
        const profile = await userStore.getUserProfile();
        this.form.username = profile.username || '';
        this.form.email = profile.email || '';
      } catch (error) {
        this.error = error.response?.data?.error || this.$t('common.error');
      } finally {
        this.loading = false;
      }
    },
    
    // 保存资料
    async saveProfile() {
      // 验证密码
      if (this.form.newPassword || this.form.confirmPassword || this.form.currentPassword) {
        if (!this.form.currentPassword) {
          this.passwordError = this.$t('profile.passwordError.empty');
          return;
        }
        
        if (this.form.newPassword !== this.form.confirmPassword) {
          this.passwordError = this.$t('profile.passwordError.mismatch');
          return;
        }
        
        if (this.form.newPassword && this.form.newPassword.length < 6) {
          this.passwordError = this.$t('profile.passwordError.tooShort');
          return;
        }
      }
      
      this.submitLoading = true;
      this.passwordError = '';
      this.error = '';
      
      try {
        const profileData = {
          username: this.form.username
        };
        
        // 如果要更改密码
        if (this.form.newPassword && this.form.currentPassword) {
          profileData.current_password = this.form.currentPassword;
          profileData.new_password = this.form.newPassword;
        }
        
        await userStore.updateProfile(profileData);
        
        // 清空密码字段
        this.form.currentPassword = '';
        this.form.newPassword = '';
        this.form.confirmPassword = '';
        
        this.showSuccessMessage(this.$t('profile.profileSaved'));
        // 延迟关闭，让用户有时间看到成功消息
        setTimeout(() => {
          this.$emit('close');
        }, 500);
      } catch (error) {
        if (error.response?.data?.error?.includes('password')) {
          this.passwordError = error.response.data.error || this.$t('common.error');
        } else {
          this.error = error.response?.data?.error || this.$t('common.error');
        }
      } finally {
        this.submitLoading = false;
      }
    },
    
    // 显示成功消息
    showSuccessMessage(message) {
     this.$message.success(message);
    }
  }
};
</script>

<style scoped>
.user-profile {
  width: 100%;
  padding: 0;
  overflow-y: auto;
  position: relative;
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
  margin-bottom: 16px;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin-top: 8px;
}

.profile-section {
  display: flex;
  flex-direction: column;
}

.section-title {
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
  padding: 0 20px;
}

.profile-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
}

.profile-label {
  font-weight: 500;
  color: var(--text-color);
}

.profile-field {
  width: 60%;
}

.text-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background-color: var(--chat-bg);
  color: var(--text-color);
  font-size: 0.95rem;
}

.text-input:disabled {
  opacity: 0.7;
  background-color: rgba(0, 0, 0, 0.02);
}

.divider {
  height: 1px;
  background-color: var(--border-color);
  margin: 0;
}

.error-panel {
  background-color: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
  padding: 12px 20px;
  border-radius: 8px;
  font-size: 14px;
  margin: 0 20px;
}

.profile-actions {
  display: flex;
  justify-content: center;
  padding: 16px 0 24px;
}

.save-btn {
  min-width: 120px;
  padding: 10px 20px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.save-btn:hover {
  background-color: var(--secondary-color);
}

.save-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.success-toast {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 20px;
  background-color: #67c23a;
  color: white;
  border-radius: 4px;
  font-size: 14px;
  z-index: 1000;
  animation: fadeOut 3s forwards;
}

@keyframes fadeOut {
  0% { opacity: 1; }
  70% { opacity: 1; }
  100% { opacity: 0; }
}
</style> 