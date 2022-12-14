/*
 * @Author: hongliu
 * @Date: 2022-09-21 10:37:45
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-21 16:55:12
 * @FilePath: \common\infra\config_source\local\local.go
 * @Description:Local配置数据源定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package local

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hongliu9527/common/infra/common"

	"github.com/hongliu9527/go-tools/logger"
	"github.com/spf13/viper"
)

// LOCAL 在系统中的标识
const LOCAL = "local"

// LocalConfigSource 数据源配置定义
type LocalConfigSource struct {
	moduleName              string             // 模块名称
	basePath                string             // 文件路径
	lastFileChangeTimestamp uint32             // 文件上次修改事件
	ctx                     context.Context    // 上下文对象
	cancel                  context.CancelFunc // 退出回调函数
}

// New 新建本地数据源
func New(moduleName, basePath string) *LocalConfigSource {
	return &LocalConfigSource{
		moduleName: moduleName,
		basePath:   basePath,
	}
}

// Init 初始化配置文件数据配置源
func (l *LocalConfigSource) Init(ctx context.Context) error {
	l.ctx, l.cancel = context.WithCancel(ctx)

	fileInfo, err := os.Stat(l.basePath)
	if err != nil {
		return fmt.Errorf("获取模块(%s)的配置文件路径(%s)的信息失败(%s)", l.moduleName, l.basePath, err.Error())
	}

	l.lastFileChangeTimestamp = uint32(fileInfo.ModTime().Unix())
	return nil
}

// Read 读取指定配置文件的配置数据
func (l *LocalConfigSource) Read(filename string, value interface{}, timeout time.Duration) error {
	viperInstance := viper.New()
	configFileName := l.basePath + filename
	viper.SetConfigFile(configFileName)

	if err := viperInstance.ReadInConfig(); err != nil {
		return fmt.Errorf("读取模块(%s)的配置文件(%s)失败(%s)", l.moduleName, configFileName, err.Error())
	}

	// 反序列化配置信息
	configData := viperInstance.AllSettings()
	err := common.DecodeConfig(configData, value)
	if err != nil {
		return fmt.Errorf("反序列化配置信息失败(%s)", err.Error())
	}

	return nil
}

// Listen 读取指定配置文件的配置数据事件，通过对比文件文件修改时间戳，判断文件是否被改变
func (l *LocalConfigSource) Listen(filename string, value interface{}, timeout time.Duration) error {
	timeoutCtx, timeoutCancel := context.WithTimeout(l.ctx, timeout)
	defer timeoutCancel()

	configFileName := l.basePath + filename
	fileInfo, err := os.Stat(configFileName)
	if err != nil {
		return fmt.Errorf("获取模块(%s)的配置文件(%s)的信息失败(%s),请及时处理", l.moduleName, configFileName, err.Error())
	}
	currentModTimestamp := uint32(fileInfo.ModTime().Unix())

	if currentModTimestamp > l.lastFileChangeTimestamp {
		logger.Debug("监听到(%s)配置数据发生变化", filename)
		return l.Read("", value, 20*time.Second)
	}

	select {
	case <-l.ctx.Done():
		return common.ErrAdvanceExit
	case <-timeoutCtx.Done():
		return common.ErrReceiveEventTimeout
	}
}
