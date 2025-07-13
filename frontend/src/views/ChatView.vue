<template>
  <div class="chat-container">
    <div class="chat-sidebar">
      <div class="user-info">
        <div class="user-avatar">
          <Icon name="user" :size="24" color="white" />
        </div>
        <div class="user-details">
          <h3>{{ currentUser }}</h3>
          <span class="user-status">在线</span>
        </div>
        <button @click="logout" class="logout-btn">
          <Icon name="logout" :size="16" color="white" />
          退出
        </button>
      </div>
      <div class="chat-list">
        <h4>聊天列表</h4>
        <div 
          class="chat-item" 
          :class="{ active: selectedChat && selectedChat.id === chat.id }"
          v-for="chat in chats" 
          :key="chat.id" 
          @click="selectChat(chat)"
        >
          <div class="chat-avatar">
            <Icon :name="chat.type === 'user' ? 'user' : 'group'" :size="20" color="var(--primary-color)" />
          </div>
          <div class="chat-info">
            <div class="chat-name">{{ chat.name }}</div>
            <div class="chat-type">
              <Icon :name="chat.type === 'user' ? 'message' : 'group'" :size="12" />
              {{ chat.type === 'user' ? '私聊' : '群聊' }}
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="chat-main">
      <div v-if="selectedChat" class="chat-header">
        <div class="chat-header-info">
          <div class="chat-header-avatar">
            <Icon :name="selectedChat.type === 'user' ? 'user' : 'group'" :size="24" color="var(--primary-color)" />
          </div>
          <div class="chat-header-details">
            <h3>{{ selectedChat.name }}</h3>
            <span class="chat-type">
              <Icon :name="selectedChat.type === 'user' ? 'message' : 'group'" :size="12" />
              {{ selectedChat.type === 'user' ? '私聊' : '群聊' }}
            </span>
          </div>
        </div>
      </div>
      
      <div class="messages-container" ref="messagesContainer">
        <div v-if="!selectedChat" class="empty-state">
          <Icon name="message" :size="64" color="var(--text-muted)" />
          <h3>选择一个聊天开始对话</h3>
          <p>从左侧聊天列表中选择一个用户或群组开始聊天</p>
        </div>
        
        <div v-else>
          <div v-for="message in messages" :key="message.id" class="message">
            <div class="message-bubble">
              <div class="message-sender">{{ message.sender }}</div>
              <div class="message-content">{{ message.content }}</div>
              <div class="message-time">{{ formatTime(message.timestamp) }}</div>
            </div>
          </div>
          
          <!-- 正在输入状态显示 -->
          <div v-if="typingUsers.length > 0" class="typing-indicator">
            <div class="typing-text">
              <Icon name="typing" :size="16" color="var(--text-secondary)" />
              {{ typingUsers.join(', ') }} 正在输入...
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="selectedChat" class="message-input">
        <input
          v-model="newMessage"
          @keyup.enter="sendMessage"
          @input="sendTyping"
          @blur="stopTyping"
          placeholder="输入消息..."
          class="message-input-field"
        />
        <button @click="sendMessage" class="send-btn">
          <Icon name="send" :size="18" color="white" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Icon from '../components/Icon.vue'
import { buildApiUrl, buildWsUrl } from '../config/api'

const router = useRouter()
const authStore = useAuthStore()

const currentUser = ref('')
const chats = ref<any[]>([])
const selectedChat = ref<any>(null)
const messages = ref<any[]>([])
const newMessage = ref('')
const messagesContainer = ref<HTMLElement>()
const typingUsers = ref<string[]>([])
const isTyping = ref(false)

let ws: WebSocket | null = null

const logout = () => {
  authStore.logout()
  router.push('/login')
}

const selectChat = (chat: any) => {
  selectedChat.value = chat
  messages.value = []
  typingUsers.value = [] // 清除正在输入状态
  stopTyping() // 停止当前的typing状态
  console.log('选择聊天:', chat)
  
  // 加载聊天历史
  if (ws && ws.readyState === WebSocket.OPEN) {
    if (chat.type === 'user') {
      // 加载私聊历史
      const historyRequest = {
        type: 'history',
        with: chat.id
      }
      console.log('请求私聊历史:', historyRequest)
      ws.send(JSON.stringify(historyRequest))
    } else {
      // 加载群聊历史
      const historyRequest = {
        type: 'history_group',
        group_id: parseInt(chat.id)
      }
      console.log('请求群聊历史:', historyRequest)
      ws.send(JSON.stringify(historyRequest))
    }
  }
}

const sendMessage = () => {
  if (!newMessage.value.trim() || !selectedChat.value || !ws) return

  let message: any = {}
  
  if (selectedChat.value.type === 'user') {
    // 私聊消息
    message = {
      type: 'send_message',
      to: selectedChat.value.id,
      content: newMessage.value
    }
  } else {
    // 群聊消息
    message = {
      type: 'send_group_message',
      group_id: parseInt(selectedChat.value.id),
      content: newMessage.value
    }
  }

  console.log('发送消息:', message)
  ws.send(JSON.stringify(message))
  newMessage.value = ''
  
  // 发送消息后停止typing状态
  stopTyping()
}

const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleTimeString()
}

let typingTimer: number | null = null

const sendTyping = () => {
  if (!selectedChat.value || !ws || ws.readyState !== WebSocket.OPEN) return
  
  if (!isTyping.value) {
    isTyping.value = true
    
    let typingMessage: any = {
      type: 'typing'
    }
    
    if (selectedChat.value.type === 'user') {
      typingMessage.to = selectedChat.value.id
    } else {
      typingMessage.group_id = parseInt(selectedChat.value.id)
    }
    
    console.log('发送正在输入状态:', typingMessage)
    ws.send(JSON.stringify(typingMessage))
  }
  
  // 清除之前的定时器
  if (typingTimer) {
    clearTimeout(typingTimer)
  }
  
  // 3秒后停止typing状态
  typingTimer = window.setTimeout(() => {
    isTyping.value = false
  }, 3000)
}

const stopTyping = () => {
  isTyping.value = false
  if (typingTimer) {
    clearTimeout(typingTimer)
    typingTimer = null
  }
}

const connectWebSocket = () => {
  const token = localStorage.getItem('token')
  if (!token) {
    console.error('没有找到token，无法连接WebSocket')
    return
  }

  console.log('正在连接WebSocket...')
  ws = new WebSocket(buildWsUrl(`/ws?token=${token}`))
  
  ws.onopen = () => {
    console.log('WebSocket 连接已建立')
  }
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      console.log('收到WebSocket消息:', data)
      
      if (data.type === 'error') {
        console.error('WebSocket错误:', data.msg)
        return
      }
      
      if (data.type === 'new_message' || data.type === 'new_group_message') {
        // 格式化消息显示
        const message = {
          id: Date.now(), // 临时ID
          sender: data.from,
          content: data.content,
          timestamp: data.ts || new Date().toISOString()
        }
        messages.value.push(message)
        
        // 滚动到底部
        setTimeout(() => {
          if (messagesContainer.value) {
            messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
          }
        }, 100)
      } else if (data.type === 'history' || data.type === 'history_group') {
        // 处理历史消息
        console.log('收到历史消息:', data)
        if (data.messages && Array.isArray(data.messages)) {
          messages.value = data.messages.map((msg: any) => ({
            id: msg.id || Date.now(),
            sender: msg.sender || msg.from,
            content: msg.content,
            timestamp: msg.timestamp || msg.ts || msg.created_at
          }))
          
          // 滚动到底部
          setTimeout(() => {
            if (messagesContainer.value) {
              messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
            }
          }, 100)
        }
      } else if (data.type === 'user_typing') {
        // 处理正在输入状态
        console.log('收到正在输入状态:', data)
        const typingUser = data.from
        
        if (typingUser && typingUser !== currentUser.value) {
          // 添加到正在输入列表
          if (!typingUsers.value.includes(typingUser)) {
            typingUsers.value.push(typingUser)
          }
          
          // 3秒后移除
          setTimeout(() => {
            const index = typingUsers.value.indexOf(typingUser)
            if (index > -1) {
              typingUsers.value.splice(index, 1)
            }
          }, 3000)
        }
      }
    } catch (error) {
      console.error('解析WebSocket消息失败:', error)
    }
  }
  
  ws.onclose = () => {
    console.log('WebSocket 连接已关闭')
  }
  
  ws.onerror = (error) => {
    console.error('WebSocket 错误:', error)
  }
}

const loadChats = async () => {
  try {
    const token = localStorage.getItem('token')
    console.log('正在加载聊天列表，token:', token ? '存在' : '不存在')
    
    const response = await fetch(buildApiUrl('/api/me/chats'), {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    console.log('聊天列表响应状态:', response.status)
    
    if (response.ok) {
      const data = await response.json()
      console.log('聊天数据:', data)
      
      // 将用户和群组合并为统一的聊天列表
      const allChats: any[] = []
      
      // 添加用户（排除自己）
      if (data.users) {
        data.users.forEach((user: any) => {
          if (user.username !== currentUser.value) {
            allChats.push({
              id: user.username,
              name: user.username,
              type: 'user'
            })
          }
        })
      }
      
      // 添加群组
      if (data.groups) {
        data.groups.forEach((group: any) => {
          allChats.push({
            id: group.id.toString(),
            name: group.name,
            type: 'group'
          })
        })
      }
      
      chats.value = allChats
      console.log('处理后的聊天列表:', allChats)
    }
  } catch (error) {
    console.error('加载聊天列表失败:', error)
  }
}

onMounted(() => {
  currentUser.value = authStore.username || 'Unknown'
  console.log('当前用户:', currentUser.value)
  console.log('认证状态:', authStore.isAuthenticated())
  connectWebSocket()
  loadChats()
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})
</script>

<style scoped>
.chat-container {
  display: flex;
  height: 100vh;
  background: var(--bg-secondary);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.chat-sidebar {
  width: 320px;
  background: var(--bg-primary);
  border-right: 1px solid var(--border-light);
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-sm);
  z-index: 10;
}

.user-info {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-light);
  display: flex;
  align-items: center;
  gap: 1rem;
  background: linear-gradient(135deg, var(--primary-color), var(--info-color));
  color: var(--text-white);
}

.user-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.user-details {
  flex: 1;
}

.user-details h3 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
}

.user-status {
  font-size: 0.8rem;
  opacity: 0.8;
}

.logout-btn {
  background: rgba(255, 255, 255, 0.2);
  color: var(--text-white);
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 0.5rem 1rem;
  border-radius: var(--border-radius);
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all var(--transition-fast);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
  transform: translateY(-1px);
}

.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 0;
}

.chat-list h4 {
  margin: 0 0 1rem 0;
  padding: 0 1rem;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.chat-item {
  padding: 1rem 1.5rem;
  margin: 0 0.5rem;
  border-radius: var(--border-radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
  position: relative;
  border: 1px solid transparent;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.chat-item:hover {
  background: var(--bg-secondary);
  border-color: var(--border-light);
  transform: translateX(4px);
}

.chat-item.active {
  background: var(--primary-light);
  border-color: var(--primary-color);
}

.chat-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 2px solid var(--border-light);
}

.chat-info {
  flex: 1;
}

.chat-name {
  font-weight: 600;
  margin-bottom: 0.25rem;
  color: var(--text-primary);
  font-size: 1rem;
}

.chat-type {
  font-size: 0.8rem;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-chat);
  position: relative;
}

.chat-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-light);
  background: var(--bg-primary);
  display: flex;
  align-items: center;
  box-shadow: var(--shadow-sm);
  z-index: 5;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.chat-header-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--border-light);
}

.chat-header-details h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
}

.chat-header-details .chat-type {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.8rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  background: linear-gradient(180deg, #f8f9fa 0%, #ffffff 100%);
  position: relative;
}

.message {
  margin-bottom: 1rem;
  max-width: 70%;
  position: relative;
  animation: messageSlideIn 0.3s ease-out;
}

.message:nth-child(odd) {
  margin-left: auto;
  margin-right: 0;
}

.message:nth-child(odd) .message-bubble {
  background: var(--bg-message-out);
  border-radius: 1.5rem 1.5rem 0.5rem 1.5rem;
}

.message:nth-child(even) .message-bubble {
  background: var(--bg-message-in);
  border-radius: 1.5rem 1.5rem 1.5rem 0.5rem;
}

.message-bubble {
  padding: 1rem 1.25rem;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  position: relative;
  backdrop-filter: blur(10px);
}

.message-sender {
  font-weight: 600;
  color: var(--primary-color);
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.message-content {
  margin-bottom: 0.5rem;
  line-height: 1.5;
  color: var(--text-primary);
  word-wrap: break-word;
}

.message-time {
  font-size: 0.75rem;
  color: var(--text-muted);
  text-align: right;
  margin-top: 0.25rem;
}

.message-input {
  padding: 1.5rem;
  border-top: 1px solid var(--border-light);
  background: var(--bg-primary);
  display: flex;
  gap: 1rem;
  align-items: flex-end;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}

.message-input-field {
  flex: 1;
  padding: 1rem 1.25rem;
  border: 1px solid var(--border-color);
  border-radius: 2rem;
  font-size: 1rem;
  background: var(--bg-secondary);
  resize: none;
  min-height: 50px;
  max-height: 120px;
  transition: all var(--transition-fast);
}

.message-input-field:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(0, 136, 204, 0.1);
  background: var(--bg-primary);
}

.send-btn {
  padding: 1rem 1.5rem;
  background: linear-gradient(135deg, var(--primary-color), var(--info-color));
  color: var(--text-white);
  border: none;
  border-radius: 2rem;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.875rem;
  transition: all var(--transition-fast);
  box-shadow: var(--shadow-sm);
  min-width: 80px;
}

.send-btn:hover {
  background: linear-gradient(135deg, var(--primary-hover), #0097a7);
  transform: translateY(-2px);
  box-shadow: var(--shadow);
}

.send-btn:active {
  transform: translateY(0);
}

.typing-indicator {
  padding: 0.75rem 1.25rem;
  margin-bottom: 1rem;
  max-width: 70%;
  animation: messageSlideIn 0.3s ease-out;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  padding: 2rem;
}

.empty-state h3 {
  margin: 1rem 0 0.5rem;
  color: var(--text-primary);
  font-size: 1.5rem;
  font-weight: 600;
}

.empty-state p {
  color: var(--text-secondary);
  font-size: 1rem;
  max-width: 300px;
  line-height: 1.6;
}

.typing-text {
  font-style: italic;
  color: var(--text-secondary);
  font-size: 0.875rem;
  background: var(--bg-tertiary);
  padding: 0.75rem 1rem;
  border-radius: 1.5rem 1.5rem 1.5rem 0.5rem;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

@keyframes messageSlideIn {
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
  .chat-sidebar {
    width: 280px;
  }
  
  .message {
    max-width: 85%;
  }
  
  .message-input {
    padding: 1rem;
    gap: 0.5rem;
  }
  
  .message-input-field {
    padding: 0.75rem 1rem;
  }
  
  .send-btn {
    padding: 0.75rem 1rem;
    min-width: 60px;
  }
}

@media (max-width: 480px) {
  .chat-sidebar {
    width: 100%;
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1000;
    transform: translateX(-100%);
    transition: transform var(--transition-normal);
  }
  
  .chat-sidebar.open {
    transform: translateX(0);
  }
  
  .chat-main {
    width: 100%;
  }
}
</style> 