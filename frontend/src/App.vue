<template>
  <router-view />
</template>

<style>
:root {
  /* 暗色主题变量 */
  --dark-primary-color: #1d9bf0;
  --dark-secondary-color: #0f7ae5;
  --dark-background-color: #0f1419;
  --dark-chat-bg: #16181c;
  --dark-user-message-bg: #1d9bf0;
  --dark-ai-message-bg: #2c2c2c;
  --dark-border-color: #38444d;
  --dark-text-color: #e7e9ea;
  --dark-text-secondary: #8b98a5;
  --dark-card-bg: #16181c;

  /* 亮色主题变量 */
  --light-primary-color: #1d9bf0;
  --light-secondary-color: #0f7ae5;
  --light-background-color: #f7f9fa;
  --light-chat-bg: #ffffff;
  --light-user-message-bg: #1d9bf0;
  --light-ai-message-bg: #f0f2f5;
  --light-border-color: #e1e8ed;
  --light-text-color: #14171a;
  --light-text-secondary: #657786;
  --light-card-bg: #ffffff;

  /* 默认使用亮色主题 */
  --primary-color: var(--light-primary-color);
  --secondary-color: var(--light-secondary-color);
  --background-color: var(--light-background-color);
  --chat-bg: var(--light-chat-bg);
  --user-message-bg: var(--light-user-message-bg);
  --ai-message-bg: var(--light-ai-message-bg);
  --border-color: var(--light-border-color);
  --text-color: var(--light-text-color);
  --text-secondary: var(--light-text-secondary);
  --card-bg: var(--light-card-bg);
  
  --font-family: 'PingFang SC', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}

/* 暗色主题样式 */
[data-theme="dark"] {
  --primary-color: var(--dark-primary-color);
  --secondary-color: var(--dark-secondary-color);
  --background-color: var(--dark-background-color);
  --chat-bg: var(--dark-chat-bg);
  --user-message-bg: var(--dark-user-message-bg);
  --ai-message-bg: var(--dark-ai-message-bg);
  --border-color: var(--dark-border-color);
  --text-color: var(--dark-text-color);
  --text-secondary: var(--dark-text-secondary);
  --card-bg: var(--dark-card-bg);
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  margin: 0;
  padding: 0;
  font-family: var(--font-family);
  background-color: var(--background-color);
  color: var(--text-color);
  width: 100%;
  height: 100vh;
  font-size: 14px;
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overflow: hidden;
}

html, body, #app {
  height: 100%;
  width: 100%;
}

#app {
  display: flex;
  flex-direction: column;
}

.app-container {
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--background-color);
  color: var(--text-color);
}

button {
  cursor: pointer;
}

/* 通用样式 */
.btn-primary {
  width: 100%;
  height: 40px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn-primary:hover {
  background-color: var(--secondary-color);
}

.btn-primary:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  color: #e74c3c;
  font-size: 13px;
  margin-bottom: 15px;
  text-align: center;
}

.field-error {
  color: #e74c3c;
  font-size: 12px;
  margin-top: 5px;
}
</style> 