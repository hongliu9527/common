/*
 * @Author: hongliu
 * @Date: 2022-09-23 15:21:26
 * @LastEditors: hongliu
 * @LastEditTime: 2022-09-24 16:30:27
 * @FilePath: \common\infra\orm\infra_implemention.go
 * @Description: orm基础设施接口实现
 *
 * Copyright (c) 2022 by 洪流, All Rights Reserved.
 */

package orm

import (
	"context"
	"fmt"
	"time"

	ormConfig "hongliu9527/common/infra/orm/config"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 常量相关定义
const (
	ormInfraName string = "orm" // Gorm基础设施名称
)

func (i *ormInfra) Name() string {
	return ormInfraName
}

// start 启动orm基础设施
func (i *ormInfra) start(ctx context.Context) error {
	// 创建基础设施上下文对象与退出回调函数
	i.ctx, i.cancel = context.WithCancel(ctx)

	// 初始化
	return i.init()
}

// init 初始化orm基础设施
func (i *ormInfra) init() error {
	// 初始化所有gorm句柄
	for _, gormConfig := range i.config.Configs {
		// 初始化gorm句柄
		db, err := connectOneGorm(i.config.LogLevel, i.config.UseExternalHost, gormConfig)
		if err != nil {
			return err
		}

		// 查询该句柄下iot平台相关表名，并添加表名-句柄哈希表
		tableList, err := queryTableNames(db, gormConfig.Type, gormConfig.DatabaseName, gormConfig.TablePrefix)
		if err != nil {
			return err
		}
		for _, tableName := range tableList {
			i.tableNameInstance[tableName] = db
		}

		// 添加实例名-配置信息哈希表
		i.nameConfig[gormConfig.Name] = gormConfig

		// 添加数据库实例名-数据库实例哈希表
		i.nameInstance[gormConfig.Name] = db
	}

	return nil
}

// connectOneGorm 初始化1个gorm连接
func connectOneGorm(level string, useExternalHost bool, config ormConfig.DataBaseConfig) (*gorm.DB, error) {
	var err error

	// 判断是否使用外网地址
	hostPort := config.InternalHostPort
	if useExternalHost {
		hostPort = config.HostPort
	}

	// 转化日志等级参数为gormLogger日志等级
	var logLevel logrus.Level
	switch level {
	case "info":
		logLevel = logrus.InfoLevel
	case "warning", "warn":
		logLevel = logrus.WarnLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		logLevel = logrus.DebugLevel
	}

	// 新建Gorm日志
	newLogger := NewLogger(LogConfig{
		Name:               config.Name,       // 服务名称
		AccessPath:         "./logs/orm/",     // 访问日志目录(相对二进制执行文件的路径)
		AccessLogSplitPace: 24,                // 日志文件分割时间（单位：小时）
		AccessLogMaxAge:    24 * 30,           // 日志文件的最长保存时间（单位：小时）
		ConsoleLevel:       logLevel,          // 终端日志等级
		FileLevel:          logrus.TraceLevel, // 文件日志等级
	})

	// 生成数据源名称
	var (
		dsn       string
		dialector gorm.Dialector
	)
	switch config.Type {
	case "mysql", "tidb":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?%s&timeout=%ds",
			config.Username,
			config.Password,
			hostPort,
			config.DatabaseName,
			"charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai",
			config.ConnectTimeout)

		dialector = mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN配置信息
			DefaultStringSize:         256,   // string类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用datetime精度，因为MySQL5.6之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并创建的方式，因为MySQL5.7之前的数据库和MariaDB不支持重命名索引
			DontSupportRenameColumn:   true,  // 用`change`重命名列，因为MySQL8之前的数据库和MariaDB不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前MySQL版本自动配置
		})

	case "clickhouse":
		dsn = fmt.Sprintf("tcp://%s?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20",
			hostPort,
			config.DatabaseName,
			config.Username,
			config.Password)

		dialector = clickhouse.New(clickhouse.Config{
			DSN:                       dsn,   // DSN配置信息
			DisableDatetimePrecision:  true,  // 禁用datetime精度
			DontSupportRenameColumn:   true,  // 用`change`重命名列
			SkipInitializeWithVersion: false, // 根据当前版本自动配置
		})
	}

	if dsn == "" || dialector == nil {
		return nil, errors.New("数据库基础配置信息缺少")
	}

	db, err := gorm.Open(dialector,
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 禁用AutoMigrate自动创建数据库外键约束
			Logger:                                   newLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   config.TablePrefix, // 表名前缀
				SingularTable: true,               // 使用单数表名
			},
		})
	if err != nil {
		return nil, errors.Wrapf(err, "初始化Gorm(%s)失败", config.Name)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrapf(err, "获取Gorm(%s)连接指针失败", config.Name)
	}

	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用最大时间

	return db, nil
}

// queryTableNames 查询Gorm实例下所有表名
func queryTableNames(db *gorm.DB, dataBaseType string, dataBaseName string, tablePrefix string) ([]string, error) {
	tableNameList := make([]string, 0, 0)
	tabelPrefix := fmt.Sprintf("%s%%", tablePrefix)

	var querySQL string
	switch dataBaseType {
	case "mysql", "tidb":
		querySQL = "select table_name from information_schema.tables where table_schema = ? and table_name like ?"
	case "clickhouse":
		querySQL = "select DISTINCT(name) from system.tables where database = ? and name like ?"
	default:
		return nil, fmt.Errorf("数据库类型未知(%s)", dataBaseType)
	}
	db.Raw(querySQL, dataBaseName, tabelPrefix).Scan(&tableNameList)
	if db.Error != nil {
		return nil, fmt.Errorf("查询数据库(%s)的表名列表失败(%s)", dataBaseName, db.Error)
	}
	return tableNameList, nil
}

// stop 关闭Gorm基础设施
func (i *ormInfra) stop() error {
	for name, db := range i.nameInstance {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("获取数据库连接(%s)失败(%s)", name, err.Error())
		}
		err = sqlDB.Close()
		if err != nil {
			return fmt.Errorf("关闭数据库(%s)失败(%s)", name, err.Error())
		}
	}

	i.nameConfig = make(map[string]ormConfig.DataBaseConfig)
	i.tableNameInstance = make(map[string]*gorm.DB)
	i.nameInstance = make(map[string]*gorm.DB)

	// 执行退出回调函数
	i.cancel()

	return nil
}
