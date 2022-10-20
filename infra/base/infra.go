/*
 * @Author: hongliu
 * @Date: 2022-09-21 15:26:52
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:53:26
 * @FilePath: \common\infra\base\infra.go
 * @Description: 基础设施基类实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package base

import (
	"context"
	"fmt"
	"time"

	"github.com/hongliu9527/common/infra/common"

	"github.com/hongliu9527/go-tools/logger"
)

// BaseInfra 基础设施基类
type BaseInfra struct {
	InfraName         string                          // 基础设施名称
	StartFunc         func(ctx context.Context) error // 基础设施启动方法
	StopFunc          func() error                    // 基础设施停止方法
	ConfigImpl        common.ConfigImpl               // 基础设施配置信息
	isListeningConfig bool                            // 是否已开启配置信息监听协程
	ctx               context.Context                 // 基础设施基类上下文对象
	cancel            context.CancelFunc              // 基础设施基类退出回调函数
}

// NewBaseInfra 创建新的基础设施基类
func NewBaseInfra(infraName string, configImpl common.ConfigImpl, startFunc func(ctx context.Context) error, stopFunc func() error) BaseInfra {
	return BaseInfra{
		InfraName:  infraName,
		ConfigImpl: configImpl,
		StartFunc:  startFunc,
		StopFunc:   stopFunc,
	}
}

// listenConfig 配置信息监听协程
func (i *BaseInfra) listenConfig(ctx context.Context) {
	// 是否已开启配置信息监听协程标志位置为true
	i.isListeningConfig = true

	// 基础设施基类上下文对象和退出回调函数
	i.ctx, i.cancel = context.WithCancel(ctx)

	logger.Info("开启基础设施(%s)配置信息监听协程...", i.InfraName)
	for {
		select {
		case <-i.ctx.Done():
			logger.Info("关闭基础设施(%s)配置信息监听协程...", i.InfraName)
			// 是否已开启配置信息监听协程标志位置为false
			i.isListeningConfig = false
			return
		case result := <-i.ConfigImpl.Listen(i.ctx, 10*time.Second):
			if result.Error != nil {
				if result.Error == common.ErrReceiveEventTimeout {
					continue
				}
				logger.Error("基础设施(%s)监听配置事件失败(%s)", i.InfraName, result.Error.Error())
				continue
			}

			// 根据事件消息类型进行处理
			var err error
			switch result.ConfigEvent.Type {
			case common.ConfigChanged:
				logger.Info("基础设施(%s)监听到一次配置信息变更事件，将会重启一次该基础设施", i.InfraName)
				err = i.Restart(i.ctx)
			default:
				logger.Warning("基础设施(%s)读取到未知事件(%s)，暂时不进行任何操作", i.InfraName, result.ConfigEvent.Type)
			}

			if err != nil {
				logger.Error("重启基础设施(%s)失败(%s)", i.InfraName, err.Error())
			}
		}
	}
}

// Start 启动基础设施
func (i *BaseInfra) Start(ctx context.Context) error {
	// 判断基础设施启动方法是否为空
	if i.StartFunc == nil {
		return fmt.Errorf("基础设施(%s)未注册启动方法", i.InfraName)
	}

	// 若基础设施的配置不为空,且之前没有开启过配置信息监听协程,则需要开启一次配置信息监听协程
	if i.ConfigImpl != nil && i.isListeningConfig == false {
		go i.listenConfig(ctx)
	}

	// 调用基础设施启动方法
	return i.StartFunc(ctx)
}

// Stop 停止基础设施
func (i *BaseInfra) Stop() error {
	// 判断基础设施停止方法是否为空
	if i.StopFunc == nil {
		return fmt.Errorf("基础设施(%s)未注册停止方法", i.InfraName)
	}

	// 调用基础设施停止方法
	return i.StopFunc()
}

// Restart 重启基础设施
func (i *BaseInfra) Restart(ctx context.Context) error {
	err := i.Stop()
	if err != nil {
		return fmt.Errorf("停止基础设施(%s)失败(%s)", i.InfraName, err.Error())
	}

	err = i.Start(ctx)
	if err != nil {
		return fmt.Errorf("启动基础设施(%s)失败(%s)", i.InfraName, err.Error())
	}

	return nil
}
