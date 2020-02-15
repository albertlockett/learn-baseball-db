package jobs

import (
	"log"

	"github.com/olivere/elastic/v7"
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
      "League": {
        "type": "keyword"
      },
      "Division": {
        "type": "keyword"
      },
      "City": {
        "type": "keyword"
			},
			"Name": {
        "type": "keyword"
			},
			"Code": {
        "type": "keyword"
      }
    }
  }
}
`

type esTeam struct {
	League   string `json:"League"`
	Division string `json:"Division"`
	City     string `json:"City"`
	Name     string `json:"Name"`
	Code     string `json:"Code"`
}

var teams = []*esTeam{
	// AL East
	&esTeam{League: "AL", Division: "east", City: "Baltimore", Name: "Orioles", Code: "BAL"},
	&esTeam{League: "AL", Division: "east", City: "Boston", Name: "Red Sox", Code: "BOS"},
	&esTeam{League: "AL", Division: "east", City: "New York", Name: "Yankees", Code: "NYY"},
	&esTeam{League: "AL", Division: "east", City: "Tampa Bay", Name: "Rays", Code: "TB"},
	&esTeam{League: "AL", Division: "east", City: "Toronto", Name: "Blue Jays", Code: "TOR"},

	// AL Central
	&esTeam{League: "AL", Division: "central", City: "Chicago", Name: "White Sox", Code: "CWS"},
	&esTeam{League: "AL", Division: "central", City: "Cleveland", Name: "Indians", Code: "CLE"},
	&esTeam{League: "AL", Division: "central", City: "Detroit", Name: "Tigers", Code: "DET"},
	&esTeam{League: "AL", Division: "central", City: "Kansas City", Name: "Royals", Code: "KC"},
	&esTeam{League: "AL", Division: "central", City: "Minnesota", Name: "Twins", Code: "MIN"},

	// AL West
	&esTeam{League: "AL", Division: "west", City: "Houston", Name: "Astros", Code: "HOU"},
	&esTeam{League: "AL", Division: "west", City: "Los Angeles", Name: "Angels", Code: "LAA"},
	&esTeam{League: "AL", Division: "west", City: "Oakland", Name: "Athletics", Code: "OAK"},
	&esTeam{League: "AL", Division: "west", City: "Seattle", Name: "Mariners", Code: "SEA"},
	&esTeam{League: "AL", Division: "west", City: "Texas", Name: "Rangers", Code: "TEX"},

	// NL East
	&esTeam{League: "NL", Division: "east", City: "Atlanta", Name: "Braves", Code: "ATL"},
	&esTeam{League: "NL", Division: "east", City: "Miami", Name: "Marlins", Code: "MIA"},
	&esTeam{League: "NL", Division: "east", City: "New York", Name: "Mets", Code: "NYM"},
	&esTeam{League: "NL", Division: "east", City: "Philadelphia", Name: "Phillies", Code: "PHI"},
	&esTeam{League: "NL", Division: "east", City: "Washington", Name: "Nationals", Code: "WSH"},

	// NL Central
	&esTeam{League: "NL", Division: "central", City: "Chicago", Name: "Cubs", Code: "CHC"},
	&esTeam{League: "NL", Division: "central", City: "Cincinati", Name: "Reds", Code: "CIN"},
	&esTeam{League: "NL", Division: "central", City: "Milwaukee", Name: "Brewers", Code: "MIL"},
	&esTeam{League: "NL", Division: "central", City: "Pittsburgh", Name: "Pirates", Code: "PIT"},
	&esTeam{League: "NL", Division: "central", City: "St. Lous", Name: "Cardinals", Code: "STL"},

	// NL West
	&esTeam{League: "NL", Division: "west", City: "Arizona", Name: "Diamondbacks", Code: "ARI"},
	&esTeam{League: "NL", Division: "west", City: "Colorado", Name: "Rockies", Code: "COL"},
	&esTeam{League: "NL", Division: "west", City: "Los Angeles", Name: "Dodgers", Code: "LAD"},
	&esTeam{League: "NL", Division: "west", City: "San Diego", Name: "Padres", Code: "SD"},
	&esTeam{League: "NL", Division: "west", City: "San Francisco", Name: "Giants", Code: "SF"},
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

	for i := 0; i < len(teams); i++ {
		row := teams[i]
		_, err := client.Index().
			Index(teamIndexName).
			Type("_doc").
			Id(row.Code).
			BodyJson(row).
			Do(ctx)
		if err != nil {
			log.Println("An error happened here")
			panic(err)
		}
	}

	return true, nil
}
