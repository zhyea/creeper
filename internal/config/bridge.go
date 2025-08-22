package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"creeper/internal/common"
	"gopkg.in/yaml.v3"
)

// ConfigStorage 配置存储接口
type ConfigStorage interface {
	Load(path string) (*Config, error)
	Save(config *Config, path string) error
	Validate(config *Config) error
	GetFormat() string
}

// ConfigValidator 配置验证接口
type ConfigValidator interface {
	Validate(config *Config) error
	GetValidationRules() map[string]interface{}
}

// ConfigBridge 配置桥接器
type ConfigBridge struct {
	storage   ConfigStorage
	validator ConfigValidator
	logger    *common.Logger
}

// NewConfigBridge 创建配置桥接器
func NewConfigBridge(storage ConfigStorage, validator ConfigValidator) *ConfigBridge {
	return &ConfigBridge{
		storage:   storage,
		validator: validator,
		logger:    common.GetLogger(),
	}
}

// LoadConfig 加载配置
func (cb *ConfigBridge) LoadConfig(path string) (*Config, error) {
	cb.logger.Info("加载配置文件:", path)
	
	// 使用存储接口加载配置
	config, err := cb.storage.Load(path)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	
	// 使用验证接口验证配置
	if err := cb.validator.Validate(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}
	
	cb.logger.Info("配置加载成功，格式:", cb.storage.GetFormat())
	return config, nil
}

// SaveConfig 保存配置
func (cb *ConfigBridge) SaveConfig(config *Config, path string) error {
	cb.logger.Info("保存配置文件:", path)
	
	// 使用验证接口验证配置
	if err := cb.validator.Validate(config); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}
	
	// 使用存储接口保存配置
	if err := cb.storage.Save(config, path); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}
	
	cb.logger.Info("配置保存成功")
	return nil
}

// YAMLConfigStorage YAML 配置存储实现
type YAMLConfigStorage struct {
	logger *common.Logger
}

// NewYAMLConfigStorage 创建 YAML 配置存储
func NewYAMLConfigStorage() *YAMLConfigStorage {
	return &YAMLConfigStorage{
		logger: common.GetLogger(),
	}
}

func (ycs *YAMLConfigStorage) Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (ycs *YAMLConfigStorage) Save(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (ycs *YAMLConfigStorage) Validate(config *Config) error {
	// 基础验证逻辑
	if config.Site.Title == "" {
		return fmt.Errorf("站点标题不能为空")
	}
	if config.InputDir == "" {
		return fmt.Errorf("输入目录不能为空")
	}
	if config.OutputDir == "" {
		return fmt.Errorf("输出目录不能为空")
	}
	return nil
}

func (ycs *YAMLConfigStorage) GetFormat() string {
	return "YAML"
}

// JSONConfigStorage JSON 配置存储实现
type JSONConfigStorage struct {
	logger *common.Logger
}

// NewJSONConfigStorage 创建 JSON 配置存储
func NewJSONConfigStorage() *JSONConfigStorage {
	return &JSONConfigStorage{
		logger: common.GetLogger(),
	}
}

func (jcs *JSONConfigStorage) Load(path string) (*Config, error) {
	// JSON 加载实现（这里简化处理）
	return nil, fmt.Errorf("JSON 配置存储暂未实现")
}

func (jcs *JSONConfigStorage) Save(config *Config, path string) error {
	// JSON 保存实现（这里简化处理）
	return fmt.Errorf("JSON 配置存储暂未实现")
}

func (jcs *JSONConfigStorage) Validate(config *Config) error {
	// 基础验证逻辑
	if config.Site.Title == "" {
		return fmt.Errorf("站点标题不能为空")
	}
	return nil
}

func (jcs *JSONConfigStorage) GetFormat() string {
	return "JSON"
}

// DefaultConfigValidator 默认配置验证器
type DefaultConfigValidator struct {
	logger *common.Logger
}

// NewDefaultConfigValidator 创建默认配置验证器
func NewDefaultConfigValidator() *DefaultConfigValidator {
	return &DefaultConfigValidator{
		logger: common.GetLogger(),
	}
}

func (dcv *DefaultConfigValidator) Validate(config *Config) error {
	// 验证站点配置
	if err := dcv.validateSiteConfig(config.Site); err != nil {
		return fmt.Errorf("站点配置验证失败: %w", err)
	}
	
	// 验证主题配置
	if err := dcv.validateThemeConfig(config.Theme); err != nil {
		return fmt.Errorf("主题配置验证失败: %w", err)
	}
	
	// 验证构建配置
	if err := dcv.validateBuildConfig(config.Build); err != nil {
		return fmt.Errorf("构建配置验证失败: %w", err)
	}
	
	// 验证部署配置
	if config.Deploy != nil {
		if err := dcv.validateDeployConfig(config.Deploy); err != nil {
			return fmt.Errorf("部署配置验证失败: %w", err)
		}
	}
	
	return nil
}

func (dcv *DefaultConfigValidator) GetValidationRules() map[string]interface{} {
	return map[string]interface{}{
		"site": map[string]interface{}{
			"title":       "required,string",
			"description": "optional,string",
			"author":      "optional,string",
			"base_url":    "optional,string",
		},
		"theme": map[string]interface{}{
			"name":             "optional,string",
			"primary_color":    "required,string",
			"secondary_color":  "optional,string",
			"background_color": "optional,string",
			"text_color":       "optional,string",
		},
		"build": map[string]interface{}{
			"minify_html": "optional,bool",
			"minify_css":  "optional,bool",
			"minify_js":   "optional,bool",
		},
		"deploy": map[string]interface{}{
			"enabled": "optional,bool",
			"type":    "optional,string",
			"config":  "optional,string",
		},
	}
}

// validateSiteConfig 验证站点配置
func (dcv *DefaultConfigValidator) validateSiteConfig(site SiteConfig) error {
	if site.Title == "" {
		return fmt.Errorf("站点标题不能为空")
	}
	
	if site.BaseURL != "" && !strings.HasPrefix(site.BaseURL, "/") {
		return fmt.Errorf("站点基础URL必须以/开头")
	}
	
	return nil
}

// validateThemeConfig 验证主题配置
func (dcv *DefaultConfigValidator) validateThemeConfig(theme ThemeConfig) error {
	if theme.PrimaryColor == "" {
		return fmt.Errorf("主题主色调不能为空")
	}
	
	// 验证颜色格式（简化验证）
	if !strings.HasPrefix(theme.PrimaryColor, "#") {
		return fmt.Errorf("主色调必须是有效的十六进制颜色值")
	}
	
	return nil
}

// validateBuildConfig 验证构建配置
func (dcv *DefaultConfigValidator) validateBuildConfig(build BuildConfig) error {
	// 构建配置验证逻辑
	return nil
}

// validateDeployConfig 验证部署配置
func (dcv *DefaultConfigValidator) validateDeployConfig(deploy *DeployConfig) error {
	if deploy.Enabled && deploy.Type == "" {
		return fmt.Errorf("启用部署时，部署类型不能为空")
	}
	
	if deploy.Enabled && deploy.Config == "" {
		return fmt.Errorf("启用部署时，部署配置文件路径不能为空")
	}
	
	return nil
}

// ConfigFactory 配置工厂
type ConfigFactory struct{}

// CreateStorage 创建配置存储
func (cf *ConfigFactory) CreateStorage(format string) (ConfigStorage, error) {
	switch strings.ToLower(format) {
	case "yaml", "yml":
		return NewYAMLConfigStorage(), nil
	case "json":
		return NewJSONConfigStorage(), nil
	default:
		return nil, fmt.Errorf("不支持的配置格式: %s", format)
	}
}

// CreateValidator 创建配置验证器
func (cf *ConfigFactory) CreateValidator(validatorType string) (ConfigValidator, error) {
	switch validatorType {
	case "default":
		return NewDefaultConfigValidator(), nil
	default:
		return nil, fmt.Errorf("不支持的验证器类型: %s", validatorType)
	}
}

// CreateBridge 创建配置桥接器
func (cf *ConfigFactory) CreateBridge(format, validatorType string) (*ConfigBridge, error) {
	storage, err := cf.CreateStorage(format)
	if err != nil {
		return nil, err
	}
	
	validator, err := cf.CreateValidator(validatorType)
	if err != nil {
		return nil, err
	}
	
	return NewConfigBridge(storage, validator), nil
}
