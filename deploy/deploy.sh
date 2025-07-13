#!/bin/bash

# Learning Telegram - Nginx反向代理部署脚本
# 使用方法: ./deploy.sh

echo "🚀 开始部署Learning Telegram项目..."

# 检查是否安装了nginx
if ! command -v nginx &> /dev/null; then
    echo "❌ 未找到nginx，正在安装..."
    
    # 检测操作系统
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Ubuntu/Debian
        if command -v apt-get &> /dev/null; then
            sudo apt-get update
            sudo apt-get install -y nginx
        # CentOS/RHEL
        elif command -v yum &> /dev/null; then
            sudo yum install -y nginx
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install nginx
        else
            echo "❌ 请先安装Homebrew: https://brew.sh/"
            exit 1
        fi
    else
        echo "❌ 不支持的操作系统，请手动安装nginx"
        exit 1
    fi
fi

echo "✅ nginx已安装"

# 停止现有的nginx服务
echo "🔄 停止现有的nginx服务..."
sudo nginx -s stop 2>/dev/null || true

# 备份原有nginx配置
if [ -f /etc/nginx/sites-available/default ]; then
    sudo cp /etc/nginx/sites-available/default /etc/nginx/sites-available/default.backup
    echo "✅ 已备份原有nginx配置"
fi

# 复制nginx配置
echo "📝 配置nginx反向代理..."
sudo cp nginx.conf /etc/nginx/sites-available/learning-telegram

# 创建软链接（如果不存在）
if [ ! -f /etc/nginx/sites-enabled/learning-telegram ]; then
    sudo ln -s /etc/nginx/sites-available/learning-telegram /etc/nginx/sites-enabled/
fi

# 删除默认配置的软链接（避免冲突）
if [ -f /etc/nginx/sites-enabled/default ]; then
    sudo rm /etc/nginx/sites-enabled/default
fi

# 测试nginx配置
echo "🔍 测试nginx配置..."
sudo nginx -t

if [ $? -eq 0 ]; then
    echo "✅ nginx配置测试通过"
else
    echo "❌ nginx配置测试失败"
    exit 1
fi

# 启动nginx
echo "🔄 启动nginx..."
sudo nginx

# 检查nginx状态
if pgrep nginx > /dev/null; then
    echo "✅ nginx已启动"
else
    echo "❌ nginx启动失败"
    exit 1
fi

echo ""
echo "🎉 部署完成！"
echo ""
echo "📋 下一步操作："
echo "1. 启动Go后端服务: cd backend && go run cmd/server/main.go"
echo "2. 启动前端开发服务: cd frontend && npm run dev"
echo "3. 访问应用: http://localhost"
echo ""
echo "📊 服务状态检查:"
echo "• nginx状态: sudo systemctl status nginx"
echo "• nginx日志: sudo tail -f /var/log/nginx/error.log"
echo "• 停止nginx: sudo nginx -s stop"
echo "• 重启nginx: sudo nginx -s reload"
echo ""
echo "🔧 故障排除:"
echo "• 如果80端口被占用，请修改nginx.conf中的端口号"
echo "• 如果权限问题，请检查nginx用户权限"
echo "• 配置文件位置: /etc/nginx/sites-available/learning-telegram" 