/*
 * @Author: hongliu
 * @Date: 2022-09-16 14:17:23
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:47:08
 * @FilePath: \common\infra\common\config_source.go
 * @Description: 配置数据源接口抽象定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package common

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hongliu9527/common/utils"

	"github.com/hongliu9527/go-tools/logger"
)

// ConfigSourceType 配置数据源类型
type ConfigSourceType string

// 配置数据源相关定义
const (
	Apollo ConfigSourceType = "apollo" // Apollo配置中心
	K8s    ConfigSourceType = "k8s"    // K8s ConfigMap存储方法
	Local  ConfigSourceType = "local"  // 配置文件本地存储
)

// 结构体标签相关定义
const (
	decodeTag  = "mapstructure" // 配置源解码标签
	defaultTag = "default"      // 字段默认值标签
	omitTag    = "omit"         // 配置为omit的时候，该字段既不反序列化，也不参与字段校验
)

// ConfigSource 配置数据源抽象接口定义
type ConfigSource interface {
	// 初始化接口需要传入服务名称和配置数据源的访问端点信息
	Init(ctx context.Context) error // 初始化配置数据源

	// 所有数据源方案通过兼容"文件名称"这样的单层概念进行数据存储
	// 文件名称通过"."进行模块划分，从而实现服务级别的唯一性，例如：infra.oss.yaml、adapter.openapi.yaml
	// 即文件名称的命名结构具有了分层的功能，所有数据源方案只需要实现单层存储即可，例如如下场景：
	// Apollo可以将每个"文件名称"按照命名空间的方案进行存储和读取
	// K8s可以以"文件名称"来命名ConfigMap的名称进行存储和读取
	// 本地配置文件存储配置文件可以按照"文件名称"为名的形式存放在本地config目录下
	Read(filename string, value interface{}, timeout time.Duration) error   // 读取指定文件名称的配置数据【带超时控制】
	Listen(filename string, value interface{}, timeout time.Duration) error // 监听指定文件的配置数据更新，配置更新时立即访问最新数据，否则超市推出该监听接口
}

// DecodeConfig 反序列化配置信息，并返回详细的错误列表
func DecodeConfig(configValue, structPointer interface{}) error {
	decodeErrors := make([]error, 0)
	// 第一个参数类型必须时map[string]interface 第二个参数必须时结构体指针
	valueMap, ok := configValue.(map[string]interface{})
	if !ok {
		decodeErrors = append(decodeErrors, errors.New("第一个参数必须是map[string]interface{}类型"))
	}
	if !utils.IsStructPointer(structPointer) {
		decodeErrors = append(decodeErrors, errors.New("第二个参数必须是结构体指针"))
	}

	// 如果参数不正确，直接返回
	if len(decodeErrors) > 0 {
		return utils.MergeErrors(decodeErrors)
	}

	// 检查结构体的Tag是否正确，如果有错误，直接返回
	checkStructTag("", structPointer, &decodeErrors)
	if len(decodeErrors) > 0 {
		return utils.MergeErrors(decodeErrors)
	}

	// 配置源多余的配置项
	spareConfigKeys := make([]string, 0)
	// 结构体多余的字段
	spareStructTags := make([]string, 0)

	// 根据tag反序列化，并记录反序列化信息
	setStructDecodeValue("", valueMap, structPointer, &spareConfigKeys, &spareStructTags)
	// 如果有多余的字段没有在配置源配置，或者配置源有多余的配置项，则打印出来
	if len(spareConfigKeys) > 0 {
		logger.Warning("配置源中(%s)字段没有被使用，请检查代码实现是否落后于配置更新或者代码中已经删除的配置项没有及时清除", strings.Join(spareConfigKeys, ","))
	}
	if len(spareStructTags) > 0 {
		logger.Warning("结构体中(%s)标记没有配置，将使用代码级默认值", strings.Join(spareStructTags, ","))
	}
	return nil
}

// checkStructTag 检查结构体的Tag定义是否正确
func checkStructTag(path string, value interface{}, decodeErrors *[]error) {
	// 定义Tag表，用于检查在一个结构体内，是否会有相同的Tag
	tagMap := make(map[string]struct{})

	reflectType := reflect.TypeOf(value).Elem()
	reflectValue := reflect.ValueOf(value).Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		tagName := strings.ToLower(reflectType.Field(i).Tag.Get(decodeTag)) // 字段的标记名
		fieldValue := reflectValue.Field(i)
		fieldKind := fieldValue.Kind()

		// 如果该字段不需要校验
		if tagName == omitTag {
			continue
		}

		// 如果在结构体的同一层出现重复tag，则记录错误
		if _, ok := tagMap[tagName]; ok {
			*decodeErrors = append(*decodeErrors, fmt.Errorf("在结构体层(%s)中，标签(%s)已经存在", path, tagName))
			continue
		} else {
			tagMap[tagName] = struct{}{}
		}

		// 判断该字段是否可以寻址和赋值
		if !fieldValue.CanAddr() || !fieldValue.CanSet() {
			*decodeErrors = append(*decodeErrors, fmt.Errorf("在结构体层(%s)中，标签(%s)对应的字段无法寻址和赋值", path, tagName))
			continue
		}

		switch fieldKind {
		case reflect.Struct: // 如果该字段类型是结构体，则递归检查
			subPath := fmt.Sprintf("%s[%s]", path, tagName)
			fieldValue := fieldValue.Addr().Interface()
			checkStructTag(subPath, fieldValue, decodeErrors)
		case reflect.Slice: // 如果该字段是切片，则递归检查
			subPath := fmt.Sprintf("%s[%s]", path, tagName)
			fieldValue := fieldValue.Addr().Interface()
			checkSliceTag(subPath, fieldValue, decodeErrors)
		default:
			// 如果字段类型是基础类型，则检查该基础类型的默认值是否正确
			if utils.IsBasicType(fieldKind) {
				defaultValue, ok := reflectType.Field(i).Tag.Lookup(defaultTag)
				if !ok {
					*decodeErrors = append(*decodeErrors, fmt.Errorf("在结构体层(%s)中，标签(%s)对应的字段默认值缺失", path, tagName))
				} else {
					// 如果有默认值则进行类型检查，目前基础类型只支持string和int
					err := checkStructDefaultValue(fieldKind, defaultValue)
					if err != nil {
						*decodeErrors = append(*decodeErrors, fmt.Errorf("在结构体层(%s)中，标签(%s)对应的默认值不合法(%s)", path,
							tagName, err.Error()))
					}
				}
			} else { // 如果该字段既不是基础类型，也不是数组和结构体，则添加不支持类型报错
				*decodeErrors = append(*decodeErrors, fmt.Errorf("在结构体层(%s)中，标签(%s)对应的字段(%v)无法解析", path, tagName, fieldKind))
			}
		}
	}
}

// checkSliceTag 检查切片元素的tag是否定义正确
func checkSliceTag(path string, value interface{}, decodeErrors *[]error) {
	reflectType := reflect.TypeOf(value).Elem()
	elemType := reflectType.Elem()
	fieldName := elemType.Name()

	// 如果是基础类型，则检查是否有tag
	if utils.IsBasicType(elemType.Kind()) {
		return
	}
	// 如果是切片类型，暂时不支持多维切片，添加错误信息
	if elemType.Kind() == reflect.Slice {
		*decodeErrors = append(*decodeErrors, fmt.Errorf("在切片层(%s)中，(%s)字段对应的类型不能为切片类型", path, fieldName))
	}
	// 如果是结构体，递归检查
	path = fmt.Sprintf("%s[%s]", path, fieldName)
	elemValue := reflect.New(elemType).Interface()
	checkStructTag(path, elemValue, decodeErrors)
}

// checkStructDefaultValue 检查结构体默认值的数据类型是否正确
func checkStructDefaultValue(kind reflect.Kind, defaultValue string) error {
	switch kind {
	case reflect.Int:
		_, err := strconv.Atoi(defaultValue)
		return err
	case reflect.Bool:
		boolStr := strings.ToLower(defaultValue)
		if boolStr != "true" && boolStr != "false" {
			return fmt.Errorf("期望值: true或false-实际值: %s", boolStr)
		}
	case reflect.String:
		return nil
	default:
		return fmt.Errorf("不支持设置默认值的类型(%s)", kind)
	}

	return nil
}

// setStructDecodeValue 将配置信息根据tag设置到对应的结构体信息中
func setStructDecodeValue(path string, configMap map[string]interface{}, structPointer interface{}, spareConfigKeys *[]string,
	spareStructTags *[]string) {
	reflectValue := reflect.ValueOf(structPointer).Elem()
	reflectType := reflect.TypeOf(structPointer).Elem()

	// 递归遍历结构体的字段
	for i := 0; i < reflectType.NumField(); i++ {
		tagName := strings.ToLower(reflectType.Field(i).Tag.Get(decodeTag)) // 获取反序列化的标记名
		defaultValue := reflectType.Field(i).Tag.Get(defaultTag)            // 获取字段的默认值
		fieldValue := reflectValue.Field(i)

		// 如果不需要赋值
		if tagName == omitTag {
			continue
		}

		// 检查配置信息表中是否有对应的tag，如果没有，则添加到对应的结构体多余列表中，并设置该字符按的值为默认值
		mapValue, ok := configMap[tagName]
		if !ok {
			*spareStructTags = append(*spareStructTags, fmt.Sprintf("%s%s", path, tagName))
			setDefaultFieldValue(fieldValue, defaultValue)
		}

		// 如果结构体的字段类型是基础类型，则设置结构体的值
		if utils.IsBasicType(fieldValue.Kind()) {
			err := setFieldValue(fieldValue, mapValue)
			if err != nil { // 如果设置失败，打印详细信息，并把该字段的值设置为默认值
				logger.Warning("在结构体层(%s)中，设置(%s)对应字段值失败，错误信息为(%s)。该字段使用默认值(%s)", path, tagName, err.Error(), defaultValue)
				setDefaultFieldValue(fieldValue, defaultValue)
			}
			delete(configMap, tagName)
			continue
		}

		// 如果是结构体，就需要递归赋值，同时需要判断需要赋的值是否为map[string]interface{}
		if fieldValue.Kind() == reflect.Struct {
			subValue := fieldValue.Addr().Interface()
			if utils.IsMap(configMap[tagName]) {
				subPath := fmt.Sprintf("%s[%s]", path, tagName)
				subConfigMap, _ := configMap[tagName].(map[string]interface{})
				// 获取子结构体和子配置表信息，递归赋值
				setStructDecodeValue(subPath, subConfigMap, subValue, spareConfigKeys, spareStructTags)
			} else {
				// 打印提示信息，并使用结构体的默认配置项
				logger.Warning("在结构体层(%s)中，配置源不是map[string]interface{}类型而是(%#v)，标记为(%s)对应的子结构体使用默认值", path,
					reflect.TypeOf(configMap[tagName]).Kind(), tagName)
				setDefaultStructValue(subValue)
			}
			delete(configMap, tagName)
		} else { // 如果是切片
			subValue := fieldValue.Addr().Interface()
			if utils.IsSlice(configMap[tagName]) {
				subPath := fmt.Sprintf("%s[%s]", path, tagName)
				subConfigSlice, _ := configMap[tagName].([]interface{})
				setSliceDecodeValue(subPath, subConfigSlice, subValue, spareConfigKeys, spareStructTags)
				delete(configMap, tagName)
			} else {
				// 打印提示信息
				logger.Warning("在结构体层(%s)中，配置源不是[]interface{}类型而是(%#v)", path, reflect.TypeOf(configMap[tagName]).Kind())
			}
		}
	}
	if len(configMap) > 0 {
		for configKey := range configMap {
			*spareConfigKeys = append(*spareConfigKeys, fmt.Sprintf("%s[%s]", path, configKey))
		}
	}
}

// setDefaultFieldValue 设置字段默认值
func setDefaultFieldValue(fieldValue reflect.Value, defaultValue string) {
	// 由于在赋值之前已经做了默认值类型检查，所以这里可以不检查
	switch fieldValue.Kind() {
	case reflect.Int: // 如果是整型
		intVal, _ := strconv.ParseInt(defaultValue, 10, 64)
		fieldValue.SetInt(intVal)
	case reflect.Bool: // 如果是布尔型
		boolString := strings.ToLower(defaultValue)
		if boolString == "true" {
			fieldValue.SetBool(true)
		} else {
			fieldValue.SetBool(false)
		}
	case reflect.String: // 如果是字符串
		fieldValue.SetString(defaultValue)
	default: // 如果是自定义类型
		fieldValue.Set(reflect.ValueOf(defaultValue).Convert(fieldValue.Type()))
	}
}

// setFieldValue 设置结构体字段值，如果设置失败，则使用默认值
func setFieldValue(fieldValue reflect.Value, value interface{}) error {
	reflectValue := reflect.ValueOf(value)
	if fieldValue.Kind() != reflectValue.Kind() { // 设置失败，使用默认值
		return fmt.Errorf("该字段需要一个(%s)类型，但是提供的是(%s)类型", fieldValue.Kind(), reflectValue.Kind())
	}
	fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
	return nil
}

// setDefaultStructValue 设置整个结构体的字段都为默认值，由于之前已经对结构体的默认值和字段类型进行了校验，这里直接递归赋值即可
func setDefaultStructValue(value interface{}) {
	reflectValue := reflect.ValueOf(value)
	reflectType := reflect.TypeOf(value)

	// 递归遍历结构体字段并赋值
	for i := 0; i < reflectType.NumField(); i++ {
		tagName := reflectType.Field(i).Tag.Get(decodeTag)       // 字段序列化名称
		defaultValue := reflectType.Field(i).Tag.Get(defaultTag) // 字段的默认值
		fieldValue := reflectValue.Field(i)

		// 如果该字段不需要赋值
		if tagName == omitTag {
			continue
		}

		// 如果字段是基本类型，则赋默认值
		if utils.IsBasicType(reflectValue.Kind()) {
			setDefaultFieldValue(fieldValue, defaultValue)
		}

		// 如果是结构体，则递归赋值
		if fieldValue.Kind() == reflect.Struct {
			subValue := fieldValue.Addr().Interface()
			setDefaultStructValue(subValue)
		}
	}
}

// setSliceDecodeValue 将配置信息根据tag设置到对应的切片中
func setSliceDecodeValue(path string, subConfigSlice []interface{}, slicePointer interface{}, spareConfigKeys *[]string,
	spareStructTags *[]string) {
	reflectType := reflect.TypeOf(slicePointer).Elem()
	elem := reflectType.Elem()
	reflectValue := reflect.ValueOf(slicePointer)

	// 如果切片指针不能被设置
	if !reflectValue.CanSet() {
		reflectValue = reflectValue.Elem()
	}

	// 由于结构体切片在未初始化的时候，只是一个空指针，因此需要获取值并且初始化
	indirectSlice := reflect.Indirect(reflectValue)
	valueSlice := reflect.MakeSlice(indirectSlice.Type(), 0, 0)

	// 如果是基础类型，直接赋值
	if utils.IsBasicType(elem.Kind()) {
		// 首先遍历值
		for _, basicElem := range subConfigSlice {
			newElem := reflect.New(elem)
			setDefaultFieldValue(reflect.Indirect(newElem), basicElem.(string))
			valueSlice = reflect.Append(valueSlice, reflect.Indirect(newElem))
		}
		// 设置反射值
		reflectValue.Set(valueSlice)
		return
	}

	// 配置信息不支持多维切片，这里按照结构体处理
	for _, configElem := range subConfigSlice {
		newElem := reflect.New(elem)
		if utils.IsMap(configElem) {
			// 创建数据表
			configMapData, _ := configElem.(map[interface{}]interface{})
			configMap := make(map[string]interface{})
			for index, value := range configMapData {
				configMap[strings.ToLower(index.(string))] = value
			}
			setStructDecodeValue(path, configMap, newElem.Interface(), spareConfigKeys, spareStructTags)
		} else {
			// 打印提示信息，并且使用结构体的默认配置项
			logger.Warning("在结构体层(%s)中，配置源不是map[string]interface{}类型而是(%v)，标记为(%s)对应的子结构体使用默认值", path,
				reflect.TypeOf(configElem).Kind(), path)
			setDefaultStructValue(newElem.Interface())
		}

		// 这里生成的值是指针，需要取到值来添加
		valueSlice = reflect.Append(valueSlice, reflect.Indirect(newElem))
		// 设置反射值
		reflectValue.Set(valueSlice)
	}
}
