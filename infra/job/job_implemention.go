/*
 * @Author: hongliu
 * @Date: 2022-11-19 16:06:11
 * @LastEditors: hongliu
 * @LastEditTime: 2022-11-19 16:18:27
 * @FilePath: \common\infra\job\job_implemention.go
 * @Description: 定时任务接口实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package job

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/hongliu9527/go-tools/logger"
)

// taskBarrierTime 同名任务多实例之间的屏障时间(等价于应用实例之间最大时差)
const taskBarrierTime = 5 * time.Second

// encodeName2Key 编码定时任务名称成键，防止中文内耗更大并且出现乱码(https://segmentfault.com/q/1010000011577694)
func encodeName2Key(name string) string {
	return fmt.Sprintf("%X", (md5.Sum([]byte(name))))
}

// AddTask 增加分布式定时任务
func (i *DistributedJobInfra) AddTask(name string, expr string, job func() error) error {
	i.taskLock.RLock()
	_, ok := i.taskEntryMap[name]
	i.taskLock.RUnlock()
	if ok {
		return fmt.Errorf("已存在同名任务(%s)，定时任务添加失败", name)
	}

	taskKey := encodeName2Key(name)
	// 定时任务保证多实例执行一次
	onceJob := func() {
		ok, err := i.redis.SetNX(taskKey, 1, taskBarrierTime)
		if err != nil {
			logger.Error("执行定时任务(%s)失败(%s)", name, err.Error())
			return
		}
		if ok {
			logger.Info("开始执行定时任务(%s)...", name)
			err := job()
			if err != nil {
				logger.Error("执行定时任务(%s)失败(%s)", name, err.Error())
			} else {
				logger.Info("定时任务(%s)执行完成", name)
			}
		}
	}

	entryID, err := i.cron.AddFunc(expr, onceJob)
	if err != nil {
		return err
	}

	i.taskLock.Lock()
	i.taskEntryMap[name] = entryID
	i.taskLock.Unlock()

	i.cron.Start()
	return nil
}

// RemoveTask 移除分布式定时任务
func (i *DistributedJobInfra) RemoveTask(name string) error {
	i.taskLock.RLock()
	entryID, ok := i.taskEntryMap[name]
	i.taskLock.RUnlock()
	if !ok {
		return fmt.Errorf("不存在该定时任务(%s)，定时任务移除失败", name)
	}

	i.taskLock.Lock()
	delete(i.taskEntryMap, name)
	i.taskLock.Unlock()

	i.cron.Remove(entryID)
	return nil
}

// AddTask 增加单机式定时任务
func (i *StandaloneJobInfra) AddTask(name string, expr string, job func() error) error {
	i.taskLock.RLock()
	_, ok := i.taskEntryMap[name]
	i.taskLock.RUnlock()
	if ok {
		return fmt.Errorf("已存在同名任务(%s)，定时任务添加失败", name)
	}

	execJob := func() {
		logger.Info("开始执行定时任务(%s)...", name)
		err := job()
		if err != nil {
			logger.Error("执行定时任务(%s)失败(%s)", name, err.Error())
		} else {
			logger.Info("定时任务(%s)执行完成", name)
		}
	}

	entryID, err := i.cron.AddFunc(expr, execJob)
	if err != nil {
		return err
	}

	i.taskLock.Lock()
	i.taskEntryMap[name] = entryID
	i.taskLock.Unlock()

	i.cron.Start()
	return nil
}

// RemoveTask 移除单机式定时任务
func (i *StandaloneJobInfra) RemoveTask(name string) error {
	i.taskLock.RLock()
	entryID, ok := i.taskEntryMap[name]
	i.taskLock.RUnlock()
	if !ok {
		return fmt.Errorf("不存在该定时任务(%s)，定时任务移除失败", name)
	}

	i.taskLock.Lock()
	delete(i.taskEntryMap, name)
	i.taskLock.Unlock()

	i.cron.Remove(entryID)
	return nil
}
