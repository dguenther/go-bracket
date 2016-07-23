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
	URL       string
	Name      string
	StartedAt *time.Time
	UpdatedAt *time.Time
	State     string
	Players   []*Player
	Matches   []*Match
}

// Player represents a participant in a tournament.
type Player struct {
	ID   string
	Name string
	Seed int
	Rank int
}

// Match represents a match in a tournament bracket.
type Match struct {
	ID                   string
	Identifier           string
	StartedAt            *time.Time
	UpdatedAt            *time.Time
	Round                int
	State                string
	Player1ID            string
	Player1PrereqMatchID *string
	Player2ID            string
	Player2PrereqMatchID *string
	WinnerID             string
	LoserID              string
	Player1Score         int
	Player2Score         int
}

// NewClient provides a convenient way to instantiate
// an API client.
func NewClient(challongeUser, challongeAPIKey string) *Client {
	return &Client{challongeUser, challongeAPIKey}
}

// FetchBracket takes a URL, calls the appropriate web service for the URL,
// and returns a bracket.
func (c Client) FetchBracket(url string) (*Bracket, error) {
	if isChallongeURL(url) {
		return fetchChallongeBracket(c.challongeUser, c.challongeAPIKey, url)
	}
	if isSmashGGURL(url) {
		return fetchSmashGGBracket(url)
	}
	return nil, nil
}
