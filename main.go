package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

const token = "s4fe-AJ0c6uMfK_7lhCNelgkVDKuPDUdKhaaGnsWDYWCtVRLwSi4jNUWu3fVPqwxO9SIFTOt1DimTe4NWvPmBQ=="
const bucket = "market_data"
const org = "rime"

type Node struct {
	Start string
	Stop  string
	First float64
	Last  float64
	Avg   float64
	Max   float64
	Min   float64
}

func main() {

	client := influxdb2.NewClient("http://192.168.88.11:8286/", token)

	defer client.Close()
	writeData(client)
	// queryData(client)

}

func writeData(client influxdb2.Client) {

	writeAPI := client.WriteAPIBlocking(org, bucket)
	start := time.Now()
	threadArray := make(chan int, 10)

	for i := 0; i < 3600; i++ {

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
			for j := 10000; j < 20000; j++ {
				p := influxdb2.NewPointWithMeasurement("market_data_v2").
					AddTag("company_id", strconv.Itoa(j)).
					AddField("pe", genFloatNum(1, 40)).
					AddField("pb", genFloatNum(1, 70)).
					AddField("ps", genFloatNum(1, 70)).
					SetTime(time.Now().AddDate(0, 0, -index))
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

func queryData(client influxdb2.Client) {

	query := `
	dataset = from(bucket: "market_data")
	|> range(start: -1y)
	|> filter(fn: (r) => r["_measurement"] == "market_data_v1" and r["_field"] == "pe" and r["company_id"] == "10202")
	|> window(every: 1w, period: 5d, offset: 2d)
	
	avg = dataset
	|> mean()
	|> set(key: "sign", value: "avg")
	|> drop(columns: ["_field", "_measurement", "entity_id", "entity_type"])
	
	first = dataset
	|> drop(columns: ["_field", "_measurement"])
	|> first()
	|> set(key: "sign", value: "first")
	
	last = dataset
	|> drop(columns: ["_field", "_measurement"])
	|> last()
	|> set(key: "sign", value: "last")
	
	max = dataset
	|> drop(columns: ["_time", "_field", "_measurement"])
	|> max()
	|> set(key: "sign", value: "max")
	
	min = dataset
	|> drop(columns: ["_time", "_field", "_measurement"])
	|> min()
	|> set(key: "sign", value: "min")
	
	union(tables: [avg, first, last, max, min])
	|> group(columns: ["sign"])
	`
	// Get query client
	queryAPI := client.QueryAPI(org)

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
