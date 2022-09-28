/*
 * @Author: hongliu
 * @Date: 2022-09-23 17:54:58
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-24 16:21:30
 * @FilePath: \common\infra\common\oss.go
 * @Description:Oss接口定义
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package common

import (
	"context"
	"io"
)

// Oss 对象存储通用接口定义
type Oss interface {
	Publish(ctx context.Context, sourcePath string, destPath string) error                                // 向oss发布文件，sourcePath为本地文件路径，destPath为目标路径
	PublishFromReader(ctx context.Context, filename string, reader io.Reader) (string, error)             // 从缓存向oss发布文件
	Copy(ctx context.Context, sourcePath string, destPath string) error                                   // 在oss中复制文件
	Get(ctx context.Context, filePath string) ([]byte, error)                                             // 从oss下载文件
	ConstructTemporarySignatureUrl(ctx context.Context, originUrl string) (string, error)                 // 构建临时签名Url
	ConstructTemporarySignatureUrls(ctx context.Context, urlMap map[uint]string) (map[uint]string, error) // 批量构建临时签名Url
	GetOriginUrl(ctx context.Context, url string) (string, error)                                         // 获取原始的url
}
