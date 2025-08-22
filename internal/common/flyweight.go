package common

import (
	"fmt"
	"sync"
)

// FlyweightFactory 享元工厂
type FlyweightFactory struct {
	flyweights map[string]Flyweight
	mutex      sync.RWMutex
}

// NewFlyweightFactory 创建享元工厂
func NewFlyweightFactory() *FlyweightFactory {
	return &FlyweightFactory{
		flyweights: make(map[string]Flyweight),
	}
}

// GetFlyweight 获取享元对象
func (ff *FlyweightFactory) GetFlyweight(key string) (Flyweight, error) {
	ff.mutex.RLock()
	if flyweight, exists := ff.flyweights[key]; exists {
		ff.mutex.RUnlock()
		return flyweight, nil
	}
	ff.mutex.RUnlock()

	ff.mutex.Lock()
	defer ff.mutex.Unlock()

	// 双重检查
	if flyweight, exists := ff.flyweights[key]; exists {
		return flyweight, nil
	}

	// 创建新的享元对象
	flyweight, err := ff.createFlyweight(key)
	if err != nil {
		return nil, err
	}

	ff.flyweights[key] = flyweight
	return flyweight, nil
}

// createFlyweight 创建享元对象
func (ff *FlyweightFactory) createFlyweight(key string) (Flyweight, error) {
	// 根据 key 创建不同类型的享元对象
	switch {
	case key == "default_theme":
		return NewThemeFlyweight("default", "#2c3e50", "#3498db", "#ffffff", "#333333"), nil
	case key == "dark_theme":
		return NewThemeFlyweight("dark", "#1a1a1a", "#4a90e2", "#2d2d2d", "#ffffff"), nil
	case key == "light_theme":
		return NewThemeFlyweight("light", "#f8f9fa", "#007bff", "#ffffff", "#212529"), nil
	default:
		return nil, fmt.Errorf("未知的享元类型: %s", key)
	}
}

// GetFlyweightCount 获取享元对象数量
func (ff *FlyweightFactory) GetFlyweightCount() int {
	ff.mutex.RLock()
	defer ff.mutex.RUnlock()
	return len(ff.flyweights)
}

// Clear 清空享元对象
func (ff *FlyweightFactory) Clear() {
	ff.mutex.Lock()
	defer ff.mutex.Unlock()
	ff.flyweights = make(map[string]Flyweight)
}

// Flyweight 享元接口
type Flyweight interface {
	GetKey() string
	GetType() string
	GetData() map[string]interface{}
}

// ThemeFlyweight 主题享元
type ThemeFlyweight struct {
	key           string
	primaryColor  string
	secondaryColor string
	backgroundColor string
	textColor     string
}

// NewThemeFlyweight 创建主题享元
func NewThemeFlyweight(key, primary, secondary, background, text string) *ThemeFlyweight {
	return &ThemeFlyweight{
		key:            key,
		primaryColor:   primary,
		secondaryColor: secondary,
		backgroundColor: background,
		textColor:      text,
	}
}

func (tf *ThemeFlyweight) GetKey() string {
	return tf.key
}

func (tf *ThemeFlyweight) GetType() string {
	return "theme"
}

func (tf *ThemeFlyweight) GetData() map[string]interface{} {
	return map[string]interface{}{
		"primary_color":    tf.primaryColor,
		"secondary_color":  tf.secondaryColor,
		"background_color": tf.backgroundColor,
		"text_color":       tf.textColor,
	}
}

// FontFlyweight 字体享元
type FontFlyweight struct {
	key      string
	family   string
	size     string
	weight   string
	style    string
}

// NewFontFlyweight 创建字体享元
func NewFontFlyweight(key, family, size, weight, style string) *FontFlyweight {
	return &FontFlyweight{
		key:    key,
		family: family,
		size:   size,
		weight: weight,
		style:  style,
	}
}

func (ff *FontFlyweight) GetKey() string {
	return ff.key
}

func (ff *FontFlyweight) GetType() string {
	return "font"
}

func (ff *FontFlyweight) GetData() map[string]interface{} {
	return map[string]interface{}{
		"family": ff.family,
		"size":   ff.size,
		"weight": ff.weight,
		"style":  ff.style,
	}
}

// CSSFlyweight CSS 享元
type CSSFlyweight struct {
	key  string
	css  string
}

// NewCSSFlyweight 创建 CSS 享元
func NewCSSFlyweight(key, css string) *CSSFlyweight {
	return &CSSFlyweight{
		key: key,
		css: css,
	}
}

func (cf *CSSFlyweight) GetKey() string {
	return cf.key
}

func (cf *CSSFlyweight) GetType() string {
	return "css"
}

func (cf *CSSFlyweight) GetData() map[string]interface{} {
	return map[string]interface{}{
		"css": cf.css,
	}
}

// FlyweightContext 享元上下文
type FlyweightContext struct {
	flyweight Flyweight
	uniqueState map[string]interface{}
}

// NewFlyweightContext 创建享元上下文
func NewFlyweightContext(flyweight Flyweight) *FlyweightContext {
	return &FlyweightContext{
		flyweight:    flyweight,
		uniqueState:  make(map[string]interface{}),
	}
}

// SetUniqueState 设置唯一状态
func (fc *FlyweightContext) SetUniqueState(key string, value interface{}) {
	fc.uniqueState[key] = value
}

// GetUniqueState 获取唯一状态
func (fc *FlyweightContext) GetUniqueState(key string) interface{} {
	return fc.uniqueState[key]
}

// GetSharedState 获取共享状态
func (fc *FlyweightContext) GetSharedState() map[string]interface{} {
	return fc.flyweight.GetData()
}

// GetAllState 获取所有状态
func (fc *FlyweightContext) GetAllState() map[string]interface{} {
	allState := make(map[string]interface{})
	
	// 添加共享状态
	for key, value := range fc.flyweight.GetData() {
		allState[key] = value
	}
	
	// 添加唯一状态
	for key, value := range fc.uniqueState {
		allState[key] = value
	}
	
	return allState
}

// FlyweightManager 享元管理器
type FlyweightManager struct {
	factory *FlyweightFactory
	logger  *Logger
}

// NewFlyweightManager 创建享元管理器
func NewFlyweightManager() *FlyweightManager {
	return &FlyweightManager{
		factory: NewFlyweightFactory(),
		logger:  GetLogger(),
	}
}

// GetTheme 获取主题
func (fm *FlyweightManager) GetTheme(themeKey string) (*FlyweightContext, error) {
	flyweight, err := fm.factory.GetFlyweight(themeKey)
	if err != nil {
		return nil, err
	}
	
	return NewFlyweightContext(flyweight), nil
}

// GetFont 获取字体
func (fm *FlyweightManager) GetFont(fontKey string) (*FlyweightContext, error) {
	flyweight, err := fm.factory.GetFlyweight(fontKey)
	if err != nil {
		return nil, err
	}
	
	return NewFlyweightContext(flyweight), nil
}

// GetCSS 获取 CSS
func (fm *FlyweightManager) GetCSS(cssKey string) (*FlyweightContext, error) {
	flyweight, err := fm.factory.GetFlyweight(cssKey)
	if err != nil {
		return nil, err
	}
	
	return NewFlyweightContext(flyweight), nil
}

// GetStatistics 获取统计信息
func (fm *FlyweightManager) GetStatistics() map[string]interface{} {
	return map[string]interface{}{
		"total_flyweights": fm.factory.GetFlyweightCount(),
		"memory_saved":     "通过享元模式节省内存使用",
	}
}

// CacheManager 缓存管理器（使用享元模式）
type CacheManager struct {
	flyweightManager *FlyweightManager
	cache           map[string]*FlyweightContext
	mutex           sync.RWMutex
	logger          *Logger
}

// NewCacheManager 创建缓存管理器
func NewCacheManager() *CacheManager {
	return &CacheManager{
		flyweightManager: NewFlyweightManager(),
		cache:           make(map[string]*FlyweightContext),
		logger:          GetLogger(),
	}
}

// GetOrCreate 获取或创建缓存项
func (cm *CacheManager) GetOrCreate(key string, creator func() (*FlyweightContext, error)) (*FlyweightContext, error) {
	cm.mutex.RLock()
	if context, exists := cm.cache[key]; exists {
		cm.mutex.RUnlock()
		return context, nil
	}
	cm.mutex.RUnlock()

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 双重检查
	if context, exists := cm.cache[key]; exists {
		return context, nil
	}

	// 创建新的上下文
	context, err := creator()
	if err != nil {
		return nil, err
	}

	cm.cache[key] = context
	cm.logger.Info("创建新的缓存项:", key)
	
	return context, nil
}

// Get 获取缓存项
func (cm *CacheManager) Get(key string) (*FlyweightContext, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	context, exists := cm.cache[key]
	return context, exists
}

// Set 设置缓存项
func (cm *CacheManager) Set(key string, context *FlyweightContext) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	cm.cache[key] = context
}

// Remove 移除缓存项
func (cm *CacheManager) Remove(key string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	delete(cm.cache, key)
}

// Clear 清空缓存
func (cm *CacheManager) Clear() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	cm.cache = make(map[string]*FlyweightContext)
}

// GetCacheSize 获取缓存大小
func (cm *CacheManager) GetCacheSize() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	return len(cm.cache)
}

// GetCacheKeys 获取缓存键列表
func (cm *CacheManager) GetCacheKeys() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	keys := make([]string, 0, len(cm.cache))
	for key := range cm.cache {
		keys = append(keys, key)
	}
	
	return keys
}
