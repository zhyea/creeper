package deploy

import (
	"os"
	"path/filepath"
	"sort"

	"creeper/internal/common"
)

// FileInfo 文件信息
type FileInfo struct {
	Path     string
	Name     string
	Size     int64
	IsDir    bool
	ModTime  int64
	Relative string
}

// FileIterator 文件迭代器接口
type FileIterator interface {
	HasNext() bool
	Next() *FileInfo
	Reset()
	GetTotal() int
	GetCurrent() int
}

// DirectoryIterator 目录迭代器
type DirectoryIterator struct {
	files       []*FileInfo
	current     int
	total       int
	rootDir     string
	logger      *common.Logger
}

// NewDirectoryIterator 创建目录迭代器
func NewDirectoryIterator(rootDir string) *DirectoryIterator {
	return &DirectoryIterator{
		rootDir: rootDir,
		logger:  common.GetLogger(),
	}
}

// LoadFiles 加载文件列表
func (di *DirectoryIterator) LoadFiles() error {
	di.logger.Info("加载文件列表:", di.rootDir)
	
	var files []*FileInfo
	
	err := filepath.Walk(di.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 计算相对路径
		relPath, err := filepath.Rel(di.rootDir, path)
		if err != nil {
			return err
		}
		
		// 跳过根目录
		if relPath == "." {
			return nil
		}
		
		fileInfo := &FileInfo{
			Path:     path,
			Name:     info.Name(),
			Size:     info.Size(),
			IsDir:    info.IsDir(),
			ModTime:  info.ModTime().Unix(),
			Relative: relPath,
		}
		
		files = append(files, fileInfo)
		return nil
	})
	
	if err != nil {
		return err
	}
	
	// 按路径排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Relative < files[j].Relative
	})
	
	di.files = files
	di.total = len(files)
	di.current = 0
	
	di.logger.Info("文件列表加载完成，共", di.total, "个文件")
	return nil
}

// HasNext 是否有下一个
func (di *DirectoryIterator) HasNext() bool {
	return di.current < di.total
}

// Next 获取下一个
func (di *DirectoryIterator) Next() *FileInfo {
	if !di.HasNext() {
		return nil
	}
	
	file := di.files[di.current]
	di.current++
	return file
}

// Reset 重置迭代器
func (di *DirectoryIterator) Reset() {
	di.current = 0
}

// GetTotal 获取总数
func (di *DirectoryIterator) GetTotal() int {
	return di.total
}

// GetCurrent 获取当前位置
func (di *DirectoryIterator) GetCurrent() int {
	return di.current
}

// GetFilesByExtension 按扩展名获取文件
func (di *DirectoryIterator) GetFilesByExtension(ext string) []*FileInfo {
	var result []*FileInfo
	
	for _, file := range di.files {
		if filepath.Ext(file.Name) == ext {
			result = append(result, file)
		}
	}
	
	return result
}

// GetFilesBySize 按大小范围获取文件
func (di *DirectoryIterator) GetFilesBySize(minSize, maxSize int64) []*FileInfo {
	var result []*FileInfo
	
	for _, file := range di.files {
		if file.Size >= minSize && file.Size <= maxSize {
			result = append(result, file)
		}
	}
	
	return result
}

// BatchIterator 批量迭代器
type BatchIterator struct {
	iterator    FileIterator
	batchSize   int
	currentBatch []*FileInfo
	logger      *common.Logger
}

// NewBatchIterator 创建批量迭代器
func NewBatchIterator(iterator FileIterator, batchSize int) *BatchIterator {
	return &BatchIterator{
		iterator:  iterator,
		batchSize: batchSize,
		logger:    common.GetLogger(),
	}
}

// HasNextBatch 是否有下一批
func (bi *BatchIterator) HasNextBatch() bool {
	return bi.iterator.HasNext()
}

// NextBatch 获取下一批
func (bi *BatchIterator) NextBatch() []*FileInfo {
	var batch []*FileInfo
	
	for i := 0; i < bi.batchSize && bi.iterator.HasNext(); i++ {
		file := bi.iterator.Next()
		if file != nil {
			batch = append(batch, file)
		}
	}
	
	bi.currentBatch = batch
	return batch
}

// GetCurrentBatch 获取当前批次
func (bi *BatchIterator) GetCurrentBatch() []*FileInfo {
	return bi.currentBatch
}

// Reset 重置迭代器
func (bi *BatchIterator) Reset() {
	bi.iterator.Reset()
	bi.currentBatch = nil
}

// GetProgress 获取进度
func (bi *BatchIterator) GetProgress() float64 {
	if bi.iterator.GetTotal() == 0 {
		return 0
	}
	return float64(bi.iterator.GetCurrent()) / float64(bi.iterator.GetTotal()) * 100
}

// FilterIterator 过滤迭代器
type FilterIterator struct {
	iterator FileIterator
	filter   func(*FileInfo) bool
	logger   *common.Logger
}

// NewFilterIterator 创建过滤迭代器
func NewFilterIterator(iterator FileIterator, filter func(*FileInfo) bool) *FilterIterator {
	return &FilterIterator{
		iterator: iterator,
		filter:   filter,
		logger:   common.GetLogger(),
	}
}

// HasNext 是否有下一个
func (fi *FilterIterator) HasNext() bool {
	// 预取下一个符合条件的文件
	for fi.iterator.HasNext() {
		file := fi.iterator.Next()
		if fi.filter(file) {
			// 将文件放回迭代器（这里简化处理）
			return true
		}
	}
	return false
}

// Next 获取下一个
func (fi *FilterIterator) Next() *FileInfo {
	for fi.iterator.HasNext() {
		file := fi.iterator.Next()
		if fi.filter(file) {
			return file
		}
	}
	return nil
}

// Reset 重置迭代器
func (fi *FilterIterator) Reset() {
	fi.iterator.Reset()
}

// GetTotal 获取总数
func (fi *FilterIterator) GetTotal() int {
	// 计算符合条件的文件总数
	count := 0
	fi.iterator.Reset()
	
	for fi.iterator.HasNext() {
		file := fi.iterator.Next()
		if fi.filter(file) {
			count++
		}
	}
	
	fi.iterator.Reset()
	return count
}

// GetCurrent 获取当前位置
func (fi *FilterIterator) GetCurrent() int {
	return fi.iterator.GetCurrent()
}

// FileCollection 文件集合
type FileCollection struct {
	iterator FileIterator
	logger   *common.Logger
}

// NewFileCollection 创建文件集合
func NewFileCollection(iterator FileIterator) *FileCollection {
	return &FileCollection{
		iterator: iterator,
		logger:   common.GetLogger(),
	}
}

// ForEach 遍历所有文件
func (fc *FileCollection) ForEach(handler func(*FileInfo) error) error {
	fc.iterator.Reset()
	
	for fc.iterator.HasNext() {
		file := fc.iterator.Next()
		if err := handler(file); err != nil {
			return err
		}
	}
	
	return nil
}

// Map 映射文件
func (fc *FileCollection) Map(mapper func(*FileInfo) interface{}) []interface{} {
	var result []interface{}
	
	fc.iterator.Reset()
	for fc.iterator.HasNext() {
		file := fc.iterator.Next()
		result = append(result, mapper(file))
	}
	
	return result
}

// Filter 过滤文件
func (fc *FileCollection) Filter(filter func(*FileInfo) bool) []*FileInfo {
	var result []*FileInfo
	
	fc.iterator.Reset()
	for fc.iterator.HasNext() {
		file := fc.iterator.Next()
		if filter(file) {
			result = append(result, file)
		}
	}
	
	return result
}

// Reduce 归约文件
func (fc *FileCollection) Reduce(initial interface{}, reducer func(interface{}, *FileInfo) interface{}) interface{} {
	result := initial
	
	fc.iterator.Reset()
	for fc.iterator.HasNext() {
		file := fc.iterator.Next()
		result = reducer(result, file)
	}
	
	return result
}

// GetStats 获取统计信息
func (fc *FileCollection) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_files": 0,
		"total_size":  int64(0),
		"directories": 0,
		"files":       0,
	}
	
	fc.iterator.Reset()
	for fc.iterator.HasNext() {
		file := fc.iterator.Next()
		stats["total_files"] = stats["total_files"].(int) + 1
		stats["total_size"] = stats["total_size"].(int64) + file.Size
		
		if file.IsDir {
			stats["directories"] = stats["directories"].(int) + 1
		} else {
			stats["files"] = stats["files"].(int) + 1
		}
	}
	
	return stats
}
