#!/bin/bash

# 后端构建脚本
echo "Building backend..."
cd backend
go mod tidy
go build -o bin/server cmd/server/main.go
cd ..

echo "Backend build completed!"
