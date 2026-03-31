#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 清理函数
cleanup() {
    echo -e "\n${YELLOW}正在停止所有服务...${NC}"
    if [ -n "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null
        echo -e "${GREEN}后端服务已停止${NC}"
    fi
    if [ -n "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null
        echo -e "${GREEN}前端服务已停止${NC}"
    fi
    exit 0
}

# 捕获中断信号
trap cleanup INT TERM

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Lik_tok 一键启动脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查依赖
check_dependency() {
    if ! command -v "$1" &> /dev/null; then
        echo -e "${RED}错误: 未找到 $1，请先安装${NC}"
        exit 1
    fi
}

echo -e "${YELLOW}检查依赖...${NC}"
check_dependency "go"
check_dependency "node"
check_dependency "npm"
echo -e "${GREEN}所有依赖已就绪${NC}"
echo ""

# 启动后端
echo -e "${BLUE}启动后端服务...${NC}"
cd "$PROJECT_ROOT/backend" || exit 1

# 检查 go.mod 是否存在
if [ ! -f "go.mod" ]; then
    echo -e "${YELLOW}初始化 Go 模块...${NC}"
    go mod init Lik_tok 2>/dev/null || true
fi

# 下载依赖
echo -e "${YELLOW}下载后端依赖...${NC}"
go mod tidy

# 启动后端
echo -e "${GREEN}后端服务启动中 (端口: 8080)...${NC}"
go run cmd/main.go &
BACKEND_PID=$!

# 等待后端启动
sleep 2

# 检查后端是否成功启动
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${RED}后端服务启动失败${NC}"
    exit 1
fi

echo -e "${GREEN}后端服务已启动 (PID: $BACKEND_PID)${NC}"
echo ""

# 启动前端
echo -e "${BLUE}启动前端服务...${NC}"
cd "$PROJECT_ROOT/frontend" || exit 1

# 检查 node_modules 是否存在
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}安装前端依赖...${NC}"
    npm install
fi

echo -e "${GREEN}前端服务启动中 (端口: 5173)...${NC}"
npm run dev &
FRONTEND_PID=$!

# 等待前端启动
sleep 3

# 检查前端是否成功启动
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    echo -e "${RED}前端服务启动失败${NC}"
    cleanup
    exit 1
fi

echo -e "${GREEN}前端服务已启动 (PID: $FRONTEND_PID)${NC}"
echo ""

# 显示访问信息
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}所有服务已启动！${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${YELLOW}前端访问:${NC} http://localhost:5173"
echo -e "${YELLOW}后端 API:${NC} http://localhost:8080"
echo -e "${YELLOW}API 代理:${NC} http://localhost:5173/api (通过前端代理)"
echo ""
echo -e "${YELLOW}按 Ctrl+C 停止所有服务${NC}"
echo -e "${BLUE}========================================${NC}"

# 等待进程
wait $BACKEND_PID $FRONTEND_PID
