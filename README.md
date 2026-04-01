# Lik_tok - 短视频社交平台

一个基于 Go + Vue3 的短视频社交平台，支持视频上传、播放、点赞、评论、关注等功能。

## 项目架构

```
Lik_tok/
├── backend/          # Go 后端服务
├── frontend/         # Vue3 前端应用
├── deploy/           # Docker 部署配置
└── test/             # 测试文件
```

## 技术栈

### 后端
- **语言**: Go 1.25
- **框架**: Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **消息队列**: RabbitMQ
- **架构模式**: 分层架构 (Handler -> Service -> Repository)

### 前端
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件**: 自定义组件
- **状态管理**: Pinia
- **路由**: Vue Router

### 部署
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx

## 核心功能

- 📹 **视频管理**: 上传、播放、删除视频
- 👍 **点赞系统**: 点赞/取消点赞，实时更新点赞数
- 💬 **评论系统**: 发表评论、查看评论列表
- 👥 **社交功能**: 关注/取消关注用户
- 🏠 **推荐 Feed**: 基于时间线的视频推荐
- 🔥 **热门视频**: 基于热度的视频排序
- 🔐 **用户认证**: JWT 认证机制

## 快速开始

### 环境要求
- Docker 20.10+
- Docker Compose 2.0+

### 启动项目

```bash
# 进入部署目录
cd deploy/docker

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps
```

服务启动后访问:
- 前端: http://localhost
- 后端 API: http://localhost:8080

### 停止项目

```bash
docker-compose down
```

## 项目结构详解

### Backend
- `cmd/` - 应用程序入口
  - `main.go` - API 服务
  - `worker/main.go` - 后台工作进程
- `internal/` - 内部包
  - `account/` - 用户账户管理
  - `video/` - 视频、点赞、评论管理
  - `feed/` - 推荐 Feed 生成
  - `social/` - 关注关系管理
  - `middleware/` - JWT、Redis、RabbitMQ 中间件
  - `worker/` - 后台任务处理

### Frontend
- `src/api/` - API 接口封装
- `src/components/` - 公共组件
- `src/views/` - 页面视图
- `src/stores/` - Pinia 状态管理
- `src/router/` - 路由配置

## 数据库设计

### 核心表
- `accounts` - 用户账户
- `videos` - 视频信息
- `likes` - 点赞记录
- `comments` - 评论记录
- `socials` - 关注关系

## 缓存策略

采用多级缓存架构:
- **L1**: 本地缓存 (5秒)
- **L2**: Redis 缓存 (1小时)
- **L3**: MySQL 数据库

## 消息队列

使用 RabbitMQ 处理异步任务:
- `like.events` - 点赞/取消点赞
- `comment.events` - 评论
- `social.events` - 关注关系
- `video.popularity.events` - 视频热度计算

## 开发指南

### Docker 开发（推荐）

```bash
# 启动所有服务
cd deploy/docker
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 本地开发

#### 1. 安装依赖

```bash
# macOS
brew install go mysql redis rabbitmq

# 启动基础服务
brew services start mysql
brew services start redis
brew services start rabbitmq
```

#### 2. 创建数据库

```bash
mysql -u root -e "CREATE DATABASE IF NOT EXISTS Lik_tok CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

#### 3. 启动后端

**重要**: 本地开发需要同时启动 API 服务和 Worker 服务

```bash
# 使用启动脚本（推荐）
./start-local.sh

# 或者手动启动
cd backend

# 终端 1: 启动 API 服务
go run cmd/main.go

# 终端 2: 启动 Worker 服务
go run cmd/worker/main.go
```

#### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

### 本地配置

本地配置文件位于 `backend/configs/config.yaml`，可根据本地环境修改：

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 3306
  user: root
  password: ""
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

### API 测试

使用 Postman 导入 `test/postman.json` 进行 API 测试。

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

[MIT](LICENSE)
