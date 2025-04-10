<template>
  <div class="chat-history">
    <div class="chat-history-header">
      <button class="new-chat-btn" @click="startNewChat">
        <span class="icon-plus"></span>
        {{ $t('history.newChat') }}
      </button>
    </div>
    
    <div v-if="loading" class="chat-history-loading">
      <div class="loading-spinner"></div>
    </div>
    
    <div v-else-if="conversations.length === 0" class="empty-history">
      <div class="empty-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="48" height="48" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>
      </div>
      <p>{{ $t('history.empty') }}</p>
      <button class="start-chat-btn" @click="startNewChat">{{ $t('history.newChat') }}</button>
    </div>
    
    <div v-else class="conversations-list">
      <div 
        v-for="(conversation, index) in conversations" 
        :key="conversation.session_id"
        class="conversation-item"
        :class="{ 'active': selectedConversationId === conversation.session_id }"
        @click="selectConversation(conversation.session_id)"
      >
        <div class="conversation-info">
          <div class="conversation-title">
            {{ conversation.title || $t('history.unnamed') }}
          </div>
          <div class="conversation-date">
            {{ formatDate(conversation.UpdatedAt) }}
          </div>
        </div>
        <div class="conversation-actions">
          <button class="action-btn" @click.stop="toggleActionMenu(conversation.session_id)">
            <i class="bi bi-three-dots-vertical"></i>
          </button>
          <div v-if="activeActionMenu === conversation.session_id" class="action-menu">
            <button class="menu-item" @click.stop="startRename(conversation)">
              <i class="bi bi-pencil"></i> {{ $t('history.rename') }}
            </button>
            <button class="menu-item delete" @click.stop="confirmDelete(conversation.session_id)">
              <i class="bi bi-trash3"></i> {{ $t('history.delete') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 重命名对话框 -->
    <div v-if="showRenameModal" class="rename-modal">
      <div class="modal-content">
        <h3>{{ $t('history.renameTitle') }}</h3>
        <input 
          type="text" 
          v-model="newConversationTitle" 
          class="rename-input"
          :placeholder="$t('history.renamePlaceholder')"
          @keyup.enter="renameConversation"
          ref="renameInput"
        >
        <div class="modal-actions">
          <button class="btn-cancel" @click="cancelRename">{{ $t('common.cancel') }}</button>
          <button class="btn-save" @click="renameConversation">{{ $t('common.save') }}</button>
        </div>
      </div>
    </div>
    
    <div v-if="showDeleteConfirm" class="delete-confirm">
      <p>{{ $t('history.deleteConfirm') }}</p>
      <div class="delete-actions">
        <button class="btn-cancel" @click="cancelDelete">{{ $t('common.cancel') }}</button>
        <button class="btn-delete" @click="deleteConversation">{{ $t('common.delete') }}</button>
      </div>
    </div>
  </div>
</template>

<script>
import { API_BASE_URL } from '../config';
import apiClient from '../services/api';

export default {
  name: 'ChatHistory',
  
  data() {
    return {
      conversations: [],
      loading: true,
      selectedConversationId: null,
      showDeleteConfirm: false,
      conversationToDelete: null,
      showRenameModal: false,
      newConversationTitle: '',
      activeActionMenu: null,
      conversationToRename: null
    }
  },
  
  created() {
    this.fetchConversations();
    // 添加全局点击事件监听器，用于关闭操作菜单
    document.addEventListener('click', this.handleClickOutside);
  },
  
  beforeUnmount() {
    // 组件销毁前移除事件监听器
    document.removeEventListener('click', this.handleClickOutside);
  },
  
  methods: {
    logout() {
      this.conversations = [];
      this.selectedConversationId = null;
      this.showDeleteConfirm = false;
      this.conversationToDelete = null;
    },
    async fetchConversations() {
      this.loading = true;
      try {
        const response = await apiClient.get('/chat/sessions');
        if (response.status !== 200) {
          throw new Error(this.$t('history.fetchError'));
        }
        const data = response.data;
        this.conversations = data.conversations || [];
      } catch (error) {
        console.error(this.$t('history.fetchErrorLog'), error);
      } finally {
        this.loading = false;
      }
    },
    selectConversation(id) {
      if (this.selectedConversationId !==id) {
        this.selectedConversationId = id;
        this.$emit('select-conversation', id);
      }
    },
    
    startNewChat() {
      this.selectedConversationId = null;
      this.$emit('new-chat');
    },
    
    confirmDelete(id) {
      this.conversationToDelete = id;
      this.showDeleteConfirm = true;
    },
    
    cancelDelete() {
      this.showDeleteConfirm = false;
      this.conversationToDelete = null;
    },
    
    async deleteConversation() {
      if (!this.conversationToDelete) return;
      
      try {
        const response = await apiClient.delete(`${API_BASE_URL}/api/chat/sessions/${this.conversationToDelete}`, {
        });
        console.log(response)
        if (response.status !== 200) {
          throw new Error(this.$t('history.deleteError'));
        }
        // 从列表中移除已删除的会话
        this.conversations = this.conversations.filter(
          conv => conv.session_id !== this.conversationToDelete
        );
        
        // 如果删除的是当前选中的会话，则清除选择
        if (this.selectedConversationId === this.conversationToDelete) {
          this.selectedConversationId = null;
          this.$emit('new-chat');
        }
      } catch (error) {
        console.error(this.$t('history.deleteErrorLog'), error);
      } finally {
        this.showDeleteConfirm = false;
        this.conversationToDelete = null;
      }
    },
    
    formatDate(dateString) {
      const date = new Date(dateString);
      const now = new Date();
      const isToday = date.toDateString() === now.toDateString();
      
      if (isToday) {
        return date.toLocaleTimeString(this.$i18n.state.currentLanguage, { hour: '2-digit', minute: '2-digit' });
      } else {
        return date.toLocaleDateString(this.$i18n.state.currentLanguage, { month: 'short', day: 'numeric' });
      }
    },
    
    // 根据ID获取会话信息
    async getConversationById(id) {
      // 首先从本地缓存查找
      const conversation = this.conversations.find(conv => conv.session_id === id);
      if (conversation) {
        return conversation;
      }
      
      // 如果本地没有找到，从服务器获取
      try {
        const response = await apiClient.get(`/chat/sessions/${id}`);
        if (response.status === 200 && response.data) {
          return response.data;
        }
      } catch (error) {
        console.error('获取会话详情失败', error);
      }
      
      return null;
    },
    
    toggleActionMenu(id) {
      this.activeActionMenu = this.activeActionMenu === id ? null : id;
    },
    
    startRename(conversation) {
      this.newConversationTitle = conversation.title || '';
      this.showRenameModal = true;
      this.conversationToRename = conversation.session_id;
      // 关闭操作菜单
      this.activeActionMenu = null;
    },
    
    cancelRename() {
      this.showRenameModal = false;
      this.newConversationTitle = '';
      this.conversationToRename = null;
    },
    
    async renameConversation() {
      if (!this.conversationToRename) return;
      
      try {
        const response = await apiClient.put(`${API_BASE_URL}/api/chat/sessions/${this.conversationToRename}`, {
          title: this.newConversationTitle
        });
        if (response.status !== 200) {
          throw new Error(this.$t('history.renameError'));
        }
        // 更新本地会话列表
        const index = this.conversations.findIndex(conv => conv.session_id === this.conversationToRename);
        if (index !== -1) {
          this.conversations[index].title = this.newConversationTitle;
        }
        
        // 如果重命名的是当前选中的会话，触发更新事件以更新页面标题
        if (this.selectedConversationId === this.conversationToRename) {
          this.$emit('conversation-updated');
        }
      } catch (error) {
        console.error(this.$t('history.renameErrorLog'), error);
      } finally {
        this.showRenameModal = false;
        this.newConversationTitle = '';
        this.conversationToRename = null;
      }
    },
    handleClickOutside(event) {
      // 如果有打开的操作菜单，且点击的不是菜单内的元素，则关闭菜单
      if (this.activeActionMenu && !event.target.closest('.conversation-actions')) {
        this.activeActionMenu = null;
      }
    },
  }
}
</script>

<style scoped>
.chat-history {
  width: 280px;
  height: 100%;
  background-color: var(--card-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  transition: background-color 0.3s, border-color 0.3s;
  position: relative;
}

.chat-history-header {
  padding: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
}

.new-chat-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background-color: var(--card-bg);
  color: var(--text-color);
  font-size: 14px;
  width: 100%;
  transition: all 0.2s;
  font-weight: 500;
}

.new-chat-btn:hover {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.icon-plus {
  display: inline-block;
  width: 16px;
  height: 16px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='currentColor' d='M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: center;
  background-size: contain;
}

.chat-history-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
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

.empty-history {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 20px;
  text-align: center;
}

.empty-icon {
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.empty-history p {
  color: var(--text-secondary);
  margin-bottom: 20px;
}

.start-chat-btn {
  padding: 8px 16px;
  border-radius: 18px;
  border: none;
  background-color: var(--primary-color);
  color: white;
  font-size: 14px;
}

.start-chat-btn:hover {
  background-color: var(--secondary-color);
}

.conversations-list {
  overflow-y: auto;
  flex-grow: 1;
  padding: 8px;
}

.conversation-item {
  position: relative;
  padding: 12px 16px;
  border-radius: 8px;
  margin: 0 8px 8px 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conversation-item.active {
  background-color: rgba(29, 155, 240, 0.1);
}

.conversation-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.conversation-info {
  flex: 1;
  min-width: 0; /* 确保flex子项可以缩小到小于其内容大小 */
}

.conversation-title {
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-color);
}

.conversation-date {
  font-size: 12px;
  color: var(--text-secondary);
}

.conversation-actions {
  display: flex;
  align-items: center;
  position: relative;
}

.action-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  padding: 0;
  cursor: pointer;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
  visibility: hidden;
  border-radius: 50%;
  transition: all 0.2s;
}

.conversation-item:hover .action-btn {
  visibility: visible;
}

.action-btn:hover {
  color: var(--primary-color);
  background-color: rgba(0, 0, 0, 0.06);
}

.action-menu {
  position: absolute;
  top: 0;
  right: 30px;
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
  width: 140px;
  text-align: left;
  border: 1px solid var(--border-color);
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 16px;
  border: none;
  background: none;
  color: var(--text-color);
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.2s;
}

.menu-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.menu-item.delete {
  color: var(--danger-color, #dc3545);
}

.rename-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: var(--card-bg);
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 300px;
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 16px;
  font-weight: 600;
}

.rename-input {
  width: 100%;
  padding: 10px;
  margin-bottom: 20px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background-color: var(--background-color);
  color: var(--text-color);
  font-size: 14px;
  transition: border-color 0.3s;
}

.rename-input:focus {
  outline: none;
  border-color: var(--primary-color);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-cancel, .btn-save {
  padding: 8px 16px;
  border-radius: 6px;
  border: none;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-cancel {
  background-color: var(--ai-message-bg, #f5f5f5);
  color: var(--text-color);
}

.btn-cancel:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.btn-save {
  background-color: var(--primary-color);
  color: white;
}

.btn-save:hover {
  background-color: var(--secondary-color, #0056b3);
}

.delete-confirm {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: var(--card-bg);
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
  width: 90%;
  max-width: 240px;
  text-align: center;
}

.delete-confirm p {
  margin-bottom: 16px;
  color: var(--text-color);
}

.delete-actions {
  display: flex;
  justify-content: space-between;
}

.btn-delete {
  background-color: #f4212e;
  color: white;
}
</style> 