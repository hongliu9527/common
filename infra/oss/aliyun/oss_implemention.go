/*
 * @Author: hongliu
 * @Date: 2022-10-17 15:06:14
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-17 15:13:09
 * @FilePath: \common\infra\oss\aliyun\oss_implemention.go
 * @Description: oss接口的实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package aliyun

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hongliu9527/go-tools/logger"
	"github.com/pkg/errors"
)

// Publish 发布文件
func (i *aliyunOssInfra) Publish(ctx context.Context, sourcePath string, targetPath string) error {

	retryCount := i.config.RetryCount
	// Oss分片上传文件方式(该方式对网络的要求相对更低，避免了大文件上传超时导致的失败)
	err := i.bucket.UploadFile(targetPath, sourcePath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	if err != nil {
		if retryCount == 0 {
			return errors.Wrap(err, "首次上传文件到阿里云Oss失败，并且不会进行重试上传操作")
		}

		logger.Error("首次上传文件到阿里云Oss失败(%s),接下来将进行%d次重试上传", err.Error(), retryCount)

		if err := i.retransmission(sourcePath, targetPath, retryCount); err != nil {
			return errors.Wrap(err, "重试上传文件到阿里云Oss失败")
		}
	}

	logger.Info("上传文件到阿里云Oss成功(%s)", targetPath)

	return nil
}

// PublishFromReader 从io.Reader上传文件到oss
func (i *aliyunOssInfra) PublishFromReader(ctx context.Context, filename string, reader io.Reader) (string, error) {

	err := i.bucket.PutObject(filename, reader)
	if err != nil {
		logger.Error("上传文件(%s)到桶(%s)失败(%s)", filename, i.config.Bucket, err.Error())

		return "", errors.Wrap(err, "上传文件到阿里云Oss失败")
	}
	logger.Info("文件(%s)上传成功", filename)

	host := strings.Trim(i.config.Endpoint, "http://")
	return fmt.Sprintf("http://%s.%s/%s", i.config.Bucket, host, filename), nil
}

// Copy 在oss中复制文件
func (i *aliyunOssInfra) Copy(ctx context.Context, sourcePath string, destPath string) error {
	return errors.Wrapf(
		i.bucket.PutSymlink(destPath, sourcePath),
		"对文件(%s)执行软链接操作失败", sourcePath)
}

// Get 从oss下载文件
func (i *aliyunOssInfra) Get(ctx context.Context, filePath string) ([]byte, error) {
	originUrl, err := i.GetOriginUrl(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("获取原始url失败(%s)", err.Error())
	}

	prefix := "http://" + i.config.Bucket + "." + strings.Trim(i.config.Endpoint, "http://") + "/"
	path := strings.TrimPrefix(originUrl, prefix)

	reader, err := i.bucket.GetObject(path)
	if err != nil {
		return nil, errors.Wrapf(err, "从阿里云oss下载文件失败(%s)", originUrl)
	}
	defer reader.Close()

	content, err := ioutil.ReadAll(reader)
	return content, errors.Wrapf(err, "阿里云oss下载文件后,读取失败(%s)", originUrl)
}

// ConstructTemporarySignatureUrl 构建临时签名url
// 阿里oss参考文档地址：https://help.aliyun.com/document_detail/100624.htm?spm=a2c4g.11186623.0.0.2deb4f77JjIsXR#concept-xzh-nzk-2gb
// 另外需要配置AliyunOSSFullAccess权限
func (i *aliyunOssInfra) ConstructTemporarySignatureUrl(ctx context.Context, originUrl string) (string, error) {
	// 当传入的原始url中的bucket与当前bucket不同时，不进行处理
	infos := strings.Split(originUrl, ".")
	if i.config.Bucket != strings.TrimPrefix(infos[0], "http://") {
		return originUrl, nil
	}

	// 构建生成sts客户端
	stsClient, err := sts.NewClientWithAccessKey(strings.TrimPrefix(strings.Split(strings.TrimPrefix(i.config.Endpoint, "http://"), ".")[0], "oss-"), i.config.AccessKeyID, i.config.AccessKeySecret)
	if err != nil {
		return "", fmt.Errorf("生成阿里oss临时访问客户端失败(%s)", err.Error())
	}

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = i.config.RoleARN
	request.RoleSessionName = "SessionTest"
	response, err := stsClient.AssumeRole(request)
	if err != nil {
		return "", fmt.Errorf("获取临时访问阿里oss Token令牌失败(%s)", err.Error())
	}

	// 使用sts构建oss client
	client, err := oss.New(i.config.Endpoint, response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, oss.SecurityToken(response.Credentials.SecurityToken))
	if err != nil {
		return "", fmt.Errorf("创建阿里oss客户端失败(%s)", err.Error())
	}

	// 构建bucket
	bucketName := i.config.Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("获取阿里oss bucket失败(%s)", err.Error())
	}
	expireTime, err := strconv.Atoi(i.config.SignatureExpiresTime)
	if err != nil {
		return "", fmt.Errorf("临时访问阿里oss过期时间转换失败(%s)", err.Error())
	}

	prefix := "http://" + i.config.Bucket + "." + strings.TrimPrefix(i.config.Endpoint, "http://") + "/"
	url, err := bucket.SignURL(strings.TrimPrefix(originUrl, prefix), oss.HTTPGet, int64(expireTime))
	if err != nil {
		return "", fmt.Errorf("获取访问阿里oss临时url失败(%s)", err.Error())
	}

	return url, nil
}

// ConstructTemporarySignatureUrls 批量构建临时签名
func (i *aliyunOssInfra) ConstructTemporarySignatureUrls(ctx context.Context, urlMap map[uint]string) (map[uint]string, error) {
	// 构建生成sts客户端
	stsClient, err := sts.NewClientWithAccessKey(strings.TrimPrefix(strings.Split(strings.TrimPrefix(i.config.Endpoint, "http://"), ".")[0], "oss-"), i.config.AccessKeyID, i.config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("生成阿里oss临时访问客户端失败(%s)", err.Error())
	}

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = i.config.RoleARN
	request.RoleSessionName = "SessionTest"
	response, err := stsClient.AssumeRole(request)
	if err != nil {
		return nil, fmt.Errorf("获取临时访问阿里oss Token令牌失败(%s)", err.Error())
	}

	// 使用sts构建oss client
	client, err := oss.New(i.config.Endpoint, response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, oss.SecurityToken(response.Credentials.SecurityToken))
	if err != nil {
		return nil, fmt.Errorf("创建阿里oss客户端失败(%s)", err.Error())
	}

	// 构建bucket
	bucketName := i.config.Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("获取阿里oss bucket失败(%s)", err.Error())
	}
	expireTime, err := strconv.Atoi(i.config.SignatureExpiresTime)
	if err != nil {
		return nil, fmt.Errorf("临时访问阿里oss过期时间转换失败(%s)", err.Error())
	}

	signatureMap := make(map[uint]string)

	for id, url := range urlMap {
		infos := strings.Split(url, ".")
		if i.config.Bucket != strings.TrimPrefix(infos[0], "http://") {
			signatureMap[id] = url
			continue
		}

		prefix := "http://" + i.config.Bucket + "." + strings.TrimPrefix(i.config.Endpoint, "http://") + "/"
		signatureUrl, err := bucket.SignURL(strings.TrimPrefix(url, prefix), oss.HTTPGet, int64(expireTime))
		if err != nil {
			logger.Error("构建临时签名URL(原始url为:%s)失败(%s)", url, err.Error())
			continue
		}
		signatureMap[id] = signatureUrl
	}

	return signatureMap, nil
}

// GetOriginUrl 获取原始的url
func (i *aliyunOssInfra) GetOriginUrl(ctx context.Context, signatureURL string) (string, error) {
	signatureURL, err := url.QueryUnescape(signatureURL)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(signatureURL)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(signatureURL, "?"+u.RawQuery), nil
}

// retransmission 实现文件重发
// 初次发布失败后进行指定次数的重发操作
func (i *aliyunOssInfra) retransmission(sourcePath string, targetPath string, maxRetryCount int) error {
	for retryCount := 0; retryCount < maxRetryCount; retryCount++ {
		time.Sleep(time.Duration(500*(retryCount+1)) * time.Millisecond)

		err := i.bucket.UploadFile(targetPath, sourcePath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
		if err != nil {
			logger.Error("第%d次上传文件到阿里云Oss失败(%s)", retryCount+1, err.Error())
			continue
		}

		logger.Info("第%d次上传文件到阿里云Oss成功(%s)", retryCount+1, targetPath)

		// 当重新上传后以正常返回值提前退出并记录相关成功上传日志
		return nil
	}

	return errors.Errorf("上传文件到阿里云Oss重试%d次后仍然失败", maxRetryCount)
}
