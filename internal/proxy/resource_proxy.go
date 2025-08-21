package proxy

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"creeper/internal/common"
)

// ResourceInterface 资源接口
type ResourceInterface interface {
	Load(path string) ([]byte, error)
	Save(path string, data []byte) error
	Exists(path string) bool
	GetInfo(path string) (os.FileInfo, error)
}

// RealResource 真实资源
type RealResource struct{}

func (rr *RealResource) Load(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (rr *RealResource) Save(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (rr *RealResource) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (rr *RealResource) GetInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// CachedResource 缓存资源代理
type CachedResource struct {
	realResource ResourceInterface
	cache        map[string]CacheEntry
	mutex        sync.RWMutex
	maxCacheSize int
	logger       *common.Logger
}

// CacheEntry 缓存条目
type CacheEntry struct {
	Data      []byte
	Timestamp time.Time
	FileInfo  os.FileInfo
}

// NewCachedResource 创建缓存资源代理
func NewCachedResource(maxCacheSize int) *CachedResource {
	return &CachedResource{
		realResource: &RealResource{},
		cache:        make(map[string]CacheEntry),
		maxCacheSize: maxCacheSize,
		logger:       common.GetLogger(),
	}
}

func (cr *CachedResource) Load(path string) ([]byte, error) {
	cr.mutex.RLock()
	
	// 检查缓存
	if entry, exists := cr.cache[path]; exists {
		// 检查文件是否有更新
		if info, err := cr.realResource.GetInfo(path); err == nil {
			if info.ModTime().Equal(entry.FileInfo.ModTime()) {
				cr.mutex.RUnlock()
				cr.logger.Debug("从缓存加载:", path)
				return entry.Data, nil
			}
		}
	}
	
	cr.mutex.RUnlock()
	
	// 缓存未命中或文件已更新，从真实资源加载
	data, err := cr.realResource.Load(path)
	if err != nil {
		return nil, err
	}
	
	// 更新缓存
	if info, err := cr.realResource.GetInfo(path); err == nil {
		cr.updateCache(path, data, info)
	}
	
	cr.logger.Debug("从文件系统加载:", path)
	return data, nil
}

func (cr *CachedResource) Save(path string, data []byte) error {
	err := cr.realResource.Save(path, data)
	if err != nil {
		return err
	}
	
	// 更新缓存
	if info, err := cr.realResource.GetInfo(path); err == nil {
		cr.updateCache(path, data, info)
	}
	
	return nil
}

func (cr *CachedResource) Exists(path string) bool {
	cr.mutex.RLock()
	
	// 检查缓存
	if _, exists := cr.cache[path]; exists {
		cr.mutex.RUnlock()
		return true
	}
	
	cr.mutex.RUnlock()
	
	// 检查真实资源
	return cr.realResource.Exists(path)
}

func (cr *CachedResource) GetInfo(path string) (os.FileInfo, error) {
	cr.mutex.RLock()
	
	// 检查缓存
	if entry, exists := cr.cache[path]; exists {
		cr.mutex.RUnlock()
		return entry.FileInfo, nil
	}
	
	cr.mutex.RUnlock()
	
	// 从真实资源获取
	return cr.realResource.GetInfo(path)
}

// updateCache 更新缓存
func (cr *CachedResource) updateCache(path string, data []byte, info os.FileInfo) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	
	// 检查缓存大小
	if len(cr.cache) >= cr.maxCacheSize {
		cr.evictOldestEntry()
	}
	
	cr.cache[path] = CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
		FileInfo:  info,
	}
}

// evictOldestEntry 淘汰最旧的缓存条目
func (cr *CachedResource) evictOldestEntry() {
	var oldestPath string
	var oldestTime time.Time
	
	for path, entry := range cr.cache {
		if oldestPath == "" || entry.Timestamp.Before(oldestTime) {
			oldestPath = path
			oldestTime = entry.Timestamp
		}
	}
	
	if oldestPath != "" {
		delete(cr.cache, oldestPath)
		cr.logger.Debug("淘汰缓存条目:", oldestPath)
	}
}

// GetCacheStats 获取缓存统计
func (cr *CachedResource) GetCacheStats() map[string]interface{} {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()
	
	return map[string]interface{}{
		"cached_files": len(cr.cache),
		"max_size":     cr.maxCacheSize,
		"cache_paths":  cr.getCachePaths(),
	}
}

// getCachePaths 获取缓存路径列表
func (cr *CachedResource) getCachePaths() []string {
	paths := make([]string, 0, len(cr.cache))
	for path := range cr.cache {
		paths = append(paths, path)
	}
	return paths
}

// ClearCache 清空缓存
func (cr *CachedResource) ClearCache() {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	
	cr.cache = make(map[string]CacheEntry)
	cr.logger.Info("缓存已清空")
}

// SecurityResourceProxy 安全资源代理
type SecurityResourceProxy struct {
	realResource   ResourceInterface
	allowedPaths   []string
	blockedPaths   []string
	logger         *common.Logger
}

// NewSecurityResourceProxy 创建安全资源代理
func NewSecurityResourceProxy(realResource ResourceInterface) *SecurityResourceProxy {
	return &SecurityResourceProxy{
		realResource: realResource,
		allowedPaths: []string{},
		blockedPaths: []string{},
		logger:       common.GetLogger(),
	}
}

// AddAllowedPath 添加允许的路径
func (srp *SecurityResourceProxy) AddAllowedPath(path string) {
	srp.allowedPaths = append(srp.allowedPaths, path)
}

// AddBlockedPath 添加阻止的路径
func (srp *SecurityResourceProxy) AddBlockedPath(path string) {
	srp.blockedPaths = append(srp.blockedPaths, path)
}

func (srp *SecurityResourceProxy) Load(path string) ([]byte, error) {
	if !srp.isPathAllowed(path) {
		srp.logger.Warn("访问被拒绝:", path)
		return nil, fmt.Errorf("访问路径被拒绝: %s", path)
	}
	
	return srp.realResource.Load(path)
}

func (srp *SecurityResourceProxy) Save(path string, data []byte) error {
	if !srp.isPathAllowed(path) {
		srp.logger.Warn("写入被拒绝:", path)
		return fmt.Errorf("写入路径被拒绝: %s", path)
	}
	
	return srp.realResource.Save(path, data)
}

func (srp *SecurityResourceProxy) Exists(path string) bool {
	if !srp.isPathAllowed(path) {
		return false
	}
	
	return srp.realResource.Exists(path)
}

func (srp *SecurityResourceProxy) GetInfo(path string) (os.FileInfo, error) {
	if !srp.isPathAllowed(path) {
		return nil, fmt.Errorf("访问路径被拒绝: %s", path)
	}
	
	return srp.realResource.GetInfo(path)
}

// isPathAllowed 检查路径是否被允许
func (srp *SecurityResourceProxy) isPathAllowed(path string) bool {
	// 检查阻止列表
	for _, blocked := range srp.blockedPaths {
		if strings.Contains(path, blocked) {
			return false
		}
	}
	
	// 如果有允许列表，检查是否在允许列表中
	if len(srp.allowedPaths) > 0 {
		for _, allowed := range srp.allowedPaths {
			if strings.Contains(path, allowed) {
				return true
			}
		}
		return false
	}
	
	// 默认允许
	return true
}

// ResourceProxyFactory 资源代理工厂
type ResourceProxyFactory struct{}

// CreateCachedProxy 创建缓存代理
func (rpf *ResourceProxyFactory) CreateCachedProxy(maxCacheSize int) ResourceInterface {
	return NewCachedResource(maxCacheSize)
}

// CreateSecurityProxy 创建安全代理
func (rpf *ResourceProxyFactory) CreateSecurityProxy(realResource ResourceInterface, allowedPaths, blockedPaths []string) ResourceInterface {
	proxy := NewSecurityResourceProxy(realResource)
	
	for _, path := range allowedPaths {
		proxy.AddAllowedPath(path)
	}
	
	for _, path := range blockedPaths {
		proxy.AddBlockedPath(path)
	}
	
	return proxy
}

// CreateCompositeProxy 创建复合代理
func (rpf *ResourceProxyFactory) CreateCompositeProxy(maxCacheSize int, allowedPaths, blockedPaths []string) ResourceInterface {
	// 创建缓存代理
	cachedProxy := rpf.CreateCachedProxy(maxCacheSize)
	
	// 在缓存代理外包装安全代理
	return rpf.CreateSecurityProxy(cachedProxy, allowedPaths, blockedPaths)
}
