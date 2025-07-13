# Learning Telegram - Windows Nginx反向代理部署脚本
# 使用方法: .\deploy.ps1

Write-Host "🚀 开始部署Learning Telegram项目..." -ForegroundColor Green

# 检查是否安装了nginx
$nginxPath = Get-Command nginx -ErrorAction SilentlyContinue
if (-not $nginxPath) {
    Write-Host "❌ 未找到nginx" -ForegroundColor Red
    Write-Host "请手动安装nginx:" -ForegroundColor Yellow
    Write-Host "1. 访问 http://nginx.org/en/download.html"
    Write-Host "2. 下载Windows版本"
    Write-Host "3. 解压到 C:\nginx"
    Write-Host "4. 将 C:\nginx 添加到系统PATH"
    Write-Host "5. 重新运行此脚本"
    exit 1
}

Write-Host "✅ nginx已安装" -ForegroundColor Green

# 获取nginx安装路径
$nginxDir = Split-Path $nginxPath.Source
Write-Host "nginx安装路径: $nginxDir" -ForegroundColor Cyan

# 停止现有的nginx进程
Write-Host "🔄 停止现有的nginx服务..." -ForegroundColor Yellow
try {
    & nginx -s stop
    Start-Sleep -Seconds 2
} catch {
    Write-Host "没有运行中的nginx进程" -ForegroundColor Gray
}

# 备份原有nginx配置
$nginxConfPath = Join-Path $nginxDir "conf\nginx.conf"
if (Test-Path $nginxConfPath) {
    $backupPath = Join-Path $nginxDir "conf\nginx.conf.backup"
    Copy-Item $nginxConfPath $backupPath -Force
    Write-Host "✅ 已备份原有nginx配置" -ForegroundColor Green
}

# 创建新的nginx配置
Write-Host "📝 配置nginx反向代理..." -ForegroundColor Yellow
$newConfig = @"
# Learning Telegram 项目的 Nginx 反向代理配置
# 用于解决跨域问题

events {
    worker_connections 1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;

    server {
        listen 80;
        server_name localhost;
        
        # 前端Vue.js应用
        location / {
            proxy_pass http://localhost:5173;
            proxy_set_header Host `$host;
            proxy_set_header X-Real-IP `$remote_addr;
            proxy_set_header X-Forwarded-For `$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto `$scheme;
        }
        
        # API接口代理到Go后端
        location /api/ {
            proxy_pass http://localhost:8080/api/;
            proxy_set_header Host `$host;
            proxy_set_header X-Real-IP `$remote_addr;
            proxy_set_header X-Forwarded-For `$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto `$scheme;
            
            # 支持WebSocket升级
            proxy_http_version 1.1;
            proxy_set_header Upgrade `$http_upgrade;
            proxy_set_header Connection "upgrade";
        }
        
        # WebSocket连接
        location /ws {
            proxy_pass http://localhost:8080/ws;
            proxy_http_version 1.1;
            proxy_set_header Upgrade `$http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host `$host;
            proxy_set_header X-Real-IP `$remote_addr;
            proxy_set_header X-Forwarded-For `$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto `$scheme;
            
            # WebSocket保持连接
            proxy_read_timeout 3600;
            proxy_send_timeout 3600;
        }
        
        # 错误页面
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }
}
"@

# 写入新配置
Set-Content -Path $nginxConfPath -Value $newConfig -Encoding UTF8
Write-Host "✅ nginx配置已更新" -ForegroundColor Green

# 测试nginx配置
Write-Host "🔍 测试nginx配置..." -ForegroundColor Yellow
try {
    & nginx -t
    Write-Host "✅ nginx配置测试通过" -ForegroundColor Green
} catch {
    Write-Host "❌ nginx配置测试失败" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

# 启动nginx
Write-Host "🔄 启动nginx..." -ForegroundColor Yellow
try {
    Start-Process nginx -WindowStyle Hidden
    Start-Sleep -Seconds 2
    
    # 检查nginx是否启动
    $nginxProcess = Get-Process nginx -ErrorAction SilentlyContinue
    if ($nginxProcess) {
        Write-Host "✅ nginx已启动" -ForegroundColor Green
    } else {
        Write-Host "❌ nginx启动失败" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ nginx启动失败" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "🎉 部署完成！" -ForegroundColor Green
Write-Host ""
Write-Host "📋 下一步操作：" -ForegroundColor Cyan
Write-Host "1. 启动Go后端服务: go run cmd/server/main.go"
Write-Host "2. 启动前端开发服务: cd frontend && npm run dev"
Write-Host "3. 访问应用: http://localhost"
Write-Host ""
Write-Host "📊 服务状态检查:" -ForegroundColor Cyan
Write-Host "• 查看nginx进程: Get-Process nginx"
Write-Host "• 停止nginx: nginx -s stop"
Write-Host "• 重启nginx: nginx -s reload"
Write-Host ""
Write-Host "🔧 故障排除:" -ForegroundColor Cyan
Write-Host "• 如果80端口被占用，请修改nginx.conf中的端口号"
Write-Host "• 配置文件位置: $nginxConfPath"
Write-Host "• 日志文件位置: $nginxDir\logs\" 