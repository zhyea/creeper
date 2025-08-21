package parser

import (
	"fmt"
)

// ParseState 解析状态接口
type ParseState interface {
	HandleLine(context *ParseContext, line string, lineNumber int) error
	GetStateName() string
}

// ParseContext 解析上下文
type ParseContext struct {
	novel           *Novel
	currentChapter  *Chapter
	contentLines    []string
	chapterID       int
	volumeID        int
	sectionID       int
	state           ParseState
	txtFormat       *TxtFormat
	chapterTypes    []ChapterType
	inMetadata      bool
	metadataLines   int
}

// NewParseContext 创建解析上下文
func NewParseContext(novel *Novel) *ParseContext {
	ctx := &ParseContext{
		novel:        novel,
		contentLines: make([]string, 0),
		txtFormat:    NewTxtFormat(),
		chapterTypes: make([]ChapterType, 0),
	}
	
	// 设置初始状态
	ctx.SetState(NewMetadataState())
	
	return ctx
}

// SetState 设置状态
func (pc *ParseContext) SetState(state ParseState) {
	pc.state = state
}

// HandleLine 处理行
func (pc *ParseContext) HandleLine(line string, lineNumber int) error {
	return pc.state.HandleLine(pc, line, lineNumber)
}

// GetCurrentState 获取当前状态
func (pc *ParseContext) GetCurrentState() string {
	if pc.state != nil {
		return pc.state.GetStateName()
	}
	return "Unknown"
}

// SaveCurrentChapter 保存当前章节
func (pc *ParseContext) SaveCurrentChapter() {
	if pc.currentChapter != nil && len(pc.contentLines) > 0 {
		pc.currentChapter.Content = pc.txtFormat.CleanContent(strings.Join(pc.contentLines, "\n"))
		pc.currentChapter.WordCount = len([]rune(pc.currentChapter.Content))
		pc.novel.Chapters = append(pc.novel.Chapters, pc.currentChapter)
		pc.contentLines = make([]string, 0)
	}
}

// MetadataState 元数据状态
type MetadataState struct{}

func NewMetadataState() *MetadataState {
	return &MetadataState{}
}

func (ms *MetadataState) GetStateName() string {
	return "Metadata"
}

func (ms *MetadataState) HandleLine(context *ParseContext, line string, lineNumber int) error {
	line = strings.TrimSpace(line)
	
	// 空行跳过
	if line == "" {
		return nil
	}
	
	// 检查是否为分隔符，如果是则切换到内容状态
	if context.txtFormat.IsSeparator(line) {
		context.inMetadata = false
		context.metadataLines = lineNumber
		context.SetState(NewContentState())
		return nil
	}
	
	// 尝试提取元数据
	if key, value := context.txtFormat.ExtractMetadata(line); key != "" {
		switch key {
		case "title":
			context.novel.Title = value
		case "author":
			context.novel.Author = value
		case "description":
			context.novel.Description = value
		}
		return nil
	}
	
	// 如果不是元数据且超过一定行数，切换到内容状态
	if lineNumber > 20 {
		context.SetState(NewContentState())
		return context.state.HandleLine(context, line, lineNumber)
	}
	
	return nil
}

// ContentState 内容状态
type ContentState struct{}

func NewContentState() *ContentState {
	return &ContentState{}
}

func (cs *ContentState) GetStateName() string {
	return "Content"
}

func (cs *ContentState) HandleLine(context *ParseContext, line string, lineNumber int) error {
	line = strings.TrimSpace(line)
	
	// 检查分隔符
	if context.txtFormat.IsSeparator(line) {
		return nil
	}
	
	// 识别章节类型
	chapterType, title := context.txtFormat.IdentifyChapterType(line)
	
	if chapterType != ChapterTypeUnknown {
		// 保存上一章节
		context.SaveCurrentChapter()
		
		// 更新计数器
		switch chapterType {
		case ChapterTypeVolume:
			context.volumeID++
			context.chapterID = 0
			context.sectionID = 0
		case ChapterTypeChapter:
			context.chapterID++
			context.sectionID = 0
		case ChapterTypeSection:
			context.sectionID++
		}
		
		// 创建新章节
		context.currentChapter = &Chapter{
			ID:    len(context.novel.Chapters) + 1,
			Title: title,
			Path:  fmt.Sprintf("chapter-%d", len(context.novel.Chapters)+1),
		}
		
		// 记录章节类型
		context.chapterTypes = append(context.chapterTypes, chapterType)
		
		// 根据章节类型切换状态
		switch chapterType {
		case ChapterTypePrologue:
			context.SetState(NewPrologueState())
		case ChapterTypeEpilogue:
			context.SetState(NewEpilogueState())
		case ChapterTypeVolume:
			context.SetState(NewVolumeState())
		default:
			context.SetState(NewChapterState())
		}
		
		return nil
	}
	
	// 普通内容行
	if line != "" {
		context.contentLines = append(context.contentLines, line)
	} else {
		context.contentLines = append(context.contentLines, "")
	}
	
	return nil
}

// ChapterState 章节状态
type ChapterState struct{}

func NewChapterState() *ChapterState {
	return &ChapterState{}
}

func (cs *ChapterState) GetStateName() string {
	return "Chapter"
}

func (cs *ChapterState) HandleLine(context *ParseContext, line string, lineNumber int) error {
	// 章节状态下的内容处理
	return NewContentState().HandleLine(context, line, lineNumber)
}

// VolumeState 卷状态
type VolumeState struct{}

func NewVolumeState() *VolumeState {
	return &VolumeState{}
}

func (vs *VolumeState) GetStateName() string {
	return "Volume"
}

func (vs *VolumeState) HandleLine(context *ParseContext, line string, lineNumber int) error {
	// 卷状态下的内容处理
	return NewContentState().HandleLine(context, line, lineNumber)
}

// PrologueState 序言状态
type PrologueState struct{}

func NewPrologueState() *PrologueState {
	return &PrologueState{}
}

func (ps *PrologueState) GetStateName() string {
	return "Prologue"
}

func (ps *PrologueState) HandleLine(context *ParseContext, line string, lineNumber int) error {
	// 序言状态下的内容处理
	return NewContentState().HandleLine(context, line, lineNumber)
}

// EpilogueState 后记状态
type EpilogueState struct{}

func NewEpilogueState() *EpilogueState {
	return &EpilogueState{}
}

func (es *EpilogueState) GetStateName() string {
	return "Epilogue"
}

func (es *EpilogueState) HandleLine(context *ParseContext, line string, lineNumber int) error {
	// 后记状态下的内容处理
	return NewContentState().HandleLine(context, line, lineNumber)
}

// StatefulTxtParser 状态化 TXT 解析器
type StatefulTxtParser struct {
	txtFormat *TxtFormat
	notifier  *ParseNotifier
}

// NewStatefulTxtParser 创建状态化 TXT 解析器
func NewStatefulTxtParser() *StatefulTxtParser {
	return &StatefulTxtParser{
		txtFormat: NewTxtFormat(),
		notifier:  NewParseNotifier(),
	}
}

// Subscribe 订阅解析事件
func (stp *StatefulTxtParser) Subscribe(observer ParseObserver) {
	stp.notifier.Subscribe(observer)
}

// ParseWithState 使用状态模式解析
func (stp *StatefulTxtParser) ParseWithState(novel *Novel, lines []string) error {
	context := NewParseContext(novel)
	
	stp.notifier.NotifyObservers(&ParseEventData{
		Event:   ParseEventStart,
		Message: fmt.Sprintf("开始状态化解析: %s", novel.Title),
	})
	
	totalLines := len(lines)
	
	for i, line := range lines {
		if err := context.HandleLine(line, i); err != nil {
			stp.notifier.NotifyObservers(&ParseEventData{
				Event:   ParseEventError,
				Message: fmt.Sprintf("解析第%d行失败", i+1),
				Error:   err,
			})
			return err
		}
		
		// 报告进度
		if i%100 == 0 || i == totalLines-1 {
			progress := float64(i+1) / float64(totalLines) * 100
			stp.notifier.NotifyObservers(&ParseEventData{
				Event:    ParseEventProgress,
				Message:  fmt.Sprintf("已处理 %d/%d 行", i+1, totalLines),
				Progress: progress,
			})
		}
	}
	
	// 保存最后一章节
	context.SaveCurrentChapter()
	
	stp.notifier.NotifyObservers(&ParseEventData{
		Event:   ParseEventComplete,
		Message: fmt.Sprintf("解析完成: 共%d章", len(novel.Chapters)),
	})
	
	return nil
}
