/*
 * @Author: hongliu
 * @Date: 2022-09-20 16:36:03
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-21 11:26:50
 * @FilePath: \common\infra\common\error.go
 * @Description: 公共错误信息定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package common

import "errors"

// "自定义错误"相关定义
var (
	ErrReceiveEventTimeout = errors.New("接收事件出现超时")
	ErrReceiveDataTimeout  = errors.New("接收数据出现超时")
	ErrAdvanceExit         = errors.New("提前退出")
)
