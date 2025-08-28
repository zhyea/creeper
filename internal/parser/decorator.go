package parser

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ParserDecorator 解析器装饰器接口
type ParserDecorator interface {
	ParseNovel(path string) (*Novel, error)
}

// BaseParserDecorator 基础解析器装饰器
type BaseParserDecorator struct {
	parser *Parser
}

// NewBaseParserDecorator 创建基础装饰器
func NewBaseParserDecorator(parser *Parser) *BaseParserDecorator {
	return &BaseParserDecorator{parser: parser}
}

func (bpd *BaseParserDecorator) ParseNovel(path string) (*Novel, error) {
	return bpd.parser.ParseNovel(path)
}

// LoggingParserDecorator 日志装饰器
type LoggingParserDecorator struct {
	ParserDecorator
	observer ParseObserver
}

// NewLoggingParserDecorator 创建日志装饰器
func NewLoggingParserDecorator(decorator ParserDecorator, observer ParseObserver) *LoggingParserDecorator {
	return &LoggingParserDecorator{
		ParserDecorator: decorator,
		observer:        observer,
	}
}

func (lpd *LoggingParserDecorator) ParseNovel(path string) (*Novel, error) {
	start := time.Now()

	if lpd.observer != nil {
		lpd.observer.OnParseEvent(&ParseEventData{
			Event:   ParseEventStart,
			Message: fmt.Sprintf("开始解析小说: %s", path),
		})
	}

	novel, err := lpd.ParserDecorator.ParseNovel(path)

	duration := time.Since(start)

	if err != nil {
		if lpd.observer != nil {
			lpd.observer.OnParseEvent(&ParseEventData{
				Event:   ParseEventError,
				Message: fmt.Sprintf("解析失败: %s (耗时: %v)", path, duration),
				Error:   err,
			})
		}
		return nil, err
	}

	if lpd.observer != nil {
		lpd.observer.OnParseEvent(&ParseEventData{
			Event:   ParseEventComplete,
			Message: fmt.Sprintf("解析完成: %s (耗时: %v, %d章)", novel.Title, duration, len(novel.Chapters)),
		})
	}

	return novel, nil
}

// CachingParserDecorator 缓存装饰器
type CachingParserDecorator struct {
	ParserDecorator
	cache map[string]*CachedNovel
}

// CachedNovel 缓存的小说
type CachedNovel struct {
	Novel     *Novel
	Timestamp time.Time
	FileInfo  FileInfo
}

// FileInfo 文件信息
type FileInfo struct {
	Path    string
	ModTime time.Time
	Size    int64
}

// NewCachingParserDecorator 创建缓存装饰器
func NewCachingParserDecorator(decorator ParserDecorator) *CachingParserDecorator {
	return &CachingParserDecorator{
		ParserDecorator: decorator,
		cache:           make(map[string]*CachedNovel),
	}
}

func (cpd *CachingParserDecorator) ParseNovel(path string) (*Novel, error) {
	// 获取文件信息
	fileInfo, err := cpd.getFileInfo(path)
	if err != nil {
		return nil, err
	}

	// 检查缓存
	if cached, exists := cpd.cache[path]; exists {
		if cached.FileInfo.ModTime.Equal(fileInfo.ModTime) && cached.FileInfo.Size == fileInfo.Size {
			// 缓存有效，返回缓存的结果
			return cpd.cloneNovel(cached.Novel), nil
		}
	}

	// 缓存无效或不存在，重新解析
	novel, err := cpd.ParserDecorator.ParseNovel(path)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	cpd.cache[path] = &CachedNovel{
		Novel:     cpd.cloneNovel(novel),
		Timestamp: time.Now(),
		FileInfo:  fileInfo,
	}

	return novel, nil
}

// getFileInfo 获取文件信息
func (cpd *CachingParserDecorator) getFileInfo(path string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Path:    path,
		ModTime: info.ModTime(),
		Size:    info.Size(),
	}, nil
}

// cloneNovel 克隆小说对象
func (cpd *CachingParserDecorator) cloneNovel(original *Novel) *Novel {
	clone := &Novel{
		Title:       original.Title,
		Author:      original.Author,
		Description: original.Description,
		Cover:       original.Cover,
		CreatedAt:   original.CreatedAt,
		UpdatedAt:   original.UpdatedAt,
		Path:        original.Path,
		Chapters:    make([]*Chapter, len(original.Chapters)),
	}

	for i, chapter := range original.Chapters {
		clone.Chapters[i] = &Chapter{
			ID:          chapter.ID,
			Title:       chapter.Title,
			Content:     chapter.Content,
			HTMLContent: chapter.HTMLContent,
			WordCount:   chapter.WordCount,
			CreatedAt:   chapter.CreatedAt,
			Path:        chapter.Path,
		}
	}

	return clone
}

// ClearCache 清空缓存
func (cpd *CachingParserDecorator) ClearCache() {
	cpd.cache = make(map[string]*CachedNovel)
}

// GetCacheStats 获取缓存统计
func (cpd *CachingParserDecorator) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cached_novels": len(cpd.cache),
		"cache_keys":    cpd.getCacheKeys(),
	}
}

// getCacheKeys 获取缓存键列表
func (cpd *CachingParserDecorator) getCacheKeys() []string {
	keys := make([]string, 0, len(cpd.cache))
	for key := range cpd.cache {
		keys = append(keys, key)
	}
	return keys
}

// ValidationParserDecorator 验证装饰器
type ValidationParserDecorator struct {
	ParserDecorator
	validator *ValidationVisitor
}

// NewValidationParserDecorator 创建验证装饰器
func NewValidationParserDecorator(decorator ParserDecorator) *ValidationParserDecorator {
	return &ValidationParserDecorator{
		ParserDecorator: decorator,
		validator:       NewValidationVisitor(),
	}
}

func (vpd *ValidationParserDecorator) ParseNovel(path string) (*Novel, error) {
	novel, err := vpd.ParserDecorator.ParseNovel(path)
	if err != nil {
		return nil, err
	}

	// 验证小说内容
	vpd.validator.ClearErrors()

	for _, chapter := range novel.Chapters {
		// 根据章节标题推断类型
		chapterType := vpd.inferChapterType(chapter.Title)
		element := NewStandardChapter(chapter, chapterType)

		if err := element.Accept(vpd.validator); err != nil {
			return nil, fmt.Errorf("验证章节失败: %w", err)
		}
	}

	// 检查验证错误
	if vpd.validator.HasErrors() {
		errors := vpd.validator.GetErrors()
		fmt.Printf("⚠️  发现 %d 个验证问题:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("   - %s\n", err)
		}
	}

	return novel, nil
}

// inferChapterType 推断章节类型
func (vpd *ValidationParserDecorator) inferChapterType(title string) ChapterType {
	title = strings.ToLower(title)

	if strings.Contains(title, "序") || strings.Contains(title, "楔子") || strings.Contains(title, "引子") {
		return ChapterTypePrologue
	}

	if strings.Contains(title, "后记") || strings.Contains(title, "尾声") || strings.Contains(title, "结语") {
		return ChapterTypeEpilogue
	}

	if strings.Contains(title, "卷") || strings.Contains(title, "volume") {
		return ChapterTypeVolume
	}

	if strings.Contains(title, "节") || strings.Contains(title, "section") {
		return ChapterTypeSection
	}

	return ChapterTypeChapter
}

// EnhancedParser 增强解析器
type EnhancedParser struct {
	baseParser     *Parser
	decorators     []func(ParserDecorator) ParserDecorator
	finalDecorator ParserDecorator
}

// NewEnhancedParser 创建增强解析器
func NewEnhancedParser() *EnhancedParser {
	baseParser := New()

	return &EnhancedParser{
		baseParser: baseParser,
		decorators: make([]func(ParserDecorator) ParserDecorator, 0),
	}
}

// WithLogging 添加日志装饰器
func (ep *EnhancedParser) WithLogging(observer ParseObserver) *EnhancedParser {
	ep.decorators = append(ep.decorators, func(decorator ParserDecorator) ParserDecorator {
		return NewLoggingParserDecorator(decorator, observer)
	})
	return ep
}

// WithCaching 添加缓存装饰器
func (ep *EnhancedParser) WithCaching() *EnhancedParser {
	ep.decorators = append(ep.decorators, func(decorator ParserDecorator) ParserDecorator {
		return NewCachingParserDecorator(decorator)
	})
	return ep
}

// WithValidation 添加验证装饰器
func (ep *EnhancedParser) WithValidation() *EnhancedParser {
	ep.decorators = append(ep.decorators, func(decorator ParserDecorator) ParserDecorator {
		return NewValidationParserDecorator(decorator)
	})
	return ep
}

// Build 构建最终的解析器
func (ep *EnhancedParser) Build() ParserDecorator {
	if ep.finalDecorator != nil {
		return ep.finalDecorator
	}

	// 从基础解析器开始
	var decorator ParserDecorator = NewBaseParserDecorator(ep.baseParser)

	// 应用所有装饰器
	for _, decoratorFunc := range ep.decorators {
		decorator = decoratorFunc(decorator)
	}

	ep.finalDecorator = decorator
	return decorator
}

// ParseNovel 解析小说（便捷方法）
func (ep *EnhancedParser) ParseNovel(path string) (*Novel, error) {
	decorator := ep.Build()
	return decorator.ParseNovel(path)
}
