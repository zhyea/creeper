package parser

import (
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

// ContentAdapter 内容适配器接口
type ContentAdapter interface {
	ConvertToHTML(content string) string
	GetContentType() string
	PreprocessContent(content string) string
	PostprocessContent(content string) string
}

// MarkdownAdapter Markdown 内容适配器
type MarkdownAdapter struct{}

// NewMarkdownAdapter 创建 Markdown 适配器
func NewMarkdownAdapter() *MarkdownAdapter {
	return &MarkdownAdapter{}
}

func (ma *MarkdownAdapter) GetContentType() string {
	return "markdown"
}

func (ma *MarkdownAdapter) PreprocessContent(content string) string {
	// Markdown 预处理
	content = strings.TrimSpace(content)
	
	// 确保段落之间有适当的空行
	lines := strings.Split(content, "\n")
	var processedLines []string
	
	for i, line := range lines {
		processedLines = append(processedLines, line)
		
		// 在非空行后面如果跟着非空行，且不是列表或标题，添加空行
		if i < len(lines)-1 && strings.TrimSpace(line) != "" && 
		   strings.TrimSpace(lines[i+1]) != "" &&
		   !strings.HasPrefix(strings.TrimSpace(lines[i+1]), "#") &&
		   !strings.HasPrefix(strings.TrimSpace(lines[i+1]), "-") &&
		   !strings.HasPrefix(strings.TrimSpace(lines[i+1]), "*") {
			// 检查是否已经有空行
			if i < len(lines)-2 && strings.TrimSpace(lines[i+1]) != "" {
				processedLines = append(processedLines, "")
			}
		}
	}
	
	return strings.Join(processedLines, "\n")
}

func (ma *MarkdownAdapter) ConvertToHTML(content string) string {
	preprocessed := ma.PreprocessContent(content)
	html := string(blackfriday.Run([]byte(preprocessed)))
	return ma.PostprocessContent(html)
}

func (ma *MarkdownAdapter) PostprocessContent(content string) string {
	// Markdown 后处理
	// 可以在这里添加自定义的 HTML 处理逻辑
	return content
}

// TxtAdapter TXT 内容适配器
type TxtAdapter struct{}

// NewTxtAdapter 创建 TXT 适配器
func NewTxtAdapter() *TxtAdapter {
	return &TxtAdapter{}
}

func (ta *TxtAdapter) GetContentType() string {
	return "txt"
}

func (ta *TxtAdapter) PreprocessContent(content string) string {
	// TXT 预处理
	content = strings.TrimSpace(content)
	
	// 将 TXT 内容转换为类似 Markdown 的格式
	lines := strings.Split(content, "\n")
	var processedLines []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if trimmed == "" {
			processedLines = append(processedLines, "")
			continue
		}
		
		// 处理特殊格式
		if ta.isDialogue(trimmed) {
			// 对话处理
			processedLines = append(processedLines, "> "+trimmed)
		} else if ta.isEmphasis(trimmed) {
			// 强调处理
			processedLines = append(processedLines, "**"+trimmed+"**")
		} else {
			processedLines = append(processedLines, trimmed)
		}
	}
	
	return strings.Join(processedLines, "\n\n")
}

func (ta *TxtAdapter) ConvertToHTML(content string) string {
	preprocessed := ta.PreprocessContent(content)
	
	// 将预处理后的内容转换为 Markdown，然后转换为 HTML
	html := string(blackfriday.Run([]byte(preprocessed)))
	
	return ta.PostprocessContent(html)
}

func (ta *TxtAdapter) PostprocessContent(content string) string {
	// TXT 后处理
	// 添加段落缩进样式
	content = strings.ReplaceAll(content, "<p>", `<p class="txt-paragraph">`)
	
	// 处理对话样式
	content = strings.ReplaceAll(content, "<blockquote>", `<blockquote class="dialogue">`)
	
	return content
}

// isDialogue 判断是否为对话
func (ta *TxtAdapter) isDialogue(line string) bool {
	return strings.HasPrefix(line, "\"") || strings.HasPrefix(line, "'")
}

// isEmphasis 判断是否需要强调
func (ta *TxtAdapter) isEmphasis(line string) bool {
	// 判断是否为心理描写或重要内容
	return strings.Contains(line, "心想") || strings.Contains(line, "暗道") ||
		   strings.Contains(line, "！！") || strings.Contains(line, "？？")
}

// ContentAdapterFactory 内容适配器工厂
type ContentAdapterFactory struct {
	adapters map[string]ContentAdapter
}

// NewContentAdapterFactory 创建内容适配器工厂
func NewContentAdapterFactory() *ContentAdapterFactory {
	factory := &ContentAdapterFactory{
		adapters: make(map[string]ContentAdapter),
	}
	
	// 注册默认适配器
	factory.RegisterAdapter("md", NewMarkdownAdapter())
	factory.RegisterAdapter("markdown", NewMarkdownAdapter())
	factory.RegisterAdapter("txt", NewTxtAdapter())
	factory.RegisterAdapter("text", NewTxtAdapter())
	
	return factory
}

// RegisterAdapter 注册适配器
func (caf *ContentAdapterFactory) RegisterAdapter(fileType string, adapter ContentAdapter) {
	caf.adapters[strings.ToLower(fileType)] = adapter
}

// GetAdapter 获取适配器
func (caf *ContentAdapterFactory) GetAdapter(fileType string) ContentAdapter {
	if adapter, exists := caf.adapters[strings.ToLower(fileType)]; exists {
		return adapter
	}
	// 默认返回 TXT 适配器
	return NewTxtAdapter()
}

// GetSupportedTypes 获取支持的文件类型
func (caf *ContentAdapterFactory) GetSupportedTypes() []string {
	types := make([]string, 0, len(caf.adapters))
	for fileType := range caf.adapters {
		types = append(types, fileType)
	}
	return types
}

// ChapterAdapter 章节适配器
type ChapterAdapter struct {
	contentFactory *ContentAdapterFactory
	notifier       *ParseNotifier
}

// NewChapterAdapter 创建章节适配器
func NewChapterAdapter() *ChapterAdapter {
	return &ChapterAdapter{
		contentFactory: NewContentAdapterFactory(),
		notifier:       NewParseNotifier(),
	}
}

// Subscribe 订阅解析事件
func (ca *ChapterAdapter) Subscribe(observer ParseObserver) {
	ca.notifier.Subscribe(observer)
}

// ConvertChapter 转换章节内容
func (ca *ChapterAdapter) ConvertChapter(chapter *Chapter, sourceType string) error {
	ca.notifier.NotifyObservers(&ParseEventData{
		Event:       ParseEventStart,
		Message:     fmt.Sprintf("开始转换章节: %s", chapter.Title),
		ChapterInfo: chapter,
	})
	
	// 获取适配器
	adapter := ca.contentFactory.GetAdapter(sourceType)
	
	// 转换内容
	chapter.HTMLContent = adapter.ConvertToHTML(chapter.Content)
	chapter.WordCount = len([]rune(chapter.Content))
	chapter.CreatedAt = time.Now()
	
	ca.notifier.NotifyObservers(&ParseEventData{
		Event:       ParseEventComplete,
		Message:     fmt.Sprintf("章节转换完成: %s", chapter.Title),
		ChapterInfo: chapter,
	})
	
	return nil
}

// ConvertNovel 转换整部小说
func (ca *ChapterAdapter) ConvertNovel(novel *Novel, sourceType string) error {
	totalChapters := len(novel.Chapters)
	
	ca.notifier.NotifyObservers(&ParseEventData{
		Event:   ParseEventStart,
		Message: fmt.Sprintf("开始转换小说: %s (%d章)", novel.Title, totalChapters),
	})
	
	for i, chapter := range novel.Chapters {
		if err := ca.ConvertChapter(chapter, sourceType); err != nil {
			ca.notifier.NotifyObservers(&ParseEventData{
				Event:   ParseEventError,
				Message: fmt.Sprintf("转换章节失败: %s", chapter.Title),
				Error:   err,
			})
			return err
		}
		
		// 更新进度
		progress := float64(i+1) / float64(totalChapters) * 100
		ca.notifier.NotifyObservers(&ParseEventData{
			Event:    ParseEventProgress,
			Message:  fmt.Sprintf("已完成 %d/%d 章", i+1, totalChapters),
			Progress: progress,
		})
	}
	
	ca.notifier.NotifyObservers(&ParseEventData{
		Event:   ParseEventComplete,
		Message: fmt.Sprintf("小说转换完成: %s", novel.Title),
	})
	
	return nil
}
