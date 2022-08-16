package mysql

import (
	"fmt"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/lilh/influx-demo/internal/utils"
	"gitlab.mvalley.com/datapack/cain/pkg/bedrock"
	"gitlab.mvalley.com/datapack/cain/pkg/essentials"
	"gitlab.mvalley.com/secondary-market-datapack/a-share/pkg/models/security/market_data"
)

var BatchSize = 5000

func SyncMysqlData() error {
	companyNameList, err := getCompanyNameList()
	if err != nil {
		return err
	}

	var sema = utils.NewSemaPhore(20)

	errList := make([]error, 0)
	for index, companyName := range companyNameList {
		if index == 0 {
			continue
		}

		marketDatas := make([]market_data.MarketData, 0)
		sema.Add(1)

		go func(Name string, cur int) {
			defer func() {
				sema.Done()
			}()

			entityID := uuid.NewV4().String()
			for i := 0; i < 30; i++ {

				dateOn := time.Now().AddDate(0, 0, -i).Format("20060102")

				pe := utils.GenFloatNum(0, 100)
				pb := utils.GenFloatNum(0, 80)
				ps := utils.GenFloatNum(0, 70)
				pcf := utils.GenFloatNum(0, 70)

				marketData := market_data.MarketData{
					BasicModel: essentials.BasicModel{
						RecID: uuid.NewV4().String(),
					},

					SecurityID: entityID,

					CurrencyWindID: "0",

					PE:  pe,
					PB:  pb,
					PS:  ps,
					PCF: pcf,

					TransactionDate:           bedrock.Date(dateOn),
					TransactionDateConfidence: 0,
				}
				marketDatas = append(marketDatas, marketData)
			}
			err := utils.SecondaryMarketDB.CreateInBatches(marketDatas, BatchSize).Error

			if err != nil {
				errList = append(errList, err)
			}

			fmt.Printf("第 %d 个公司, 名字是 %s 前 3 年数据插入完成\n", cur, Name)

		}(companyName, index)

	}
	sema.Wait()
	return utils.ErrListToError(errList)

}

func getCompanyNameList() ([]string, error) {
	f, err := excelize.OpenFile("./data1.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Get all the rows in the Sheet1.
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cols[2], nil
}
