package core

import (
	"fmt"
	"github.com/kgip/go-spring/configuration"
	errors "github.com/kgip/go-spring/error"
	"reflect"
	"strings"
)

const (
	injectTag = "autowired"  //true, false
	configTag = "autoconfig" //true, false

	beanNameTag = "name"

	configPrefixTag          = "prefix"
	configKeyTag             = "key"
	pathSplitChar            = "."
	configKeySplitChar       = " "
	configKeySubKeySplitChar = "="

	configKeySubKeyValue   = "value"
	configKeySubKeyDefault = "default"
)

var (
	configKeySubKeys = map[string]bool{configKeySubKeyValue: true, configKeySubKeyDefault: true}
	instanceHandlers = []InstanceHandler{
		&ConfigInstanceHandler{
			fieldHandler: &DefaultConfigFieldHandler{},
		},
		&DefaultInstanceHandler{},
	}
)

type InstanceHandler interface {
	IsSupport(instance interface{}) bool
	Handle(c *Container, instance interface{})
}

type ConfigInstanceHandler struct {
	fieldHandler ConfigFieldHandler
}

func (*ConfigInstanceHandler) IsSupport(instance interface{}) bool {
	if _, ok := instance.(configuration.Storage); ok {
		return true
	}
	rt := reflect.TypeOf(instance).Elem()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if _, ok := f.Tag.Lookup(configPrefixTag); ok {
			return true
		}
		if _, ok := f.Tag.Lookup(configKeyTag); ok {
			return true
		}
	}
	return false
}

func (handler *ConfigInstanceHandler) Handle(c *Container, instance interface{}) {
	var prefix string
	if store, ok := instance.(configuration.Storage); ok {
		prefix = store.ConfigurationPrefix()
	}
	rt := reflect.TypeOf(instance).Elem()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if autoconfig, ok := f.Tag.Lookup(configTag); ok && autoconfig == "false" {
			continue
		}
		handler.fieldHandler.Handle(c.GetConfiguration(), &f, prefix)
	}
}

type ConfigFieldHandler interface {
	Handle(c configuration.Provider, field *reflect.StructField, prefix string)
}

type DefaultConfigFieldHandler struct{}

func (handler *DefaultConfigFieldHandler) resolveConfigKeySubKey(keyValueStr string, index int) (key, value string) {
	key = keyValueStr[:index]
	if configKeySubKeys[key] {
		value = keyValueStr[index+1:]
	} else {
		panic(errors.UnknownConfigKeySubKeyError.Detail(fmt.Sprintf("unknown sub key '%s'", key)))
	}
	return
}

//configKey??????
//1.??????????????????key???????????????value:?????? `configKey:"path"`
//2.??????key??????????????????key???????????????key????????????????????? `configKey:"value=path default=10.4.68.144:3306"`
func (handler *DefaultConfigFieldHandler) resolveConfigKey(configKey string) map[string]string {
	splits := strings.Split(configKey, configKeySplitChar)
	var keyValues []string
	var keyValuesMap = map[string]string{}
	//???????????????
	for _, split := range splits {
		if split != " " {
			keyValues = append(keyValues, split)
		}
	}
	if len(keyValues) <= 0 {
	} else if len(keyValues) == 1 {
		if index := strings.Index(keyValues[0], configKeySubKeySplitChar); index > -1 {
			k, v := handler.resolveConfigKeySubKey(keyValues[0], index)
			keyValuesMap[k] = v
		} else {
			keyValuesMap[configKeySubKeyValue] = keyValues[0]
		}
	} else { //more than one sub key
		for _, keyValue := range keyValues {
			if index := strings.Index(keyValue, configKeySubKeySplitChar); index > -1 {
				k, v := handler.resolveConfigKeySubKey(keyValue, index)
				keyValuesMap[k] = v
			} else {
				panic(errors.ConfigKeySubKeyResolveError.Detail(fmt.Sprintf("error key '%s'", keyValue)))
			}
		}
	}
	return keyValuesMap
}

func (handler *DefaultConfigFieldHandler) Handle(c configuration.Provider, field *reflect.StructField, prefix string) {
	var configKey string
	var defaultValue string
	var configPrefix string
	if key, ok := field.Tag.Lookup(configKeyTag); ok {
		if key == "" {
			panic(errors.ConfigKeyError.Detail(fmt.Sprintf("config key of '%s' can't be empty", field.Name)))
		}
		keyMap := handler.resolveConfigKey(key)
		if keyMap[configKeySubKeyDefault] != "" {
			defaultValue = keyMap[configKeySubKeyDefault]
		}
		if keyMap[configKeySubKeyValue] != "" {
			configKey = prefix + pathSplitChar + keyMap[configKeySubKeyValue]
		} else {
			configKey = prefix + pathSplitChar + field.Name
		}
	} else if prefixTagValue, ok := field.Tag.Lookup(configPrefixTag); ok {
		if prefixTagValue != "" {
			configPrefix = prefix + pathSplitChar + prefixTagValue
		} else {
			configPrefix = prefix
		}
	} else {
		configKey = prefix + pathSplitChar + field.Name
	}
	fmt.Println(configPrefix, defaultValue)
	if configKey != "" {
		if value := c.GetConfig(configKey); value == nil {
			panic(errors.UnknownConfigKeyError.Detail(configKey))
		} else {

		}
	}
}

type DefaultInstanceHandler struct{}

func (*DefaultInstanceHandler) IsSupport(instance interface{}) bool {
	return true
}

func (*DefaultInstanceHandler) Handle(c *Container, instance interface{}) {
	rv := reflect.ValueOf(instance).Elem()
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		//????????????????????????????????????field?????????
		if inject, ok := f.Tag.Lookup(injectTag); ok && inject == "false" {
			continue
		}
		var fieldInstance interface{}
		var hasBeanNameTag bool
		var tagName string
		if tagName, hasBeanNameTag = f.Tag.Lookup(beanNameTag); hasBeanNameTag {
			if fieldInstance = c.GetBeanInstanceByName(tagName); fieldInstance == nil {
				panic(errors.UnknownBeanNameError.Detail(fmt.Sprintf("unknown bean name: %s", tagName)))
			}
		} else if !f.Anonymous {
			fieldInstance = c.GetBeanInstanceByName(f.Name)
		}
		if fieldInstance == nil {
			fieldInstance = c.GetInstance(rv.Field(i).Type())
		} else {
			t := f.Type
			fieldInstance = reflect.ValueOf(fieldInstance).Elem().Interface()
			for t.Kind() != reflect.Ptr {
				t = t.Elem()
				fieldInstance = &fieldInstance
			}
			if t.Kind() != reflect.Struct && hasBeanNameTag {
				panic(errors.TypeNotMatchError.Detail(fmt.Sprintf("field %s of instance %v is not a struct type", tagName, instance)))
			}
		}
		if fieldInstance != nil {
			rv.Set(reflect.ValueOf(fieldInstance))
		}
	}
}

// AssignBeanPostProcessor ???bean????????????????????????
type AssignBeanPostProcessor struct{}

func (*AssignBeanPostProcessor) PostProcess(c *Container, instance interface{}) {
	for _, handler := range instanceHandlers {
		if handler.IsSupport(instance) {
			handler.Handle(c, instance)
			return
		}
	}
}

func (*AssignBeanPostProcessor) GetPriority() int {
	return 9999
}
