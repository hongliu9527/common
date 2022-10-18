/*
 * @Author: hongliu
 * @Date: 2022-10-17 15:51:23
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-17 15:53:01
 * @FilePath: \common\utils\encryption.go
 * @Description:数据编解码相关工具函数
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// HmacSha1 HmacSha1加密
func HmacSha1(data, secrect string) string {
	mac := hmac.New(sha1.New, []byte(secrect))
	mac.Write([]byte(data))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// MD5WithBase64 先使用MD5加密 然后使用base64加密
func MD5WithBase64(data []byte) string {
	content := md5.Sum(data)
	return base64.StdEncoding.EncodeToString(content[:])
}

// Sha256hash sha256哈希编码
func Sha256hash(data []byte) string {
	s := sha256.New()
	s.Write(data)

	return hex.EncodeToString(s.Sum(nil))
}

// HmacSha256加密
func HmacSha256(data, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return mac.Sum(nil)
}
