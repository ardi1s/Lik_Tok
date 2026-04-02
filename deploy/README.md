# Lik_tok Deployment 🚀

> Containerized deployment with Docker Compose

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Nginx)                          │
│                           Port: 80                               │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Backend API (Go Gin)                         │
│                           Port: 8080                             │
└─────────────────────────────────────────────────────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│    MySQL    │      │    Redis    │      │  RabbitMQ   │
│    :3306    │      │    :6379    │      │  :5672      │
└─────────────┘      └─────────────┘      └─────────────┘
                                                  │
                                                  ▼
                                        ┌─────────────────┐
                                        │     Worker      │
                                        │  (Background)   │
                                        └─────────────────┘
```

## 📦 Services

| Service | Container | Image | Ports | Purpose |
|---------|-----------|-------|-------|---------|
| MySQL | `lik_tok_mysql` | mysql:8.0 | 3306 | Primary database |
| Redis | `lik_tok_redis` | redis:7-alpine | 6379 | Cache layer |
| RabbitMQ | `lik_tok_rabbitmq` | rabbitmq:3-management-alpine | 5672, 15672 | Message broker |
| Backend | `lik_tok_backend` | docker-backend | 8080 | API server |
| Worker | `lik_tok_worker` | docker-worker | - | Background processor |
| Frontend | `lik_tok_frontend` | docker-frontend | 80 | Web UI |

## 🔐 Credentials

| Service | Username | Password |
|---------|----------|----------|
| MySQL | root | `123456` |
| Redis | - | `123456` |
| RabbitMQ | guest | `guest` |

## 🚀 Quick Start

```bash
# Navigate to deploy directory
cd deploy/docker

# Build and start all services
docker-compose up -d --build

# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

## 🌐 Service URLs

| Service | URL |
|---------|-----|
| Frontend | http://localhost |
| Backend API | http://localhost:8080 |
| RabbitMQ Management | http://localhost:15672 |

---

## 📁 File Structure

```
deploy/
└── docker/
    ├── Dockerfile.backend      # Backend multi-stage build
    ├── Dockerfile.frontend    # Frontend build + Nginx
    ├── Dockerfile.worker      # Worker multi-stage build
    ├── docker-compose.yml    # Service orchestration
    └── nginx.conf           # Nginx configuration
```

## 🔧 docker-compose.yml

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: lik_tok_mysql
    environment:
      MYSQL_ROOT_PASSWORD: "123456"
      MYSQL_DATABASE: Lik_tok
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: lik_tok_redis
    command: redis-server --requirepass 123456
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"

  rabbitmq:
    image: rabbitmq:3-management-alpine
    container_name: lik_tok_rabbitmq
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
    container_name: lik_tok_backend
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
    container_name: lik_tok_worker
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
    container_name: lik_tok_frontend
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  mysql_data:
  redis_data:
```

---

## 🐳 Dockerfiles

### Backend (Dockerfile.backend)

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/configs/config.docker.yaml ./configs/config.yaml
EXPOSE 8080
CMD ["./server"]
```

### Worker (Dockerfile.worker)

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker/main.go

# Stage 2: Runtime
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/worker .
COPY --from=builder /app/configs/config.docker.yaml ./configs/config.yaml
CMD ["./worker"]
```

### Frontend (Dockerfile.frontend)

```dockerfile
# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Runtime
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

---

## 📋 Common Commands

### Service Management

```bash
# Start services
docker-compose up -d

# Rebuild and start
docker-compose up -d --build

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Restart specific service
docker-compose restart backend
docker-compose restart worker
```

### Log Viewing

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f worker
docker-compose logs -f mysql

# Last 100 lines
docker-compose logs --tail 100 backend
```

### Container Access

```bash
# Backend shell
docker exec -it lik_tok_backend sh

# MySQL CLI
docker exec -it lik_tok_mysql mysql -u root -p123456

# Redis CLI
docker exec -it lik_tok_redis redis-cli -a 123456

# RabbitMQ CLI
docker exec -it lik_tok_rabbitmq rabbitmqctl status
```

### Data Operations

```bash
# Backup database
docker exec lik_tok_mysql mysqldump -u root -p123456 Lik_tok > backup.sql

# Restore database
docker exec -i lik_tok_mysql mysql -u root -p123456 Lik_tok < backup.sql

# Clear Redis cache
docker exec lik_tok_redis redis-cli -a 123456 FLUSHALL

# Check RabbitMQ queues
docker exec lik_tok_rabbitmq rabbitmqctl list_queues
```

---

## 🐛 Troubleshooting

### Service Won't Start

```bash
# Check port conflicts
lsof -i :8080
lsof -i :3306
lsof -i :80

# Check logs
docker-compose logs --no-color > logs.txt
```

### Database Connection Issues

```bash
# Verify MySQL is running
docker-compose ps mysql

# Check MySQL logs
docker-compose logs mysql

# Test connection
docker exec -it lik_tok_mysql mysql -u root -p123456 -e "SELECT 1"
```

### Cache Problems

```bash
# Clear Redis
docker exec lik_tok_redis redis-cli -a 123456 FLUSHALL

# Restart backend
docker-compose restart backend worker
```

### RabbitMQ Authentication Failed

```bash
# Rebuild RabbitMQ with fresh data
docker-compose down -v rabbitmq
docker-compose up -d rabbitmq
```

### Full Reset

```bash
# Stop everything and remove volumes
docker-compose down -v

# Fresh start
docker-compose up -d --build
```

---

## 📊 Health Checks

Each service has health checks configured:

| Service | Check | Interval |
|---------|-------|----------|
| MySQL | `mysqladmin ping` | 10s |
| Backend | HTTP `localhost:8080` | - |
| Frontend | HTTP `localhost:80` | - |

---

## 🌱 Environment Variables

### MySQL
| Variable | Value |
|----------|-------|
| `MYSQL_ROOT_PASSWORD` | `123456` |
| `MYSQL_DATABASE` | `Lik_tok` |

### Redis
| Variable | Value |
|----------|-------|
| `requirepass` | `123456` |

### RabbitMQ
| Variable | Value |
|----------|-------|
| `RABBITMQ_DEFAULT_USER` | `guest` |
| `RABBITMQ_DEFAULT_PASS` | `guest` |

---

## 🔒 Production Checklist

- [ ] Change default passwords
- [ ] Configure SSL/TLS certificates
- [ ] Set up resource limits (CPU/memory)
- [ ] Configure log rotation
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure backup strategy
- [ ] Use external storage volumes
- [ ] Set up reverse proxy with load balancing
