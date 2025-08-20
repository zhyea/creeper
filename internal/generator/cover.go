package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"creeper/internal/parser"
)

// generateNovelCover 为小说生成带标题的封面
func (g *Generator) generateNovelCover(novel *parser.Novel) error {
	// 如果没有指定封面，使用默认封面
	coverPath := novel.Cover
	if coverPath == "" {
		coverPath = "static/images/default-cover.svg"
	}

	// 读取原始封面文件
	originalCoverPath := filepath.Join(g.config.InputDir, "..", coverPath)
	if _, err := os.Stat(originalCoverPath); os.IsNotExist(err) {
		// 如果封面文件不存在，尝试从项目根目录查找
		originalCoverPath = coverPath
		if _, err := os.Stat(originalCoverPath); os.IsNotExist(err) {
			// 如果仍然不存在，跳过封面生成
			return nil
		}
	}

	svgContent, err := os.ReadFile(originalCoverPath)
	if err != nil {
		return fmt.Errorf("读取封面文件失败: %v", err)
	}

	// 为小说生成带标题的封面
	modifiedSVG := g.addTitleToCover(string(svgContent), novel.Title, novel.Author)

	// 生成输出路径
	novelDir := filepath.Join(g.config.OutputDir, "novels", g.sanitizeFileName(novel.Title))
	coverOutputPath := filepath.Join(novelDir, "cover.svg")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(coverOutputPath), 0755); err != nil {
		return fmt.Errorf("创建封面目录失败: %v", err)
	}

	// 写入修改后的封面
	return os.WriteFile(coverOutputPath, []byte(modifiedSVG), 0644)
}

// addTitleToCover 在封面上添加标题
func (g *Generator) addTitleToCover(svgContent, title, author string) string {
	// 检测封面风格
	style := g.detectCoverStyle(svgContent)

	// 根据不同风格选择合适的标题样式
	titleElement := g.generateTitleElement(title, author, style)

	// 在 </svg> 标签前插入标题元素
	svgEndRegex := regexp.MustCompile(`</svg>\s*$`)
	return svgEndRegex.ReplaceAllString(svgContent, titleElement+"\n</svg>")
}

// detectCoverStyle 检测封面风格
func (g *Generator) detectCoverStyle(svgContent string) string {
	content := strings.ToLower(svgContent)
	
	if strings.Contains(content, "fantasy") || strings.Contains(content, "castle") || strings.Contains(content, "magic") {
		return "fantasy"
	} else if strings.Contains(content, "scifi") || strings.Contains(content, "neon") || strings.Contains(content, "#00ffff") {
		return "scifi"
	} else if strings.Contains(content, "classical") || strings.Contains(content, "serif") || strings.Contains(content, "#8b4513") {
		return "classical"
	} else if strings.Contains(content, "modern") || strings.Contains(content, "geometric") {
		return "modern"
	}
	
	return "default"
}

// generateTitleElement 生成标题元素
func (g *Generator) generateTitleElement(title, author, style string) string {
	// 限制标题长度，避免溢出
	displayTitle := title
	if len([]rune(title)) > 12 {
		displayTitle = string([]rune(title)[:12]) + "..."
	}

	switch style {
	case "fantasy":
		return g.generateFantasyTitle(displayTitle, author)
	case "scifi":
		return g.generateScifiTitle(displayTitle, author)
	case "classical":
		return g.generateClassicalTitle(displayTitle, author)
	case "modern":
		return g.generateModernTitle(displayTitle, author)
	default:
		return g.generateDefaultTitle(displayTitle, author)
	}
}

// generateDefaultTitle 生成默认风格标题
func (g *Generator) generateDefaultTitle(title, author string) string {
	return fmt.Sprintf(`
  <!-- 动态标题 -->
  <g id="dynamic-title">
    <rect x="40" y="320" width="220" height="60" fill="#000000" opacity="0.4" rx="10"/>
    <text x="150" y="345" text-anchor="middle" fill="#ffffff" font-family="Arial, sans-serif" font-size="18" font-weight="bold">
      %s
    </text>
    %s
  </g>`, title, g.generateAuthorText(author, "150", "365", "#ecf0f1", "Arial, sans-serif", "12"))
}

// generateFantasyTitle 生成奇幻风格标题
func (g *Generator) generateFantasyTitle(title, author string) string {
	return fmt.Sprintf(`
  <!-- 动态标题 - 奇幻风格 -->
  <g id="dynamic-title">
    <rect x="30" y="320" width="240" height="60" fill="#000000" opacity="0.5" rx="15"/>
    <rect x="30" y="320" width="240" height="60" fill="none" stroke="#ffffff" stroke-width="1" opacity="0.3" rx="15"/>
    <text x="150" y="345" text-anchor="middle" fill="#ffffff" font-family="serif" font-size="18" font-weight="bold">
      %s
    </text>
    %s
  </g>`, title, g.generateAuthorText(author, "150", "365", "#ecf0f1", "serif", "11"))
}

// generateScifiTitle 生成科幻风格标题
func (g *Generator) generateScifiTitle(title, author string) string {
	return fmt.Sprintf(`
  <!-- 动态标题 - 科幻风格 -->
  <g id="dynamic-title">
    <rect x="35" y="320" width="230" height="60" fill="#00ffff" opacity="0.1" rx="12"/>
    <rect x="35" y="320" width="230" height="60" fill="none" stroke="#00ffff" stroke-width="1" opacity="0.6" rx="12"/>
    <text x="150" y="345" text-anchor="middle" fill="#00ffff" font-family="monospace" font-size="16" font-weight="bold">
      %s
    </text>
    %s
  </g>`, title, g.generateAuthorText(author, "150", "365", "#00ffff", "monospace", "10"))
}

// generateClassicalTitle 生成古典风格标题
func (g *Generator) generateClassicalTitle(title, author string) string {
	return fmt.Sprintf(`
  <!-- 动态标题 - 古典风格 -->
  <g id="dynamic-title">
    <rect x="50" y="320" width="200" height="60" fill="#f5deb3" opacity="0.9" rx="8"/>
    <rect x="50" y="320" width="200" height="60" fill="none" stroke="#8b4513" stroke-width="2" rx="8"/>
    <text x="150" y="345" text-anchor="middle" fill="#8b4513" font-family="serif" font-size="17" font-weight="bold">
      %s
    </text>
    %s
  </g>`, title, g.generateAuthorText(author, "150", "365", "#a0522d", "serif", "11"))
}

// generateModernTitle 生成现代风格标题
func (g *Generator) generateModernTitle(title, author string) string {
	return fmt.Sprintf(`
  <!-- 动态标题 - 现代风格 -->
  <g id="dynamic-title">
    <rect x="40" y="320" width="220" height="60" fill="#ffffff" opacity="0.15" rx="20"/>
    <rect x="40" y="320" width="220" height="60" fill="none" stroke="#ffffff" stroke-width="1" opacity="0.4" rx="20"/>
    <text x="150" y="345" text-anchor="middle" fill="#ffffff" font-family="Arial, sans-serif" font-size="17" font-weight="300">
      %s
    </text>
    %s
  </g>`, title, g.generateAuthorText(author, "150", "365", "#ffffff", "Arial, sans-serif", "11"))
}

// generateAuthorText 生成作者文本
func (g *Generator) generateAuthorText(author, x, y, color, fontFamily, fontSize string) string {
	if author == "" {
		return ""
	}
	
	// 限制作者名长度
	displayAuthor := author
	if len([]rune(author)) > 15 {
		displayAuthor = string([]rune(author)[:15]) + "..."
	}
	
	return fmt.Sprintf(`
    <text x="%s" y="%s" text-anchor="middle" fill="%s" font-family="%s" font-size="%s" opacity="0.8">
      %s
    </text>`, x, y, color, fontFamily, fontSize, displayAuthor)
}
