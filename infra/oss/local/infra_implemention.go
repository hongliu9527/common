/*
 * @Author: hongliu
 * @Date: 2022-10-17 17:26:08
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-17 17:26:13
 * @FilePath: \common\infra\oss\local\infra_implemention.go
 * @Description:基础设施接口的实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package local

import "context"

// 常量相关定义
const (
	LocalOssInfraName string = "LocalOss" // 本地Oss基础设施名称
)

// Name 获取基础设施名
func (i *localOssInfra) Name() string {
	return LocalOssInfraName
}

// Start 启动基础设施
func (i *localOssInfra) start(ctx context.Context) error {
	return nil
}

// Stop 停止基础设施
func (i *localOssInfra) stop() error {
	return nil
}
