/*
 * @Author: tj
 * @Date: 2022-11-23 15:26:58
 * @LastEditors: tj
 * @LastEditTime: 2022-11-24 17:33:38
 * @FilePath: \iogormdemo\cmd\main.go
 */
package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/Slary1260/iogormdemo/conf"
	"github.com/Slary1260/iogormdemo/internal/dbgorm"
	"github.com/Slary1260/iogormdemo/internal/dbgorm/model"
	"github.com/Slary1260/iogormdemo/public/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/plugin/dbresolver"
)

var (
	// -c=true
	cfgFile = flag.String("c", "server.toml", "specified the config file path")

	log = logrus.WithFields(logrus.Fields{
		"main": "",
	})
)

func flagParse() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Errorf("Error: can not reach the work directory. %v.\n", err)
		os.Exit(1)
	}

	if !filepath.IsAbs(*cfgFile) {
		newCfg := filepath.Join(wd, *cfgFile)
		cfgFile = &newCfg
	}
}

func main() {
	logger.DefaultLogrusLogger()
	logrus.SetLevel(logrus.InfoLevel)

	flagParse()

	defer func() {
		if err := recover(); err != nil {
			log.Errorln("main recover error:", err)
		}
	}()

	cfg, err := conf.GetConfigInfo(*cfgFile)
	if err != nil {
		log.Error("GetConfigInfo error:", err)
		os.Exit(1)
	}

	// mysql 数据库
	gormDB, err := dbgorm.InitGormDB(cfg.DatabaseConfigs)
	if err != nil {
		log.Error("InitGormDB error:", err)
		os.Exit(1)
	}

	queryCfg := &model.CommConfig{ConfigID: "name"}
	err = gormDB.Where(queryCfg).First(&queryCfg).Error
	if err != nil {
		log.Errorln("queryCfg error:", err)
		return
	}

	// 使用主库写数据
	insertCfg := &model.CommConfig{ConfigID: "name1"}
	err = gormDB.Clauses(dbresolver.Write).Create(insertCfg).Error
	if err != nil {
		log.Errorln("insertCfg error:", err)
		return
	}

	queryAllCfgs := make([]*model.CommConfig, 0, 8)
	err = gormDB.Where(&model.CommConfig{}).Find(&queryAllCfgs).Error
	if err != nil {
		log.Errorln("queryAllCfgs error:", err)
		return
	}

	for _, v := range queryAllCfgs {
		log.Infoln(v)
	}
}
