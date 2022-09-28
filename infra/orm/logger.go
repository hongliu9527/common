/*
 * @Author: hongliu
 * @Date: 2022-09-23 10:11:21
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-23 15:19:58
 * @FilePath: \common\infra\orm\logger.go
 * @Description:gorm 日志封装
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */
package orm

import (
	"context"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// LogConfig 日志配置信息
type LogConfig struct {
	Name               string       // 服务名称
	AccessPath         string       // 访问日志目录
	AccessLogSplitPace int          // 日志文件分割时间（单位：小时）
	AccessLogMaxAge    int          // 日志文件的最长保存时间（单位：小时）
	ConsoleLevel       logrus.Level // 终端日志等级
	FileLevel          logrus.Level // 文件日志等级
}

// OrmLogger 日志定义
type OrmLogger struct {
	consolelogger *logrus.Logger
	fileLogger    *logrus.Logger
}

// NewLogger 获取gorm的logger实例接口
func NewLogger(config LogConfig) logger.Interface {
	var l OrmLogger

	// 终端日志
	l.consolelogger = logrus.New()
	l.consolelogger.SetFormatter(&consoleFormatter)
	// 终端最低等级为Debug，故将Trace等级并入Debug级别
	if config.ConsoleLevel == logrus.DebugLevel {
		l.consolelogger.SetLevel(logrus.TraceLevel)
	} else {
		l.consolelogger.SetLevel(config.ConsoleLevel)
	}
	l.consolelogger.Out = os.Stdout

	// 文件日志
	l.fileLogger = logrus.New()
	l.fileLogger.SetLevel(config.FileLevel)
	filePath := config.AccessPath
	scr, err := os.OpenFile(filePath+"/"+config.Name+".log", os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		l.consolelogger.Error("%s", err.Error())
	}
	l.fileLogger.Out = scr

	logWriter, _ := rotatelogs.New(
		filePath+"/"+config.Name+"-%Y-%m-%d-%H.log",
		rotatelogs.WithMaxAge(time.Duration(config.AccessLogMaxAge)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(config.AccessLogSplitPace)*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.TraceLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	l.fileLogger.AddHook(Hook)
	return &l
}

func (l *OrmLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

// Info 信息级别日志
func (l *OrmLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if len(data) > 0 {
		l.consolelogger.Infof(msg, data...)
		l.fileLogger.Infof(msg, data...)
	} else {
		l.consolelogger.Info(msg)
		l.fileLogger.Info(msg)
	}
}

// Warn 警告级别日志
func (l *OrmLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if len(data) > 0 {
		l.consolelogger.Warnf(msg, data...)
		l.fileLogger.Warnf(msg, data...)
	} else {
		l.consolelogger.Warn(msg)
		l.fileLogger.Warn(msg)
	}
}

// Error 错误级别日志
func (l *OrmLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if len(data) > 0 {
		l.consolelogger.Errorf(msg, data...)
		l.fileLogger.Errorf(msg, data...)
	} else {
		l.consolelogger.Error(msg)
		l.fileLogger.Error(msg)
	}
}

// Trace 追踪级别日志
func (l *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	if rows == -1 {
		l.consolelogger.Debugf("%s", sql)
		l.fileLogger.Tracef("%s", sql)
	} else {
		l.consolelogger.WithField("rows", rows).Debugf("%s", sql)
		l.fileLogger.WithField("rows", rows).Tracef("%s", sql)
	}
}
