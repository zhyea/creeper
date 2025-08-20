package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// CoverTheme 封面主题配置
type CoverTheme struct {
	Name        string   `json:"name"`
	BgGradient  []string `json:"bg_gradient"`
	TextColor   string   `json:"text_color"`
	AccentColor string   `json:"accent_color"`
	Style       string   `json:"style"`
	Description string   `json:"description"`
}

// CoverGenerator 封面生成器
type CoverGenerator struct {
	themes map[string]CoverTheme
}

// NewCoverGenerator 创建新的封面生成器
func NewCoverGenerator() *CoverGenerator {
	return &CoverGenerator{
		themes: getDefaultThemes(),
	}
}

// getDefaultThemes 获取默认主题配置
func getDefaultThemes() map[string]CoverTheme {
	return map[string]CoverTheme{
		"default": {
			Name:        "default",
			BgGradient:  []string{"#2c3e50", "#3498db"},
			TextColor:   "#ffffff",
			AccentColor: "#f1c40f",
			Style:       "modern",
			Description: "简洁现代的设计风格",
		},
		"fantasy": {
			Name:        "fantasy",
			BgGradient:  []string{"#8e44ad", "#2c3e50", "#1a1a2e"},
			TextColor:   "#ffffff",
			AccentColor: "#e74c3c",
			Style:       "fantasy",
			Description: "奇幻魔法主题，适合玄幻小说",
		},
		"modern": {
			Name:        "modern",
			BgGradient:  []string{"#667eea", "#764ba2"},
			TextColor:   "#ffffff",
			AccentColor: "#ffffff",
			Style:       "geometric",
			Description: "现代几何风格，简约时尚",
		},
		"classical": {
			Name:        "classical",
			BgGradient:  []string{"#8b4513", "#a0522d", "#654321"},
			TextColor:   "#8b4513",
			AccentColor: "#ffd700",
			Style:       "ornate",
			Description: "古典文学风格，典雅庄重",
		},
		"scifi": {
			Name:        "scifi",
			BgGradient:  []string{"#0a0a23", "#1a1a2e", "#000000"},
			TextColor:   "#00ffff",
			AccentColor: "#0080ff",
			Style:       "tech",
			Description: "科幻未来主题，霓虹科技感",
		},
		"wuxia": {
			Name:        "wuxia",
			BgGradient:  []string{"#f5f5dc", "#e6ddd4", "#d2b48c"},
			TextColor:   "#2f4f4f",
			AccentColor: "#dc143c",
			Style:       "traditional",
			Description: "武侠江湖风格，水墨山水意境",
		},
	}
}

// Config 命令行配置
type Config struct {
	Title      string
	Subtitle   string
	Theme      string
	Output     string
	Width      int
	Height     int
	ListThemes bool
}

// parseFlags 解析命令行参数
func parseFlags() *Config {
	config := &Config{}
	
	flag.StringVar(&config.Title, "title", "", "小说标题 (必需)")
	flag.StringVar(&config.Subtitle, "subtitle", "", "副标题")
	flag.StringVar(&config.Theme, "theme", "default", "主题风格")
	flag.StringVar(&config.Output, "output", "", "输出文件名")
	flag.IntVar(&config.Width, "width", 300, "宽度 (像素)")
	flag.IntVar(&config.Height, "height", 400, "高度 (像素)")
	flag.BoolVar(&config.ListThemes, "list-themes", false, "列出所有主题")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Creeper 封面生成器\n\n")
		fmt.Fprintf(os.Stderr, "用法: %s [选项]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  %s -title \"我的小说\" -theme fantasy\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -title \"科幻故事\" -theme scifi -subtitle \"未来世界\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -list-themes\n", os.Args[0])
	}
	
	flag.Parse()
	return config
}

// validateConfig 验证配置
func (c *Config) validate() error {
	if c.Title == "" && !c.ListThemes {
		return fmt.Errorf("必须指定标题")
	}
	
	if c.Width <= 0 || c.Height <= 0 {
		return fmt.Errorf("宽度和高度必须大于0")
	}
	
	if len(c.Title) > 50 {
		return fmt.Errorf("标题长度不能超过50个字符")
	}
	
	if len(c.Subtitle) > 30 {
		return fmt.Errorf("副标题长度不能超过30个字符")
	}
	
	return nil
}

func main() {
	config := parseFlags()
	
	if err := config.validate(); err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	
	generator := NewCoverGenerator()
	
	if config.ListThemes {
		generator.listThemes()
		return
	}
	
	if err := generator.generateCover(config); err != nil {
		log.Fatalf("生成封面失败: %v", err)
	}
}

// listThemes 列出所有可用主题
func (g *CoverGenerator) listThemes() {
	fmt.Println("🎨 可用主题:")
	fmt.Println()
	
	// 按名称排序
	names := make([]string, 0, len(g.themes))
	for name := range g.themes {
		names = append(names, name)
	}
	sort.Strings(names)
	
	for _, name := range names {
		theme := g.themes[name]
		fmt.Printf("  %-12s %s\n", theme.Name+":", theme.Description)
		fmt.Printf("  %-12s 风格: %s\n", "", theme.Style)
		fmt.Printf("  %-12s 颜色: %s\n", "", strings.Join(theme.BgGradient, " → "))
		fmt.Println()
	}
}

// generateCover 生成封面
func (g *CoverGenerator) generateCover(config *Config) error {
	// 检查主题是否存在
	theme, exists := g.themes[config.Theme]
	if !exists {
		return fmt.Errorf("未知主题 '%s'，可用主题: %s", config.Theme, g.getThemeNames())
	}
	
	// 生成 SVG 内容
	svgContent := g.generateSVGCover(config.Title, config.Subtitle, theme, config.Width, config.Height)
	
	// 确定输出文件名
	outputFile := config.Output
	if outputFile == "" {
		safeTitle := sanitizeFileName(config.Title)
		outputFile = fmt.Sprintf("static/images/%s-cover.svg", safeTitle)
	}
	
	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(outputFile, []byte(svgContent), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	
	// 输出成功信息
	g.printSuccess(outputFile, config, theme)
	return nil
}

// generateSVGCover 生成 SVG 封面
func (g *CoverGenerator) generateSVGCover(title, subtitle string, theme CoverTheme, width, height int) string {
	gradient := g.createGradient(theme.BgGradient)
	decorations := g.generateDecorations(theme.Name, theme.AccentColor)
	
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">
    <defs>%s</defs>
    
    <!-- 背景 -->
    <rect width="%d" height="%d" fill="url(#bgGradient)"/>
    
    %s
    
    <!-- 装饰元素 -->
    <g transform="translate(%d, %d)">
        <circle cx="0" cy="0" r="8" fill="%s" opacity="0.4"/>
        <circle cx="0" cy="0" r="4" fill="%s" opacity="0.7"/>
    </g>
</svg>`, width, height, width, height, gradient, width, height, decorations, width/2, height-50, theme.AccentColor, theme.AccentColor)
}

// createGradient 创建渐变定义
func (g *CoverGenerator) createGradient(colors []string) string {
	gradient := ""
	
	if len(colors) == 2 {
		gradient = fmt.Sprintf(`
        <linearGradient id="bgGradient" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
            <stop offset="0%%" style="stop-color:%s;stop-opacity:1" />
            <stop offset="100%%" style="stop-color:%s;stop-opacity:1" />
        </linearGradient>`, colors[0], colors[1])
	} else if len(colors) >= 3 {
		gradient = fmt.Sprintf(`
        <radialGradient id="bgGradient" cx="50%%" cy="30%%" r="80%%">
            <stop offset="0%%" style="stop-color:%s;stop-opacity:1" />
            <stop offset="50%%" style="stop-color:%s;stop-opacity:1" />
            <stop offset="100%%" style="stop-color:%s;stop-opacity:1" />
        </radialGradient>`, colors[0], colors[1], colors[2])
	} else {
		// 默认单色
		gradient = fmt.Sprintf(`
        <linearGradient id="bgGradient">
            <stop offset="0%%" style="stop-color:%s;stop-opacity:1" />
        </linearGradient>`, colors[0])
	}
	
	// 为现代主题添加额外的渐变
	gradient += `
        <linearGradient id="accentGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:#f093fb;stop-opacity:0.8" />
            <stop offset="100%" style="stop-color:#f5576c;stop-opacity:0.6" />
        </linearGradient>`
	
	return gradient
}

// printSuccess 输出成功信息
func (g *CoverGenerator) printSuccess(outputFile string, config *Config, theme CoverTheme) {
	fmt.Printf("✅ 封面已生成: %s\n", outputFile)
	fmt.Printf("📖 标题: %s\n", config.Title)
	if config.Subtitle != "" {
		fmt.Printf("📝 副标题: %s\n", config.Subtitle)
	}
	fmt.Printf("🎨 主题: %s (%s)\n", config.Theme, theme.Description)
	fmt.Printf("📐 尺寸: %dx%d 像素\n", config.Width, config.Height)
}

// getThemeNames 获取所有主题名称
func (g *CoverGenerator) getThemeNames() string {
	names := make([]string, 0, len(g.themes))
	for name := range g.themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// generateDecorations 根据主题生成装饰元素
func (g *CoverGenerator) generateDecorations(themeName, accentColor string) string {
	switch themeName {
	case "fantasy":
		return `
        <!-- 星星装饰 -->
        <g fill="#ffffff" opacity="0.8">
            <circle cx="50" cy="60" r="1"/>
            <circle cx="250" cy="80" r="1.5"/>
            <circle cx="80" cy="320" r="1"/>
        </g>
        <!-- 城堡剪影 -->
        <g fill="#000000" opacity="0.3">
            <rect x="120" y="280" width="60" height="50" rx="5"/>
            <polygon points="140,280 150,260 160,280"/>
        </g>`

	case "scifi":
		return `
        <!-- 星空 -->
        <g fill="#ffffff">
            <circle cx="50" cy="50" r="0.5" opacity="0.8"/>
            <circle cx="250" cy="80" r="1" opacity="0.6"/>
            <circle cx="80" cy="300" r="0.5" opacity="0.9"/>
        </g>
        <!-- 科技线条 -->
        <g stroke="#00ffff" stroke-width="1" fill="none" opacity="0.6">
            <path d="M50 250 L100 270 L150 250 L200 270 L250 250"/>
        </g>`

	case "modern":
		return `
        <!-- 现代几何装饰 -->
        <g fill="#ffffff" opacity="0.2">
            <circle cx="220" cy="120" r="80"/>
            <circle cx="100" cy="300" r="50"/>
        </g>
        <g fill="url(#accentGradient)" opacity="0.4">
            <rect x="180" y="80" width="50" height="50" rx="8" transform="rotate(15 205 105)"/>
            <rect x="70" y="250" width="35" height="35" rx="6" transform="rotate(-20 87 267)"/>
        </g>`

	case "classical":
		return `
        <!-- 装饰边框 -->
        <rect x="30" y="30" width="240" height="340" fill="none" stroke="#ffd700" stroke-width="2" rx="10"/>
        <!-- 装饰花纹 -->
        <g fill="#ffd700" opacity="0.6">
            <circle cx="150" cy="80" r="20" fill="none" stroke="#ffd700" stroke-width="2"/>
        </g>`

	case "wuxia":
		return `
        <!-- 远山剪影 -->
        <g fill="#696969" opacity="0.4">
            <path d="M0,200 Q50,180 100,190 Q150,170 200,185 Q250,175 300,190 L300,400 L0,400 Z"/>
            <path d="M0,220 Q60,200 120,210 Q180,195 240,205 Q270,200 300,210 L300,400 L0,400 Z"/>
        </g>
        
        <!-- 竹林 -->
        <g fill="#2f4f4f" opacity="0.3">
            <rect x="50" y="200" width="3" height="80" rx="1"/>
            <rect x="60" y="190" width="3" height="90" rx="1"/>
            <rect x="70" y="205" width="3" height="75" rx="1"/>
        </g>
        
        <!-- 剑影 -->
        <g transform="translate(150, 150) rotate(-15)">
            <rect x="-1" y="-40" width="2" height="80" fill="#c0c0c0" opacity="0.6"/>
            <polygon points="0,-42 -2,-40 2,-40" fill="#c0c0c0" opacity="0.6"/>
        </g>
        
        <!-- 印章 -->
        <g transform="translate(230, 320)">
            <circle cx="0" cy="0" r="15" fill="#dc143c" opacity="0.7"/>
            <rect x="-6" y="-6" width="4" height="4" fill="#ffffff" opacity="0.9"/>
            <rect x="2" y="-6" width="4" height="4" fill="#ffffff" opacity="0.9"/>
            <rect x="-6" y="2" width="4" height="4" fill="#ffffff" opacity="0.9"/>
            <rect x="2" y="2" width="4" height="4" fill="#ffffff" opacity="0.9"/>
        </g>`

	default:
		return `
        <!-- 书本形状 -->
        <g transform="translate(125, 150)">
            <rect x="0" y="0" width="50" height="40" fill="#ffffff" opacity="0.8" rx="3"/>
            <rect x="5" y="5" width="40" height="30" fill="none" stroke="` + accentColor + `" stroke-width="2" rx="2"/>
            <line x1="10" y1="15" x2="40" y2="15" stroke="` + accentColor + `" stroke-width="1"/>
            <line x1="10" y1="20" x2="35" y2="20" stroke="` + accentColor + `" stroke-width="1"/>
            <line x1="10" y1="25" x2="30" y2="25" stroke="` + accentColor + `" stroke-width="1"/>
        </g>`
	}
}

// sanitizeFileName 清理文件名，保留中文字符
func sanitizeFileName(name string) string {
	var result strings.Builder
	
	for _, r := range name {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			result.WriteRune(r)
		case unicode.IsSpace(r):
			result.WriteRune('-')
		case r == '-' || r == '_':
			result.WriteRune(r)
		default:
			// 跳过其他特殊字符
		}
	}
	
	cleaned := result.String()
	
	// 移除连续的横线和下划线
	for strings.Contains(cleaned, "--") {
		cleaned = strings.ReplaceAll(cleaned, "--", "-")
	}
	for strings.Contains(cleaned, "__") {
		cleaned = strings.ReplaceAll(cleaned, "__", "_")
	}
	for strings.Contains(cleaned, "-_") || strings.Contains(cleaned, "_-") {
		cleaned = strings.ReplaceAll(cleaned, "-_", "-")
		cleaned = strings.ReplaceAll(cleaned, "_-", "-")
	}
	
	// 移除首尾的特殊字符
	cleaned = strings.Trim(cleaned, "-_")
	
	// 如果结果为空，使用默认名称
	if cleaned == "" {
		cleaned = "cover"
	}
	
	// 限制长度
	if len([]rune(cleaned)) > 30 {
		runes := []rune(cleaned)
		cleaned = string(runes[:30])
	}
	
	return cleaned
}
