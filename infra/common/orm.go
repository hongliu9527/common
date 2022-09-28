/*
 * @Author: hongliu
 * @Date: 2022-09-21 15:58:54
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-21 15:59:05
 * @FilePath: \common\infra\common\orm.go
 * @Description:Orm基础设施接口定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package common

import "strings"

// Orm Orm基础设施接口定义
type Orm interface {
	DB(name string) Orm                                   // 执行设置数据库方法
	Model(value interface{}) Orm                          // 查询模型方法
	Create(value interface{}) Orm                         // 执行创建方法
	CreateInBatches(value interface{}, batchSize int) Orm // 执行按批次创建方法
	Delete(value interface{}, conds ...interface{}) Orm   // 执行删除方法
	Select(query interface{}, args ...interface{}) Orm    // 执行设置字段方法
	Group(name string) Orm                                // 执行设置分组方法
	Update(column string, value interface{}) Orm          // 执行更新单个字段方法
	Updates(value interface{}) Orm                        // 执行更新多个字段方法
	Save(value interface{}) Orm                           // 执行整个字段更新
	Scan(value interface{}) Orm                           // 执行查询方法
	Where(query interface{}, args ...interface{}) Orm     // 执行条件方法
	Find(value interface{}) Orm                           // 执行查询方法
	Count(count *int64) Orm                               // 执行统计个数方法
	Debug() Orm                                           // 执行调试信息方法
	Table(name string, args ...interface{}) Orm           // 执行设置表名方法
	Limit(limit int) Orm                                  // 执行设置查询条数方法
	Offset(offset int) Orm                                // 执行设置从第几条数据开始查的方法
	Order(value interface{}) Orm                          // 执行设置结果排序方法
	Error() error                                         // 查询最新的Orm执行错误信息，执行完数据库操作方法后可以调用该方法检查执行错误
	Distinct(args ...interface{}) Orm                     // 执行去重方法
	Pluck(column string, dest interface{}) Orm            // 执行查询单列方法
	Raw(sql string, values ...interface{}) Orm            // 执行原生查询方法
	Filter(params []FilterParam) Orm                      // 执行多字段筛选
	Sort(params []SortParm) Orm                           // 执行多字段排序
	Joins(query string, args ...interface{}) Orm          // 执行联合查询
	Exec(sql string, values ...interface{}) Orm           // 执行原生语句方法
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
