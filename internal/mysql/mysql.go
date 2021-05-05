package mysql

import (
	"fmt"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/lee0720/nuwa/pkg/essentials"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/lilh/influx-demo/internal/utils"
	"gitlab.com/lilh/influx-demo/models"
)

var BatchSize = 300

func SyncMysqlData() error {
	companyNameList, err := getCompanyNameList()
	if err != nil {
		return err
	}
	threadArray := make(chan struct{}, 10)

	marketDatas := make([]models.MarketData, 0)
	errList := make([]error, 0)
	for index, companyName := range companyNameList {
		threadArray <- struct{}{}
		entityID := uuid.NewV4().String()
		go func(Name string, ID string, cur int) {
			defer func() {
				<-threadArray
			}()

			for i := 0; i < 1000; i++ {
				recID := uuid.NewV4().String()
				marketData := models.MarketData{

					BasicModel: essentials.BasicModel{
						RecID: recID,
					},
					EntityID:   ID,
					EntityName: Name,
					PE:         utils.GenFloatNum(0, 100),
					PB:         utils.GenFloatNum(0, 80),
					PS:         utils.GenFloatNum(0, 70),
					DateOn:     time.Now().AddDate(0, 0, -i).Format("20060102"),
				}
				marketDatas = append(marketDatas, marketData)
			}
			err := utils.SecondaryMarketDB.CreateInBatches(marketDatas, BatchSize).Error

			if err != nil {
				errList = append(errList, err)
			}

			fmt.Printf("第 %d 个公司, 名字是 %s 前 3 年数据插入完成\n", cur, Name)

		}(companyName, entityID, index)

	}
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
	// for _, v := range cols[2] {
	// 	fmt.Println(v)
	// }
	return cols[2], nil
}