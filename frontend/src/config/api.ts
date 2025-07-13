// API配置文件
// 支持开发环境和生产环境（nginx代理）的不同配置

interface ApiConfig {
  baseUrl: string
  wsUrl: string
}

// 检测是否在nginx反向代理环境中
const isProxyEnvironment = () => {
  // 如果当前域名不是localhost:5173，说明可能在代理环境中
  return window.location.host !== 'localhost:5173' && window.location.port !== '5173'
}

// 获取API配置
const getApiConfig = (): ApiConfig => {
  if (isProxyEnvironment()) {
    // 使用nginx反向代理时，使用相对路径
    return {
      baseUrl: '', // 空字符串表示使用当前域名
      wsUrl: window.location.protocol === 'https:' ? 'wss://' : 'ws://' + window.location.host
    }
  } else {
    // 开发环境，直接访问后端服务
    return {
      baseUrl: 'http://localhost:8080',
      wsUrl: 'ws://localhost:8080'
    }
  }
}

export const apiConfig = getApiConfig()

// 构建完整的API URL
export const buildApiUrl = (path: string): string => {
  const base = apiConfig.baseUrl
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return base + cleanPath
}

// 构建WebSocket URL
export const buildWsUrl = (path: string): string => {
  const base = apiConfig.wsUrl
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return base + cleanPath
}

// 调试信息
console.log('API配置:', {
  host: window.location.host,
  protocol: window.location.protocol,
  isProxy: isProxyEnvironment(),
  config: apiConfig
}) 