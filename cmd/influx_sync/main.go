package main

import (
	"flag"
	"fmt"
	"strings"

	config "github.com/lee0720/nuwa/pkg/config"
	cfg "gitlab.com/lilh/influx-demo/internal/config"
	"gitlab.com/lilh/influx-demo/internal/influx"
	"gitlab.com/lilh/influx-demo/internal/utils"
)

var configFileName = flag.String("cfn", "config", "name of config file")
var configFilePath = flag.String("cfp", "./configs/", "path of config file")

func main() {
	flag.Parse()
	err := config.InitConfiguration(*configFileName, strings.Split(*configFilePath, ","), &cfg.CONFIG)
	if err != nil {
		panic(err)
	}
	utils.InitSecondaryMarketMysql()
	utils.InitSecondaryMarketInflux()

	err = influx.SyncInflux()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}

	fmt.Println("influx 插入数据完成")
}
