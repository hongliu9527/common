/*
 * @Author: hongliu
 * @Date: 2022-09-21 16:59:20
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-20 14:54:17
 * @FilePath: \common\infra\orm\orm.go
 * @Description:Orm基础设施实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package orm

import (
	"context"

	"github.com/hongliu9527/common/infra/base"
	"github.com/hongliu9527/common/infra/common"
	"github.com/hongliu9527/common/infra/orm/config"

	"github.com/jmoiron/sqlx"
)

// 编译期保证接口实现的一致性
var _ common.OrmInfra = (*ormInfra)(nil)

// ormInfra orm基础设施定义类型定义
type ormInfra struct {
	base.BaseInfra                                     // 基础设施基类
	config            *config.OrmInfraConfig           // 数据库配置信息
	nameConfig        map[string]config.DataBaseConfig // 数据库实例名-配置信息哈希表
	tableNameInstance map[string]*sqlx.DB              // 数据库表名-数据库实例哈希表
	nameInstance      map[string]*sqlx.DB              // 数据库实例名-数据库实例哈希表
	lastError         error                            // 实例的最新错误信息

	ctx    context.Context    // 上下文对象
	cancel context.CancelFunc // 取消回调函数
}

// orm基础设施单例
var singleton ormInfra

// New 创建orm基础设施
func New(ormConfig *config.OrmInfraConfig) common.OrmInfra {

	singleton.config = ormConfig
	singleton.nameConfig = make(map[string]config.DataBaseConfig)
	singleton.tableNameInstance = make(map[string]*sqlx.DB)
	singleton.nameInstance = make(map[string]*sqlx.DB)

	// 构建基础设施基类
	singleton.BaseInfra = base.NewBaseInfra(singleton.Name(), ormConfig, singleton.start, singleton.stop)

	return &singleton
}
