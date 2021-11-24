package main

import (
	"flag"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
)

// @title Taskman Server New
// @version 1.0
// @description 任务服务管理后台
func main() {
	configFile := flag.String("c", "conf/default.json", "config file path")
	flag.Parse()
	if initConfigMessage := models.InitConfig(*configFile); initConfigMessage != "" {
		fmt.Printf("Init config file error,%s \n", initConfigMessage)
		return
	}
	log.InitLogger()
	if initDbError := db.InitDatabase(); initDbError != nil {
		return
	}
	go db.StartCornJob()
	//start http
	api.InitHttpServer()
}
