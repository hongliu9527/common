/*
 * @Author: hongliu
 * @Date: 2022-11-19 16:07:47
 * @LastEditors: hongliu
 * @LastEditTime: 2022-11-19 16:07:57
 * @FilePath: \common\infra\job\job.go
 * @Description: 定时任务实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package job

import (
	"context"
	"sync"

	"github.com/hongliu9527/common/infra/base"
	"github.com/hongliu9527/common/infra/common"
	"github.com/robfig/cron/v3"
)

// DistributedJobInfra 分布式定时任务基础设施
type DistributedJobInfra struct {
	base.BaseInfra                         // 基础设施基类
	redis          common.RedisInfra       // redis缓存基础设施
	taskEntryMap   map[string]cron.EntryID // 任务列表
	taskLock       sync.RWMutex            // 任务列表读写锁
	cron           *cron.Cron              // 定时任务
	ctx            context.Context         // 上下文对象
	cancel         context.CancelFunc      // 取消回调函数
}

// StandaloneJobInfra 单机定时任务基础设施
type StandaloneJobInfra struct {
	base.BaseInfra                         // 基础设施基类
	taskEntryMap   map[string]cron.EntryID // 任务列表
	taskLock       sync.RWMutex            // 任务列表读写锁
	cron           *cron.Cron              // 定时任务
	ctx            context.Context         // 上下文对象
	cancel         context.CancelFunc      // 取消回调函数
}

// job基础设施单例
var (
	distributedJobInfra DistributedJobInfra
	standaloneJobInfra  StandaloneJobInfra
)

// NewDistributedJobInfra 创建分布式job基础设施
func NewDistributedJobInfra(redisInfra common.RedisInfra) common.JobInfra {
	distributedJobInfra.taskEntryMap = make(map[string]cron.EntryID)

	if distributedJobInfra.cron == nil {
		distributedJobInfra.cron = cron.New()
	}
	distributedJobInfra.redis = redisInfra
	return &distributedJobInfra
}

// NewStandaloneJobInfra 创建单机式job基础设施
func NewStandaloneJobInfra() common.JobInfra {
	standaloneJobInfra.taskEntryMap = make(map[string]cron.EntryID)

	if standaloneJobInfra.cron == nil {
		standaloneJobInfra.cron = cron.New()
	}
	return &standaloneJobInfra
}
