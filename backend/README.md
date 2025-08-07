# etcd-admin Backend

基于Go + Gin框架的etcd管理平台后端服务。

## 功能特性

- ✅ JWT认证系统
- ✅ 数据库迁移功能
- ✅ 用户注册/登录
- ✅ 角色权限管理
- ✅ CORS跨域支持
- ✅ 结构化配置管理
- 🚧 etcd集群管理（待实现）

## 技术栈

- **Go 1.23+** - 编程语言
- **Gin** - HTTP Web框架
- **GORM** - ORM数据库操作
- **MySQL** - 数据库
- **JWT** - 身份认证
- **Migrate** - 数据库迁移

## 项目结构

```
backend/
├── cmd/server/          # 应用入口
├── internal/
│   ├── config/         # 配置管理
│   ├── handlers/       # HTTP处理器
│   ├── middleware/     # 中间件
│   └── models/         # 数据模型
├── pkg/database/       # 数据库相关
├── migrations/         # 数据库迁移文件
├── bin/               # 编译输出
└── scripts/           # 辅助脚本
```

## 快速开始

### 1. 环境准备

确保已安装：
- Go 1.23+
- MySQL 5.7+
- Git

### 2. 克隆项目

```bash
git clone <your-repo>
cd etcd-admin/backend
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 配置环境

复制环境变量配置文件：
```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接信息：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=etcd_admin
JWT_SECRET=your-secret-key
```

### 5. 创建数据库

在MySQL中创建数据库：
```sql
CREATE DATABASE etcd_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 6. 运行数据库迁移

```bash
# 方法1：使用脚本
./scripts/migrate.sh migrate

# 方法2：直接运行
go run cmd/server/main.go -migrate
```

### 7. 启动服务

```bash
# 开发模式
go run cmd/server/main.go

# 或编译后运行
go build -o bin/server cmd/server/main.go
./bin/server
```

服务将在 `http://localhost:8080` 启动。

## API文档

### 认证相关

#### 用户注册
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

#### 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

响应：
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user",
      "is_active": true
    }
  }
}
```

#### 获取用户信息
```
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

#### 用户登出
```
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

### 健康检查

```
GET /health
```

## 数据库迁移

### 运行迁移
```bash
# 使用脚本
./scripts/migrate.sh migrate

# 直接运行
go run cmd/server/main.go -migrate
```

### 回滚迁移
```bash
# 使用脚本
./scripts/migrate.sh rollback

# 直接运行
go run cmd/server/main.go -rollback
```

### 重置数据库
```bash
./scripts/migrate.sh reset
```

## 默认用户

迁移完成后，系统会创建一个默认管理员用户：
- 用户名: `admin`
- 邮箱: `admin@example.com`
- 密码: `admin123`
- 角色: `admin`

**⚠️ 生产环境请立即更改默认密码！**

## 开发指南

### 添加新的API端点

1. 在 `internal/handlers/` 中创建处理器
2. 在 `internal/handlers/routes.go` 中注册路由
3. 如需要数据库操作，在 `internal/models/` 中定义模型

### 添加数据库迁移

1. 在 `migrations/` 目录下创建新的迁移文件：
   ```
   000002_your_migration_name.up.sql
   000002_your_migration_name.down.sql
   ```

2. 运行迁移：
   ```bash
   ./scripts/migrate.sh migrate
   ```

### 中间件使用

#### JWT认证中间件
```go
protected := r.Group("/api/v1/protected")
protected.Use(middleware.JWTAuth(cfg))
```

#### 角色验证中间件
```go
admin := protected.Group("/admin")
admin.Use(middleware.RequireRole("admin"))
```

## 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USERNAME | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | password |
| DB_DATABASE | 数据库名 | etcd_admin |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 | (空) |
| SERVER_PORT | 服务端口 | 8080 |
| JWT_SECRET | JWT密钥 | your-secret-key |
| GIN_MODE | Gin模式 | debug |

## 部署

### 编译
```bash
go build -o bin/server cmd/server/main.go
```

### Docker部署
```bash
# 构建镜像
docker build -f docker/Dockerfile.backend -t etcd-admin-backend .

# 运行容器
docker run -p 8080:8080 etcd-admin-backend
```

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库是否启动
   - 验证连接参数是否正确
   - 确认数据库用户权限

2. **迁移失败**
   - 检查数据库是否存在
   - 验证用户是否有DDL权限
   - 查看具体错误信息

3. **JWT认证失败**
   - 检查JWT_SECRET是否配置
   - 验证token是否过期
   - 确认请求头格式正确

## 许可证

MIT License
