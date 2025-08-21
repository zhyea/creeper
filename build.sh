#!/bin/bash

# Creeper 构建脚本

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

# 构建主程序
echo "🔨 构建主程序..."
go build -o creeper main.go
if [ $? -ne 0 ]; then
    echo "❌ 主程序构建失败"
    exit 1
fi

# 构建封面生成器
echo "🎨 构建封面生成器..."
go build -o cover-gen cmd/cover/main.go
if [ $? -ne 0 ]; then
    echo "❌ 封面生成器构建失败"
    exit 1
fi

# 构建部署工具
echo "🚀 构建部署工具..."
go build -o deploy-tool cmd/deploy/main.go
if [ $? -ne 0 ]; then
    echo "❌ 部署工具构建失败"
    exit 1
fi

echo "✅ 构建完成！"
echo ""
echo "🎉 可用工具："
echo "  ./creeper                    # 主程序 - 生成静态站点"
echo "  ./creeper -serve             # 生成并启动本地服务器"
echo "  ./creeper -serve -port 3000  # 在指定端口启动服务器"
echo "  ./creeper -deploy            # 生成并自动部署"
echo "  ./cover-gen                  # 封面生成器"
echo "  ./deploy-tool                # 部署工具"
echo ""
echo "📚 封面生成器使用示例："
echo "  ./cover-gen -title \"我的小说\" -theme fantasy"
echo "  ./cover-gen -title \"科幻故事\" -theme scifi -subtitle \"未来世界\""
echo "  ./cover-gen -list-themes     # 查看所有可用主题"
echo ""
echo "🚀 部署工具使用示例："
echo "  ./deploy-tool -init -type cloudflare  # 初始化 Cloudflare 配置"
echo "  ./deploy-tool -config deploy-config.yaml -site dist  # 部署站点"
echo "  ./deploy-tool -status                  # 查看部署状态"
echo "  ./deploy-tool -list                    # 查看部署历史"
echo ""
echo "📚 示例小说文件已准备好，位于 novels/ 目录"
echo "🌐 生成的站点将保存在 dist/ 目录"
