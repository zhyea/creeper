package generator

import (
	"fmt"
	"html/template"
)

// TemplateType 模板类型
type TemplateType string

const (
	IndexTemplate   TemplateType = "index"
	NovelTemplate   TemplateType = "novel"
	ChapterTemplate TemplateType = "chapter"
)

// TemplateBuilder 模板构建器接口
type TemplateBuilder interface {
	Build(funcMap template.FuncMap) (*template.Template, error)
	GetType() TemplateType
}

// BaseTemplateBuilder 基础模板构建器
type BaseTemplateBuilder struct {
	templateType TemplateType
	baseTemplate string
}

func (b *BaseTemplateBuilder) GetType() TemplateType {
	return b.templateType
}

// IndexTemplateBuilder 首页模板构建器
type IndexTemplateBuilder struct {
	*BaseTemplateBuilder
}

func NewIndexTemplateBuilder(baseTemplate string) *IndexTemplateBuilder {
	return &IndexTemplateBuilder{
		BaseTemplateBuilder: &BaseTemplateBuilder{
			templateType: IndexTemplate,
			baseTemplate: baseTemplate,
		},
	}
}

func (b *IndexTemplateBuilder) Build(funcMap template.FuncMap) (*template.Template, error) {
	indexContent := `
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

	templateContent := b.baseTemplate + indexContent
	return template.New("index").Funcs(funcMap).Parse(templateContent)
}

// NovelTemplateBuilder 小说模板构建器
type NovelTemplateBuilder struct {
	*BaseTemplateBuilder
}

func NewNovelTemplateBuilder(baseTemplate string) *NovelTemplateBuilder {
	return &NovelTemplateBuilder{
		BaseTemplateBuilder: &BaseTemplateBuilder{
			templateType: NovelTemplate,
			baseTemplate: baseTemplate,
		},
	}
}

func (b *NovelTemplateBuilder) Build(funcMap template.FuncMap) (*template.Template, error) {
	novelContent := `
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

	templateContent := b.baseTemplate + novelContent
	return template.New("novel").Funcs(funcMap).Parse(templateContent)
}

// ChapterTemplateBuilder 章节模板构建器
type ChapterTemplateBuilder struct {
	*BaseTemplateBuilder
}

func NewChapterTemplateBuilder(baseTemplate string) *ChapterTemplateBuilder {
	return &ChapterTemplateBuilder{
		BaseTemplateBuilder: &BaseTemplateBuilder{
			templateType: ChapterTemplate,
			baseTemplate: baseTemplate,
		},
	}
}

func (b *ChapterTemplateBuilder) Build(funcMap template.FuncMap) (*template.Template, error) {
	chapterContent := `
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

	templateContent := b.baseTemplate + chapterContent
	return template.New("chapter").Funcs(funcMap).Parse(templateContent)
}

// TemplateFactory 模板工厂
type TemplateFactory struct {
	baseTemplate string
	builders     map[TemplateType]TemplateBuilder
}

// NewTemplateFactory 创建模板工厂
func NewTemplateFactory(baseTemplate string) *TemplateFactory {
	factory := &TemplateFactory{
		baseTemplate: baseTemplate,
		builders:     make(map[TemplateType]TemplateBuilder),
	}
	
	// 注册模板构建器
	factory.RegisterBuilder(NewIndexTemplateBuilder(baseTemplate))
	factory.RegisterBuilder(NewNovelTemplateBuilder(baseTemplate))
	factory.RegisterBuilder(NewChapterTemplateBuilder(baseTemplate))
	
	return factory
}

// RegisterBuilder 注册模板构建器
func (f *TemplateFactory) RegisterBuilder(builder TemplateBuilder) {
	f.builders[builder.GetType()] = builder
}

// CreateTemplate 创建模板
func (f *TemplateFactory) CreateTemplate(templateType TemplateType, funcMap template.FuncMap) (*template.Template, error) {
	builder, exists := f.builders[templateType]
	if !exists {
		return nil, fmt.Errorf("未知模板类型: %s", templateType)
	}
	
	return builder.Build(funcMap)
}

// GetAvailableTypes 获取可用的模板类型
func (f *TemplateFactory) GetAvailableTypes() []TemplateType {
	types := make([]TemplateType, 0, len(f.builders))
	for templateType := range f.builders {
		types = append(types, templateType)
	}
	return types
}
