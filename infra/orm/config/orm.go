/*
 * @Author: hongliu
 * @Date: 2022-09-21 16:01:47
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:52:54
 * @FilePath: \common\infra\orm\config\orm.go
 * @Description:Orm基础设施配置格式定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package config

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/hongliu9527/common/infra/base"
	"github.com/hongliu9527/common/infra/common"

	"github.com/pkg/errors"
)

// 常量相关定义
const (
	OrmModuleName          = "orm"            // 配置模块名
	OrmInfraConfigFileName = "infra.orm.yaml" // Orm基础设施配置文件名称
)

// 为了方便调试,还是需要支持通过程序配置日志等级和orm访问方式
var (
	// OrmLogLevel Orm基础设施日志等级
	OrmLogLevel string

	// useExternal 使用外网地址,用于本地调试
	useExternal bool

	// 只执行一次
	once sync.Once

	// singleton Orm基础设施配置单例对象
	singleton OrmInfraConfig
)

// OrmInfraConfig Orm基础设施配置结构定义
type OrmInfraConfig struct {
	Configs         []DataBaseConfig               `mapstructure:"configList"` // 数据库基础设施配置列表
	LogLevel        string                         `mapstructure:"omit"`       // orm基础设施日志等级
	UseExternalHost bool                           `mapstructure:"omit"`       // 使用外网地址(默认为false)
	base.BaseConfig `mapstructure:"omit" yaml:"-"` // 基础配置信息
}

// DataBaseConfig 数据库配置结构定义
type DataBaseConfig struct {
	Name             string `mapstructure:"name" default:"iotplatform.mysql"`          // 配置信息名称，用于区分不同的数据库实例
	Type             string `mapstructure:"type" default:"mysql"`                      // 数据库类型
	HostPort         string `mapstructure:"hostPort" default:"127.0.0.1:3306"`         // 数据库外网主机名称或访问地址和访问端口，例如：127.0.0.1:3306
	InternalHostPort string `mapstructure:"internalHostPort" default:"127.0.0.1:3306"` // 数据库内网主机名称或访问地址和访问端口，例如：127.0.0.1:3306
	DatabaseName     string `mapstructure:"databaseName" default:"my-blog"`            // 数据库名称
	Username         string `mapstructure:"username" default:"main"`                   // 数据库访问用户名
	Password         string `mapstructure:"password" default:"hongliu-2016"`           // 数据库访问密码
	TablePrefix      string `mapstructure:"tablePrefix" default:"blog_"`               // 表名前缀
	ConnectTimeout   int    `mapstructure:"connectTimeout" default:"10"`               // 连接超时时间，单位(秒)
}

// New 创建Orm基础设施配置
func New(source common.ConfigSource, logLevel string, useExternalHost bool) (*OrmInfraConfig, error) {

	singleton.BaseConfig = base.NewBaseConfig(OrmModuleName, OrmInfraConfigFileName)

	err := source.Read(OrmInfraConfigFileName, &singleton, 20*time.Second)
	if err != nil {
		return nil, errors.WithMessage(err, "读取Orm基础设施配置信息失败")
	}

	singleton.LogLevel = logLevel
	singleton.UseExternalHost = useExternalHost

	if useExternal {
		singleton.UseExternalHost = true
	}

	if len(OrmLogLevel) != 0 {
		singleton.LogLevel = OrmLogLevel
	}

	// 开启数据源监听协程
	once.Do(func() {
		go singleton.ListenSource(context.TODO(), source, &singleton)
	})

	return &singleton, nil
}

// SetCommamdConfig 设置有命令行传入的配置信息
func SetCommamdConfig(logLevel string, external string) {
	OrmLogLevel = logLevel
	useExternal = false

	outerList := strings.Split(external, ",")
	for _, infraName := range outerList {
		if infraName == "orm" {
			useExternal = true
		}
	}
}
