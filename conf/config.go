package conf

import (
	"log"
	"os"
)

var (
	CommitHash string
	BuildTime  string

	ApiToken         = os.Getenv("API_TOKEN")
	HttpTimeout      = 5
	TibberGraphQLUrl = "https://api.tibber.com/v1-beta/gql"
	//Query            = `{"query":"{viewer { homes { consumption(resolution: HOURLY, last: 1) { nodes { from to cost unitPrice unitPriceVAT consumption  } } } }}"}`
	Query = `{"query":"{\n  viewer {\n    homes {\n      currentSubscription{\n        priceInfo{\n          today {\n            total\n            energy\n            tax\n            startsAt\n          }\n          tomorrow {\n            total\n            energy\n            tax\n            startsAt\n          }\n        }\n      }\n    }\n  }\n}\n"}`
)

func EnvironmentComplete() {
	envComplete := true

	if len(ApiToken) == 0 {
		log.Print("missing envvar \"API_TOKEN\"")
		envComplete = false
	}

	if !envComplete {
		log.Fatal("one or more required envvars missing, aborting...")
	}
}
