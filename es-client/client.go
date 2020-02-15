package client

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic"
)

var ctx = context.Background()

// CreateClient creates a new client
func CreateClient() *elastic.Client {
	client, connError := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)

	if connError != nil {
		panic(connError)
	}

	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	return client
}
