<template>
  <div class="user-profile">
    <h2>用户资料</h2>
    
    <div v-if="loading" class="loading">
      加载中...
    </div>
    
    <div v-else-if="error" class="error-message">
      {{ error }}
    </div>
    
    <form v-else @submit.prevent="updateProfile" class="profile-form">
      <div class="avatar-section">
        <div class="avatar">
          <img 
            :src="form.avatar || '/default-avatar.png'" 
            alt="用户头像" 
          />
        </div>
        <div class="avatar-upload">
          <label for="avatar">上传头像</label>
          <input 
            type="file" 
            id="avatar" 
            accept="image/*" 
            @change="handleAvatarUpload" 
          />
        </div>
      </div>
      
      <div class="form-group">
        <label for="username">用户名</label>
        <input 
          type="text" 
          id="username" 
          v-model="form.username" 
          placeholder="请输入用户名"
          required
          minlength="3"
          maxlength="50"
        />
      </div>
      
      <div class="form-group">
        <label for="email">邮箱</label>
        <input 
          type="email" 
          id="email" 
          v-model="form.email" 
          placeholder="未设置"
          disabled
        />
      </div>
      
      <div class="form-group">
        <label for="phone">手机号</label>
        <input 
          type="tel" 
          id="phone" 
          v-model="form.phone" 
          placeholder="未设置"
          disabled
        />
      </div>
      
      <div class="form-actions">
        <button 
          type="submit" 
          class="btn-primary" 
          :disabled="submitLoading"
        >
          {{ submitLoading ? '保存中...' : '保存修改' }}
        </button>
        <button 
          type="button" 
          class="btn-secondary" 
          @click="showPasswordModal = true"
        >
          修改密码
        </button>
      </div>
    </form>
    
    <!-- 修改密码的模态框 -->
    <div v-if="showPasswordModal" class="password-modal">
      <div class="modal-content">
        <h3>修改密码</h3>
        <form @submit.prevent="updatePassword">
          <div class="form-group">
            <label for="old-password">原密码</label>
            <input 
              type="password" 
              id="old-password" 
              v-model="passwordForm.oldPassword" 
              placeholder="请输入原密码"
              required
            />
          </div>
          
          <div class="form-group">
            <label for="new-password">新密码</label>
            <input 
              type="password" 
              id="new-password" 
              v-model="passwordForm.newPassword" 
              placeholder="请输入新密码"
              required
              minlength="6"
            />
          </div>
          
          <div v-if="passwordError" class="error-message">
            {{ passwordError }}
          </div>
          
          <div class="form-actions">
            <button 
              type="submit" 
              class="btn-primary" 
              :disabled="passwordLoading"
            >
              {{ passwordLoading ? '提交中...' : '确认修改' }}
            </button>
            <button 
              type="button" 
              class="btn-secondary" 
              @click="showPasswordModal = false"
            >
              取消
            </button>
          </div>
        </form>
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
        avatar: '',
        email: '',
        phone: ''
      },
      loading: true,
      submitLoading: false,
      error: '',
      
      // 修改密码
      showPasswordModal: false,
      passwordForm: {
        oldPassword: '',
        newPassword: ''
      },
      passwordLoading: false,
      passwordError: '',
      
      // 成功消息
      successMessage: ''
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
        const user = await userStore.getUserProfile();
        this.form.username = user.username;
        this.form.avatar = user.avatar?.String || '';
        this.form.email = user.email?.String || '';
        this.form.phone = user.phone?.String || '';
      } catch (error) {
        this.error = error.response?.data?.error || '获取用户资料失败';
      } finally {
        this.loading = false;
      }
    },
    
    // 处理头像上传
    handleAvatarUpload(event) {
      const file = event.target.files[0];
      if (!file) return;
      
      // 简单实现，实际项目中应该上传到服务器
      // 这里仅为演示，使用本地URL
      const reader = new FileReader();
      reader.onload = (e) => {
        this.form.avatar = e.target.result;
      };
      reader.readAsDataURL(file);
    },
    
    // 更新用户资料
    async updateProfile() {
      this.submitLoading = true;
      this.error = '';
      
      try {
        await userStore.updateProfile({
          username: this.form.username,
          avatar: this.form.avatar
        });
        
        this.successMessage = '资料更新成功';
        setTimeout(() => {
          this.successMessage = '';
        }, 3000);
      } catch (error) {
        this.error = error.response?.data?.error || '更新资料失败';
      } finally {
        this.submitLoading = false;
      }
    },
    
    // 更新密码
    async updatePassword() {
      this.passwordLoading = true;
      this.passwordError = '';
      
      try {
        await userStore.updatePassword({
          old_password: this.passwordForm.oldPassword,
          new_password: this.passwordForm.newPassword
        });
        
        // 密码修改成功，关闭模态框
        this.showPasswordModal = false;
        
        // 重置表单
        this.passwordForm = {
          oldPassword: '',
          newPassword: ''
        };
        
        this.successMessage = '密码修改成功';
        setTimeout(() => {
          this.successMessage = '';
        }, 3000);
      } catch (error) {
        this.passwordError = error.response?.data?.error || '修改密码失败';
      } finally {
        this.passwordLoading = false;
      }
    }
  }
};
</script>

<style scoped>
.user-profile {
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

.profile-form {
  background-color: var(--card-bg);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.avatar-section {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 20px;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-upload input {
  display: none;
}

.avatar-upload label {
  display: inline-block;
  padding: 6px 12px;
  background: #409eff;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: var(--text-color);
}

input {
  width: 100%;
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--input-bg);
  color: var(--text-color);
}

input:disabled {
  background-color: rgba(0, 0, 0, 0.05);
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.btn-primary, .btn-secondary {
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  border: none;
}

.btn-primary {
  background: #409eff;
  color: white;
}

.btn-primary:hover {
  background: #66b1ff;
}

.btn-primary:disabled {
  background: #a0cfff;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f4f4f5;
  color: #606266;
}

.btn-secondary:hover {
  background: #e9e9eb;
}

/* 密码修改模态框 */
.password-modal {
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
  background: var(--card-bg);
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
}

.modal-content h3 {
  text-align: center;
  margin-bottom: 20px;
  color: var(--text-color);
}
</style> 