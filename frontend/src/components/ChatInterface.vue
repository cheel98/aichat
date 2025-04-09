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
            <div v-if="message.role === 'user'">
              <!-- 编辑模式下显示输入框，否则显示消息内容 -->
              <div v-if="editingMessageIndex === index" class="user-message-content edit-input-container">
                <textarea
                  v-model="editingContent"
                  @keydown="handleEditKeyDown($event, index, message.message_id)"
                  rows="1"
                  :ref="`editTextarea-${index}`"
                  class="edit-textarea"
                  @input="autoResizeEdit(index)"
                ></textarea>
                <div class="edit-actions">
                  <button @click="cancelEdit" class="edit-cancel-btn">取消</button>
                  <button @click="saveEdit(index, message.message_id)" class="edit-save-btn">保存并发送</button>
                </div>
              </div>
              <div v-else class="user-message-content">
                <span class="message-text">{{ message.content }}</span>
                <!-- 用户消息控制按钮 -->
                <div class="user-message-controls">
                  <button 
                    @click="editUserMessage(index, message.content, message.message_id)" 
                    class="edit-button" 
                    :disabled="loading"
                    title="修改问题"
                  >
                    <i class="bi bi-pencil"></i>
                  </button>
                  <button 
                    @click="regenerateResponse(message.message_id, index)" 
                    class="regenerate-button" 
                    :disabled="loading"
                    title="重新生成回答"
                  >
                    <i class="bi bi-arrow-repeat"></i>
                  </button>
                </div>
              </div>
            </div>
            <div v-else>
              <!-- 检查是否有思考内容 -->
              <template v-if="message.thinkingContent">
                <div class="thinking-container" :class="{ 'collapsed': !collapsedThinking[index] }">
                  <div class="thinking-header">
                    <i class="thinking-icon bi bi-lightbulb"></i>
                    <span v-if="index !== thinkingIndex">{{ $t('chat.thinkend') }}</span>
                    <span v-else>{{ $t('chat.thinking') }}</span>

                    <button @click="toggleThinking(index)" class="collapse-btn">
                      <i class="bi" :class="!collapsedThinking[index] ? 'bi-chevron-down' : 'bi-chevron-up'"></i>
                      {{ !collapsedThinking[index] ? $t('chat.showThinking') : $t('chat.hideThinking') }}
                    </button>
                  </div>
                  <div v-show="collapsedThinking[index]" v-html="renderMarkdown(message.thinkingContent)" class="markdown-content thinking-content"></div>
                </div>
                <div class="answer-container">
                  <div v-html="renderMarkdown(message.content)" class="markdown-content"></div>
                </div>
              </template>
              <template v-else>
                <div v-html="renderMarkdown(message.content)" class="markdown-content"></div>
              </template>
              
              <!-- AI回复的控制按钮 -->
              <div class="ai-response-controls">
                <!-- 只有当有多个回答时才显示版本选择 -->
                <div v-if="message.alternative_responses && message.alternative_responses.length > 0" class="response-versions">
                  <span class="versions-label">选择回答版本:</span>
                  <div class="version-buttons">
                    <!-- 原始回答 -->
                    <button 
                      @click="selectResponseVersion(message.message_id, 1)" 
                      :class="['version-btn', message.is_active ? 'active' : '']"
                      title="原始回答"
                    >
                      1
                    </button>
                    
                    <!-- 其他回答 -->
                    <button 
                      v-for="(response, rIndex) in message.alternative_responses" 
                      :key="'r-' + rIndex"
                      @click="selectResponseVersion(message.message_id, response.version)" 
                      :class="['version-btn', response.is_active ? 'active' : '']"
                      :title="`回答 ${response.version}`"
                    >
                      {{ response.version }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
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
          :class="{ 'active': deepThink }"
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
import { v4 as uuidv4 } from 'uuid'

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
      deepThink: false,
      collapsedThinking: {}, // 追踪每个消息的折叠状态
      thinkingIndex: -1, // 思考内容索引
      retryingMessageId: null, // 当前正在重试的消息ID
      pendingResponseVersion: null, // 等待设置为激活的响应版本
      editingMessageIndex: -1, // 当前正在编辑的消息索引
      editingMessageId: null, // 当前正在编辑的消息ID
      editingContent: '' // 正在编辑的内容
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
    logout() {
      this.messages = [];
      this.currentConversationId = null;
      this.inputMessage = '';
      this.collapsedThinking = {};
      this.thinkingIndex = -1;
    },
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
      
      // 正常发送新消息
      // 添加用户消息
      this.messages.push({
        role: 'user',
        content: this.inputMessage,
        message_id: uuidv4() // 生成客户端消息ID
      })
      
      this.loading = true

      // 添加一个空的AI回复消息，用于流式更新
      const aiMessageIndex = this.messages.length
      
      let msg = {
        content: this.inputMessage,
        thinking: this.deepThink
      }

      // 保存用户输入，然后清空输入框
      const userInput = this.inputMessage
      
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
        // 标记当前是否在思考模式
        if (this.deepThink){
          this.thinkingIndex = this.messages.length;
        }
        let currentContent = '';
        let thinkingContent = '';
        let messageId = '';
        let responseVersion = null;
        
        const response = await apiClient.post(apiUrl, {
          content: msg.content,
          thinking: this.deepThink
        }, {
          responseType: 'stream',
          onDownloadProgress: (progressEvent) => {
            const chunk = progressEvent.event.target.response;
            if (chunk) {
              this.loading = false;
              
              // 检查是否包含消息ID或版本标识
              if (chunk.includes('$messageId$')) {
                const parts = chunk.split('$messageId$');
                messageId = parts[1] || '';
                // 更新消息，将消息ID添加到内容中
                currentContent = parts[0];
              } else if (chunk.includes('$responseVersion$')) {
                const parts = chunk.split('$responseVersion$');
                responseVersion = parseInt(parts[1], 10) || null;
                // 只保留内容部分
                currentContent = parts[0];
              } else if (chunk.includes('$thinkEnd$') && this.thinkingIndex !== -1) {
                // 分离思考内容和回答内容
                const parts = chunk.split('$thinkEnd$');
                currentContent = parts[1] || '';
                this.thinkingIndex = -1;
              } else if (this.thinkingIndex !== -1) {
                // 仍在思考模式，所有内容都是思考内容
                thinkingContent = chunk;
                currentContent = '';
              } else if (chunk.includes('$thinkEnd$')){
                // 分离思考内容和回答内容
                const parts = chunk.split('$thinkEnd$');
                currentContent = parts[1] || '';
              } else {
                // 不在思考模式，所有内容都是回答内容
                currentContent = chunk;
              }
              
              // 更新AI回复消息
              this.messages[aiMessageIndex] = {
                role: 'ai',
                content: currentContent,
                thinkingContent: thinkingContent,
                message_id: messageId,
                is_active: true
              };
              
              // 滚动到底部
              this.scrollToBottom();
            }
          }
        });
        
        // 如果是重试消息的响应
        if (this.retryingMessageId && responseVersion) {
          this.pendingResponseVersion = {
            messageId: this.retryingMessageId,
            version: responseVersion
          };
          
          // 重置状态
          this.retryingMessageId = null;
          
          // 更新消息列表
          await this.loadConversation(this.currentConversationId);
        }
        
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
          const processedMessages = data.messages.map(message => {
            if (message.role === 'ai' && message.content) {
              return {
                ...message,
                thinkingContent: message.think_content,
                content: message.content
              };
            }
            return message;
          });
          this.messages = processedMessages;
          this.currentConversationId = conversationId;
          
          // 如果有待处理的版本设置请求，立即处理
          if (this.pendingResponseVersion) {
            await this.selectResponseVersion(
              this.pendingResponseVersion.messageId, 
              this.pendingResponseVersion.version
            );
            this.pendingResponseVersion = null;
          }
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
        this.$message.error(this.$t('chat.saveError') || '保存对话失败');
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
      this.deepThink = !this.deepThink;
    },

    // 切换思考内容的显示/隐藏
    toggleThinking(index) {
      this.collapsedThinking[index] = !this.collapsedThinking[index];
    },

    // 编辑用户消息
    editUserMessage(index, content, messageId) {
      this.editingMessageIndex = index;
      this.editingMessageId = messageId;
      this.editingContent = content;
      
      // 聚焦输入框
      this.$nextTick(() => {
        const editTextarea = this.$refs[`editTextarea-${index}`];
        if (editTextarea && editTextarea[0]) {
          editTextarea[0].focus();
          this.autoResizeEdit(index);
        }
      });
    },
    
    // 取消编辑
    cancelEdit() {
      this.editingMessageIndex = -1;
      this.editingMessageId = null;
      this.editingContent = '';
    },
    
    // 处理编辑框中的键盘事件
    handleEditKeyDown(event, index, messageId) {
      // Enter键发送消息，Shift+Enter添加换行符
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault();
        this.saveEdit(index, messageId);
      }
    },
    
    // 编辑框自动调整高度
    autoResizeEdit(index) {
      const editTextarea = this.$refs[`editTextarea-${index}`];
      if (!editTextarea || !editTextarea[0]) return;
      
      const textarea = editTextarea[0];
      textarea.style.height = 'auto';
      const newHeight = Math.min(textarea.scrollHeight, 120);
      textarea.style.height = newHeight + 'px';
    },
    
    // 保存编辑并发送
    async saveEdit(index, messageId) {
      if (!this.editingContent.trim() || this.loading) {
        return;
      }
      
      const originalContent = this.messages[index].content;
      
      // 更新用户消息内容
      this.messages[index].content = this.editingContent.trim();
      
      // 如果内容没有变化，则不重新生成回答
      if (originalContent === this.editingContent.trim()) {
        this.cancelEdit();
        return;
      }
      
      // 删除此消息后面的所有回答
      const userMessages = this.messages.filter(msg => msg.role === 'user');
      const userIndex = userMessages.findIndex(msg => msg.message_id === messageId);
      
      if (userIndex !== -1 && userIndex < userMessages.length - 1) {
        // 找到当前用户消息在所有消息中的索引
        const nextUserMsgIndex = this.messages.findIndex(
          (msg, i) => i > index && msg.role === 'user'
        );
        
        if (nextUserMsgIndex !== -1) {
          // 删除从当前消息的AI回复到下一条用户消息之前的所有消息
          this.messages.splice(index + 1, nextUserMsgIndex - (index + 1));
        } else {
          // 删除从当前消息的AI回复到结尾的所有消息
          this.messages.splice(index + 1);
        }
      } else {
        // 如果是最后一条用户消息，只保留自己，删除后面的AI回复
        if (index < this.messages.length - 1) {
          this.messages.splice(index + 1);
        }
      }
      
      // 准备编辑后的内容
      const editedContent = this.editingContent.trim();
      
      // 退出编辑模式
      this.cancelEdit();
      
      // 保存会话以更新用户消息内容
      await this.saveConversation();
      
      // 重新发送消息
      this.loading = true;
      
      try {
        // 准备API请求URL
        const apiUrl = `${API_BASE_URL}/api/chat/sessions/${this.currentConversationId}`;
        
        // 标记当前是否在思考模式
        if (this.deepThink) {
          this.thinkingIndex = index + 1;
        }
        
        let currentContent = '';
        let thinkingContent = '';
        let responseMessageId = '';
        
        const response = await apiClient.post(apiUrl, {
          content: editedContent,
          thinking: this.deepThink
        }, {
          responseType: 'stream',
          onDownloadProgress: (progressEvent) => {
            const chunk = progressEvent.event.target.response;
            if (chunk) {
              this.loading = false;
              
              // 检查是否包含消息ID或版本标识
              if (chunk.includes('$messageId$')) {
                const parts = chunk.split('$messageId$');
                responseMessageId = parts[1] || '';
                // 更新内容，去除ID部分
                currentContent = parts[0] || '';
              } else if (chunk.includes('$thinkEnd$') && this.thinkingIndex !== -1) {
                // 分离思考内容和回答内容
                const parts = chunk.split('$thinkEnd$');
                thinkingContent = thinkingContent || '';
                currentContent = parts[1] || '';
                this.thinkingIndex = -1;
              } else if (this.thinkingIndex !== -1) {
                // 仍在思考模式，所有内容都是思考内容
                thinkingContent += chunk;
              } else if (chunk.includes('$thinkEnd$')) {
                // 分离思考内容和回答内容
                const parts = chunk.split('$thinkEnd$');
                currentContent = parts[1] || '';
              } else {
                // 不在思考模式，所有内容都是回答内容
                currentContent += chunk;
              }
              
              // 如果是第一次接收响应，添加一个AI消息
              if (!this.messages[index + 1] || this.messages[index + 1].role !== 'ai') {
                this.messages.splice(index + 1, 0, {
                  role: 'ai',
                  content: currentContent,
                  thinkingContent: thinkingContent,
                  message_id: responseMessageId,
                  is_active: true
                });
              } else {
                // 更新现有AI消息
                this.messages[index + 1] = {
                  ...this.messages[index + 1],
                  content: currentContent,
                  thinkingContent: thinkingContent,
                  message_id: responseMessageId,
                  is_active: true
                };
              }
              
              // 滚动到底部
              this.scrollToBottom();
            }
          }
        });
        
        // 流式传输完成，保存对话
        await this.saveConversation();
      } catch (error) {
        console.error(this.$t('chat.streamError'), error);
        if (!this.messages[index + 1] || this.messages[index + 1].role !== 'ai') {
          this.messages.splice(index + 1, 0, {
            role: 'ai',
            content: this.$t('chat.errorMessage') || '请求失败，请重试',
            message_id: ''
          });
        }
        this.loading = false;
      }
    },
    
    // 重新生成回答
    regenerateResponse(messageId, userMessageIndex) {
      // 找到关联的AI消息
      const aiMessageIndex = userMessageIndex + 1;
      
      if (aiMessageIndex < this.messages.length && this.messages[aiMessageIndex].role === 'ai') {
        // 提取AI消息ID
        const aiMessageId = this.messages[aiMessageIndex].message_id;
        
        // 调用重试生成功能
        this.retryMessage(aiMessageId);
      } else {
        // 如果没有找到AI消息，尝试直接使用用户消息重新生成
        this.retryMessage(messageId);
      }
    },

    // 重试生成新回答
    async retryMessage(messageId) {
      if (this.loading) return;
      
      this.loading = true;
      this.retryingMessageId = messageId;
      
      try {
        await apiClient.post('/chat/retry', {
          message_id: messageId,
          thinking: this.deepThink
        }, {
          responseType: 'stream',
          onDownloadProgress: (progressEvent) => {
            const chunk = progressEvent.event.target.response;
            if (chunk) {
              this.loading = false;
              
              // 检查是否包含版本信息
              if (chunk.includes('$responseVersion$')) {
                const parts = chunk.split('$responseVersion$');
                const version = parseInt(parts[1], 10) || null;
                
                if (version) {
                  this.pendingResponseVersion = {
                    messageId: messageId,
                    version: version
                  };
                }
              }
            }
          }
        });
        
        // 重新加载会话以显示新回答
        await this.loadConversation(this.currentConversationId);
      } catch (error) {
        console.error('重试消息失败:', error);
        ElMessage.error(this.$t('chat.retryError') || '重试失败');
        this.loading = false;
        this.retryingMessageId = null;
      }
    },
    
    // 选择特定版本的回答
    async selectResponseVersion(messageId, version) {
      if (this.loading) return;
      
      try {
        const response = await apiClient.put('/chat/response/active', {
          message_id: messageId,
          version: version
        });
        
        if (response.status === 200) {
          // 重新加载会话以显示选中的回答
          await this.loadConversation(this.currentConversationId);
        } else {
          throw new Error('设置活跃版本失败');
        }
      } catch (error) {
        console.error('设置活跃版本失败:', error);
        ElMessage.error(this.$t('chat.setActiveError') || '切换版本失败');
      }
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
  width: auto;
  max-width: max-content;
  box-sizing: border-box;
}

.user-message > div {
  width: 100%;
  word-break: break-word;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.user-message > div > span {
  align-self: flex-start;
}

.user-message > div > .user-message-controls {
  width: 100%;
  justify-content: flex-end;
}

.user-message-content {
  display: flex;
  flex-direction: column;
  width: 100%;
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
  position: relative;
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

/* 思考容器样式 */
.thinking-container {
  margin-bottom: 16px;
  background-color: var(--thinking-bg, #f5f7fa);
  border-radius: 12px;
  padding: 16px;
  position: relative;
  border: 1px solid var(--thinking-border, #e1e4e8);
  transition: all 0.3s ease;
  max-height: 2000px;
  overflow: hidden;
}

.thinking-container.collapsed {
  padding: 12px 16px;
  max-height: 50px; /* 折叠时的高度 */
}

[data-theme="dark"] .thinking-container {
  background-color: var(--thinking-bg, rgba(55, 65, 81, 0.2));
  border-color: var(--thinking-border, rgba(75, 85, 99, 0.4));
}

.thinking-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 12px;
  font-size: 0.9rem;
}

.thinking-icon {
  color: var(--thinking-icon, #f59e0b);
}

.thinking-content {
  color: var(--thinking-text, #6b7280);
  font-size: 0.95rem;
  line-height: 1.5;
  transition: opacity 0.3s ease;
}

/* 最终回答容器 */
.answer-container {
  padding: 4px 0;
}

.thinking-container.collapsed .thinking-content {
  opacity: 0;
}

/* 添加折叠按钮样式 */
.collapse-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  font-size: 0.8rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0;
  margin-left: auto;
  transition: color 0.2s ease;
}

.collapse-btn:hover {
  color: var(--primary-color);
}

@media (max-width: 768px) {
  .message {
    max-width: 85%;
  }
  
  .user-message {
    max-width: max-content;
  }
  
  .edit-input-container {
    max-width: 100%;
  }
  
  .thinking-container {
    padding: 12px;
  }
  
  .thinking-header {
    font-size: 0.8rem;
  }
  
  .collapse-btn {
    font-size: 0.7rem;
  }
}

.ai-response-controls {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
  margin-top: 8px;
}

.retry-button {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.retry-button:hover:not(:disabled) {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--primary-color);
}

.retry-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.response-versions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.versions-label {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.version-buttons {
  display: flex;
  gap: 4px;
}

.version-btn {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  background-color: var(--chat-bg);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.version-btn:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.version-btn.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.user-message-content {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.message-text {
  text-align: left;
  align-self: flex-start;
  margin-bottom: 4px;
}

.user-message-controls {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  opacity: 0;
  transition: opacity 0.2s ease;
  width: 100%;
}

.user-message:hover .user-message-controls {
  opacity: 1;
}

.edit-button, .regenerate-button {
  background: transparent;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: rgba(255, 255, 255, 0.8);
}

.edit-button:hover, .regenerate-button:hover {
  background-color: rgba(255, 255, 255, 0.15);
  color: white;
}

.edit-button:disabled, .regenerate-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.edit-input-container {
  width: 100%;
  position: relative;
  max-width: 100%; /* 使用100%而不是inherit以确保与父元素宽度一致 */
  min-width: unset; /* 移除最小宽度限制，使其完全跟随气泡宽度 */
}

.edit-textarea {
  width: 100%;
  background: transparent;
  border: none;
  color: white;
  font-size: 1rem;
  line-height: 1.5;
  resize: none;
  padding: 4px 0;
  outline: none;
  font-family: var(--font-family);
  min-height: 24px;
  max-width: 100%;
  overflow-wrap: break-word;
  word-break: break-word;
  box-sizing: border-box;
  margin: 0; /* 移除默认边距 */
  display: block; /* 防止奇怪的内联块行为 */
}

.edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
  width: 100%; /* 确保按钮区域与文本框等宽 */
}

.edit-cancel-btn, .edit-save-btn {
  border: none;
  border-radius: 4px;
  padding: 5px 10px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.2s;
}

.edit-cancel-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.edit-save-btn {
  background: white;
  color: var(--user-message-bg);
  font-weight: 500;
}

.edit-cancel-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.edit-save-btn:hover {
  background: rgba(255, 255, 255, 0.9);
}
</style> 