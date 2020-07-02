package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"s32x.com/httpclient"
)

const crawlerURL = "https://raw.githubusercontent.com/monperrus/crawler-user-agents/master/crawler-user-agents.json"

// Client is a struct that contains the
type Client struct{ botRegexps []*regexp.Regexp }

// Crawlers is a struct that represents the JSON response from monperrus'
// crawler data-set
type Crawlers struct {
	Pattern   string   `json:"pattern"`
	URL       string   `json:"url"`
	Instances []string `json:"instances"`
}

// New creates a new crawler client
func New() (*Client, error) {
	// Retrieve the crawler dataset for map population
	var crawlers []Crawlers
	if err := httpclient.New().
		Get(crawlerURL).
		WithExpectedStatus(200).
		WithRetry(5).
		JSON(&crawlers); err != nil {
		return nil, fmt.Errorf("failed to retrieve crawler dataset: %w", err)
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
