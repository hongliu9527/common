/*
 * @Author: hongliu
 * @Date: 2022-09-21 15:58:54
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-08 18:56:46
 * @FilePath: \common\infra\common\orm.go
 * @Description:Orm基础设施接口定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package common

import "strings"

// Orm Orm基础设施接口定义
type Orm interface {
	Select(dest interface{}, query string, args ...interface{}) error // 查询多条数据，dest必须是切片
	Get(dest interface{}, query string, args ...interface{}) error    // 查询单个数据
	Update()
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
