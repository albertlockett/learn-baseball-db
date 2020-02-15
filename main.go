package main

import (
	"context"
	"flag"
	"fmt"

	client "github.com/albertlockett/learn-baseball-db/es-client"
	"github.com/albertlockett/learn-baseball-db/jobs"
)

var ctx = context.Background()

func main() {
	dbjob := flag.String("job", "all", "which DB thing to do")
	eshost := flag.String("eshost", "http://127.0.0.1:9200", "host of elasticserch")
	flag.Parse()

	client := client.CreateClient(eshost)
	if "all" == *dbjob {
		fmt.Println("doall")
		jobs.LoadTeams(client)
		jobs.LoadPlayers(client)
	}

	if "players" == *dbjob {
		fmt.Println("players")
		jobs.LoadPlayers(client)
	}

	if "teams" == *dbjob {
		fmt.Println("teams")
		jobs.LoadTeams(client)
	}
}
