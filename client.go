package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
)

const crawlerURL = "https://raw.githubusercontent.com/monperrus/crawler-user-agents/master/crawler-user-agents.json"

// Client is a struct that contains the
type Client struct{ botRegexps []*regexp.Regexp }

// Crawler is a struct that represents the JSON response from monperrus'
// crawler data-set
type Crawler struct {
	Pattern      string   `json:"pattern"`
	URL          string   `json:"url,omitempty"`
	Instances    []string `json:"instances"`
	AdditionDate string   `json:"addition_date,omitempty"`
	DependsOn    []string `json:"depends_on,omitempty"`
	Description  string   `json:"description,omitempty"`
}

// New creates a new crawler client
func New() (*Client, error) {
	// Retrieve and decode the crawler dataset into a new slice of Crawler
	var crawlers []Crawler
	if _, err := resty.New().SetRetryCount(5).R().
		ForceContentType("application/json").
		SetResult(&crawlers).Get(crawlerURL); err != nil {
		return nil, fmt.Errorf("retrieving crawler dataset: %w", err)
	}

	// Iterate over all crawlers and return a fully populated Client containing
	// all crawlers in the Clients uamap
	botRegexps := []*regexp.Regexp{}
	for _, crawler := range crawlers {
		// Normalize the pattern so we can use it in Go
		pattern := strings.Replace(crawler.Pattern, `\\`, `\`, -1)

		// Compile the regexp expression to set on the bot botRegexps slice
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile regexp pattern (%s): %w",
				pattern, err)
		}
		botRegexps = append(botRegexps, re)
	}
	return &Client{botRegexps: botRegexps}, nil
}

// IsCrawler returns true if the useragent is found in the map of crawlsers
func (c *Client) IsCrawler(useragent string) bool {
	for _, re := range c.botRegexps {
		if re.Match([]byte(useragent)) {
			return true
		}
	}
	return false
}
