package utils

import (
	"log"

	"github.com/lee0720/nuwa/pkg/client"
	"gitlab.com/lilh/influx-demo/internal/config"
	"gorm.io/gorm"
)

var BatchSize = 200 * 10
var SecondaryMarketDB *gorm.DB

func InitSecondaryMarketMysql() {
	if SecondaryMarketDB != nil {
		return
	}

	db, err := client.InitGormV2(config.CONFIG.SecondaryMarketConfig)
	if err != nil {
		panic(err)
	}

	log.Println("InitSecondaryMarketDB success")
	SecondaryMarketDB = db
}
