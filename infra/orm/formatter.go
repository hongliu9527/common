/*
 * @Author: hongliu
 * @Date: 2022-09-23 10:10:16
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-23 10:10:21
 * @FilePath: \common\infra\orm\formatter.go
 * @Description: orm封装日志
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package orm

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// consoleFormatter 终端日志样式
var consoleFormatter = Formatter{
	TimestampFormat: "2006-01-02 15:04:05.000", // 设置时间戳格式
	NoFieldsColors:  false,                     // 设置单行所有内容都显示颜色
	HideKeys:        true,                      // 是否隐藏键值
}

// 对应终端字体颜色
const (
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
	colorBlue   = 36
)

// getColorByLevel 定义不同日志等级颜色
func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorGreen
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

// Formatter 自定义格式控制选项
type Formatter struct {
	FieldsOrder     []string // default: fields sorted alphabetically
	TimestampFormat string   // default: time.StampMilli = "Jan _2 15:04:05.000"
	HideKeys        bool     // show [fieldValue] instead of [fieldKey:fieldValue]
	NoColors        bool     // disable colors
	NoFieldsColors  bool     // color only level, default is level + fields
	ShowFullLevel   bool     // true to show full level [WARNING] instead [WARN]
	TrimMessages    bool     // true to trim whitespace on messages
}

// Format 自定义格式函数
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	// 设置时间字符串格式
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	// 按照日志等级填充ASC-II颜色特殊符号
	levelColor := getColorByLevel(entry.Level)
	level := strings.ToUpper(entry.Level.String())
	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	// 填充日志等级，使用"[]"包裹
	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("] ")

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// 填充时间格式字符串，使用"[]"包裹
	b.WriteString("[")
	b.WriteString(entry.Time.Format(timestampFormat))
	b.WriteString("] ")

	// 填充消息字段，用于结构化日志
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	// 填充日志内容
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}
	b.WriteByte('\n')

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	return b.Bytes(), nil
}

// writeFields 写入每个日志字段
func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

// writeOrderedFields 按顺序写入日志字段
func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if !foundFieldsMap[field] {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v] ", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v] ", field, entry.Data[field])
	}
}
