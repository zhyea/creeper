package common

import (
	"fmt"
	"sync"
	"time"
)

// ResourceComponent 资源组件接口
type ResourceComponent interface {
	GetName() string
	GetType() string
	GetSize() int64
	GetCreatedTime() time.Time
	GetModifiedTime() time.Time
	IsDirectory() bool
	GetChildren() []ResourceComponent
	AddChild(component ResourceComponent) error
	RemoveChild(name string) error
	GetChild(name string) ResourceComponent
	GetPath() string
	GetMetadata() map[string]interface{}
}

// FileResource 文件资源
type FileResource struct {
	name         string
	path         string
	size         int64
	createdTime  time.Time
	modifiedTime time.Time
	metadata     map[string]interface{}
	mutex        sync.RWMutex
}

// NewFileResource 创建文件资源
func NewFileResource(name, path string, size int64) *FileResource {
	now := time.Now()
	return &FileResource{
		name:         name,
		path:         path,
		size:         size,
		createdTime:  now,
		modifiedTime: now,
		metadata:     make(map[string]interface{}),
	}
}

func (fr *FileResource) GetName() string {
	return fr.name
}

func (fr *FileResource) GetType() string {
	return "file"
}

func (fr *FileResource) GetSize() int64 {
	return fr.size
}

func (fr *FileResource) GetCreatedTime() time.Time {
	return fr.createdTime
}

func (fr *FileResource) GetModifiedTime() time.Time {
	return fr.modifiedTime
}

func (fr *FileResource) IsDirectory() bool {
	return false
}

func (fr *FileResource) GetChildren() []ResourceComponent {
	return nil
}

func (fr *FileResource) AddChild(component ResourceComponent) error {
	return fmt.Errorf("文件资源不能添加子组件")
}

func (fr *FileResource) RemoveChild(name string) error {
	return fmt.Errorf("文件资源不能移除子组件")
}

func (fr *FileResource) GetChild(name string) ResourceComponent {
	return nil
}

func (fr *FileResource) GetPath() string {
	return fr.path
}

func (fr *FileResource) GetMetadata() map[string]interface{} {
	fr.mutex.RLock()
	defer fr.mutex.RUnlock()
	
	metadata := make(map[string]interface{})
	for key, value := range fr.metadata {
		metadata[key] = value
	}
	return metadata
}

// SetMetadata 设置元数据
func (fr *FileResource) SetMetadata(key string, value interface{}) {
	fr.mutex.Lock()
	defer fr.mutex.Unlock()
	
	fr.metadata[key] = value
	fr.modifiedTime = time.Now()
}

// DirectoryResource 目录资源
type DirectoryResource struct {
	name         string
	path         string
	children     map[string]ResourceComponent
	createdTime  time.Time
	modifiedTime time.Time
	metadata     map[string]interface{}
	mutex        sync.RWMutex
}

// NewDirectoryResource 创建目录资源
func NewDirectoryResource(name, path string) *DirectoryResource {
	now := time.Now()
	return &DirectoryResource{
		name:         name,
		path:         path,
		children:     make(map[string]ResourceComponent),
		createdTime:  now,
		modifiedTime: now,
		metadata:     make(map[string]interface{}),
	}
}

func (dr *DirectoryResource) GetName() string {
	return dr.name
}

func (dr *DirectoryResource) GetType() string {
	return "directory"
}

func (dr *DirectoryResource) GetSize() int64 {
	dr.mutex.RLock()
	defer dr.mutex.RUnlock()
	
	var totalSize int64
	for _, child := range dr.children {
		totalSize += child.GetSize()
	}
	return totalSize
}

func (dr *DirectoryResource) GetCreatedTime() time.Time {
	return dr.createdTime
}

func (dr *DirectoryResource) GetModifiedTime() time.Time {
	return dr.modifiedTime
}

func (dr *DirectoryResource) IsDirectory() bool {
	return true
}

func (dr *DirectoryResource) GetChildren() []ResourceComponent {
	dr.mutex.RLock()
	defer dr.mutex.RUnlock()
	
	children := make([]ResourceComponent, 0, len(dr.children))
	for _, child := range dr.children {
		children = append(children, child)
	}
	return children
}

func (dr *DirectoryResource) AddChild(component ResourceComponent) error {
	dr.mutex.Lock()
	defer dr.mutex.Unlock()
	
	if _, exists := dr.children[component.GetName()]; exists {
		return fmt.Errorf("子组件已存在: %s", component.GetName())
	}
	
	dr.children[component.GetName()] = component
	dr.modifiedTime = time.Now()
	return nil
}

func (dr *DirectoryResource) RemoveChild(name string) error {
	dr.mutex.Lock()
	defer dr.mutex.Unlock()
	
	if _, exists := dr.children[name]; !exists {
		return fmt.Errorf("子组件不存在: %s", name)
	}
	
	delete(dr.children, name)
	dr.modifiedTime = time.Now()
	return nil
}

func (dr *DirectoryResource) GetChild(name string) ResourceComponent {
	dr.mutex.RLock()
	defer dr.mutex.RUnlock()
	
	return dr.children[name]
}

func (dr *DirectoryResource) GetPath() string {
	return dr.path
}

func (dr *DirectoryResource) GetMetadata() map[string]interface{} {
	dr.mutex.RLock()
	defer dr.mutex.RUnlock()
	
	metadata := make(map[string]interface{})
	for key, value := range dr.metadata {
		metadata[key] = value
	}
	return metadata
}

// SetMetadata 设置元数据
func (dr *DirectoryResource) SetMetadata(key string, value interface{}) {
	dr.mutex.Lock()
	defer dr.mutex.Unlock()
	
	dr.metadata[key] = value
	dr.modifiedTime = time.Now()
}

// GetChildCount 获取子组件数量
func (dr *DirectoryResource) GetChildCount() int {
	dr.mutex.RLock()
	defer dr.mutex.RUnlock()
	
	return len(dr.children)
}

// GetChildByType 按类型获取子组件
func (dr *DirectoryResource) GetChildByType(componentType string) []ResourceComponent {
	dr.mutex.RLock()
	defer dr.mutex.RUnlock()
	
	var result []ResourceComponent
	for _, child := range dr.children {
		if child.GetType() == componentType {
			result = append(result, child)
		}
	}
	return result
}

// ResourceTree 资源树
type ResourceTree struct {
	root   ResourceComponent
	logger *Logger
}

// NewResourceTree 创建资源树
func NewResourceTree(root ResourceComponent) *ResourceTree {
	return &ResourceTree{
		root:   root,
		logger: GetLogger(),
	}
}

// Traverse 遍历资源树
func (rt *ResourceTree) Traverse(visitor func(ResourceComponent, int) error) error {
	return rt.traverseNode(rt.root, 0, visitor)
}

// traverseNode 遍历节点
func (rt *ResourceTree) traverseNode(node ResourceComponent, depth int, visitor func(ResourceComponent, int) error) error {
	if err := visitor(node, depth); err != nil {
		return err
	}
	
	if node.IsDirectory() {
		for _, child := range node.GetChildren() {
			if err := rt.traverseNode(child, depth+1, visitor); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// FindByPath 按路径查找资源
func (rt *ResourceTree) FindByPath(path string) ResourceComponent {
	return rt.findByPathRecursive(rt.root, path)
}

// findByPathRecursive 递归查找路径
func (rt *ResourceTree) findByPathRecursive(node ResourceComponent, path string) ResourceComponent {
	if node.GetPath() == path {
		return node
	}
	
	if node.IsDirectory() {
		for _, child := range node.GetChildren() {
			if result := rt.findByPathRecursive(child, path); result != nil {
				return result
			}
		}
	}
	
	return nil
}

// GetStatistics 获取统计信息
func (rt *ResourceTree) GetStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"total_files":     0,
		"total_directories": 0,
		"total_size":      int64(0),
		"max_depth":       0,
	}
	
	rt.Traverse(func(component ResourceComponent, depth int) error {
		if component.IsDirectory() {
			stats["total_directories"] = stats["total_directories"].(int) + 1
		} else {
			stats["total_files"] = stats["total_files"].(int) + 1
			stats["total_size"] = stats["total_size"].(int64) + component.GetSize()
		}
		
		if depth > stats["max_depth"].(int) {
			stats["max_depth"] = depth
		}
		
		return nil
	})
	
	return stats
}

// ResourceManager 资源管理器
type ResourceManager struct {
	tree   *ResourceTree
	logger *Logger
}

// NewResourceManager 创建资源管理器
func NewResourceManager() *ResourceManager {
	root := NewDirectoryResource("root", "/")
	tree := NewResourceTree(root)
	
	return &ResourceManager{
		tree:   tree,
		logger: GetLogger(),
	}
}

// AddResource 添加资源
func (rm *ResourceManager) AddResource(path string, size int64) error {
	// 解析路径
	components := rm.parsePath(path)
	
	// 创建或获取父目录
	parent := rm.tree.root
	for i := 0; i < len(components)-1; i++ {
		child := parent.GetChild(components[i])
		if child == nil {
			child = NewDirectoryResource(components[i], rm.buildPath(components[:i+1]))
			parent.AddChild(child)
		}
		parent = child
	}
	
	// 创建文件资源
	fileName := components[len(components)-1]
	fileResource := NewFileResource(fileName, path, size)
	
	return parent.AddChild(fileResource)
}

// RemoveResource 移除资源
func (rm *ResourceManager) RemoveResource(path string) error {
	// 解析路径
	components := rm.parsePath(path)
	
	// 查找父目录
	parent := rm.tree.root
	for i := 0; i < len(components)-1; i++ {
		child := parent.GetChild(components[i])
		if child == nil {
			return fmt.Errorf("路径不存在: %s", path)
		}
		parent = child
	}
	
	// 移除资源
	fileName := components[len(components)-1]
	return parent.RemoveChild(fileName)
}

// GetResource 获取资源
func (rm *ResourceManager) GetResource(path string) ResourceComponent {
	return rm.tree.FindByPath(path)
}

// GetStatistics 获取统计信息
func (rm *ResourceManager) GetStatistics() map[string]interface{} {
	return rm.tree.GetStatistics()
}

// parsePath 解析路径
func (rm *ResourceManager) parsePath(path string) []string {
	// 简化路径解析
	var components []string
	current := ""
	
	for _, char := range path {
		if char == '/' {
			if current != "" {
				components = append(components, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		components = append(components, current)
	}
	
	return components
}

// buildPath 构建路径
func (rm *ResourceManager) buildPath(components []string) string {
	if len(components) == 0 {
		return "/"
	}
	
	path := "/"
	for _, component := range components {
		path += component + "/"
	}
	
	return path[:len(path)-1]
}
