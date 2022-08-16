package es

import (
	"fmt"
	"sync"

	"github.com/lee0720/nuwa/pkg/es"
	"gitlab.com/lilh/influx-demo/internal/config"
	"gitlab.com/lilh/influx-demo/internal/utils"
	"gitlab.com/lilh/influx-demo/models"
	"gorm.io/gorm"
)

func SyncES() error {
	var (
		totalCount int64
	)
	err := utils.SecondaryMarketDB.
		Model(&models.MarketData{}).Count(&totalCount).Error
	if err != nil {
		return err
	}

	count := int(totalCount)

	var wg sync.WaitGroup

	errList := make([]error, 0)

	for i := 0; i < count; i += utils.BatchSize {
		wg.Add(1)
		records := make([]models.MarketData, 0)
		err := getMarketDataFromDB(utils.SecondaryMarketDB, i, utils.BatchSize).
			Scan(&records).Error

		if err != nil {
			return err
		}

		go func(index int, recordList []models.MarketData) {

			defer wg.Done()

			enterpriseDocs := MakeDocuments(recordList)
			err := es.CreateInBatches(config.CONFIG.SecondaryMarketIndex.EnterpriseIndex, enterpriseDocs, utils.ES)
			if err != nil {
				errList = append(errList, err)
			}

		}(i, records)

	}
	wg.Wait()

	fmt.Println("SyncES done!")
	return nil
}

func getMarketDataFromDB(db *gorm.DB, limit, offset int) *gorm.DB {
	return db
}

func MakeDocuments(recordList []models.MarketData) []interface{} {
	docs := make([]interface{}, 0)
	for _, v := range recordList {
		doc := makeEnterpriseDocument(v)
		docs = append(docs, doc)
	}
	return docs
}

func makeEnterpriseDocument(record models.MarketData) models.EnterpriseDocument {
	doc := models.EnterpriseDocument{
		CommonDocument: models.CommonDocument{
			CreatedAt:   record.DateOn,
			EntityID:    record.EntityID,
			Keywords:    nil,
			PrimaryName: record.EntityName,
		},
		PE: float32(record.PE),
		PB: float32(record.PB),
		PS: float32(record.PS),
	}

	keywords := make([]string, 0)
	keywords = append(keywords, record.EntityID)
	keywords = append(keywords, record.EntityName)
	doc.Keywords = keywords
	return doc
}
