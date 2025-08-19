# Creeper - 静态小说站点生成器

Creeper 是一个用 Go 语言开发的静态小说站点生成器，能够读取 Markdown 格式的小说文件并生成美观的静态网站。

## ✨ 特性

- 📖 **多格式支持**：支持单文件和多文件章节模式
- 🎨 **美观界面**：现代化响应式设计，支持移动端
- 🔍 **实时搜索**：支持小说和章节的全文搜索
- ⌨️ **键盘导航**：支持快捷键快速阅读
- 📊 **阅读进度**：显示章节阅读进度条
- 🎯 **高度可配置**：支持自定义主题色彩和样式
- ⚡ **快速生成**：高效的静态站点生成

## 🚀 快速开始

### 安装依赖

```bash
# 使用构建脚本（推荐）
./build.sh

# 或手动安装
go mod tidy
go build -o creeper main.go
```

### 基本用法

1. **准备小说文件**：项目已包含示例小说文件在 `novels` 目录
   - `示例小说.md` - 展示单文件多卷结构
   - `多章节小说/` - 展示简单多文件结构  
   - `多卷小说/` - 展示复杂多卷多文件结构

2. **生成静态站点**：
```bash
./creeper
```

3. **启动本地服务器**：
```bash
./creeper -serve
```

4. **访问网站**：打开浏览器访问 `http://localhost:8080`

### 命令行选项

```bash
./creeper [选项]

选项：
  -config string    配置文件路径 (默认 "config.yaml")
  -input string     小说文件输入目录 (默认 "novels")  
  -output string    静态站点输出目录 (默认 "dist")
  -serve           生成后启动本地服务器
  -port int        本地服务器端口 (默认 8080)
```

## 📚 小说文件格式

Creeper 支持两种小说组织方式：

### 单文件模式

将整部小说写在一个 `.md` 文件中，使用标题来分隔卷和章节：

```markdown
---
title: 我的小说
author: 作者姓名
description: 小说简介
---

# 第一卷 起源

## 第一章 开始

这里是第一章的内容...

## 第二章 发展

这里是第二章的内容...

# 第二卷 成长

## 第三章 转折

这里是第三章的内容...
```

**层次结构说明：**
- `# 第X卷 卷名` - 卷标题（一级标题）
- `## 第X章 章节名` - 章节标题（二级标题）
- `### 小节名` - 小节标题（三级标题）

### 多文件模式

支持两种多文件组织方式：

#### 简单多文件模式
为每个章节创建独立的 `.md` 文件，放在同一个目录下：

```
我的小说/
├── meta.md          # 小说元数据
├── 01-第一章.md     # 第一章
├── 02-第二章.md     # 第二章
└── 03-第三章.md     # 第三章
```

#### 多卷多文件模式
为复杂的长篇小说创建卷和章节的层次结构：

```
我的小说/
├── meta.md                    # 小说元数据
├── 第一卷-起源/
│   ├── 01-觉醒.md            # 第一卷第一章
│   ├── 02-启程.md            # 第一卷第二章
│   └── 03-试炼.md            # 第一卷第三章
├── 第二卷-成长/
│   ├── 01-挑战.md            # 第二卷第一章
│   ├── 02-突破.md            # 第二卷第二章
│   └── 03-蜕变.md            # 第二卷第三章
└── 第三卷-命运/
    ├── 01-真相.md            # 第三卷第一章
    ├── 02-决战.md            # 第三卷第二章
    └── 03-新生.md            # 第三卷第三章
```

**meta.md 示例：**
```markdown
---
title: 我的小说
author: 作者姓名
description: 小说简介
cover: static/images/cover.jpg
---
```

**章节文件示例：**
```markdown
# 第一章 开始

这里是章节内容...
```

## ⚙️ 配置文件

创建 `config.yaml` 文件来自定义站点设置：

```yaml
# 站点基本信息
site:
  title: "我的小说站点"
  description: "静态小说阅读站点" 
  author: "作者"
  base_url: "/"

# 目录配置
input_dir: "novels"
output_dir: "dist"

# 主题配置
theme:
  name: "default"
  primary_color: "#2c3e50"
  secondary_color: "#3498db"
  background_color: "#ffffff"
  text_color: "#333333"
  font_family: "'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif"
  font_size: "16px"
  line_height: "1.6"

# 构建配置
build:
  minify_html: true
  minify_css: true
  minify_js: true
```

## 🎨 主题定制

你可以通过修改配置文件中的 `theme` 部分来定制站点外观：

- `primary_color`: 主色调（导航栏、按钮等）
- `secondary_color`: 辅色调（链接、强调色等）
- `background_color`: 背景色
- `text_color`: 文字颜色
- `font_family`: 字体族
- `font_size`: 基础字体大小
- `line_height`: 行高

## 📱 响应式设计

生成的站点完全支持响应式设计，在桌面、平板和手机上都有良好的阅读体验。

## ⌨️ 键盘快捷键

在章节阅读页面支持以下快捷键：

- `Ctrl + ←`: 上一章
- `Ctrl + →`: 下一章  
- `Ctrl + ↑`: 返回目录
- `Esc`: 关闭搜索框

## 🔍 搜索功能

站点支持实时搜索功能：

- 搜索小说标题、作者
- 搜索章节标题
- 支持中文搜索
- 键盘导航搜索结果

## 📁 输出结构

生成的静态站点结构如下：

```
dist/
├── index.html              # 首页
├── novels/                 # 小说目录
│   ├── 小说1/
│   │   ├── index.html      # 小说目录页
│   │   ├── chapter-1.html  # 章节页面
│   │   └── ...
│   └── 小说2/
│       └── ...
└── static/                 # 静态资源
    ├── css/
    │   └── style.css       # 样式文件
    ├── js/
    │   ├── main.js         # 主脚本
    │   └── search-data.json # 搜索数据
    └── images/             # 图片资源
```

## 🛠️ 开发

### 项目结构

```
creeper/
├── main.go                 # 主程序入口
├── internal/
│   ├── config/            # 配置管理
│   ├── parser/            # Markdown 解析
│   └── generator/         # 静态站点生成
├── novels/                # 示例小说文件
├── config.yaml            # 配置文件
└── README.md
```

### 构建

```bash
# 构建可执行文件
go build -o creeper main.go

# 运行
./creeper -serve
```

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

如果你在使用过程中遇到问题，请：

1. 查看本文档
2. 检查示例文件
3. 提交 Issue

---

**享受阅读！** 📖
