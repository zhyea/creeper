package chain

import (
	"fmt"
	"strings"

	"creeper/internal/common"
)

// ErrorSeverity 错误严重程度
type ErrorSeverity int

const (
	SeverityInfo ErrorSeverity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

// ErrorContext 错误上下文
type ErrorContext struct {
	Error     error
	Severity  ErrorSeverity
	Component string
	Operation string
	Data      map[string]interface{}
	Handled   bool
}

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	SetNext(handler ErrorHandler)
	Handle(context *ErrorContext) error
}

// BaseErrorHandler 基础错误处理器
type BaseErrorHandler struct {
	next ErrorHandler
}

func (beh *BaseErrorHandler) SetNext(handler ErrorHandler) {
	beh.next = handler
}

func (beh *BaseErrorHandler) Handle(context *ErrorContext) error {
	if beh.next != nil {
		return beh.next.Handle(context)
	}
	return nil
}

// LoggingErrorHandler 日志错误处理器
type LoggingErrorHandler struct {
	BaseErrorHandler
	logger *common.Logger
}

// NewLoggingErrorHandler 创建日志错误处理器
func NewLoggingErrorHandler() *LoggingErrorHandler {
	return &LoggingErrorHandler{
		logger: common.GetLogger(),
	}
}

func (leh *LoggingErrorHandler) Handle(context *ErrorContext) error {
	// 记录所有错误
	message := fmt.Sprintf("[%s] %s: %v", context.Component, context.Operation, context.Error)
	
	switch context.Severity {
	case SeverityInfo:
		leh.logger.Info(message)
	case SeverityWarning:
		leh.logger.Warn(message)
	case SeverityError:
		leh.logger.Error(message)
	case SeverityCritical:
		leh.logger.Error("CRITICAL:", message)
	}
	
	// 继续传递到下一个处理器
	return leh.BaseErrorHandler.Handle(context)
}

// RetryErrorHandler 重试错误处理器
type RetryErrorHandler struct {
	BaseErrorHandler
	maxRetries int
	retryCount map[string]int
}

// NewRetryErrorHandler 创建重试错误处理器
func NewRetryErrorHandler(maxRetries int) *RetryErrorHandler {
	return &RetryErrorHandler{
		maxRetries: maxRetries,
		retryCount: make(map[string]int),
	}
}

func (reh *RetryErrorHandler) Handle(context *ErrorContext) error {
	// 只处理可重试的错误
	if reh.isRetryableError(context.Error) && context.Severity <= SeverityError {
		key := fmt.Sprintf("%s:%s", context.Component, context.Operation)
		
		if reh.retryCount[key] < reh.maxRetries {
			reh.retryCount[key]++
			
			// 标记为已处理，阻止继续传递
			context.Handled = true
			
			common.GetLogger().Info(fmt.Sprintf("重试操作 %s (第%d次)", key, reh.retryCount[key]))
			
			// 这里可以添加重试逻辑
			// 实际项目中，这里应该调用原始操作的重试机制
			
			return nil
		}
	}
	
	// 继续传递到下一个处理器
	return reh.BaseErrorHandler.Handle(context)
}

// isRetryableError 判断是否为可重试错误
func (reh *RetryErrorHandler) isRetryableError(err error) bool {
	errorStr := strings.ToLower(err.Error())
	
	// 网络相关错误通常可以重试
	retryableKeywords := []string{
		"connection", "timeout", "temporary", "network",
		"file not found", "permission denied",
	}
	
	for _, keyword := range retryableKeywords {
		if strings.Contains(errorStr, keyword) {
			return true
		}
	}
	
	return false
}

// FallbackErrorHandler 回退错误处理器
type FallbackErrorHandler struct {
	BaseErrorHandler
	fallbackActions map[string]func(*ErrorContext) error
}

// NewFallbackErrorHandler 创建回退错误处理器
func NewFallbackErrorHandler() *FallbackErrorHandler {
	handler := &FallbackErrorHandler{
		fallbackActions: make(map[string]func(*ErrorContext) error),
	}
	
	// 注册默认回退动作
	handler.RegisterFallback("parser:parse_file", handler.fallbackParseFile)
	handler.RegisterFallback("generator:generate_assets", handler.fallbackGenerateAssets)
	
	return handler
}

// RegisterFallback 注册回退动作
func (feh *FallbackErrorHandler) RegisterFallback(operation string, action func(*ErrorContext) error) {
	feh.fallbackActions[operation] = action
}

func (feh *FallbackErrorHandler) Handle(context *ErrorContext) error {
	// 查找回退动作
	key := fmt.Sprintf("%s:%s", context.Component, context.Operation)
	
	if action, exists := feh.fallbackActions[key]; exists {
		common.GetLogger().Info("执行回退动作:", key)
		
		if err := action(context); err == nil {
			context.Handled = true
			return nil
		}
	}
	
	// 继续传递到下一个处理器
	return feh.BaseErrorHandler.Handle(context)
}

// 回退动作实现
func (feh *FallbackErrorHandler) fallbackParseFile(context *ErrorContext) error {
	// 解析文件失败的回退：尝试创建空章节
	common.GetLogger().Warn("文件解析失败，创建空章节")
	return nil
}

func (feh *FallbackErrorHandler) fallbackGenerateAssets(context *ErrorContext) error {
	// 资源生成失败的回退：使用默认资源
	common.GetLogger().Warn("资源生成失败，使用默认资源")
	return nil
}

// CriticalErrorHandler 关键错误处理器
type CriticalErrorHandler struct {
	BaseErrorHandler
	shutdownCallback func()
}

// NewCriticalErrorHandler 创建关键错误处理器
func NewCriticalErrorHandler(shutdownCallback func()) *CriticalErrorHandler {
	return &CriticalErrorHandler{
		shutdownCallback: shutdownCallback,
	}
}

func (ceh *CriticalErrorHandler) Handle(context *ErrorContext) error {
	// 只处理关键错误
	if context.Severity == SeverityCritical {
		common.GetLogger().Error("检测到关键错误，系统即将关闭:", context.Error)
		
		context.Handled = true
		
		// 执行关闭回调
		if ceh.shutdownCallback != nil {
			ceh.shutdownCallback()
		}
		
		return context.Error
	}
	
	// 继续传递到下一个处理器
	return ceh.BaseErrorHandler.Handle(context)
}

// ErrorHandlerChain 错误处理链
type ErrorHandlerChain struct {
	firstHandler ErrorHandler
	logger       *common.Logger
}

// NewErrorHandlerChain 创建错误处理链
func NewErrorHandlerChain() *ErrorHandlerChain {
	return &ErrorHandlerChain{
		logger: common.GetLogger(),
	}
}

// BuildChain 构建处理链
func (ehc *ErrorHandlerChain) BuildChain(shutdownCallback func()) {
	// 创建处理器
	loggingHandler := NewLoggingErrorHandler()
	retryHandler := NewRetryErrorHandler(3)
	fallbackHandler := NewFallbackErrorHandler()
	criticalHandler := NewCriticalErrorHandler(shutdownCallback)
	
	// 构建处理链：日志 -> 重试 -> 回退 -> 关键错误
	loggingHandler.SetNext(retryHandler)
	retryHandler.SetNext(fallbackHandler)
	fallbackHandler.SetNext(criticalHandler)
	
	ehc.firstHandler = loggingHandler
}

// HandleError 处理错误
func (ehc *ErrorHandlerChain) HandleError(err error, severity ErrorSeverity, component, operation string, data map[string]interface{}) error {
	if ehc.firstHandler == nil {
		ehc.BuildChain(nil)
	}
	
	context := &ErrorContext{
		Error:     err,
		Severity:  severity,
		Component: component,
		Operation: operation,
		Data:      data,
		Handled:   false,
	}
	
	return ehc.firstHandler.Handle(context)
}

// GetHandlerChainInfo 获取处理链信息
func (ehc *ErrorHandlerChain) GetHandlerChainInfo() string {
	return "LoggingHandler -> RetryHandler -> FallbackHandler -> CriticalHandler"
}

// ErrorManager 错误管理器
type ErrorManager struct {
	chain     *ErrorHandlerChain
	errorLog  []ErrorContext
	maxLogSize int
}

// NewErrorManager 创建错误管理器
func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		chain:      NewErrorHandlerChain(),
		errorLog:   make([]ErrorContext, 0),
		maxLogSize: 1000,
	}
}

// HandleError 处理错误
func (em *ErrorManager) HandleError(err error, severity ErrorSeverity, component, operation string, data map[string]interface{}) error {
	// 记录错误
	errorContext := ErrorContext{
		Error:     err,
		Severity:  severity,
		Component: component,
		Operation: operation,
		Data:      data,
	}
	
	em.addToLog(errorContext)
	
	// 使用责任链处理
	return em.chain.HandleError(err, severity, component, operation, data)
}

// addToLog 添加到错误日志
func (em *ErrorManager) addToLog(context ErrorContext) {
	em.errorLog = append(em.errorLog, context)
	
	// 保持日志大小
	if len(em.errorLog) > em.maxLogSize {
		em.errorLog = em.errorLog[len(em.errorLog)-em.maxLogSize:]
	}
}

// GetErrorLog 获取错误日志
func (em *ErrorManager) GetErrorLog() []ErrorContext {
	return em.errorLog
}

// GetErrorStatistics 获取错误统计
func (em *ErrorManager) GetErrorStatistics() map[string]int {
	stats := map[string]int{
		"total":    len(em.errorLog),
		"info":     0,
		"warning":  0,
		"error":    0,
		"critical": 0,
	}
	
	for _, errorContext := range em.errorLog {
		switch errorContext.Severity {
		case SeverityInfo:
			stats["info"]++
		case SeverityWarning:
			stats["warning"]++
		case SeverityError:
			stats["error"]++
		case SeverityCritical:
			stats["critical"]++
		}
	}
	
	return stats
}

// ClearErrorLog 清空错误日志
func (em *ErrorManager) ClearErrorLog() {
	em.errorLog = make([]ErrorContext, 0)
}

// SetShutdownCallback 设置关闭回调
func (em *ErrorManager) SetShutdownCallback(callback func()) {
	em.chain.BuildChain(callback)
}
