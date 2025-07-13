<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <Icon name="message" :size="48" color="white" />
        </div>
        <h1 class="app-title">Telegram Chat</h1>
        <p class="app-subtitle">开始您的聊天之旅</p>
      </div>
      
      <div class="login-form">
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="username" class="form-label">用户名</label>
            <input
              id="username"
              v-model="username"
              type="text"
              placeholder="请输入用户名"
              required
              class="form-input"
              :class="{ 'is-invalid': error }"
            />
          </div>
          
          <div class="form-group">
            <label for="password" class="form-label">密码</label>
            <input
              id="password"
              v-model="password"
              type="password"
              placeholder="请输入密码"
              required
              class="form-input"
              :class="{ 'is-invalid': error }"
            />
          </div>
          
          <button type="submit" class="submit-btn" :disabled="loading">
            <Icon v-if="loading" name="loading" :size="20" color="white" />
            {{ loading ? '处理中...' : (isLogin ? '登录' : '注册') }}
          </button>
        </form>
        
        <div class="form-footer">
          <p class="switch-mode">
            {{ isLogin ? '还没有账户？' : '已有账户？' }}
            <a href="#" @click.prevent="toggleMode" class="switch-link">
              {{ isLogin ? '立即注册' : '立即登录' }}
            </a>
          </p>
        </div>
        
        <div v-if="error" class="error-message">
          <Icon name="error" :size="16" color="currentColor" />
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Icon from '../components/Icon.vue'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const isLogin = ref(true)
const loading = ref(false)
const error = ref('')

const toggleMode = () => {
  isLogin.value = !isLogin.value
  error.value = ''
}

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    if (isLogin.value) {
      await authStore.login(username.value, password.value)
    } else {
      await authStore.register(username.value, password.value)
    }
    router.push('/chat')
  } catch (err: any) {
    error.value = err.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="50" cy="50" r="0.5" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
  pointer-events: none;
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 500px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

.login-header {
  text-align: center;
  padding: 2rem 2rem 1rem;
  background: linear-gradient(135deg, #007bff, #00bcd4);
  color: white;
  position: relative;
  z-index: 1;
}

.login-header .logo {
  margin-bottom: 0.5rem;
}

.login-header .logo svg {
  fill: white;
}

.login-header .app-title {
  font-size: 2rem;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.login-header .app-subtitle {
  font-size: 1rem;
  opacity: 0.8;
}

.login-form {
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  position: relative;
  z-index: 2;
}

.form-group {
  position: relative;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  color: #555;
}

.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  box-sizing: border-box;
  transition: all 0.3s ease;
}

.form-input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.2);
}

.form-input.is-invalid {
  border-color: #dc3545;
  box-shadow: 0 0 0 3px rgba(220, 53, 69, 0.2);
}

.submit-btn {
  width: 100%;
  padding: 0.75rem;
  background: linear-gradient(135deg, #007bff, #00bcd4);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  min-height: 50px;
}

.submit-btn:hover {
  background: linear-gradient(135deg, #0056b3, #0097a7);
  transform: translateY(-2px);
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  color: #888;
}



.form-footer {
  text-align: center;
  margin-top: 1.5rem;
}

.switch-mode {
  font-size: 0.9rem;
  color: #555;
}

.switch-mode .switch-link {
  color: #007bff;
  text-decoration: none;
  font-weight: bold;
}

.switch-mode .switch-link:hover {
  text-decoration: underline;
}

.error-message {
  color: #dc3545;
  text-align: center;
  margin-top: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}


</style> 