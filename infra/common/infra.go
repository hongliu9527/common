/*
 * @Author: hongliu
 * @Date: 2022-09-21 15:46:30
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-19 14:00:18
 * @FilePath: \common\infra\common\infra.go
 * @Description:基础设施接口定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package common

import "context"

// Infra 基础设施接口定义
type Infra interface {
	Start(ctx context.Context) error   // 启动基础设施
	Stop() error                       // 关闭基础设施
	Restart(ctx context.Context) error // 重启基础设施
	Name() string                      // 基础设施名称
}

// OrmInfra orm基础设施接口定义
type OrmInfra interface {
	Orm
	Infra
}

// OssInfra oss基础设施接口定义
type OssInfra interface {
	Oss
	Infra
}

// RedisInfra redis基础设施接口定义
type RedisInfra interface {
	Redis
	Infra
}
