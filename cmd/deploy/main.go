package main

import (
	"flag"
	"fmt"
	"log"

	"creeper/internal/deploy"
)

func main() {
	var (
		configPath = flag.String("config", "deploy-config.yaml", "部署配置文件路径")
		siteDir    = flag.String("site", "dist", "站点目录")
		init       = flag.Bool("init", false, "初始化部署配置")
		deployType = flag.String("type", "cloudflare", "部署类型 (cloudflare|github|vercel|netlify)")
		status     = flag.Bool("status", false, "查看部署状态")
		list       = flag.Bool("list", false, "列出部署历史")
	)
	flag.Parse()

	// 初始化部署配置
	if *init {
		if err := initDeployConfig(*deployType, *configPath); err != nil {
			log.Fatalf("初始化部署配置失败: %v", err)
		}
		fmt.Printf("✅ 部署配置已创建: %s\n", *configPath)
		fmt.Printf("📝 请编辑配置文件并填入正确的参数\n")
		return
	}

	// 加载部署配置
	deployConfig, err := deploy.LoadDeployConfig(*configPath)
	if err != nil {
		log.Fatalf("加载部署配置失败: %v", err)
	}

	// 创建部署管理器
	deployManager := deploy.NewDeployManager(deployConfig)

	// 初始化部署管理器
	if err := deployManager.Initialize(); err != nil {
		log.Fatalf("初始化部署管理器失败: %v", err)
	}

	// 查看部署状态
	if *status {
		status, err := deployManager.GetStatus()
		if err != nil {
			log.Fatalf("获取部署状态失败: %v", err)
		}
		fmt.Printf("📊 部署状态:\n")
		for key, value := range status {
			fmt.Printf("  %s: %v\n", key, value)
		}
		return
	}

	// 列出部署历史
	if *list {
		fmt.Printf("⚠️  当前部署类型不支持列出部署历史\n")
		return
	}

	// 执行部署
	fmt.Printf("🚀 开始部署到 %s...\n", deployConfig.Type)
	fmt.Printf("📁 站点目录: %s\n", *siteDir)

	if err := deployManager.Deploy(*siteDir); err != nil {
		log.Fatalf("部署失败: %v", err)
	}

	deploymentURL := deployManager.GetDeploymentURL()
	fmt.Printf("✅ 部署完成！\n")
	fmt.Printf("🌐 访问地址: %s\n", deploymentURL)
}

// initDeployConfig 初始化部署配置
func initDeployConfig(deployType, configPath string) error {
	var deployTypeEnum deploy.DeployType
	switch deployType {
	case "cloudflare":
		deployTypeEnum = deploy.CloudflarePages
	case "github":
		deployTypeEnum = deploy.GitHubPages
	case "vercel":
		deployTypeEnum = deploy.Vercel
	case "netlify":
		deployTypeEnum = deploy.Netlify
	default:
		return fmt.Errorf("不支持的部署类型: %s", deployType)
	}

	// 创建默认配置
	config := deploy.CreateDefaultDeployConfig(deployTypeEnum)

	// 保存配置
	if err := deploy.SaveDeployConfig(config, configPath); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	return nil
}
