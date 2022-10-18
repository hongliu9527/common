/*
 * @Author: hongliu
 * @Date: 2022-09-23 15:21:26
 * @LastEditors: hongliu
 * @LastEditTime: 2022-10-17 10:28:09
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

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// 常量相关定义
const (
	ormInfraName string = "orm" // Orm基础设施名称
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
	// 初始化所有sqlx句柄
	for _, sqlxConfig := range i.config.Configs {
		// 初始化sqlx句柄
		db, err := connectOneSqlx(i.config.LogLevel, i.config.UseExternalHost, sqlxConfig)
		if err != nil {
			return err
		}

		// 查询该句柄下iot平台相关表名，并添加表名-句柄哈希表
		tableList, err := queryTableNames(db, sqlxConfig.Type, sqlxConfig.DatabaseName, sqlxConfig.TablePrefix)
		if err != nil {
			return err
		}
		for _, tableName := range tableList {
			i.tableNameInstance[tableName] = db
		}

		// 添加实例名-配置信息哈希表
		i.nameConfig[sqlxConfig.Name] = sqlxConfig

		// 添加数据库实例名-数据库实例哈希表
		i.nameInstance[sqlxConfig.Name] = db
	}

	return nil
}

// connectOneSqlx 初始化1个sqlx连接
func connectOneSqlx(level string, useExternalHost bool, config ormConfig.DataBaseConfig) (*sqlx.DB, error) {

	// 判断是否使用外网地址
	hostPort := config.InternalHostPort
	if useExternalHost {
		hostPort = config.HostPort
	}

	// 生成数据源名称
	var (
		dsn string
		db  *sqlx.DB
		err error
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
		db, err = sqlx.Open("mysql", dsn)
	case "clickhouse":
		dsn = fmt.Sprintf("tcp://%s?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20",
			hostPort,
			config.DatabaseName,
			config.Username,
			config.Password)
		db, err = sqlx.Open("clickhouse", dsn)
	}

	if dsn == "" {
		return nil, errors.New("数据库基础配置信息缺少")
	}

	if err != nil {
		return nil, errors.Wrapf(err, "初始化%s-%s连接失败", config.Name, config.Type)
	}

	// 确认sqlx实例正常
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "初始化%s-%s连接失败", config.Name, config.Type)
	}

	db.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	db.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	db.SetConnMaxLifetime(time.Hour) // 设置连接可复用最大时间

	return db, nil
}

// queryTableNames 查询sqlx实例下所有表名
func queryTableNames(db *sqlx.DB, dataBaseType string, dataBaseName string, tablePrefix string) ([]string, error) {
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

	rows, err := db.Query(querySQL, dataBaseName, tabelPrefix)
	if err != nil {
		return nil, fmt.Errorf("查询数据库(%s)的表名列表失败(%s)", dataBaseName, err)
	}
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tableNameList = append(tableNameList, tableName)
	}

	return tableNameList, nil
}

// stop 关闭orm基础设施
func (i *ormInfra) stop() error {
	for name, db := range i.nameInstance {
		err := db.Close()
		if err != nil {
			return fmt.Errorf("关闭数据库(%s)失败(%s)", name, err.Error())
		}
	}

	i.nameConfig = make(map[string]ormConfig.DataBaseConfig)
	i.tableNameInstance = make(map[string]*sqlx.DB)
	i.nameInstance = make(map[string]*sqlx.DB)

	// 执行退出回调函数
	i.cancel()

	return nil
}
