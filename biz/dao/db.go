//Package dao
/**
get:
	1. if not record, return errno.NotExists(if ignore not record use Find not Take or First)
exists:
	1. has->true, hasn't->false, if db err return false and err
update:
	1. need check RowsAffected
*/
package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"tiktok/biz/config"
	"time"
)

var Db *gorm.DB

func Init(logLevel logger.LogLevel, logWriter logger.Writer) {
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.C.Mysql.Username,
		config.C.Mysql.Password,
		config.C.Mysql.Host,
		config.C.Mysql.Port,
		config.C.Mysql.Database,
	)
	Db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		//Logger: logger.Default.LogMode(logLevel),
		Logger: logger.New(logWriter, logger.Config{
			LogLevel:                  logLevel,
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			Colorful:                  logLevel == logger.Info,
		}),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名
			SingularTable: true,
		},
	})

	if err != nil {
		panic(fmt.Sprintf("连接数据库失败，请检查参数：%s", err.Error()))
	}

	// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
	//_ = Db.AutoMigrate(&Video{}, &User{}, &Favorite{}, &Follow{}, &Comment{})

	sqlDB, _ := Db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(10)

	// SetConnMaxLifetime 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

}
