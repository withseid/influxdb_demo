package influx

import (
	"context"
	"log"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	cfg "gitlab.com/lilh/influx-demo/internal/config"
	"gitlab.com/lilh/influx-demo/internal/utils"
	"gitlab.com/lilh/influx-demo/models"
	"gorm.io/gorm"
)

func SyncInflux() error {
	var totocalCount int64

	err := utils.SecondaryMarketDB.Model(&models.MarketData{}).
		Count(&totocalCount).Error
	if err != nil {
		return err
	}

	count := int(totocalCount)

	var wg sync.WaitGroup

	errList := make([]error, 0)

	for i := 0; i < count; i += utils.BatchSize {
		wg.Add(1)
		records := make([]models.MarketData, 0)
		err := getMarketData(utils.SecondaryMarketDB, i, utils.BatchSize).
			Scan(&records).Error

		if err != nil {
			return err
		}

		go func(index int, dataArray []models.MarketData) {

			defer wg.Done()

			points := make([]*write.Point, 0)
			for _, val := range dataArray {

				dateOn, err := utils.ToTrueTime(val.DateOn)

				if err != nil {
					errList = append(errList, err)
					continue
				}

				p := influxdb2.NewPointWithMeasurement("market_data_v1").
					AddTag("entity_id", val.EntityID).
					AddField("pe", val.PE).
					AddField("pb", val.PB).
					AddField("ps", val.PB).
					SetTime(dateOn.Add(15 * time.Hour))
				points = append(points, p)
			}

			err := utils.SecondaryMarketInfluxClient.
				WriteAPIBlocking(cfg.CONFIG.InfluxConfig.Org, cfg.CONFIG.InfluxConfig.Bucket).
				WritePoint(context.Background(), points...)

			if err != nil {
				errList = append(errList, err)
			}

			log.Printf("插入数据 %d - %d. Total: %d\n", index, index+len(points), totocalCount)

		}(i, records)

	}

	wg.Wait()
	if errList != nil {
		return utils.ErrListToError(errList)
	}
	return nil

}

func getMarketData(db *gorm.DB, offset, limit int) *gorm.DB {
	return db.Raw(`select * from market_data where rec_id >= 
	(select rec_id from market_data order by rec_id  limit ?,1) 
	order by rec_id limit ?`, offset, limit)
}
