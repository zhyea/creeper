package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"creeper/internal/common"
)

// DeploymentMemento 部署备忘录
type DeploymentMemento struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Status        string                 `json:"status"`
	URL           string                 `json:"url"`
	SiteDir       string                 `json:"site_dir"`
	FileCount     int                    `json:"file_count"`
	TotalSize     int64                  `json:"total_size"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Duration      time.Duration          `json:"duration"`
	Metadata      map[string]interface{} `json:"metadata"`
	Error         string                 `json:"error,omitempty"`
	RollbackPoint string                 `json:"rollback_point,omitempty"`
}

// DeploymentCaretaker 部署状态管理者
type DeploymentCaretaker struct {
	mementos    []*DeploymentMemento
	maxHistory  int
	storagePath string
	logger      *common.Logger
}

// NewDeploymentCaretaker 创建部署状态管理者
func NewDeploymentCaretaker(storagePath string, maxHistory int) *DeploymentCaretaker {
	return &DeploymentCaretaker{
		mementos:    make([]*DeploymentMemento, 0),
		maxHistory:  maxHistory,
		storagePath: storagePath,
		logger:      common.GetLogger(),
	}
}

// SaveMemento 保存备忘录
func (dc *DeploymentCaretaker) SaveMemento(memento *DeploymentMemento) error {
	dc.mementos = append(dc.mementos, memento)
	
	// 限制历史记录数量
	if len(dc.mementos) > dc.maxHistory {
		dc.mementos = dc.mementos[len(dc.mementos)-dc.maxHistory:]
	}
	
	// 保存到文件
	return dc.saveToFile()
}

// GetMemento 获取备忘录
func (dc *DeploymentCaretaker) GetMemento(id string) (*DeploymentMemento, error) {
	for _, memento := range dc.mementos {
		if memento.ID == id {
			return memento, nil
		}
	}
	return nil, fmt.Errorf("备忘录不存在: %s", id)
}

// GetLatestMemento 获取最新的备忘录
func (dc *DeploymentCaretaker) GetLatestMemento() (*DeploymentMemento, error) {
	if len(dc.mementos) == 0 {
		return nil, fmt.Errorf("没有部署历史")
	}
	return dc.mementos[len(dc.mementos)-1], nil
}

// GetAllMementos 获取所有备忘录
func (dc *DeploymentCaretaker) GetAllMementos() []*DeploymentMemento {
	return dc.mementos
}

// RestoreMemento 恢复备忘录
func (dc *DeploymentCaretaker) RestoreMemento(id string) (*DeploymentMemento, error) {
	memento, err := dc.GetMemento(id)
	if err != nil {
		return nil, err
	}
	
	dc.logger.Info("恢复部署状态:", id)
	return memento, nil
}

// saveToFile 保存到文件
func (dc *DeploymentCaretaker) saveToFile() error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(dc.storagePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	
	// 序列化数据
	data, err := json.MarshalIndent(dc.mementos, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(dc.storagePath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	
	return nil
}

// loadFromFile 从文件加载
func (dc *DeploymentCaretaker) loadFromFile() error {
	if _, err := os.Stat(dc.storagePath); os.IsNotExist(err) {
		// 文件不存在，使用空列表
		return nil
	}
	
	data, err := os.ReadFile(dc.storagePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	
	if err := json.Unmarshal(data, &dc.mementos); err != nil {
		return fmt.Errorf("反序列化失败: %w", err)
	}
	
	return nil
}

// LoadHistory 加载历史记录
func (dc *DeploymentCaretaker) LoadHistory() error {
	return dc.loadFromFile()
}

// ClearHistory 清空历史记录
func (dc *DeploymentCaretaker) ClearHistory() error {
	dc.mementos = make([]*DeploymentMemento, 0)
	return dc.saveToFile()
}

// GetDeploymentStats 获取部署统计
func (dc *DeploymentCaretaker) GetDeploymentStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_deployments": len(dc.mementos),
		"successful":        0,
		"failed":           0,
		"total_duration":   time.Duration(0),
		"total_files":      0,
		"total_size":       int64(0),
	}
	
	for _, memento := range dc.mementos {
		if memento.Status == "success" {
			stats["successful"] = stats["successful"].(int) + 1
		} else {
			stats["failed"] = stats["failed"].(int) + 1
		}
		
		stats["total_duration"] = stats["total_duration"].(time.Duration) + memento.Duration
		stats["total_files"] = stats["total_files"].(int) + memento.FileCount
		stats["total_size"] = stats["total_size"].(int64) + memento.TotalSize
	}
	
	return stats
}

// DeploymentOriginator 部署发起者
type DeploymentOriginator struct {
	caretaker *DeploymentCaretaker
	logger    *common.Logger
}

// NewDeploymentOriginator 创建部署发起者
func NewDeploymentOriginator(caretaker *DeploymentCaretaker) *DeploymentOriginator {
	return &DeploymentOriginator{
		caretaker: caretaker,
		logger:    common.GetLogger(),
	}
}

// CreateMemento 创建备忘录
func (do *DeploymentOriginator) CreateMemento(deploymentType, siteDir string) *DeploymentMemento {
	return &DeploymentMemento{
		ID:        fmt.Sprintf("deploy_%d", time.Now().Unix()),
		Type:      deploymentType,
		Status:    "pending",
		SiteDir:   siteDir,
		StartTime: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// UpdateMemento 更新备忘录
func (do *DeploymentOriginator) UpdateMemento(memento *DeploymentMemento, status string, metadata map[string]interface{}) {
	memento.Status = status
	memento.EndTime = time.Now()
	memento.Duration = memento.EndTime.Sub(memento.StartTime)
	
	if metadata != nil {
		for key, value := range metadata {
			memento.Metadata[key] = value
		}
	}
	
	// 保存更新
	if err := do.caretaker.SaveMemento(memento); err != nil {
		do.logger.Warn("保存部署状态失败:", err)
	}
}

// SetError 设置错误信息
func (do *DeploymentOriginator) SetError(memento *DeploymentMemento, err error) {
	memento.Status = "failed"
	memento.Error = err.Error()
	memento.EndTime = time.Now()
	memento.Duration = memento.EndTime.Sub(memento.StartTime)
	
	// 保存更新
	if err := do.caretaker.SaveMemento(memento); err != nil {
		do.logger.Warn("保存部署状态失败:", err)
	}
}

// SetSuccess 设置成功状态
func (do *DeploymentOriginator) SetSuccess(memento *DeploymentMemento, url string, fileCount int, totalSize int64) {
	memento.Status = "success"
	memento.URL = url
	memento.FileCount = fileCount
	memento.TotalSize = totalSize
	memento.EndTime = time.Now()
	memento.Duration = memento.EndTime.Sub(memento.StartTime)
	
	// 保存更新
	if err := do.caretaker.SaveMemento(memento); err != nil {
		do.logger.Warn("保存部署状态失败:", err)
	}
}
