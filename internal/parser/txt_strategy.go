package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

// TxtFileStrategy TXT 文件解析策略
type TxtFileStrategy struct {
	parser         *Parser
	txtFormat      *TxtFormat
	enhancedParser *EnhancedParser
	chapterAdapter *ChapterAdapter
	statefulParser *StatefulTxtParser
}

// NewTxtFileStrategy 创建 TXT 文件策略
func NewTxtFileStrategy(parser *Parser) *TxtFileStrategy {
	strategy := &TxtFileStrategy{
		parser:         parser,
		txtFormat:      NewTxtFormat(),
		chapterAdapter: NewChapterAdapter(),
		statefulParser: NewStatefulTxtParser(),
	}

	// 创建增强解析器
	consoleObserver := NewConsoleObserver(false) // 设置为 false 减少输出
	strategy.enhancedParser = NewEnhancedParser().
		WithLogging(consoleObserver).
		WithCaching().
		WithValidation()

	// 订阅解析事件
	strategy.chapterAdapter.Subscribe(consoleObserver)
	strategy.statefulParser.Subscribe(consoleObserver)

	return strategy
}

func (s *TxtFileStrategy) CanHandle(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".txt")
}

func (s *TxtFileStrategy) GetName() string {
	return "TxtFile"
}

func (s *TxtFileStrategy) Parse(novel *Novel, path string) error {
	// 使用文件名作为默认标题
	novel.Title = strings.TrimSuffix(filepath.Base(path), ".txt")

	// 读取文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")

	// 使用状态化解析器
	if err := s.statefulParser.ParseWithState(novel, lines); err != nil {
		return fmt.Errorf("状态化解析失败: %v", err)
	}

	// 使用适配器转换内容
	if err := s.chapterAdapter.ConvertNovel(novel, "txt"); err != nil {
		return fmt.Errorf("内容转换失败: %v", err)
	}

	return nil
}

// parseContentOld 旧的解析方法（保留作为备用）
func (s *TxtFileStrategy) parseContentOld(novel *Novel, lines []string) error {
	var txtChapters []*TxtChapter
	var currentChapter *TxtChapter
	var contentLines []string

	currentVolumeID := 0
	currentChapterID := 0
	currentSectionID := 0

	inMetadata := true // 文件开头可能有元数据
	metadataEndLine := 0

	// 第一遍：识别结构和元数据
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行
		if s.txtFormat.IsEmptyLine(line) {
			if inMetadata && len(contentLines) == 0 {
				continue
			}
			contentLines = append(contentLines, "")
			continue
		}

		// 检查分隔符
		if s.txtFormat.IsSeparator(line) {
			if inMetadata {
				metadataEndLine = i
				inMetadata = false
			}
			continue
		}

		// 在元数据阶段，尝试提取元数据
		if inMetadata && i < 50 { // 只在前50行查找元数据
			if key, value := s.txtFormat.ExtractMetadata(line); key != "" {
				s.setNovelMetadata(novel, key, value)
				continue
			}
		}

		// 识别章节类型
		chapterType, title := s.txtFormat.IdentifyChapterType(line)

		if chapterType != ChapterTypeUnknown {
			// 保存上一章节
			if currentChapter != nil {
				currentChapter.Content = s.txtFormat.CleanContent(strings.Join(contentLines, "\n"))
				currentChapter.LineEnd = i - 1
				txtChapters = append(txtChapters, currentChapter)
			}

			// 更新层级计数
			switch chapterType {
			case ChapterTypeVolume:
				currentVolumeID++
				currentChapterID = 0
				currentSectionID = 0
			case ChapterTypeChapter:
				currentChapterID++
				currentSectionID = 0
			case ChapterTypeSection:
				currentSectionID++
			}

			// 创建新章节
			currentChapter = &TxtChapter{
				Type:      chapterType,
				Level:     s.getChapterLevel(chapterType),
				VolumeID:  currentVolumeID,
				ChapterID: currentChapterID,
				SectionID: currentSectionID,
				Title:     title,
				LineStart: i,
			}

			contentLines = make([]string, 0)
			inMetadata = false
		} else {
			// 普通内容行
			if !inMetadata {
				contentLines = append(contentLines, line)
			}
		}
	}

	// 保存最后一章节
	if currentChapter != nil {
		currentChapter.Content = s.txtFormat.CleanContent(strings.Join(contentLines, "\n"))
		currentChapter.LineEnd = len(lines) - 1
		txtChapters = append(txtChapters, currentChapter)
	}

	// 如果没有找到任何章节，将整个文件作为一章
	if len(txtChapters) == 0 {
		allContent := strings.Join(lines[metadataEndLine:], "\n")
		txtChapters = append(txtChapters, &TxtChapter{
			Type:      ChapterTypeChapter,
			Level:     2,
			VolumeID:  1,
			ChapterID: 1,
			SectionID: 0,
			Title:     "正文",
			Content:   s.txtFormat.CleanContent(allContent),
			LineStart: metadataEndLine,
			LineEnd:   len(lines) - 1,
		})
	}

	// 转换为标准章节格式
	return s.convertToStandardChapters(novel, txtChapters)
}

// setNovelMetadata 设置小说元数据
func (s *TxtFileStrategy) setNovelMetadata(novel *Novel, key, value string) {
	switch key {
	case "title":
		novel.Title = value
	case "author":
		novel.Author = value
	case "description":
		novel.Description = value
	case "category":
		novel.Category = value
	case "tags":
		// 处理标签，支持逗号、分号、空格分隔
		tags := strings.FieldsFunc(value, func(r rune) bool {
			return r == ',' || r == ';' || r == '，' || r == '；'
		})
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		novel.Tags = tags
	}
}

// getChapterLevel 获取章节层级
func (s *TxtFileStrategy) getChapterLevel(chapterType ChapterType) int {
	switch chapterType {
	case ChapterTypeVolume:
		return 1
	case ChapterTypeChapter:
		return 2
	case ChapterTypeSection:
		return 3
	default:
		return 2 // 默认为章节级别
	}
}

// convertToStandardChapters 转换为标准章节格式
func (s *TxtFileStrategy) convertToStandardChapters(novel *Novel, txtChapters []*TxtChapter) error {
	chapterID := 0

	for _, txtChapter := range txtChapters {
		// 跳过空内容的章节
		if strings.TrimSpace(txtChapter.Content) == "" {
			continue
		}

		chapterID++

		// 生成章节标题
		title := s.generateChapterTitle(txtChapter)

		// 创建标准章节
		chapter := &Chapter{
			ID:          chapterID,
			Title:       title,
			Content:     txtChapter.Content,
			HTMLContent: string(blackfriday.Run([]byte(txtChapter.Content))),
			WordCount:   len([]rune(txtChapter.Content)),
			CreatedAt:   time.Now(),
			Path:        fmt.Sprintf("chapter-%d", chapterID),
		}

		novel.Chapters = append(novel.Chapters, chapter)
	}

	return nil
}

// generateChapterTitle 生成章节标题
func (s *TxtFileStrategy) generateChapterTitle(txtChapter *TxtChapter) string {
	switch txtChapter.Type {
	case ChapterTypeVolume:
		if txtChapter.Title != "" {
			return fmt.Sprintf("第%d卷 %s", txtChapter.VolumeID, txtChapter.Title)
		}
		return fmt.Sprintf("第%d卷", txtChapter.VolumeID)

	case ChapterTypeChapter:
		if txtChapter.VolumeID > 0 && txtChapter.Title != "" {
			return fmt.Sprintf("第%d章 %s", txtChapter.ChapterID, txtChapter.Title)
		} else if txtChapter.Title != "" {
			return fmt.Sprintf("第%d章 %s", txtChapter.ChapterID, txtChapter.Title)
		}
		return fmt.Sprintf("第%d章", txtChapter.ChapterID)

	case ChapterTypeSection:
		if txtChapter.Title != "" {
			return fmt.Sprintf("第%d节 %s", txtChapter.SectionID, txtChapter.Title)
		}
		return fmt.Sprintf("第%d节", txtChapter.SectionID)

	case ChapterTypePrologue:
		return txtChapter.Title

	case ChapterTypeEpilogue:
		return txtChapter.Title

	case ChapterTypeSummary:
		return txtChapter.Title

	default:
		if txtChapter.Title != "" {
			return txtChapter.Title
		}
		return fmt.Sprintf("第%d章", txtChapter.ChapterID)
	}
}

// TxtDirectoryStrategy TXT 目录解析策略
type TxtDirectoryStrategy struct {
	parser    *Parser
	txtFormat *TxtFormat
}

// NewTxtDirectoryStrategy 创建 TXT 目录策略
func NewTxtDirectoryStrategy(parser *Parser) *TxtDirectoryStrategy {
	return &TxtDirectoryStrategy{
		parser:    parser,
		txtFormat: NewTxtFormat(),
	}
}

func (s *TxtDirectoryStrategy) CanHandle(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	// 检查是否包含 .txt 文件
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	hasTxt := false
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".txt") {
			hasTxt = true
			break
		}
	}

	return hasTxt
}

func (s *TxtDirectoryStrategy) GetName() string {
	return "TxtDirectory"
}

func (s *TxtDirectoryStrategy) Parse(novel *Novel, path string) error {
	// 查找元数据文件
	metaFiles := []string{"meta.txt", "info.txt", "简介.txt"}
	for _, metaFile := range metaFiles {
		metaPath := filepath.Join(path, metaFile)
		if _, err := os.Stat(metaPath); err == nil {
			if err := s.parseMetadataFile(novel, metaPath); err != nil {
				fmt.Printf("警告：解析元数据文件失败: %v\n", err)
			}
			break
		}
	}

	// 如果没有找到元数据，使用目录名作为标题
	if novel.Title == "" {
		novel.Title = filepath.Base(path)
	}

	// 读取所有 TXT 文件
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	var txtFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".txt") {
			// 跳过元数据文件
			skip := false
			for _, metaFile := range metaFiles {
				if entry.Name() == metaFile {
					skip = true
					break
				}
			}
			if !skip {
				txtFiles = append(txtFiles, filepath.Join(path, entry.Name()))
			}
		}
	}

	// 按文件名排序
	s.sortTxtFiles(txtFiles)

	// 解析每个文件
	chapterID := 0
	for _, txtFile := range txtFiles {
		chapters, err := s.parseTextFile(txtFile, &chapterID)
		if err != nil {
			fmt.Printf("警告：解析文件 %s 失败: %v\n", txtFile, err)
			continue
		}
		novel.Chapters = append(novel.Chapters, chapters...)
	}

	return nil
}

// parseMetadataFile 解析元数据文件
func (s *TxtDirectoryStrategy) parseMetadataFile(novel *Novel, metaPath string) error {
	file, err := os.Open(metaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var descriptionLines []string
	inDescription := false

	for scanner.Scan() {
		line := scanner.Text()

		if key, value := s.txtFormat.ExtractMetadata(line); key != "" {
			switch key {
			case "title":
				novel.Title = value
			case "author":
				novel.Author = value
			case "description":
				if value != "" {
					novel.Description = value
				} else {
					inDescription = true
				}
			case "category":
				novel.Category = value
			case "tags":
				// 处理标签，支持逗号、分号、空格分隔
				tags := strings.FieldsFunc(value, func(r rune) bool {
					return r == ',' || r == ';' || r == '，' || r == '；'
				})
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
				novel.Tags = tags
			}
		} else if inDescription && strings.TrimSpace(line) != "" {
			descriptionLines = append(descriptionLines, strings.TrimSpace(line))
		}
	}

	// 如果简介是多行的
	if len(descriptionLines) > 0 && novel.Description == "" {
		novel.Description = strings.Join(descriptionLines, "\n")
	}

	return scanner.Err()
}

// sortTxtFiles 排序 TXT 文件
func (s *TxtDirectoryStrategy) sortTxtFiles(files []string) {
	// 简单的数字排序
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if s.compareFilenames(files[i], files[j]) > 0 {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
}

// compareFilenames 比较文件名
func (s *TxtDirectoryStrategy) compareFilenames(a, b string) int {
	nameA := filepath.Base(a)
	nameB := filepath.Base(b)

	// 提取数字
	numA := s.extractNumberFromFilename(nameA)
	numB := s.extractNumberFromFilename(nameB)

	if numA != numB {
		return numA - numB
	}

	// 如果数字相同，按字符串排序
	if nameA < nameB {
		return -1
	} else if nameA > nameB {
		return 1
	}
	return 0
}

// extractNumberFromFilename 从文件名提取数字
func (s *TxtDirectoryStrategy) extractNumberFromFilename(filename string) int {
	// 移除扩展名
	name := strings.TrimSuffix(filename, ".txt")

	// 尝试提取开头的数字
	for i, char := range name {
		if char >= '0' && char <= '9' {
			// 找到数字的结束位置
			j := i
			for j < len(name) && name[j] >= '0' && name[j] <= '9' {
				j++
			}
			if num, err := strconv.Atoi(name[i:j]); err == nil {
				return num
			}
		}
	}

	return 0
}

// parseTextFile 解析单个文本文件
func (s *TxtDirectoryStrategy) parseTextFile(filePath string, chapterID *int) ([]*Chapter, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// 使用 TXT 文件策略解析
	txtStrategy := NewTxtFileStrategy(s.parser)

	// 创建临时小说对象
	tempNovel := &Novel{
		Path:     filePath,
		Chapters: make([]*Chapter, 0),
	}

	if err := txtStrategy.parseContentOld(tempNovel, lines); err != nil {
		return nil, err
	}

	// 重新分配章节ID
	var chapters []*Chapter
	for _, chapter := range tempNovel.Chapters {
		*chapterID++
		chapter.ID = *chapterID
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}
