package facade

import (
	"fmt"
	"path/filepath"
	"time"

	"creeper/internal/config"
	"creeper/internal/generator"
	"creeper/internal/parser"
	"creeper/internal/common"
	"creeper/internal/deploy"
)

// CreeperFacade Creeper 外观类
type CreeperFacade struct {
	config          *config.Config
	parser          *parser.Parser
	generator       *generator.Generator
	enhancedParser  *parser.EnhancedParser
	deployManager   *deploy.DeployManager
	logger          *common.Logger
	resourceManager *common.ResourceManager
	configCache     *common.ConfigCache
}

// NewCreeperFacade 创建 Creeper 外观
func NewCreeperFacade(configPath string) (*CreeperFacade, error) {
	facade := &CreeperFacade{
		logger:          common.GetLogger(),
		resourceManager: common.GetResourceManager(),
		configCache:     common.GetConfigCache(),
	}
	
	// 初始化配置
	if err := facade.initializeConfig(configPath); err != nil {
		return nil, fmt.Errorf("初始化配置失败: %w", err)
	}
	
	// 初始化组件
	facade.initializeComponents()
	
	facade.logger.Info("Creeper 外观初始化完成")
	
	return facade, nil
}

// initializeConfig 初始化配置
func (cf *CreeperFacade) initializeConfig(configPath string) error {
	// 尝试从缓存加载
	if cached, exists := cf.configCache.Get(configPath); exists {
		if config, ok := cached.(*config.Config); ok {
			cf.config = config
			cf.logger.Info("从缓存加载配置:", configPath)
			return nil
		}
	}
	
	// 加载配置文件
	cfg, err := config.Load(configPath)
	if err != nil {
		cf.logger.Warn("无法读取配置文件，使用默认配置:", err)
		cfg = config.Default()
	}
	
	cf.config = cfg
	
	// 缓存配置
	cf.configCache.Set(configPath, cfg)
	
	return nil
}

// initializeComponents 初始化组件
func (cf *CreeperFacade) initializeComponents() {
	// 创建解析器
	cf.parser = parser.New()
	
	// 创建增强解析器
	consoleObserver := parser.NewConsoleObserver(true)
	cf.enhancedParser = parser.NewEnhancedParser().
		WithLogging(consoleObserver).
		WithCaching().
		WithValidation()
	
	// 创建生成器
	cf.generator = generator.New(cf.config)
	
	// 初始化部署管理器
	if cf.config.Deploy != nil && cf.config.Deploy.Enabled {
		if err := cf.initializeDeployManager(); err != nil {
			cf.logger.Warn("部署管理器初始化失败:", err)
		}
	}
	
	// 设置资源管理器
	cf.resourceManager.Set("config", cf.config)
	cf.resourceManager.Set("parser", cf.parser)
	cf.resourceManager.Set("generator", cf.generator)
	if cf.deployManager != nil {
		cf.resourceManager.Set("deployManager", cf.deployManager)
	}
}

// GenerateWebsite 生成网站（主要功能入口）
func (cf *CreeperFacade) GenerateWebsite() error {
	startTime := time.Now()
	
	cf.logger.Info("开始生成静态网站")
	cf.logger.Info("输入目录:", cf.config.InputDir)
	cf.logger.Info("输出目录:", cf.config.OutputDir)
	
	// 执行生成
	if err := cf.generator.Generate(); err != nil {
		cf.logger.Error("网站生成失败:", err)
		return fmt.Errorf("网站生成失败: %w", err)
	}
	
	duration := time.Since(startTime)
	cf.logger.Info("网站生成完成，耗时:", duration)
	
	return nil
}

// ServeWebsite 启动服务器
func (cf *CreeperFacade) ServeWebsite(port int) error {
	cf.logger.Info("启动本地服务器，端口:", port)
	
	if err := cf.generator.Serve(port); err != nil {
		cf.logger.Error("服务器启动失败:", err)
		return fmt.Errorf("服务器启动失败: %w", err)
	}
	
	return nil
}

// ParseNovel 解析单个小说
func (cf *CreeperFacade) ParseNovel(novelPath string) (*parser.Novel, error) {
	cf.logger.Info("解析小说:", novelPath)
	
	// 使用增强解析器
	decorator := cf.enhancedParser.Build()
	novel, err := decorator.ParseNovel(novelPath)
	
	if err != nil {
		cf.logger.Error("小说解析失败:", err)
		return nil, fmt.Errorf("小说解析失败: %w", err)
	}
	
	cf.logger.Info("小说解析完成:", novel.Title, "共", len(novel.Chapters), "章")
	
	return novel, nil
}

// GetNovelList 获取小说列表
func (cf *CreeperFacade) GetNovelList() ([]*parser.Novel, error) {
	cf.logger.Info("获取小说列表")
	
	// 这里可以扩展为更复杂的小说发现逻辑
	// 目前简单返回生成器中的小说列表
	
	return nil, fmt.Errorf("功能待实现")
}

// UpdateConfig 更新配置
func (cf *CreeperFacade) UpdateConfig(updates map[string]interface{}) error {
	cf.logger.Info("更新配置")
	
	// 使用建造者模式更新配置
	builder := config.Builder()
	
	for key, value := range updates {
		switch key {
		case "site.title":
			if title, ok := value.(string); ok {
				builder.WithSiteTitle(title)
			}
		case "site.description":
			if desc, ok := value.(string); ok {
				builder.WithSiteDescription(desc)
			}
		case "site.author":
			if author, ok := value.(string); ok {
				builder.WithSiteAuthor(author)
			}
		case "theme.primary_color":
			if color, ok := value.(string); ok {
				builder.WithThemePrimaryColor(color)
			}
		case "input_dir":
			if dir, ok := value.(string); ok {
				builder.WithInputDir(dir)
			}
		case "output_dir":
			if dir, ok := value.(string); ok {
				builder.WithOutputDir(dir)
			}
		}
	}
	
	// 应用更新
	updatedConfig := builder.Build()
	cf.config = updatedConfig
	
	// 更新缓存
	cf.configCache.Set("current", updatedConfig)
	
	// 重新初始化生成器
	cf.generator = generator.New(cf.config)
	cf.resourceManager.Set("generator", cf.generator)
	
	cf.logger.Info("配置更新完成")
	
	return nil
}

// GetSystemStatus 获取系统状态
func (cf *CreeperFacade) GetSystemStatus() map[string]interface{} {
	status := map[string]interface{}{
		"config": map[string]interface{}{
			"input_dir":  cf.config.InputDir,
			"output_dir": cf.config.OutputDir,
			"site_title": cf.config.Site.Title,
		},
		"resources": map[string]interface{}{
			"count": cf.resourceManager.Count(),
			"keys":  cf.resourceManager.Keys(),
		},
		"cache": map[string]interface{}{
			"config_cache_exists": cf.configCache.Exists("current"),
		},
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	
	return status
}

// ValidateSetup 验证系统设置
func (cf *CreeperFacade) ValidateSetup() error {
	cf.logger.Info("验证系统设置")
	
	// 验证输入目录
	if cf.config.InputDir == "" {
		return fmt.Errorf("输入目录未设置")
	}
	
	// 验证输出目录
	if cf.config.OutputDir == "" {
		return fmt.Errorf("输出目录未设置")
	}
	
	// 验证主题配置
	if cf.config.Theme.PrimaryColor == "" {
		return fmt.Errorf("主题主色调未设置")
	}
	
	cf.logger.Info("系统设置验证通过")
	
	return nil
}

// GetSupportedFormats 获取支持的文件格式
func (cf *CreeperFacade) GetSupportedFormats() []string {
	return []string{
		"Markdown (.md)",
		"Text (.txt)",
		"Markdown Directory",
		"Text Directory",
		"Multi-Volume Structure",
	}
}

// CleanCache 清理缓存
func (cf *CreeperFacade) CleanCache() {
	cf.logger.Info("清理系统缓存")
	
	cf.configCache.Clear()
	cf.resourceManager.Clear()
	
	cf.logger.Info("缓存清理完成")
}

// GetVersion 获取版本信息
func (cf *CreeperFacade) GetVersion() map[string]string {
	return map[string]string{
		"version":     "1.0.0",
		"build_date":  "2024-01-01",
		"go_version":  "1.21+",
		"description": "Creeper 静态小说站点生成器",
	}
}

// Shutdown 优雅关闭
func (cf *CreeperFacade) Shutdown() error {
	cf.logger.Info("开始优雅关闭系统")
	
	// 清理资源
	cf.CleanCache()
	
	// 保存重要状态
	if err := cf.saveSystemState(); err != nil {
		cf.logger.Warn("保存系统状态失败:", err)
	}
	
	cf.logger.Info("系统关闭完成")
	
	return nil
}

// initializeDeployManager 初始化部署管理器
func (cf *CreeperFacade) initializeDeployManager() error {
	cf.logger.Info("初始化部署管理器")
	
	// 加载部署配置
	deployConfig, err := deploy.LoadDeployConfig(cf.config.Deploy.Config)
	if err != nil {
		return fmt.Errorf("加载部署配置失败: %w", err)
	}
	
	// 创建部署管理器
	cf.deployManager = deploy.NewDeployManager(deployConfig)
	
	// 初始化部署管理器
	if err := cf.deployManager.Initialize(); err != nil {
		return fmt.Errorf("初始化部署管理器失败: %w", err)
	}
	
	cf.logger.Info("部署管理器初始化完成")
	return nil
}

// DeployWebsite 部署网站
func (cf *CreeperFacade) DeployWebsite() error {
	if cf.deployManager == nil {
		return fmt.Errorf("部署管理器未初始化，请检查部署配置")
	}
	
	cf.logger.Info("开始部署网站")
	
	// 执行部署
	if err := cf.deployManager.Deploy(cf.config.OutputDir); err != nil {
		cf.logger.Error("网站部署失败:", err)
		return fmt.Errorf("网站部署失败: %w", err)
	}
	
	deploymentURL := cf.deployManager.GetDeploymentURL()
	cf.logger.Info("网站部署完成，访问地址:", deploymentURL)
	
	return nil
}

// GetDeploymentStatus 获取部署状态
func (cf *CreeperFacade) GetDeploymentStatus() (map[string]interface{}, error) {
	if cf.deployManager == nil {
		return nil, fmt.Errorf("部署管理器未初始化")
	}
	
	return cf.deployManager.GetStatus()
}

// GetDeploymentURL 获取部署 URL
func (cf *CreeperFacade) GetDeploymentURL() string {
	if cf.deployManager == nil {
		return ""
	}
	return cf.deployManager.GetDeploymentURL()
}

// saveSystemState 保存系统状态
func (cf *CreeperFacade) saveSystemState() error {
	// 这里可以保存重要的系统状态
	// 比如最后的配置、缓存信息等
	
	stateFile := filepath.Join(cf.config.OutputDir, ".creeper_state")
	status := cf.GetSystemStatus()
	
	// 简单的状态保存（实际项目中可能需要更复杂的序列化）
	cf.resourceManager.Set("last_state", status)
	
	return nil
}
