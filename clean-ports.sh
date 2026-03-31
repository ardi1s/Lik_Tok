#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Lik_tok 端口清理工具${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 需要清理的端口
PORTS=("3306" "6379" "8080" "5672" "15672")

echo -e "${YELLOW}检查并清理占用端口的进程...${NC}"
echo ""

for port in "${PORTS[@]}"; do
    pid=$(lsof -ti:$port 2>/dev/null)
    if [ -n "$pid" ]; then
        echo -e "端口 ${YELLOW}$port${NC} 被占用 (PID: $pid)，正在终止..."
        kill -9 $pid 2>/dev/null
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓${NC} 端口 $port 已释放"
        else
            echo -e "${RED}✗${NC} 端口 $port 释放失败，可能需要 sudo"
        fi
    else
        echo -e "${GREEN}✓${NC} 端口 $port 未被占用"
    fi
done

echo ""
echo -e "${YELLOW}停止本地 MySQL 服务...${NC}"
brew services stop mysql 2>/dev/null && echo -e "${GREEN}✓${NC} MySQL 已停止" || echo -e "${YELLOW}!${NC} MySQL 未运行或不是通过 brew 安装"

echo ""
echo -e "${YELLOW}停止本地 Redis 服务...${NC}"
brew services stop redis 2>/dev/null && echo -e "${GREEN}✓${NC} Redis 已停止" || echo -e "${YELLOW}!${NC} Redis 未运行或不是通过 brew 安装"

echo ""
echo -e "${YELLOW}停止本地 RabbitMQ 服务...${NC}"
brew services stop rabbitmq 2>/dev/null && echo -e "${GREEN}✓${NC} RabbitMQ 已停止" || echo -e "${YELLOW}!${NC} RabbitMQ 未运行或不是通过 brew 安装"

echo ""
echo -e "${YELLOW}停止 Lik_tok Docker 容器...${NC}"
cd "$(dirname "$0")/deploy/docker" 2>/dev/null && docker-compose down 2>/dev/null && echo -e "${GREEN}✓${NC} Docker 容器已停止" || echo -e "${YELLOW}!${NC} Docker 容器未运行"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  端口清理完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "现在可以运行: ${YELLOW}./docker.sh up${NC}"
