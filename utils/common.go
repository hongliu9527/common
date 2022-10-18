/*
 * @Author: hongliu
 * @Date: 2022-09-16 14:57:22
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-18 15:47:15
 * @FilePath: \common\utils\common.go
 * @Description: 定义一些全局公用的方法和类型
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// 全局常量定义
const (
	// 日期格式字符串
	DateFormat = "2006-01-02"

	// 时间格式字符串
	TimeFormat = "2006-01-02 15:04:05"
)

// IsFilePathExist 判断某个文件路径是否存在
func IsFilePathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// CreatFilePath 创建文件路径
func CreatFilePath(path string) error {
	if !IsFilePathExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// IsStructPointer 判断是否为结构体指针
func IsStructPointer(value interface{}) bool {
	reflectType := reflect.TypeOf(value)
	return reflectType.Kind() == reflect.Ptr && reflectType.Elem().Kind() == reflect.Struct
}

// IsSlicePointer 判断是否为切片指针
func IsSlicePointer(value interface{}) bool {
	reflectType := reflect.TypeOf(value)
	return reflectType.Kind() == reflect.Ptr && reflectType.Elem().Kind() == reflect.Slice
}

// IsMap 判断是否为map
func IsMap(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Map
}

// IsSlice 判断参数是否为slice
func IsSlice(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

// IsBasicType 判断参数类型是否为基础类型
func IsBasicType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

// MergeErrors 合并错误列表为一个错误，使用换行符连接报错信息
func MergeErrors(errorLisst []error) error {
	var errMessage string
	if len(errorLisst) == 0 {
		return nil
	}

	for index, err := range errorLisst {
		if index == 0 {
			errMessage = err.Error()
		} else {
			errMessage += "\n" + err.Error()
		}
	}
	return errors.New(errMessage)
}

// Compare 比较两个数字字符串的大小
func Compare(first, second string) (bool, error) {
	firstNum, err := strconv.ParseInt(first, 10, 64)
	if err != nil {
		return false, fmt.Errorf("第一个参数(%s)转为int64失败(%s)", first, err.Error())
	}
	secondNum, err := strconv.ParseInt(second, 10, 64)
	if err != nil {
		return false, fmt.Errorf("第二个参数(%s)转为int64失败(%s)", second, err.Error())
	}
	return firstNum > secondNum, nil
}
