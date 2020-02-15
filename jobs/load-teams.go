package jobs

import (
	"log"

	"github.com/olivere/elastic"
)

const teamIndexName = "teams"

const teamIndexMapping = `
{
  "settings": {
		"number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "dynamic": false,
    "properties": {
      "league": {
        "type": "keyword"
      },
      "division": {
        "type": "keyword"
      },
      "city": {
        "type": "keyword"
			},
			"name": {
        "type": "keyword"
			},
			"code": {
        "type": "keyword"
      }
    }
  }
}
`

type esTeam struct {
	league   string `json:"league"`
	division string `json:"division"`
	city     string `json:"city"`
	name     string `json:"name"`
	code     string `json:"code"`
}

var teams = []*esTeam{
	// AL East
	&esTeam{league: "AL", division: "east", city: "Baltimore", name: "Orioles", code: "BAL"},
	&esTeam{league: "AL", division: "east", city: "Boston", name: "Red Sox", code: "BOS"},
	&esTeam{league: "AL", division: "east", city: "New York", name: "Yankees", code: "NYY"},
	&esTeam{league: "AL", division: "east", city: "Tampa Bay", name: "Rays", code: "TB"},
	&esTeam{league: "AL", division: "east", city: "Toronto", name: "Blue Jays", code: "TOR"},

	// AL Central
	&esTeam{league: "AL", division: "central", city: "Chicago", name: "White Sox", code: "CSW"},
	&esTeam{league: "AL", division: "central", city: "Cleveland", name: "Indians", code: "CLE"},
	&esTeam{league: "AL", division: "central", city: "Detroit", name: "Tigers", code: "DET"},
	&esTeam{league: "AL", division: "central", city: "Kansas City", name: "Royals", code: "KC"},
	&esTeam{league: "AL", division: "central", city: "Minnesota", name: "Twins", code: "MIN"},

	// AL West
	&esTeam{league: "AL", division: "west", city: "Houston", name: "Astros", code: "HOU"},
	&esTeam{league: "AL", division: "west", city: "Los Angeles", name: "Angels", code: "LAA"},
	&esTeam{league: "AL", division: "west", city: "Oakland", name: "Athletics", code: "OAK"},
	&esTeam{league: "AL", division: "west", city: "Seattle", name: "Mariners", code: "STL"},
	&esTeam{league: "AL", division: "west", city: "Texas", name: "Rangers", code: "TEX"},

	// NL East
	&esTeam{league: "NL", division: "east", city: "Atlanta", name: "Braves", code: "ATL"},
	&esTeam{league: "NL", division: "east", city: "Miami", name: "Marlins", code: "MIA"},
	&esTeam{league: "NL", division: "east", city: "New York", name: "Mets", code: "NYM"},
	&esTeam{league: "NL", division: "east", city: "Philadelphia", name: "Phillies", code: "PHI"},
	&esTeam{league: "NL", division: "east", city: "Washington", name: "Nationals", code: "WAS"},

	// NL Central
	&esTeam{league: "NL", division: "central", city: "Chicago", name: "Cubs", code: "CHC"},
	&esTeam{league: "NL", division: "central", city: "Cincinati", name: "Reds", code: "CIN"},
	&esTeam{league: "NL", division: "central", city: "Milwaukee", name: "Brewers", code: "MIL"},
	&esTeam{league: "NL", division: "central", city: "Pittsburgh", name: "Pirates", code: "PIT"},
	&esTeam{league: "NL", division: "central", city: "St. Lous", name: "Cardinals", code: "STL"},

	// NL West
	&esTeam{league: "NL", division: "west", city: "Arizona", name: "Diamondbacks", code: "ARI"},
	&esTeam{league: "NL", division: "west", city: "Colorado", name: "Rockies", code: "COL"},
	&esTeam{league: "NL", division: "west", city: "Los Angeles", name: "Dodgers", code: "LAD"},
	&esTeam{league: "NL", division: "west", city: "San Diego", name: "Padres", code: "SD"},
	&esTeam{league: "NL", division: "west", city: "San Francisco", name: "Giants", code: "SF"},
}

// LoadTeams load all the teams
func LoadTeams(client *elastic.Client) (bool, error) {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(teamIndexName).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	//if index does not exist, create a new one with the specified mapping
	if !exists {
		createIndex, err := client.CreateIndex(teamIndexName).BodyString(teamIndexMapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			log.Println(createIndex)
		} else {
			log.Println("successfully created teams index")
		}
	} else {
		log.Println("teams index already exist")
	}

	for i := 1; i < len(teams); i++ {
		row := teams[i]
		_, err := client.Index().
			Index(teamIndexName).
			Type("_doc").
			Id(row.code).
			BodyJson(row).
			Do(ctx)
		if err != nil {
			log.Println("An error happened here")
			panic(err)
		}
	}

	return true, nil
}
