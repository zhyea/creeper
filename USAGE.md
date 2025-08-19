# Creeper 使用指南

## 🚀 快速开始

### 1. 环境要求

- Go 1.21 或更高版本
- 支持的操作系统：Windows、macOS、Linux

### 2. 安装 Go（如果尚未安装）

#### macOS
```bash
# 使用 Homebrew
brew install go

# 或下载安装包
# 访问 https://golang.org/dl/
```

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install golang-go
```

#### CentOS/RHEL
```bash
sudo yum install golang
```

#### Windows
访问 https://golang.org/dl/ 下载安装包

### 3. 构建项目

```bash
# 使用构建脚本（推荐）
./build.sh

# 或手动构建
go mod tidy
go build -o creeper main.go
```

### 4. 运行程序

```bash
# 生成静态站点
./creeper

# 生成并启动本地服务器
./creeper -serve

# 在指定端口启动服务器
./creeper -serve -port 3000
```

## 📚 创建你的第一部小说

### 方式一：单文件模式

1. 在 `novels/` 目录下创建 `我的小说.md`：

```markdown
---
title: 我的第一部小说
author: 我的名字
description: 这是我的第一部小说
---

# 第一卷 开始

## 第一章 序幕

这里写第一章的内容...

## 第二章 起航

这里写第二章的内容...

# 第二卷 冒险

## 第三章 挑战

这里写第三章的内容...
```

### 方式二：多文件模式

1. 创建小说目录：`novels/我的小说/`
2. 创建元数据文件 `meta.md`：

```markdown
---
title: 我的第一部小说
author: 我的名字
description: 这是我的第一部小说
---
```

3. 创建章节文件：
   - `01-第一章.md`
   - `02-第二章.md`
   - `03-第三章.md`

### 方式三：多卷多文件模式

```
我的小说/
├── meta.md
├── 第一卷-开始/
│   ├── 01-序幕.md
│   ├── 02-起航.md
│   └── 03-初试.md
└── 第二卷-冒险/
    ├── 01-挑战.md
    ├── 02-成长.md
    └── 03-收获.md
```

## ⚙️ 配置文件

修改 `config.yaml` 来自定义你的站点：

```yaml
site:
  title: "我的小说站点"
  description: "欢迎来到我的小说世界"
  author: "我的名字"
  base_url: "/"

theme:
  primary_color: "#2c3e50"    # 主色调
  secondary_color: "#3498db"  # 辅色调
  font_family: "'PingFang SC', 'Microsoft YaHei', sans-serif"
```

## 🎨 自定义主题

你可以通过修改配置文件来自定义站点外观：

### 颜色主题
- `primary_color`: 导航栏、按钮的主色调
- `secondary_color`: 链接、强调内容的辅色调
- `background_color`: 页面背景色
- `text_color`: 正文字体颜色

### 字体设置
- `font_family`: 字体族
- `font_size`: 基础字体大小
- `line_height`: 行高

### 预设主题

#### 深蓝主题
```yaml
theme:
  primary_color: "#1e3a8a"
  secondary_color: "#3b82f6"
  background_color: "#f8fafc"
  text_color: "#1f2937"
```

#### 暖色主题
```yaml
theme:
  primary_color: "#dc2626"
  secondary_color: "#f59e0b"
  background_color: "#fefefe"
  text_color: "#374151"
```

#### 森林主题
```yaml
theme:
  primary_color: "#166534"
  secondary_color: "#059669"
  background_color: "#f0fdf4"
  text_color: "#1f2937"
```

## 📱 功能特性

### 搜索功能
- 实时搜索小说标题、作者、章节
- 支持中文搜索
- 键盘导航支持

### 阅读体验
- 响应式设计，支持手机、平板、桌面
- 阅读进度条
- 键盘快捷键：
  - `Ctrl + ←`: 上一章
  - `Ctrl + →`: 下一章
  - `Ctrl + ↑`: 返回目录
  - `Esc`: 关闭搜索框

### 导航功能
- 自动生成小说目录
- 章节间快速跳转
- 面包屑导航

## 🚀 部署到服务器

### 1. 生成静态文件
```bash
./creeper
```

### 2. 上传 dist 目录
将生成的 `dist/` 目录上传到你的服务器。

### 3. 配置 Web 服务器

#### Nginx 配置示例
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

#### Apache 配置示例
```apache
<VirtualHost *:80>
    ServerName your-domain.com
    DocumentRoot /path/to/dist
    
    <Directory /path/to/dist>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>
```

## 🔧 故障排除

### 常见问题

#### Q: 生成的页面显示乱码
A: 确保你的 Markdown 文件使用 UTF-8 编码保存。

#### Q: 章节顺序不对
A: 检查文件名是否按照 `01-`, `02-` 的格式命名。

#### Q: 搜索功能不工作
A: 确保 `static/js/search-data.json` 文件已正确生成。

#### Q: 样式显示异常
A: 检查配置文件中的颜色值格式是否正确（如 `#ffffff`）。

### 调试模式

运行时添加 `-v` 参数查看详细信息：
```bash
./creeper -v
```

### 重新生成
如果遇到问题，可以删除 `dist/` 目录后重新生成：
```bash
rm -rf dist/
./creeper
```

## 💡 最佳实践

1. **文件命名**：使用有意义的文件名，如 `01-序章.md`
2. **目录结构**：复杂小说使用多卷结构组织
3. **元数据**：为每部小说添加完整的元数据信息
4. **内容格式**：善用 Markdown 语法丰富内容表现
5. **定期备份**：定期备份你的小说文件

## 🆘 获取帮助

如果遇到问题，请：

1. 查看本使用指南
2. 检查示例文件格式
3. 在项目仓库提交 Issue

---

祝你创作愉快！📖✨
