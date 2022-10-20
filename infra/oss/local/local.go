/*
 * @Author: hongliu
 * @Date: 2022-10-17 15:54:33
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:53:37
 * @FilePath: \common\infra\oss\local\local.go
 * @Description:本地存储oss实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package local

import (
	"github.com/hongliu9527/common/infra/base"
	"github.com/hongliu9527/common/infra/common"
	"github.com/hongliu9527/common/infra/oss/config"
)

// 本地文件oss单例
var singleton localOssInfra

// localOssInfra 本地Oss基础设施类型定义
type localOssInfra struct {
	base.BaseInfra                        // 基础设施基类
	config         *config.OssInfraConfig // Oss配置信息
}

// New 创建Oss基础设施
func New(config *config.OssInfraConfig) common.OssInfra {
	singleton.config = config

	// 构建基础设施基类
	singleton.BaseInfra = base.NewBaseInfra(singleton.Name(), config, singleton.start, singleton.stop)

	return &singleton
}
