/*
 * @Author: hongliu
 * @Date: 2022-09-16 10:22:29
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:53:20
 * @FilePath: \common\infra\base\config.go
 * @Description: 基础配置
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package base

import (
	"context"
	"time"

	"github.com/hongliu9527/common/infra/common"

	"github.com/hongliu9527/go-tools/logger"
)

// BaseConfig 基础配置类
type BaseConfig struct {
	Name              string                  // 配置模块名称
	FileName          string                  // 配置源文件名称
	ConfigEventBuffer chan common.ConfigEvent // 配置文件通道
}

// eventBufferLength 事件缓存长度
const eventBufferLength = 64

// readSourceTimeout 读取配置源超时(秒)
const readSourceTimeout = 20

// NewBaseConfig 新建基本配置对象
func NewBaseConfig(name, fileName string) BaseConfig {
	return BaseConfig{
		Name:              name,
		FileName:          fileName,
		ConfigEventBuffer: make(chan common.ConfigEvent, eventBufferLength),
	}
}

// ListenSource 监听配置源变更信息
func (c *BaseConfig) ListenSource(ctx context.Context, source common.ConfigSource, configInstance interface{}) {
	logger.Info("(%s)配置开启数据源监听协程...", c.Name)
	for {
		select {
		case <-ctx.Done(): // 提前退出
			close(c.ConfigEventBuffer)
			logger.Info("(%s)配置关闭数据源监听协程...", c.Name)
		default:
			time.Sleep(10 * time.Millisecond)
		}

		// 读取事件源
		err := source.Listen(c.FileName, configInstance, 20*time.Second)
		if err != nil {
			if err == common.ErrReceiveEventTimeout {
				continue
			}
			logger.Error("(%s)配置读取数据源(%s)事件失败(%s)", c.Name, c.FileName, err)
			continue
		}

		// 重新读取配置信息
		err = source.Read(c.FileName, configInstance, readSourceTimeout*time.Second)
		if err != nil {
			logger.Error("读取(%s)配置信息失败(%s)", c.Name, err.Error())
			continue
		}

		event := common.ConfigEvent{
			Type: common.ConfigChanged,
		}

		c.ConfigEventBuffer <- event
	}
}

// Listen 外部调用的配置监听接口
func (c *BaseConfig) Listen(ctx context.Context, timeout time.Duration) <-chan common.ConfigEventListenResult {
	result := make(chan common.ConfigEventListenResult, 8)

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, timeout)
	defer timeoutCancel()

	for {
		select {
		case <-ctx.Done(): // 提前退出
			result <- common.ConfigEventListenResult{
				ConfigEvent: common.ConfigEvent{},
				Error:       common.ErrAdvanceExit,
			}
			return result
		case <-timeoutCtx.Done(): // 超时退出
			result <- common.ConfigEventListenResult{
				ConfigEvent: common.ConfigEvent{},
				Error:       common.ErrReceiveEventTimeout,
			}
			return result
		case event := <-c.ConfigEventBuffer:
			result <- common.ConfigEventListenResult{
				ConfigEvent: event,
				Error:       nil,
			}
			return result
		}
	}
}
