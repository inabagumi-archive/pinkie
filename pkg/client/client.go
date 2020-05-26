package client

import (
	"log"
	"sync"
	"time"

	algolia "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
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

func (c *Client) Scrape(channelID string, opts *scraper.ScrapeOptions) interface{} {
	return c.Scraper.Scrape(channelID, opts)
}

func (c *Client) Crawl(channelID string, all bool) (algolia.GroupBatchRes, error) {
	opts := &scraper.ScrapeOptions{
		All:             all,
		PublishedBefore: time.Now(),
	}

	results := c.Scrape(channelID, opts)

	res, err := c.Index.SaveObjects(results)
	if err != nil {
		return algolia.GroupBatchRes{}, err
	}

	return res, nil
}

func (c *Client) Run(channels []string, all bool) {
	var wg sync.WaitGroup

	count := 0
	for _, channel := range channels {
		wg.Add(1)

		go func(channel string) {
			defer wg.Done()

			res, err := c.Crawl(channel, all)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}

			for _, batchRes := range res.Responses {
				count += len(batchRes.ObjectIDs)
			}
		}(channel)
	}

	wg.Wait()

	log.Printf("Successfully indexed %d videos.", count)
}
