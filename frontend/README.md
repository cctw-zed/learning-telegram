# Telegram Chat 前端

基于Vue 3 + TypeScript的即时通讯前端应用。

## 技术栈

- Vue 3 (Composition API)
- TypeScript
- Vue Router 4 (路由管理)
- Pinia (状态管理)
- Vite (构建工具)
- JWT (身份验证)

## 项目结构

```
src/
├── components/     # 可复用组件
├── views/         # 页面组件
│   ├── LoginView.vue    # 登录/注册页面
│   └── ChatView.vue     # 聊天页面
├── stores/        # Pinia状态管理
│   └── auth.ts          # 认证状态
├── router/        # 路由配置
│   └── index.ts         # 路由定义
├── types/         # TypeScript类型定义
│   └── index.ts         # 通用类型
├── services/      # API服务
├── App.vue        # 根组件
└── main.ts        # 应用入口
```

## 主要功能

- 用户注册/登录
- JWT身份验证
- 实时聊天 (WebSocket)
- 私聊和群聊
- 聊天列表
- 消息历史
- 响应式界面

## 开发命令

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

## 后端API接口

- POST `/api/register` - 用户注册
- POST `/api/login` - 用户登录
- GET `/api/me/chats` - 获取聊天列表
- WebSocket `/ws?token=<jwt_token>` - 实时通讯

## 环境要求

- Node.js >= 16
- npm >= 8
- 后端服务运行在 `http://localhost:8080`
