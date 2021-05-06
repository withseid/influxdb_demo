package utils

import (
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/lee0720/nuwa/pkg/client"
	"gitlab.com/lilh/influx-demo/internal/config"
)

var SecondaryMarketInfluxClient influxdb2.Client

func InitSecondaryMarketInflux() {
	if SecondaryMarketInfluxClient != nil {
		return
	}

	influxClient, err := client.InitInflux(config.CONFIG.InfluxConfig)
	if err != nil {
		panic(err)
	}

	log.Println("InitSecondaryMarketInflux success")
	SecondaryMarketInfluxClient = influxClient

}
