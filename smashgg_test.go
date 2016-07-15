package bracket

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsSmashGGURL(t *testing.T) {
	url := "https://smash.gg/tournament/super-smash-sundays-48/brackets/14221/50133/165583"
	assert.True(t, isSmashGGURL(url))
}

func TestGetSmashGGAPIURL(t *testing.T) {
	url := "https://smash.gg/tournament/super-smash-sundays-48/brackets/14221/50133/165583"
	apiURL := "https://smash.gg/api/-/resource/gg_api./phase_group/165583;expand=%5B%22sets%22%2C%22seeds%22%2C%22standings%22%5D;mutations=%5B%22playerData%22%5D;reset=false"
	assert.Equal(t, apiURL, getSmashGGAPIURL(url))
}

func TestDecodeSmashGGData(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/smashgg.json")
	if err != nil {
		t.Error(err)
	}
	resp, err := decodeSmashGGData(b)
	if err != nil {
		t.Error(err)
	}
	e := resp.Entities
	// Groups
	g := e.Groups
	assert.Equal(t, 171722, g.ID)
	assert.Equal(t, 50132, g.PhaseID)
	assert.Equal(t, 8322, g.WaveID)
	assert.Equal(t, 3, g.State)

	// Sets
	s := e.Sets
	assert.Len(t, s, 2)

	set := s[0]
	assert.Equal(t, 4689059, set.ID)
	assert.Equal(t, "A", set.Identifier)
	assert.Equal(t, 1, set.Round)
	assert.Equal(t, 3, set.State)
	assert.Equal(t, 211768, set.Entrant1ID)
	assert.Equal(t, 2426316, set.Entrant1PrereqID)
	assert.Equal(t, "seed", set.Entrant1PrereqType)
	assert.Equal(t, 2, set.Entrant1Score)
	assert.Equal(t, 212928, set.Entrant2ID)
	assert.Equal(t, 2428388, set.Entrant2PrereqID)
	assert.Equal(t, "seed", set.Entrant2PrereqType)
	assert.Equal(t, 0, set.Entrant2Score)
	assert.Equal(t, 211768, set.WinnerID)
	assert.Equal(t, 212928, set.LoserID)
	assert.EqualValues(t, 1468185969, set.UpdatedAt)

	set = s[1]
	assert.Equal(t, 4689067, set.ID)
	assert.Equal(t, "I", set.Identifier)
	assert.Equal(t, 2, set.Round)
	assert.Equal(t, 3, set.State)
	assert.Equal(t, 211768, set.Entrant1ID)
	assert.Equal(t, 4689059, set.Entrant1PrereqID)
	assert.Equal(t, "set", set.Entrant1PrereqType)
	assert.Equal(t, 2, set.Entrant1Score)
	assert.Equal(t, 211974, set.Entrant2ID)
	assert.Equal(t, 4689060, set.Entrant2PrereqID)
	assert.Equal(t, "set", set.Entrant2PrereqType)
	assert.Equal(t, 0, set.Entrant2Score)
	assert.Equal(t, 211768, set.WinnerID)
	assert.Equal(t, 211974, set.LoserID)
	assert.EqualValues(t, 1468187020, set.UpdatedAt)

	// Seeds
	p := e.Seeds
	assert.Len(t, p, 3)

	seed := e.Seeds[0]
	assert.Equal(t, 2426316, seed.ID)
	assert.Equal(t, 211768, seed.EntrantID)
	assert.Equal(t, 3, seed.Placement)
	assert.Equal(t, 7, seed.SeedNum)
	entrants := seed.Mutations.Entrants
	assert.Len(t, entrants, 1)
	entrant := entrants["211768"]
	assert.Equal(t, 211768, entrant.ID)
	assert.Equal(t, "TA | CDK", entrant.Name)
	assert.Equal(t, []int{238181}, entrant.ParticipantIds)
	assert.Equal(t, map[string]int{"238181": 1092}, entrant.PlayerIds)
	players := seed.Mutations.Players
	assert.Len(t, players, 1)
	player := players["1092"]
	assert.Equal(t, 1092, player.ID)
	assert.Equal(t, "CDK", player.GamerTag)
	assert.Equal(t, "TA", player.Prefix)
	assert.Equal(t, "CA", player.State)
	assert.Equal(t, "United States", player.Country)
	assert.Equal(t, "Connor Nguyen", player.Name)
	participants := seed.Mutations.Participants
	assert.Len(t, participants, 1)
	participant := participants["238181"]
	assert.Equal(t, 238181, participant.ID)
	assert.Equal(t, "CDK", participant.GamerTag)
	assert.Equal(t, "TA", participant.Prefix)
	assert.Equal(t, "Connor", participant.ContactInfo.NameFirst)
	assert.Equal(t, "Nguyen", participant.ContactInfo.NameLast)

	seed = e.Seeds[1]
	assert.Equal(t, 2426510, seed.ID)
	assert.Equal(t, 211974, seed.EntrantID)
	assert.Equal(t, 5, seed.Placement)
	assert.Equal(t, 70, seed.SeedNum)
	entrants = seed.Mutations.Entrants
	assert.Len(t, entrants, 1)
	entrant = entrants["211974"]
	assert.Equal(t, 211974, entrant.ID)
	assert.Equal(t, "A-Dar", entrant.Name)
	assert.Equal(t, []int{238347}, entrant.ParticipantIds)
	assert.Equal(t, map[string]int{"238347": 14453}, entrant.PlayerIds)
	players = seed.Mutations.Players
	assert.Len(t, players, 1)
	player = players["14453"]
	assert.Equal(t, 14453, player.ID)
	assert.Equal(t, "A-Dar", player.GamerTag)
	assert.Equal(t, "", player.Prefix)
	assert.Equal(t, "CA", player.State)
	assert.Equal(t, "US", player.Country)
	assert.Equal(t, "Brandon Panapanaan", player.Name)
	participants = seed.Mutations.Participants
	assert.Len(t, participants, 1)
	participant = participants["238347"]
	assert.Equal(t, 238347, participant.ID)
	assert.Equal(t, "A-Dar", participant.GamerTag)
	assert.Equal(t, "", participant.Prefix)
	assert.Equal(t, "Brandon", participant.ContactInfo.NameFirst)
	assert.Equal(t, "Panapanaan", participant.ContactInfo.NameLast)

	seed = e.Seeds[2]
	assert.Equal(t, 2428388, seed.ID)
	assert.Equal(t, 212928, seed.EntrantID)
	assert.Equal(t, 9, seed.Placement)
	assert.Equal(t, 122, seed.SeedNum)
	entrants = seed.Mutations.Entrants
	assert.Len(t, entrants, 1)
	entrant = entrants["212928"]
	assert.Equal(t, 212928, entrant.ID)
	assert.Equal(t, "Slime", entrant.Name)
	assert.Equal(t, []int{239067}, entrant.ParticipantIds)
	assert.Equal(t, map[string]int{"239067": 13963}, entrant.PlayerIds)
	players = seed.Mutations.Players
	assert.Len(t, players, 1)
	player = players["13963"]
	assert.Equal(t, 13963, player.ID)
	assert.Equal(t, "Slime", player.GamerTag)
	assert.Equal(t, "", player.Prefix)
	assert.Equal(t, "CA", player.State)
	assert.Equal(t, "United States", player.Country)
	assert.Equal(t, "Anthony Bruno", player.Name)
	participants = seed.Mutations.Participants
	assert.Len(t, participants, 1)
	participant = participants["239067"]
	assert.Equal(t, 239067, participant.ID)
	assert.Equal(t, "Slime", participant.GamerTag)
	assert.Equal(t, "", participant.Prefix)
	assert.Equal(t, "Anthony", participant.ContactInfo.NameFirst)
	assert.Equal(t, "Bruno", participant.ContactInfo.NameLast)
}

func TestConvertSmashGGData(t *testing.T) {
	// Ideally would be creating an API response here rather than
	// building it from a file
	b, err := ioutil.ReadFile("testdata/smashgg.json")
	if err != nil {
		t.Error(err)
	}
	resp, err := decodeSmashGGData(b)
	if err != nil {
		t.Error(err)
	}
	bracket := convertSmashGGData(resp)

	assert.Equal(t, "", bracket.Name)
	assert.Equal(t, "", bracket.URL)
	updatedAt := time.Unix(1468187020, 0)
	assert.Equal(t, &updatedAt, bracket.UpdatedAt)

	// Players
	players := bracket.Players
	assert.Len(t, players, 3)
	player := players[0]
	assert.Equal(t, "211768", player.ID)
	assert.Equal(t, "TA | CDK", player.Name)
	assert.Equal(t, 3, player.Rank)
	assert.Equal(t, 7, player.Seed)
	player = players[1]
	assert.Equal(t, "211974", player.ID)
	assert.Equal(t, "A-Dar", player.Name)
	assert.Equal(t, 5, player.Rank)
	assert.Equal(t, 70, player.Seed)
	player = players[2]
	assert.Equal(t, "212928", player.ID)
	assert.Equal(t, "Slime", player.Name)
	assert.Equal(t, 9, player.Rank)
	assert.Equal(t, 122, player.Seed)

	// Matches
	matches := bracket.Matches
	assert.Len(t, matches, 2)
	match := bracket.Matches[0]
	updatedAt = time.Unix(1468185969, 0)
	assert.Equal(t, "4689059", match.ID)
	assert.Equal(t, "A", match.Identifier)
	assert.Equal(t, 1, match.Round)
	assert.Equal(t, &updatedAt, match.UpdatedAt)
	assert.Equal(t, "complete", match.State)
	assert.Equal(t, "211768", match.Player1ID)
	assert.Equal(t, 2, match.Player1Score)
	assert.Equal(t, "212928", match.Player2ID)
	assert.Equal(t, 0, match.Player2Score)
	assert.Equal(t, "211768", match.WinnerID)
	assert.Equal(t, "212928", match.LoserID)
	match = bracket.Matches[1]
	updatedAt = time.Unix(1468187020, 0)
	assert.Equal(t, "4689067", match.ID)
	assert.Equal(t, "I", match.Identifier)
	assert.Equal(t, 2, match.Round)
	assert.Equal(t, &updatedAt, match.UpdatedAt)
	assert.Equal(t, "complete", match.State)
	assert.Equal(t, "211768", match.Player1ID)
	assert.Equal(t, 2, match.Player1Score)
	assert.Equal(t, "211974", match.Player2ID)
	assert.Equal(t, 0, match.Player2Score)
	assert.Equal(t, "211768", match.WinnerID)
	assert.Equal(t, "211974", match.LoserID)
}
