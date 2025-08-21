package mediator

import (
	"fmt"
	"sync"

	"creeper/internal/config"
	"creeper/internal/parser"
	"creeper/internal/generator"
	"creeper/internal/common"
)

// ComponentType 组件类型
type ComponentType string

const (
	ParserComponent    ComponentType = "parser"
	GeneratorComponent ComponentType = "generator"
	ConfigComponent    ComponentType = "config"
	LoggerComponent    ComponentType = "logger"
)

// Message 消息结构
type Message struct {
	Type      string
	From      ComponentType
	To        ComponentType
	Data      interface{}
	Timestamp int64
	ID        string
}

// Component 组件接口
type Component interface {
	GetType() ComponentType
	HandleMessage(message *Message) error
	SetMediator(mediator Mediator)
}

// Mediator 中介者接口
type Mediator interface {
	Register(component Component)
	Unregister(componentType ComponentType)
	Send(message *Message) error
	Broadcast(message *Message) error
	GetComponent(componentType ComponentType) Component
}

// CreeperMediator Creeper 中介者实现
type CreeperMediator struct {
	components map[ComponentType]Component
	logger     *common.Logger
	mutex      sync.RWMutex
}

// NewCreeperMediator 创建中介者
func NewCreeperMediator() *CreeperMediator {
	return &CreeperMediator{
		components: make(map[ComponentType]Component),
		logger:     common.GetLogger(),
	}
}

// Register 注册组件
func (cm *CreeperMediator) Register(component Component) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	componentType := component.GetType()
	cm.components[componentType] = component
	component.SetMediator(cm)
	
	cm.logger.Info("注册组件:", componentType)
}

// Unregister 注销组件
func (cm *CreeperMediator) Unregister(componentType ComponentType) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	delete(cm.components, componentType)
	cm.logger.Info("注销组件:", componentType)
}

// Send 发送消息
func (cm *CreeperMediator) Send(message *Message) error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	target, exists := cm.components[message.To]
	if !exists {
		return fmt.Errorf("目标组件不存在: %s", message.To)
	}
	
	cm.logger.Debug("发送消息:", message.From, "->", message.To, "类型:", message.Type)
	
	return target.HandleMessage(message)
}

// Broadcast 广播消息
func (cm *CreeperMediator) Broadcast(message *Message) error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	cm.logger.Debug("广播消息:", message.From, "类型:", message.Type)
	
	for componentType, component := range cm.components {
		if componentType != message.From {
			if err := component.HandleMessage(message); err != nil {
				cm.logger.Warn("组件处理消息失败:", componentType, err)
			}
		}
	}
	
	return nil
}

// GetComponent 获取组件
func (cm *CreeperMediator) GetComponent(componentType ComponentType) Component {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	return cm.components[componentType]
}

// ParserComponent 解析器组件
type ParserComponentImpl struct {
	parser   *parser.Parser
	mediator Mediator
}

// NewParserComponent 创建解析器组件
func NewParserComponent(parser *parser.Parser) *ParserComponentImpl {
	return &ParserComponentImpl{
		parser: parser,
	}
}

func (pc *ParserComponentImpl) GetType() ComponentType {
	return ParserComponent
}

func (pc *ParserComponentImpl) SetMediator(mediator Mediator) {
	pc.mediator = mediator
}

func (pc *ParserComponentImpl) HandleMessage(message *Message) error {
	switch message.Type {
	case "parse_novel":
		if path, ok := message.Data.(string); ok {
			novel, err := pc.parser.ParseNovel(path)
			if err != nil {
				return err
			}
			
			// 通知生成器
			response := &Message{
				Type: "novel_parsed",
				From: ParserComponent,
				To:   GeneratorComponent,
				Data: novel,
			}
			
			return pc.mediator.Send(response)
		}
		
	case "get_supported_formats":
		// 返回支持的格式
		formats := []string{"md", "txt"}
		response := &Message{
			Type: "supported_formats",
			From: ParserComponent,
			To:   message.From,
			Data: formats,
		}
		
		return pc.mediator.Send(response)
	}
	
	return nil
}

// GeneratorComponent 生成器组件
type GeneratorComponentImpl struct {
	generator *generator.Generator
	mediator  Mediator
}

// NewGeneratorComponent 创建生成器组件
func NewGeneratorComponent(generator *generator.Generator) *GeneratorComponentImpl {
	return &GeneratorComponentImpl{
		generator: generator,
	}
}

func (gc *GeneratorComponentImpl) GetType() ComponentType {
	return GeneratorComponent
}

func (gc *GeneratorComponentImpl) SetMediator(mediator Mediator) {
	gc.mediator = mediator
}

func (gc *GeneratorComponentImpl) HandleMessage(message *Message) error {
	switch message.Type {
	case "novel_parsed":
		// 处理解析完成的小说
		if novel, ok := message.Data.(*parser.Novel); ok {
			// 这里可以添加小说处理逻辑
			common.GetLogger().Info("接收到解析完成的小说:", novel.Title)
		}
		
	case "generate_website":
		return gc.generator.Generate()
		
	case "start_server":
		if port, ok := message.Data.(int); ok {
			return gc.generator.Serve(port)
		}
	}
	
	return nil
}

// ConfigComponent 配置组件
type ConfigComponentImpl struct {
	config   *config.Config
	mediator Mediator
}

// NewConfigComponent 创建配置组件
func NewConfigComponent(config *config.Config) *ConfigComponentImpl {
	return &ConfigComponentImpl{
		config: config,
	}
}

func (cc *ConfigComponentImpl) GetType() ComponentType {
	return ConfigComponent
}

func (cc *ConfigComponentImpl) SetMediator(mediator Mediator) {
	cc.mediator = mediator
}

func (cc *ConfigComponentImpl) HandleMessage(message *Message) error {
	switch message.Type {
	case "get_config":
		response := &Message{
			Type: "config_data",
			From: ConfigComponent,
			To:   message.From,
			Data: cc.config,
		}
		
		return cc.mediator.Send(response)
		
	case "update_config":
		if updates, ok := message.Data.(map[string]interface{}); ok {
			// 使用建造者模式更新配置
			builder := config.Builder().WithDefaults()
			
			// 应用更新
			for key, value := range updates {
				switch key {
				case "site_title":
					if title, ok := value.(string); ok {
						builder.WithSiteTitle(title)
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
			
			cc.config = builder.Build()
			
			// 广播配置更新消息
			broadcast := &Message{
				Type: "config_updated",
				From: ConfigComponent,
				Data: cc.config,
			}
			
			return cc.mediator.Broadcast(broadcast)
		}
	}
	
	return nil
}

// LoggerComponent 日志组件
type LoggerComponentImpl struct {
	logger   *common.Logger
	mediator Mediator
}

// NewLoggerComponent 创建日志组件
func NewLoggerComponent() *LoggerComponentImpl {
	return &LoggerComponentImpl{
		logger: common.GetLogger(),
	}
}

func (lc *LoggerComponentImpl) GetType() ComponentType {
	return LoggerComponent
}

func (lc *LoggerComponentImpl) SetMediator(mediator Mediator) {
	lc.mediator = mediator
}

func (lc *LoggerComponentImpl) HandleMessage(message *Message) error {
	switch message.Type {
	case "log_info":
		if msg, ok := message.Data.(string); ok {
			lc.logger.Info(msg)
		}
		
	case "log_warn":
		if msg, ok := message.Data.(string); ok {
			lc.logger.Warn(msg)
		}
		
	case "log_error":
		if msg, ok := message.Data.(string); ok {
			lc.logger.Error(msg)
		}
	}
	
	return nil
}

// MediatorFactory 中介者工厂
type MediatorFactory struct{}

// CreateMediator 创建完整的中介者系统
func (mf *MediatorFactory) CreateMediator(config *config.Config) *CreeperMediator {
	mediator := NewCreeperMediator()
	
	// 创建并注册组件
	parserComp := NewParserComponent(parser.New())
	generatorComp := NewGeneratorComponent(generator.New(config))
	configComp := NewConfigComponent(config)
	loggerComp := NewLoggerComponent()
	
	mediator.Register(parserComp)
	mediator.Register(generatorComp)
	mediator.Register(configComp)
	mediator.Register(loggerComp)
	
	return mediator
}
