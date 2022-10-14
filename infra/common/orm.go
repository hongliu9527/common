/*
 * @Author: hongliu
 * @Date: 2022-09-21 15:58:54
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-14 17:01:01
 * @FilePath: \common\infra\common\orm.go
 * @Description:Orm基础设施接口定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package common

import "strings"

// Orm Orm基础设施接口定义
type Orm interface {
	Conn(tableName string) (orm Orm, err error)                        // 获取数据库连接
	Get(dest interface{}, query string, args ...interface{}) error     // 查询单个数据
	Select(dest interface{}, query string, args ...interface{}) error  // 查询多个数据
	Insert(value interface{}) (uint64, error)                          // 插入单个数据
	BatchInsert(values interface{}) error                              // 批量插入
	Update(condition string, updateValue map[string]interface{}) error // 更新数据
	Delete(condition string, args ...interface{}) error                // 删除数据
	Exec(query string, args ...interface{}) error                      // 执行原生sql
	Begin() error                                                      // 开启事务
	Rollback() error                                                   // 回滚事务
	Commit() error                                                     // 执行事务
}

// FilterParam 筛选参数
type FilterParam struct {
	Name     string   // 筛选参数名称
	Values   []string // 筛选参数值列表
	Operator string   // 筛选条件操作
}

// SortParm 排序参数
type SortParm struct {
	Name      string // 排序条件名称
	Direction string // 排序方向
}

// ToSql 筛选参数转换为SQL语句
func (p FilterParam) ToSql(tableName ...string) string {
	if len(p.Values) == 0 {
		return " 1 = 1 "
	}

	var tablePrefix string
	if len(tableName) > 0 {
		tablePrefix = tableName[0] + "."
	}

	switch p.Operator {
	case "=", "<>", ">", "<", ">=", "<=":
		return " " + tablePrefix + p.Name + " " + p.Operator + " '" + p.Values[0] + "' "
	case "in":
		return " " + tablePrefix + p.Name + " " + p.Operator + " ('" + strings.Join(p.Values, "','") + "') "
	case "like":
		return " " + tablePrefix + p.Name + " " + p.Operator + " '%" + p.Values[0] + "%' "
	}
	return " 1 = 1 "
}

// ToSql 排序参数转换为SQL语句
func (p SortParm) ToSql(tableName ...string) string {
	var tablePrefix string
	if len(tableName) > 0 {
		tablePrefix = tableName[0] + "."
	}
	return " " + tablePrefix + p.Name + " " + p.Direction + " "
}
