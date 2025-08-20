package common

import (
	"log"
	"sync"
)

// Logger 全局日志记录器
type Logger struct {
	logger *log.Logger
}

var (
	loggerInstance *Logger
	loggerOnce     sync.Once
)

// GetLogger 获取日志记录器单例
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		loggerInstance = &Logger{
			logger: log.Default(),
		}
	})
	return loggerInstance
}

// Info 记录信息日志
func (l *Logger) Info(v ...interface{}) {
	l.logger.Println("[INFO]", v)
}

// Warn 记录警告日志
func (l *Logger) Warn(v ...interface{}) {
	l.logger.Println("[WARN]", v)
}

// Error 记录错误日志
func (l *Logger) Error(v ...interface{}) {
	l.logger.Println("[ERROR]", v)
}

// Debug 记录调试日志
func (l *Logger) Debug(v ...interface{}) {
	l.logger.Println("[DEBUG]", v)
}

// ResourceManager 资源管理器
type ResourceManager struct {
	resources map[string]interface{}
	mutex     sync.RWMutex
}

var (
	resourceInstance *ResourceManager
	resourceOnce     sync.Once
)

// GetResourceManager 获取资源管理器单例
func GetResourceManager() *ResourceManager {
	resourceOnce.Do(func() {
		resourceInstance = &ResourceManager{
			resources: make(map[string]interface{}),
		}
	})
	return resourceInstance
}

// Set 设置资源
func (rm *ResourceManager) Set(key string, value interface{}) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.resources[key] = value
}

// Get 获取资源
func (rm *ResourceManager) Get(key string) (interface{}, bool) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	value, exists := rm.resources[key]
	return value, exists
}

// GetString 获取字符串资源
func (rm *ResourceManager) GetString(key string) (string, bool) {
	value, exists := rm.Get(key)
	if !exists {
		return "", false
	}
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// GetInt 获取整数资源
func (rm *ResourceManager) GetInt(key string) (int, bool) {
	value, exists := rm.Get(key)
	if !exists {
		return 0, false
	}
	if num, ok := value.(int); ok {
		return num, true
	}
	return 0, false
}

// Delete 删除资源
func (rm *ResourceManager) Delete(key string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	delete(rm.resources, key)
}

// Clear 清空所有资源
func (rm *ResourceManager) Clear() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.resources = make(map[string]interface{})
}

// Keys 获取所有键
func (rm *ResourceManager) Keys() []string {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	keys := make([]string, 0, len(rm.resources))
	for key := range rm.resources {
		keys = append(keys, key)
	}
	return keys
}

// Count 获取资源数量
func (rm *ResourceManager) Count() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	return len(rm.resources)
}

// ConfigCache 配置缓存
type ConfigCache struct {
	cache map[string]interface{}
	mutex sync.RWMutex
}

var (
	configInstance *ConfigCache
	configOnce     sync.Once
)

// GetConfigCache 获取配置缓存单例
func GetConfigCache() *ConfigCache {
	configOnce.Do(func() {
		configInstance = &ConfigCache{
			cache: make(map[string]interface{}),
		}
	})
	return configInstance
}

// Set 设置缓存
func (cc *ConfigCache) Set(key string, value interface{}) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	cc.cache[key] = value
}

// Get 获取缓存
func (cc *ConfigCache) Get(key string) (interface{}, bool) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()
	value, exists := cc.cache[key]
	return value, exists
}

// Delete 删除缓存
func (cc *ConfigCache) Delete(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	delete(cc.cache, key)
}

// Clear 清空缓存
func (cc *ConfigCache) Clear() {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	cc.cache = make(map[string]interface{})
}

// Exists 检查键是否存在
func (cc *ConfigCache) Exists(key string) bool {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()
	_, exists := cc.cache[key]
	return exists
}
