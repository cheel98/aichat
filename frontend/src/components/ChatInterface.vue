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
        <h2>欢迎使用 DeepSeek AI 聊天</h2>
        <p>开始与AI对话，DeepSeek将为您提供智能、有趣的回答</p>
        <div class="example-questions">
          <button @click="useExampleQuestion('介绍一下量子计算的基本原理')" class="example-btn">介绍一下量子计算的基本原理</button>
          <button @click="useExampleQuestion('帮我写一篇关于人工智能的短文')" class="example-btn">帮我写一篇关于人工智能的短文</button>
          <button @click="useExampleQuestion('如何提高编程效率？')" class="example-btn">如何提高编程效率？</button>
        </div>
      </div>
      <div class="messages-list" v-if="messages.length > 0">
        <div 
          v-for="(message, index) in messages" 
          :key="index" 
          :class="['message-wrapper', message.role === 'user' ? 'user-wrapper' : 'ai-wrapper']"
        >
          <div class="message-avatar">
            <div v-if="message.role === 'user'" class="user-avatar">
              <span>我</span>
            </div>
            <div v-else class="ai-avatar">
              <span>AI</span>
            </div>
          </div>
          <div :class="['message', message.role === 'user' ? 'user-message' : 'ai-message']">
            <p>{{ message.content }}</p>
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
          @keyup.enter.ctrl="sendMessage"
          @keyup.enter.meta="sendMessage"
          placeholder="输入消息，按 Ctrl+Enter 发送"
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
      <div class="input-info">
        按 <kbd>Ctrl</kbd>+<kbd>Enter</kbd> 发送
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'ChatInterface',
  
  data() {
    return {
      inputMessage: '',
      messages: [],
      loading: false,
      exampleQuestions: [
        '介绍一下量子计算的基本原理',
        '帮我写一篇关于人工智能的短文',
        '如何提高编程效率？'
      ]
    }
  },
  
  methods: {
    sendMessage() {
      if (!this.inputMessage.trim() || this.loading) {
        return
      }
      
      // 添加用户消息
      this.messages.push({
        role: 'user',
        content: this.inputMessage
      })
      
      this.loading = true
      
      // 发送消息到后端
      axios.post('http://localhost:8080/api/chat', {
        message: this.inputMessage
      })
        .then(response => {
          // 添加AI回复
          this.messages.push({
            role: 'ai',
            content: response.data.reply
          })
        })
        .catch(error => {
          console.error('发送消息出错:', error)
          this.messages.push({
            role: 'ai',
            content: '抱歉，发生了错误，请重试。'
          })
        })
        .finally(() => {
          this.loading = false
          this.inputMessage = ''
          this.scrollToBottom()
          this.$nextTick(() => {
            this.autoResize()
          })
        })
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
    }
  },
  
  watch: {
    messages() {
      this.scrollToBottom()
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
</style> 