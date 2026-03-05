#!/bin/bash

# 快速启动脚本

set -e

echo "🚀 MagicKingdom Go 快速启动"
echo "=========================="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误：未安装 Go，请先安装 Go 1.24+"
    exit 1
fi

echo "✅ Go 版本: $(go version)"

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo "❌ 错误：配置文件不存在，请先创建 configs/config.yaml"
    exit 1
fi

echo "✅ 配置文件存在"

# 安装依赖
echo ""
echo "📦 安装依赖..."
go mod download
go mod tidy

echo "✅ 依赖安装完成"

# 检查是否需要启动 MySQL
read -p "是否使用 Docker Compose 启动 MySQL? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if ! command -v docker-compose &> /dev/null; then
        echo "❌ 错误：未安装 docker-compose"
    else
        echo "🐳 启动 MySQL..."
        docker-compose up -d mysql
        echo "✅ MySQL 已启动"
        echo "⏳ 等待 MySQL 就绪..."
        sleep 10
    fi
fi

# 编译项目
echo ""
echo "🔨 编译项目..."
go build -o bin/magickingdom-go main.go
echo "✅ 编译完成"

# 运行项目
echo ""
echo "🎉 启动服务..."
echo "=========================="
./bin/magickingdom-go

