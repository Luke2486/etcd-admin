#!/bin/bash

# API测试脚本

set -e

# 服务器地址
SERVER_URL="http://localhost:8080"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== etcd-admin Backend API Tests ===${NC}"

# 检查服务器是否运行
echo -e "${BLUE}Checking server status...${NC}"
if curl -s -f "$SERVER_URL/health" > /dev/null; then
    echo -e "${GREEN}✓ Server is running${NC}"
else
    echo -e "${RED}✗ Server is not running. Please start the server first.${NC}"
    exit 1
fi

# 测试健康检查
echo -e "\n${BLUE}Testing health endpoint...${NC}"
curl -s "$SERVER_URL/health" | jq '.'

# 测试用户注册
echo -e "\n${BLUE}Testing user registration...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')

echo "$REGISTER_RESPONSE" | jq '.'

# 测试用户登录
echo -e "\n${BLUE}Testing user login...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq '.'

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // empty')

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ Login successful, token obtained${NC}"
    
    # 测试获取用户信息
    echo -e "\n${BLUE}Testing get profile with token...${NC}"
    curl -s -X GET "$SERVER_URL/api/v1/auth/profile" \
      -H "Authorization: Bearer $TOKEN" | jq '.'
    
    # 测试管理员登录
    echo -e "\n${BLUE}Testing admin login...${NC}"
    ADMIN_LOGIN_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/auth/login" \
      -H "Content-Type: application/json" \
      -d '{
        "username": "admin",
        "password": "admin123"
      }')
    
    echo "$ADMIN_LOGIN_RESPONSE" | jq '.'
    
    # 提取管理员token
    ADMIN_TOKEN=$(echo "$ADMIN_LOGIN_RESPONSE" | jq -r '.data.token // empty')
    
    if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "null" ]; then
        echo -e "${GREEN}✓ Admin login successful${NC}"
        
        # 测试管理员端点
        echo -e "\n${BLUE}Testing admin endpoint...${NC}"
        curl -s -X GET "$SERVER_URL/api/v1/admin/users" \
          -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
        
        # 测试连接管理
        echo -e "\n${BLUE}Testing connection management...${NC}"
        # 列出现有连接
        curl -s -X GET "$SERVER_URL/api/v1/connections" -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
        # 创建测试连接（假设本地 etcd）
        CREATE_CONN_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/connections" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $ADMIN_TOKEN" \
          -d '{"name":"local-etcd","endpoints":["localhost:2379"]}')
        echo "$CREATE_CONN_RESPONSE" | jq '.'
        # 提取 connection_id
        CONN_ID=$(echo "$CREATE_CONN_RESPONSE" | jq -r '.data.id // empty')

        echo -e "Connection ID: $CONN_ID"

        # 测试 KV 操作
        echo -e "\n${BLUE}Testing KV operations...${NC}"
        # 列出 keys
        curl -s -X GET "$SERVER_URL/api/v1/kv?connection_id=$CONN_ID" -H "Authorization: Bearer $TOKEN" | jq '.'
        # 设置一个 key
        curl -s -X PUT "$SERVER_URL/api/v1/kv/sampleKey?connection_id=$CONN_ID" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d '{"value":{"foo":"bar"}}' | jq '.'
        # 获取该 key
        curl -s -X GET "$SERVER_URL/api/v1/kv/sampleKey?connection_id=$CONN_ID" -H "Authorization: Bearer $TOKEN" | jq '.'

        # 测试导出备份
        echo -e "\n${BLUE}Testing backup export...${NC}"
        curl -s -X GET "$SERVER_URL/api/v1/backup/export?connection_id=$CONN_ID" -H "Authorization: Bearer $ADMIN_TOKEN" -o backup.json
        cat backup.json | jq '.'

        # 提示请审阅脚本内容后执行测试
        echo -e "\n${YELLOW}Tests prepared. Please review the new etcd tests above before running.${NC}"
    else
        echo -e "${YELLOW}⚠ Admin login failed or token not found${NC}"
    fi
    
else
    echo -e "${YELLOW}⚠ Login failed or token not found${NC}"
fi

echo -e "\n${GREEN}=== API Tests Completed ===${NC}"
