package generator

import (
	"fmt"
	"html/template"
	"creeper/internal/parser"
)

// loadTemplates 加载所有模板
func (g *Generator) loadTemplates() error {
	// 创建模板函数
	funcMap := template.FuncMap{
		"sanitizeFileName": g.sanitizeFileName,
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"formatWordCount": func(count int) string {
			if count < 1000 {
				return fmt.Sprintf("%d字", count)
			} else if count < 10000 {
				return fmt.Sprintf("%.1f千字", float64(count)/1000)
			} else {
				return fmt.Sprintf("%.1f万字", float64(count)/10000)
			}
		},
		"totalWordCount": func(chapters []*parser.Chapter) string {
			// 计算章节总字数
			total := 0
			for _, chapter := range chapters {
				total += chapter.WordCount
			}
			// 使用 formatWordCount 格式化
			if total < 1000 {
				return fmt.Sprintf("%d字", total)
			} else if total < 10000 {
				return fmt.Sprintf("%.1f千字", float64(total)/1000)
			} else {
				return fmt.Sprintf("%.1f万字", float64(total)/10000)
			}
		},
	}

	// 基础模板
	baseTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <meta name="description" content="{{.Config.Site.Description}}">
    <meta name="author" content="{{.Config.Site.Author}}">
    <link rel="stylesheet" href="{{.Config.Site.BaseURL}}static/css/style.css">
    <link rel="icon" type="image/x-icon" href="{{.Config.Site.BaseURL}}static/images/favicon.ico">
</head>
<body>
    <header class="header">
        <div class="container">
            <h1 class="site-title">
                <a href="{{.Config.Site.BaseURL}}">{{.Config.Site.Title}}</a>
            </h1>
            <nav class="nav">
                <a href="{{.Config.Site.BaseURL}}" class="nav-link">首页</a>
                <div class="search-box">
                    <input type="text" id="search-input" placeholder="搜索小说或章节...">
                    <div id="search-results" class="search-results"></div>
                </div>
            </nav>
        </div>
    </header>

    <main class="main">
        <div class="container">
            {{template "content" .}}
        </div>
    </main>

    <footer class="footer">
        <div class="container">
            <p>&copy; 2024 {{.Config.Site.Title}}. 由 Creeper 生成</p>
        </div>
    </footer>

    <script src="{{.Config.Site.BaseURL}}static/js/main.js"></script>
</body>
</html>`

	// 首页模板
	indexTemplate := baseTemplate + `
{{define "content"}}
<div class="hero">
    <h2>{{.Config.Site.Description}}</h2>
    <p>共收录 {{len .Novels}} 部小说</p>
</div>

<div class="novels-grid">
    {{range .Novels}}
    <div class="novel-card">
        <div class="novel-cover">
            <img src="{{$.Config.Site.BaseURL}}novels/{{sanitizeFileName .Title}}/cover.svg" alt="{{.Title}} 封面" 
                 onerror="this.src='{{$.Config.Site.BaseURL}}static/images/default-cover.svg'">
        </div>
        <div class="novel-info">
            <h3 class="novel-title">
                <a href="{{$.Config.Site.BaseURL}}novels/{{sanitizeFileName .Title}}/">{{.Title}}</a>
            </h3>
            {{if .Author}}
            <p class="novel-author">作者：{{.Author}}</p>
            {{end}}
            {{if .Description}}
            <p class="novel-description">{{.Description}}</p>
            {{end}}
            <div class="novel-stats">
                <span class="chapter-count">{{len .Chapters}} 章</span>
                <span class="word-count">{{totalWordCount .Chapters}} 总字数</span>
            </div>
        </div>
    </div>
    {{end}}
</div>
{{end}}`

	// 小说目录模板
	novelTemplate := baseTemplate + `
{{define "content"}}
<div class="novel-header">
    <div class="novel-meta">
        <div class="novel-cover-large">
            <img src="{{$.Config.Site.BaseURL}}novels/{{sanitizeFileName .Novel.Title}}/cover.svg" alt="{{.Novel.Title}} 封面"
                 onerror="this.src='{{$.Config.Site.BaseURL}}static/images/default-cover.svg'">
        </div>
        <div class="novel-details">
            <h1 class="novel-title">{{.Novel.Title}}</h1>
            {{if .Novel.Author}}
            <p class="novel-author">作者：{{.Novel.Author}}</p>
            {{end}}
            {{if .Novel.Description}}
            <p class="novel-description">{{.Novel.Description}}</p>
            {{end}}
            <div class="novel-stats">
                <span class="chapter-count">共 {{len .Novel.Chapters}} 章</span>
                <span class="update-time">更新于 {{.Novel.UpdatedAt.Format "2006-01-02"}}</span>
            </div>
            <div class="novel-actions">
                <a href="chapter-1.html" class="btn btn-primary">开始阅读</a>
            </div>
        </div>
    </div>
</div>

<div class="chapters-list">
    <h2>章节目录</h2>
    <div class="chapters-grid">
        {{range .Novel.Chapters}}
        <div class="chapter-item">
            <a href="chapter-{{.ID}}.html" class="chapter-link">
                <span class="chapter-title">{{.Title}}</span>
                <span class="chapter-stats">{{formatWordCount .WordCount}}</span>
            </a>
        </div>
        {{end}}
    </div>
</div>
{{end}}`

	// 章节阅读模板
	chapterTemplate := baseTemplate + `
{{define "content"}}
<div class="chapter-header">
    <nav class="breadcrumb">
        <a href="{{$.Config.Site.BaseURL}}">首页</a>
        <span class="separator">/</span>
        <a href="./index.html">{{.Novel.Title}}</a>
        <span class="separator">/</span>
        <span class="current">{{.Chapter.Title}}</span>
    </nav>
    
    <h1 class="chapter-title">{{.Chapter.Title}}</h1>
    
    <div class="chapter-nav">
        {{if gt .Chapter.ID 1}}
        <a href="chapter-{{sub .Chapter.ID 1}}.html" class="btn btn-nav">上一章</a>
        {{end}}
        <a href="./index.html" class="btn btn-nav">目录</a>
        {{if lt .Chapter.ID (len .Novel.Chapters)}}
        <a href="chapter-{{add .Chapter.ID 1}}.html" class="btn btn-nav">下一章</a>
        {{end}}
    </div>
</div>

<article class="chapter-content">
    {{.Chapter.HTMLContent | printf "%s" | safeHTML}}
</article>

<div class="chapter-footer">
    <div class="chapter-info">
        <p>字数：{{formatWordCount .Chapter.WordCount}}</p>
    </div>
    
    <div class="chapter-nav">
        {{if gt .Chapter.ID 1}}
        <a href="chapter-{{sub .Chapter.ID 1}}.html" class="btn btn-nav">上一章</a>
        {{end}}
        <a href="./index.html" class="btn btn-nav">目录</a>
        {{if lt .Chapter.ID (len .Novel.Chapters)}}
        <a href="chapter-{{add .Chapter.ID 1}}.html" class="btn btn-nav">下一章</a>
        {{end}}
    </div>
</div>
{{end}}`

	// 编译模板
	var err error

	g.templates["index"], err = template.New("index").Funcs(funcMap).Parse(indexTemplate)
	if err != nil {
		return fmt.Errorf("编译首页模板失败: %v", err)
	}

	g.templates["novel"], err = template.New("novel").Funcs(funcMap).Parse(novelTemplate)
	if err != nil {
		return fmt.Errorf("编译小说模板失败: %v", err)
	}

	// 为章节模板添加 safeHTML 函数
	chapterFuncMap := template.FuncMap{
		"sanitizeFileName": g.sanitizeFileName,
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"formatWordCount": func(count int) string {
			if count < 1000 {
				return fmt.Sprintf("%d字", count)
			} else if count < 10000 {
				return fmt.Sprintf("%.1f千字", float64(count)/1000)
			} else {
				return fmt.Sprintf("%.1f万字", float64(count)/10000)
			}
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	g.templates["chapter"], err = template.New("chapter").Funcs(chapterFuncMap).Parse(chapterTemplate)
	if err != nil {
		return fmt.Errorf("编译章节模板失败: %v", err)
	}

	return nil
}
