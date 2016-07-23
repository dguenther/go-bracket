package bracket

import (
	"encoding/json"
	"io/ioutil"
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
	StartedAt          *int64 `json:"startedAt"`
	State              int    `json:"state"`
	Entrant1ID         int    `json:"entrant1Id"`
	Entrant1Score      int    `json:"entrant1Score"`
	Entrant2Score      int    `json:"entrant2Score"`
	Entrant2ID         int    `json:"entrant2Id"`
	WinnerID           int    `json:"winnerId"`
	LoserID            int    `json:"loserId"`
	Entrant1PrereqType string `json:"entrant1PrereqType"`
	Entrant1PrereqID   *int   `json:"entrant1PrereqId"`
	Entrant2PrereqType string `json:"entrant2PrereqType"`
	Entrant2PrereqID   *int   `json:"entrant2PrereqId"`
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

func fetchSmashGGData(apiURL string) (*smashGGAPIResponse, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return decodeSmashGGData(body)
}

func decodeSmashGGData(body []byte) (*smashGGAPIResponse, error) {
	var decoded smashGGAPIResponse
	err := json.Unmarshal(body, &decoded)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}

func convertSmashGGState(stateNum int, hasEntrants bool) string {
	if stateNum == 3 {
		return "complete"
	}
	if stateNum == 1 && hasEntrants {
		return "open"
	}
	return "pending"
}

func convertSmashGGMatches(resp *smashGGAPIResponse) []*Match {
	// smash gg seems to return a lot of junk matches, so let's
	// filter them out.
	// In particular, cases where there are byes in round 1
	var filteredSets []*smashGGSet
	for _, s := range resp.Entities.Sets {
		if s.Round == 1 || s.Round == -1 {
			if s.Entrant1PrereqType != "bye" && s.Entrant2PrereqType != "bye" {
				filteredSets = append(filteredSets, s)
			}
		} else {
			filteredSets = append(filteredSets, s)
		}
	}

	matches := make([]*Match, len(filteredSets))
	for i, s := range filteredSets {
		updatedAt := time.Unix(s.UpdatedAt, 0)
		var startedAt *time.Time
		if s.StartedAt != nil {
			time := time.Unix(*s.StartedAt, 0)
			startedAt = &time
		}

		var p1prereq *string
		var p2prereq *string
		if s.Entrant1PrereqID != nil {
			p1prereq = new(string)
			*p1prereq = strconv.Itoa(*s.Entrant1PrereqID)
		}
		if s.Entrant2PrereqID != nil {
			p2prereq = new(string)
			*p2prereq = strconv.Itoa(*s.Entrant2PrereqID)
		}

		matches[i] = &Match{
			ID:                   strconv.Itoa(s.ID),
			Identifier:           s.Identifier,
			StartedAt:            startedAt,
			UpdatedAt:            &updatedAt,
			Round:                s.Round,
			State:                convertSmashGGState(s.State, s.Entrant1ID != 0 && s.Entrant2ID != 0),
			Player1ID:            strconv.Itoa(s.Entrant1ID),
			Player1Score:         s.Entrant1Score,
			Player1PrereqMatchID: p1prereq,
			Player2ID:            strconv.Itoa(s.Entrant2ID),
			Player2Score:         s.Entrant2Score,
			Player2PrereqMatchID: p2prereq,
			WinnerID:             strconv.Itoa(s.WinnerID),
			LoserID:              strconv.Itoa(s.LoserID),
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
	// Build tournament state
	state := ""
	if resp.Entities != nil && resp.Entities.Groups != nil {
		state = convertSmashGGState(resp.Entities.Groups.State, true)
	}

	b := &Bracket{
		Name:    "", // API does not return a tournament name
		URL:     "", // API does not return a tournament URL
		State:   state,
		Matches: convertSmashGGMatches(resp),
		Players: convertSmashGGPlayers(resp),
	}

	// API does not return updatedAt or startedAt on tournament, so
	// attempt to pull that off of the matches
	for _, m := range b.Matches {
		if m.UpdatedAt != nil && (b.UpdatedAt == nil || m.UpdatedAt.After(*b.UpdatedAt)) {
			b.UpdatedAt = m.UpdatedAt
		}
		if m.StartedAt != nil && (b.StartedAt == nil || m.StartedAt.Before(*b.StartedAt)) {
			b.StartedAt = m.StartedAt
		}
	}

	return b
}

func fetchSmashGGBracket(url string) (*Bracket, error) {
	apiURL := getSmashGGAPIURL(url)
	resp, err := fetchSmashGGData(apiURL)
	if err != nil {
		return nil, err
	}

	b := convertSmashGGData(resp)
	b.URL = url
	return b, nil
}
