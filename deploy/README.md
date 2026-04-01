# Lik_tok 部署指南

本项目使用 Docker Compose 进行容器化部署。

## 目录结构

```
deploy/
└── docker/
    ├── Dockerfile.backend    # 后端 API 服务镜像
    ├── Dockerfile.frontend   # 前端应用镜像
    ├── Dockerfile.worker    # 后台工作进程镜像
    ├── docker-compose.yml   # 编排配置
    └── nginx.conf           # Nginx 配置
```

## 服务组成

| 服务 | 容器名 | 镜像 | 端口 | 说明 |
|------|--------|------|------|------|
| mysql | lik_tok_mysql | mysql:8.0 | 3306 | 数据库 |
| redis | lik_tok_redis | redis:7-alpine | 6379 | 缓存 |
| rabbitmq | lik_tok_rabbitmq | rabbitmq:3-management-alpine | 5672, 15672 | 消息队列 |
| backend | lik_tok_backend | docker-backend | 8080 | API 服务 |
| worker | lik_tok_worker | docker-worker | - | 后台任务 |
| frontend | lik_tok_frontend | docker-frontend | 80 | Web 前端 |

## 统一密码配置

| 服务 | 用户名 | 密码 |
|------|--------|------|
| MySQL | root | 123456 |
| Redis | - | 123456 |
| RabbitMQ | guest | guest |

## 快速部署

### 1. 环境准备

确保已安装：
- Docker 20.10+
- Docker Compose 2.0+

### 2. 启动服务

```bash
cd deploy/docker

# 首次构建并启动
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看所有日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f worker
```

### 3. 访问服务

- **Web 应用**: http://localhost
- **后端 API**: http://localhost:8080
- **RabbitMQ 管理界面**: http://localhost:15672 (guest/guest)

### 4. 停止服务

```bash
docker-compose down

# 停止并删除数据卷（谨慎使用）
docker-compose down -v
```

## 配置文件

### docker-compose.yml 核心配置

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: "123456"
      MYSQL_DATABASE: Lik_tok
    ports:
      - "3306:3306"

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass 123456
    ports:
      - "6379:6379"

  rabbitmq:
    image: rabbitmq:3-management-alpine
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    ports:
      - "5672:5672"
      - "15672:15672"

  backend:
    build:
      context: ../../backend
      dockerfile: ../deploy/docker/Dockerfile.backend
    ports:
      - "8080:8080"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  worker:
    build:
      context: ../../backend
      dockerfile: ../deploy/docker/Dockerfile.worker
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  frontend:
    build:
      context: ../../frontend
      dockerfile: ../deploy/docker/Dockerfile.frontend
    ports:
      - "80:80"
    depends_on:
      backend:
        condition: service_healthy
```

## 镜像构建

### 后端 API 镜像 (Dockerfile.backend)

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/configs/config.docker.yaml ./configs/config.yaml
EXPOSE 8080
CMD ["./server"]
```

### Worker 镜像 (Dockerfile.worker)

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/worker .
COPY --from=builder /app/configs/config.docker.yaml ./configs/config.yaml
CMD ["./worker"]
```

### 前端镜像 (Dockerfile.frontend)

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 常用命令

### 服务管理

```bash
# 进入 Docker 目录
cd deploy/docker

# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 重启特定服务
docker-compose restart backend
docker-compose restart worker

# 停止所有服务
docker-compose down

# 强制重建特定服务
docker-compose up -d --build backend
```

### 日志查看

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看后端日志
docker-compose logs -f backend

# 查看 Worker 日志
docker-compose logs -f worker

# 查看最近 100 行日志
docker-compose logs --tail 100 backend
```

### 进入容器

```bash
# 进入后端容器
docker exec -it lik_tok_backend sh

# 进入 MySQL 容器
docker exec -it lik_tok_mysql mysql -u root -p123456

# 进入 Redis 容器
docker exec -it lik_tok_redis redis-cli -a 123456

# 进入 RabbitMQ 容器
docker exec -it lik_tok_rabbitmq rabbitmqctl status
```

### 数据管理

```bash
# 备份数据库
docker exec lik_tok_mysql mysqldump -u root -p123456 Lik_tok > backup.sql

# 恢复数据库
docker exec -i lik_tok_mysql mysql -u root -p123456 Lik_tok < backup.sql

# 清理 Redis 缓存
docker exec lik_tok_redis redis-cli -a 123456 FLUSHALL

# 查看 RabbitMQ 队列
docker exec lik_tok_rabbitmq rabbitmqctl list_queues

# 查看资源使用
docker stats
```

## 故障排查

### 服务无法启动

```bash
# 检查端口占用
lsof -i :8080
lsof -i :3306

# 查看详细日志
docker-compose logs --no-color > logs.txt
```

### 数据库连接失败

```bash
# 检查 MySQL 状态
docker-compose ps mysql
docker-compose logs mysql

# 进入 MySQL
docker exec -it lik_tok_mysql mysql -u root -p123456
```

### 缓存问题

```bash
# 清理 Redis
docker exec lik_tok_redis redis-cli -a 123456 FLUSHALL

# 重启后端和 Worker
docker-compose restart backend worker
```

### RabbitMQ 连接失败

```bash
# 检查 RabbitMQ 状态
docker-compose ps rabbitmq
docker-compose logs rabbitmq

# 完全重建 RabbitMQ
docker-compose down -v rabbitmq
docker-compose up -d rabbitmq
```

## 生产环境部署

### 1. 修改配置

- 更改默认密码
- 配置域名和 SSL 证书
- 调整资源限制（CPU、内存）

### 2. 使用外部存储

```yaml
volumes:
  mysql_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/mysql
```

### 3. 配置反向代理

使用 Nginx 或 Traefik 进行反向代理和负载均衡。

### 4. 监控和日志

- 配置 Prometheus + Grafana 监控
- 使用 ELK 或 Loki 收集日志

## 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build

# 清理旧镜像
docker image prune -f
```
