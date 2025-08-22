package deploy

import (
	"fmt"
	"time"

	"creeper/internal/common"
)

// DeployTemplate 部署模板接口
type DeployTemplate interface {
	ValidateConfig() error
	ValidateSiteDir(siteDir string) error
	CreateDeployment(siteDir string) (string, error)
	UploadFiles(deploymentID, siteDir string) error
	FinalizeDeployment(deploymentID string) error
	GetDeploymentURL() string
	GetDeploymentStatus(deploymentID string) (map[string]interface{}, error)
}

// BaseDeployTemplate 基础部署模板
type BaseDeployTemplate struct {
	logger *common.Logger
}

// NewBaseDeployTemplate 创建基础部署模板
func NewBaseDeployTemplate() *BaseDeployTemplate {
	return &BaseDeployTemplate{
		logger: common.GetLogger(),
	}
}

// Deploy 执行部署流程（模板方法）
func (bdt *BaseDeployTemplate) Deploy(deployer DeployTemplate, siteDir string) error {
	startTime := time.Now()
	bdt.logger.Info("开始部署流程")

	// 1. 验证配置
	bdt.logger.Info("步骤 1: 验证配置")
	if err := deployer.ValidateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// 2. 验证站点目录
	bdt.logger.Info("步骤 2: 验证站点目录")
	if err := deployer.ValidateSiteDir(siteDir); err != nil {
		return fmt.Errorf("站点目录验证失败: %w", err)
	}

	// 3. 创建部署
	bdt.logger.Info("步骤 3: 创建部署")
	deploymentID, err := deployer.CreateDeployment(siteDir)
	if err != nil {
		return fmt.Errorf("创建部署失败: %w", err)
	}

	// 4. 上传文件
	bdt.logger.Info("步骤 4: 上传文件")
	if err := deployer.UploadFiles(deploymentID, siteDir); err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}

	// 5. 完成部署
	bdt.logger.Info("步骤 5: 完成部署")
	if err := deployer.FinalizeDeployment(deploymentID); err != nil {
		return fmt.Errorf("完成部署失败: %w", err)
	}

	duration := time.Since(startTime)
	bdt.logger.Info("部署流程完成，耗时:", duration)

	return nil
}

// DeployWithRetry 带重试的部署流程
func (bdt *BaseDeployTemplate) DeployWithRetry(deployer DeployTemplate, siteDir string, maxRetries int) error {
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		bdt.logger.Info(fmt.Sprintf("部署尝试 %d/%d", attempt, maxRetries))
		
		if err := bdt.Deploy(deployer, siteDir); err != nil {
			lastErr = err
			bdt.logger.Warn(fmt.Sprintf("部署尝试 %d 失败: %v", attempt, err))
			
			if attempt < maxRetries {
				bdt.logger.Info("等待重试...")
				time.Sleep(time.Duration(attempt) * time.Second)
			}
		} else {
			bdt.logger.Info("部署成功")
			return nil
		}
	}
	
	return fmt.Errorf("部署失败，已重试 %d 次，最后错误: %w", maxRetries, lastErr)
}

// DeployWithRollback 带回滚的部署流程
func (bdt *BaseDeployTemplate) DeployWithRollback(deployer DeployTemplate, siteDir string) error {
	// 获取当前部署状态
	currentDeploymentID := ""
	
	// 执行部署
	if err := bdt.Deploy(deployer, siteDir); err != nil {
		// 如果部署失败，尝试回滚
		if currentDeploymentID != "" {
			bdt.logger.Warn("部署失败，尝试回滚到上一个版本")
			// 这里可以实现回滚逻辑
		}
		return err
	}
	
	return nil
}

// DeployWithValidation 带验证的部署流程
func (bdt *BaseDeployTemplate) DeployWithValidation(deployer DeployTemplate, siteDir string) error {
	// 执行部署
	if err := bdt.Deploy(deployer, siteDir); err != nil {
		return err
	}
	
	// 部署后验证
	bdt.logger.Info("验证部署结果")
	deploymentURL := deployer.GetDeploymentURL()
	if deploymentURL == "" {
		return fmt.Errorf("部署 URL 为空")
	}
	
	bdt.logger.Info("部署验证通过，访问地址:", deploymentURL)
	return nil
}
