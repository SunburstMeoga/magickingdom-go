# MagicKingdom Go

一个基于 Go 的 RESTful API 服务，为微信小程序提供用户管理功能。

## 📚 文档

- [开发指南](docs/DEVELOPMENT.md) - 详细的开发文档和最佳实践
- [API 测试](docs/API_TESTING.md) - API 接口测试示例
- [部署指南](docs/DEPLOYMENT.md) - 生产环境部署说明

## 🚀 技术栈

- **Web 框架**: Gin
- **配置管理**: Viper
- **数据库**: MySQL + GORM
- **日志**: Logrus
- **认证**: JWT (golang-jwt/jwt)
- **数据库迁移**: golang-migrate/migrate
- **容器化**: Docker

## 项目结构

```
magickingdom-go/
├── cmd/                    # 命令行工具（可选）
├── configs/                # 配置文件
│   └── config.yaml        # 开发环境配置
├── internal/              # 内部代码
│   ├── config/           # 配置加载
│   ├── database/         # 数据库初始化
│   ├── dto/              # 数据传输对象
│   ├── handler/          # HTTP 处理器
│   ├── logger/           # 日志工具
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── repository/       # 数据访问层
│   ├── response/         # 统一响应格式
│   ├── router/           # 路由配置
│   ├── service/          # 业务逻辑层
│   └── utils/            # 工具函数
├── migrations/            # 数据库迁移文件
├── .air.toml             # Air 热重载配置
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── main.go               # 程序入口
├── Makefile              # 构建脚本
└── README.md
```

## 快速开始

### 前置要求

- Go 1.24+
- MySQL 5.7+
- (可选) Docker
- (可选) golang-migrate CLI

### 安装依赖

```bash
make deps
```

或手动执行：

```bash
go mod download
go mod tidy
```

### 配置

1. 复制配置文件并修改：

```bash
cp configs/config.yaml configs/config.local.yaml
```

2. 修改 `configs/config.local.yaml` 中的配置：
   - 数据库连接信息
   - 微信小程序 AppID 和 AppSecret
   - JWT Secret

### 数据库迁移

使用 golang-migrate 工具执行迁移：

```bash
# 安装 migrate 工具
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 执行迁移
make migrate-up
```

或者直接运行程序，GORM 会自动创建表（仅开发环境）。

### 运行

#### 开发模式（热重载）

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 运行
make dev
```

#### 普通模式

```bash
make run
```

或：

```bash
go run main.go
```

#### 编译运行

```bash
make build
./bin/magickingdom-go
```

### Docker 运行

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

或使用 docker-compose（需要创建 docker-compose.yml）。

## API 文档

### 认证相关

#### 微信小程序登录

```
POST /api/v1/auth/wechat/login
```

请求体：

```json
{
  "code": "微信登录凭证"
}
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "JWT Token",
    "user": {
      "id": 1,
      "open_id": "xxx",
      "nickname": "用户昵称",
      ...
    }
  }
}
```

### 用户相关（需要 JWT 认证）

所有用户相关接口需要在 Header 中携带 Token：

```
Authorization: Bearer <token>
```

#### 获取用户信息

```
GET /api/v1/user/info
```

#### 更新用户信息

```
PUT /api/v1/user/info
```

请求体：

```json
{
  "nickname": "新昵称",
  "avatar_url": "头像URL",
  "gender": 1,
  "phone": "手机号",
  "email": "邮箱"
}
```

## 开发指南

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释
- 错误处理使用 error wrapping

### 项目架构

项目采用分层架构：

1. **Handler 层**: 处理 HTTP 请求和响应
2. **Service 层**: 业务逻辑处理
3. **Repository 层**: 数据访问
4. **Model 层**: 数据模型定义

使用手动依赖注入，在 `main.go` 中组装各层依赖。

### 添加新功能

1. 在 `models` 中定义数据模型
2. 在 `repository` 中实现数据访问接口
3. 在 `service` 中实现业务逻辑
4. 在 `handler` 中实现 HTTP 处理
5. 在 `router` 中注册路由
6. 在 `main.go` 中注入依赖

## 测试

```bash
# 运行所有测试
make test

# 生成测试覆盖率报告
make test-coverage
```

## 常用命令

```bash
make help          # 查看所有可用命令
make build         # 编译项目
make run           # 运行项目
make test          # 运行测试
make clean         # 清理编译文件
make deps          # 安装依赖
make dev           # 开发模式运行
```

## 环境变量

可以通过环境变量覆盖配置文件：

- `CONFIG_PATH`: 配置文件路径（默认: configs/config.yaml）

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
