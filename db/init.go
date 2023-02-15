package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	DB *gorm.DB
)

func InitDB() {
	var err error

	dsn := "user:pass@tcp(127.0.0.1:3306)/Doushengdb?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Gorm Open Error", err.Error())
	}

	// AutoMigrate会创建表，缺失的外键，约束，列和索引。如果大小，精度，是否为空，可以更改，则AutoMigrate会改变列的类型。出于保护您数据的目的，它不会删除未使用的列
	// 刷新数据库的表格，使其保持最新。即如果我在旧表的基础上增加一个字段age，那么调用autoMigrate后，旧表会自动多出一列age，值为空
	if err := DB.AutoMigrate(&User{}, &Video{}, &Comment{}, &Relation{}); err != nil {
		log.Panic("DB AutoMigrate Error", err.Error())
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := DB.DB()
	if err != nil {
		log.Panic("sqlDB Open Error", err.Error())
	}

	if err := sqlDB.Ping(); err != nil {
		log.Panic("sqlDB Ping Error", err.Error())
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(200)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(20)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
