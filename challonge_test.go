package bracket

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const url = "http://challonge.com/xyfuz5c3"
const hash = "xyfuz5c3"

const orgURL = "http://smashchateau.challonge.com/melee_halloween"
const orgHash = "smashchateau-melee_halloween"

func TestIsChallongeUrl(t *testing.T) {
	assert.True(t, isChallongeURL(url))
}

func TestGetChallongeHash(t *testing.T) {
	assert.Equal(t, getChallongeHash(url), hash)
}

func TestGetChallongeOrgHash(t *testing.T) {
	assert.Equal(t, getChallongeHash(orgURL), orgHash)
}

func TestFetchChallonge(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/challonge.json")
	if err != nil {
		t.Error(err)
	}
	resp := decodeChallongeData(b)
	assert.Equal(t, resp.Tournament.ID, 2385234)
	assert.Equal(t, resp.Tournament.Name, "Missouri River Arcadian - The Sequel: Smash4 Top 16")
	assert.Equal(t, resp.Tournament.FullChallongeURL, "http://HSCSmashNE.challonge.com/MRA2_s4s_t16")
	assert.Equal(t, resp.Tournament.State, "complete")
	startedAt, _ := time.Parse(time.RFC3339, "2016-04-02T21:02:39.766-06:00")
	assert.Equal(t, resp.Tournament.StartedAt, &startedAt)
	completedAt, _ := time.Parse(time.RFC3339, "2016-04-03T00:00:43.525-06:00")
	assert.Equal(t, resp.Tournament.CompletedAt, &completedAt)
	createdAt, _ := time.Parse(time.RFC3339, "2016-04-02T19:54:30.732-06:00")
	assert.Equal(t, resp.Tournament.CreatedAt, &createdAt)
	updatedAt, _ := time.Parse(time.RFC3339, "2016-04-03T00:00:43.621-06:00")
	assert.Equal(t, resp.Tournament.UpdatedAt, &updatedAt)

	// Participants
	assert.Len(t, resp.Tournament.Participants, 2)
	participant := resp.Tournament.Participants[0].Participant
	assert.Equal(t, participant.ID, 38172466)
	assert.Equal(t, participant.Name, "(P1W) DPS|Dr. Pizza")
	assert.Equal(t, participant.Seed, 1)
	assert.Equal(t, participant.DisplayName, "(P1W) DPS|Dr. Pizza")
	assert.Equal(t, participant.FinalRank, 9)

	participant = resp.Tournament.Participants[1].Participant
	assert.Equal(t, participant.ID, 38172533)
	assert.Equal(t, participant.Name, "(P1L) YCL|Hite")
	assert.Equal(t, participant.Seed, 16)
	assert.Equal(t, participant.DisplayName, "(P1L) YCL|Hite")
	assert.Equal(t, participant.FinalRank, 13)

	// Matches
	assert.Len(t, resp.Tournament.Matches, 1)
	match := resp.Tournament.Matches[0].Match
	assert.Equal(t, match.ID, 58296521)
	assert.Equal(t, match.State, "complete")
	assert.Equal(t, match.Identifier, "A")
	assert.Equal(t, match.Round, 1)
	startedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:39.812-06:00")
	assert.Equal(t, match.StartedAt, &startedAt)
	completedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:49.417-06:00")
	assert.Equal(t, match.CompletedAt, &completedAt)
	createdAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:39.653-06:00")
	assert.Equal(t, match.CreatedAt, &createdAt)
	updatedAt, _ = time.Parse(time.RFC3339, "2016-04-02T21:02:49.397-06:00")
	assert.Equal(t, match.UpdatedAt, &updatedAt)
	assert.Equal(t, match.Player1ID, 38172466)
	assert.Equal(t, match.Player2ID, 38172533)
	assert.Equal(t, match.Player1PrereqMatchID, 0)
	assert.Equal(t, match.Player2PrereqMatchID, 0)
	assert.Equal(t, match.LoserID, 38172533)
	assert.Equal(t, match.WinnerID, 38172466)
	assert.Equal(t, match.ScoresCsv, "0--1")
}
