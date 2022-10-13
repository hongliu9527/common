/*
 * @Author: hongliu
 * @Date: 2022-09-23 10:13:07
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-13 15:11:31
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

// Conn 获取数据库查询句柄
func (i *ormInfra) Conn(name string) common.Orm {
	// 根据表明查询实例
	db, ok := i.tableNameInstance[name]
	if !ok {
		i.lastError = fmt.Errorf("根据表名(%s)无法找到对应的Gorm实例", name)
		return i
	}

	// 创建新的查询会话
	return &dbConnection{
		db:          db,
		tx:          nil,
		serviceName: i.InfraName,
		tableName:   name,
		lastError:   nil,
	}
}

// Get 查询单个数据
func (i *ormInfra) Get(dest interface{}, query string, args ...interface{}) error {
	i.mustStartWithConn()
	return i.lastError
}

// Select 查询多个数据
func (i *ormInfra) Select(dest interface{}, query string, args ...interface{}) error {
	i.mustStartWithConn()
	return i.lastError
}

// Insert 创建单个数据
func (i *ormInfra) Insert(data interface{}) (uint64, error) {
	i.mustStartWithConn()
	return 0, i.lastError
}

// BatchInsert 批量插入
func (i *ormInfra) BatchInsert(datas []interface{}) error {
	i.mustStartWithConn()
	return i.lastError
}

// Update 更新数据
func (i *ormInfra) Update(query string, arg map[string]interface{}) error {
	i.mustStartWithConn()
	return i.lastError
}

// Exec 执行原生sql
func (i *ormInfra) Exec(qeury string, args ...interface{}) error {
	i.mustStartWithConn()
	return i.lastError
}

// Delete 删除数据
func (i *ormInfra) Delete(qeury string, args ...interface{}) error {
	i.mustStartWithConn()
	return i.lastError
}

// Begin 开启事务
func (i *ormInfra) Begin() (common.Orm, error) {
	i.mustStartWithConn()
	return nil, i.lastError
}

// RollBack 事务回滚
func (i *ormInfra) RollBack() {
	i.mustStartWithConn()
}

// Commit 执行事务
func (i *ormInfra) Commit() error {
	i.mustStartWithConn()
	return i.lastError
}

func (i *ormInfra) mustStartWithConn() {
	if i.lastError == nil {
		i.lastError = fmt.Errorf("数据库查询连接未创建，请检查Orm是否已经最先调用Conn方法")
	}
}
