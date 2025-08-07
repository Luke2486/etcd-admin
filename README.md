# etcd-admin

🚀 現代化的 etcd 管理平台，提供直觀的 Web 界面用於管理和監控 etcd 集群。

## ✨ 項目特色

- 🖥️ **現代化前端**: Vue 3 + TypeScript + TailwindCSS
- 🚀 **高性能後端**: Go + Gin + GORM
- 💾 **零配置數據庫**: 默認使用 SQLite，開箱即用
- 🔐 **安全認證**: JWT 身份驗證系統
- 🌳 **可視化界面**: JSON 樹狀視圖，支持節點選擇和複製
- 📦 **容器化支持**: Docker 一鍵部署
- 🔧 **靈活配置**: 支持 SQLite/MySQL 數據庫切換

## 🚀 技术栈

### 前端
- **Vite** - 现代化构建工具，支持热模块替换
- **Vue 3** - 渐进式JavaScript框架，使用Composition API
- **TailwindCSS** - 实用优先的CSS框架
- **TypeScript** - 类型安全的JavaScript开发

### 后端
- **Go 1.23** - 静态类型编程语言
- **Gin 1.10** - 高性能HTTP Web框架
- **GORM 1.30** - 对象关系映射库
- **SQLite/MySQL** - 关系型数据库（默认使用SQLite进行开发）
- **JWT 4.5** - JSON Web Token用于身份验证
- **etcd Client v3.6** - 官方etcd客户端库
- **Redis** - 内存数据结构存储（可选，用于缓存）

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
- Go 1.23+
- SQLite (自動創建，無需額外安裝)
- etcd 集群 (用於測試和管理)
- Docker & Docker Compose (可選，用於容器化部署)

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

Docker 環境提供完整的生產級部署，包含所有依賴服務：

```bash
# 构建并启动所有服务
docker-compose -f docker/docker-compose.yml up -d

# 停止服务
docker-compose -f docker/docker-compose.yml down

# 查看服务状态
docker-compose -f docker/docker-compose.yml ps

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f [service_name]
```

### Docker 服務配置

| 服務 | 端口 | 說明 |
|------|------|------|
| 前端 (Nginx) | <http://localhost:3000> | Vue.js 應用，Nginx 提供服務 |
| 後端 (Go) | <http://localhost:8080> | Gin API 服務器 |
| MySQL | localhost:3306 | 生產數據庫 |
| Redis | localhost:6379 | 緩存和會話存儲 |
| etcd | <http://localhost:2379> | 測試用 etcd 集群 |

### 部署環境差異

| 配置項 | 本地開發 | Docker 部署 |
|--------|----------|-------------|
| **數據庫** | SQLite (零配置) | MySQL 8.0 |
| **緩存** | 可選 Redis | Redis 7 |
| **前端服務** | Vite 開發服務器 | Nginx 生產服務器 |
| **etcd** | 外部連接 | 內建測試集群 |

> **重要說明**:
>
> - 🔧 **本地開發**: 使用 SQLite，快速啟動，無需額外依賴
> - 🚀 **Docker 部署**: 使用 MySQL + Redis，完整生產環境模擬

## � 數據庫配置

本項目支持兩種數據庫配置：

### SQLite (默認推薦)

- ✅ **零配置**: 無需額外安裝或設置
- ✅ **輕量級**: 適合開發和小型部署
- ✅ **便攜性**: 數據庫文件可輕松遷移
- ✅ **快速啟動**: 專案啟動時自動創建數據庫

數據庫文件位置：`backend/data/etcd-admin.db`

### MySQL (生產環境可選)

- ✅ **高性能**: 適合大型部署和高並發
- ✅ **可擴展**: 支持集群和複製
- ✅ **豐富功能**: 完整的 SQL 支持

切換到 MySQL：更新 `.env` 文件中的 `DB_TYPE=mysql` 並提供 MySQL 連接信息。

## 📝 開發指南

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

## �📝 開發指南

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

## 🔧 環境變量

### 本地開發配置 (SQLite 默認)

```bash
# 服務器配置
SERVER_PORT=8080
GIN_MODE=debug

# 數據庫配置 (默認使用 SQLite)
DB_TYPE=sqlite
DB_PATH=data/etcd-admin.db

# JWT配置
JWT_SECRET=your-very-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=24h

# Redis配置 (可選，用於緩存)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# CORS配置
CORS_ORIGINS=http://localhost:5173,http://localhost:3000

# 默認管理員賬戶（首次運行時創建）
ADMIN_USERNAME=admin
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=admin123
```

### Docker 生產配置 (MySQL + Redis)

Docker 環境自動使用以下配置：

```bash
# 服務器配置
SERVER_PORT=8080
GIN_MODE=release

# 數據庫配置 (Docker 環境使用 MySQL)
DB_TYPE=mysql
DB_HOST=mysql
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=rootpassword
DB_DATABASE=etcd_admin

# Redis配置 (Docker 環境內建 Redis)
REDIS_HOST=redis
REDIS_PORT=6379

# JWT配置
JWT_SECRET=your-production-jwt-secret-key
JWT_EXPIRES_IN=24h
```

### MySQL 手動配置 (可選)

如果您希望在本地開發中使用 MySQL 而非 SQLite：

```bash
# 數據庫配置 (本地 MySQL)
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_DATABASE=etcd_admin
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
