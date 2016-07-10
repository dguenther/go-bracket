package bracket

import "time"

// Client holds API keys and data necessary to make
// calls to different bracket services.
type Client struct {
	challongeUser   string
	challongeAPIKey string
}

// Bracket represents a tournament bracket.
type Bracket struct {
	URL     string
	Name    string
	Players []*Player
	Matches []*Match
}

// Player represents a participant in a tournament.
type Player struct {
	Name string
	Seed int
	Rank int
}

// Match represents a match in a tournament bracket.
type Match struct {
	ID           string
	Identifier   string
	UpdatedAt    *time.Time
	Round        int
	State        string
	Player1ID    string
	Player2ID    string
	WinnerID     string
	LoserID      string
	Player1Score int
	Player2Score int
}

// NewClient provides a convenient way to instantiate
// an API client.
func NewClient(challongeUser, challongeAPIKey string) *Client {
	return &Client{challongeUser, challongeAPIKey}
}

// FetchBracket takes a URL, calls the appropriate web service for the URL,
// and returns a bracket.
func (c Client) FetchBracket(url string) *Bracket {
	if isChallongeURL(url) {
		return fetchChallongeBracket(c.challongeUser, c.challongeAPIKey, url)
	}
	return nil
}
