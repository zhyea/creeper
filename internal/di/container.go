package di

import (
	"fmt"
	"reflect"
	"sync"
)

// ServiceLifetime 服务生命周期
type ServiceLifetime int

const (
	Transient ServiceLifetime = iota // 瞬态
	Singleton                        // 单例
	Scoped                          // 作用域
)

// ServiceDescriptor 服务描述符
type ServiceDescriptor struct {
	ServiceType reflect.Type
	Lifetime    ServiceLifetime
	Factory     func(container *Container) (interface{}, error)
	Instance    interface{}
}

// Container 依赖注入容器
type Container struct {
	services map[reflect.Type]*ServiceDescriptor
	mutex    sync.RWMutex
}

// NewContainer 创建新的容器
func NewContainer() *Container {
	return &Container{
		services: make(map[reflect.Type]*ServiceDescriptor),
	}
}

// RegisterTransient 注册瞬态服务
func (c *Container) RegisterTransient(serviceType interface{}, factory func(*Container) (interface{}, error)) {
	c.register(serviceType, factory, Transient)
}

// RegisterSingleton 注册单例服务
func (c *Container) RegisterSingleton(serviceType interface{}, factory func(*Container) (interface{}, error)) {
	c.register(serviceType, factory, Singleton)
}

// RegisterInstance 注册实例
func (c *Container) RegisterInstance(serviceType interface{}, instance interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	t := reflect.TypeOf(serviceType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	c.services[t] = &ServiceDescriptor{
		ServiceType: t,
		Lifetime:    Singleton,
		Instance:    instance,
	}
}

// register 注册服务
func (c *Container) register(serviceType interface{}, factory func(*Container) (interface{}, error), lifetime ServiceLifetime) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	t := reflect.TypeOf(serviceType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	c.services[t] = &ServiceDescriptor{
		ServiceType: t,
		Lifetime:    lifetime,
		Factory:     factory,
	}
}

// Resolve 解析服务
func (c *Container) Resolve(serviceType interface{}) (interface{}, error) {
	t := reflect.TypeOf(serviceType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	c.mutex.RLock()
	descriptor, exists := c.services[t]
	c.mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("服务未注册: %s", t.Name())
	}
	
	switch descriptor.Lifetime {
	case Singleton:
		if descriptor.Instance != nil {
			return descriptor.Instance, nil
		}
		
		// 创建单例实例
		c.mutex.Lock()
		defer c.mutex.Unlock()
		
		// 双重检查
		if descriptor.Instance != nil {
			return descriptor.Instance, nil
		}
		
		instance, err := descriptor.Factory(c)
		if err != nil {
			return nil, err
		}
		
		descriptor.Instance = instance
		return instance, nil
		
	case Transient:
		return descriptor.Factory(c)
		
	default:
		return descriptor.Factory(c)
	}
}

// MustResolve 必须解析服务（panic on error）
func (c *Container) MustResolve(serviceType interface{}) interface{} {
	service, err := c.Resolve(serviceType)
	if err != nil {
		panic(fmt.Sprintf("无法解析服务: %v", err))
	}
	return service
}

// GetRegisteredServices 获取已注册的服务
func (c *Container) GetRegisteredServices() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	services := make([]string, 0, len(c.services))
	for serviceType := range c.services {
		services = append(services, serviceType.Name())
	}
	
	return services
}

// Clear 清空容器
func (c *Container) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.services = make(map[reflect.Type]*ServiceDescriptor)
}

// ServiceBuilder 服务构建器
type ServiceBuilder struct {
	container *Container
}

// NewServiceBuilder 创建服务构建器
func NewServiceBuilder() *ServiceBuilder {
	return &ServiceBuilder{
		container: NewContainer(),
	}
}

// AddTransient 添加瞬态服务
func (sb *ServiceBuilder) AddTransient(serviceType interface{}, factory func(*Container) (interface{}, error)) *ServiceBuilder {
	sb.container.RegisterTransient(serviceType, factory)
	return sb
}

// AddSingleton 添加单例服务
func (sb *ServiceBuilder) AddSingleton(serviceType interface{}, factory func(*Container) (interface{}, error)) *ServiceBuilder {
	sb.container.RegisterSingleton(serviceType, factory)
	return sb
}

// AddInstance 添加实例
func (sb *ServiceBuilder) AddInstance(serviceType interface{}, instance interface{}) *ServiceBuilder {
	sb.container.RegisterInstance(serviceType, instance)
	return sb
}

// Build 构建容器
func (sb *ServiceBuilder) Build() *Container {
	return sb.container
}

// ServiceLocator 服务定位器
type ServiceLocator struct {
	container *Container
	mutex     sync.RWMutex
}

var (
	serviceLocatorInstance *ServiceLocator
	serviceLocatorOnce     sync.Once
)

// GetServiceLocator 获取服务定位器单例
func GetServiceLocator() *ServiceLocator {
	serviceLocatorOnce.Do(func() {
		serviceLocatorInstance = &ServiceLocator{
			container: NewContainer(),
		}
	})
	return serviceLocatorInstance
}

// SetContainer 设置容器
func (sl *ServiceLocator) SetContainer(container *Container) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()
	sl.container = container
}

// GetService 获取服务
func (sl *ServiceLocator) GetService(serviceType interface{}) (interface{}, error) {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()
	return sl.container.Resolve(serviceType)
}

// MustGetService 必须获取服务
func (sl *ServiceLocator) MustGetService(serviceType interface{}) interface{} {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()
	return sl.container.MustResolve(serviceType)
}
