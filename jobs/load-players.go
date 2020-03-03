package jobs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

const playersIndexName = "players"

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
			},
			"fantasyRank": {
				"type": "integer"
			},
			"bats": {
				"type": "keyword"
			},
			"throws": {
				"type": "keyword"
			},
			"debut": {
				"type": "date"
			},
			"born": {
				"type": "date"
			},
			"birthCity": {
				"type": "keyword"
			},
			"birthState": {
				"type": "keyword"
			},
			"birthCountry": {
				"type": "keyword"
			},
			"playerId": {
				"type": "keyword"
			}
    }
  }
}
`

var ctx = context.Background()

type row struct {
	Position     string `json:"position"`
	Name         string `json:"name_display_first_last"`
	Team         string `json:"team_abbrev"`
	Bats         string `json:"bats"`
	Throws       string `json:"throws"`
	Debut        string `json:"pro_debut_date"`
	Born         string `json:"birth_date"`
	BirthCity    string `json:"birth_city"`
	BirthState   string `json:"birth_state"`
	BirthCountry string `json:"birth_country"`
	PlayerID     string `json:"player_id"`
}

type queryResults struct {
	Created   string `json:"created"`
	TotalSize string `json:"totalSize"`
	Rows      []row  `json:"row"`
}

type searchPlayerAll struct {
	CopyRight string       `json:"copyRight"`
	QR        queryResults `json:"queryResults"`
}

type result struct {
	SPA searchPlayerAll `json:"search_player_all"`
}

type esPlayer struct {
	Name         string `json:"name"`
	Position     string `json:"position"`
	Team         string `json:"team"`
	FantasyRank  string `json:"fantasyRank"`
	Bats         string `json:"bats"`
	Throws       string `json:"throws"`
	Debut        string `json:"debut"`
	Born         string `json:"born"`
	BirthCity    string `json:"birthCity"`
	BirthState   string `json:"birthState"`
	BirthCountry string `json:"birthCountry"`
	PlayerID     string `json:"playerId"`
}

func rowToPlayer(r row) esPlayer {
	ob2 := esPlayer{}
	ob2.Name = r.Name
	ob2.Position = r.Position
	ob2.Team = r.Team
	ob2.Bats = r.Bats
	ob2.Throws = r.Throws
	ob2.Debut = r.Debut
	ob2.Born = r.Born
	ob2.BirthCity = r.BirthCity
	ob2.BirthState = r.BirthState
	ob2.BirthCountry = r.BirthCountry
	ob2.PlayerID = r.PlayerID
	if ob2.Debut == "" {
		ob2.Debut = "1970-01-01T00:00:00"
	}
	return ob2
}

// LoadPlayers list of Players into the database
func LoadPlayers(client *elastic.Client) (bool, error) {

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(playersIndexName).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	//if index does not exist, create a new one with the specified mapping
	if !exists {
		createIndex, err := client.CreateIndex(playersIndexName).BodyString(indexMapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			log.Println(createIndex)
		} else {
			log.Println("successfully created players index")
		}
	} else {
		log.Println("players index already exist")
	}

	file, err := os.Open("./data/fantasy_rankings.csv")
	if err != nil {
		panic(err)
	}

	fantasyRanks := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		segments := strings.Split(line, ",")
		rank := segments[0]
		playerName := segments[1]
		fantasyRanks[playerName] = rank
	}

	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for _, char1 := range alphabet {
		for _, char2 := range alphabet {

			url := "http://lookup-service-prod.mlb.com/json/named.search_player_all.bam?sport_code=%27mlb%27&active_sw=%27Y%27&name_part=%27" + string(char1) + string(char2) + "%25%27"
			httpClient := http.Client{
				Timeout: time.Second * 10, // Maximum of 2 secs
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {

				log.Fatal(err)
			}

			res, getErr := httpClient.Do(req)
			if getErr != nil {

				log.Fatal(getErr)
			}

			body, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {

				log.Fatal(readErr)
			}

			var info result
			// err := json.Unmarshal([]byte(text), &info);
			jsonErr := json.Unmarshal(body, &info)
			if jsonErr != nil {
				fmt.Println(string(body))
				log.Println("An error happened here")
				log.Println(jsonErr)
			} else {

				rows := info.SPA.QR.Rows
				for i := 1; i < len(rows); i++ {
					row := rows[i]
					esPlayer := rowToPlayer(row)
					frank, frpresent := fantasyRanks[esPlayer.Name]
					if frpresent {
						esPlayer.FantasyRank = frank
					}

					_, err := client.Index().
						Index(playersIndexName).
						Type("_doc").
						Id(row.Name).
						BodyJson(esPlayer).
						Do(ctx)
					if err != nil {
						log.Println("An error happened here")
						panic(err)
					}
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	}

	return true, nil
}
