import { defineStore } from 'pinia'
import { ref } from 'vue'
import { jwtDecode } from 'jwt-decode'
import { buildApiUrl } from '../config/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const username = ref<string>('')

  // 初始化时从token中解析用户信息
  if (token.value) {
    try {
      const decoded: any = jwtDecode(token.value)
      username.value = decoded.username || ''
      console.log('从token解析的用户名:', decoded.username)
    } catch (error) {
      console.error('Token解析失败:', error)
      token.value = null
      localStorage.removeItem('token')
    }
  }

  const login = async (usernameInput: string, password: string) => {
    const response = await fetch(buildApiUrl('/api/login'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: usernameInput,
        password: password
      })
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(errorText || '登录失败')
    }

    const data = await response.json()
    token.value = data.token
    username.value = usernameInput
    localStorage.setItem('token', data.token)
  }

  const register = async (usernameInput: string, password: string) => {
    const response = await fetch(buildApiUrl('/api/register'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: usernameInput,
        password: password
      })
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(errorText || '注册失败')
    }

    const data = await response.json()
    token.value = data.token
    username.value = usernameInput
    localStorage.setItem('token', data.token)
  }

  const logout = () => {
    token.value = null
    username.value = ''
    localStorage.removeItem('token')
  }

  const isAuthenticated = () => {
    return !!token.value
  }

  return {
    token,
    username,
    login,
    register,
    logout,
    isAuthenticated
  }
}) 