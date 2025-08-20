package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config 配置结构
type Config struct {
	// 站点基本信息
	Site SiteConfig `yaml:"site"`

	// 目录配置
	InputDir  string `yaml:"input_dir"`
	OutputDir string `yaml:"output_dir"`

	// 主题配置
	Theme ThemeConfig `yaml:"theme"`

	// 构建配置
	Build BuildConfig `yaml:"build"`
}

// SiteConfig 站点配置
type SiteConfig struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	BaseURL     string `yaml:"base_url"`
}

// ThemeConfig 主题配置
type ThemeConfig struct {
	Name            string `yaml:"name"`
	PrimaryColor    string `yaml:"primary_color"`
	SecondaryColor  string `yaml:"secondary_color"`
	BackgroundColor string `yaml:"background_color"`
	TextColor       string `yaml:"text_color"`
	FontFamily      string `yaml:"font_family"`
	FontSize        string `yaml:"font_size"`
	LineHeight      string `yaml:"line_height"`
}

// BuildConfig 构建配置
type BuildConfig struct {
	MinifyHTML bool `yaml:"minify_html"`
	MinifyCSS  bool `yaml:"minify_css"`
	MinifyJS   bool `yaml:"minify_js"`
}

// Default 返回默认配置
func Default() *Config {
	return &Config{
		Site: SiteConfig{
			Title:       "我的小说站点",
			Description: "静态小说阅读站点",
			Author:      "作者",
			BaseURL:     "/",
		},
		InputDir:  "novels",
		OutputDir: "dist",
		Theme: ThemeConfig{
			Name:            "default",
			PrimaryColor:    "#2c3e50",
			SecondaryColor:  "#3498db",
			BackgroundColor: "#ffffff",
			TextColor:       "#333333",
			FontFamily:      "'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif",
			FontSize:        "16px",
			LineHeight:      "1.6",
		},
		Build: BuildConfig{
			MinifyHTML: true,
			MinifyCSS:  true,
			MinifyJS:   true,
		},
	}
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
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

// Save 保存配置到文件
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
