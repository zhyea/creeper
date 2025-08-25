package generator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"creeper/internal/config"
	"creeper/internal/parser"
)

// Generator 静态站点生成器
type Generator struct {
	config   *config.Config
	parser   *parser.Parser
	novels   []*parser.Novel
	templates map[string]*template.Template
}

// New 创建新的生成器
func New(cfg *config.Config) *Generator {
	return &Generator{
		config:    cfg,
		parser:    parser.New(),
		novels:    make([]*parser.Novel, 0),
		templates: make(map[string]*template.Template),
	}
}

// Generate 生成静态站点
func (g *Generator) Generate() error {
	// 1. 解析所有小说
	if err := g.parseNovels(); err != nil {
		return fmt.Errorf("解析小说失败: %v", err)
	}

	// 2. 加载模板
	if err := g.loadTemplates(); err != nil {
		return fmt.Errorf("加载模板失败: %v", err)
	}

	// 3. 创建输出目录
	if err := g.createOutputDir(); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 4. 生成静态资源
	if err := g.generateAssets(); err != nil {
		return fmt.Errorf("生成静态资源失败: %v", err)
	}

	// 5. 生成首页
	if err := g.generateIndex(); err != nil {
		return fmt.Errorf("生成首页失败: %v", err)
	}

	// 6. 生成小说页面
	for _, novel := range g.novels {
		// 生成带标题的封面
		if err := g.generateNovelCover(novel); err != nil {
			fmt.Printf("警告：生成小说 %s 的封面失败: %v\n", novel.Title, err)
		}
		
		if err := g.generateNovel(novel); err != nil {
			return fmt.Errorf("生成小说 %s 失败: %v", novel.Title, err)
		}
	}

	// 7. 生成搜索数据
	if err := g.generateSearchData(); err != nil {
		return fmt.Errorf("生成搜索数据失败: %v", err)
	}

	// 8. 生成分类页面
	if err := g.generateCategoryPages(); err != nil {
		return fmt.Errorf("生成分类页面失败: %v", err)
	}

	// 9. 生成作者页面
	if err := g.generateAuthorPages(); err != nil {
		return fmt.Errorf("生成作者页面失败: %v", err)
	}

	return nil
}

// parseNovels 解析所有小说
func (g *Generator) parseNovels() error {
	inputDir := g.config.InputDir

	// 检查输入目录是否存在
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return fmt.Errorf("输入目录不存在: %s", inputDir)
	}

	// 遍历输入目录
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("读取输入目录失败: %v", err)
	}

	for _, entry := range entries {
		path := filepath.Join(inputDir, entry.Name())

		// 跳过隐藏文件和目录
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		var novel *parser.Novel
		if entry.IsDir() {
			// 目录模式
			novel, err = g.parser.ParseNovel(path)
		} else if strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			// 单文件模式
			novel, err = g.parser.ParseNovel(path)
		} else {
			continue
		}

		if err != nil {
			fmt.Printf("警告：解析 %s 失败: %v\n", path, err)
			continue
		}

		if len(novel.Chapters) > 0 {
			g.novels = append(g.novels, novel)
		}
	}

	// 按标题排序
	sort.Slice(g.novels, func(i, j int) bool {
		return g.novels[i].Title < g.novels[j].Title
	})

	fmt.Printf("成功解析 %d 部小说\n", len(g.novels))
	return nil
}

// createOutputDir 创建输出目录
func (g *Generator) createOutputDir() error {
	outputDir := g.config.OutputDir

	// 如果目录存在，先删除
	if _, err := os.Stat(outputDir); err == nil {
		if err := os.RemoveAll(outputDir); err != nil {
			return fmt.Errorf("删除旧输出目录失败: %v", err)
		}
	}

	// 创建目录结构
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "novels"),
		filepath.Join(outputDir, "static"),
		filepath.Join(outputDir, "static", "css"),
		filepath.Join(outputDir, "static", "js"),
		filepath.Join(outputDir, "static", "images"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}
	}

	return nil
}

// generateIndex 生成首页
func (g *Generator) generateIndex() error {
	data := map[string]interface{}{
		"Config": g.config,
		"Novels": g.novels,
		"Title":  g.config.Site.Title,
	}

	return g.renderTemplate("index", "index.html", data)
}

// generateNovel 生成小说页面
func (g *Generator) generateNovel(novel *parser.Novel) error {
	novelDir := filepath.Join(g.config.OutputDir, "novels", g.sanitizeFileName(novel.Title))
	if err := os.MkdirAll(novelDir, 0755); err != nil {
		return fmt.Errorf("创建小说目录失败: %v", err)
	}

	// 生成小说目录页
	data := map[string]interface{}{
		"Config": g.config,
		"Novel":  novel,
		"Title":  novel.Title,
	}

	indexPath := filepath.Join(novelDir, "index.html")
	if err := g.renderTemplateToFile("novel", indexPath, data); err != nil {
		return fmt.Errorf("生成小说目录页失败: %v", err)
	}

	// 生成每个章节页面
	for _, chapter := range novel.Chapters {
		chapterData := map[string]interface{}{
			"Config":  g.config,
			"Novel":   novel,
			"Chapter": chapter,
			"Title":   fmt.Sprintf("%s - %s", chapter.Title, novel.Title),
		}

		chapterPath := filepath.Join(novelDir, fmt.Sprintf("chapter-%d.html", chapter.ID))
		if err := g.renderTemplateToFile("chapter", chapterPath, chapterData); err != nil {
			return fmt.Errorf("生成章节 %d 失败: %v", chapter.ID, err)
		}
	}

	return nil
}

// generateSearchData 生成搜索数据
func (g *Generator) generateSearchData() error {
	searchData := make([]map[string]interface{}, 0)

	for _, novel := range g.novels {
		novelData := map[string]interface{}{
			"type":        "novel",
			"title":       novel.Title,
			"author":      novel.Author,
			"description": novel.Description,
			"url":         fmt.Sprintf("novels/%s/", g.sanitizeFileName(novel.Title)),
		}
		searchData = append(searchData, novelData)

		for _, chapter := range novel.Chapters {
			chapterData := map[string]interface{}{
				"type":   "chapter",
				"title":  chapter.Title,
				"novel":  novel.Title,
				"author": novel.Author,
				"url":    fmt.Sprintf("novels/%s/chapter-%d.html", g.sanitizeFileName(novel.Title), chapter.ID),
			}
			searchData = append(searchData, chapterData)
		}
	}

	data, err := json.MarshalIndent(searchData, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化搜索数据失败: %v", err)
	}

	searchPath := filepath.Join(g.config.OutputDir, "static", "js", "search-data.json")
	return os.WriteFile(searchPath, data, 0644)
}

// renderTemplate 渲染模板到默认位置
func (g *Generator) renderTemplate(templateName, fileName string, data interface{}) error {
	outputPath := filepath.Join(g.config.OutputDir, fileName)
	return g.renderTemplateToFile(templateName, outputPath, data)
}

// renderTemplateToFile 渲染模板到指定文件
func (g *Generator) renderTemplateToFile(templateName, outputPath string, data interface{}) error {
	tmpl, exists := g.templates[templateName]
	if !exists {
		return fmt.Errorf("模板 %s 不存在", templateName)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件 %s 失败: %v", outputPath, err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// sanitizeFileName 清理文件名
func (g *Generator) sanitizeFileName(name string) string {
	// 替换不安全的字符
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "*", "-")
	name = strings.ReplaceAll(name, "?", "-")
	name = strings.ReplaceAll(name, "\"", "-")
	name = strings.ReplaceAll(name, "<", "-")
	name = strings.ReplaceAll(name, ">", "-")
	name = strings.ReplaceAll(name, "|", "-")
	return strings.TrimSpace(name)
}

// Serve 启动本地服务器
func (g *Generator) Serve(port int) error {
	handler := http.FileServer(http.Dir(g.config.OutputDir))
	http.Handle("/", handler)

	fmt.Printf("服务器运行在: http://localhost:%d\n", port)
	fmt.Printf("按 Ctrl+C 停止服务器\n")

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// generateCategoryPages 生成分类页面
func (g *Generator) generateCategoryPages() error {
	// 按分类组织小说
	categoryMap := make(map[string][]*parser.Novel)
	
	for _, novel := range g.novels {
		category := novel.Category
		if category == "" {
			category = "未分类"
		}
		categoryMap[category] = append(categoryMap[category], novel)
	}

	// 生成分类列表页面
	categories := make([]map[string]interface{}, 0)
	for category, novels := range categoryMap {
		categoryData := map[string]interface{}{
			"name":        category,
			"count":       len(novels),
			"description": g.getCategoryDescription(category),
			"color":       g.getCategoryColor(category),
			"icon":        g.getCategoryIcon(category),
		}
		categories = append(categories, categoryData)
	}

	// 生成分类列表页面
	categoryListData := map[string]interface{}{
		"Config":      g.config,
		"Categories":  categories,
		"Title":       "分类浏览",
		"Description": "按分类浏览所有小说",
	}

	if err := g.renderTemplate("category-list", filepath.Join(g.config.OutputDir, "categories.html"), categoryListData); err != nil {
		return fmt.Errorf("生成分类列表页面失败: %v", err)
	}

	// 生成每个分类的详情页面
	for category, novels := range categoryMap {
		categoryData := map[string]interface{}{
			"Config":      g.config,
			"Category":    category,
			"Novels":      novels,
			"Count":       len(novels),
			"Description": g.getCategoryDescription(category),
			"Color":       g.getCategoryColor(category),
			"Icon":        g.getCategoryIcon(category),
			"Title":       fmt.Sprintf("%s - 分类浏览", category),
		}

		categoryPath := filepath.Join(g.config.OutputDir, "categories", fmt.Sprintf("%s.html", g.sanitizeFileName(category)))
		
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(categoryPath), 0755); err != nil {
			return fmt.Errorf("创建分类目录失败: %v", err)
		}

		if err := g.renderTemplate("category", categoryPath, categoryData); err != nil {
			return fmt.Errorf("生成分类 %s 页面失败: %v", category, err)
		}
	}

	return nil
}

// generateAuthorPages 生成作者页面
func (g *Generator) generateAuthorPages() error {
	// 按作者组织小说
	authorMap := make(map[string][]*parser.Novel)
	
	for _, novel := range g.novels {
		author := novel.Author
		if author == "" {
			author = "未知作者"
		}
		authorMap[author] = append(authorMap[author], novel)
	}

	// 生成作者列表页面
	authors := make([]map[string]interface{}, 0)
	for author, novels := range authorMap {
		authorData := map[string]interface{}{
			"name":        author,
			"count":       len(novels),
			"totalWords":  g.calculateTotalWords(novels),
			"lastUpdated": g.getLastUpdated(novels),
		}
		authors = append(authors, authorData)
	}

	// 生成作者列表页面
	authorListData := map[string]interface{}{
		"Config":      g.config,
		"Authors":     authors,
		"Title":       "作者作品",
		"Description": "按作者浏览所有作品",
	}

	if err := g.renderTemplate("author-list", filepath.Join(g.config.OutputDir, "authors.html"), authorListData); err != nil {
		return fmt.Errorf("生成作者列表页面失败: %v", err)
	}

	// 生成每个作者的详情页面
	for author, novels := range authorMap {
		authorData := map[string]interface{}{
			"Config":      g.config,
			"Author":      author,
			"Novels":      novels,
			"Count":       len(novels),
			"TotalWords":  g.calculateTotalWords(novels),
			"LastUpdated": g.getLastUpdated(novels),
			"Title":       fmt.Sprintf("%s - 作者作品", author),
		}

		authorPath := filepath.Join(g.config.OutputDir, "authors", fmt.Sprintf("%s.html", g.sanitizeFileName(author)))
		
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(authorPath), 0755); err != nil {
			return fmt.Errorf("创建作者目录失败: %v", err)
		}

		if err := g.renderTemplate("author", authorPath, authorData); err != nil {
			return fmt.Errorf("生成作者 %s 页面失败: %v", author, err)
		}
	}

	return nil
}

// getCategoryDescription 获取分类描述
func (g *Generator) getCategoryDescription(category string) string {
	descriptions := map[string]string{
		"科幻":     "探索未来科技与宇宙奥秘的科幻小说",
		"现代":     "反映现代都市生活的现实主义小说",
		"武侠":     "传统武侠文化与江湖恩怨的武侠小说",
		"玄幻":     "融合东方玄幻元素的奇幻小说",
		"历史":     "基于历史背景的历史小说",
		"悬疑":     "充满谜团与推理的悬疑小说",
		"言情":     "浪漫爱情故事",
		"未分类":   "暂无分类的小说",
	}
	return descriptions[category]
}

// getCategoryColor 获取分类颜色
func (g *Generator) getCategoryColor(category string) string {
	colors := map[string]string{
		"科幻":   "#3498db",
		"现代":   "#2ecc71",
		"武侠":   "#e74c3c",
		"玄幻":   "#9b59b6",
		"历史":   "#f39c12",
		"悬疑":   "#34495e",
		"言情":   "#e91e63",
		"未分类": "#95a5a6",
	}
	return colors[category]
}

// getCategoryIcon 获取分类图标
func (g *Generator) getCategoryIcon(category string) string {
	icons := map[string]string{
		"科幻":   "🚀",
		"现代":   "🏢",
		"武侠":   "⚔️",
		"玄幻":   "🐉",
		"历史":   "📚",
		"悬疑":   "🔍",
		"言情":   "💕",
		"未分类": "📖",
	}
	return icons[category]
}

// calculateTotalWords 计算总字数
func (g *Generator) calculateTotalWords(novels []*parser.Novel) int {
	total := 0
	for _, novel := range novels {
		for _, chapter := range novel.Chapters {
			total += chapter.WordCount
		}
	}
	return total
}

// getLastUpdated 获取最后更新时间
func (g *Generator) getLastUpdated(novels []*parser.Novel) string {
	if len(novels) == 0 {
		return ""
	}
	
	latest := novels[0].UpdatedAt
	for _, novel := range novels {
		if novel.UpdatedAt.After(latest) {
			latest = novel.UpdatedAt
		}
	}
	
	return latest.Format("2006-01-02")
}
