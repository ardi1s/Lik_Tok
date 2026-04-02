# Lik_tok Backend ⚙️

> Go-based API server with layered architecture

## 🏗️ Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                          Handler                             │
│              HTTP Request / Response / Validation             │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                          Service                             │
│                   Business Logic Processing                   │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                        Repository                            │
│                    Data Access Layer (DAO)                    │
└──────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
         ┌─────────┐    ┌─────────┐    ┌─────────┐
         │  MySQL  │    │  Redis  │    │ RabbitMQ │
         │ Storage │    │  Cache  │    │    MQ    │
         └─────────┘    └─────────┘    └─────────┘
```

## 📁 Directory Structure

```
backend/
├── cmd/                               # Application entry points
│   ├── main.go                        # API server entry
│   └── worker/
│       └── main.go                    # Background worker entry
│
├── configs/                           # Configuration files
│   ├── config.yaml                    # Local development
│   └── config.docker.yaml             # Docker deployment
│
└── internal/                          # Internal packages
    ├── account/                       # User account module
    │   ├── entity.go                 # Data models
    │   ├── handler.go                # HTTP handlers
    │   ├── repo.go                   # Data access
    │   └── service.go                # Business logic
    │
    ├── video/                         # Video module
    │   ├── video_entity.go            # Video model
    │   ├── like_entity.go            # Like model
    │   ├── comment_entity.go          # Comment model
    │   ├── *_handler.go              # HTTP handlers
    │   ├── *_repo.go                 # Data access
    │   └── *_service.go              # Business logic
    │
    ├── feed/                          # Feed module
    │   ├── entity.go                 # Feed models
    │   ├── handler.go
    │   ├── repo.go
    │   └── service.go
    │
    ├── social/                        # Social module
    │   ├── entity.go
    │   ├── handler.go
    │   ├── repo.go
    │   └── service.go
    │
    ├── middleware/                    # Middleware
    │   ├── jwt/                      # JWT authentication
    │   ├── redis/                    # Redis cache client
    │   └── rabbitmq/                 # RabbitMQ client
    │
    ├── worker/                        # Background workers
    │   ├── likeworker.go             # Like processor
    │   ├── commentworker.go          # Comment processor
    │   ├── socialworker.go            # Social processor
    │   └── popularityworker.go       # Popularity calculator
    │
    ├── auth/                          # Authentication
    ├── config/                        # Configuration loader
    ├── db/                            # Database connection
    └── observability/                  # Monitoring / PProf
```

## 🔧 Core Modules

### Account Module
| Component | Description |
|-----------|-------------|
| `entity.go` | Account data model with validation |
| `handler.go` | Register, Login, Logout, Password change |
| `service.go` | Business logic: token generation, validation |
| `repo.go` | MySQL operations via GORM |

### Video Module
| Component | Description |
|-----------|-------------|
| `video_entity.go` | Video metadata, streaming URLs |
| `like_entity.go` | Like records with unique constraint |
| `comment_entity.go` | Comment with threading support |
| `*_service.go` | Upload, publish, delete, count update |
| `*_repo.go` | Batch operations, pagination |

### Feed Module
| Component | Description |
|-----------|-------------|
| `entity.go` | Feed response models |
| `service.go` | Multi-level caching strategy |
| `repo.go` | SQL optimization for feeds |

### Social Module
| Component | Description |
|-----------|-------------|
| `entity.go` | Follower/Following relationships |
| `service.go` | Follow/unfollow with timeline updates |
| `repo.go` | Bidirectional query optimization |

## ⚙️ Configuration

### Local Development (`configs/config.yaml`)

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 3306
  user: root
  password: "123456"
  dbname: Lik_tok

redis:
  host: localhost
  port: 6379
  password: "123456"
  db: 0

rabbitmq:
  host: localhost
  port: 5672
  username: guest
  password: guest
```

### Docker Deployment (`configs/config.docker.yaml`)

```yaml
server:
  port: 8080

database:
  host: mysql
  port: 3306
  user: root
  password: "123456"
  dbname: Lik_tok

redis:
  host: redis
  port: 6379
  password: "123456"
  db: 0

rabbitmq:
  host: rabbitmq
  port: 5672
  username: guest
  password: guest
```

## 🗄️ Database Schema

### Core Tables

```sql
-- User Account
CREATE TABLE accounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Video
CREATE TABLE videos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    author_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    play_url VARCHAR(500) NOT NULL,
    cover_url VARCHAR(500) NOT NULL,
    likes_count INT DEFAULT 0,
    comments_count INT DEFAULT 0,
    popularity INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_author_created (author_id, created_at)
);

-- Like (Unique constraint prevents duplicate likes)
CREATE TABLE likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY idx_like_video_account (video_id, account_id)
);

-- Comment
CREATE TABLE comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_video_created (video_id, created_at)
);

-- Social (Follow relationship)
CREATE TABLE socials (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    follower_id BIGINT NOT NULL,
    followed_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY idx_social (follower_id, followed_id)
);
```

## 🔌 API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/account/register` | Create new account |
| POST | `/account/login` | Authenticate user |
| POST | `/account/changePassword` | Update password |

### Video
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/video/publish` | Upload video |
| POST | `/video/listByAuthorID` | Get author's videos |
| POST | `/video/getDetail` | Get video details |
| POST | `/video/delete` | Delete video |

### Interaction
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/like/like` | Like video |
| POST | `/like/unlike` | Unlike video |
| POST | `/like/isLiked` | Check like status |
| POST | `/comment/publish` | Add comment |
| POST | `/comment/list` | Get comments |

### Feed
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/feed/listLatest` | Latest videos |
| POST | `/feed/listByPopularity` | Popular videos |
| POST | `/feed/listByFollowing` | Following feed |

### Social
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/social/follow` | Follow user |
| POST | `/social/unfollow` | Unfollow user |
| POST | `/social/getAllFollowers` | Get followers |
| POST | `/social/getAllVloggers` | Get following |

## 📊 Caching Strategy

```
Request ──► L1 Cache ──► L2 Cache ──► Database
           (Local)      (Redis)       (MySQL)
             │             │             │
         5 seconds      1 hour      Permanent
```

### Cache Keys
| Key Pattern | TTL | Description |
|------------|-----|-------------|
| `video:entity:{id}` | 1 hour | Video details |
| `feed:latest` | 5 seconds | Latest feed |
| `feed:popular` | 10 minutes | Popular feed |

## 📬 Message Queues

### Queue Architecture

| Queue | Consumer | Purpose |
|-------|----------|---------|
| `like.events` | LikeWorker | Handle like/unlike async |
| `comment.events` | CommentWorker | Comment notifications |
| `social.events` | SocialWorker | Follow timeline updates |
| `video.popularity.events` | PopularityWorker | Recalculate popularity |

### Message Format

```json
{
  "action": "like",
  "user_id": 123,
  "video_id": 456,
  "timestamp": "2026-04-01T10:00:00Z"
}
```

## 🚀 Development

### Run API Server

```bash
cd backend
go run cmd/main.go
```

### Run Worker

```bash
cd backend
go run cmd/worker/main.go
```

### Docker Development

```bash
cd ../deploy/docker
docker-compose up -d backend worker
```

### Add New Module

```bash
# 1. Create module directory
mkdir internal/newmodule

# 2. Implement components
touch internal/newmodule/entity.go
touch internal/newmodule/repo.go
touch internal/newmodule/service.go
touch internal/newmodule/handler.go

# 3. Register routes
# Edit internal/http/router.go
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific module
go test ./internal/video/...
```

## 📈 Performance Optimizations

| Optimization | Implementation |
|-------------|-----------------|
| Connection Pooling | GORM connection pool |
| Batch Operations | Bulk insert/update |
| Index Optimization | Composite indexes |
| Query Caching | Redis L2 cache |
| Async Processing | RabbitMQ workers |
