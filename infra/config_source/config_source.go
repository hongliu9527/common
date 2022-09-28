/*
 * @Author: hongliu
 * @Date: 2022-09-21 10:23:17
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-21 11:37:16
 * @FilePath: \common\infra\config_source\config_source.go
 * @Description: 配置源构造器
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package configsource

import (
	"context"
	"sync"

	"hongliu9527/common/infra/common"
	"hongliu9527/common/infra/config_source/apollo"
	"hongliu9527/common/infra/config_source/local"

	"github.com/hongliu9527/go-tools/logger"
)

var (
	// singleton 配置数据源单例对象
	singleton common.ConfigSource

	// 配置数据源类型和选项
	configSourceType   common.ConfigSourceType
	configSourceOption string
	configServiceName  string

	// 只执行一次
	once sync.Once
)

// SetConfigSource 设置"配置数据源"
func SetConfigSource(sourceType common.ConfigSourceType, sourceOption string, serviceName string) error {
	configSourceType = sourceType
	configSourceOption = sourceOption
	configServiceName = serviceName

	return nil
}

// New 创建配置数据源
func New() common.ConfigSource {
	singleton = nil

	// 根据配置数据源类型创建数据源并进行初始化
	switch configSourceType {
	case common.Apollo:
		singleton = apollo.New(configServiceName, configSourceOption)
	case common.Local:
		singleton = local.New(configServiceName, configSourceOption)
	default:
		logger.Error("创建配置数据源失败(不支持数据源类型：%s)，因为配置加载失败整个服务将无法正常运行", configSourceType)
	}

	once.Do(func() {
		singleton.Init(context.TODO())
	})

	return singleton
}
