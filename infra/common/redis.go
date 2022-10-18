/*
 * @Author: hongliu
 * @Date: 2022-10-18 15:49:08
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-18 15:50:45
 * @FilePath: \common\infra\common\redis.go
 * @Description:Redis接口定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package common

import "time"

// Redis Redis接口定义
type Redis interface {
	Set(key string, value interface{}, expiration time.Duration) error           // 设置单个字符串数据
	Get(key string) (string, error)                                              // 读取单个字符串数据
	ReadAllKeys(match string) ([]string, error)                                  // 读取所有的key
	Delete(key string) error                                                     // 删除
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error) // 单个键值数据不存在时设置
}
