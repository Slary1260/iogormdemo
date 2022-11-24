/*
 * @Author: tj
 * @Date: 2022-10-10 09:53:13
 * @LastEditors: tj
 * @LastEditTime: 2022-11-24 17:33:30
 * @FilePath: \iogormdemo\internal\dbgorm\db.go
 */
package dbgorm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Slary1260/iogormdemo/conf"
	"github.com/Slary1260/iogormdemo/internal/dbgorm/model"
	myLogger "github.com/Slary1260/iogormdemo/public/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// InitGormDB InitGormDB
func InitGormDB(cfgs []*conf.DatabaseCfg) (*gorm.DB, error) {
	if len(cfgs) != 1 {
		return ConnectRWDB(cfgs)
	}

	return ConnectSingleDB(cfgs[0])
}

func ConnectSingleDB(cfg *conf.DatabaseCfg) (*gorm.DB, error) {
	out := myLogger.GetOutWriter()

	// 为了处理time.Time，需要包括parseTime作为参数
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", cfg.User, cfg.PassWd, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
	db, err := gorm.Open(mysql.Open(dbDSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger: logger.New(log.New(out, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(200)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = checkDatabaseAndTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 读写分离
func ConnectRWDB(cfgs []*conf.DatabaseCfg) (*gorm.DB, error) {
	out := myLogger.GetOutWriter()

	var instance *gorm.DB
	replicas := []gorm.Dialector{}
	masterDbDSN := ""

	for _, cfg := range cfgs {
		if DbMode(cfg.Mode) == Mastermode {
			masterDbDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", cfg.User, cfg.PassWd, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
			db, err := gorm.Open(mysql.Open(masterDbDSN), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{SingularTable: true},
				Logger: logger.New(log.New(out, "\r\n", log.LstdFlags), logger.Config{
					SlowThreshold:             200 * time.Millisecond,
					LogLevel:                  logger.Silent,
					IgnoreRecordNotFoundError: false,
					Colorful:                  true,
				}),
			})
			if err != nil {
				return nil, err
			}

			instance = db
		} else if DbMode(cfg.Mode) == Slavemode {
			dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", cfg.User, cfg.PassWd, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
			replicas = append(replicas, mysql.New(mysql.Config{DSN: dbDSN}))
		} else {
			return nil, errors.New("unknown mode:" + cfg.Mode)
		}
	}

	plugin := dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.New(mysql.Config{
			DSN: masterDbDSN,
		})},
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	})

	plugin.SetMaxIdleConns(10)
	plugin.SetConnMaxLifetime(time.Hour)
	plugin.SetMaxOpenConns(200)

	instance.Use(plugin)

	err := checkDatabaseAndTables(instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func checkDatabaseAndTables(dbMgr *gorm.DB) error {
	if dbMgr == nil {
		return os.ErrInvalid
	}

	if !dbMgr.Migrator().HasTable(&model.CommConfig{}) {
		dbMgr.Migrator().CreateTable(&model.CommConfig{})
	}

	return nil
}
