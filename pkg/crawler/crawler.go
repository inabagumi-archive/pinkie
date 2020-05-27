package crawler

import (
	algolia "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/inabagumi/pinkie/pkg/scraper"
	"google.golang.org/api/option"
)

type Crawler struct {
	index   *algolia.Index
	scraper *scraper.Scraper
}

type Options struct {
	AlgoliaAPIKey        string
	AlgoliaApplicationID string
	AlgoliaIndexName     string
	GoogleAPIKey         string
}

func New(opts *Options) (*Crawler, error) {
	client := algolia.NewClient(opts.AlgoliaApplicationID, opts.AlgoliaAPIKey)
	index := client.InitIndex(opts.AlgoliaIndexName)

	scraper, err := scraper.New(option.WithAPIKey(opts.GoogleAPIKey))
	if err != nil {
		return nil, err
	}

	return &Crawler{scraper: scraper, index: index}, nil
}

func (c *Crawler) Scrape(channelID string, opts *scraper.ScrapeOptions) interface{} {
	return c.scraper.Scrape(channelID, opts)
}

func (c *Crawler) Crawl(channelID string, opts *scraper.ScrapeOptions) (algolia.GroupBatchRes, error) {
	results := c.Scrape(channelID, opts)

	res, err := c.index.SaveObjects(results)
	if err != nil {
		return algolia.GroupBatchRes{}, err
	}

	return res, nil
}
