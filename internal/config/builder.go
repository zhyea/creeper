package config

import "fmt"

// ConfigBuilder 配置建造者
type ConfigBuilder struct {
	config *Config
}

// NewConfigBuilder 创建配置建造者
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			Site:  SiteConfig{},
			Theme: ThemeConfig{},
			Build: BuildConfig{},
		},
	}
}

// WithSite 设置站点配置
func (b *ConfigBuilder) WithSite(title, description, author, baseURL string) *ConfigBuilder {
	b.config.Site = SiteConfig{
		Title:       title,
		Description: description,
		Author:      author,
		BaseURL:     baseURL,
	}
	return b
}

// WithSiteTitle 设置站点标题
func (b *ConfigBuilder) WithSiteTitle(title string) *ConfigBuilder {
	b.config.Site.Title = title
	return b
}

// WithSiteDescription 设置站点描述
func (b *ConfigBuilder) WithSiteDescription(description string) *ConfigBuilder {
	b.config.Site.Description = description
	return b
}

// WithSiteAuthor 设置站点作者
func (b *ConfigBuilder) WithSiteAuthor(author string) *ConfigBuilder {
	b.config.Site.Author = author
	return b
}

// WithSiteBaseURL 设置站点基础URL
func (b *ConfigBuilder) WithSiteBaseURL(baseURL string) *ConfigBuilder {
	b.config.Site.BaseURL = baseURL
	return b
}

// WithDirectories 设置目录配置
func (b *ConfigBuilder) WithDirectories(inputDir, outputDir string) *ConfigBuilder {
	b.config.InputDir = inputDir
	b.config.OutputDir = outputDir
	return b
}

// WithInputDir 设置输入目录
func (b *ConfigBuilder) WithInputDir(inputDir string) *ConfigBuilder {
	b.config.InputDir = inputDir
	return b
}

// WithOutputDir 设置输出目录
func (b *ConfigBuilder) WithOutputDir(outputDir string) *ConfigBuilder {
	b.config.OutputDir = outputDir
	return b
}

// WithTheme 设置主题配置
func (b *ConfigBuilder) WithTheme(name, primaryColor, secondaryColor, backgroundColor, textColor string) *ConfigBuilder {
	b.config.Theme.Name = name
	b.config.Theme.PrimaryColor = primaryColor
	b.config.Theme.SecondaryColor = secondaryColor
	b.config.Theme.BackgroundColor = backgroundColor
	b.config.Theme.TextColor = textColor
	return b
}

// WithThemeName 设置主题名称
func (b *ConfigBuilder) WithThemeName(name string) *ConfigBuilder {
	b.config.Theme.Name = name
	return b
}

// WithThemeColors 设置主题颜色
func (b *ConfigBuilder) WithThemeColors(primary, secondary, background, text string) *ConfigBuilder {
	b.config.Theme.PrimaryColor = primary
	b.config.Theme.SecondaryColor = secondary
	b.config.Theme.BackgroundColor = background
	b.config.Theme.TextColor = text
	return b
}

// WithThemePrimaryColor 设置主色调
func (b *ConfigBuilder) WithThemePrimaryColor(color string) *ConfigBuilder {
	b.config.Theme.PrimaryColor = color
	return b
}

// WithThemeSecondaryColor 设置辅色调
func (b *ConfigBuilder) WithThemeSecondaryColor(color string) *ConfigBuilder {
	b.config.Theme.SecondaryColor = color
	return b
}

// WithThemeBackgroundColor 设置背景色
func (b *ConfigBuilder) WithThemeBackgroundColor(color string) *ConfigBuilder {
	b.config.Theme.BackgroundColor = color
	return b
}

// WithThemeTextColor 设置文字颜色
func (b *ConfigBuilder) WithThemeTextColor(color string) *ConfigBuilder {
	b.config.Theme.TextColor = color
	return b
}

// WithFont 设置字体配置
func (b *ConfigBuilder) WithFont(family, size, lineHeight string) *ConfigBuilder {
	b.config.Theme.FontFamily = family
	b.config.Theme.FontSize = size
	b.config.Theme.LineHeight = lineHeight
	return b
}

// WithFontFamily 设置字体族
func (b *ConfigBuilder) WithFontFamily(family string) *ConfigBuilder {
	b.config.Theme.FontFamily = family
	return b
}

// WithFontSize 设置字体大小
func (b *ConfigBuilder) WithFontSize(size string) *ConfigBuilder {
	b.config.Theme.FontSize = size
	return b
}

// WithLineHeight 设置行高
func (b *ConfigBuilder) WithLineHeight(lineHeight string) *ConfigBuilder {
	b.config.Theme.LineHeight = lineHeight
	return b
}

// WithBuild 设置构建配置
func (b *ConfigBuilder) WithBuild(minifyHTML, minifyCSS, minifyJS bool) *ConfigBuilder {
	b.config.Build = BuildConfig{
		MinifyHTML: minifyHTML,
		MinifyCSS:  minifyCSS,
		MinifyJS:   minifyJS,
	}
	return b
}

// WithMinifyHTML 设置HTML压缩
func (b *ConfigBuilder) WithMinifyHTML(minify bool) *ConfigBuilder {
	b.config.Build.MinifyHTML = minify
	return b
}

// WithMinifyCSS 设置CSS压缩
func (b *ConfigBuilder) WithMinifyCSS(minify bool) *ConfigBuilder {
	b.config.Build.MinifyCSS = minify
	return b
}

// WithMinifyJS 设置JS压缩
func (b *ConfigBuilder) WithMinifyJS(minify bool) *ConfigBuilder {
	b.config.Build.MinifyJS = minify
	return b
}

// WithDefaults 使用默认配置
func (b *ConfigBuilder) WithDefaults() *ConfigBuilder {
	defaultConfig := Default()
	b.config = defaultConfig
	return b
}

// Build 构建配置
func (b *ConfigBuilder) Build() *Config {
	// 验证必要字段
	if b.config.Site.Title == "" {
		b.config.Site.Title = "我的小说站点"
	}
	if b.config.Site.Description == "" {
		b.config.Site.Description = "静态小说阅读站点"
	}
	if b.config.Site.BaseURL == "" {
		b.config.Site.BaseURL = "/"
	}
	if b.config.InputDir == "" {
		b.config.InputDir = "novels"
	}
	if b.config.OutputDir == "" {
		b.config.OutputDir = "dist"
	}
	
	// 设置默认主题
	if b.config.Theme.Name == "" {
		b.config.Theme.Name = "default"
	}
	if b.config.Theme.PrimaryColor == "" {
		b.config.Theme.PrimaryColor = "#2c3e50"
	}
	if b.config.Theme.SecondaryColor == "" {
		b.config.Theme.SecondaryColor = "#3498db"
	}
	if b.config.Theme.BackgroundColor == "" {
		b.config.Theme.BackgroundColor = "#ffffff"
	}
	if b.config.Theme.TextColor == "" {
		b.config.Theme.TextColor = "#333333"
	}
	if b.config.Theme.FontFamily == "" {
		b.config.Theme.FontFamily = "'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif"
	}
	if b.config.Theme.FontSize == "" {
		b.config.Theme.FontSize = "16px"
	}
	if b.config.Theme.LineHeight == "" {
		b.config.Theme.LineHeight = "1.6"
	}
	
	return b.config
}

// Validate 验证配置
func (b *ConfigBuilder) Validate() error {
	if b.config.Site.Title == "" {
		return fmt.Errorf("站点标题不能为空")
	}
	if b.config.InputDir == "" {
		return fmt.Errorf("输入目录不能为空")
	}
	if b.config.OutputDir == "" {
		return fmt.Errorf("输出目录不能为空")
	}
	return nil
}

// Clone 克隆配置建造者
func (b *ConfigBuilder) Clone() *ConfigBuilder {
	newConfig := &Config{
		Site:      b.config.Site,
		Theme:     b.config.Theme,
		Build:     b.config.Build,
		InputDir:  b.config.InputDir,
		OutputDir: b.config.OutputDir,
	}
	return &ConfigBuilder{config: newConfig}
}

// Reset 重置配置
func (b *ConfigBuilder) Reset() *ConfigBuilder {
	b.config = &Config{
		Site:  SiteConfig{},
		Theme: ThemeConfig{},
		Build: BuildConfig{},
	}
	return b
}
