package utils

import (
	"errors"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/lee0720/nuwa/pkg/client"
	cfg "gitlab.com/lilh/influx-demo/internal/config"
)

// ErrPerformESBulkRespondBody ...
var ErrPerformESBulkRespondBody = errors.New("Response body contains errors")

var ES *elasticsearch.Client

// InitElasticSearch ...
func InitElasticSearch() {
	log.Println("connecting to es...")

	es, err := client.InitElasticsearch(cfg.CONFIG.ESConfig)
	if err != nil {
		panic(err)
	}
	ES = es
	log.Println("es connected!")
}
