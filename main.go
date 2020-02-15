package main

import (
	"context"

	client "github.com/albertlockett/learn-baseball-db/es-client"
	"github.com/albertlockett/learn-baseball-db/jobs"
)

var ctx = context.Background()

const indexName = "players"

const indexMapping = `
{
  "settings": {
		"number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "dynamic": false,
    "properties": {
      "name": {
        "type": "keyword"
      },
      "team": {
        "type": "keyword"
      },
      "position": {
        "type": "keyword"
      }
    }
  }
}
`

func main() {
	client := client.CreateClient()
	jobs.LoadPlayers(client)
}
