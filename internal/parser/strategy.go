package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseStrategy 解析策略接口
type ParseStrategy interface {
	Parse(novel *Novel, path string) error
	CanHandle(path string) bool
	GetName() string
}

// SingleFileStrategy 单文件解析策略
type SingleFileStrategy struct {
	parser *Parser
}

// NewSingleFileStrategy 创建单文件策略
func NewSingleFileStrategy(parser *Parser) *SingleFileStrategy {
	return &SingleFileStrategy{parser: parser}
}

func (s *SingleFileStrategy) CanHandle(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md")
}

func (s *SingleFileStrategy) GetName() string {
	return "SingleFile"
}

func (s *SingleFileStrategy) Parse(novel *Novel, path string) error {
	_, err := s.parser.parseNovelFromFile(novel)
	return err
}

// MultiFileStrategy 多文件解析策略
type MultiFileStrategy struct {
	parser *Parser
}

// NewMultiFileStrategy 创建多文件策略
func NewMultiFileStrategy(parser *Parser) *MultiFileStrategy {
	return &MultiFileStrategy{parser: parser}
}

func (s *MultiFileStrategy) CanHandle(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	// 检查是否包含 .md 文件但不包含子目录结构
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	hasMarkdown := false
	hasSubDirs := false

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			hasSubDirs = true
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".md") && entry.Name() != "meta.md" {
			hasMarkdown = true
		}
	}

	return hasMarkdown && !hasSubDirs
}

func (s *MultiFileStrategy) GetName() string {
	return "MultiFile"
}

func (s *MultiFileStrategy) Parse(novel *Novel, path string) error {
	_, err := s.parser.parseNovelFromDir(novel)
	return err
}

// MultiVolumeStrategy 多卷解析策略
type MultiVolumeStrategy struct {
	parser *Parser
}

// NewMultiVolumeStrategy 创建多卷策略
func NewMultiVolumeStrategy(parser *Parser) *MultiVolumeStrategy {
	return &MultiVolumeStrategy{parser: parser}
}

func (s *MultiVolumeStrategy) CanHandle(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	// 检查是否包含卷目录
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	volumeDirs := 0
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			// 检查目录名是否包含"卷"字符
			if strings.Contains(entry.Name(), "卷") ||
				strings.Contains(strings.ToLower(entry.Name()), "volume") {
				volumeDirs++
			}
		}
	}

	return volumeDirs > 0
}

func (s *MultiVolumeStrategy) GetName() string {
	return "MultiVolume"
}

func (s *MultiVolumeStrategy) Parse(novel *Novel, path string) error {
	// 解析多卷结构
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	// 首先读取元数据
	metaFile := filepath.Join(path, "meta.md")
	if _, err := os.Stat(metaFile); err == nil {
		if err := s.parser.parseNovelMeta(novel, metaFile); err != nil {
			return fmt.Errorf("解析元数据失败: %v", err)
		}
	}

	// 如果没有找到元数据，使用目录名作为标题
	if novel.Title == "" {
		novel.Title = filepath.Base(path)
	}

	// 解析各卷
	volumeDirs := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") && entry.Name() != "meta.md" {
			volumeDirs = append(volumeDirs, entry.Name())
		}
	}

	// 按名称排序卷
	for _, volumeDir := range volumeDirs {
		volumePath := filepath.Join(path, volumeDir)
		if err := s.parseVolume(novel, volumePath); err != nil {
			return fmt.Errorf("解析卷 %s 失败: %v", volumeDir, err)
		}
	}

	// 重新分配章节ID
	for i, chapter := range novel.Chapters {
		chapter.ID = i + 1
	}

	return nil
}

// parseVolume 解析单个卷
func (s *MultiVolumeStrategy) parseVolume(novel *Novel, volumePath string) error {
	entries, err := os.ReadDir(volumePath)
	if err != nil {
		return fmt.Errorf("读取卷目录失败: %v", err)
	}

	// 收集章节文件
	chapterFiles := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			chapterFiles = append(chapterFiles, filepath.Join(volumePath, entry.Name()))
		}
	}

	// 解析每个章节
	for _, chapterFile := range chapterFiles {
		chapter, err := s.parser.parseChapterFile(chapterFile)
		if err != nil {
			return fmt.Errorf("解析章节文件 %s 失败: %v", chapterFile, err)
		}
		novel.Chapters = append(novel.Chapters, chapter)
	}

	return nil
}

// StrategyManager 策略管理器
type StrategyManager struct {
	strategies []ParseStrategy
}

// NewStrategyManager 创建策略管理器
func NewStrategyManager(parser *Parser) *StrategyManager {
	return &StrategyManager{
		strategies: []ParseStrategy{
			NewTxtDirectoryStrategy(parser), // 优先检查 TXT 目录
			NewTxtFileStrategy(parser),      // 然后检查 TXT 文件
			NewMultiVolumeStrategy(parser),  // 接着检查多卷 Markdown
			NewMultiFileStrategy(parser),    // 然后检查多文件 Markdown
			NewSingleFileStrategy(parser),   // 最后检查单文件 Markdown
		},
	}
}

// SelectStrategy 选择合适的解析策略
func (sm *StrategyManager) SelectStrategy(path string) ParseStrategy {
	for _, strategy := range sm.strategies {
		if strategy.CanHandle(path) {
			return strategy
		}
	}
	// 默认返回单文件策略
	return sm.strategies[len(sm.strategies)-1]
}

// GetAvailableStrategies 获取所有可用策略
func (sm *StrategyManager) GetAvailableStrategies() []string {
	names := make([]string, len(sm.strategies))
	for i, strategy := range sm.strategies {
		names[i] = strategy.GetName()
	}
	return names
}
