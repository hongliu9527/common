/*
 * @Author: hongliu
 * @Date: 2022-09-16 10:25:38
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-16 14:16:44
 * @FilePath: \common\infra\common\config.go
 * @Description: 配置模块通用数据定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package common

import (
	"context"
	"time"
)

// ConfigEventTye 配置事件类型
type ConfigEventTye string

// 配置事件类型相关定义
const (
	ConfigChanged ConfigEventTye = "configChanged" // 配置发生变化
)

// ConfigImpl 配置信息接口定义
type ConfigImpl interface {
	Listen(context.Context, time.Duration) <-chan ConfigEventListenResult
}

// ConfigEventListenResult 配置事件监听结果
type ConfigEventListenResult struct {
	ConfigEvent ConfigEvent // 配置事件
	Error       error       // 错误
}

// ConfigEvent 配置事件定义
type ConfigEvent struct {
	Type ConfigEventTye // 事件类型
}
