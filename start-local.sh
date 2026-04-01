#!/bin/bash

# 本地开发启动脚本

echo "Starting Lik_tok local development environment..."

# 检查 MySQL
if ! mysql -u root -p123456 -e "SELECT 1" &> /dev/null; then
    echo "Error: MySQL is not running or connection failed"
    echo "Please ensure MySQL is running: brew services start mysql"
    exit 1
fi

# 检查 Redis
if ! redis-cli ping &> /dev/null; then
    echo "Error: Redis is not running"
    echo "Please ensure Redis is running: brew services start redis"
    exit 1
fi

# 检查 RabbitMQ
if ! rabbitmqctl status &> /dev/null; then
    echo "Error: RabbitMQ is not running"
    echo "Please ensure RabbitMQ is running: brew services start rabbitmq"
    exit 1
fi

echo "✓ Infrastructure services are running"

# 创建数据库
mysql -u root -p123456 -e "CREATE DATABASE IF NOT EXISTS Lik_tok CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "✓ Database ready"

# 启动后端服务
cd backend

echo "Starting API server..."
go run cmd/main.go &
API_PID=$!

echo "Starting Worker..."
go run cmd/worker/main.go &
WORKER_PID=$!

cd ..

# 启动前端服务
cd frontend

echo "Starting Frontend..."
npm run dev &
FRONTEND_PID=$!

cd ..

echo ""
echo "=========================================="
echo "Lik_tok is running!"
echo "API Server: http://localhost:8080"
echo "Frontend: http://localhost:5173"
echo "=========================================="
echo "Press Ctrl+C to stop all services"

trap "kill $API_PID $WORKER_PID $FRONTEND_PID 2>/dev/null; exit" INT
wait
