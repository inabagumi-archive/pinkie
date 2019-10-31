package client

import (
	"time"

	algolia "github.com/algolia/algoliasearch-client-go/algolia/search"
	"github.com/inabagumi/pinkie/pkg/scraper"
	"google.golang.org/api/option"
)

type Client struct {
	IndexName string
	Scraper   *scraper.Scraper

	algoliaClient *algolia.Client
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

	c := &Client{
		IndexName:     opts.AlgoliaIndexName,
		Scraper:       scraper,
		algoliaClient: algoliaClient,
	}

	return c, nil
}

func (c *Client) Crawl(channelID string, all bool) (int, error) {
	opts := &scraper.ScrapeOptions{
		All: all,
		PublishedBefore: time.Now(),
	}

	videos := c.Scraper.Scrape(channelID, opts)

	if len(videos) < 1 {
		return 0, nil
	}

	index := c.algoliaClient.InitIndex(c.IndexName)

	res, err := index.SaveObjects(videos)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, batch := range res.Responses {
		count += len(batch.ObjectIDs)
	}

	return count, nil
}
