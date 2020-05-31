package client

import (
	"log"
	"sync"
	"time"

	algolia "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/inabagumi/pinkie/v4/pkg/crawler"
	"github.com/inabagumi/pinkie/v4/pkg/scraper"
)

type Client struct {
	crawler *crawler.Crawler
}

func New(opts *crawler.Options) (*Client, error) {
	crawler, err := crawler.New(opts)
	if err != nil {
		return nil, err
	}

	return &Client{crawler: crawler}, nil
}

func (c *Client) Crawl(channelID string, all bool) (algolia.GroupBatchRes, error) {
	opts := &scraper.ScrapeOptions{
		All:   all,
		Until: time.Now(),
	}

	return c.crawler.Crawl(channelID, opts)
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
