package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	config "github.com/lee0720/nuwa/pkg/config"
	cfg "gitlab.com/lilh/influx-demo/internal/config"
	"gitlab.com/lilh/influx-demo/internal/utils"
)

// const token = "p0_sYwo-A72-AVyqGxjtrfINB52iyFRmhJayigj70G-bZvNTC2lCxH1SvLBwO4-ayHRfv-61D9YyXrYcDN2HTg=="

// const bucket = "secondary_market"
// const bucket = "test1"
// const org = "rime"

type Node struct {
	Start string
	Stop  string
	First float64
	Last  float64
	Avg   float64
	Max   float64
	Min   float64
}

var configFileName = flag.String("cfn", "config", "name of config file")
var configFilePath = flag.String("cfp", "./configs/", "path of config file")

func main() {
	flag.Parse()
	err := config.InitConfiguration(*configFileName, strings.Split(*configFilePath, ","), &cfg.CONFIG)
	if err != nil {
		panic(err)
	}
	utils.InitSecondaryMarketInflux()

	defer utils.SecondaryMarketInfluxClient.Close()

	queryData4(utils.SecondaryMarketInfluxClient)

	// writeData(utils.SecondaryMarketInfluxClient)

	// DeleteData(utils.SecondaryMarketInfluxClient)
	// queryData(client)
	// queryData2(client)
	// date := "20190101"
	// cur, err := toTrueTime(date)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(cur)
}

func queryData4(client influxdb2.Client) {

	query := `

	dataset = from(bucket: "secondary_market")
	|> range(start: -3y)
	|> filter(fn: (r) => r["_measurement"] == "market_data_v1")
	|> filter(fn: (r) => r["_field"] == "pe")
	|> filter(fn: (r) => r["security_entity_id"] == "01cc2c54-089b-494b-9c20-96bfd3028206")
	|> drop(columns: ["_field", "_measurement"])
	|> window(every: 1mo)
  
	dataset
	|> mean()  
	`

	// Get query client
	queryAPI := client.QueryAPI(cfg.CONFIG.InfluxConfig.Org)

	// nodes := make([]Node, 0)

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	if err != nil {
		return
	}

	count := 0

	for result.Next() {
		count++
	}
	fmt.Println(count)
}

func queryData3(client influxdb2.Client) {

	query := `

	dataset = from(bucket: "secondary_market")
	|> range(start: -3y)
	|> filter(fn: (r) => r["_measurement"] == "market_data_v1")
	|> filter(fn: (r) => r["_field"] == "pe")
	|> filter(fn: (r) => r["security_entity_id"] == "01cc2c54-089b-494b-9c20-96bfd3028206")
	|> drop(columns: ["_field", "_measurement"])
	|> window(every: 3mo)
  
  avg = dataset
	|> mean()
  
  first = dataset
	|> first()
  
  last = dataset
	|> last()
  
  max = dataset
	|> drop(columns: ["_time"])
	|> max()
  
  min = dataset
	|> drop(columns: ["_time"])
	|> min()
  
  temp1 = join(tables: {d1: first,d2: last}, on :["security_entity_id","_start","_stop"])
	|> rename(columns: {_value_d1: "first", _value_d2: "last"})
  
  temp2 = join(tables: {d3:max,d4:min}, on: ["security_entity_id","_start","_stop"])
  |> rename(columns: {_value_d3: "max", _value_d4: "min"})
  
  temp3 = join(tables:{t1: temp2,d5:avg}, on:["security_entity_id","_start","_stop"])
  |> rename(columns: {_value: "avg"})
  
  temp4 = join(tables:{t4:temp1,t3:temp3}, on:["security_entity_id","_start","_stop"])
  |> rename(columns: {_time_d1: "start", _time_d2: "stop"})
  |> drop(columns: ["_start","_stop"])
  
  temp4
	`

	// Get query client
	queryAPI := client.QueryAPI(cfg.CONFIG.InfluxConfig.Org)

	nodes := make([]Node, 0)

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	if err != nil {
		return
	}

	for result.Next() {
		node := Node{
			Start: result.Record().ValueByKey("start").(time.Time).Format("2006-01-02"),
			Stop:  result.Record().ValueByKey("stop").(time.Time).Format("2006-01-02"),
			First: toFixed(result.Record().ValueByKey("first").(float64), 2),
			Last:  toFixed(result.Record().ValueByKey("last").(float64), 2),
			Min:   toFixed(result.Record().ValueByKey("min").(float64), 2),
			Max:   toFixed(result.Record().ValueByKey("max").(float64), 2),
			Avg:   toFixed(result.Record().ValueByKey("avg").(float64), 2),
		}
		nodes = append(nodes, node)
	}
	for _, node := range nodes {
		fmt.Printf("node: %+v\n", node)
	}
}

func writeData(client influxdb2.Client) {

	writeAPI := client.WriteAPIBlocking(cfg.CONFIG.InfluxConfig.Org, cfg.CONFIG.InfluxConfig.Bucket)
	start := time.Now()
	threadArray := make(chan int, 10)

	for i := 0; i < 20; i++ {

		day := int(time.Now().AddDate(0, 0, -i).Weekday())
		if day == 0 || day == 6 {
			continue
		}

		threadArray <- i
		go func(index int) {

			defer func() {
				<-threadArray
			}()

			points := make([]*write.Point, 0)
			for j := 10000; j < 10010; j++ {
				p := influxdb2.NewPointWithMeasurement("market_data_v1").
					AddTag("company_id", strconv.Itoa(j)).
					AddField("pe", genFloatNum(1, 40)).
					AddField("pb", genFloatNum(1, 70)).
					AddField("ps", genFloatNum(1, 70)).
					AddField("pa", genFloatNum(1, 90)).
					SetTime(time.Now().AddDate(0, 0, -index))
					// Setime
				points = append(points, p)
			}
			fmt.Printf("day %d insert into data start\n", index)
			fmt.Printf("total count: %+v\n", len(points))
			writeAPI.WritePoint(context.Background(), points...)
			fmt.Printf("day %d insert into data end\n", index)

		}(i)
	}

	elapsed := time.Since(start)
	fmt.Println("插入10000 * 7000 条数据花费时间:", elapsed)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func queryData(client influxdb2.Client) {

	query := `
	dataset = from(bucket: "market_data")
	|> range(start: -1y)
	|> filter(fn: (r) => r["_measurement"] == "market_data_v2" and r["company_id"] == "10005")
	|> filter(fn: (r) => r["_field"] == "pe")
	|> drop(columns: ["_field", "_measurement"])
	|> window(every: 1w, period: 5d, offset: 2d)
  
  avg = dataset
	|> mean()
	|> set(key: "sign", value: "avg")
  
  fi = dataset
	|> first()
	|> set(key: "sign", value: "fi")
  
  la = dataset
	|> last()
	|> set(key: "sign", value: "la")
  
  ma = dataset
	|> max()
	|> set(key: "sign", value: "max")
  
  mi = dataset
	|> min()
	|> set(key: "sign", value: "min")
  
  union(tables: [avg, fi, la, ma, mi])
	|> group(columns: ["sign"])
	`

	// Get query client
	queryAPI := client.QueryAPI(cfg.CONFIG.InfluxConfig.Org)

	maps := make(map[string]*Node, 0)

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	if err != nil {
		return
	}

	for result.Next() {
		start := ""
		stop := ""
		key1 := result.Record().Start().Format("2006-01-02")
		key2 := result.Record().Stop().AddDate(0, 0, -1).Format("2006-01-02")

		sign := result.Record().ValueByKey("sign")
		if sign == "first" {
			start = result.Record().Time().Format("2006-01-02")
		}
		if sign == "last" {
			stop = result.Record().Time().Format("2006-01-02")
		}
		key := fmt.Sprintf("%s_%s", key1, key2)
		value := result.Record().Value()

		if node, ok := maps[key]; ok {
			AssignmentNode(sign, value, node)
			if node.Start == "" {
				node.Start = start
			}
			if node.Stop == "" {
				node.Stop = stop
			}
			maps[key] = node
		} else {
			node := new(Node)
			if node.Start == "" {
				node.Start = start
			}
			if node.Stop == "" {
				node.Stop = stop
			}
			AssignmentNode(sign, value, node)
			maps[key] = node
		}

	}

	nodes := make([]Node, 0)
	for _, v := range maps {
		nodes = append(nodes, *v)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Start < nodes[j].Start
	})

	for _, node := range nodes {
		fmt.Printf("node: %+v\n", node)
	}

}

func genFloatNum(min, max int) float64 {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return float64(randNum)
}

func genIntNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

func AssignmentNode(sign interface{}, value interface{}, node *Node) {
	switch sign {
	case "avg":
		node.Avg = value.(float64)
	case "first":
		node.First = value.(float64)
	case "last":
		node.Last = value.(float64)
	case "min":
		node.Min = value.(float64)
	case "max":
		node.Max = value.(float64)
	default:

	}
}

func toTrueTime(date string) (time.Time, error) {
	timeString := "20060102"
	loc, _ := time.LoadLocation("Local")
	if len(date) == 8 {
		return time.ParseInLocation(timeString, string(date), loc)
	}
	if len(date) == 6 {
		return time.ParseInLocation(timeString[:6], string(date), loc)
	}
	if len(date) == 4 {
		return time.ParseInLocation(timeString[:4], string(date), loc)
	}
	return time.Time{}, errors.New("time transcode error")
}

func queryData2(client influxdb2.Client) {
	startTime := int64(1557646863)
	endTime := int64(1620891663)
	fmt.Println(utils.UnixFormat(startTime), utils.UnixFormat(endTime))
	// a := 1
	// query := fmt.Sprintf(`
	// from(bucket: "secondary_market")
	// 	|> range(start: %s, stop: %s)
	// 	|> filter(fn: (r) => r["_measurement"] == "market_data_v1")
	// 	|> filter(fn: (r) => r["_field"] == "pb")
	// 	|> filter(fn: (r) => r["security_entity_id"] == "01cc2c54-089b-494b-9c20-96bfd3028206")`, utils.UnixFormat(startTime), utils.UnixFormat(endTime))

	query := fmt.Sprintf(`
	from(bucket: "secondary_market")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "market_data_v1")
		|> filter(fn: (r) => r["_field"] == "pb")
		|> filter(fn: (r) => r["security_entity_id"] == "01cc2c54-089b-494b-9c20-96bfd3028206")
		|> mean()
		`, utils.UnixFormat(startTime), utils.UnixFormat(endTime))

	// if a > 0 {
	// 	query = query + fmt.Sprintf(" |> filter(fn: (r) => r[\"_value\"] > 50)")
	// }

	// Get query client
	queryAPI := client.QueryAPI(cfg.CONFIG.InfluxConfig.Org)

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	if err != nil {
		log.Println(err)
		return
	}

	for result.Next() {
		fmt.Printf("field: %+v, value: %+v\n", result.Record().Field(), Decimal(result.Record().Value().(float64)))
		// fmt.Printf("time: %+v,field: %+v, value: %+v\n", result.Record().Time().Format("2006-01-02"), result.Record().Field(), result.Record().Value())
	}
}

func Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

func DeleteData(client influxdb2.Client) {

	client.BucketsAPI().CreateBucketWithNameWithID(context.Background(), cfg.CONFIG.InfluxConfig.Org, cfg.CONFIG.InfluxConfig.Bucket)

	client.Options().HTTPOptions().HTTPClient().Timeout = 5 * time.Minute

	start := time.Now()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))
	defer cancel()

	go doStuff(ctx)
	deleteAPI := client.DeleteAPI()

	err := deleteAPI.DeleteWithName(ctx, cfg.CONFIG.InfluxConfig.Org, cfg.CONFIG.InfluxConfig.Bucket, time.Now().AddDate(-20, 0, 0), time.Now(), "_measurement=market_data_v3")
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		panic(err)
	}

	duration := time.Since(start)
	fmt.Printf("duration: %+v\n", duration)

}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
func doStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			log.Printf("done")
			return
		default:
			log.Printf("work")
		}
	}
}

func writeData2(client influxdb2.Client) {
	// dates := []string{"20210505", "20210506"}
	writeAPI := client.WriteAPIBlocking(cfg.CONFIG.InfluxConfig.Org, cfg.CONFIG.InfluxConfig.Bucket)
	// points := make([]*write.Point, 0)
	// for i := 0; i < len(dates); i++ {
	// 	toTrueDate, _ := toTrueTime(dates[i])
	// 	p := influxdb2.NewPointWithMeasurement("market_data_v6").
	// 		AddTag("company_id", "10001").
	// 		AddField("pe", i+2).
	// 		AddField("pb", i+2).
	// 		AddField("ps", i+2).
	// 		AddField("pa", i+2).
	// 		SetTime(toTrueDate)
	// 	points = append(points, p)
	// }
	// for i := 0; i < len(points); i++ {
	// 	fmt.Println(points[i].Time())
	// }
	// writeAPI.WritePoint(context.Background(), points...)
	start := time.Now()
	threadArray := make(chan int, 10)

	for i := 0; i < 20; i++ {

		day := int(time.Now().AddDate(0, 0, -i).Weekday())
		if day == 0 || day == 6 {
			continue
		}

		threadArray <- i
		go func(index int) {

			defer func() {
				<-threadArray
			}()

			points := make([]*write.Point, 0)
			for j := 10000; j < 10010; j++ {
				p := influxdb2.NewPointWithMeasurement("market_data_v6").
					AddTag("company_id", strconv.Itoa(j)).
					AddField("pe", genFloatNum(1, 40)).
					AddField("pb", genFloatNum(1, 70)).
					AddField("ps", genFloatNum(1, 70)).
					AddField("pa", genFloatNum(1, 90)).
					SetTime(time.Now().AddDate(0, 0, -index))
					// Setime
				points = append(points, p)
			}
			fmt.Printf("day %d insert into data start\n", index)
			fmt.Printf("total count: %+v\n", len(points))
			writeAPI.WritePoint(context.Background(), points...)
			fmt.Printf("day %d insert into data end\n", index)

		}(i)
	}

	elapsed := time.Since(start)
	fmt.Println("插入10000 * 7000 条数据花费时间:", elapsed)

}
