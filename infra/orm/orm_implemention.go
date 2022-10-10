/*
 * @Author: hongliu
 * @Date: 2022-09-23 10:13:07
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-10 11:53:44
 * @FilePath: \common\infra\orm\orm_implemention.go
 * @Description:orm接口实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package orm

import (
	"fmt"
	"hongliu9527/common/infra/common"
)

// Model 查询模型方法
func (i *ormInfra) Model(value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Select 设置查询字段方法
func (i *ormInfra) Select(query interface{}, args ...interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Update 设置更新单个字段的方法
func (i *ormInfra) Update(column string, value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Updates 设置更新多个字段的方法
func (i *ormInfra) Updates(value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Delete 删除方法
func (i *ormInfra) Delete(value interface{}, conds ...interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Create 创建方法
func (i *ormInfra) Create(value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// CreateInBatches 批量创建方法
func (i *ormInfra) CreateInBatches(value interface{}, batchSize int) common.Orm {
	i.mustStartWithConn()

	return i
}

// Group 分组方法
func (i *ormInfra) Group(name string) common.Orm {
	i.mustStartWithConn()

	return i
}

// Where 条件方法
func (i *ormInfra) Where(query interface{}, args ...interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Find 查询方法
func (i *ormInfra) Find(value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Scan 查询方法(无需绑定表名)
func (i *ormInfra) Scan(value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Count 统计个数方法
func (i *ormInfra) Count(count *int64) common.Orm {
	i.mustStartWithConn()

	return i
}

// Debug 执行调试信息
func (i *ormInfra) Debug() common.Orm {
	i.mustStartWithConn()

	return i
}

// Table 执行设置表名方法
func (i *ormInfra) Table(name string, args ...interface{}) common.Orm {
	// 根据表明查询实例
	db, ok := i.tableNameInstance[name]
	if !ok {
		i.lastError = fmt.Errorf("根据表名(%s)无法找到对应的Gorm实例", name)
		return i
	}

	// 创建新的查询会话
	return &session{
		tx: db.WithContext(i.ctx).Table(name, args...),
	}
}

// DB 按名称指定DB实例
func (i *ormInfra) DB(name string) common.Orm {
	db, ok := i.nameInstance[name]
	if !ok {
		i.lastError = fmt.Errorf("根据数据库名(%s)无法找到对应的Gorm实例", name)
		return i
	}
	return &session{
		tx: db.WithContext(i.ctx),
	}
}

// Limit 执行设置查询记录条数方法
func (i *ormInfra) Limit(limit int) common.Orm {
	i.mustStartWithConn()

	return i
}

// Offset 执行设置跳过多少条记录开始查询的方法
func (i *ormInfra) Offset(offset int) common.Orm {
	i.mustStartWithConn()

	return i
}

// Order 执行设置查询结果排序方法
func (i *ormInfra) Order(value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Error 查询最新的Orm执行错误信息
func (i *ormInfra) Error() error {
	i.mustStartWithConn()
	return i.lastError
}

// Distinct 执行结果集去重
func (i *ormInfra) Distinct(args ...interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Pluck 执行查询单列
func (i *ormInfra) Pluck(column string, dest interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Save 执行更新全部字段
func (i *ormInfra) Save(value interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Raw 原生sql查询方法
func (i *ormInfra) Raw(sql string, values ...interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Filter 执行多字段筛选
func (i *ormInfra) Filter(params []common.FilterParam) common.Orm {
	i.mustStartWithConn()

	return i
}

// Sort 执行多字段排序
func (i *ormInfra) Sort(params []common.SortParm) common.Orm {
	i.mustStartWithConn()

	return i
}

// Joins 执行联合查询
func (i *ormInfra) Joins(query string, args ...interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

// Exec 执行原生语句方法
func (i *ormInfra) Exec(sql string, values ...interface{}) common.Orm {
	i.mustStartWithConn()

	return i
}

func (i *ormInfra) mustStartWithConn() {
	if i.lastError == nil {
		i.lastError = fmt.Errorf("数据库查询会话未创建，请检查Orm是否已经最先调用Conn方法")
	}
}
