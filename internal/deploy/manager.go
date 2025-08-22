package deploy

import (
	"fmt"
	"os"
	"path/filepath"

	"creeper/internal/common"
	"gopkg.in/yaml.v3"
)

// DeployType 部署类型
type DeployType string

const (
	CloudflarePages DeployType = "cloudflare"
	GitHubPages     DeployType = "github"
	Vercel          DeployType = "vercel"
	Netlify         DeployType = "netlify"
)

// DeployConfig 部署配置
type DeployConfig struct {
	Type       DeployType                `yaml:"type"`
	Cloudflare *CloudflareConfig         `yaml:"cloudflare,omitempty"`
	GitHub     *GitHubPagesConfig        `yaml:"github,omitempty"`
	Vercel     *VercelConfig             `yaml:"vercel,omitempty"`
	Netlify    *NetlifyConfig            `yaml:"netlify,omitempty"`
	Options    map[string]interface{}    `yaml:"options,omitempty"`
}

// GitHubPagesConfig GitHub Pages 配置
type GitHubPagesConfig struct {
	Repository string `yaml:"repository"`
	Branch     string `yaml:"branch"`
	Token      string `yaml:"token"`
	Username   string `yaml:"username"`
}

// VercelConfig Vercel 配置
type VercelConfig struct {
	ProjectID   string `yaml:"project_id"`
	Token       string `yaml:"token"`
	TeamID      string `yaml:"team_id,omitempty"`
	Framework   string `yaml:"framework"`
	BuildCommand string `yaml:"build_command"`
}

// NetlifyConfig Netlify 配置
type NetlifyConfig struct {
	SiteID      string `yaml:"site_id"`
	Token       string `yaml:"token"`
	BuildCommand string `yaml:"build_command"`
	PublishDir  string `yaml:"publish_dir"`
}

// Deployer 部署器接口
type Deployer interface {
	Deploy(siteDir string) error
	GetStatus() (map[string]interface{}, error)
	GetDeploymentURL() string
}

// DeployManager 部署管理器
type DeployManager struct {
	config           *DeployConfig
	logger           *common.Logger
	deployer         Deployer
	template         *BaseDeployTemplate
	eventManager     *DeploymentEventManager
	caretaker        *DeploymentCaretaker
	originator       *DeploymentOriginator
	fileIterator     FileIterator
}

// NewDeployManager 创建部署管理器
func NewDeployManager(config *DeployConfig) *DeployManager {
	// 创建事件管理器
	eventManager := NewDeploymentEventManager()
	
	// 创建状态管理者
	caretaker := NewDeploymentCaretaker(".creeper/deployments.json", 100)
	
	// 创建发起者
	originator := NewDeploymentOriginator(caretaker)
	
	// 创建模板
	template := NewBaseDeployTemplate()
	
	return &DeployManager{
		config:       config,
		logger:       common.GetLogger(),
		template:     template,
		eventManager: eventManager,
		caretaker:    caretaker,
		originator:   originator,
	}
}

// Initialize 初始化部署器
func (dm *DeployManager) Initialize() error {
	dm.logger.Info("初始化部署管理器，类型:", dm.config.Type)

	// 加载部署历史
	if err := dm.caretaker.LoadHistory(); err != nil {
		dm.logger.Warn("加载部署历史失败:", err)
	}

	// 添加默认观察者
	dm.eventManager.Attach(NewConsoleObserver("console"))
	dm.eventManager.Attach(NewMetricsObserver("metrics"))

	switch dm.config.Type {
	case CloudflarePages:
		if dm.config.Cloudflare == nil {
			return fmt.Errorf("Cloudflare 配置不能为空")
		}
		dm.deployer = NewCloudflareDeployer(dm.config.Cloudflare)

	case GitHubPages:
		if dm.config.GitHub == nil {
			return fmt.Errorf("GitHub Pages 配置不能为空")
		}
		dm.deployer = NewGitHubPagesDeployer(dm.config.GitHub)

	case Vercel:
		if dm.config.Vercel == nil {
			return fmt.Errorf("Vercel 配置不能为空")
		}
		dm.deployer = NewVercelDeployer(dm.config.Vercel)

	case Netlify:
		if dm.config.Netlify == nil {
			return fmt.Errorf("Netlify 配置不能为空")
		}
		dm.deployer = NewNetlifyDeployer(dm.config.Netlify)

	default:
		return fmt.Errorf("不支持的部署类型: %s", dm.config.Type)
	}

	dm.logger.Info("部署管理器初始化完成")
	return nil
}

// Deploy 执行部署
func (dm *DeployManager) Deploy(siteDir string) error {
	if dm.deployer == nil {
		return fmt.Errorf("部署器未初始化")
	}

	dm.logger.Info("开始部署到", dm.config.Type)
	dm.logger.Info("站点目录:", siteDir)

	// 创建部署备忘录
	memento := dm.originator.CreateMemento(string(dm.config.Type), siteDir)

	// 发送部署开始事件
	dm.eventManager.Notify(NewDeploymentEventBuilder(EventDeploymentStarted).
		WithData("site_dir", siteDir).
		WithData("deployment_type", dm.config.Type).
		Build())

	// 初始化文件迭代器
	dm.fileIterator = NewDirectoryIterator(siteDir)
	if err := dm.fileIterator.(*DirectoryIterator).LoadFiles(); err != nil {
		dm.originator.SetError(memento, err)
		return fmt.Errorf("加载文件列表失败: %w", err)
	}

	// 验证站点目录
	if err := dm.validateSiteDir(siteDir); err != nil {
		dm.originator.SetError(memento, err)
		return fmt.Errorf("站点目录验证失败: %w", err)
	}

	// 使用模板方法执行部署
	if err := dm.template.DeployWithValidation(dm.deployer, siteDir); err != nil {
		dm.originator.SetError(memento, err)
		
		// 发送部署失败事件
		dm.eventManager.Notify(NewDeploymentEventBuilder(EventDeploymentFailed).
			WithData("site_dir", siteDir).
			WithError(err).
			Build())
		
		return fmt.Errorf("部署失败: %w", err)
	}

	// 获取部署 URL
	deploymentURL := dm.deployer.GetDeploymentURL()
	
	// 获取文件统计
	fileCollection := NewFileCollection(dm.fileIterator)
	stats := fileCollection.GetStats()
	
	// 设置成功状态
	dm.originator.SetSuccess(memento, deploymentURL, stats["files"].(int), stats["total_size"].(int64))
	
	// 发送部署完成事件
	dm.eventManager.Notify(NewDeploymentEventBuilder(EventDeploymentCompleted).
		WithData("deployment_url", deploymentURL).
		WithData("file_count", stats["files"]).
		WithData("total_size", stats["total_size"]).
		Build())

	dm.logger.Info("部署完成，访问地址:", deploymentURL)
	return nil
}

// GetStatus 获取部署状态
func (dm *DeployManager) GetStatus() (map[string]interface{}, error) {
	if dm.deployer == nil {
		return nil, fmt.Errorf("部署器未初始化")
	}

	return dm.deployer.GetStatus()
}

// GetDeploymentURL 获取部署 URL
func (dm *DeployManager) GetDeploymentURL() string {
	if dm.deployer == nil {
		return ""
	}
	return dm.deployer.GetDeploymentURL()
}

// GetDeploymentHistory 获取部署历史
func (dm *DeployManager) GetDeploymentHistory() []*DeploymentMemento {
	return dm.caretaker.GetAllMementos()
}

// GetDeploymentStats 获取部署统计
func (dm *DeployManager) GetDeploymentStats() map[string]interface{} {
	return dm.caretaker.GetDeploymentStats()
}

// GetEventMetrics 获取事件指标
func (dm *DeployManager) GetEventMetrics() map[string]interface{} {
	if metricsObserver, ok := dm.eventManager.observers["metrics"]; ok {
		if mo, ok := metricsObserver.(*MetricsObserver); ok {
			return mo.GetMetrics()
		}
	}
	return make(map[string]interface{})
}

// AddObserver 添加观察者
func (dm *DeployManager) AddObserver(observer DeploymentObserver) {
	dm.eventManager.Attach(observer)
}

// RemoveObserver 移除观察者
func (dm *DeployManager) RemoveObserver(observer DeploymentObserver) {
	dm.eventManager.Detach(observer)
}

// GetObserverCount 获取观察者数量
func (dm *DeployManager) GetObserverCount() int {
	return dm.eventManager.GetObserverCount()
}

// DeployWithRetry 带重试的部署
func (dm *DeployManager) DeployWithRetry(siteDir string, maxRetries int) error {
	if dm.deployer == nil {
		return fmt.Errorf("部署器未初始化")
	}

	return dm.template.DeployWithRetry(dm.deployer, siteDir, maxRetries)
}

// DeployWithRollback 带回滚的部署
func (dm *DeployManager) DeployWithRollback(siteDir string) error {
	if dm.deployer == nil {
		return fmt.Errorf("部署器未初始化")
	}

	return dm.template.DeployWithRollback(dm.deployer, siteDir)
}

// validateSiteDir 验证站点目录
func (dm *DeployManager) validateSiteDir(siteDir string) error {
	info, err := os.Stat(siteDir)
	if err != nil {
		return fmt.Errorf("站点目录不存在: %s", siteDir)
	}
	if !info.IsDir() {
		return fmt.Errorf("站点路径不是目录: %s", siteDir)
	}

	// 检查是否有 index.html
	indexPath := filepath.Join(siteDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		return fmt.Errorf("站点目录缺少 index.html: %s", indexPath)
	}

	return nil
}

// LoadDeployConfig 从文件加载部署配置
func LoadDeployConfig(path string) (*DeployConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取部署配置文件失败: %w", err)
	}

	var config DeployConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析部署配置文件失败: %w", err)
	}

	return &config, nil
}

// SaveDeployConfig 保存部署配置到文件
func SaveDeployConfig(config *DeployConfig, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化部署配置失败: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入部署配置文件失败: %w", err)
	}

	return nil
}

// CreateDefaultDeployConfig 创建默认部署配置
func CreateDefaultDeployConfig(deployType DeployType) *DeployConfig {
	config := &DeployConfig{
		Type: deployType,
		Options: map[string]interface{}{
			"auto_deploy": true,
			"preview":     false,
		},
	}

	switch deployType {
	case CloudflarePages:
		config.Cloudflare = &CloudflareConfig{
			Branch:    "main",
			Framework: "none",
			OutputDir: ".",
		}

	case GitHubPages:
		config.GitHub = &GitHubPagesConfig{
			Branch: "gh-pages",
		}

	case Vercel:
		config.Vercel = &VercelConfig{
			Framework:    "other",
			BuildCommand: "",
		}

	case Netlify:
		config.Netlify = &NetlifyConfig{
			BuildCommand: "",
			PublishDir:   ".",
		}
	}

	return config
}

// GitHubPagesDeployer GitHub Pages 部署器
type GitHubPagesDeployer struct {
	config *GitHubPagesConfig
	logger *common.Logger
}

// NewGitHubPagesDeployer 创建 GitHub Pages 部署器
func NewGitHubPagesDeployer(config *GitHubPagesConfig) *GitHubPagesDeployer {
	return &GitHubPagesDeployer{
		config: config,
		logger: common.GetLogger(),
	}
}

func (gpd *GitHubPagesDeployer) Deploy(siteDir string) error {
	gpd.logger.Info("GitHub Pages 部署功能待实现")
	return fmt.Errorf("GitHub Pages 部署功能暂未实现")
}

func (gpd *GitHubPagesDeployer) GetStatus() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type": "github_pages",
		"status": "not_implemented",
	}, nil
}

func (gpd *GitHubPagesDeployer) GetDeploymentURL() string {
	return fmt.Sprintf("https://%s.github.io/%s", gpd.config.Username, gpd.config.Repository)
}

// VercelDeployer Vercel 部署器
type VercelDeployer struct {
	config *VercelConfig
	logger *common.Logger
}

// NewVercelDeployer 创建 Vercel 部署器
func NewVercelDeployer(config *VercelConfig) *VercelDeployer {
	return &VercelDeployer{
		config: config,
		logger: common.GetLogger(),
	}
}

func (vd *VercelDeployer) Deploy(siteDir string) error {
	vd.logger.Info("Vercel 部署功能待实现")
	return fmt.Errorf("Vercel 部署功能暂未实现")
}

func (vd *VercelDeployer) GetStatus() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type": "vercel",
		"status": "not_implemented",
	}, nil
}

func (vd *VercelDeployer) GetDeploymentURL() string {
	return "https://vercel.com"
}

// NetlifyDeployer Netlify 部署器
type NetlifyDeployer struct {
	config *NetlifyConfig
	logger *common.Logger
}

// NewNetlifyDeployer 创建 Netlify 部署器
func NewNetlifyDeployer(config *NetlifyConfig) *NetlifyDeployer {
	return &NetlifyDeployer{
		config: config,
		logger: common.GetLogger(),
	}
}

func (nd *NetlifyDeployer) Deploy(siteDir string) error {
	nd.logger.Info("Netlify 部署功能待实现")
	return fmt.Errorf("Netlify 部署功能暂未实现")
}

func (nd *NetlifyDeployer) GetStatus() (map[string]interface{}, error) {
	return map[string]interface{}{
		"type": "netlify",
		"status": "not_implemented",
	}, nil
}

func (nd *NetlifyDeployer) GetDeploymentURL() string {
	return "https://netlify.com"
}
