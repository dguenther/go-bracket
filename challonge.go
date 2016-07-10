package bracket

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type challongeAPIResponse struct {
	Tournament *challongeTournament `json:"tournament"`
}

type challongeTournament struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	State            string     `json:"state"`
	StartedAt        *time.Time `json:"started_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	FullChallongeURL string     `json:"full_challonge_url"`

	Participants []*challongeParticipantWrap `json:"participants"`
	Matches      []*challongeMatchWrap       `json:"matches"`
}

type challongeParticipantWrap struct {
	Participant *challongeParticipant `json:"participant"`
}

type challongeParticipant struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Seed        int    `json:"seed"`
	FinalRank   int    `json:"final_rank"`
	DisplayName string `json:"display_name"`
}

type challongeMatchWrap struct {
	Match *challongeMatch `json:"match"`
}

type challongeMatch struct {
	ID                   int        `json:"id"`
	Identifier           string     `json:"identifier"`
	Round                int        `json:"round"`
	StartedAt            *time.Time `json:"started_at,omitempty"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
	State                string     `json:"state"`
	Player1ID            int        `json:"player1_id"`
	Player2ID            int        `json:"player2_id"`
	Player1PrereqMatchID int        `json:"player1_prereq_match_id"`
	Player2PrereqMatchID int        `json:"player2_prereq_match_id"`
	WinnerID             int        `json:"winner_id"`
	LoserID              int        `json:"loser_id"`
	ScoresCsv            string     `json:"scores_csv"`
}

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

func fetchChallongeData(user, apiKey, apiURL string) *challongeAPIResponse {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var decoded challongeAPIResponse
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		log.Fatal(err)
	}
	return &decoded
}

func convertChallongePlayers(data []*challongeParticipantWrap) []*Player {
	players := make([]*Player, len(data))
	for i, d := range data {
		players[i] = &Player{
			Name: d.Participant.DisplayName,
			Seed: d.Participant.Seed,
			Rank: d.Participant.FinalRank,
		}
	}
	return players
}

func convertChallongeMatches(data []*challongeMatchWrap) []*Match {
	matches := make([]*Match, len(data))
	for i, d := range data {
		p1score := 0
		p2score := 0
		// sum up the set results, since we're not tracking sets yet
		for _, set := range strings.Split(d.Match.ScoresCsv, ",") {
			scoreSplit := strings.SplitN(set, "-", 2)
			p1setscore, _ := strconv.Atoi(scoreSplit[0])
			p2setscore, _ := strconv.Atoi(scoreSplit[1])
			p1score += p1setscore
			p2score += p2setscore
		}

		matches[i] = &Match{
			ID:           strconv.Itoa(d.Match.ID),
			Identifier:   d.Match.Identifier,
			UpdatedAt:    d.Match.UpdatedAt,
			Round:        d.Match.Round,
			State:        d.Match.State,
			Player1ID:    strconv.Itoa(d.Match.Player1ID),
			Player2ID:    strconv.Itoa(d.Match.Player2ID),
			WinnerID:     strconv.Itoa(d.Match.WinnerID),
			LoserID:      strconv.Itoa(d.Match.LoserID),
			Player1Score: p1score,
			Player2Score: p2score,
		}
	}
	return matches
}

func convertChallongeData(data *challongeAPIResponse) *Bracket {
	return &Bracket{
		URL:     data.Tournament.FullChallongeURL,
		Name:    data.Tournament.Name,
		Players: convertChallongePlayers(data.Tournament.Participants),
		Matches: convertChallongeMatches(data.Tournament.Matches),
	}
}

func fetchChallongeBracket(user, apiKey, url string) *Bracket {
	resp := fetchChallongeData(user, apiKey, url)
	return convertChallongeData(resp)
}
