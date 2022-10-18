/*
 * @Author: hongliu
 * @Date: 2022-10-18 16:48:10
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-18 16:49:33
 * @FilePath: \common\infra\redis\redis_implemention.go
 * @Description:redis 接口实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package redis_infra

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Set 设置单个字符串数据
func (i *RedisInfra) Set(key string, value interface{}, expiration time.Duration) error {
	return errors.Wrapf(i.client.Set(i.ctx, key, value, expiration).Err(), "redis写失败(%s)", key)
}

// Get 查询单个字符串数据
func (i *RedisInfra) Get(key string) (string, error) {
	value, err := i.client.Get(i.ctx, key).Result()

	if err == redis.Nil {
		return "", errors.Errorf("(%s)键不存在", key)
	}

	return value, errors.Wrapf(err, "redis读失败(%s)", key)
}

// ReadAllKeys 读取所有的key
func (i *RedisInfra) ReadAllKeys(match string) ([]string, error) {
	keys, _, err := i.client.Scan(i.ctx, 0, match, 1000000).Result()
	return keys, err
}

// Delete 删除
func (i *RedisInfra) Delete(key string) error {
	return i.client.Del(i.ctx, key).Err()
}

// SetNX 单个键值数据不存在时设置
func (i *RedisInfra) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return i.client.SetNX(i.ctx, key, value, expiration).Result()
}
