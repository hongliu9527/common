/*
 * @Author: hongliu
 * @Date: 2022-09-24 16:25:22
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-18 16:46:44
 * @FilePath: \common\infra\oss\oss.go
 * @Description: oss 构造器
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package oss

import (
	"hongliu9527/common/infra/common"
	"hongliu9527/common/infra/oss/aliyun"
	"hongliu9527/common/infra/oss/config"
	"hongliu9527/common/infra/oss/local"

	"github.com/hongliu9527/go-tools/logger"
)

// ossInfra Oss基础设施单例对象
var ossInfra common.OssInfra

// New 创建Oss基础设施
func New(ossConfig *config.OssInfraConfig) common.OssInfra {
	switch ossConfig.ServiceVendor {
	case config.AliYun:
		ossInfra = aliyun.New(ossConfig)
	// case config.Ctyun:
	// 	ossInfra = ctyun.New(ossConfig)
	case config.Local:
		ossInfra = local.New(ossConfig)
	default:
		logger.Error("创建Oss基础设施失败(不支持的Oss厂商类型：%s)", ossConfig.ServiceVendor)
	}

	return ossInfra
}
