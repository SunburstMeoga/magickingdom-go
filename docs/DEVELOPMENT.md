# 开发指南

## 项目架构

本项目采用经典的分层架构，遵循依赖倒置原则：

```
┌─────────────────────────────────────────┐
│           Handler Layer                  │  HTTP 请求处理
│  (接收请求、参数验证、返回响应)           │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│           Service Layer                  │  业务逻辑
│  (业务规则、数据转换、第三方调用)         │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         Repository Layer                 │  数据访问
│  (数据库操作、缓存操作)                   │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│            Database                      │  数据存储
└─────────────────────────────────────────┘
```

## 目录结构说明

```
internal/
├── config/         # 配置管理
│   └── config.go   # 配置结构和加载逻辑
├── database/       # 数据库初始化
│   └── database.go # 数据库连接和配置
├── dto/            # 数据传输对象
│   └── user.go     # 用户相关 DTO
├── handler/        # HTTP 处理器
│   └── user_handler.go
├── logger/         # 日志工具
│   └── logger.go
├── middleware/     # 中间件
│   ├── auth.go     # JWT 认证
│   ├── cors.go     # 跨域处理
│   ├── logger.go   # 日志记录
│   └── recovery.go # 错误恢复
├── models/         # 数据模型
│   └── user.go
├── repository/     # 数据访问层
│   └── user_repository.go
├── response/       # 统一响应格式
│   └── response.go
├── router/         # 路由配置
│   └── router.go
├── service/        # 业务逻辑层
│   └── user_service.go
└── utils/          # 工具函数
    └── jwt.go      # JWT 工具
```

## 添加新功能的步骤

### 示例：添加文章管理功能

#### 1. 定义数据模型 (models/article.go)

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Article struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    UserID  uint   `gorm:"not null;index" json:"user_id"`
    Title   string `gorm:"size:200;not null" json:"title"`
    Content string `gorm:"type:text" json:"content"`
    Status  int    `gorm:"default:1" json:"status"`
    
    // 关联
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

#### 2. 创建 DTO (dto/article.go)

```go
package dto

type CreateArticleRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content" binding:"required"`
}

type ArticleDTO struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Status    int       `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 3. 实现 Repository (repository/article_repository.go)

```go
package repository

import (
    "gorm.io/gorm"
    "magickingdom-go/internal/models"
)

type ArticleRepository interface {
    Create(article *models.Article) error
    FindByID(id uint) (*models.Article, error)
    FindByUserID(userID uint) ([]*models.Article, error)
    Update(article *models.Article) error
    Delete(id uint) error
}

type articleRepository struct {
    db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
    return &articleRepository{db: db}
}

func (r *articleRepository) Create(article *models.Article) error {
    return r.db.Create(article).Error
}

// ... 实现其他方法
```

#### 4. 实现 Service (service/article_service.go)

```go
package service

import (
    "magickingdom-go/internal/dto"
    "magickingdom-go/internal/models"
    "magickingdom-go/internal/repository"
)

type ArticleService interface {
    CreateArticle(userID uint, req *dto.CreateArticleRequest) (*dto.ArticleDTO, error)
    GetArticle(id uint) (*dto.ArticleDTO, error)
    // ... 其他方法
}

type articleService struct {
    articleRepo repository.ArticleRepository
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
    return &articleService{articleRepo: articleRepo}
}

func (s *articleService) CreateArticle(userID uint, req *dto.CreateArticleRequest) (*dto.ArticleDTO, error) {
    article := &models.Article{
        UserID:  userID,
        Title:   req.Title,
        Content: req.Content,
        Status:  1,
    }
    
    if err := s.articleRepo.Create(article); err != nil {
        return nil, err
    }
    
    return s.modelToDTO(article), nil
}

// ... 实现其他方法
```

#### 5. 实现 Handler (handler/article_handler.go)

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "magickingdom-go/internal/dto"
    "magickingdom-go/internal/middleware"
    "magickingdom-go/internal/response"
    "magickingdom-go/internal/service"
)

type ArticleHandler struct {
    articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
    return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    var req dto.CreateArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }
    
    article, err := h.articleService.CreateArticle(userID, &req)
    if err != nil {
        response.Error(c, 500, err.Error())
        return
    }
    
    response.Success(c, article)
}
```

#### 6. 注册路由 (router/router.go)

```go
// 在 SetupRouter 函数中添加
article := v1.Group("/article")
article.Use(middleware.AuthMiddleware(jwtUtil))
{
    article.POST("", articleHandler.CreateArticle)
    article.GET("/:id", articleHandler.GetArticle)
}
```

#### 7. 依赖注入 (main.go)

```go
// 在 main 函数中添加
articleRepo := repository.NewArticleRepository(db)
articleService := service.NewArticleService(articleRepo)
articleHandler := handler.NewArticleHandler(articleService)

// 传递给路由
r := router.SetupRouter(userHandler, articleHandler, jwtUtil)
```

## 代码规范

### 命名规范

- 包名：小写，简短，有意义
- 文件名：小写，下划线分隔
- 接口：名词，如 `UserService`
- 实现：小写开头，如 `userService`
- 方法：驼峰命名，如 `GetUserInfo`

### 错误处理

使用 Go 1.13+ 的 error wrapping：

```go
if err != nil {
    return fmt.Errorf("操作失败: %w", err)
}
```

### 日志记录

```go
logger.GetLogger().WithFields(map[string]interface{}{
    "user_id": userID,
    "action": "login",
}).Info("用户登录")
```

## 测试

### 单元测试

```go
func TestUserService_GetUserInfo(t *testing.T) {
    // 准备
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, nil, nil)
    
    // 执行
    user, err := service.GetUserInfo(1)
    
    // 断言
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/service/...

# 生成覆盖率报告
make test-coverage
```

## 数据库迁移

### 创建迁移文件

```bash
make migrate-create name=create_articles_table
```

### 编写迁移 SQL

up 文件：
```sql
CREATE TABLE articles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ...
);
```

down 文件：
```sql
DROP TABLE IF EXISTS articles;
```

### 执行迁移

```bash
make migrate-up
```

## 常见问题

### 1. 如何添加新的配置项？

在 `internal/config/config.go` 中添加字段，然后在 `configs/config.yaml` 中添加配置值。

### 2. 如何处理事务？

```go
func (s *service) DoSomething() error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 在事务中执行操作
        return nil
    })
}
```

### 3. 如何添加新的中间件？

在 `internal/middleware/` 中创建新文件，然后在 `router.go` 中使用。

