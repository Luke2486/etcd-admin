#!/bin/bash

# 前端构建脚本
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

echo "Frontend build completed!"
