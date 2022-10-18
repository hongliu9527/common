/*
 * @Author: hongliu
 * @Date: 2022-10-17 17:30:01
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-18 15:46:58
 * @FilePath: \common\infra\oss\local\oss_implemention.go
 * @Description:oss接口实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"hongliu9527/common/utils"

	"github.com/hongliu9527/go-tools/logger"
)

// Publish 发布文件,在local存储为
func (i *localOssInfra) Publish(ctx context.Context, sourcePath string, targetPath string) error {
	if !utils.IsFilePathExist(sourcePath) {
		logger.Error("文件路径(%s)不存在", sourcePath)
		return fmt.Errorf("文件路径(%s)不存在", sourcePath)
	}

	// 如果目标文件存在，则重命名目标文件
	targetFilePath := i.config.Endpoint + "/" + targetPath
	isBackup := false
	if utils.IsFilePathExist(targetFilePath) {
		err := os.Rename(targetFilePath, targetFilePath+".bak")
		if err != nil {
			logger.Error("重命名(%s.bak)文件失败(%s)", targetFilePath, err.Error())
			return err
		}
		isBackup = true
	}

	// 如果文件夹不存在则创建文件夹
	targetDir := filepath.Dir(targetFilePath)
	if !utils.IsFilePathExist(targetDir) {
		err := os.MkdirAll(targetDir, 0666)
		if err != nil {
			if isBackup {
				os.Rename(targetFilePath+".bak", targetFilePath)
			}
			logger.Error("创建目录(%s)失败(%s)", targetDir, err.Error())
			return err
		}
	}

	// 拷贝文件内容
	content, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		if isBackup {
			os.Rename(targetFilePath+".bak", targetFilePath)
		}
		logger.Error("读取文件(%s)内容失败(%s)", sourcePath, err.Error())
		return err
	}

	// 写入到目标文件
	err = ioutil.WriteFile(targetFilePath, content, 0666)
	if err != nil {
		if isBackup {
			os.Rename(targetFilePath+".bak", targetFilePath)
		}
		logger.Error("写入数据到文件(%s)失败(%s)", targetFilePath, err.Error())
		return err
	}

	if isBackup {
		os.Remove(targetFilePath + ".bak")
	}

	return nil
}

// PublishFromReader 从缓存向oss发布数据,这里的
func (i *localOssInfra) PublishFromReader(ctx context.Context, filename string, reader io.Reader) (string, error) {
	// 如果目标文件存在，则重命名目标文件
	isBackup := false
	if utils.IsFilePathExist(filename) {
		err := os.Rename(filename, filename+".bak")
		if err != nil {
			logger.Error("重命名(%s.bak)文件失败(%s)", filename, err.Error())
			return "", err
		}
		isBackup = true
	}

	targetFilePath := i.config.Endpoint + "/" + filename
	// 如果文件夹不存在则创建文件夹
	targetDir := filepath.Dir(targetFilePath)
	if !utils.IsFilePathExist(targetDir) {
		err := os.MkdirAll(targetDir, 0666)
		if err != nil {
			if isBackup {
				os.Rename(targetFilePath+".bak", targetFilePath)
			}
			logger.Error("创建目录(%s)失败(%s)", targetDir, err.Error())
			return "", err
		}
	}

	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		if isBackup {
			os.Rename(targetFilePath+".bak", targetFilePath)
		}
		logger.Error("创建文件(%s),失败(%s)", targetFilePath, err.Error())
		targetFile.Close()
		return "", err
	}

	_, err = io.Copy(targetFile, reader)
	if err != nil {
		if isBackup {
			os.Rename(targetFilePath+".bak", targetFilePath)
		}
		logger.Error("写入文件(%s),失败(%s)", targetFilePath, err.Error())
		targetFile.Close()
		return "", err
	}

	if isBackup {
		os.Remove(targetFilePath + ".bak")
	}
	targetFile.Close()

	return filename, nil
}

// Copy 在oss中复制文件
func (i *localOssInfra) Copy(ctx context.Context, sourcePath string, destPath string) error {
	sourceFilePath := i.config.Endpoint + "/" + sourcePath

	if !utils.IsFilePathExist(sourceFilePath) {
		logger.Error("源文件(%s)不存在", sourceFilePath)
		return fmt.Errorf("源文件(%s)不存在", sourceFilePath)
	}

	targetFilePath := i.config.Endpoint + "/" + destPath

	// 如果目标文件夹不存在，则创建文件夹
	targetDir := filepath.Dir(targetFilePath)
	if !utils.IsFilePathExist(targetDir) {
		err := os.MkdirAll(targetDir, 0666)
		if err != nil {
			logger.Error("创建目录(%s)失败(%s)", targetDir, err.Error())
			return err
		}
	}

	// 如果目标文件夹存在，则先修改文件名备份
	isBackup := false
	if utils.IsFilePathExist(targetFilePath) {
		err := os.Rename(targetFilePath, targetFilePath+".bak")
		if err != nil {
			logger.Error("重命名(%s.bak)文件失败(%s)", targetFilePath, err.Error())
			return err
		}
		isBackup = true
	}

	// 创建目标文件
	targetFile, err := os.Create(targetFilePath)
	defer targetFile.Close()
	if err != nil {
		if isBackup {
			os.Rename(targetFilePath+".bak", targetFilePath)
		}
		logger.Error("创建文件(%s),失败(%s)", targetFilePath, err.Error())
		return err
	}

	// 打开源文件
	sourceFile, err := os.Open(sourceFilePath)
	defer sourceFile.Close()
	if err != nil {
		if isBackup {
			os.Rename(targetFilePath+".bak", targetFilePath)
		}
		logger.Error("打开文件(%s),失败(%s)", sourceFile, err.Error())
		return err
	}

	// 拷贝数据
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		if isBackup {
			os.Rename(targetFilePath+".bak", targetFilePath)
		}

		logger.Error("拷贝文件(%s)到目标文件(%s),失败(%s)", sourceFilePath, targetFilePath, err.Error())
		return err
	}

	if isBackup {
		os.Remove(targetFilePath + "dir")
	}
	return nil
}

// Get 从oss下载文件
func (i *localOssInfra) Get(ctx context.Context, filePath string) ([]byte, error) {
	sourceFilePath := i.config.Endpoint + "/" + filePath
	content, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		logger.Error("读取文件(%s)内容失败(%s)", sourceFilePath, err.Error())
		return nil, err
	}

	return content, nil
}

// ConstructTemporarySignatureUrl 构建临时签名url
func (i *localOssInfra) ConstructTemporarySignatureUrl(ctx context.Context, originUrl string) (string, error) {
	return "", errors.New("本地文件暂不支持构建临时签名url接口")
}

// ConstructTemporarySignatureUrls 批量构建临时签名url接口
func (i *localOssInfra) ConstructTemporarySignatureUrls(ctx context.Context, urlMap map[uint]string) (map[uint]string, error) {
	return nil, errors.New("本地文件暂不支持批量构建临时签名url接口")
}

// GetOriginUrl 获取原始的url
func (i *localOssInfra) GetOriginUrl(ctx context.Context, url string) (string, error) {
	return "", errors.New("本地文件暂不支持获取原始的url接口")
}
