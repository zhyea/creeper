package deploy

import (
	"fmt"
	"sync"
	"time"

	"creeper/internal/common"
)

// DeploymentEvent 部署事件
type DeploymentEvent struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Error     error                  `json:"error,omitempty"`
}

// DeploymentObserver 部署观察者接口
type DeploymentObserver interface {
	OnDeploymentEvent(event *DeploymentEvent)
	GetName() string
}

// DeploymentSubject 部署主题接口
type DeploymentSubject interface {
	Attach(observer DeploymentObserver)
	Detach(observer DeploymentObserver)
	Notify(event *DeploymentEvent)
}

// DeploymentEventManager 部署事件管理器
type DeploymentEventManager struct {
	observers map[string]DeploymentObserver
	mutex     sync.RWMutex
	logger    *common.Logger
}

// NewDeploymentEventManager 创建部署事件管理器
func NewDeploymentEventManager() *DeploymentEventManager {
	return &DeploymentEventManager{
		observers: make(map[string]DeploymentObserver),
		logger:    common.GetLogger(),
	}
}

// Attach 添加观察者
func (dem *DeploymentEventManager) Attach(observer DeploymentObserver) {
	dem.mutex.Lock()
	defer dem.mutex.Unlock()
	
	dem.observers[observer.GetName()] = observer
	dem.logger.Info("添加部署观察者:", observer.GetName())
}

// Detach 移除观察者
func (dem *DeploymentEventManager) Detach(observer DeploymentObserver) {
	dem.mutex.Lock()
	defer dem.mutex.Unlock()
	
	delete(dem.observers, observer.GetName())
	dem.logger.Info("移除部署观察者:", observer.GetName())
}

// Notify 通知所有观察者
func (dem *DeploymentEventManager) Notify(event *DeploymentEvent) {
	dem.mutex.RLock()
	defer dem.mutex.RUnlock()
	
	dem.logger.Info("通知部署事件:", event.Type, "观察者数量:", len(dem.observers))
	
	for name, observer := range dem.observers {
		go func(obs DeploymentObserver, evt *DeploymentEvent) {
			defer func() {
				if r := recover(); r != nil {
					dem.logger.Error("观察者处理事件时发生错误:", name, r)
				}
			}()
			
			obs.OnDeploymentEvent(evt)
		}(observer, event)
	}
}

// GetObserverCount 获取观察者数量
func (dem *DeploymentEventManager) GetObserverCount() int {
	dem.mutex.RLock()
	defer dem.mutex.RUnlock()
	
	return len(dem.observers)
}

// GetObserverNames 获取观察者名称列表
func (dem *DeploymentEventManager) GetObserverNames() []string {
	dem.mutex.RLock()
	defer dem.mutex.RUnlock()
	
	names := make([]string, 0, len(dem.observers))
	for name := range dem.observers {
		names = append(names, name)
	}
	
	return names
}

// ConsoleObserver 控制台观察者
type ConsoleObserver struct {
	name   string
	logger *common.Logger
}

// NewConsoleObserver 创建控制台观察者
func NewConsoleObserver(name string) *ConsoleObserver {
	return &ConsoleObserver{
		name:   name,
		logger: common.GetLogger(),
	}
}

func (co *ConsoleObserver) OnDeploymentEvent(event *DeploymentEvent) {
	message := fmt.Sprintf("[%s] %s", event.Type, event.Timestamp.Format("2006-01-02 15:04:05"))
	
	if event.Error != nil {
		co.logger.Error(message, "错误:", event.Error)
	} else {
		co.logger.Info(message)
	}
	
	// 输出详细信息
	for key, value := range event.Data {
		co.logger.Debug(fmt.Sprintf("  %s: %v", key, value))
	}
}

func (co *ConsoleObserver) GetName() string {
	return co.name
}

// FileObserver 文件观察者
type FileObserver struct {
	name       string
	filePath   string
	logger     *common.Logger
	eventCount int
}

// NewFileObserver 创建文件观察者
func NewFileObserver(name, filePath string) *FileObserver {
	return &FileObserver{
		name:     name,
		filePath: filePath,
		logger:   common.GetLogger(),
	}
}

func (fo *FileObserver) OnDeploymentEvent(event *DeploymentEvent) {
	fo.eventCount++
	
	// 这里可以实现将事件写入文件的逻辑
	fo.logger.Info(fmt.Sprintf("文件观察者 %s 记录事件 #%d: %s", fo.name, fo.eventCount, event.Type))
}

func (fo *FileObserver) GetName() string {
	return fo.name
}

// WebhookObserver Webhook 观察者
type WebhookObserver struct {
	name     string
	webhookURL string
	logger   *common.Logger
}

// NewWebhookObserver 创建 Webhook 观察者
func NewWebhookObserver(name, webhookURL string) *WebhookObserver {
	return &WebhookObserver{
		name:       name,
		webhookURL: webhookURL,
		logger:     common.GetLogger(),
	}
}

func (wo *WebhookObserver) OnDeploymentEvent(event *DeploymentEvent) {
	// 这里可以实现发送 Webhook 的逻辑
	wo.logger.Info(fmt.Sprintf("Webhook 观察者 %s 发送事件到: %s", wo.name, wo.webhookURL))
}

func (wo *WebhookObserver) GetName() string {
	return wo.name
}

// MetricsObserver 指标观察者
type MetricsObserver struct {
	name        string
	metrics     map[string]interface{}
	mutex       sync.RWMutex
	logger      *common.Logger
}

// NewMetricsObserver 创建指标观察者
func NewMetricsObserver(name string) *MetricsObserver {
	return &MetricsObserver{
		name:    name,
		metrics: make(map[string]interface{}),
		logger:  common.GetLogger(),
	}
}

func (mo *MetricsObserver) OnDeploymentEvent(event *DeploymentEvent) {
	mo.mutex.Lock()
	defer mo.mutex.Unlock()
	
	// 更新指标
	eventType := event.Type
	if _, exists := mo.metrics[eventType]; !exists {
		mo.metrics[eventType] = 0
	}
	
	mo.metrics[eventType] = mo.metrics[eventType].(int) + 1
	
	// 记录最后事件时间
	mo.metrics["last_event_time"] = event.Timestamp
	mo.metrics["last_event_type"] = event.Type
	
	mo.logger.Debug(fmt.Sprintf("指标观察者 %s 更新指标: %s", mo.name, eventType))
}

func (mo *MetricsObserver) GetName() string {
	return mo.name
}

// GetMetrics 获取指标
func (mo *MetricsObserver) GetMetrics() map[string]interface{} {
	mo.mutex.RLock()
	defer mo.mutex.RUnlock()
	
	// 复制指标
	metrics := make(map[string]interface{})
	for key, value := range mo.metrics {
		metrics[key] = value
	}
	
	return metrics
}

// DeploymentEventBuilder 部署事件构建器
type DeploymentEventBuilder struct {
	event *DeploymentEvent
}

// NewDeploymentEventBuilder 创建部署事件构建器
func NewDeploymentEventBuilder(eventType string) *DeploymentEventBuilder {
	return &DeploymentEventBuilder{
		event: &DeploymentEvent{
			Type:      eventType,
			Timestamp: time.Now(),
			Data:      make(map[string]interface{}),
		},
	}
}

// WithData 添加数据
func (deb *DeploymentEventBuilder) WithData(key string, value interface{}) *DeploymentEventBuilder {
	deb.event.Data[key] = value
	return deb
}

// WithError 添加错误
func (deb *DeploymentEventBuilder) WithError(err error) *DeploymentEventBuilder {
	deb.event.Error = err
	return deb
}

// Build 构建事件
func (deb *DeploymentEventBuilder) Build() *DeploymentEvent {
	return deb.event
}

// DeploymentEventTypes 部署事件类型常量
const (
	EventDeploymentStarted   = "deployment_started"
	EventDeploymentProgress  = "deployment_progress"
	EventDeploymentCompleted = "deployment_completed"
	EventDeploymentFailed    = "deployment_failed"
	EventFileUploaded        = "file_uploaded"
	EventValidationPassed    = "validation_passed"
	EventValidationFailed    = "validation_failed"
	EventRollbackStarted     = "rollback_started"
	EventRollbackCompleted   = "rollback_completed"
)
