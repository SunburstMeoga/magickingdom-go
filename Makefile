.PHONY: help build run test clean docker-build docker-run migrate-up migrate-down migrate-create

help: ## 显示帮助信息
	@echo "可用命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 编译项目
	@echo "编译项目..."
	go build -o bin/magickingdom-go main.go

run: ## 运行项目
	@echo "运行项目..."
	go run main.go

test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "生成测试覆盖率报告..."
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

clean: ## 清理编译文件
	@echo "清理编译文件..."
	rm -rf bin/
	rm -f coverage.txt coverage.html

deps: ## 安装依赖
	@echo "安装依赖..."
	go mod download
	go mod tidy

docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t magickingdom-go:latest .

docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	docker run -p 8080:8080 --env-file .env magickingdom-go:latest

migrate-up: ## 执行数据库迁移（升级）
	@echo "执行数据库迁移..."
	migrate -path migrations -database "mysql://root:your_password@tcp(localhost:3306)/magickingdom?charset=utf8mb4&parseTime=True&loc=Local" up

migrate-down: ## 回滚数据库迁移
	@echo "回滚数据库迁移..."
	migrate -path migrations -database "mysql://root:your_password@tcp(localhost:3306)/magickingdom?charset=utf8mb4&parseTime=True&loc=Local" down

migrate-create: ## 创建新的迁移文件 (使用: make migrate-create name=create_users_table)
	@echo "创建迁移文件..."
	migrate create -ext sql -dir migrations -seq $(name)

dev: ## 开发模式运行（使用 air 热重载）
	@echo "开发模式运行..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air 未安装，请运行: go install github.com/cosmtrek/air@latest"; \
		echo "或直接运行: make run"; \
	fi

