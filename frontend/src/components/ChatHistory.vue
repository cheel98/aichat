<template>
  <div class="chat-history">
    <div class="chat-history-header">
      <h2>{{ $t('history.title') }}</h2>
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
        <div class="conversation-title">
          {{ conversation.title || $t('history.unnamed') }}
        </div>
        <div class="conversation-date">
          {{ formatDate(conversation.updated_at) }}
        </div>
        <button class="delete-btn" @click.stop="confirmDelete(conversation.session_id)">
          <span class="icon-delete"></span>
        </button>
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
      conversationToDelete: null
    }
  },
  
  created() {
    this.fetchConversations();
  },
  
  methods: {
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
      this.selectedConversationId = id;
      this.$emit('select-conversation', id);
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
    }
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
}

.chat-history-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
}

.chat-history-header h2 {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-color);
  margin: 0;
}

.new-chat-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 18px;
  border: none;
  background-color: var(--primary-color);
  color: white;
  font-size: 14px;
}

.new-chat-btn:hover {
  background-color: var(--secondary-color);
}

.icon-plus {
  display: inline-block;
  width: 14px;
  height: 14px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='white' d='M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6z'/%3E%3C/svg%3E");
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
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 4px;
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
  color: var(--text-color);
  position: relative;
}

.conversation-item:hover {
  background-color: rgba(29, 155, 240, 0.1);
}

.conversation-item.active {
  background-color: rgba(29, 155, 240, 0.2);
}

.conversation-title {
  font-size: 14px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding-right: 24px;
}

.conversation-date {
  font-size: 12px;
  color: var(--text-secondary);
}

.delete-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  background: transparent;
  border: none;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
  border-radius: 4px;
}

.conversation-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background-color: rgba(244, 33, 46, 0.1);
}

.icon-delete {
  display: inline-block;
  width: 16px;
  height: 16px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23f4212e' d='M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: center;
  background-size: contain;
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

.btn-cancel, .btn-delete {
  padding: 8px 16px;
  border-radius: 4px;
  border: none;
  font-size: 14px;
}

.btn-cancel {
  background-color: var(--ai-message-bg);
  color: var(--text-color);
}

.btn-delete {
  background-color: #f4212e;
  color: white;
}
</style> 