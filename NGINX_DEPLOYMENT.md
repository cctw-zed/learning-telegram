# Nginx 反向代理部署指南

本指南将帮助您在 Learning Telegram 项目中实现 Nginx 反向代理，解决跨域问题。

## 🎯 为什么需要反向代理？

### 跨域问题
- **前端**: `http://localhost:5173` (Vue.js)
- **后端**: `http://localhost:8080` (Go)
- **问题**: 端口不同导致跨域限制

### 反向代理解决方案
- **统一入口**: `http://localhost:80`
- **透明转发**: nginx 自动转发请求到对应服务
- **消除跨域**: 浏览器只看到同域名请求

## 🚀 自动化部署

### Linux/macOS
```bash
# 运行自动化部署脚本
./deploy.sh
```

### Windows
```powershell
# 运行PowerShell脚本
.\deploy.ps1
```

## 📋 手动部署步骤

### 1. 安装 Nginx

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install nginx
```

#### CentOS/RHEL
```bash
sudo yum install nginx
```

#### macOS
```bash
brew install nginx
```

#### Windows
1. 访问 http://nginx.org/en/download.html
2. 下载 Windows 版本
3. 解压到 `C:\nginx`
4. 添加到系统 PATH

### 2. 配置 Nginx

#### Linux/macOS
```bash
# 复制配置文件
sudo cp nginx.conf /etc/nginx/sites-available/learning-telegram

# 创建软链接
sudo ln -s /etc/nginx/sites-available/learning-telegram /etc/nginx/sites-enabled/

# 删除默认配置（避免冲突）
sudo rm /etc/nginx/sites-enabled/default
```

#### Windows
```powershell
# 替换默认配置文件
# 将 nginx.conf 内容复制到 C:\nginx\conf\nginx.conf
```

### 3. 测试配置
```bash
# 测试配置语法
sudo nginx -t
```

### 4. 启动 Nginx
```bash
# 启动
sudo nginx

# 重启
sudo nginx -s reload

# 停止
sudo nginx -s stop
```

## 🔧 配置详解

### 核心配置结构
```nginx
server {
    listen 80;                    # 监听80端口
    server_name localhost;        # 服务器名称
    
    # 前端应用代理
    location / {
        proxy_pass http://localhost:5173;
        # 设置代理头部...
    }
    
    # API接口代理
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        # 设置代理头部...
    }
    
    # WebSocket代理
    location /ws {
        proxy_pass http://localhost:8080/ws;
        # WebSocket特殊配置...
    }
}
```

### 关键代理头部
```nginx
proxy_set_header Host $host;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
```

### WebSocket特殊配置
```nginx
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
proxy_read_timeout 3600;
proxy_send_timeout 3600;
```

## 🎭 前端代码适配

### 自动适配逻辑
项目已经实现了自动适配逻辑，位于 `frontend/src/config/api.ts`：

```typescript
// 检测是否在代理环境
const isProxyEnvironment = () => {
  return window.location.host !== 'localhost:5173'
}

// 自动配置API地址
const getApiConfig = (): ApiConfig => {
  if (isProxyEnvironment()) {
    // 代理环境：使用相对路径
    return {
      baseUrl: '',
      wsUrl: 'ws://' + window.location.host
    }
  } else {
    // 开发环境：直接访问后端
    return {
      baseUrl: 'http://localhost:8080',
      wsUrl: 'ws://localhost:8080'
    }
  }
}
```

### 请求示例
```typescript
// 自动适配的API请求
const response = await fetch(buildApiUrl('/api/login'), {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username, password })
})

// 自动适配的WebSocket连接
const ws = new WebSocket(buildWsUrl(`/ws?token=${token}`))
```

## 🔄 完整的启动流程

### 1. 启动后端服务
```bash
# 在项目根目录
go run cmd/server/main.go
```

### 2. 启动前端服务
```bash
# 在 frontend 目录
cd frontend
npm run dev
```

### 3. 启动 Nginx
```bash
# 启动 nginx
sudo nginx
```

### 4. 访问应用
```
http://localhost
```

## 🌊 请求流程对比

### 传统方式（有跨域问题）
```
浏览器 → 前端(5173) → 浏览器 → 后端(8080) ❌跨域
```

### 反向代理方式（无跨域问题）
```
浏览器 → Nginx(80) → 前端(5173) ✅同域
浏览器 → Nginx(80) → 后端(8080) ✅同域
```

## 🛠️ 故障排除

### 常见问题

#### 1. 端口占用
```bash
# 查看端口占用
lsof -i :80
netstat -tlnp | grep :80

# 修改nginx配置中的端口
listen 8080;  # 改为其他端口
```

#### 2. 权限问题
```bash
# 检查nginx用户权限
ps aux | grep nginx
```

#### 3. 配置错误
```bash
# 查看nginx错误日志
sudo tail -f /var/log/nginx/error.log
```

#### 4. 服务状态
```bash
# 检查服务状态
sudo systemctl status nginx
```

### Windows特有问题

#### 1. 路径问题
```powershell
# 确保nginx路径正确
$env:PATH += ";C:\nginx"
```

#### 2. 防火墙
```
确保防火墙允许80端口访问
```

## 📊 性能优化

### 缓存配置
```nginx
# 静态文件缓存
location ~* \.(js|css|png|jpg|jpeg|gif|svg|woff|woff2)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

### 压缩配置
```nginx
# 启用gzip压缩
gzip on;
gzip_types text/plain text/css application/json application/javascript;
```

### 负载均衡
```nginx
# 后端服务器池
upstream backend {
    server localhost:8080;
    server localhost:8081;  # 多个后端实例
}

location /api/ {
    proxy_pass http://backend/api/;
}
```

## 🚦 监控和日志

### 访问日志
```nginx
# 自定义日志格式
log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                '$status $body_bytes_sent "$http_referer" '
                '"$http_user_agent" "$http_x_forwarded_for"';

access_log /var/log/nginx/access.log main;
```

### 实时监控
```bash
# 实时查看访问日志
tail -f /var/log/nginx/access.log

# 实时查看错误日志
tail -f /var/log/nginx/error.log
```

## 🔒 安全配置

### HTTPS配置
```nginx
server {
    listen 443 ssl;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # SSL优化配置
    ssl_session_cache shared:SSL:1m;
    ssl_session_timeout 5m;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
}
```

### 安全头部
```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
```

## 🎉 总结

通过 Nginx 反向代理，您可以：

1. ✅ **解决跨域问题** - 统一入口，消除跨域限制
2. ✅ **简化前端代码** - 使用相对路径，自动适配
3. ✅ **提高性能** - 缓存、压缩、负载均衡
4. ✅ **增强安全性** - 隐藏后端架构，添加安全头部
5. ✅ **便于部署** - 统一域名，便于生产环境部署

现在您可以通过 `http://localhost` 访问完整的应用，享受无跨域问题的开发体验！ 