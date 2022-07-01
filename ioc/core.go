package ioc

import (
	"go-spring/configuration"
	"log"
	"reflect"
	"sort"
	"sync"
)

const (
	defaultConfigPath = "./config.yaml"
	defaultConfigType = "yaml"
)

var (
	ioc *Container
)

func init() {
	ioc = NewContainer()
}

func RegistryBeans(beans ...*Bean) *Container {
	ioc.AddBeans(beans...)
	return ioc
}

func SetConfigPath(path string) *Container {
	ioc.Configuration.Path = path
}

func SetConfigType(configType string) {
	ioc.Configuration.ConfigType = configType
}

// Container ioc容器
type Container struct {
	beans                    map[string]*Bean
	Configuration            *configuration.Configuration
	globalBeanPreProcessors  []BeanPreProcessor
	globalBeanPostProcessors []BeanPostProcessor
	containerPreProcessors   []ContainerPreProcessor
	containerPostProcessors  []ContainerPostProcessor
	isInited                 bool //是否已经初始化
	Logger                   *log.Logger
	lock                     *sync.Mutex
}

func NewContainer() *Container {
	return &Container{
		beans:         map[string]*Bean{},
		lock:          &sync.Mutex{},
		Configuration: &configuration.Configuration{Path: defaultConfigPath, ConfigType: defaultConfigType, Logger: log.Default()},
		Logger:        log.Default()}
}

// Init 容器初始化方法
func (c *Container) Init() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Logger.Println("Ioc container start init....")
	if c.containerPreProcessors != nil {
		sort.Slice(c.containerPreProcessors, func(i, j int) bool {
			return c.getProcessorPriority(c.containerPreProcessors[i]) > c.getProcessorPriority(c.containerPreProcessors[j])
		})
		for _, processor := range c.containerPreProcessors {
			processor.PreProcess(c)
		}
	}
	//加载配置
	c.Logger.Printf("Start loading the configuration from the path %s", c.Configuration.Path)
	c.Configuration.Load()
	c.Logger.Println("Load configuration complete")
	//实例化单例bean
	for name := range c.beans {
		c.GetBeanInstanceByName(name)
	}
	c.Logger.Println("Ioc container instance beans complete")
	c.isInited = true
	if c.containerPostProcessors != nil {
		sort.Slice(c.containerPostProcessors, func(i, j int) bool {
			return c.getProcessorPriority(c.containerPostProcessors[i]) > c.getProcessorPriority(c.containerPostProcessors)
		})
		for _, processor := range c.containerPostProcessors {
			processor.PostProcess(c)
		}
	}
	c.Logger.Println("Ioc container init complete")
}

// GetBeanInstanceByName 获取bean
func (c *Container) GetBeanInstanceByName(name string) interface{} {
	if bean := c.beans[name]; bean != nil {
		if bean.isSingleton && bean.instance != nil {
			return bean.instance
		}
		return c.InstanceBean(bean)
	}
	return nil
}

// GetBeanInstanceByStruct 通过结构体获取实例化的bean
func (c *Container) GetBeanInstanceByStruct(value interface{}) (interface{}, error) {
	//获取bean名称
	if beanNameProvider, ok := value.(BeanNameProvider); ok {
		return c.GetBeanInstanceByName(beanNameProvider.GetBeanName()), nil
	}
	rt := reflect.TypeOf(value)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return nil, TypeNotMatchError
	}
	return c.GetBeanInstanceByName(rt.Name()), nil
}

// AddBeans 添加bean
func (c *Container) AddBeans(beans ...*Bean) {
	for _, bean := range beans {
		if bean == nil {
			panic(NilError)
		}
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, bean := range beans {
		if c.beans[bean.name] != nil && c.beans[bean.name].priority > bean.priority {
			return
		}
		if c.globalBeanPreProcessors != nil {
			if bean.beanPreProcessors == nil {
				bean.beanPreProcessors = c.globalBeanPreProcessors
			} else {
				bean.beanPreProcessors = append(c.globalBeanPreProcessors, bean.beanPreProcessors...)
				sort.Slice(bean.beanPreProcessors, func(i, j int) bool {
					return c.getProcessorPriority(bean.beanPreProcessors[i]) > c.getProcessorPriority(bean.beanPreProcessors[j])
				})
			}
		}
		if c.globalBeanPostProcessors != nil {
			if bean.beanPostProcessors == nil {
				bean.beanPostProcessors = c.globalBeanPostProcessors
			} else {
				bean.beanPostProcessors = append(c.globalBeanPostProcessors, bean.beanPostProcessors...)
				sort.Slice(bean.beanPostProcessors, func(i, j int) bool {
					return c.getProcessorPriority(bean.beanPostProcessors[i]) > c.getProcessorPriority(bean.beanPostProcessors[j])
				})
			}
		}
		c.beans[bean.name] = bean
	}
}

func (c *Container) getProcessorPriority(processor interface{}) int {
	if priorityProcessor, ok := processor.(Priority); ok {
		return priorityProcessor.GetPriority()
	}
	return 0
}

func (c *Container) AddBeanPreProcessor(processor BeanPreProcessor) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if processor != nil {
		if c.globalBeanPreProcessors == nil {
			c.globalBeanPreProcessors = []BeanPreProcessor{processor}
		} else {
			c.globalBeanPreProcessors = append(c.globalBeanPreProcessors, processor)
		}
	}
}

func (c *Container) AddBeanPostProcessor(processor BeanPostProcessor) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if processor != nil {
		if c.globalBeanPostProcessors == nil {
			c.globalBeanPostProcessors = []BeanPostProcessor{processor}
		} else {
			c.globalBeanPostProcessors = append(c.globalBeanPostProcessors, processor)
		}
	}
}

func (c *Container) AddContainerPreProcessor(processor ContainerPreProcessor) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if processor != nil {
		if c.containerPreProcessors == nil {
			c.containerPreProcessors = []ContainerPreProcessor{processor}
		} else {
			c.containerPreProcessors = append(c.containerPreProcessors, processor)
		}
	}
}

func (c *Container) AddContainerPostProcessor(processor ContainerPostProcessor) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if processor != nil {
		if c.containerPostProcessors == nil {
			c.containerPostProcessors = []ContainerPostProcessor{processor}
		} else {
			c.containerPostProcessors = append(c.containerPostProcessors, processor)
		}
	}
}

// InstanceBean 实例化bean
func (c *Container) InstanceBean(bean *Bean) interface{} {
	c.Logger.Printf("start create bean:%s", bean.name)
	//1.调用前置处理器
	if bean.beanPreProcessors != nil {
		for _, processor := range bean.beanPreProcessors {
			processor.PreProcess(c)
		}
	}

	//2.初始化bean实例
	if bean.factoryMethod != nil {
		method := reflect.ValueOf(bean.factoryMethod)
		//实例化方法参数
		args := make([]reflect.Value, method.Type().NumIn())
		for i := 0; i < method.Type().NumIn(); i++ {
			in := method.Type().In(i)
			isPtr := in.Kind() == reflect.Ptr
			if isPtr {
				in = in.Elem()
			}
			if in.Kind() == reflect.Struct {
				if instance, err := c.GetBeanInstanceByStruct(reflect.New(in).Interface()); err != nil {
					panic(err)
				} else {
					args[i] = reflect.ValueOf(instance)
					if !isPtr {
						args[i] = args[i].Elem()
					}
				}
			} else {
				args[i] = reflect.ValueOf(c.instanceNonStructArgs(method.Type().In(i)))
			}
		}
		//调用工厂方法
		values := method.Call(args)
		instanceRv := values[0]
		//保存指针值
		bean.instance = instanceRv.Interface()
		if instanceRv.Kind() == reflect.Struct {
			bean.instance = &bean.instance
		}
	} else {
		rt := reflect.TypeOf(bean.value)
		bean.instance = reflect.New(rt).Interface()
	}
	//调用初始化方法
	if initializer, ok := bean.instance.(Initializer); ok {
		initializer.Init(c)
	}

	//3.调用后置处理器
	if bean.beanPostProcessors != nil {
		for _, processor := range bean.beanPostProcessors {
			processor.PostProcess(c, bean.instance)
		}
	}
	c.Logger.Printf("create bean:%s finished", bean.name)
	return bean.instance
}

func (c *Container) instanceNonStructArgs(rt reflect.Type) interface{} {
	switch rt.Kind() {
	case reflect.Ptr:
		return nil
	case reflect.String:
		return ""
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		return 0
	case reflect.Float64:
	case reflect.Float32:
		return 0.0
	case reflect.Bool:
		return false
	case reflect.Slice:
		return reflect.MakeSlice(rt, 0, 0).Interface()
	case reflect.Map:
		return reflect.MakeMap(rt).Interface()
	case reflect.Array:
		return reflect.New(rt).Elem().Interface()
	case reflect.Chan:
		return reflect.MakeChan(rt, 0).Interface()
	case reflect.Func:
		return reflect.MakeFunc(rt, func(args []reflect.Value) (results []reflect.Value) {
			return nil
		}).Interface()
	}
	return nil
}

// Bean 表示一个对象
type Bean struct {
	name               string
	priority           int
	value              interface{} //原始对象,struct
	instance           interface{} //创建完成后并赋值后的实例
	factoryMethod      interface{} //实例化工厂方法
	isSingleton        bool        //是否单例
	beanPreProcessors  []BeanPreProcessor
	beanPostProcessors []BeanPostProcessor
	lock               *sync.Mutex
}

func NewBean() *Bean {
	return &Bean{}
}

type BeanNameProvider interface {
	GetBeanName() string
}

func (bean *Bean) SetName(name string) {
	if name != "" {
		bean.name = name
	}
}

func (bean *Bean) SetPriority(priority int) {
	bean.priority = priority
}

func (bean *Bean) SetValue(value interface{}) {
	if value == nil {
		panic(NilError)
	}
	rt := reflect.TypeOf(value)
	if rt.Kind() == reflect.Ptr && rt.Elem().Kind() != reflect.Struct || rt.Kind() != reflect.Struct {
		panic(TypeNotMatchError)
	}
	bean.value = value
}

func (bean *Bean) SetFactoryMethod(method interface{}) {
	if method == nil {
		panic(NilError)
	}
	rt := reflect.TypeOf(method)
	if rt.NumIn() != 1 {
		panic(FactoryMethodReturnsError)
	} else {
		returnRt := rt.In(0)
		if returnRt.Kind() == reflect.Ptr && returnRt.Elem().Kind() != reflect.Struct || returnRt.Kind() != reflect.Struct {
			panic(TypeNotMatchError)
		}
	}
	bean.factoryMethod = method
}

func (bean *Bean) SetIsSingleton(isSingleton bool) {
	bean.isSingleton = isSingleton
}

func (bean *Bean) AddBeanPreProcessor(processor BeanPreProcessor) {
	bean.lock.Lock()
	defer bean.lock.Unlock()
	if processor != nil {
		if bean.beanPreProcessors == nil {
			bean.beanPreProcessors = []BeanPreProcessor{processor}
		} else {
			bean.beanPreProcessors = append(bean.beanPreProcessors, processor)
		}
	}
}

func (bean *Bean) AddBeanPostProcessor(processor BeanPostProcessor) {
	bean.lock.Lock()
	defer bean.lock.Unlock()
	if processor != nil {
		if bean.beanPostProcessors == nil {
			bean.beanPostProcessors = []BeanPostProcessor{processor}
		} else {
			bean.beanPostProcessors = append(bean.beanPostProcessors, processor)
		}
	}
}
