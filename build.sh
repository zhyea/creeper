#!/bin/bash

# Creeper 静态小说站点生成器构建脚本

echo "🚀 开始构建 Creeper..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误：未找到 Go 编译器"
    echo "请先安装 Go 语言环境："
    echo "  - 访问 https://golang.org/dl/ 下载安装包"
    echo "  - 或使用包管理器安装："
    echo "    macOS: brew install go"
    echo "    Ubuntu: sudo apt install golang-go"
    echo "    CentOS: sudo yum install golang"
    exit 1
fi

echo "✅ 检测到 Go 版本: $(go version)"

# 下载依赖
echo "📦 下载依赖..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ 依赖下载失败"
    exit 1
fi

# 构建可执行文件
echo "🔨 构建可执行文件..."
go build -o creeper main.go

if [ $? -ne 0 ]; then
    echo "❌ 构建失败"
    exit 1
fi

echo "✅ 构建完成！"
echo ""
echo "🎉 使用方法："
echo "  ./creeper                    # 生成静态站点"
echo "  ./creeper -serve             # 生成并启动本地服务器"
echo "  ./creeper -serve -port 3000  # 在指定端口启动服务器"
echo ""
echo "📚 示例小说文件已准备好，位于 novels/ 目录"
echo "🌐 生成的站点将保存在 dist/ 目录"
