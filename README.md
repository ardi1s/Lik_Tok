# Lik_tok 🎬

> A high-performance short video social platform built with Go and Vue3

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4-green.svg)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-purple.svg)](LICENSE)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🎥 **Video Management** | Upload, play, and delete videos with streaming support |
| 👍 **Like System** | Real-time like/unlike with accurate count updates |
| 💬 **Comments** | Threaded comments with real-time updates |
| 👥 **Social Graph** | Follow/unfollow users with activity feeds |
| 🏠 **Smart Feed** | Timeline-based video recommendations |
| 🔥 **Popularity Ranking** | Engagement-based video ranking system |
| 🔐 **Authentication** | Secure JWT-based authentication |

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend                                 │
│                     Vue 3 + TypeScript                          │
│                    Vite + Pinia + Vue Router                     │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway                                │
│                      Gin + JWT Middleware                         │
└─────────────────────────────────────────────────────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   MySQL 8   │      │    Redis    │      │  RabbitMQ   │
│  (Storage)  │      │   (Cache)   │      │    (MQ)     │
└─────────────┘      └─────────────┘      └─────────────┘
                                                  │
                                                  ▼
                                        ┌─────────────────┐
                                        │     Worker      │
                                        │  (Background)   │
                                        └─────────────────┘
```

## 📁 Project Structure

```
lik_tok/
├── backend/                     # Go Backend API
│   ├── cmd/                     # Application entry points
│   │   ├── main.go              # API server
│   │   └── worker/              # Background worker
│   ├── configs/                 # Configuration files
│   │   ├── config.yaml          # Local development
│   │   └── config.docker.yaml   # Docker deployment
│   └── internal/
│       ├── account/             # User management
│       ├── video/               # Video, likes, comments
│       ├── feed/                # Feed generation
│       ├── social/              # Follow relationships
│       ├── middleware/          # JWT, Redis, RabbitMQ
│       └── worker/              # Async task processors
│
├── frontend/                    # Vue3 Frontend
│   └── src/
│       ├── api/                 # API clients
│       ├── components/          # Reusable components
│       ├── views/               # Page views
│       ├── stores/              # Pinia stores
│       └── router/              # Route definitions
│
├── deploy/                      # Deployment configs
│   └── docker/                  # Docker Compose setup
│
└── README.md                    # This file
```

---

## 🚀 Quick Start

### Prerequisites

| Component | Version | Install |
|----------|---------|---------|
| Docker | 20.10+ | [Guide](https://docs.docker.com/get-docker/) |
| Docker Compose | 2.0+ | [Guide](https://docs.docker.com/compose/install/) |

### One-Command Launch

```bash
# Clone the repository
git clone https://github.com/yourusername/lik_tok.git
cd lik_tok

# Start all services
cd deploy/docker && docker-compose up -d

# Access the application
open http://localhost
```

---

## 🛠️ Development Setup

### Option 1: Docker Development (Recommended)

```bash
cd deploy/docker

# Build and start
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Option 2: Local Development

**Requirements:**
- Go 1.25+
- Node.js 20+
- MySQL 8.0
- Redis 7+
- RabbitMQ 3.12

**Launch Script:**

```bash
# Start backend (API + Worker) and frontend
./start-local.sh

# Or manual startup
# Terminal 1: API Server
cd backend && go run cmd/main.go

# Terminal 2: Worker
cd backend && go run cmd/worker/main.go

# Terminal 3: Frontend
cd frontend && npm install && npm run dev
```

> ⚠️ **Note:** Docker and local development cannot run simultaneously due to port conflicts.

---

## ⚙️ Configuration

### Unified Credentials

| Service | Username | Password |
|---------|----------|----------|
| MySQL | root | `123456` |
| Redis | - | `123456` |
| RabbitMQ | guest | `guest` |

### Configuration Files

| Environment | File Path |
|-------------|-----------|
| Local | `backend/configs/config.yaml` |
| Docker | `backend/configs/config.docker.yaml` |

### Configuration Example

```yaml
server:
  port: 8080

database:
  host: localhost      # or 'mysql' for Docker
  port: 3306
  user: root
  password: "123456"
  dbname: Lik_tok

redis:
  host: localhost      # or 'redis' for Docker
  port: 6379
  password: "123456"

rabbitmq:
  host: localhost      # or 'rabbitmq' for Docker
  port: 5672
  username: guest
  password: guest
```

---

## 🔌 API Reference

### Authentication

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/account/register` | POST | Register new user |
| `/account/login` | POST | User login |
| `/account/changePassword` | POST | Change password |

### Video

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/video/publish` | POST | Upload video |
| `/video/listByAuthorID` | POST | Get user's videos |
| `/video/delete` | POST | Delete video |

### Interaction

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/like/like` | POST | Like a video |
| `/like/unlike` | POST | Unlike a video |
| `/like/isLiked` | POST | Check like status |
| `/comment/publish` | POST | Add comment |
| `/comment/list` | POST | Get comments |

### Feed

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/feed/listLatest` | POST | Latest videos |
| `/feed/listByPopularity` | POST | Popular videos |

---

## 📊 Data Architecture

### Multi-Level Cache

```
Request → L1 (Local) → L2 (Redis) → L3 (MySQL)
         5 seconds     1 hour        Source
```

### Message Queues

| Queue | Purpose |
|-------|---------|
| `like.events` | Like/unlike processing |
| `comment.events` | Comment notifications |
| `social.events` | Follow relationship updates |
| `video.popularity.events` | Popularity calculation |

---

## 🔧 Service Ports

| Service | Port | URL |
|---------|------|-----|
| Frontend | 80 | http://localhost |
| Backend API | 8080 | http://localhost:8080 |
| RabbitMQ | 15672 | http://localhost:15672 |

---

## 🐛 Troubleshooting

### Database Connection Failed

```bash
# Check MySQL status
docker-compose ps mysql
mysql -u root -p123456

# Reset password
ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Clear Cache

```bash
# Redis (Docker)
docker exec lik_tok_redis redis-cli -a 123456 FLUSHALL

# Rebuild services
docker-compose up -d --build
```

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Gin Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://vuejs.org/)
- [GORM](https://gorm.io/)
- [RabbitMQ](https://www.rabbitmq.com/)
