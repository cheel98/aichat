<template>
  <div class="chat-container">
    <div class="chat-messages" ref="messageContainer">
      <div v-if="messages.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="64" height="64" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
        </div>
        <h2>{{ $t('chat.welcome') }}</h2>
        <p>{{ $t('chat.welcomeDesc') }}</p>
        <div class="example-questions">
          <button @click="useExampleQuestion($t('chat.exampleQuestions.0'))" class="example-btn">{{ $t('chat.exampleQuestions.0') }}</button>
          <button @click="useExampleQuestion($t('chat.exampleQuestions.1'))" class="example-btn">{{ $t('chat.exampleQuestions.1') }}</button>
          <button @click="useExampleQuestion($t('chat.exampleQuestions.2'))" class="example-btn">{{ $t('chat.exampleQuestions.2') }}</button>
        </div>
      </div>
      <div class="messages-list" v-if="messages.length > 0">
        <div 
          v-for="(message, index) in messages" 
          :key="index" 
          :class="['message-wrapper', message.role === 'user' ? 'user-wrapper' : 'ai-wrapper']"
        >
          <div class="message-avatar">
            <div v-if="message.role === 'user'" class="display-none">
              <!-- <span>{{ $t('chat.me') }}</span> -->
            </div>
            <div v-else class="ai-avatar">
              <span>AI</span>
            </div>
          </div>
          <div :class="['message', message.role === 'user' ? 'user-message' : 'ai-message']">
            <div v-if="message.role === 'user'">{{ message.content }}</div>
            <div v-else v-html="renderMarkdown(message.content)" class="markdown-content"></div>
          </div>
        </div>
        <div v-if="loading" class="message-wrapper ai-wrapper">
          <div class="message-avatar">
            <div class="ai-avatar">
              <span>AI</span>
            </div>
          </div>
          <div class="message ai-message loading-message">
            <div class="typing-indicator">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="chat-input-container">
      <div class="chat-input">
        <textarea 
          v-model="inputMessage" 
          @keydown="handleKeyDown"
          :placeholder="$t('chat.placeholder')"
          :disabled="loading"
          rows="1"
          ref="messageInput"
          @input="autoResize"
        ></textarea>
        <button 
          @click="sendMessage" 
          class="send-button"
          :disabled="loading || !inputMessage.trim()"
          :class="{ 'active': inputMessage.trim() }"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="22" y1="2" x2="11" y2="13"></line>
            <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
          </svg>
        </button>
      </div>
      <div class="input-options">
        <button 
          @click="toggleThinkMode" 
          class="think-button"
          :class="{ 'active': isDeepThinking }"
        >
          <span class="think-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2a8 8 0 0 0-8 8c0 2.2.9 4.2 2.3 5.6.4.4.6.9.6 1.4v1a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2v-1c0-.5.2-1 .6-1.4A7.95 7.95 0 0 0 20 10c0-4.4-3.6-8-8-8z"></path>
              <path d="M10 16v2"></path>
              <path d="M14 16v2"></path>
            </svg>
          </span>
          {{ $t('chat.think') }}
        </button>
      </div>
      <div class="input-info">
        {{ $t('chat.enterHint') }}
      </div>
    </div>
  </div>
</template>

<script>
import { marked } from 'marked'
import { API_BASE_URL } from '../config'
import apiClient from '../services/api'
import { ElMessage } from 'element-plus'

export default {
  name: 'ChatInterface',
  
  props: {
    conversationId: {
      type: String,
      default: null
    }
  },
  
  data() {
    return {
      inputMessage: '',
      messages: [],
      loading: false,
      currentConversationId: null,
      isDeepThinking: false
    }
  },
  
  watch: {
    conversationId: {
      immediate: true,
      handler(newId) {
        if (newId) {
          this.loadConversation(newId);
        }
      }
    },
    
    messages() {
      this.scrollToBottom();
    }
  },
  
  methods: {
    renderMarkdown(content) {
      try {
        return marked(content, { breaks: true, gfm: true })
      } catch (error) {
        console.error(this.$t('chat.markdownError'), error)
        return content
      }
    },
    
    handleKeyDown(event) {
      // Enter键发送消息，Shift+Enter添加换行符
      if (event.key === 'Enter') {
        if (event.shiftKey) {
          // Shift+Enter换行
          return
        } else {
          // Enter键发送消息
          event.preventDefault()
          this.sendMessage()
        }
      }
    },
    
    async sendMessage() {
      if (!this.inputMessage.trim() || this.loading) {
        return
      }
      
      // 添加用户消息
      this.messages.push({
        role: 'user',
        content: this.inputMessage
      })
      
      this.loading = true

      // 添加一个空的AI回复消息，用于流式更新
      const aiMessageIndex = this.messages.length
      
      let msg = {
        message: this.inputMessage,
        conversation_id: this.currentConversationId,
        deep_thinking: this.isDeepThinking
      }

      // 保存用户输入，然后清空输入框
      const userInput = this.inputMessage
      
      // // 根据界面选择的语言添加提示词告知AI以什么语言回复
      // if (this.$i18n.state.currentLanguage === 'zh-CN') {
      //   msg.message = "请用中文回复以下内容：" + msg.message;
      // } else if (this.$i18n.state.currentLanguage === 'en-US') {
      //   msg.message = "Please reply in English to the following: " + msg.message;
      // }
      this.inputMessage = ''
      this.scrollToBottom()
      this.$nextTick(() => {
        this.autoResize()
      })
      
      // 准备API请求URL
      let apiUrl;
      if (this.currentConversationId) {
        // 如果已有会话ID，使用现有会话
        apiUrl = `${API_BASE_URL}/api/chat/sessions/${this.currentConversationId}`;
      } else {
        // 如果是新会话，先创建会话
        const sessionRes = await apiClient.post('/chat/sessions', {
          title: userInput.length > 30 ? userInput.substring(0, 30) + '...' : userInput
        });
        if (sessionRes.status !== 201) {
          ElMessage.error(sessionRes.statusText);
          this.messages[aiMessageIndex] = {
            role: 'ai',
            content: sessionRes.statusText
          };
          this.loading = false;
          return;
        }
        console.log(sessionRes.data)
        
        const sessionData = sessionRes.data;
        this.currentConversationId = sessionData.session_id;
        apiUrl = `${API_BASE_URL}/api/chat/sessions/${this.currentConversationId}`;
        
        // 更新URL以包含会话ID
        this.updateUrlWithSessionId();
      }
      
      try {
        console.log(apiUrl)
        const response = await apiClient.post(apiUrl, {
          content: msg.message
        }, {
          responseType: 'stream',
          onDownloadProgress: (progressEvent) => {
            const chunk = progressEvent.event.target.response;
            if (chunk) {
              this.loading = false;
              // 更新AI回复消息
              this.messages[aiMessageIndex] = {
                role: 'ai',
                content: chunk
              };
              
              // 滚动到底部
              this.scrollToBottom();
            }
          }
        });
        // 流式传输完成，保存对话
        this.saveConversation();
      } catch (error) {
        console.error(this.$t('chat.streamError'), error);
        this.messages[aiMessageIndex] = {
          role: 'ai',
          content: this.$t('chat.errorMessage')
        };
        this.loading = false;
      }
    },
    
    scrollToBottom() {
      this.$nextTick(() => {
        if (this.$refs.messageContainer) {
          this.$refs.messageContainer.scrollTop = this.$refs.messageContainer.scrollHeight
        }
      })
    },

    useExampleQuestion(question) {
      this.inputMessage = question
      this.$nextTick(() => {
        this.autoResize()
      })
    },

    autoResize() {
      const textarea = this.$refs.messageInput
      if (!textarea) return
      
      textarea.style.height = 'auto'
      const newHeight = Math.min(textarea.scrollHeight, 120)
      textarea.style.height = newHeight + 'px'
    },
    
    // 加载特定对话的消息历史
    async loadConversation(conversationId) {
      this.clearMessages();
      this.loading = true;
      
      try {
        const response = await apiClient.get(`/chat/sessions/${conversationId}`);
        if (response.status !== 200) {
          throw new Error(this.$t('chat.loadHistoryError'));
        }
        
        const data = response.data;
        if (data.messages && Array.isArray(data.messages)) {
          this.messages = data.messages;
          this.currentConversationId = conversationId;
        }
      } catch (error) {
        console.error(this.$t('chat.loadHistoryErrorLog'), error);
      } finally {
        this.loading = false;
        this.scrollToBottom();
      }
    },
    
    // 清空消息
    clearMessages() {
      this.messages = [];
      this.currentConversationId = null;
    },
    
    // 保存对话
    async saveConversation() {
      // 如果消息为空，不保存
      if (this.messages.length === 0) return;
      
      try {
        // 获取第一条消息作为对话标题
        const firstUserMessage = this.messages.find(m => m.role === 'user');
        const title = firstUserMessage ? 
          (firstUserMessage.content.length > 30 ? 
            firstUserMessage.content.substring(0, 30) + '...' : 
            firstUserMessage.content) : 
          '未命名对话';
        
        const conversationData = {
          id: this.currentConversationId,
          title: title,
          messages: this.messages
        };
        
        // 根据是否有ID决定创建新对话还是更新现有对话
        const method = this.currentConversationId ? 'PUT' : 'POST';
        const url = this.currentConversationId ? 
          `${API_BASE_URL}/api/chat/sessions/${this.currentConversationId}` : 
          `${API_BASE_URL}/api/chat/sessions`;
        let response;
        if (method === 'PUT') {
          response = await apiClient.put(url, conversationData);
        } else {
          response = await apiClient.post(url, conversationData);
        }
        
        if (response.status !== 200) {
          throw new Error('保存对话失败');
        }
        
        // 如果是新创建的对话，获取返回的ID
        if (!this.currentConversationId) {
          const data = await response.json();
          this.currentConversationId = data.session_id;
          
          // 更新URL以包含会话ID
          this.updateUrlWithSessionId();
        }
        
        // 通知父组件对话已更新
        this.$emit('conversation-updated');
        
      } catch (error) {
        console.error('保存对话出错:', error);
        ElMessage.error(this.$t('chat.saveError') || '保存对话失败');
      }
    },
    
    // 更新URL以包含会话ID
    updateUrlWithSessionId() {
      if (this.$router && this.currentConversationId) {
        const currentRoute = this.$router.currentRoute.value;
        if (currentRoute.path === '/chat' || currentRoute.path === '/') {
          this.$router.replace({
            path: '/',
            query: { id: this.currentConversationId }
          });
        }
      }
    },
    
    toggleThinkMode() {
      this.isDeepThinking = !this.isDeepThinking;
    }
  },

  mounted() {
    this.$nextTick(() => {
      this.autoResize()
    })
  }
}
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  flex-grow: 1;
  background-color: var(--chat-bg);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

[data-theme="light"] .chat-container {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.chat-messages {
  flex-grow: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  height: 100%;
  padding: 0 20px;
}

.empty-icon {
  color: var(--text-secondary);
  margin-bottom: 20px;
  transition: color 0.3s;
}

.empty-state h2 {
  font-size: 1.8rem;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--text-color);
  transition: color 0.3s;
}

.empty-state p {
  color: var(--text-secondary);
  font-size: 1.1rem;
  margin-bottom: 24px;
  max-width: 500px;
  transition: color 0.3s;
}

.example-questions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  max-width: 500px;
}

.example-btn {
  background-color: transparent;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-color);
  padding: 12px 16px;
  text-align: left;
  font-size: 0.9rem;
  transition: all 0.2s ease, color 0.3s, border-color 0.3s;
}

.example-btn:hover {
  background-color: rgba(29, 155, 240, 0.1);
  border-color: var(--primary-color);
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
}

.message-wrapper {
  display: flex;
  gap: 12px;
  width: 100%;
}

.user-wrapper {
  justify-content: flex-end;
}

.ai-wrapper {
  justify-content: flex-start;
}

.message-avatar {
  display: flex;
  align-items: flex-start;
  flex-shrink: 0;
}

.user-avatar, .ai-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  transition: background-color 0.3s, color 0.3s;
}

.user-avatar {
  background-color: var(--user-message-bg);
  color: white;
}

.ai-avatar {
  background-color: var(--ai-message-bg);
  color: var(--text-color);
}

.message {
  padding: 14px 18px;
  border-radius: 18px;
  max-width: 80%;
  word-break: break-word;
  transition: background-color 0.3s, color 0.3s;
}

.user-message {
  background-color: var(--user-message-bg);
  color: white;
  border-top-right-radius: 4px;
}

.ai-message {
  background-color: var(--ai-message-bg);
  color: var(--text-color);
  border-top-left-radius: 4px;
}

.loading-message {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 14px 24px;
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background-color: var(--text-secondary);
  border-radius: 50%;
  display: inline-block;
  animation: typing 1.4s infinite ease-in-out;
  transition: background-color 0.3s;
}

.typing-indicator span:nth-child(1) {
  animation-delay: 0s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.6;
  }
  30% {
    transform: translateY(-6px);
    opacity: 1;
  }
}

.chat-input-container {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  background-color: var(--chat-bg);
  transition: background-color 0.3s, border-color 0.3s;
}

.chat-input {
  display: flex;
  align-items: flex-end;
  background-color: var(--ai-message-bg);
  border-radius: 24px;
  padding: 10px 16px;
  transition: background-color 0.3s, box-shadow 0.3s;
}

.chat-input:focus-within {
  box-shadow: 0 0 0 2px var(--primary-color);
}

textarea {
  flex-grow: 1;
  background: transparent;
  border: none;
  color: var(--text-color);
  font-size: 1rem;
  line-height: 1.5;
  max-height: 120px;
  resize: none;
  padding: 4px 0;
  outline: none;
  font-family: var(--font-family);
  transition: color 0.3s;
}

textarea::placeholder {
  color: var(--text-secondary);
  transition: color 0.3s;
}

.send-button {
  background-color: transparent;
  border: none;
  color: var(--text-secondary);
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 8px;
  padding: 0;
  cursor: pointer;
  transition: all 0.2s, color 0.3s, background-color 0.3s;
}

.send-button.active {
  color: var(--primary-color);
  background-color: rgba(29, 155, 240, 0.1);
}

.send-button.active:hover {
  background-color: rgba(29, 155, 240, 0.2);
}

.send-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.input-options {
  display: flex;
  justify-content: flex-start;
  padding: 4px 16px;
  margin-top: 4px;
}

.think-button {
  display: flex;
  align-items: center;
  gap: 6px;
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 20px;
  color: var(--text-secondary);
  padding: 6px 14px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.think-button:hover {
  background-color: rgba(29, 155, 240, 0.1);
}

.think-button.active {
  color: var(--primary-color);
  background-color: rgba(29, 155, 240, 0.15);
  border-color: var(--primary-color);
}

.think-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.think-icon svg {
  stroke: currentColor;
}

.input-info {
  text-align: center;
  margin-top: 8px;
  font-size: 0.8rem;
  color: var(--text-secondary);
  transition: color 0.3s;
}

kbd {
  background-color: var(--ai-message-bg);
  border-radius: 4px;
  padding: 2px 5px;
  font-size: 0.8rem;
  font-family: monospace;
  transition: background-color 0.3s;
}

.markdown-content {
  line-height: 1.6;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin-top: 1em;
  margin-bottom: 0.5em;
  font-weight: 600;
}

.markdown-content :deep(p) {
  margin-bottom: 0.75em;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  padding-left: 1.5em;
  margin-bottom: 0.75em;
}

.markdown-content :deep(code) {
  font-family: Consolas, Monaco, 'Andale Mono', monospace;
  background-color: rgba(0, 0, 0, 0.1);
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-size: 0.9em;
  white-space: pre-wrap;
}

.markdown-content :deep(pre) {
  background-color: rgba(0, 0, 0, 0.1);
  padding: 1em;
  border-radius: 5px;
  overflow-x: auto;
  margin: 0.75em 0;
}

.markdown-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
  white-space: pre;
  word-break: normal;
}

.markdown-content :deep(blockquote) {
  border-left: 3px solid var(--border-color);
  padding-left: 1em;
  margin-left: 0;
  margin-right: 0;
  color: var(--text-secondary);
}

.markdown-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 1em;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid var(--border-color);
  padding: 8px;
  text-align: left;
}

.markdown-content :deep(th) {
  background-color: rgba(0, 0, 0, 0.05);
  font-weight: 600;
}

.markdown-content :deep(a) {
  color: var(--primary-color);
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(img) {
  max-width: 100%;
  border-radius: 5px;
}
</style> 