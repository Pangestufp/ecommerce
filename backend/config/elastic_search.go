package config

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

func ConnectElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", ENV.ElasticDBHost, ENV.ElasticDBPort),
		},
		Username: ENV.ElasticDBUser,
		Password: ENV.ElasticDBPassword,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("Cannot connect to Elasticsearch:", err)
	}

	res, err := client.Info()
	if err != nil {
		log.Fatal("Elasticsearch info error:", err)
	}
	defer res.Body.Close()

	ESClient = client
	log.Println("Connected to Elasticsearch at", ENV.ElasticDBPort)
}
