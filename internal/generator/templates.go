package generator

import (
	"fmt"
	"html/template"
	"creeper/internal/parser"
)

// loadTemplates 加载所有模板
func (g *Generator) loadTemplates() error {
	// 创建模板函数
	funcMap := g.createTemplateFuncs()

	// 基础模板
	baseTemplate := g.getBaseTemplate()
	
	// 使用工厂模式创建模板
	factory := NewTemplateFactory(baseTemplate)
	
	// 创建各种类型的模板
	templateTypes := []TemplateType{IndexTemplate, NovelTemplate, ChapterTemplate}
	
	for _, templateType := range templateTypes {
		tmpl, err := factory.CreateTemplate(templateType, funcMap)
		if err != nil {
			return fmt.Errorf("创建%s模板失败: %v", templateType, err)
		}
		g.templates[string(templateType)] = tmpl
	}

	return nil
}

// getBaseTemplate 获取基础模板
func (g *Generator) getBaseTemplate() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <meta name="description" content="{{.Config.Site.Description}}">
    <meta name="author" content="{{.Config.Site.Author}}">
    <link rel="stylesheet" href="{{.Config.Site.BaseURL}}static/css/style.css">
    <link rel="stylesheet" href="{{.Config.Site.BaseURL}}static/css/reading-enhanced.css">
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
    <script src="{{.Config.Site.BaseURL}}static/js/reading-enhanced.js"></script>
</body>
</html>`
}

// createTemplateFuncs 创建模板函数映射
func (g *Generator) createTemplateFuncs() template.FuncMap {
	return template.FuncMap{
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
			total := 0
			for _, chapter := range chapters {
				total += chapter.WordCount
			}
			return g.formatWordCount(total)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
}

// formatWordCount 格式化字数显示
func (g *Generator) formatWordCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d字", count)
	} else if count < 10000 {
		return fmt.Sprintf("%.1f千字", float64(count)/1000)
	} else {
		return fmt.Sprintf("%.1f万字", float64(count)/10000)
	}
}
