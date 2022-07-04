package core

import (
	"github.com/kgip/go-spring/configuration"
	errors "github.com/kgip/go-spring/error"
	"log"
	"reflect"
	"sort"
	"sync"
)

var (
	once = &sync.Once{}
)

func GetPriority(o interface{}) int {
	if priority, ok := o.(PriorityProvider); ok {
		return priority.GetPriority()
	}
	return 0
}

// Container ioc容器
type Container struct {
	beans                    map[string]*Bean
	configuration            configuration.Provider
	globalBeanPreProcessors  []BeanPreProcessor
	globalBeanPostProcessors []BeanPostProcessor
	containerPreProcessors   []ContainerPreProcessor
	containerPostProcessors  []ContainerPostProcessor
	isInited                 bool //是否已经初始化
	logger                   *log.Logger
	lock                     *sync.Mutex
	rv                       *reflect.Value
}

func NewContainer(configurationProvider configuration.Provider, logger *log.Logger) *Container {
	c := &Container{
		beans:         map[string]*Bean{},
		lock:          &sync.Mutex{},
		configuration: configurationProvider,
		logger:        logger}
	rv := reflect.ValueOf(c)
	c.rv = &rv
	return c
}

// Init 容器初始化方法
func (c *Container) Init() {
	once.Do(func() {
		c.isInited = true
		c.logger.Println("Ioc container start init....")
		if c.containerPreProcessors != nil {
			sort.Slice(c.containerPreProcessors, func(i, j int) bool {
				return GetPriority(c.containerPreProcessors[i]) > GetPriority(c.containerPreProcessors[j])
			})
			for _, processor := range c.containerPreProcessors {
				processor.PreProcess(c)
			}
		}
		//加载配置
		c.logger.Printf("Start loading the configuration")
		c.configuration.Load()
		c.logger.Println("Load configuration complete")
		//实例化单例bean
		for name := range c.beans {
			c.GetBeanInstanceByName(name)
		}
		c.logger.Println("Ioc container instance beans complete")
		c.isInited = true
		if c.containerPostProcessors != nil {
			sort.Slice(c.containerPostProcessors, func(i, j int) bool {
				return GetPriority(c.containerPostProcessors[i]) > GetPriority(c.containerPostProcessors)
			})
			for _, processor := range c.containerPostProcessors {
				processor.PostProcess(c)
			}
		}
		c.logger.Println("Ioc container init complete")
	})
}

// GetBeanInstanceByName 获取bean
func (c *Container) GetBeanInstanceByName(name string) interface{} {
	if bean := c.beans[name]; bean != nil {
		if bean.isCreating && bean.factoryMethod != nil {
			panic(errors.CircularReferenceError)
		}
		if bean.isSingleton && bean.instance != nil {
			return bean.instance
		}
		return c.instanceBean(bean)
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
		return nil, errors.TypeNotMatchError
	}
	return c.GetBeanInstanceByName(rt.Name()), nil
}

func (c *Container) GetConfiguration() configuration.Provider {
	return c.configuration
}

func (c *Container) SetConfiguration(provider configuration.Provider) {
	c.configuration = provider
}

func (c *Container) checkInited() {
	if c.isInited {
		panic(errors.ContainerUpdateError)
	}
}

// AddBean 添加bean
func (c *Container) AddBean(bean *Bean) bool {
	c.checkInited()
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.beans[bean.name] != nil && GetPriority(c.beans[bean.name]) > GetPriority(bean) {
		return false
	}
	if c.globalBeanPreProcessors != nil {
		if bean.beanPreProcessors == nil {
			bean.beanPreProcessors = c.globalBeanPreProcessors
		} else {
			bean.beanPreProcessors = append(c.globalBeanPreProcessors, bean.beanPreProcessors...)
			sort.Slice(bean.beanPreProcessors, func(i, j int) bool {
				return GetPriority(bean.beanPreProcessors[i]) > GetPriority(bean.beanPreProcessors[j])
			})
		}
	}
	if c.globalBeanPostProcessors != nil {
		if bean.beanPostProcessors == nil {
			bean.beanPostProcessors = c.globalBeanPostProcessors
		} else {
			bean.beanPostProcessors = append(c.globalBeanPostProcessors, bean.beanPostProcessors...)
			sort.Slice(bean.beanPostProcessors, func(i, j int) bool {
				return GetPriority(bean.beanPostProcessors[i]) > GetPriority(bean.beanPostProcessors[j])
			})
		}
	}
	c.beans[bean.name] = bean
	return true
}

func (c *Container) AddBeanPreProcessor(processor BeanPreProcessor) {
	c.checkInited()
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
	c.checkInited()
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
	c.checkInited()
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
	c.checkInited()
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

// instanceBean 实例化bean
func (c *Container) instanceBean(bean *Bean) interface{} {
	c.logger.Printf("start creating bean:%s", bean.name)
	bean.isCreating = true
	//调用前置处理器
	if bean.beanPreProcessors != nil {
		for _, processor := range bean.beanPreProcessors {
			processor.PreProcess(c, bean)
		}
	}

	//初始化bean实例
	if bean.factoryMethod != nil {
		method := reflect.ValueOf(bean.factoryMethod)
		//实例化方法参数
		args := make([]reflect.Value, method.Type().NumIn())
		for i := 0; i < method.Type().NumIn(); i++ {
			in := method.Type().In(i)
			//如果接收容器指针作为参数，则将参数设置为容器指针
			if in.Kind() == reflect.Ptr && c.rv.Type().AssignableTo(in) {
				args[i] = *c.rv
			} else {
				args[i] = reflect.ValueOf(c.GetInstance(in))
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
		rt := reflect.TypeOf(bean.model).Elem()
		bean.instance = reflect.New(rt).Interface()
	}
	//调用初始化方法
	if initializer, ok := bean.instance.(Initializer); ok {
		initializer.Init(c)
	}

	//调用后置处理器
	if bean.beanPostProcessors != nil {
		for _, processor := range bean.beanPostProcessors {
			processor.PostProcess(c, bean.instance)
		}
	}
	c.logger.Printf("create bean:%s complete", bean.name)
	bean.isCreating = false
	return bean.instance
}

func (c *Container) GetInstance(rt reflect.Type) interface{} {
	if rt.Kind() == reflect.Ptr {
		instance := c.GetInstance(rt.Elem())
		if instance == nil {
			return nil
		}
		return &instance
	} else {
		return c.instanceByType(rt)
	}
}

func (c *Container) instanceByType(rt reflect.Type) interface{} {
	switch rt.Kind() {
	case reflect.Struct:
		if instance, err := c.GetBeanInstanceByStruct(reflect.New(rt).Interface()); err != nil {
			panic(err)
		} else if instance != nil {
			return reflect.ValueOf(instance).Elem().Interface()
		} else { //容器中不存在则创建一个默认对象
			return reflect.New(rt).Elem().Interface()
		}
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
	case reflect.Array:
		return reflect.New(rt).Elem().Interface()
	}
	return nil
}

// Bean 表示一个对象
type Bean struct {
	name               string
	isCreating         bool //是否正在被创建
	priority           int
	model              interface{} //原始对象,struct指针
	instance           interface{} //创建完成后并赋值后的实例指针
	factoryMethod      interface{} //实例化工厂方法
	isSingleton        bool        //是否单例
	beanPreProcessors  []BeanPreProcessor
	beanPostProcessors []BeanPostProcessor
	lock               *sync.Mutex
}

func NewBean(model interface{}) *Bean {
	bean := &Bean{lock: &sync.Mutex{}, isSingleton: true}
	bean.SetModel(model)
	if provider, ok := bean.model.(BeanNameProvider); ok {
		bean.SetName(provider.GetBeanName())
	}
	if bean.name == "" {
		bean.SetName(reflect.TypeOf(bean.model).Elem().Name())
	}
	return bean
}

func NewFactoryBean(factoryMethod interface{}) *Bean {
	bean := &Bean{lock: &sync.Mutex{}, isSingleton: true}
	bean.SetFactoryMethod(factoryMethod)
	if provider, ok := bean.factoryMethod.(BeanNameProvider); ok {
		bean.SetName(provider.GetBeanName())
	}
	if bean.name == "" {
		rt := reflect.TypeOf(bean.factoryMethod).In(0)
		if rt.Kind() == reflect.Ptr {
			bean.SetName(rt.Elem().Name())
		} else {
			bean.SetName(rt.Name())
		}
	}
	return bean
}

func (bean *Bean) SetName(name string) *Bean {
	if name != "" {
		bean.name = name
	} else {
		panic(errors.NameEmptyError)
	}
	return bean
}

func (bean *Bean) SetPriority(priority int) *Bean {
	bean.priority = priority
	return bean
}

func (bean *Bean) SetModel(model interface{}) *Bean {
	if model == nil {
		panic(errors.NilError)
	}
	rt := reflect.TypeOf(model)
	if rt.Kind() == reflect.Ptr && rt.Elem().Kind() != reflect.Struct || rt.Kind() != reflect.Struct {
		panic(errors.TypeNotMatchError)
	}
	if rt.Kind() == reflect.Struct {
		bean.model = &model
	} else {
		bean.model = model
	}
	return bean
}

func (bean *Bean) SetFactoryMethod(method interface{}) *Bean {
	if method == nil {
		panic(errors.NilError)
	}
	rt := reflect.TypeOf(method)
	if rt.Kind() != reflect.Func {
		panic(errors.TypeNotMatchError)
	}
	if rt.NumIn() != 1 {
		panic(errors.FactoryMethodReturnsError)
	} else {
		returnRt := rt.In(0)
		if returnRt.Kind() == reflect.Ptr && returnRt.Elem().Kind() != reflect.Struct || returnRt.Kind() != reflect.Struct {
			panic(errors.TypeNotMatchError)
		}
	}
	bean.factoryMethod = method
	return bean
}

func (bean *Bean) SetIsSingleton(isSingleton bool) *Bean {
	bean.isSingleton = isSingleton
	return bean
}

func (bean *Bean) AddBeanPreProcessor(processor BeanPreProcessor) *Bean {
	bean.lock.Lock()
	defer bean.lock.Unlock()
	if processor != nil {
		if bean.beanPreProcessors == nil {
			bean.beanPreProcessors = []BeanPreProcessor{processor}
		} else {
			bean.beanPreProcessors = append(bean.beanPreProcessors, processor)
		}
	}
	return bean
}

func (bean *Bean) AddBeanPostProcessor(processor BeanPostProcessor) *Bean {
	bean.lock.Lock()
	defer bean.lock.Unlock()
	if processor != nil {
		if bean.beanPostProcessors == nil {
			bean.beanPostProcessors = []BeanPostProcessor{processor}
		} else {
			bean.beanPostProcessors = append(bean.beanPostProcessors, processor)
		}
	}
	return bean
}

func (bean *Bean) GetName() string {
	return bean.name
}

func (bean *Bean) GetModel() interface{} {
	return bean.model
}

func (bean *Bean) GetFactoryMethod() interface{} {
	return bean.factoryMethod
}

func (bean *Bean) GetPriority() int {
	return bean.priority
}
