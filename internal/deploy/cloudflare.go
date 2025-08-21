package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"creeper/internal/common"
)

// CloudflareConfig Cloudflare 配置
type CloudflareConfig struct {
	AccountID    string `yaml:"account_id"`
	ProjectName  string `yaml:"project_name"`
	APIKey       string `yaml:"api_key"`
	Email        string `yaml:"email"`
	Branch       string `yaml:"branch"`
	Framework    string `yaml:"framework"`
	BuildCommand string `yaml:"build_command"`
	OutputDir    string `yaml:"output_dir"`
}

// CloudflareDeployer Cloudflare 部署器
type CloudflareDeployer struct {
	config *CloudflareConfig
	logger *common.Logger
	client *http.Client
}

// NewCloudflareDeployer 创建 Cloudflare 部署器
func NewCloudflareDeployer(config *CloudflareConfig) *CloudflareDeployer {
	return &CloudflareDeployer{
		config: config,
		logger: common.GetLogger(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Deploy 部署到 Cloudflare Pages
func (cd *CloudflareDeployer) Deploy(siteDir string) error {
	cd.logger.Info("开始部署到 Cloudflare Pages")
	cd.logger.Info("项目名称:", cd.config.ProjectName)
	cd.logger.Info("站点目录:", siteDir)

	// 1. 验证配置
	if err := cd.validateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// 2. 检查站点目录
	if err := cd.validateSiteDir(siteDir); err != nil {
		return fmt.Errorf("站点目录验证失败: %w", err)
	}

	// 3. 创建部署
	deploymentID, err := cd.createDeployment(siteDir)
	if err != nil {
		return fmt.Errorf("创建部署失败: %w", err)
	}

	// 4. 上传文件
	if err := cd.uploadFiles(deploymentID, siteDir); err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}

	// 5. 完成部署
	if err := cd.finalizeDeployment(deploymentID); err != nil {
		return fmt.Errorf("完成部署失败: %w", err)
	}

	cd.logger.Info("Cloudflare Pages 部署完成")
	return nil
}

// validateConfig 验证配置
func (cd *CloudflareDeployer) validateConfig() error {
	if cd.config.AccountID == "" {
		return fmt.Errorf("Account ID 不能为空")
	}
	if cd.config.ProjectName == "" {
		return fmt.Errorf("项目名称不能为空")
	}
	if cd.config.APIKey == "" {
		return fmt.Errorf("API Key 不能为空")
	}
	if cd.config.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	// 设置默认值
	if cd.config.Branch == "" {
		cd.config.Branch = "main"
	}
	if cd.config.Framework == "" {
		cd.config.Framework = "none"
	}
	if cd.config.OutputDir == "" {
		cd.config.OutputDir = "."
	}

	return nil
}

// validateSiteDir 验证站点目录
func (cd *CloudflareDeployer) validateSiteDir(siteDir string) error {
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

// createDeployment 创建部署
func (cd *CloudflareDeployer) createDeployment(siteDir string) (string, error) {
	cd.logger.Info("创建 Cloudflare Pages 部署")

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects/%s/deployments", 
		cd.config.AccountID, cd.config.ProjectName)

	payload := map[string]interface{}{
		"branch":        cd.config.Branch,
		"framework":     cd.config.Framework,
		"build_command": cd.config.BuildCommand,
		"output_dir":    cd.config.OutputDir,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cd.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", cd.config.Email)

	resp, err := cd.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("创建部署失败: %s, 响应: %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result["success"] != true {
		return "", fmt.Errorf("API 调用失败: %v", result["errors"])
	}

	deploymentID := result["result"].(map[string]interface{})["id"].(string)
	cd.logger.Info("部署创建成功，ID:", deploymentID)

	return deploymentID, nil
}

// uploadFiles 上传文件
func (cd *CloudflareDeployer) uploadFiles(deploymentID, siteDir string) error {
	cd.logger.Info("开始上传文件")

	// 获取所有文件
	files, err := cd.getAllFiles(siteDir)
	if err != nil {
		return err
	}

	cd.logger.Info("需要上传", len(files), "个文件")

	// 分批上传文件
	batchSize := 100
	for i := 0; i < len(files); i += batchSize {
		end := i + batchSize
		if end > len(files) {
			end = len(files)
		}

		batch := files[i:end]
		if err := cd.uploadBatch(deploymentID, batch); err != nil {
			return fmt.Errorf("上传批次 %d 失败: %w", i/batchSize+1, err)
		}

		cd.logger.Info(fmt.Sprintf("已上传 %d/%d 个文件", end, len(files)))
	}

	return nil
}

// getAllFiles 获取所有文件
func (cd *CloudflareDeployer) getAllFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 计算相对路径
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}

		return nil
	})

	return files, err
}

// uploadBatch 上传一批文件
func (cd *CloudflareDeployer) uploadBatch(deploymentID string, files []string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects/%s/deployments/%s/files", 
		cd.config.AccountID, cd.config.ProjectName, deploymentID)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件
	for _, file := range files {
		filePath := filepath.Join(cd.config.OutputDir, file)
		
		// 读取文件内容
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取文件失败 %s: %w", file, err)
		}

		// 创建表单字段
		part, err := writer.CreateFormFile(file, file)
		if err != nil {
			return fmt.Errorf("创建表单字段失败: %w", err)
		}

		// 写入文件内容
		if _, err := part.Write(content); err != nil {
			return fmt.Errorf("写入文件内容失败: %w", err)
		}
	}

	writer.Close()

	// 发送请求
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cd.config.APIKey)
	req.Header.Set("X-Auth-Email", cd.config.Email)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := cd.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("上传失败: %s, 响应: %s", resp.Status, string(body))
	}

	return nil
}

// finalizeDeployment 完成部署
func (cd *CloudflareDeployer) finalizeDeployment(deploymentID string) error {
	cd.logger.Info("完成部署")

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects/%s/deployments/%s", 
		cd.config.AccountID, cd.config.ProjectName, deploymentID)

	payload := map[string]interface{}{
		"status": "ready",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cd.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", cd.config.Email)

	resp, err := cd.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("完成部署失败: %s, 响应: %s", resp.Status, string(body))
	}

	// 获取部署 URL
	deploymentURL := fmt.Sprintf("https://%s.pages.dev", cd.config.ProjectName)
	cd.logger.Info("部署完成，访问地址:", deploymentURL)

	return nil
}

// GetDeploymentStatus 获取部署状态
func (cd *CloudflareDeployer) GetDeploymentStatus(deploymentID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects/%s/deployments/%s", 
		cd.config.AccountID, cd.config.ProjectName, deploymentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cd.config.APIKey)
	req.Header.Set("X-Auth-Email", cd.config.Email)

	resp, err := cd.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取部署状态失败: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListDeployments 列出部署历史
func (cd *CloudflareDeployer) ListDeployments() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/pages/projects/%s/deployments", 
		cd.config.AccountID, cd.config.ProjectName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cd.config.APIKey)
	req.Header.Set("X-Auth-Email", cd.config.Email)

	resp, err := cd.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取部署列表失败: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result["success"] != true {
		return nil, fmt.Errorf("API 调用失败: %v", result["errors"])
	}

	deployments := result["result"].([]interface{})
	resultList := make([]map[string]interface{}, len(deployments))
	
	for i, deployment := range deployments {
		resultList[i] = deployment.(map[string]interface{})
	}

	return resultList, nil
}
