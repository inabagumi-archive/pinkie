package client

import (
	"time"

	algolia "github.com/algolia/algoliasearch-client-go/algolia/search"
	"github.com/inabagumi/pinkie/pkg/scraper"
	"google.golang.org/api/option"
)

type Client struct {
	Index   *algolia.Index
	Scraper *scraper.Scraper
}

type Options struct {
	AlgoliaAPIKey        string
	AlgoliaApplicationID string
	AlgoliaIndexName     string
	GoogleAPIKey         string
}

func New(opts *Options) (*Client, error) {
	scraper, err := scraper.New(option.WithAPIKey(opts.GoogleAPIKey))
	if err != nil {
		return nil, err
	}

	algoliaClient := algolia.NewClient(opts.AlgoliaApplicationID, opts.AlgoliaAPIKey)
	index := algoliaClient.InitIndex(opts.AlgoliaIndexName)

	c := &Client{
		Index:   index,
		Scraper: scraper,
	}

	return c, nil
}

func (c *Client) Crawl(channelID string, all bool) (algolia.GroupBatchRes, error) {
	opts := &scraper.ScrapeOptions{
		All:             all,
		PublishedBefore: time.Now(),
	}

	videos := c.Scraper.Scrape(channelID, opts)

	res, err := c.Index.SaveObjects(videos)
	if err != nil {
		return algolia.GroupBatchRes{}, err
	}

	return res, nil
}
