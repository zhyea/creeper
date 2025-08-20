#!/bin/bash

# Creeper 工具构建脚本

echo "🔨 构建 Creeper 工具..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误：未找到 Go 编译器"
    echo "请先安装 Go 语言环境"
    exit 1
fi

echo "✅ 检测到 Go 版本: $(go version)"

# 构建主程序
echo "📦 构建主程序..."
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

echo "✅ 构建完成！"
echo ""
echo "🎉 可用工具："
echo "  ./creeper                    # 主程序 - 生成静态站点"
echo "  ./creeper -serve             # 生成并启动本地服务器"
echo "  ./cover-gen                  # 封面生成器"
echo ""
echo "📚 封面生成器使用示例："
echo "  ./cover-gen -title \"我的小说\" -theme fantasy"
echo "  ./cover-gen -title \"科幻故事\" -theme scifi -subtitle \"未来世界\""
echo "  ./cover-gen -list-themes     # 查看所有可用主题"
