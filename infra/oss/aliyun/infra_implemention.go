/*
 * @Author: hongliu
 * @Date: 2022-10-17 15:04:36
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-17 15:07:37
 * @FilePath: \common\infra\oss\aliyun\infra_implemention.go
 * @Description: 基础设施接口的实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package aliyun

import (
	"context"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

// 常量相关定义
const (
	AliyunOssInfraName string = "AliyunOss" // 阿里云Oss基础设施名称
)

// Name 获取基础设施名
func (i *aliyunOssInfra) Name() string {
	return AliyunOssInfraName
}

// Start 启动基础设施
func (i *aliyunOssInfra) start(ctx context.Context) error {
	client, err := oss.New(i.config.Endpoint, i.config.AccessKeyID, i.config.AccessKeySecret)
	if err != nil {
		return errors.Wrap(err, "创建阿里云Oss客户端失败")
	}

	// 获取存储空间
	bucket, err := client.Bucket(i.config.Bucket)
	if err != nil {
		return errors.Wrap(err, "获取阿里云Oss存储空间失败")
	}
	i.bucket = bucket

	return nil
}

// Stop 停止基础设施
func (i *aliyunOssInfra) stop() error {
	return nil
}
