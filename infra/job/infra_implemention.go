/*
 * @Author: hongliu
 * @Date: 2022-11-19 16:05:07
 * @LastEditors: hongliu
 * @LastEditTime: 2022-11-19 16:19:40
 * @FilePath: \common\infra\job\infra_implemention.go
 * @Description: 基础设施接口实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package job

import "context"

// 常量相关定义
const (
	distributedJobInfraName string = "DistributedJob" // 分布式job基础设施名称
	standaloneJobInfraName  string = "StandaloneJob"  // 单机式job基础设施名称
)

// Name 基础设施名称
func (i *DistributedJobInfra) Name() string {
	return distributedJobInfraName
}

// start 基础设施启动
func (i *DistributedJobInfra) Start(ctx context.Context) error {
	// 创建基础设施上下文对象与退出回调函数
	i.ctx, i.cancel = context.WithCancel(ctx)
	if i.cron != nil {
		i.cron.Start()
	}
	return nil
}

// stop 停止基础设施
func (i *DistributedJobInfra) Stop() error {
	if i.cron != nil {
		i.cron.Stop()
	}
	return nil
}

// Name 基础设施名称
func (i *StandaloneJobInfra) Name() string {
	return standaloneJobInfraName
}

// start 基础设施启动
func (i *StandaloneJobInfra) Start(ctx context.Context) error {
	// 创建基础设施上下文对象与退出回调函数
	i.ctx, i.cancel = context.WithCancel(ctx)
	if i.cron != nil {
		i.cron.Start()
	}
	return nil
}

// stop 停止基础设施
func (i *StandaloneJobInfra) Stop() error {
	if i.cron != nil {
		i.cron.Stop()
	}
	return nil
}
