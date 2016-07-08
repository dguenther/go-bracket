package bracket

import (
	"log"
	"net/http"
	"strings"
)

func isChallongeURL(url string) bool {
	return strings.Contains(url, "challonge")
}

func getChallongeHash(url string) string {
	tourneyHash := url[strings.LastIndex(url, "/")+1 : len(url)]
	tourneyHash = strings.TrimSpace(tourneyHash)

	//If tournament belongs to an organization,
	//it must be specified in the request
	if len(strings.Split(url, "."))-1 > 1 {
		orgHash := url[strings.LastIndex(url, "://")+3 : strings.Index(url, ".")]
		return orgHash + "-" + tourneyHash
	}

	//Standard tournament
	return tourneyHash
}

func getChallongeAPIURL(url string) string {
	hash := getChallongeHash(url)
	return "https://api.challonge.com/v1/tournaments/" + hash + ".json?include_matches=1&include_participants=1"
}

func fetchChallongeData(user, apiKey, apiURL string) {
	req, err := http.NewRequest("asdf", apiURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, apiKey)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchChallongeBracket(user, apiKey, url string) *Bracket {
	return &Bracket{
		url: url,
	}
}
