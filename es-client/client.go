package client

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

var ctx = context.Background()

// CreateClient creates a new client
func CreateClient(host *string) *elastic.Client {
	client, connError := elastic.NewClient(
		elastic.SetURL(*host),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)

	if connError != nil {
		panic(connError)
	}

	info, code, err := client.Ping(*host).Do(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	return client
}
