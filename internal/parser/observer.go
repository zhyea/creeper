package parser

import (
	"fmt"
	"sync"
)

// ParseEvent 解析事件类型
type ParseEvent int

const (
	ParseEventStart ParseEvent = iota
	ParseEventProgress
	ParseEventChapterFound
	ParseEventVolumeFound
	ParseEventMetadataFound
	ParseEventComplete
	ParseEventError
)

// ParseEventData 解析事件数据
type ParseEventData struct {
	Event       ParseEvent
	Message     string
	Progress    float64 // 0-100
	ChapterInfo *Chapter
	Error       error
	Metadata    map[string]interface{}
}

// ParseObserver 解析观察者接口
type ParseObserver interface {
	OnParseEvent(data *ParseEventData)
}

// ParseSubject 解析主题接口
type ParseSubject interface {
	Subscribe(observer ParseObserver)
	Unsubscribe(observer ParseObserver)
	NotifyObservers(data *ParseEventData)
}

// ParseNotifier 解析通知器
type ParseNotifier struct {
	observers []ParseObserver
	mutex     sync.RWMutex
}

// NewParseNotifier 创建解析通知器
func NewParseNotifier() *ParseNotifier {
	return &ParseNotifier{
		observers: make([]ParseObserver, 0),
	}
}

// Subscribe 订阅观察者
func (pn *ParseNotifier) Subscribe(observer ParseObserver) {
	pn.mutex.Lock()
	defer pn.mutex.Unlock()
	pn.observers = append(pn.observers, observer)
}

// Unsubscribe 取消订阅观察者
func (pn *ParseNotifier) Unsubscribe(observer ParseObserver) {
	pn.mutex.Lock()
	defer pn.mutex.Unlock()
	
	for i, obs := range pn.observers {
		if obs == observer {
			pn.observers = append(pn.observers[:i], pn.observers[i+1:]...)
			break
		}
	}
}

// NotifyObservers 通知所有观察者
func (pn *ParseNotifier) NotifyObservers(data *ParseEventData) {
	pn.mutex.RLock()
	defer pn.mutex.RUnlock()
	
	for _, observer := range pn.observers {
		go observer.OnParseEvent(data) // 异步通知
	}
}

// ConsoleObserver 控制台观察者
type ConsoleObserver struct {
	verbose bool
}

// NewConsoleObserver 创建控制台观察者
func NewConsoleObserver(verbose bool) *ConsoleObserver {
	return &ConsoleObserver{verbose: verbose}
}

// OnParseEvent 处理解析事件
func (co *ConsoleObserver) OnParseEvent(data *ParseEventData) {
	switch data.Event {
	case ParseEventStart:
		fmt.Printf("🚀 开始解析: %s\n", data.Message)
		
	case ParseEventProgress:
		if co.verbose {
			fmt.Printf("📊 解析进度: %.1f%% - %s\n", data.Progress, data.Message)
		}
		
	case ParseEventChapterFound:
		if co.verbose && data.ChapterInfo != nil {
			fmt.Printf("📖 发现章节: %s\n", data.ChapterInfo.Title)
		}
		
	case ParseEventVolumeFound:
		fmt.Printf("📚 发现卷: %s\n", data.Message)
		
	case ParseEventMetadataFound:
		if co.verbose {
			fmt.Printf("📋 发现元数据: %s\n", data.Message)
		}
		
	case ParseEventComplete:
		fmt.Printf("✅ 解析完成: %s\n", data.Message)
		
	case ParseEventError:
		fmt.Printf("❌ 解析错误: %s\n", data.Message)
		if data.Error != nil && co.verbose {
			fmt.Printf("   详细错误: %v\n", data.Error)
		}
	}
}

// ProgressObserver 进度观察者
type ProgressObserver struct {
	callback func(progress float64, message string)
}

// NewProgressObserver 创建进度观察者
func NewProgressObserver(callback func(float64, string)) *ProgressObserver {
	return &ProgressObserver{callback: callback}
}

// OnParseEvent 处理解析事件
func (po *ProgressObserver) OnParseEvent(data *ParseEventData) {
	if data.Event == ParseEventProgress && po.callback != nil {
		po.callback(data.Progress, data.Message)
	}
}

// StatisticsObserver 统计观察者
type StatisticsObserver struct {
	ChapterCount int
	VolumeCount  int
	WordCount    int
	ErrorCount   int
	mutex        sync.Mutex
}

// NewStatisticsObserver 创建统计观察者
func NewStatisticsObserver() *StatisticsObserver {
	return &StatisticsObserver{}
}

// OnParseEvent 处理解析事件
func (so *StatisticsObserver) OnParseEvent(data *ParseEventData) {
	so.mutex.Lock()
	defer so.mutex.Unlock()
	
	switch data.Event {
	case ParseEventChapterFound:
		so.ChapterCount++
		if data.ChapterInfo != nil {
			so.WordCount += data.ChapterInfo.WordCount
		}
		
	case ParseEventVolumeFound:
		so.VolumeCount++
		
	case ParseEventError:
		so.ErrorCount++
	}
}

// GetStatistics 获取统计信息
func (so *StatisticsObserver) GetStatistics() map[string]int {
	so.mutex.Lock()
	defer so.mutex.Unlock()
	
	return map[string]int{
		"chapters": so.ChapterCount,
		"volumes":  so.VolumeCount,
		"words":    so.WordCount,
		"errors":   so.ErrorCount,
	}
}

// Reset 重置统计
func (so *StatisticsObserver) Reset() {
	so.mutex.Lock()
	defer so.mutex.Unlock()
	
	so.ChapterCount = 0
	so.VolumeCount = 0
	so.WordCount = 0
	so.ErrorCount = 0
}
