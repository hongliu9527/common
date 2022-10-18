/*
 * @Author: hongliu
 * @Date: 2022-09-24 16:26:04
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-17 15:03:43
 * @FilePath: \common\infra\oss\config\oss.go
 * @Description: oss 配置信息
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package config

import (
	"context"
	"sync"
	"time"

	"hongliu9527/common/infra/base"
	"hongliu9527/common/infra/common"

	"github.com/hongliu9527/go-tools/uuid"
	"github.com/pkg/errors"
)

// OssServiceVendor Oss运营商类型定义
type OssServiceVendor string

const (
	// Oss运营商相关定义
	AliYun      OssServiceVendor = "aliyun"      // 阿里云
	HuaWeiCloud OssServiceVendor = "huaweicloud" // 华为云
	Ctyun       OssServiceVendor = "ctyun"       // 天翼云
	Local       OssServiceVendor = "local"       // 本地oss

	// iot oos 目录结构
	Root            = "iot"            // oss桶下的一级目录
	SystemDirectory = Root + "/system" // iot下的二级目录, 存放系统资源文件
	DeviceDirectory = Root + "/device" // iot下的二级目录, 存放设备产生的媒体文件

	// 配置文件相关
	OssModuleName          = "Oss"            // 配置模块名
	OssInfraConfigFileName = "infra.oss.yaml" // Oss基础设施配置文件名称
)

var (
	// singleton Oss基础设施配置单例对象
	singleton OssInfraConfig

	// 只执行一次
	once sync.Once
)

// SystemFileDirectory 获取系统资源文件对应的oss目录
func SystemFileDirectory() string {
	return SystemDirectory + "/" + uuid.UUID() + "/"
}

// DeviceFileDirectory 获取媒体文件对应的oss目录
func DeviceFileDirectory() string {
	return DeviceDirectory + "/" + time.Now().Format("20060102") + "/"
}

// OssInfraConfig Oss基础设施配置结构定义
type OssInfraConfig struct {
	ServiceVendor OssServiceVendor `mapstructure:"serviceVendor" default:"aliyun"` // Oss运营商
	Compress      bool             `mapstructure:"compress" default:"true"`        // 上传文件是否压缩
	RetryCount    int              `mapstructure:"retryCount" default:"3"`         // 上传失败重试次数
	// Oss访问地址和密钥相关配置
	AccessKeyID          string                `mapstructure:"accessKeyID" default:"ID"`             // 数据访问KEY标识
	AccessKeySecret      string                `mapstructure:"accessKeySecret" default:"keySecret"`  // 数据访问密钥
	Endpoint             string                `mapstructure:"endpoint" default:"127.0.0.1"`         // 数据挂载点名称
	Bucket               string                `mapstructure:"bucket" default:"bucket-02"`           // 数据仓库名称
	RoleARN              string                `mapstructure:"roleARN" default:"testRole"`           // 临时角色访问ARN
	SignatureExpiresTime string                `mapstructure:"signatureExpiresTime" default:"10800"` // 签名过期时间
	base.BaseConfig      `mapstructure:"omit"` // 基础配置信息
}

// New 创建Oss基础设施配置
func New(source common.ConfigSource) (*OssInfraConfig, error) {
	singleton.BaseConfig = base.NewBaseConfig(OssModuleName, OssInfraConfigFileName)

	err := source.Read(OssInfraConfigFileName, &singleton, 20*time.Second)
	if err != nil {
		return nil, errors.WithMessage(err, "读取Oss基础设施配置信息失败")
	}

	once.Do(func() {
		go singleton.ListenSource(context.TODO(), source, &singleton)
	})

	return &singleton, nil
}
