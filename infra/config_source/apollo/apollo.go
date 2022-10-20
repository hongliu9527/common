/*
 * @Author: hongliu
 * @Date: 2022-09-21 10:30:45
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:52:46
 * @FilePath: \common\infra\config_source\apollo\apollo.go
 * @Description: Apollo配置数据源定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package apollo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hongliu9527/common/infra/common"

	"github.com/hongliu9527/go-tools/logger"
	"github.com/shima-park/agollo"
	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
)

// 编译期保证接口实现的一致性
var _ common.ConfigSource = (*ApolloConfigSource)(nil)

// APOLLO 在系统中的标识
const APOLLO = "apollo"

// ApolloConfigSource 数据配置源定义
type ApolloConfigSource struct {
	errChannel <-chan *agollo.LongPollerError // Apollo错误传递通道
	appID      string                         // Apollo的AppId
	endPoint   string                         // Apollo访问地址
	ctx        context.Context                // 上下文对象
	cancel     context.CancelFunc             // 退出回调函数
}

// New 创建Apollo配置数据源
func New(serviceName string, endpoint string) *ApolloConfigSource {
	return &ApolloConfigSource{
		appID:    serviceName,
		endPoint: endpoint,
	}
}

// Init 初始化Apollo数据配置源
func (a *ApolloConfigSource) Init(ctx context.Context) error {
	// 初始化Apollo数据配置源
	apolloCtx, apolloCancel := context.WithCancel(ctx)
	a.ctx = apolloCtx
	a.cancel = apolloCancel

	// 初始化Apollo客户端
	err := agollo.Init(a.endPoint, a.appID,
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))), // 打印Apollo日志信息
		// agollo.PreloadNamespaces("TEST.Namespace"),        // 预先加载的namespace列表，如果是通过配置启动，会在app.properties配置的基础上追加
		agollo.AutoFetchOnCacheMiss(),       // 在配置未找到时，去Apollo的带缓存的获取配置接口，获取配置
		agollo.FailTolerantOnBackupExists(), // 在连接Apollo失败时，如果在配置的目录下存在.Apollo备份配置，会读取备份在服务器无法连接的情况下
	)

	// 启动Apollo客户端
	a.errChannel = agollo.Start()

	return err
}

// Read 读取指定配置文件的配置数据
func (a *ApolloConfigSource) Read(filename string, value interface{}, timeout time.Duration) error {
	// 设置远程AppId
	remote.SetAppID(a.appID)
	viperInstance := viper.New()

	// Apollo默认的配置文件是properties格式，为了获取完整的配置信息，需要设置为yaml
	viperInstance.SetConfigType("yaml")

	// 添加Apollo数据配置
	err := viperInstance.AddRemoteProvider(APOLLO, a.endPoint, filename)
	if err != nil {
		return fmt.Errorf("增加Apollo(%s)命名空间(%s)Provider失败(%s)", a.endPoint, filename, err.Error())
	}

	// 读取Apollo远程配置信息
	err = viperInstance.ReadRemoteConfig()
	if err != nil {
		return fmt.Errorf("读取Apollo(%s)命名空间配置数据(%s)失败(%s)", a.endPoint, filename, err.Error())
	}

	// 反序列化配置信息
	configData := viperInstance.AllSettings()
	err = common.DecodeConfig(configData, value)
	if err != nil {
		return fmt.Errorf("反序列化配置信息失败(%s)", err.Error())
	}

	return nil
}

// Listen 读取指定配置文件的配置数据事件
func (a *ApolloConfigSource) Listen(filename string, value interface{}, timeout time.Duration) error {
	// 获取Apollo配置变更监听通道
	stop := make(chan bool)
	configChannel := agollo.WatchNamespace(filename, stop)

	// 添加Apollo配置信息读取超时控制
	timeoutCtx, timeoutCancel := context.WithTimeout(a.ctx, timeout)
	defer timeoutCancel()

	select {
	case err := <-a.errChannel:
		return fmt.Errorf("Apollo配置中心出现严重错误(%s)，请及时进行处理", err.Err.Error())
	case resp := <-configChannel:
		logger.Debug("监听到(%s)配置数据发生变化(%v)", filename, resp)
		return a.Read(filename, value, 20*time.Second)
	case <-a.ctx.Done():
		return common.ErrAdvanceExit
	case <-timeoutCtx.Done():
		return common.ErrReceiveEventTimeout
	}
}
