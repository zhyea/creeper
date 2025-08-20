package main

import (
	"flag"
	"fmt"
	"log"

	"creeper/internal/config"
	"creeper/internal/generator"
)

func main() {
	var (
		configPath = flag.String("config", "config.yaml", "配置文件路径")
		inputDir   = flag.String("input", "novels", "小说文件输入目录")
		outputDir  = flag.String("output", "dist", "静态站点输出目录")
		serve      = flag.Bool("serve", false, "生成后启动本地服务器")
		port       = flag.Int("port", 8080, "本地服务器端口")
	)
	flag.Parse()

	// 读取配置文件
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("警告：无法读取配置文件 %s，使用默认配置: %v", *configPath, err)
		cfg = config.Default()
	}

	// 覆盖命令行参数
	if *inputDir != "novels" {
		cfg.InputDir = *inputDir
	}
	if *outputDir != "dist" {
		cfg.OutputDir = *outputDir
	}

	// 创建生成器
	gen := generator.New(cfg)

	// 生成静态站点
	fmt.Printf("开始生成静态小说站点...\n")
	fmt.Printf("输入目录: %s\n", cfg.InputDir)
	fmt.Printf("输出目录: %s\n", cfg.OutputDir)

	if err := gen.Generate(); err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("✅ 静态站点生成完成！\n")

	// 启动本地服务器
	if *serve {
		fmt.Printf("启动本地服务器 http://localhost:%d\n", *port)
		if err := gen.Serve(*port); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}
}
