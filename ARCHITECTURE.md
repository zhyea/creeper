# Creeper 系统架构设计文档

## 🏛️ 企业级架构概览

Creeper 静态小说站点生成器采用了企业级架构设计，应用了 20+ 种经典设计模式，构建了一个高度模块化、可扩展、易维护的系统。

## 🎯 设计模式应用总览

### 1. 创建型模式 (Creational Patterns)

#### 🏭 抽象工厂模式 (Abstract Factory Pattern)
- **位置**: `internal/factory/abstract_factory.go`
- **用途**: 管理不同类型的生成器
- **实现**: `AbstractGeneratorFactory`、`StaticGeneratorFactory`、`EnhancedGeneratorFactory`、`MinimalGeneratorFactory`
- **优势**: 支持多种生成器配置，易于扩展新类型

#### 🏗️ 建造者模式 (Builder Pattern)
- **位置**: `internal/config/builder.go`
- **用途**: 构建配置对象
- **实现**: `ConfigBuilder`、链式调用
- **优势**: 灵活的配置构建，支持默认值

#### 🏭 工厂方法模式 (Factory Method Pattern)
- **位置**: `internal/generator/template_factory.go`
- **用途**: 创建不同类型的模板
- **实现**: `TemplateFactory`、模板类型管理
- **优势**: 统一的模板创建接口

#### 🏭 单例模式 (Singleton Pattern)
- **位置**: `internal/common/singleton.go`
- **用途**: 全局资源管理
- **实现**: `Logger`、`ResourceManager`、`ConfigCache`
- **优势**: 确保全局唯一实例

### 2. 结构型模式 (Structural Patterns)

#### 🌉 适配器模式 (Adapter Pattern)
- **位置**: `internal/parser/adapter.go`
- **用途**: 统一 TXT 和 Markdown 输出格式
- **实现**: `ContentAdapter`、`ChapterAdapter`
- **优势**: 不同格式的统一处理

#### 🌉 桥接模式 (Bridge Pattern)
- **位置**: `internal/config/bridge.go`
- **用途**: 分离配置存储和验证
- **实现**: `ConfigStorage`、`ConfigValidator`、`ConfigBridge`
- **优势**: 存储和验证的解耦

#### 🌳 组合模式 (Composite Pattern)
- **位置**: `internal/common/composite.go`
- **用途**: 资源树形结构管理
- **实现**: `ResourceComponent`、`FileResource`、`DirectoryResource`
- **优势**: 统一的资源操作接口

#### 🦋 享元模式 (Flyweight Pattern)
- **位置**: `internal/common/flyweight.go`
- **用途**: 优化内存使用
- **实现**: `FlyweightFactory`、`ThemeFlyweight`、`FontFlyweight`
- **优势**: 减少内存占用，提高性能

#### 🎭 外观模式 (Facade Pattern)
- **位置**: `internal/facade/creeper_facade.go`
- **用途**: 统一系统入口
- **实现**: `CreeperFacade`
- **优势**: 简化复杂子系统调用

#### 🛡️ 代理模式 (Proxy Pattern)
- **位置**: `internal/proxy/resource_proxy.go`
- **用途**: 资源访问控制
- **实现**: `CachedResource`、`SecurityResourceProxy`
- **优势**: 缓存、安全控制

### 3. 行为型模式 (Behavioral Patterns)

#### ⛓️ 责任链模式 (Chain of Responsibility Pattern)
- **位置**: `internal/chain/error_handler.go`
- **用途**: 错误处理链
- **实现**: `ErrorHandler`、`LoggingErrorHandler`、`RetryErrorHandler`
- **优势**: 分级错误处理

#### 💉 依赖注入 (Dependency Injection)
- **位置**: `internal/di/container.go`
- **用途**: 控制反转
- **实现**: `Container`、`ServiceBuilder`、`ServiceLocator`
- **优势**: 松耦合，易于测试

#### 🧮 解释器模式 (Interpreter Pattern)
- **位置**: `internal/config/interpreter.go`
- **用途**: 配置表达式解析
- **实现**: `Expression`、`ConfigInterpreter`、`ConfigParser`
- **优势**: 灵活的配置表达式

#### 🔄 迭代器模式 (Iterator Pattern)
- **位置**: `internal/deploy/iterator.go`
- **用途**: 文件遍历
- **实现**: `FileIterator`、`DirectoryIterator`、`BatchIterator`
- **优势**: 统一的遍历接口

#### 💾 备忘录模式 (Memento Pattern)
- **位置**: `internal/deploy/memento.go`
- **用途**: 部署状态保存
- **实现**: `DeploymentMemento`、`DeploymentCaretaker`
- **优势**: 状态恢复和历史记录

#### 🤝 中介者模式 (Mediator Pattern)
- **位置**: `internal/mediator/mediator.go`
- **用途**: 组件间协调
- **实现**: `CreeperMediator`、组件通信
- **优势**: 松耦合的组件交互

#### 👁️ 观察者模式 (Observer Pattern)
- **位置**: `internal/deploy/observer.go`
- **用途**: 事件通知
- **实现**: `DeploymentObserver`、`DeploymentEventManager`
- **优势**: 事件驱动的架构

#### 🎯 状态模式 (State Pattern)
- **位置**: `internal/parser/state.go`
- **用途**: 解析状态管理
- **实现**: `ParseState`、`StatefulTxtParser`
- **优势**: 状态转换管理

#### 📋 策略模式 (Strategy Pattern)
- **位置**: `internal/parser/strategy.go`
- **用途**: 解析策略选择
- **实现**: `ParseStrategy`、`StrategyManager`
- **优势**: 灵活的解析策略

#### 📋 模板方法模式 (Template Method Pattern)
- **位置**: `internal/deploy/template.go`
- **用途**: 部署流程模板
- **实现**: `DeployTemplate`、`BaseDeployTemplate`
- **优势**: 统一的部署流程

#### 🎭 访问者模式 (Visitor Pattern)
- **位置**: `internal/parser/visitor.go`
- **用途**: 章节处理
- **实现**: `ChapterVisitor`、章节类型处理
- **优势**: 扩展章节处理逻辑

## 🏗️ 系统架构层次

### 1. 表现层 (Presentation Layer)
- **外观模式**: 统一系统入口
- **代理模式**: 资源访问控制
- **观察者模式**: 事件通知

### 2. 业务逻辑层 (Business Logic Layer)
- **策略模式**: 解析策略选择
- **状态模式**: 状态管理
- **模板方法模式**: 流程模板
- **责任链模式**: 错误处理

### 3. 数据访问层 (Data Access Layer)
- **桥接模式**: 配置管理
- **适配器模式**: 格式转换
- **迭代器模式**: 数据遍历

### 4. 基础设施层 (Infrastructure Layer)
- **单例模式**: 全局资源
- **享元模式**: 内存优化
- **组合模式**: 资源管理
- **备忘录模式**: 状态保存

## 🚀 架构优势

### 1. 高度模块化
- 每个模块职责单一
- 模块间松耦合
- 易于独立开发和测试

### 2. 强扩展性
- 新功能易于添加
- 支持插件化架构
- 配置驱动的扩展

### 3. 高性能
- 享元模式优化内存
- 缓存机制提升速度
- 异步处理提高并发

### 4. 高可靠性
- 完整的错误处理链
- 状态恢复机制
- 事件驱动的监控

### 5. 易维护性
- 清晰的代码结构
- 统一的编码规范
- 完善的文档说明

## 📊 系统特性

### 1. 多格式支持
- Markdown 文件解析
- TXT 文件解析
- 多卷结构支持
- 自定义格式扩展

### 2. 智能部署
- Cloudflare Pages 部署
- GitHub Pages 部署
- Vercel 部署
- Netlify 部署

### 3. 主题系统
- 多主题支持
- 动态主题切换
- 自定义主题配置
- 响应式设计

### 4. 性能优化
- 静态站点生成
- 资源压缩
- 缓存机制
- CDN 支持

### 5. 监控统计
- 部署状态监控
- 性能指标统计
- 错误日志记录
- 用户行为分析

## 🎯 设计原则

### 1. SOLID 原则
- **单一职责原则**: 每个类只有一个职责
- **开闭原则**: 对扩展开放，对修改关闭
- **里氏替换原则**: 子类可以替换父类
- **接口隔离原则**: 接口要小而专一
- **依赖倒置原则**: 依赖抽象而非具体

### 2. DRY 原则
- 避免重复代码
- 提取公共组件
- 统一配置管理

### 3. KISS 原则
- 保持简单
- 易于理解
- 易于维护

## 🔮 未来扩展

### 1. 新功能支持
- 更多文件格式
- 更多部署平台
- 更多主题样式
- 更多插件功能

### 2. 性能优化
- 并行处理
- 增量构建
- 智能缓存
- 资源优化

### 3. 用户体验
- 可视化配置
- 实时预览
- 拖拽操作
- 智能提示

## 📚 总结

Creeper 通过应用 20+ 种设计模式，构建了一个企业级的静态站点生成器。系统具有高度的模块化、可扩展性和可维护性，为静态站点生成提供了完整的解决方案。

这种架构设计不仅满足了当前的需求，还为未来的功能扩展和性能优化奠定了坚实的基础。通过设计模式的应用，Creeper 实现了代码的复用、解耦和扩展，达到了企业级应用的标准。
