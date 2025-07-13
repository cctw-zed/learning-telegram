<template>
  <div class="test-container">
    <div class="test-header">
      <h1>🚀 测试页面</h1>
      <p class="subtitle">开发者调试工具</p>
    </div>
    
    <div class="test-cards">
      <div class="test-card">
        <h3>基础功能测试</h3>
        <p>如果您能看到这个页面，说明Vue应用正常工作！</p>
        <button @click="testClick" class="btn-primary">点击测试</button>
        <div v-if="clicked" class="success-message">
          ✅ 按钮点击成功！
        </div>
      </div>
      
      <div class="test-card">
        <h3>认证状态</h3>
        <div class="status-info">
          <div class="status-item">
            <span class="label">当前用户:</span>
            <span class="value">{{ authStore.username || '未登录' }}</span>
          </div>
          <div class="status-item">
            <span class="label">Token状态:</span>
            <span class="value" :class="{ 'status-active': authStore.token }">
              {{ authStore.token ? '✅ 已认证' : '❌ 未认证' }}
            </span>
          </div>
        </div>
        <button @click="testLogin" class="btn-secondary">测试登录</button>
      </div>
      
      <div class="test-card">
        <h3>API测试</h3>
        <div class="button-group">
          <button @click="testChats" class="btn-outline">测试聊天列表</button>
          <button @click="testWebSocket" class="btn-outline">测试WebSocket</button>
          <button @click="testSendMessage" class="btn-outline">测试发送消息</button>
          <button @click="testTyping" class="btn-outline">测试正在输入</button>
        </div>
      </div>
      
      <div v-if="chatData" class="test-card full-width">
        <h3>📊 聊天数据</h3>
        <div class="data-container">
          <pre class="data-display">{{ JSON.stringify(chatData, null, 2) }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { buildApiUrl, buildWsUrl } from '../config/api'

const clicked = ref(false)
const authStore = useAuthStore()
const chatData = ref(null)

const testClick = () => {
  clicked.value = true
}

const testLogin = async () => {
  try {
    await authStore.login('test', 'test123')
    console.log('登录成功')
  } catch (error) {
    console.error('登录失败:', error)
  }
}

const testChats = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(buildApiUrl('/api/me/chats'), {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (response.ok) {
      chatData.value = await response.json()
    } else {
      console.error('获取聊天列表失败:', response.status)
    }
  } catch (error) {
    console.error('请求失败:', error)
  }
}

let testWs: WebSocket | null = null

const testWebSocket = () => {
  const token = localStorage.getItem('token')
  if (!token) {
    console.error('没有token')
    return
  }
  
  testWs = new WebSocket(buildWsUrl(`/ws?token=${token}`))
  
  testWs.onopen = () => {
    console.log('测试WebSocket连接成功')
  }
  
  testWs.onmessage = (event) => {
    console.log('收到测试消息:', event.data)
  }
  
  testWs.onerror = (error) => {
    console.error('测试WebSocket错误:', error)
  }
  
  testWs.onclose = () => {
    console.log('测试WebSocket连接关闭')
  }
}

const testSendMessage = () => {
  if (!testWs || testWs.readyState !== WebSocket.OPEN) {
    console.error('WebSocket未连接')
    return
  }
  
  const message = {
    type: 'send_message',
    to: 'testuser',
    content: '这是一条测试消息'
  }
  
  console.log('发送测试消息:', message)
  testWs.send(JSON.stringify(message))
}

const testTyping = () => {
  if (!testWs || testWs.readyState !== WebSocket.OPEN) {
    console.error('WebSocket未连接')
    return
  }
  
  const typingMessage = {
    type: 'typing',
    to: 'testuser'
  }
  
  console.log('发送正在输入状态:', typingMessage)
  testWs.send(JSON.stringify(typingMessage))
}
</script>

<style scoped>
.test-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.test-header {
  text-align: center;
  margin-bottom: 3rem;
  color: white;
}

.test-header h1 {
  font-size: 3rem;
  margin-bottom: 0.5rem;
  font-weight: 700;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
  margin: 0;
}

.test-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.test-card {
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.test-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15);
}

.test-card.full-width {
  grid-column: 1 / -1;
}

.test-card h3 {
  color: var(--text-primary);
  margin-bottom: 1rem;
  font-size: 1.5rem;
  font-weight: 600;
}

.test-card p {
  color: var(--text-secondary);
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.status-info {
  margin-bottom: 1.5rem;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--border-light);
}

.status-item:last-child {
  border-bottom: none;
}

.label {
  font-weight: 500;
  color: var(--text-secondary);
}

.value {
  font-weight: 600;
  color: var(--text-primary);
}

.value.status-active {
  color: var(--success-color);
}

.button-group {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.75rem;
}

.success-message {
  margin-top: 1rem;
  padding: 1rem;
  background: var(--success-color);
  color: white;
  border-radius: 0.5rem;
  text-align: center;
  font-weight: 500;
  animation: slideIn 0.3s ease-out;
}

.data-container {
  background: var(--bg-secondary);
  border-radius: 0.5rem;
  padding: 1rem;
  border: 1px solid var(--border-light);
  max-height: 400px;
  overflow-y: auto;
}

.data-display {
  margin: 0;
  font-size: 0.875rem;
  color: var(--text-primary);
  line-height: 1.5;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

/* 按钮样式 */
button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 44px;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary-color), var(--info-color));
  color: white;
  box-shadow: 0 4px 12px rgba(0, 136, 204, 0.3);
}

.btn-primary:hover {
  background: linear-gradient(135deg, var(--primary-hover), #0097a7);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 136, 204, 0.4);
}

.btn-secondary {
  background: linear-gradient(135deg, var(--success-color), #20c997);
  color: white;
  box-shadow: 0 4px 12px rgba(40, 167, 69, 0.3);
}

.btn-secondary:hover {
  background: linear-gradient(135deg, #218838, #1e7e34);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(40, 167, 69, 0.4);
}

.btn-outline {
  background: transparent;
  color: var(--primary-color);
  border: 2px solid var(--primary-color);
}

.btn-outline:hover {
  background: var(--primary-color);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 136, 204, 0.3);
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .test-container {
    padding: 1rem;
  }
  
  .test-header h1 {
    font-size: 2rem;
  }
  
  .test-cards {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .test-card {
    padding: 1.5rem;
  }
  
  .button-group {
    grid-template-columns: 1fr;
  }
}
</style> 