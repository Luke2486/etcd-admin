# etcd Admin Backend API

## 快速开始

### 1. 环境配置

复制并编辑环境配置文件：
```bash
cp .env.example .env
```

### 2. 数据库迁移

```bash
# 创建SQLite数据库和表结构
go run cmd/server/main.go -migrate
```

### 3. 启动服务

```bash
# 开发模式
go run cmd/server/main.go

# 或者使用VS Code任务
# Ctrl/Cmd + Shift + P -> "Tasks: Run Task" -> "Run Backend"
```

## API 文档

### 认证相关

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `GET /api/v1/auth/profile` - 获取用户信息
- `POST /api/v1/auth/logout` - 用户登出

### etcd 连接管理

- `POST /api/v1/connections` - 创建etcd连接
- `GET /api/v1/connections` - 获取连接列表
- `GET /api/v1/connections/:id` - 获取连接详情
- `DELETE /api/v1/connections/:id` - 删除连接
- `POST /api/v1/connections/:id/test` - 测试连接

#### 创建连接示例：
```json
{
  "name": "本地etcd",
  "endpoints": ["localhost:2379"],
  "username": "",
  "password": "",
  "tls_enabled": false
}
```

### KV 操作

- `GET /api/v1/connections/:connection_id/kv` - 列出所有键（支持前缀过滤）
- `GET /api/v1/connections/:connection_id/kv/:key` - 获取键值
- `PUT /api/v1/connections/:connection_id/kv/:key` - 设置键值
- `DELETE /api/v1/connections/:connection_id/kv/:key` - 删除键

#### 查询参数：
- `prefix` - 前缀过滤（用于列表）

#### 设置键值示例：
```json
{
  "value": "{\"name\": \"test\", \"version\": \"1.0\"}"
}
```

### 备份导入导出

- `GET /api/v1/connections/:connection_id/backup/export` - 导出所有KV数据
- `POST /api/v1/connections/:connection_id/backup/import` - 导入KV数据

#### 导入数据格式：
```json
{
  "data": {
    "key1": "value1",
    "key2": "{\"json\": \"value\"}"
  },
  "overwrite": true
}
```

### 连接间传输

- `POST /api/v1/transfer` - 批量传输KV数据
- `POST /api/v1/transfer/copy/:key` - 复制单个键

#### 批量传输示例：
```json
{
  "source_connection_id": 1,
  "target_connection_id": 2,
  "keys": ["app/config", "app/version"],
  "overwrite": true,
  "key_mapping": true,
  "source_prefix": "app/",
  "target_prefix": "prod/"
}
```

## 测试本地etcd

确保本地etcd服务器运行在 `localhost:2379`：

```bash
# 使用Docker运行etcd
docker run --rm -p 2379:2379 -p 2380:2380 \
  --name etcd-server \
  quay.io/coreos/etcd:latest \
  etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
```

## 开发注意事项

1. **数据库类型**: 默认使用SQLite，可通过环境变量切换到MySQL
2. **认证**: 所有etcd相关API都需要JWT认证
3. **错误处理**: API返回统一的错误格式
4. **连接缓存**: etcd客户端连接会被缓存以提高性能
5. **日志**: 使用Gin的默认日志中间件记录请求
