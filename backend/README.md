# Lik_tok Backend

Go 实现的后端 API 服务，采用分层架构设计。

## 架构设计

```
┌─────────────┐
│   Handler   │  ← HTTP 请求处理、参数校验
├─────────────┤
│   Service   │  ← 业务逻辑处理
├─────────────┤
│ Repository  │  ← 数据访问层
├─────────────┤
│  MySQL/Redis│  ← 数据存储
└─────────────┘
```

## 目录结构

```
backend/
├── cmd/
│   ├── main.go              # API 服务入口
│   └── worker/
│       └── main.go          # 工作进程入口
├── configs/
│   └── config.docker.yaml   # Docker 环境配置
├── internal/
│   ├── account/             # 用户账户模块
│   │   ├── entity.go        # 实体定义
│   │   ├── handler.go       # HTTP 处理器
│   │   ├── repo.go          # 数据访问
│   │   └── service.go       # 业务逻辑
│   ├── auth/                # JWT 认证
│   ├── config/              # 配置管理
│   ├── db/                  # 数据库连接
│   ├── feed/                # Feed 推荐
│   ├── http/                # HTTP 路由
│   ├── middleware/          # 中间件
│   │   ├── jwt/             # JWT 认证
│   │   ├── rabbitmq/        # RabbitMQ 客户端
│   │   └── redis/           # Redis 客户端
│   ├── observability/       # 监控和性能分析
│   ├── social/              # 社交关系
│   ├── video/               # 视频、点赞、评论
│   └── worker/              # 后台任务处理器
├── go.mod
└── go.sum
```

## 核心模块

### Account (用户账户)
- 用户注册/登录
- JWT Token 生成与验证
- 密码修改

### Video (视频)
- 视频上传 (支持流式上传)
- 视频信息查询
- 视频删除
- 作者视频列表

### Like (点赞)
- 点赞/取消点赞
- 查询是否已点赞
- 查询点赞列表
- 实时更新点赞数

### Comment (评论)
- 发表评论
- 评论列表查询
- 评论数统计

### Feed (推荐)
- 基于时间线的视频推荐
- 热门视频排序
- 多级缓存架构 (L1/L2/L3)

### Social (社交)
- 关注/取消关注
- 粉丝列表
- 关注列表

## 中间件

### JWT 认证
```go
// 使用示例
r.Use(jwt.JWT())
```

### Redis 缓存
```go
// 缓存操作
cache.Set(ctx, key, value, expiration)
cache.Get(ctx, key)
cache.Del(ctx, key)
```

### RabbitMQ
```go
// 发送消息
mq.Publish(exchange, routingKey, body)

// 消费消息
mq.Consume(queue, handler)
```

## 数据库设计

### 核心表结构

```sql
-- 用户表
CREATE TABLE accounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 视频表
CREATE TABLE videos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    author_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    play_url VARCHAR(500) NOT NULL,
    cover_url VARCHAR(500) NOT NULL,
    likes_count INT DEFAULT 0,
    comments_count INT DEFAULT 0,
    popularity INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 点赞表
CREATE TABLE likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY idx_like_video_account (video_id, account_id)
);

-- 评论表
CREATE TABLE comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 社交关系表
CREATE TABLE socials (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    follower_id BIGINT NOT NULL,
    followed_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY idx_social (follower_id, followed_id)
);
```

## API 接口

### 账户相关
- `POST /account/register` - 注册
- `POST /account/login` - 登录
- `POST /account/info` - 获取用户信息
- `POST /account/changePassword` - 修改密码

### 视频相关
- `POST /video/publish` - 发布视频
- `POST /video/listByAuthorID` - 获取作者视频列表
- `POST /video/delete` - 删除视频

### 点赞相关
- `POST /like/like` - 点赞
- `POST /like/unlike` - 取消点赞
- `POST /like/isLiked` - 查询是否已点赞
- `POST /like/listLikedVideos` - 获取点赞视频列表

### 评论相关
- `POST /comment/comment` - 发表评论
- `POST /comment/list` - 获取评论列表

### Feed 相关
- `POST /feed/listLatest` - 获取最新视频
- `POST /feed/listHot` - 获取热门视频

### 社交相关
- `POST /social/follow` - 关注用户
- `POST /social/unfollow` - 取消关注
- `POST /social/isFollowing` - 查询是否已关注
- `POST /social/getAllFollowers` - 获取粉丝列表
- `POST /social/getAllVloggers` - 获取关注列表

## 开发指南

### 本地开发

```bash
# 安装依赖
go mod download

# 运行服务
go run cmd/main.go

# 运行工作进程
go run cmd/worker/main.go
```

### 配置说明

配置文件位于 `configs/config.docker.yaml`:

```yaml
server:
  port: 8080
  mode: release

database:
  host: mysql
  port: 3306
  user: root
  password: 123456
  dbname: Lik_tok

redis:
  host: redis
  port: 6379
  password: 123456

rabbitmq:
  host: rabbitmq
  port: 5672
  user: guest
  password: guest

jwt:
  secret: your-secret-key
  expiration: 720h
```

### 添加新模块

1. 在 `internal/` 下创建新目录
2. 创建 `entity.go` 定义数据结构
3. 创建 `repo.go` 实现数据访问
4. 创建 `service.go` 实现业务逻辑
5. 创建 `handler.go` 实现 HTTP 接口
6. 在 `http/router.go` 注册路由

## 性能优化

### 缓存策略
- 视频详情: Redis 缓存 1 小时
- Feed 列表: 多级缓存 (本地 5s + Redis 1h)
- 热门视频: 本地缓存 10 分钟

### 数据库优化
- 点赞表: 复合索引 `(video_id, account_id)`
- 社交表: 复合索引 `(follower_id, followed_id)`
- 视频表: 索引 `(author_id, created_at)`

### 异步处理
- 点赞计数: RabbitMQ 异步更新
- 评论计数: RabbitMQ 异步更新
- 热度计算: 定时任务 + RabbitMQ

## 测试

```bash
# 运行单元测试
go test ./...

# 运行基准测试
go test -bench=. ./...
```
