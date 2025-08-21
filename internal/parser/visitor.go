package parser

import (
	"fmt"
	"strings"
)

// ChapterVisitor 章节访问者接口
type ChapterVisitor interface {
	VisitPrologue(chapter *Chapter) error
	VisitEpilogue(chapter *Chapter) error
	VisitRegularChapter(chapter *Chapter) error
	VisitVolumeChapter(chapter *Chapter) error
	VisitSectionChapter(chapter *Chapter) error
}

// ChapterElement 章节元素接口
type ChapterElement interface {
	Accept(visitor ChapterVisitor) error
	GetChapterType() ChapterType
}

// StandardChapter 标准章节实现
type StandardChapter struct {
	*Chapter
	chapterType ChapterType
}

// NewStandardChapter 创建标准章节
func NewStandardChapter(chapter *Chapter, chapterType ChapterType) *StandardChapter {
	return &StandardChapter{
		Chapter:     chapter,
		chapterType: chapterType,
	}
}

// Accept 接受访问者
func (sc *StandardChapter) Accept(visitor ChapterVisitor) error {
	switch sc.chapterType {
	case ChapterTypePrologue:
		return visitor.VisitPrologue(sc.Chapter)
	case ChapterTypeEpilogue:
		return visitor.VisitEpilogue(sc.Chapter)
	case ChapterTypeVolume:
		return visitor.VisitVolumeChapter(sc.Chapter)
	case ChapterTypeSection:
		return visitor.VisitSectionChapter(sc.Chapter)
	default:
		return visitor.VisitRegularChapter(sc.Chapter)
	}
}

// GetChapterType 获取章节类型
func (sc *StandardChapter) GetChapterType() ChapterType {
	return sc.chapterType
}

// HTMLGeneratorVisitor HTML 生成访问者
type HTMLGeneratorVisitor struct {
	contentAdapter *ChapterAdapter
}

// NewHTMLGeneratorVisitor 创建 HTML 生成访问者
func NewHTMLGeneratorVisitor(contentAdapter *ChapterAdapter) *HTMLGeneratorVisitor {
	return &HTMLGeneratorVisitor{
		contentAdapter: contentAdapter,
	}
}

func (hgv *HTMLGeneratorVisitor) VisitPrologue(chapter *Chapter) error {
	// 序言特殊处理
	chapter.HTMLContent = hgv.wrapWithClass(chapter.HTMLContent, "prologue-content")
	return nil
}

func (hgv *HTMLGeneratorVisitor) VisitEpilogue(chapter *Chapter) error {
	// 后记特殊处理
	chapter.HTMLContent = hgv.wrapWithClass(chapter.HTMLContent, "epilogue-content")
	return nil
}

func (hgv *HTMLGeneratorVisitor) VisitRegularChapter(chapter *Chapter) error {
	// 普通章节处理
	chapter.HTMLContent = hgv.wrapWithClass(chapter.HTMLContent, "regular-content")
	return nil
}

func (hgv *HTMLGeneratorVisitor) VisitVolumeChapter(chapter *Chapter) error {
	// 卷章节特殊处理
	chapter.HTMLContent = hgv.wrapWithClass(chapter.HTMLContent, "volume-content")
	return nil
}

func (hgv *HTMLGeneratorVisitor) VisitSectionChapter(chapter *Chapter) error {
	// 小节特殊处理
	chapter.HTMLContent = hgv.wrapWithClass(chapter.HTMLContent, "section-content")
	return nil
}

// wrapWithClass 用 CSS 类包装内容
func (hgv *HTMLGeneratorVisitor) wrapWithClass(content, className string) string {
	return fmt.Sprintf(`<div class="%s">%s</div>`, className, content)
}

// StatisticsVisitor 统计访问者
type StatisticsVisitor struct {
	PrologueCount int
	EpilogueCount int
	RegularCount  int
	VolumeCount   int
	SectionCount  int
	TotalWords    int
}

// NewStatisticsVisitor 创建统计访问者
func NewStatisticsVisitor() *StatisticsVisitor {
	return &StatisticsVisitor{}
}

func (sv *StatisticsVisitor) VisitPrologue(chapter *Chapter) error {
	sv.PrologueCount++
	sv.TotalWords += chapter.WordCount
	return nil
}

func (sv *StatisticsVisitor) VisitEpilogue(chapter *Chapter) error {
	sv.EpilogueCount++
	sv.TotalWords += chapter.WordCount
	return nil
}

func (sv *StatisticsVisitor) VisitRegularChapter(chapter *Chapter) error {
	sv.RegularCount++
	sv.TotalWords += chapter.WordCount
	return nil
}

func (sv *StatisticsVisitor) VisitVolumeChapter(chapter *Chapter) error {
	sv.VolumeCount++
	sv.TotalWords += chapter.WordCount
	return nil
}

func (sv *StatisticsVisitor) VisitSectionChapter(chapter *Chapter) error {
	sv.SectionCount++
	sv.TotalWords += chapter.WordCount
	return nil
}

// GetReport 获取统计报告
func (sv *StatisticsVisitor) GetReport() string {
	var report strings.Builder
	
	report.WriteString("📊 章节统计报告\n")
	report.WriteString("==================\n")
	
	if sv.PrologueCount > 0 {
		report.WriteString(fmt.Sprintf("序言章节: %d\n", sv.PrologueCount))
	}
	if sv.VolumeCount > 0 {
		report.WriteString(fmt.Sprintf("卷章节: %d\n", sv.VolumeCount))
	}
	if sv.RegularCount > 0 {
		report.WriteString(fmt.Sprintf("普通章节: %d\n", sv.RegularCount))
	}
	if sv.SectionCount > 0 {
		report.WriteString(fmt.Sprintf("小节: %d\n", sv.SectionCount))
	}
	if sv.EpilogueCount > 0 {
		report.WriteString(fmt.Sprintf("后记章节: %d\n", sv.EpilogueCount))
	}
	
	total := sv.PrologueCount + sv.VolumeCount + sv.RegularCount + sv.SectionCount + sv.EpilogueCount
	report.WriteString(fmt.Sprintf("总章节数: %d\n", total))
	report.WriteString(fmt.Sprintf("总字数: %d\n", sv.TotalWords))
	
	return report.String()
}

// ValidationVisitor 验证访问者
type ValidationVisitor struct {
	errors []string
}

// NewValidationVisitor 创建验证访问者
func NewValidationVisitor() *ValidationVisitor {
	return &ValidationVisitor{
		errors: make([]string, 0),
	}
}

func (vv *ValidationVisitor) VisitPrologue(chapter *Chapter) error {
	return vv.validateChapter(chapter, "序言")
}

func (vv *ValidationVisitor) VisitEpilogue(chapter *Chapter) error {
	return vv.validateChapter(chapter, "后记")
}

func (vv *ValidationVisitor) VisitRegularChapter(chapter *Chapter) error {
	return vv.validateChapter(chapter, "章节")
}

func (vv *ValidationVisitor) VisitVolumeChapter(chapter *Chapter) error {
	return vv.validateChapter(chapter, "卷")
}

func (vv *ValidationVisitor) VisitSectionChapter(chapter *Chapter) error {
	return vv.validateChapter(chapter, "小节")
}

// validateChapter 验证章节
func (vv *ValidationVisitor) validateChapter(chapter *Chapter, chapterType string) error {
	if chapter.Title == "" {
		vv.errors = append(vv.errors, fmt.Sprintf("%s缺少标题: ID=%d", chapterType, chapter.ID))
	}
	
	if strings.TrimSpace(chapter.Content) == "" {
		vv.errors = append(vv.errors, fmt.Sprintf("%s内容为空: %s", chapterType, chapter.Title))
	}
	
	if chapter.WordCount < 10 {
		vv.errors = append(vv.errors, fmt.Sprintf("%s内容过短: %s (%d字)", chapterType, chapter.Title, chapter.WordCount))
	}
	
	return nil
}

// GetErrors 获取验证错误
func (vv *ValidationVisitor) GetErrors() []string {
	return vv.errors
}

// HasErrors 是否有错误
func (vv *ValidationVisitor) HasErrors() bool {
	return len(vv.errors) > 0
}

// ClearErrors 清空错误
func (vv *ValidationVisitor) ClearErrors() {
	vv.errors = make([]string, 0)
}

// ChapterProcessor 章节处理器
type ChapterProcessor struct {
	visitors []ChapterVisitor
}

// NewChapterProcessor 创建章节处理器
func NewChapterProcessor() *ChapterProcessor {
	return &ChapterProcessor{
		visitors: make([]ChapterVisitor, 0),
	}
}

// AddVisitor 添加访问者
func (cp *ChapterProcessor) AddVisitor(visitor ChapterVisitor) {
	cp.visitors = append(cp.visitors, visitor)
}

// ProcessChapter 处理章节
func (cp *ChapterProcessor) ProcessChapter(element ChapterElement) error {
	for _, visitor := range cp.visitors {
		if err := element.Accept(visitor); err != nil {
			return fmt.Errorf("访问者处理失败: %w", err)
		}
	}
	return nil
}

// ProcessNovel 处理整部小说
func (cp *ChapterProcessor) ProcessNovel(novel *Novel, chapterTypes []ChapterType) error {
	if len(novel.Chapters) != len(chapterTypes) {
		return fmt.Errorf("章节数量与类型数量不匹配")
	}
	
	for i, chapter := range novel.Chapters {
		chapterType := ChapterTypeChapter
		if i < len(chapterTypes) {
			chapterType = chapterTypes[i]
		}
		
		element := NewStandardChapter(chapter, chapterType)
		if err := cp.ProcessChapter(element); err != nil {
			return fmt.Errorf("处理章节 %s 失败: %w", chapter.Title, err)
		}
	}
	
	return nil
}
