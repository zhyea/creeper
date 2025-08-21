package factory

import (
	"fmt"

	"creeper/internal/config"
	"creeper/internal/parser"
	"creeper/internal/generator"
)

// GeneratorType 生成器类型
type GeneratorType string

const (
	StaticGenerator    GeneratorType = "static"
	DynamicGenerator   GeneratorType = "dynamic"
	MinimalGenerator   GeneratorType = "minimal"
	EnhancedGenerator  GeneratorType = "enhanced"
)

// AbstractGeneratorFactory 抽象生成器工厂
type AbstractGeneratorFactory interface {
	CreateParser() parser.ParseStrategy
	CreateGenerator(config *config.Config) *generator.Generator
	CreateTemplateFactory() *generator.TemplateFactory
	GetGeneratorType() GeneratorType
	GetDescription() string
}

// StaticGeneratorFactory 静态生成器工厂
type StaticGeneratorFactory struct{}

func (sgf *StaticGeneratorFactory) CreateParser() parser.ParseStrategy {
	baseParser := parser.New()
	return parser.NewSingleFileStrategy(baseParser)
}

func (sgf *StaticGeneratorFactory) CreateGenerator(config *config.Config) *generator.Generator {
	return generator.New(config)
}

func (sgf *StaticGeneratorFactory) CreateTemplateFactory() *generator.TemplateFactory {
	baseTemplate := sgf.getBaseTemplate()
	return generator.NewTemplateFactory(baseTemplate)
}

func (sgf *StaticGeneratorFactory) GetGeneratorType() GeneratorType {
	return StaticGenerator
}

func (sgf *StaticGeneratorFactory) GetDescription() string {
	return "标准静态站点生成器，支持基础功能"
}

func (sgf *StaticGeneratorFactory) getBaseTemplate() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="{{.Config.Site.BaseURL}}static/css/style.css">
</head>
<body>
    {{template "content" .}}
</body>
</html>`
}

// EnhancedGeneratorFactory 增强生成器工厂
type EnhancedGeneratorFactory struct{}

func (egf *EnhancedGeneratorFactory) CreateParser() parser.ParseStrategy {
	baseParser := parser.New()
	
	// 创建策略管理器，支持所有格式
	strategyManager := parser.NewStrategyManager(baseParser)
	
	// 返回多格式策略（这里简化为返回 TXT 策略）
	return parser.NewTxtFileStrategy(baseParser)
}

func (egf *EnhancedGeneratorFactory) CreateGenerator(config *config.Config) *generator.Generator {
	return generator.New(config)
}

func (egf *EnhancedGeneratorFactory) CreateTemplateFactory() *generator.TemplateFactory {
	baseTemplate := egf.getEnhancedBaseTemplate()
	return generator.NewTemplateFactory(baseTemplate)
}

func (egf *EnhancedGeneratorFactory) GetGeneratorType() GeneratorType {
	return EnhancedGenerator
}

func (egf *EnhancedGeneratorFactory) GetDescription() string {
	return "增强型生成器，支持多格式、多主题、增强阅读体验"
}

func (egf *EnhancedGeneratorFactory) getEnhancedBaseTemplate() string {
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

// MinimalGeneratorFactory 最小化生成器工厂
type MinimalGeneratorFactory struct{}

func (mgf *MinimalGeneratorFactory) CreateParser() parser.ParseStrategy {
	baseParser := parser.New()
	return parser.NewSingleFileStrategy(baseParser)
}

func (mgf *MinimalGeneratorFactory) CreateGenerator(config *config.Config) *generator.Generator {
	// 创建最小化配置
	minimalConfig := config.Builder().
		WithSiteTitle(config.Site.Title).
		WithInputDir(config.InputDir).
		WithOutputDir(config.OutputDir).
		WithBuild(false, false, false). // 不压缩
		Build()
	
	return generator.New(minimalConfig)
}

func (mgf *MinimalGeneratorFactory) CreateTemplateFactory() *generator.TemplateFactory {
	baseTemplate := mgf.getMinimalTemplate()
	return generator.NewTemplateFactory(baseTemplate)
}

func (mgf *MinimalGeneratorFactory) GetGeneratorType() GeneratorType {
	return MinimalGenerator
}

func (mgf *MinimalGeneratorFactory) GetDescription() string {
	return "最小化生成器，生成简洁的静态站点"
}

func (mgf *MinimalGeneratorFactory) getMinimalTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        .container { max-width: 800px; margin: 0 auto; }
        h1, h2, h3 { color: #333; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="container">
        {{template "content" .}}
    </div>
</body>
</html>`
}

// GeneratorFactoryRegistry 生成器工厂注册表
type GeneratorFactoryRegistry struct {
	factories map[GeneratorType]AbstractGeneratorFactory
}

// NewGeneratorFactoryRegistry 创建工厂注册表
func NewGeneratorFactoryRegistry() *GeneratorFactoryRegistry {
	registry := &GeneratorFactoryRegistry{
		factories: make(map[GeneratorType]AbstractGeneratorFactory),
	}
	
	// 注册默认工厂
	registry.RegisterFactory(StaticGenerator, &StaticGeneratorFactory{})
	registry.RegisterFactory(EnhancedGenerator, &EnhancedGeneratorFactory{})
	registry.RegisterFactory(MinimalGenerator, &MinimalGeneratorFactory{})
	
	return registry
}

// RegisterFactory 注册工厂
func (gfr *GeneratorFactoryRegistry) RegisterFactory(generatorType GeneratorType, factory AbstractGeneratorFactory) {
	gfr.factories[generatorType] = factory
}

// GetFactory 获取工厂
func (gfr *GeneratorFactoryRegistry) GetFactory(generatorType GeneratorType) (AbstractGeneratorFactory, error) {
	factory, exists := gfr.factories[generatorType]
	if !exists {
		return nil, fmt.Errorf("未知生成器类型: %s", generatorType)
	}
	return factory, nil
}

// GetAvailableTypes 获取可用类型
func (gfr *GeneratorFactoryRegistry) GetAvailableTypes() []GeneratorType {
	types := make([]GeneratorType, 0, len(gfr.factories))
	for generatorType := range gfr.factories {
		types = append(types, generatorType)
	}
	return types
}

// GetFactoryDescriptions 获取工厂描述
func (gfr *GeneratorFactoryRegistry) GetFactoryDescriptions() map[GeneratorType]string {
	descriptions := make(map[GeneratorType]string)
	for generatorType, factory := range gfr.factories {
		descriptions[generatorType] = factory.GetDescription()
	}
	return descriptions
}

// CreateGeneratorSuite 创建生成器套件
func (gfr *GeneratorFactoryRegistry) CreateGeneratorSuite(generatorType GeneratorType, config *config.Config) (*GeneratorSuite, error) {
	factory, err := gfr.GetFactory(generatorType)
	if err != nil {
		return nil, err
	}
	
	return &GeneratorSuite{
		Parser:          factory.CreateParser(),
		Generator:       factory.CreateGenerator(config),
		TemplateFactory: factory.CreateTemplateFactory(),
		Type:            generatorType,
		Description:     factory.GetDescription(),
	}, nil
}

// GeneratorSuite 生成器套件
type GeneratorSuite struct {
	Parser          parser.ParseStrategy
	Generator       *generator.Generator
	TemplateFactory *generator.TemplateFactory
	Type            GeneratorType
	Description     string
}

// GetInfo 获取套件信息
func (gs *GeneratorSuite) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"type":        gs.Type,
		"description": gs.Description,
		"parser":      gs.Parser.GetName(),
		"templates":   gs.TemplateFactory.GetAvailableTypes(),
	}
}
