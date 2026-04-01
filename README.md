# Lik_tok - 短视频社交平台

一个基于 Go + Vue3 的短视频社交平台，支持视频上传、播放、点赞、评论、关注等功能。

## 项目架构

```
Lik_tok/
├── backend/          # Go 后端服务
├── frontend/         # Vue3 前端应用
├── deploy/           # Docker 部署配置
├── test/             # 测试文件
├── start-local.sh     # 本地开发启动脚本
└── README.md          # 项目文档
```

## 技术栈

### 后端
- **语言**: Go 1.25
- **框架**: Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **消息队列**: RabbitMQ
- **架构模式**: 分层架构 (Handler → Service → Repository)

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
- Docker 20.10+（Docker 开发）
- Go 1.25+、Node.js 20+（本地开发）
- MySQL 8.0、Redis 7+、RabbitMQ 3.12+

---

## 启动方式

### 方式一：Docker 开发（推荐）

```bash
# 进入部署目录
cd deploy/docker

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down
```

**服务访问地址：**
- 前端: http://localhost
- 后端 API: http://localhost:8080
- RabbitMQ 管理: http://localhost:15672 (guest/guest)

---

### 方式二：本地开发

**前提条件：**
```bash
# 安装基础服务 (macOS)
brew install go mysql redis rabbitmq

# 启动基础服务
brew services start mysql
brew services start redis
brew services start rabbitmq
```

**一键启动：**
```bash
# 启动后端 (API + Worker) 和前端
./start-local.sh

# 停止服务: Ctrl+C
```

**手动启动：**
```bash
# 终端 1: 启动后端 API
cd backend
go run cmd/main.go

# 终端 2: 启动 Worker
cd backend
go run cmd/worker/main.go

# 终端 3: 启动前端
cd frontend
npm install  # 首次需要安装依赖
npm run dev
```

> ⚠️ **注意**: Docker 和本地开发不能同时运行，会出现端口冲突。请根据需要选择一种方式。

---

## 配置说明

### 统一密码配置

MySQL、Redis、RabbitMQ 使用统一密码：

| 服务 | 用户名 | 密码 |
|------|--------|------|
| MySQL | root | 123456 |
| Redis | - | 123456 |
| RabbitMQ | guest | guest |

### 配置文件位置

- **Docker 配置**: `backend/configs/config.docker.yaml`
- **本地配置**: `backend/configs/config.yaml`

### 本地配置示例

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

---

## 项目结构详解

### Backend
```
backend/
├── cmd/
│   ├── main.go              # API 服务入口
│   └── worker/
│       └── main.go          # 后台工作进程入口
├── configs/
│   ├── config.yaml          # 本地开发配置
│   └── config.docker.yaml   # Docker 配置
├── internal/
│   ├── account/             # 用户账户模块
│   ├── video/               # 视频、点赞、评论管理
│   ├── feed/                # 推荐 Feed 生成
│   ├── social/              # 关注关系管理
│   ├── middleware/          # 中间件
│   │   ├── jwt/             # JWT 认证
│   │   ├── rabbitmq/        # RabbitMQ 客户端
│   │   └── redis/           # Redis 客户端
│   └── worker/              # 后台任务处理器
└── go.mod
```

### Frontend
```
frontend/
├── src/
│   ├── api/                 # API 接口封装
│   ├── components/          # 公共组件
│   ├── views/               # 页面视图
│   ├── stores/              # Pinia 状态管理
│   └── router/              # 路由配置
└── package.json
```

---

## 数据库设计

### 核心表

| 表名 | 说明 |
|------|------|
| accounts | 用户账户 |
| videos | 视频信息 |
| likes | 点赞记录 |
| comments | 评论记录 |
| socials | 关注关系 |

---

## 缓存策略

采用多级缓存架构：
- **L1**: 本地缓存 (5秒)
- **L2**: Redis 缓存 (1小时)
- **L3**: MySQL 数据库

---

## 消息队列

使用 RabbitMQ 处理异步任务：

| 队列 | 说明 |
|------|------|
| like.events | 点赞/取消点赞 |
| comment.events | 评论 |
| social.events | 关注关系 |
| video.popularity.events | 视频热度计算 |

---

## 常用命令

### Docker 相关
```bash
# 进入 Docker 目录
cd deploy/docker

# 重新构建并启动
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
docker-compose logs -f worker

# 重启特定服务
docker-compose restart backend

# 清理并重建
docker-compose down -v
docker-compose up -d
```

### 数据库操作
```bash
# 连接 MySQL (本地)
mysql -u root -p123456

# 备份数据库
mysqldump -u root -p123456 Lik_tok > backup.sql

# 清理 Redis 缓存 (Docker)
docker exec lik_tok_redis redis-cli -a 123456 FLUSHALL
```

---

## API 测试

使用 Postman 导入 `test/postman.json` 进行 API 测试。

---

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

---

## 许可证

[MIT](LICENSE)
