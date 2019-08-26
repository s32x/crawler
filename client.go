package crawler

import (
	"github.com/pkg/errors"
	"s32x.com/httpclient"
)

const crawlerURL = "https://raw.githubusercontent.com/monperrus/crawler-user-agents/master/crawler-user-agents.json"

// Client is a struct that represents an in memory crawler database
type Client struct{ uamap map[string]string }

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
		return nil, errors.Wrap(err, "Failed to retrieve crawler dataset")
	}

	// Iterate over all crawlers and return a fully populated Client containing
	// all crawlers in the Clients uamap
	uamap := make(map[string]string)
	for _, crawler := range crawlers {
		for _, useragent := range crawler.Instances {
			uamap[useragent] = crawler.Pattern
		}
	}
	return &Client{uamap: uamap}, nil
}

// IsCrawler returns true if the useragent is found in the map of crawlsers
func (c *Client) IsCrawler(useragent string) bool {
	_, ok := c.uamap[useragent]
	return ok
}

// Crawler attempts to find the useragent in the uamap and returns a true bool
// if the useragent is a crawler and the name of the crawler
func (c *Client) Crawler(useragent string) (string, bool) {
	name, ok := c.uamap[useragent]
	return name, ok
}
