# etcd-admin

一个现代化的etcd管理平台，提供直观的Web界面用于管理和监控etcd集群。

## 🚀 技术栈

### 前端
- **Vite** - 现代化构建工具，支持热模块替换
- **Vue 3** - 渐进式JavaScript框架，使用Composition API
- **TailwindCSS** - 实用优先的CSS框架
- **TypeScript** - 类型安全的JavaScript开发

### 后端
- **Go** - 静态类型编程语言
- **Gin** - 高性能HTTP Web框架
- **MySQL** - 关系型数据库，用于持久化数据
- **Redis** - 内存数据结构存储，用于缓存
- **etcd** - 分布式键值存储（目标管理系统）

## 📁 项目结构

```
etcd-admin/
├── frontend/           # Vue 3 + Vite前端
│   ├── src/           # 源代码
│   ├── public/        # 静态资源
│   └── dist/          # 构建输出
├── backend/           # Go + Gin后端
│   ├── cmd/           # 应用程序入口
│   ├── internal/      # 内部包
│   └── pkg/           # 可重用包
├── docker/            # Docker配置
├── docs/              # 项目文档
├── scripts/           # 构建和部署脚本
└── .github/           # GitHub配置
```

## 🛠️ 开发环境设置

### 前置要求
- Node.js 20+
- Go 1.21+
- MySQL 8.0+
- Redis 7+
- Docker & Docker Compose

### 快速开始

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd etcd-admin
   ```

2. **启动开发环境**
   ```bash
   # 启动前端开发服务器
   cd frontend
   npm install
   npm run dev
   
   # 启动后端服务器（新终端）
   cd backend
   go run cmd/server/main.go
   ```

3. **使用Docker启动完整环境**
   ```bash
   docker-compose -f docker/docker-compose.yml up -d
   ```

### 开发服务器
- 前端：http://localhost:5173 (Vite开发服务器)
- 后端：http://localhost:8080 (Gin服务器)
- 前端使用Vite代理配置处理API调用

## 🐳 Docker部署

```bash
# 构建并启动所有服务
docker-compose -f docker/docker-compose.yml up -d

# 停止服务
docker-compose -f docker/docker-compose.yml down
```

### 服务端口
- 前端：http://localhost:3000
- 后端：http://localhost:8080
- MySQL：localhost:3306
- Redis：localhost:6379
- etcd：http://localhost:2379

## 📝 开发指南

### 前端开发
- 使用Vue 3 Composition API和`<script setup>`语法
- 优先使用`ref()`和`reactive()`进行响应式状态管理
- 使用TailwindCSS类进行样式设计，避免自定义CSS
- 遵循TypeScript最佳实践，提供适当的类型定义

### 后端开发
- 遵循Go约定和标准项目布局
- 使用Gin中间件处理CORS、认证、日志
- 实现适当的错误处理和HTTP状态码
- 使用结构体标签进行JSON序列化
- 遵循RESTful API设计原则

### API设计
- 使用RESTful端点和适当的HTTP方法
- 实现一致的响应格式（status、data、message字段）
- 添加输入验证和清理
- 使用JWT令牌进行认证
- 实现适当的错误响应

## 🔧 环境变量

### 后端配置
```bash
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_DATABASE=etcd_admin

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

SERVER_PORT=8080
JWT_SECRET=your-secret-key
GIN_MODE=debug
```

## 📚 API文档

API文档将在后续版本中提供。

## 🤝 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 📄 许可证

此项目使用MIT许可证。详见LICENSE文件。
