/*
 * @Author: hongliu
 * @Date: 2022-10-18 16:13:44
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:54:07
 * @FilePath: \common\infra\redis\config\redis.go
 * @Description:redis基础设施配置定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package config

import (
	"context"
	"sync"
	"time"

	"github.com/hongliu9527/common/infra/base"
	"github.com/hongliu9527/common/infra/common"

	"github.com/pkg/errors"
)

// 常量相关定义
const (
	RedisModuleName          = "Redis"            // 配置模块名
	RedisInfraConfigFileName = "infra.redis.yaml" // Redis基础设施配置文件名称
)

var (
	// redisInfraConfig Redis基础设施配置单例对象
	redisInfraConfig RedisInfraConfig

	// 只执行一次
	once sync.Once
)

// RedisInfraConfig Redis基础设施配置结构定义
type RedisInfraConfig struct {
	UseExternalHost  bool                  `mapstructure:"useExternalHost" default:"false"`            // 是否使用外网连接                                                             // 使用外网地址(默认为false)
	HostPort         string                `mapstructure:"hostPort" default:"127.0.0.1:6379" `         // Redis外网主机名称或访问地址和访问端口
	InternalHostPort string                `mapstructure:"internalHostPort" default:"127.0.0.1:6379" ` // Redis内网主机名称或访问地址和访问端口
	Password         string                `mapstructure:"password" default:"" `                       // 登录密码，默认为空
	DB               int                   `mapstructure:"db" default:"0" `                            // 数据库索引(从0开始),根据数据库个数逐个递增
	base.BaseConfig  `mapstructure:"omit"` // 基础配置信息
}

// New 创建Redis基础设施配置
func New(source common.ConfigSource, useExternalHost bool) (*RedisInfraConfig, error) {
	redisInfraConfig.BaseConfig = base.NewBaseConfig(RedisModuleName, RedisInfraConfigFileName)
	err := source.Read(RedisInfraConfigFileName, &redisInfraConfig, 20*time.Second)
	if err != nil {
		return nil, errors.WithMessage(err, "读取Redis基础设施配置信息失败")
	}
	// 调整终端命令配置参数优先级高于apollo远程配置
	redisInfraConfig.UseExternalHost = useExternalHost

	once.Do(func() {
		go redisInfraConfig.ListenSource(context.TODO(), source, &redisInfraConfig)
	})

	return &redisInfraConfig, nil
}
