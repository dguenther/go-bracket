package bracket

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsChallongeUrl(t *testing.T) {
	url := "http://challonge.com/xyfuz5c3"
	assert.True(t, isChallongeURL(url))
}

func TestGetChallongeHash(t *testing.T) {
	url := "http://challonge.com/xyfuz5c3"
	hash := "xyfuz5c3"
	assert.Equal(t, getChallongeHash(url), hash)
}

func TestGetChallongeOrgHash(t *testing.T) {
	orgURL := "http://smashchateau.challonge.com/melee_halloween"
	orgHash := "smashchateau-melee_halloween"
	assert.Equal(t, getChallongeHash(orgURL), orgHash)
}

func TestDecodeChallonge(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/challonge.json")
	if err != nil {
		t.Error(err)
	}
	resp, err := decodeChallongeData(b)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2385234, resp.Tournament.ID)
	assert.Equal(t, "Missouri River Arcadian - The Sequel: Smash4 Top 16", resp.Tournament.Name)
	assert.Equal(t, "http://HSCSmashNE.challonge.com/MRA2_s4s_t16", resp.Tournament.FullChallongeURL)
	assert.Equal(t, "complete", resp.Tournament.State)
	startedAt, _ := time.Parse(time.RFC3339, "2016-04-02T21:02:39.766-06:00")
	assert.Equal(t, &startedAt, resp.Tournament.StartedAt)
	completedAt, _ := time.Parse(time.RFC3339, "2016-04-03T00:00:43.525-06:00")
	assert.Equal(t, &completedAt, resp.Tournament.CompletedAt)
	createdAt, _ := time.Parse(time.RFC3339, "2016-04-02T19:54:30.732-06:00")
	assert.Equal(t, &createdAt, resp.Tournament.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, "2016-04-03T00:00:43.621-06:00")
	assert.Equal(t, &updatedAt, resp.Tournament.UpdatedAt)

	// Participants
	assert.Len(t, resp.Tournament.Participants, 2)
	participant := resp.Tournament.Participants[0].Participant
	assert.Equal(t, 38172466, participant.ID)
	assert.Equal(t, "(P1W) DPS|Dr. Pizza", participant.Name)
	assert.Equal(t, 1, participant.Seed)
	assert.Equal(t, "(P1W) DPS|Dr. Pizza", participant.DisplayName)
	assert.Equal(t, 9, participant.FinalRank)

	participant = resp.Tournament.Participants[1].Participant
	assert.Equal(t, 38172533, participant.ID)
	assert.Equal(t, "(P1L) YCL|Hite", participant.Name)
	assert.Equal(t, 16, participant.Seed)
	assert.Equal(t, "(P1L) YCL|Hite", participant.DisplayName)
	assert.Equal(t, 13, participant.FinalRank)

	// Matches
	assert.Len(t, resp.Tournament.Matches, 1)
	match := resp.Tournament.Matches[0].Match
	assert.Equal(t, 58296521, match.ID)
	assert.Equal(t, "complete", match.State)
	assert.Equal(t, "A", match.Identifier)
	assert.Equal(t, 1, match.Round)
	startedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:39.812-06:00")
	assert.Equal(t, &startedAt, match.StartedAt)
	completedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:49.417-06:00")
	assert.Equal(t, &completedAt, match.CompletedAt)
	createdAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:39.653-06:00")
	assert.Equal(t, &createdAt, match.CreatedAt)
	updatedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:49.397-06:00")
	assert.Equal(t, &updatedAt, match.UpdatedAt)
	assert.Equal(t, 38172466, match.Player1ID)
	assert.Equal(t, 38172533, match.Player2ID)
	assert.Nil(t, match.Player1PrereqMatchID)
	assert.Nil(t, match.Player2PrereqMatchID)
	assert.Equal(t, 38172533, match.LoserID)
	assert.Equal(t, 38172466, match.WinnerID)
	assert.Equal(t, "0--1", match.ScoresCsv)
}

func TestConvertChallongeData(t *testing.T) {
	// Ideally would be creating an API response here rather than
	// building it from a file
	b, err := ioutil.ReadFile("testdata/challonge.json")
	if err != nil {
		t.Error(err)
	}
	resp, err := decodeChallongeData(b)
	if err != nil {
		t.Error(err)
	}
	bracket := convertChallongeData(resp)

	assert.Equal(t, "Missouri River Arcadian - The Sequel: Smash4 Top 16", bracket.Name)
	assert.Equal(t, "http://HSCSmashNE.challonge.com/MRA2_s4s_t16", bracket.URL)
	assert.Equal(t, "complete", bracket.State)
	updatedAt, _ := time.Parse(time.RFC3339, "2016-04-03T00:00:43.621-06:00")
	assert.Equal(t, &updatedAt, bracket.UpdatedAt)
	startedAt, _ := time.Parse(time.RFC3339, "2016-04-02T21:02:39.766-06:00")
	assert.Equal(t, &startedAt, bracket.StartedAt)
	players := bracket.Players
	assert.Len(t, players, 2)
	player := players[0]
	assert.Equal(t, "38172466", player.ID)
	assert.Equal(t, "(P1W) DPS|Dr. Pizza", player.Name)
	assert.Equal(t, 9, player.Rank)
	assert.Equal(t, 1, player.Seed)
	player = players[1]
	assert.Equal(t, "38172533", player.ID)
	assert.Equal(t, "(P1L) YCL|Hite", player.Name)
	assert.Equal(t, 13, player.Rank)
	assert.Equal(t, 16, player.Seed)
	matches := bracket.Matches
	assert.Len(t, matches, 1)
	match := bracket.Matches[0]
	updatedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:49.397-06:00")
	startedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:39.812-06:00")
	assert.Equal(t, "58296521", match.ID)
	assert.Equal(t, "A", match.Identifier)
	assert.Equal(t, 1, match.Round)
	assert.Equal(t, &updatedAt, match.UpdatedAt)
	assert.Equal(t, &startedAt, match.StartedAt)
	assert.Equal(t, "complete", match.State)
	assert.Equal(t, "38172466", match.Player1ID)
	assert.Equal(t, 0, match.Player1Score)
	assert.Equal(t, "38172533", match.Player2ID)
	assert.Equal(t, -1, match.Player2Score)
	assert.Equal(t, "38172466", match.WinnerID)
	assert.Equal(t, "38172533", match.LoserID)
}
