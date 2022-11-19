/*
 * @Author: hongliu
 * @Date: 2022-11-19 15:55:10
 * @LastEditors: hongliu
 * @LastEditTime: 2022-11-19 15:55:20
 * @FilePath: \common\infra\common\job.go
 * @Description: 定时任务基础设施接口定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package common

type Job interface {
	AddTask(name string, expr string, job func() error) error // 增加定时任务
	RemoveTask(name string) error                             // 移除定时任务
}
