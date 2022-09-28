/*
 * @Author: hongliu
 * @Date: 2022-09-23 10:15:54
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-23 10:22:02
 * @FilePath: \common\infra\orm\session.go
 * @Description: Orm会话封装
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package orm

import (
	"fmt"
	"strings"

	"hongliu9527/common/infra/common"
	"hongliu9527/common/utils"

	"gorm.io/gorm"
)

// session 查询会话
type session struct {
	tx          *gorm.DB // 查询过程实例
	serviceName string   // 服务名称，用于日志记录
	lastError   error    // 实例的最新错误信息
}

// Model 查询模型方法
func (s *session) Model(value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Model(value)
	s.lastError = s.tx.Error

	return s
}

// Select 设置查询字段方法
func (s *session) Select(query interface{}, args ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Select(query, args...)
	s.lastError = s.tx.Error

	return s
}

// Update 设置更新单个字段的方法
func (s *session) Update(column string, value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx.Update(column, value)
	s.lastError = s.tx.Error

	return s
}

// Updates 设置更新多个字段的方法
func (s *session) Updates(value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	// 校验参数是否是map[string]interface{}
	updateData, ok := value.(map[string]interface{})
	if !ok {
		s.lastError = fmt.Errorf("更新多个字段参数必须是map[string]interface{}")
		return s
	}

	s.tx.Updates(updateData)
	s.lastError = s.tx.Error

	return s
}

// Delete 删除方法
func (s *session) Delete(value interface{}, conds ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx.Delete(value, conds...)
	s.lastError = s.tx.Error

	return s
}

// Create 创建方法
func (s *session) Create(value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx.Create(value)
	s.lastError = s.tx.Error

	return s
}

// CreateInBatches 批量创建方法
func (s *session) CreateInBatches(value interface{}, batchSize int) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx.CreateInBatches(value, batchSize)
	s.lastError = s.tx.Error

	return s
}

// Where 条件方法
func (s *session) Where(query interface{}, args ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Where(query, args...)
	s.lastError = s.tx.Error

	return s
}

// Find 查询方法
func (s *session) Find(value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx.Find(value)
	s.lastError = s.tx.Error

	return s
}

// Scan 查询方法 无需绑定表明
func (s *session) Scan(value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx.Scan(value)
	s.lastError = s.tx.Error

	return s
}

// Count 统计个数方法
func (s *session) Count(count *int64) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx.Count(count)
	s.lastError = s.tx.Error

	return s
}

// Debug 执行调试信息
func (s *session) Debug() common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Debug()
	s.lastError = s.tx.Error

	return s
}

// Table 执行设置表名方法
func (s *session) Table(name string, args ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Table(name, args...)
	s.lastError = s.tx.Error

	return s
}

// Group 执行分组方法
func (s *session) Group(name string) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Group(name)
	s.lastError = s.tx.Error

	return s
}

// Limit 执行设置查询记录条数方法
func (s *session) Limit(limit int) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Limit(limit)
	s.lastError = s.tx.Error

	return s
}

// Offset 执行设置跳过多少条记录开始查询的方法
func (s *session) Offset(offset int) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Offset(offset)
	s.lastError = s.tx.Error

	return s
}

// Order 执行设置查询结果排序的方法
func (s *session) Order(value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Order(value)
	s.lastError = s.tx.Error

	return s
}

// Error 查询最新的Orm执行错误信息
func (s *session) Error() error {
	return s.lastError
}

// Distinct 执行结果集去重
func (s *session) Distinct(args ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Distinct(args)
	s.lastError = s.tx.Error

	return s
}

// Pluck 执行查询单列
func (s *session) Pluck(column string, dest interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Pluck(column, dest)
	s.lastError = s.tx.Error

	return s
}

// Save 执行更新全部字段
func (s *session) Save(value interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Save(value)
	s.lastError = s.tx.Error

	return s
}

// Raw 原生查询方法
func (s *session) Raw(sql string, values ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Raw(sql, values...)
	s.lastError = s.tx.Error
	return s
}

// Filter 执行多字段筛选
func (s *session) Filter(params []common.FilterParam) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	for _, param := range params {
		if strings.ToLower(param.Operator) == "in" {
			s.tx.Where(param.Name+" "+param.Operator+" ?", param.Values)
		} else if strings.ToLower(param.Operator) == "like" {
			s.tx.Where(param.Name+" "+param.Operator+" ?", "%"+param.Values[0]+"%")
		} else if strings.ToLower(param.Operator) == "between" {
			// 当为between时，小的必须在前
			isGreater, err := utils.Compare(param.Values[0], param.Values[1])
			if err != nil {
				s.lastError = fmt.Errorf("传入的值不支持比较(%s)", err.Error())
				return s
			}
			if isGreater { // 当Values[0] 大于Values[1]时，交换两值 此处不能用其他方法交换，防止string或者浮点类型值不支持
				temp := param.Values[0]
				param.Values[0] = param.Values[1]
				param.Values[1] = temp
			}
			s.tx.Where(param.Name+" "+param.Operator+" ? and ?", param.Values[0], param.Values[1])
		} else {
			s.tx.Where(param.Name+" "+param.Operator+" ?", param.Values[0])
		}
		if s.tx.Error != nil {
			s.lastError = s.tx.Error
		}
	}
	return s
}

// Sort 执行多字段排序
func (s *session) Sort(params []common.SortParm) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	for _, param := range params {
		s.tx.Order(param.Name + " " + param.Direction)
		if s.tx.Error != nil {
			s.lastError = s.tx.Error
		}
	}
	return s
}

// Joins 执行联合查询
func (s *session) Joins(query string, args ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}

	s.tx = s.tx.Joins(query, args...)
	s.lastError = s.tx.Error
	return s
}

// DB 按名称指定DB实例
func (s *session) DB(name string) common.Orm {
	return s
}

// Exec 执行原生sql
func (s *session) Exec(sql string, values ...interface{}) common.Orm {
	if s.tx == nil {
		s.lastError = fmt.Errorf("数据库查询会话为空，请调用GormInfra.Table生成新的查询会话")
		return s
	}
	s.tx = s.tx.Exec(sql, values...)
	s.lastError = s.tx.Error
	return s
}
