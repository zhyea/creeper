package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"creeper/internal/config"
	"creeper/internal/facade"
	"creeper/internal/mediator"
	"creeper/internal/factory"
	"creeper/internal/chain"
	"creeper/internal/di"
	"creeper/internal/common"
	"creeper/internal/parser"
	"creeper/internal/generator"
)

// Application 应用程序
type Application struct {
	facade       *facade.CreeperFacade
	mediator     mediator.Mediator
	errorManager *chain.ErrorManager
	container    *di.Container
	logger       *common.Logger
}

// NewApplication 创建应用程序
func NewApplication() *Application {
	return &Application{
		logger: common.GetLogger(),
	}
}

// Initialize 初始化应用程序
func (app *Application) Initialize(configPath string, generatorType factory.GeneratorType) error {
	app.logger.Info("初始化 Creeper 应用程序")
	
	// 1. 初始化依赖注入容器
	if err := app.initializeDI(configPath, generatorType); err != nil {
		return fmt.Errorf("初始化依赖注入失败: %w", err)
	}
	
	// 2. 初始化错误处理链
	app.initializeErrorHandling()
	
	// 3. 初始化中介者
	if err := app.initializeMediator(); err != nil {
		return fmt.Errorf("初始化中介者失败: %w", err)
	}
	
	// 4. 初始化外观
	if err := app.initializeFacade(configPath); err != nil {
		return fmt.Errorf("初始化外观失败: %w", err)
	}
	
	// 5. 设置信号处理
	app.setupSignalHandling()
	
	app.logger.Info("应用程序初始化完成")
	
	return nil
}

// initializeDI 初始化依赖注入
func (app *Application) initializeDI(configPath string, generatorType factory.GeneratorType) error {
	builder := di.NewServiceBuilder()
	
	// 注册配置服务
	builder.AddSingleton((*config.Config)(nil), func(container *di.Container) (interface{}, error) {
		cfg, err := config.Load(configPath)
		if err != nil {
			app.logger.Warn("使用默认配置:", err)
			cfg = config.Default()
		}
		return cfg, nil
	})
	
	// 注册解析器服务
	builder.AddTransient((*parser.Parser)(nil), func(container *di.Container) (interface{}, error) {
		return parser.New(), nil
	})
	
	// 注册生成器服务
	builder.AddSingleton((*generator.Generator)(nil), func(container *di.Container) (interface{}, error) {
		cfg, err := container.Resolve((*config.Config)(nil))
		if err != nil {
			return nil, err
		}
		
		// 使用抽象工厂创建生成器
		factoryRegistry := factory.NewGeneratorFactoryRegistry()
		suite, err := factoryRegistry.CreateGeneratorSuite(generatorType, cfg.(*config.Config))
		if err != nil {
			return nil, err
		}
		
		return suite.Generator, nil
	})
	
	// 注册错误管理器
	builder.AddSingleton((*chain.ErrorManager)(nil), func(container *di.Container) (interface{}, error) {
		return chain.NewErrorManager(), nil
	})
	
	app.container = builder.Build()
	
	// 设置服务定位器
	di.GetServiceLocator().SetContainer(app.container)
	
	return nil
}

// initializeErrorHandling 初始化错误处理
func (app *Application) initializeErrorHandling() {
	errorManager, err := app.container.Resolve((*chain.ErrorManager)(nil))
	if err != nil {
		log.Fatal("无法解析错误管理器:", err)
	}
	
	app.errorManager = errorManager.(*chain.ErrorManager)
	
	// 设置关闭回调
	app.errorManager.SetShutdownCallback(func() {
		app.logger.Error("检测到关键错误，系统即将关闭")
		app.Shutdown()
		os.Exit(1)
	})
}

// initializeMediator 初始化中介者
func (app *Application) initializeMediator() error {
	// 获取服务
	parserService, err := app.container.Resolve((*parser.Parser)(nil))
	if err != nil {
		return err
	}
	
	generatorService, err := app.container.Resolve((*generator.Generator)(nil))
	if err != nil {
		return err
	}
	
	configService, err := app.container.Resolve((*config.Config)(nil))
	if err != nil {
		return err
	}
	
	// 创建中介者
	app.mediator = mediator.NewCreeperMediator()
	
	// 注册组件
	app.mediator.Register(mediator.NewParserComponent(parserService.(*parser.Parser)))
	app.mediator.Register(mediator.NewGeneratorComponent(generatorService.(*generator.Generator)))
	app.mediator.Register(mediator.NewConfigComponent(configService.(*config.Config)))
	app.mediator.Register(mediator.NewLoggerComponent())
	
	return nil
}

// initializeFacade 初始化外观
func (app *Application) initializeFacade(configPath string) error {
	facade, err := facade.NewCreeperFacade(configPath)
	if err != nil {
		return err
	}
	
	app.facade = facade
	
	// 验证设置
	if err := app.facade.ValidateSetup(); err != nil {
		return app.errorManager.HandleError(err, chain.SeverityError, "application", "validate_setup", nil)
	}
	
	return nil
}

// setupSignalHandling 设置信号处理
func (app *Application) setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		sig := <-sigChan
		app.logger.Info("接收到信号:", sig)
		app.Shutdown()
		os.Exit(0)
	}()
}

// Generate 生成网站
func (app *Application) Generate() error {
	app.logger.Info("开始生成网站")
	
	if err := app.facade.GenerateWebsite(); err != nil {
		return app.errorManager.HandleError(err, chain.SeverityError, "application", "generate", nil)
	}
	
	return nil
}

// Serve 启动服务器
func (app *Application) Serve(port int) error {
	app.logger.Info("启动服务器，端口:", port)
	
	if err := app.facade.ServeWebsite(port); err != nil {
		return app.errorManager.HandleError(err, chain.SeverityCritical, "application", "serve", map[string]interface{}{
			"port": port,
		})
	}
	
	return nil
}

// GetStatus 获取应用状态
func (app *Application) GetStatus() map[string]interface{} {
	status := app.facade.GetSystemStatus()
	
	// 添加错误统计
	status["errors"] = app.errorManager.GetErrorStatistics()
	
	// 添加服务信息
	status["services"] = app.container.GetRegisteredServices()
	
	return status
}

// Shutdown 关闭应用程序
func (app *Application) Shutdown() {
	app.logger.Info("开始关闭应用程序")
	
	if app.facade != nil {
		app.facade.Shutdown()
	}
	
	if app.errorManager != nil {
		app.errorManager.ClearErrorLog()
	}
	
	app.logger.Info("应用程序关闭完成")
}

func main() {
	var (
		configPath    = flag.String("config", "config.yaml", "配置文件路径")
		inputDir      = flag.String("input", "novels", "小说文件输入目录")
		outputDir     = flag.String("output", "dist", "静态站点输出目录")
		serve         = flag.Bool("serve", false, "生成后启动本地服务器")
		port          = flag.Int("port", 8080, "本地服务器端口")
		generatorType = flag.String("generator", "enhanced", "生成器类型 (static|enhanced|minimal)")
		verbose       = flag.Bool("verbose", false, "详细输出")
		status        = flag.Bool("status", false, "显示系统状态")
	)
	flag.Parse()
	
	// 创建应用程序
	app := NewApplication()
	
	// 解析生成器类型
	var genType factory.GeneratorType
	switch *generatorType {
	case "static":
		genType = factory.StaticGenerator
	case "minimal":
		genType = factory.MinimalGenerator
	case "enhanced":
		genType = factory.EnhancedGenerator
	default:
		genType = factory.EnhancedGenerator
	}
	
	// 初始化应用程序
	if err := app.Initialize(*configPath, genType); err != nil {
		log.Fatalf("应用程序初始化失败: %v", err)
	}
	
	// 如果只是查看状态
	if *status {
		status := app.GetStatus()
		fmt.Printf("🎯 Creeper 系统状态\n")
		fmt.Printf("==================\n")
		for key, value := range status {
			fmt.Printf("%s: %v\n", key, value)
		}
		return
	}
	
	// 更新配置（如果通过命令行指定）
	if *inputDir != "novels" || *outputDir != "dist" {
		updates := map[string]interface{}{
			"input_dir":  *inputDir,
			"output_dir": *outputDir,
		}
		
		if err := app.facade.UpdateConfig(updates); err != nil {
			log.Fatalf("更新配置失败: %v", err)
		}
	}
	
	// 生成网站
	if err := app.Generate(); err != nil {
		log.Fatalf("生成网站失败: %v", err)
	}
	
	fmt.Printf("✅ 静态站点生成完成！\n")
	fmt.Printf("📁 输出目录: %s\n", *outputDir)
	fmt.Printf("🎨 生成器类型: %s\n", genType)
	
	if *verbose {
		status := app.GetStatus()
		fmt.Printf("📊 系统状态: %v\n", status)
	}
	
	// 启动服务器
	if *serve {
		fmt.Printf("🚀 启动本地服务器 http://localhost:%d\n", *port)
		fmt.Printf("按 Ctrl+C 停止服务器\n")
		
		if err := app.Serve(*port); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}
}
