package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

// Novel 小说结构
type Novel struct {
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Description string     `json:"description"`
	Cover       string     `json:"cover"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Chapters    []*Chapter `json:"chapters"`
	Path        string     `json:"path"`
}

// Chapter 章节结构
type Chapter struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	HTMLContent string    `json:"html_content"`
	WordCount   int       `json:"word_count"`
	CreatedAt   time.Time `json:"created_at"`
	Path        string    `json:"path"`
}

// Parser Markdown解析器
type Parser struct {
	chapterRegex *regexp.Regexp
	metaRegex    *regexp.Regexp
}

// New 创建新的解析器
func New() *Parser {
	return &Parser{
		// 匹配章节标题，支持多种格式，包括卷和章节
		chapterRegex: regexp.MustCompile(`^#+\s*(?:第[0-9一二三四五六七八九十百千万]+[卷章回]|Chapter\s*\d+|Volume\s*\d+|[0-9]+\.)\s*(.+)`),
		// 匹配元数据
		metaRegex: regexp.MustCompile(`^---\s*$`),
	}
}

// ParseNovel 解析小说目录
func (p *Parser) ParseNovel(novelPath string) (*Novel, error) {
	info, err := os.Stat(novelPath)
	if err != nil {
		return nil, fmt.Errorf("无法访问路径 %s: %v", novelPath, err)
	}

	novel := &Novel{
		Path:      novelPath,
		CreatedAt: info.ModTime(),
		UpdatedAt: info.ModTime(),
		Chapters:  make([]*Chapter, 0),
	}

	if info.IsDir() {
		// 目录模式：每个文件是一个章节
		return p.parseNovelFromDir(novel)
	} else {
		// 单文件模式：一个文件包含所有章节
		return p.parseNovelFromFile(novel)
	}
}

// parseNovelFromDir 从目录解析小说
func (p *Parser) parseNovelFromDir(novel *Novel) (*Novel, error) {
	// 读取目录中的所有 .md 文件
	files, err := filepath.Glob(filepath.Join(novel.Path, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	// 查找元数据文件
	metaFile := filepath.Join(novel.Path, "meta.md")
	if _, err := os.Stat(metaFile); err == nil {
		if err := p.parseNovelMeta(novel, metaFile); err != nil {
			return nil, fmt.Errorf("解析元数据失败: %v", err)
		}
		// 从文件列表中移除元数据文件
		for i, file := range files {
			if file == metaFile {
				files = append(files[:i], files[i+1:]...)
				break
			}
		}
	}

	// 如果没有找到元数据，使用目录名作为标题
	if novel.Title == "" {
		novel.Title = filepath.Base(novel.Path)
	}

	// 解析每个章节文件
	for _, file := range files {
		chapter, err := p.parseChapterFile(file)
		if err != nil {
			return nil, fmt.Errorf("解析章节文件 %s 失败: %v", file, err)
		}
		novel.Chapters = append(novel.Chapters, chapter)
	}

	// 按文件名排序章节
	sort.Slice(novel.Chapters, func(i, j int) bool {
		return novel.Chapters[i].Path < novel.Chapters[j].Path
	})

	// 重新分配章节ID
	for i, chapter := range novel.Chapters {
		chapter.ID = i + 1
	}

	return novel, nil
}

// parseNovelFromFile 从单个文件解析小说
func (p *Parser) parseNovelFromFile(novel *Novel) (*Novel, error) {
	content, err := ioutil.ReadFile(novel.Path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var currentChapter *Chapter
	var contentLines []string
	inMeta := false
	chapterID := 0

	// 使用文件名作为默认标题
	novel.Title = strings.TrimSuffix(filepath.Base(novel.Path), ".md")

	for i, line := range lines {
		// 检查是否是元数据分隔符
		if p.metaRegex.MatchString(line) {
			if i == 0 {
				inMeta = true
				continue
			} else if inMeta {
				inMeta = false
				continue
			}
		}

		// 处理元数据
		if inMeta {
			p.parseMetaLine(novel, line)
			continue
		}

		// 检查是否是章节标题
		if matches := p.chapterRegex.FindStringSubmatch(line); matches != nil {
			// 保存上一章节
			if currentChapter != nil {
				currentChapter.Content = strings.Join(contentLines, "\n")
				currentChapter.HTMLContent = string(blackfriday.Run([]byte(currentChapter.Content)))
				currentChapter.WordCount = len([]rune(currentChapter.Content))
				novel.Chapters = append(novel.Chapters, currentChapter)
			}

			// 开始新章节
			chapterID++
			currentChapter = &Chapter{
				ID:    chapterID,
				Title: strings.TrimSpace(matches[1]),
				Path:  fmt.Sprintf("chapter-%d", chapterID),
			}
			contentLines = make([]string, 0)
		} else if currentChapter != nil {
			// 添加内容到当前章节
			contentLines = append(contentLines, line)
		}
	}

	// 保存最后一个章节
	if currentChapter != nil {
		currentChapter.Content = strings.Join(contentLines, "\n")
		currentChapter.HTMLContent = string(blackfriday.Run([]byte(currentChapter.Content)))
		currentChapter.WordCount = len([]rune(currentChapter.Content))
		novel.Chapters = append(novel.Chapters, currentChapter)
	}

	return novel, nil
}

// parseNovelMeta 解析小说元数据
func (p *Parser) parseNovelMeta(novel *Novel, metaFile string) error {
	file, err := os.Open(metaFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inMeta := false

	for scanner.Scan() {
		line := scanner.Text()

		if p.metaRegex.MatchString(line) {
			if !inMeta {
				inMeta = true
				continue
			} else {
				break
			}
		}

		if inMeta {
			p.parseMetaLine(novel, line)
		}
	}

	return scanner.Err()
}

// parseMetaLine 解析元数据行
func (p *Parser) parseMetaLine(novel *Novel, line string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch strings.ToLower(key) {
	case "title", "标题":
		novel.Title = value
	case "author", "作者":
		novel.Author = value
	case "description", "简介", "描述":
		novel.Description = value
	case "cover", "封面":
		novel.Cover = value
	}
}

// parseChapterFile 解析章节文件
func (p *Parser) parseChapterFile(filePath string) (*Chapter, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// 提取章节编号和标题
	fileName := strings.TrimSuffix(filepath.Base(filePath), ".md")
	chapterID := 0
	title := fileName

	// 尝试从文件名提取章节编号
	if matches := regexp.MustCompile(`^(\d+)[-_.]?(.*)$`).FindStringSubmatch(fileName); matches != nil {
		if id, err := strconv.Atoi(matches[1]); err == nil {
			chapterID = id
			if matches[2] != "" {
				title = matches[2]
			}
		}
	}

	// 解析内容
	lines := strings.Split(string(content), "\n")
	var contentLines []string
	inMeta := false

	for i, line := range lines {
		// 检查元数据分隔符
		if p.metaRegex.MatchString(line) {
			if i == 0 {
				inMeta = true
				continue
			} else if inMeta {
				inMeta = false
				continue
			}
		}

		if inMeta {
			// 处理章节元数据
			if strings.HasPrefix(strings.ToLower(line), "title:") || strings.HasPrefix(strings.ToLower(line), "标题:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					title = strings.TrimSpace(parts[1])
				}
			}
		} else {
			// 检查是否第一行是标题
			if i == 0 || (i == 1 && inMeta) {
				if matches := p.chapterRegex.FindStringSubmatch(line); matches != nil {
					title = strings.TrimSpace(matches[1])
					continue
				}
			}
			contentLines = append(contentLines, line)
		}
	}

	contentText := strings.Join(contentLines, "\n")
	
	chapter := &Chapter{
		ID:          chapterID,
		Title:       title,
		Content:     contentText,
		HTMLContent: string(blackfriday.Run([]byte(contentText))),
		WordCount:   len([]rune(contentText)),
		CreatedAt:   info.ModTime(),
		Path:        strings.TrimSuffix(fileName, ".md"),
	}

	return chapter, nil
}
