#!/bin/bash

# 开发环境启动脚本
echo "Starting development environment..."

# 启动后端服务器（在后台）
echo "Starting backend server..."
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

# 等待后端启动
sleep 3

# 启动前端开发服务器
echo "Starting frontend dev server..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo "Development servers started!"
echo "Backend PID: $BACKEND_PID"
echo "Frontend PID: $FRONTEND_PID"
echo "Frontend: http://localhost:5173"
echo "Backend: http://localhost:8080"

# 等待用户输入以停止服务器
read -p "Press any key to stop servers..."

# 停止服务器
kill $BACKEND_PID
kill $FRONTEND_PID

echo "Development servers stopped."
