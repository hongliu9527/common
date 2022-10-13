/*
 * @Author: hongliu
 * @Date: 2022-09-23 10:15:54
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-13 18:02:52
 * @FilePath: \common\infra\orm\db_connection.go
 * @Description: Orm实例封装
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"hongliu9527/common/infra/common"
	"hongliu9527/common/utils"

	"github.com/jmoiron/sqlx"
)

// dbConnection 查询连接
type dbConnection struct {
	db          *sqlx.DB // 查询实例
	tx          *sqlx.Tx // 事务实例
	serviceName string   // 服务名称，用于日志记录
	tableName   string   // 当前需要操作的表名称
	lastError   error    // 实例的最新错误信息
}

// Conn 获取数据库连接句柄
func (c *dbConnection) Conn(name string) common.Orm {
	return c
}

// Get 查询单个数据
func (c *dbConnection) Get(dest interface{}, query string, args ...interface{}) error {
	// 如果事务实例存在，则使用事务实例执行查询
	if c.tx != nil {
		return c.tx.Get(dest, query, args...)
	}

	return c.db.Get(dest, query, args...)
}

// Select 查询多个数据
func (c *dbConnection) Select(dest interface{}, query string, args ...interface{}) error {
	// 传入的参数必须是切片的指针类型
	if utils.IsSlicePointer(dest) == false {
		return errors.New("传入的参数必须是切片指针")
	}

	// 如果事务实例存在，则使用事务实例进行查询
	if c.tx != nil {
		return c.tx.Select(dest, query, args...)
	}

	return c.db.Select(dest, query, args...)
}

// Insert 插入单个数据
func (c *dbConnection) Insert(value interface{}) (uint64, error) {
	query, err := c.createInsertSql(value)
	if err != nil {
		return 0, err
	}

	// 如果事务实例存在，则使用事务实例进行插入
	if c.tx != nil {
		txResult, err := c.tx.NamedExec(query, value)
		if err != nil {
			return 0, nil
		}
		id, _ := txResult.LastInsertId()
		return id, nil
	}

	dbResult, err := c.db.NamedExec(query, value)
	if err != nil {
		return 0, nil
	}
	id, _ := dbResult.LastInsertId()
	return id, nil
}

// createInsertSql 根据传入的带db标记的结构体指针新增数据
func (c *dbConnection) createInsertSql(value interface{}) (string, error) {
	tags, fields, err := dbFields(value)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", c.tableName, strings.Join(tags, ","), strings.Join(fields, ","))
	return query, nil
}

// dbFields 根据结构体获取db标签名数组
func dbFields(value interface{}) ([]string, []string, error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	tagNames := []string{}
	fields := []string{}
	if v.Kind() != reflect.Struct {
		return nil, nil, errors.New("参数必须是结构体或者结构体指针")
	}
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Type().Field(i)
		tagName := fieldValue.Tag.Get("db")
		field := fieldValue.Name
		if tagName != "" {
			tagNames = append(tagNames, tagName)
			fields = append(fields, ":"+field)
		}
	}
	return tagNames, fields, nil
}

// BatchInsert 批量插入
func (c *dbConnection) BatchInsert(values interface{}) error {
	sliceValues := reflect.ValueOf(values)
	if sliceValues.Kind() != reflect.Slice {
		return errors.New("批量插入的参数必须为结构体切片")
	}

	if sliceValues.Len() == 0 {
		return errors.New("批量插入的参数切片不能为空")
	}

	query, err := c.createInsertSql(sliceValues.Index(0).Interface())
	if err != nil {
		return err
	}

	// 如果事务实例存在，则使用事务批量插入
	if c.tx != nil {
		_, err = c.tx.NamedExec(query, values)
		return err
	}

	_, err = c.db.NamedExec(query, values)
	return err
}

// Update 更新数据
func (c *dbConnection) Update(query string, arg map[string]interface{}) error {
	// 如果事务实例存在，则使用事务更新
	if c.tx != nil {
		_, err := c.tx.NamedExec()
	}

}
