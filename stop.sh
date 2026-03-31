#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Lik_tok 停止所有服务${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 停止后端 (go run)
echo -e "${YELLOW}停止后端服务...${NC}"
pkill -f "go run cmd/main.go" 2>/dev/null && echo -e "${GREEN}后端已停止${NC}" || echo -e "${YELLOW}后端未运行${NC}"

# 停止前端 (vite)
echo -e "${YELLOW}停止前端服务...${NC}"
pkill -f "vite" 2>/dev/null && echo -e "${GREEN}前端已停止${NC}" || echo -e "${YELLOW}前端未运行${NC}"

# 检查端口并强制停止
echo -e "${YELLOW}检查端口占用...${NC}"

for port in 8080 5173 5174 5175 5176; do
    pid=$(lsof -ti:$port 2>/dev/null)
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}释放端口 $port (PID: $pid)${NC}"
        kill -9 $pid 2>/dev/null && echo -e "${GREEN}端口 $port 已释放${NC}" || echo -e "${YELLOW}端口 $port 释放失败${NC}"
    fi
done

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  所有服务已停止${NC}"
echo -e "${GREEN}========================================${NC}"
