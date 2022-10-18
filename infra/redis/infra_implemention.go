/*
 * @Author: hongliu
 * @Date: 2022-10-18 16:48:48
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-18 16:49:55
 * @FilePath: \common\infra\redis\infra_implemention.go
 * @Description:基础设施接口实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package redis_infra

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// 常量相关定义
const (
	RedisInfraName = "Redis" // Redis基础设施名称
)

// Name 查询基础设施名称
func (i *RedisInfra) Name() string {
	return RedisInfraName
}

// start 启动基础设施
func (i *RedisInfra) start(ctx context.Context) error {
	infraCtx, infraCancel := context.WithCancel(ctx)
	i.ctx = infraCtx
	i.cancel = infraCancel

	// 判断是否使用外网地址
	hostPort := i.config.InternalHostPort
	if i.config.UseExternalHost {
		hostPort = i.config.HostPort
	}

	client := redis.NewClient(&redis.Options{
		Addr:         hostPort,
		Password:     i.config.Password,
		DB:           i.config.DB,
		MinIdleConns: 5, // 最小空闲连接数
	})
	i.client = client

	// TODO: 增加限流设置
	// 参考：https://redis.uptrace.dev/guide/rate-limiting.html

	return nil
}

// stop 停止基础设施
func (i *RedisInfra) stop() error {
	err := i.client.Close()
	if err != nil {
		return errors.Wrap(err, "断开Redis哨兵失败")
	}

	return nil
}
