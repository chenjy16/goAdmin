# Admin 管理系统

一个基于 Go 语言开发的现代化后端管理系统，采用清洁架构设计，提供用户管理、身份验证等核心功能。

## ✨ 特性

- 🚀 **高性能**: 基于 Gin 框架，提供高性能的 HTTP 服务
- 🏗️ **清洁架构**: 采用分层架构设计，代码结构清晰，易于维护
- 🔐 **安全认证**: 集成 JWT 身份验证机制
- 📊 **结构化日志**: 使用 Zap 提供结构化日志记录
- 🗄️ **数据库支持**: 支持 SQLite 数据库，使用 SQLC 生成类型安全的数据库操作代码
- ⚡ **依赖注入**: 使用 Google Wire 进行依赖注入管理
- ✅ **数据验证**: 集成强大的数据验证功能
- 🔧 **配置管理**: 使用 Viper 进行灵活的配置管理

## 🛠️ 技术栈

### 核心框架
- **Go 1.24.0** - 编程语言
- **Gin v1.11.0** - HTTP Web 框架
- **SQLite3** - 轻量级数据库

### 主要依赖
- **SQLC** - 类型安全的 SQL 代码生成器
- **Google Wire v0.7.0** - 依赖注入框架
- **Zap v1.27.0** - 结构化日志库
- **Viper v1.17.0** - 配置管理
- **JWT v5.3.0** - JSON Web Token 认证
- **Validator v10.28.0** - 数据验证
- **Crypto** - 密码加密

### 开发工具
- **Air** - 热重载开发工具（推荐）
- **Wire** - 依赖注入代码生成

## 📁 项目结构

```
admin/
├── cmd/                    # 应用程序入口
│   └── main.go
├── internal/               # 内部包（不对外暴露）
│   ├── config/            # 配置管理
│   ├── controllers/       # 控制器层
│   ├── database/          # 数据库相关
│   │   ├── connection.go  # 数据库连接
│   │   ├── curd/         # SQL 查询文件
│   │   └── generated/    # SQLC 生成的代码
│   ├── dto/              # 数据传输对象
│   ├── errors/           # 错误处理
│   ├── middleware/       # 中间件
│   ├── repository/       # 数据访问层
│   ├── response/         # 响应格式化
│   ├── route/           # 路由配置
│   ├── services/        # 业务逻辑层
│   ├── utils/           # 工具函数
│   └── wire/            # 依赖注入配置
├── schemas/             # 数据库模式文件
│   └── users/
├── config.yaml         # 配置文件
├── sqlc.yaml          # SQLC 配置
└── go.mod             # Go 模块文件
```

## 🚀 快速开始

### 环境要求

- Go 1.24.0 或更高版本
- SQLite3

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd admin
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **安装开发工具**
   ```bash
   # 安装 SQLC（用于生成数据库代码）
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   
   # 安装 Wire（用于依赖注入）
   go install github.com/google/wire/cmd/wire@latest
   
   # 安装 Air（可选，用于热重载开发）
   go install github.com/air-verse/air@latest
   ```

4. **生成代码**
   ```bash
   # 生成数据库访问代码
   sqlc generate
   
   # 生成依赖注入代码
   cd internal/wire && wire
   ```

5. **初始化数据库**
   ```bash
   # 创建数据目录
   mkdir -p data
   
   # 初始化数据库
   sqlite3 data/admin.db < schemas/users/001_create_users_table.sql
   ```

6. **配置应用**
   
   复制并修改配置文件：
   ```bash
   cp config.yaml config.local.yaml
   ```
   
   编辑 `config.local.yaml` 根据需要修改配置：
   ```yaml
   server:
     host: "localhost"
     port: "8080"
     mode: "debug"  # debug, release, test
   
   database:
     driver: "sqlite3"
     dsn: "./data/admin.db"
   
   jwt:
     secret: "your-secret-key-change-this-in-production"
     expire_time: 24  # hours
   ```

7. **运行应用**
   ```bash
   # 开发模式（推荐使用 Air 进行热重载）
   air
   
   # 或者直接运行
   go run cmd/main.go
   ```

应用将在 `http://localhost:8080` 启动。

### 验证安装

访问健康检查端点：
```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

## 📚 API 文档

### 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`

### 用户管理 API

#### 1. 创建用户
```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User"
}
```

#### 2. 获取用户列表
```http
GET /api/v1/users?page=1&limit=10
```

#### 3. 获取单个用户
```http
GET /api/v1/users/{id}
```

#### 4. 更新用户
```http
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "email": "newemail@example.com",
  "full_name": "New Name",
  "is_active": true
}
```

#### 5. 删除用户
```http
DELETE /api/v1/users/{id}
```

### 响应格式

成功响应：
```json
{
  "code": 200,
  "message": "Success message",
  "data": {
    // 响应数据
  }
}
```

错误响应：
```json
{
  "code": 400,
  "message": "Error message",
  "error": "ERROR_CODE"
}
```

## 🔧 开发指南

### 代码生成

当修改数据库模式或查询时，需要重新生成代码：

```bash
# 重新生成数据库访问代码
sqlc generate

# 重新生成依赖注入代码
cd internal/wire && wire
```

### 添加新的数据库模块

1. 在 `schemas/` 目录下创建新的模块目录
2. 在 `internal/database/curd/` 目录下创建对应的 SQL 查询文件
3. 更新 `sqlc.yaml` 配置文件
4. 运行 `sqlc generate` 生成代码
5. 在 `internal/database/connection.go` 中添加新的查询字段

### 项目架构说明

本项目采用清洁架构（Clean Architecture）设计：

- **Controllers**: 处理 HTTP 请求和响应
- **Services**: 业务逻辑层
- **Repository**: 数据访问层
- **Models/DTO**: 数据传输对象
- **Middleware**: 中间件（认证、日志、错误处理等）

### 错误处理

项目使用统一的错误处理机制：

- 自定义错误类型 `AppError`
- 错误中间件自动处理和格式化错误响应
- 结构化错误日志记录

### 日志记录

使用 Zap 进行结构化日志记录：

- 请求/响应日志
- 错误日志
- 业务操作日志
- 支持不同日志级别

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📦 构建和部署

### 构建

```bash
# 构建二进制文件
go build -o bin/admin cmd/main.go

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o bin/admin-linux cmd/main.go
```

### Docker 部署

创建 `Dockerfile`：
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o admin cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/admin .
COPY --from=builder /app/config.yaml .
CMD ["./admin"]
```

构建和运行：
```bash
docker build -t admin-system .
docker run -p 8080:8080 admin-system
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件至：[your-email@example.com]

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [SQLC](https://github.com/sqlc-dev/sqlc) - SQL 代码生成器
- [Wire](https://github.com/google/wire) - 依赖注入框架
- [Zap](https://github.com/uber-go/zap) - 日志库
- [Viper](https://github.com/spf13/viper) - 配置管理

---

⭐ 如果这个项目对你有帮助，请给它一个星标！