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

type smashGGAPIResponse struct {
	Entities *smashGGEntities `json:"entities"`
}

type smashGGEntities struct {
	Groups *smashGGGroup  `json:"groups"`
	Sets   []*smashGGSet  `json:"sets"`
	Seeds  []*smashGGSeed `json:"seeds"`
}

type smashGGGroup struct {
	ID      int `json:"id"`
	PhaseID int `json:"phaseId"`
	WaveID  int `json:"waveId"`
	State   int `json:"state"`
}

type smashGGSet struct {
	ID                 int    `json:"id"`
	Identifier         string `json:"identifier"`
	Round              int    `json:"round"`
	UpdatedAt          int64  `json:"updatedAt"`
	State              int    `json:"state"`
	Entrant1ID         int    `json:"entrant1Id"`
	Entrant1Score      int    `json:"entrant1Score"`
	Entrant2Score      int    `json:"entrant2Score"`
	Entrant2ID         int    `json:"entrant2Id"`
	WinnerID           int    `json:"winnerId"`
	LoserID            int    `json:"loserId"`
	Entrant1PrereqType string `json:"entrant1PrereqType"`
	Entrant1PrereqID   int    `json:"entrant1PrereqId"`
	Entrant2PrereqType string `json:"entrant2PrereqType"`
	Entrant2PrereqID   int    `json:"entrant2PrereqId"`
}

type smashGGSeed struct {
	ID        int               `json:"id"`
	EntrantID int               `json:"entrantId"`
	SeedNum   int               `json:"seedNum"`
	Placement int               `json:"placement"`
	Mutations *smashGGMutations `json:"mutations"`
}

type smashGGMutations struct {
	Participants map[string]*smashGGParticipant `json:"participants"`
	Players      map[string]*smashGGPlayer      `json:"players"`
	Entrants     map[string]*smashGGEntrant     `json:"entrants"`
}

type smashGGParticipant struct {
	ID          int                 `json:"id"`
	GamerTag    string              `json:"gamerTag"`
	Prefix      string              `json:"prefix"`
	ContactInfo *smashGGContactInfo `json:"contactInfo"`
}

type smashGGContactInfo struct {
	NameFirst string `json:"nameFirst"`
	NameLast  string `json:"nameLast"`
}

type smashGGPlayer struct {
	ID       int    `json:"id"`
	GamerTag string `json:"gamerTag"`
	Prefix   string `json:"prefix"`
	Name     string `json:"name"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

type smashGGEntrant struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ParticipantIds []int          `json:"participantIds"`
	PlayerIds      map[string]int `json:"playerIds"`
}

func isSmashGGURL(url string) bool {
	return strings.Contains(url, "smash.gg")
}

func getSmashGGAPIURL(url string) string {
	trimURL := strings.TrimRight(url, "/")
	splitURL := strings.Split(trimURL, "/")
	phaseGroup := splitURL[len(splitURL)-1]
	return "https://smash.gg/api/-/resource/gg_api./phase_group/" + phaseGroup + ";expand=%5B%22sets%22%2C%22seeds%22%2C%22standings%22%5D;mutations=%5B%22playerData%22%5D;reset=false"
}

func fetchSmashGGData(apiURL string) *smashGGAPIResponse {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return decodeSmashGGData(body)
}

func decodeSmashGGData(body []byte) *smashGGAPIResponse {
	var decoded smashGGAPIResponse
	err := json.Unmarshal(body, &decoded)
	if err != nil {
		log.Fatal(err)
	}
	return &decoded
}

func convertSmashGGState(set *smashGGSet) string {
	if set.State == 3 {
		return "complete"
	}
	if set.State == 1 &&
		set.Entrant1ID != 0 &&
		set.Entrant2ID != 0 {
		return "open"
	}
	return "pending"
}

func convertSmashGGMatches(resp *smashGGAPIResponse) []*Match {
	var filteredSets []*smashGGSet

	for _, s := range resp.Entities.Sets {
		// smash gg seems to return a lot of junk matches, so let's
		// filter them out.
		// In particular, cases where there are byes in round 1
		if s.Round == 1 || s.Round == -1 {
			if !(s.Entrant1PrereqType == "bye" || s.Entrant2PrereqType == "bye") {
				filteredSets = append(filteredSets, s)
				continue
			}
		}
		if !(s.Entrant1PrereqType == "bye" && s.Entrant2PrereqType == "bye") {
			filteredSets = append(filteredSets, s)
		}
	}
	matches := make([]*Match, len(filteredSets))
	for i, s := range filteredSets {
		updatedAt := time.Unix(s.UpdatedAt, 0)
		matches[i] = &Match{
			ID:           strconv.Itoa(s.ID),
			Identifier:   s.Identifier,
			UpdatedAt:    &updatedAt,
			Round:        s.Round,
			State:        convertSmashGGState(s),
			Player1ID:    strconv.Itoa(s.Entrant1ID),
			Player1Score: s.Entrant1Score,
			Player2ID:    strconv.Itoa(s.Entrant2ID),
			Player2Score: s.Entrant2Score,
			WinnerID:     strconv.Itoa(s.WinnerID),
			LoserID:      strconv.Itoa(s.LoserID),
		}
	}
	return matches
}

func convertSmashGGPlayers(resp *smashGGAPIResponse) []*Player {
	players := make([]*Player, len(resp.Entities.Seeds))
	for i, p := range resp.Entities.Seeds {
		entrantID := strconv.Itoa(p.EntrantID)
		players[i] = &Player{
			ID:   entrantID,
			Name: p.Mutations.Entrants[entrantID].Name,
			Seed: p.SeedNum,
			Rank: p.Placement,
		}
	}
	return players
}

func convertSmashGGData(resp *smashGGAPIResponse) *Bracket {
	return &Bracket{
		Name:    "", // API does not return a tournament name
		Matches: convertSmashGGMatches(resp),
		Players: convertSmashGGPlayers(resp),
	}
}

func fetchSmashGGBracket(url string) *Bracket {
	apiURL := getSmashGGAPIURL(url)
	resp := fetchSmashGGData(apiURL)

	b := convertSmashGGData(resp)
	b.URL = url
	return b
}
