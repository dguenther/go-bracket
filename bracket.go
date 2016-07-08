package bracket

// Client holds API keys and data necessary to make
// calls to different bracket services.
type Client struct {
	challongeUser   string
	challongeAPIKey string
}

// Bracket represents a tournament bracket.
type Bracket struct {
	url string
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
