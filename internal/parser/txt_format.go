package parser

import (
	"regexp"
	"strings"
)

// TxtFormat TXT 文件格式定义和解析规则
type TxtFormat struct {
	// 章节标题匹配规则
	VolumeRegex    *regexp.Regexp // 卷标题
	ChapterRegex   *regexp.Regexp // 章节标题
	SectionRegex   *regexp.Regexp // 小节标题
	PrologueRegex  *regexp.Regexp // 序言类
	EpilogueRegex  *regexp.Regexp // 后记类
	SummaryRegex   *regexp.Regexp // 内容概述
	
	// 分隔符规则
	SeparatorRegex *regexp.Regexp // 章节分隔符
	
	// 元数据规则
	TitleRegex     *regexp.Regexp // 小说标题
	AuthorRegex    *regexp.Regexp // 作者
	IntroRegex     *regexp.Regexp // 简介
}

// NewTxtFormat 创建 TXT 格式解析器
func NewTxtFormat() *TxtFormat {
	return &TxtFormat{
		// 卷标题：第一卷、第1卷、卷一、Volume 1 等
		VolumeRegex: regexp.MustCompile(`(?i)^\s*(?:第[0-9一二三四五六七八九十百千万]+卷|卷[0-9一二三四五六七八九十百千万]+|Volume\s*[0-9]+|VOLUME\s*[0-9]+)\s*[：:\s]*(.*)$`),
		
		// 章节标题：第一章、第1章、章节001、Chapter 1 等
		ChapterRegex: regexp.MustCompile(`(?i)^\s*(?:第[0-9一二三四五六七八九十百千万]+[章回]|[章回][0-9一二三四五六七八九十百千万]+|Chapter\s*[0-9]+|[0-9]{1,4}[\.、\s]|第[0-9]{1,4}[章回]|[0-9]{1,4}章)\s*[：:\s]*(.*)$`),
		
		// 小节标题：第一节、1.1、一、（一）等
		SectionRegex: regexp.MustCompile(`(?i)^\s*(?:第[0-9一二三四五六七八九十]+节|[0-9]+\.[0-9]+|[一二三四五六七八九十]+、|\([一二三四五六七八九十0-9]+\))\s*[：:\s]*(.*)$`),
		
		// 序言类：序、序言、楔子、引子、前言等
		PrologueRegex: regexp.MustCompile(`(?i)^\s*(?:序言?|楔子|引子|前言|开篇|缘起|题记|自序|代序)\s*[：:\s]*(.*)$`),
		
		// 后记类：后记、尾声、结语、跋等
		EpilogueRegex: regexp.MustCompile(`(?i)^\s*(?:后记|尾声|结语|跋|跋文|结尾|终章|大结局|完结感言|作者后记)\s*[：:\s]*(.*)$`),
		
		// 内容概述：简介、内容简介、故事梗概等
		SummaryRegex: regexp.MustCompile(`(?i)^\s*(?:简介|内容简介|故事梗概|内容梗概|作品简介|小说简介)\s*[：:\s]*(.*)$`),
		
		// 分隔符：多个等号、减号、星号等
		SeparatorRegex: regexp.MustCompile(`^\s*[=\-*~]{3,}\s*$`),
		
		// 小说标题：书名、标题等
		TitleRegex: regexp.MustCompile(`(?i)^\s*(?:书名|标题|小说名|作品名)\s*[：:\s]+(.+)$`),
		
		// 作者：作者、Author等
		AuthorRegex: regexp.MustCompile(`(?i)^\s*(?:作者|Author|著|编著|原著)\s*[：:\s]+(.+)$`),
		
		// 简介：可能跨多行
		IntroRegex: regexp.MustCompile(`(?i)^\s*(?:简介|内容简介|故事简介|作品简介)\s*[：:\s]*(.*)$`),
	}
}

// ChapterType 章节类型
type ChapterType int

const (
	ChapterTypeUnknown ChapterType = iota
	ChapterTypeVolume              // 卷
	ChapterTypeChapter             // 章
	ChapterTypeSection             // 节
	ChapterTypePrologue            // 序言
	ChapterTypeEpilogue            // 后记
	ChapterTypeSummary             // 概述
)

// String 返回章节类型的字符串表示
func (ct ChapterType) String() string {
	switch ct {
	case ChapterTypeVolume:
		return "卷"
	case ChapterTypeChapter:
		return "章"
	case ChapterTypeSection:
		return "节"
	case ChapterTypePrologue:
		return "序言"
	case ChapterTypeEpilogue:
		return "后记"
	case ChapterTypeSummary:
		return "概述"
	default:
		return "未知"
	}
}

// TxtChapter TXT 章节结构
type TxtChapter struct {
	Type       ChapterType `json:"type"`
	Level      int         `json:"level"`      // 层级：1=卷，2=章，3=节
	VolumeID   int         `json:"volume_id"`  // 所属卷ID
	ChapterID  int         `json:"chapter_id"` // 章节ID
	SectionID  int         `json:"section_id"` // 小节ID
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	LineStart  int         `json:"line_start"` // 起始行号
	LineEnd    int         `json:"line_end"`   // 结束行号
}

// IdentifyChapterType 识别章节类型
func (tf *TxtFormat) IdentifyChapterType(line string) (ChapterType, string) {
	line = strings.TrimSpace(line)
	
	// 检查序言类
	if matches := tf.PrologueRegex.FindStringSubmatch(line); matches != nil {
		title := strings.TrimSpace(matches[1])
		if title == "" {
			title = strings.TrimSpace(matches[0])
		}
		return ChapterTypePrologue, title
	}
	
	// 检查后记类
	if matches := tf.EpilogueRegex.FindStringSubmatch(line); matches != nil {
		title := strings.TrimSpace(matches[1])
		if title == "" {
			title = strings.TrimSpace(matches[0])
		}
		return ChapterTypeEpilogue, title
	}
	
	// 检查概述类
	if matches := tf.SummaryRegex.FindStringSubmatch(line); matches != nil {
		title := strings.TrimSpace(matches[1])
		if title == "" {
			title = strings.TrimSpace(matches[0])
		}
		return ChapterTypeSummary, title
	}
	
	// 检查卷标题
	if matches := tf.VolumeRegex.FindStringSubmatch(line); matches != nil {
		title := strings.TrimSpace(matches[1])
		if title == "" {
			// 如果没有卷名，使用整个匹配作为标题
			title = strings.TrimSpace(strings.Split(matches[0], ":")[0])
		}
		return ChapterTypeVolume, title
	}
	
	// 检查章节标题
	if matches := tf.ChapterRegex.FindStringSubmatch(line); matches != nil {
		title := strings.TrimSpace(matches[1])
		if title == "" {
			// 如果没有章节名，使用整个匹配作为标题
			title = strings.TrimSpace(strings.Split(matches[0], ":")[0])
		}
		return ChapterTypeChapter, title
	}
	
	// 检查小节标题
	if matches := tf.SectionRegex.FindStringSubmatch(line); matches != nil {
		title := strings.TrimSpace(matches[1])
		if title == "" {
			title = strings.TrimSpace(strings.Split(matches[0], ":")[0])
		}
		return ChapterTypeSection, title
	}
	
	return ChapterTypeUnknown, ""
}

// IsSeparator 判断是否为分隔符
func (tf *TxtFormat) IsSeparator(line string) bool {
	return tf.SeparatorRegex.MatchString(line)
}

// ExtractMetadata 提取元数据
func (tf *TxtFormat) ExtractMetadata(line string) (key, value string) {
	line = strings.TrimSpace(line)
	
	// 检查标题
	if matches := tf.TitleRegex.FindStringSubmatch(line); matches != nil {
		return "title", strings.TrimSpace(matches[1])
	}
	
	// 检查作者
	if matches := tf.AuthorRegex.FindStringSubmatch(line); matches != nil {
		return "author", strings.TrimSpace(matches[1])
	}
	
	// 检查简介
	if matches := tf.IntroRegex.FindStringSubmatch(line); matches != nil {
		return "description", strings.TrimSpace(matches[1])
	}
	
	return "", ""
}

// IsEmptyLine 判断是否为空行
func (tf *TxtFormat) IsEmptyLine(line string) bool {
	return strings.TrimSpace(line) == ""
}

// CleanContent 清理内容
func (tf *TxtFormat) CleanContent(content string) string {
	// 移除多余的空行
	lines := strings.Split(content, "\n")
	var cleanLines []string
	
	emptyCount := 0
	for _, line := range lines {
		if tf.IsEmptyLine(line) {
			emptyCount++
			if emptyCount <= 2 { // 最多保留2个连续空行
				cleanLines = append(cleanLines, "")
			}
		} else {
			emptyCount = 0
			cleanLines = append(cleanLines, strings.TrimSpace(line))
		}
	}
	
	// 移除开头和结尾的空行
	for len(cleanLines) > 0 && cleanLines[0] == "" {
		cleanLines = cleanLines[1:]
	}
	for len(cleanLines) > 0 && cleanLines[len(cleanLines)-1] == "" {
		cleanLines = cleanLines[:len(cleanLines)-1]
	}
	
	return strings.Join(cleanLines, "\n")
}

// ExtractNumber 从标题中提取数字
func (tf *TxtFormat) ExtractNumber(title string) int {
	// 数字提取正则
	numberRegex := regexp.MustCompile(`[0-9]+`)
	
	if matches := numberRegex.FindString(title); matches != "" {
		// 尝试转换为数字
		if num := parseChineseNumber(matches); num > 0 {
			return num
		}
	}
	
	// 中文数字提取
	chineseRegex := regexp.MustCompile(`[一二三四五六七八九十百千万]+`)
	if matches := chineseRegex.FindString(title); matches != "" {
		return parseChineseNumber(matches)
	}
	
	return 0
}

// parseChineseNumber 解析中文数字
func parseChineseNumber(chinese string) int {
	chineseNums := map[rune]int{
		'一': 1, '二': 2, '三': 3, '四': 4, '五': 5,
		'六': 6, '七': 7, '八': 8, '九': 9, '十': 10,
		'百': 100, '千': 1000, '万': 10000,
	}
	
	result := 0
	temp := 0
	
	for _, char := range chinese {
		if val, exists := chineseNums[char]; exists {
			switch char {
			case '十':
				if temp == 0 {
					temp = 1
				}
				result += temp * 10
				temp = 0
			case '百':
				result += temp * 100
				temp = 0
			case '千':
				result += temp * 1000
				temp = 0
			case '万':
				result += temp * 10000
				temp = 0
			default:
				temp = val
			}
		}
	}
	
	result += temp
	return result
}
