#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Lik_tok Docker 启动${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

case "${1:-up}" in
    up)
        echo -e "${YELLOW}构建并启动所有服务...${NC}"
        cd deploy/docker || exit 1
        docker-compose up -d --build
        echo ""
        echo -e "${GREEN}========================================${NC}"
        echo -e "${GREEN}  服务已启动！${NC}"
        echo -e "${GREEN}========================================${NC}"
        echo -e "${YELLOW}前端:${NC} http://localhost:80"
        echo -e "${YELLOW}后端:${NC} http://localhost:8080"
        echo -e "${YELLOW}RabbitMQ:${NC} http://localhost:15672"
        ;;
    down)
        cd deploy/docker || exit 1
        echo -e "${YELLOW}停止所有服务...${NC}"
        docker-compose down
        ;;
    restart)
        cd deploy/docker || exit 1
        echo -e "${YELLOW}重启所有服务...${NC}"
        docker-compose restart
        ;;
    logs)
        cd deploy/docker || exit 1
        docker-compose logs -f
        ;;
    ps)
        cd deploy/docker || exit 1
        docker-compose ps
        ;;
    *)
        echo "用法: ./docker.sh [up|down|restart|logs|ps]"
        ;;
esac
