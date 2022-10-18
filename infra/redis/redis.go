/*
 * @Author: hongliu
 * @Date: 2022-10-18 16:20:51
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-18 16:45:12
 * @FilePath: \common\infra\redis\redis.go
 * @Description:Redis基础设施定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package redis_infra

import (
	"context"

	"hongliu9527/common/infra/base"
	"hongliu9527/common/infra/common"
	"hongliu9527/common/infra/redis/config"

	"github.com/go-redis/redis/v8"
)

// RedisInfra Redis基础设施类型定义
type RedisInfra struct {
	base.BaseInfra                          // 基础设施基类
	client         *redis.Client            // Redis客户端实例
	config         *config.RedisInfraConfig // 配置信息
	ctx            context.Context          // 上下文对象
	cancel         context.CancelFunc       // 退出回调函数
}

// singleton Redis基础设施单例对象
var singleton RedisInfra

// New 创建Redis基础设施
func New(config *config.RedisInfraConfig) common.RedisInfra {
	singleton.config = config
	singleton.BaseInfra = base.NewBaseInfra(singleton.Name(), nil, singleton.start, singleton.stop)
	singleton.client = redis.NewClient(
		&redis.Options{
			Addr:     config.HostPort,
			Password: config.Password,
			DB:       config.DB,
		},
	)

	return &singleton
}
