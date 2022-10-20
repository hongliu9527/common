/*
 * @Author: hongliu
 * @Date: 2022-10-17 11:07:18
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:53:29
 * @FilePath: \common\infra\oss\aliyun\aliyun.go
 * @Description: 阿里云oss实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package aliyun

import (
	"github.com/hongliu9527/common/infra/base"
	"github.com/hongliu9527/common/infra/common"
	"github.com/hongliu9527/common/infra/oss/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 阿里云oss单例
var singleton aliyunOssInfra

// aliyunOssInfra 阿里云Oss基础设施类型定义
type aliyunOssInfra struct {
	base.BaseInfra                        // 基础设施基类
	config         *config.OssInfraConfig // Oss配置信息
	bucket         *oss.Bucket            // 数据桶实例
}

// New 创建Oss基础设施
func New(config *config.OssInfraConfig) common.OssInfra {
	singleton.config = config

	// 构建基础设施基类
	singleton.BaseInfra = base.NewBaseInfra(singleton.Name(), config, singleton.start, singleton.stop)

	return &singleton
}
