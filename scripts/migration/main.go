package main

import (
	"flag"
	"fmt"
	"strings"

	config "github.com/lee0720/nuwa/pkg/config"
	cfg "gitlab.com/lilh/influx-demo/internal/config"
	"gitlab.com/lilh/influx-demo/models"

	"gitlab.com/lilh/influx-demo/internal/utils"
	"gorm.io/gorm"
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
	tables := []interface{}{
		&models.MarketData{},
	}
	createDeltaTables(utils.SecondaryMarketDB, tables)
	//dropDeltaTables(db, tables)

}

func createDeltaTables(db *gorm.DB, tables []interface{}) {
	err := db.AutoMigrate(tables...)
	if err != nil {
		fmt.Println(err)
	}
}

func dropDeltaTables(db *gorm.DB, tables []interface{}) {
	err := db.Migrator().DropTable(tables...)
	if err != nil {
		fmt.Println(err)
	}
}
