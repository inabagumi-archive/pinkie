package client

import (
	"context"

	algolia "github.com/algolia/algoliasearch-client-go/algolia/search"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type client struct {
	algoliaClient  *algolia.Client
	indexName      string
	youtubeService *youtube.Service
}

type Options struct {
	AlgoliaAPIKey        string
	AlgoliaApplicationID string
	AlgoliaIndexName     string
	GoogleAPIKey         string
}

func New(opts *Options) (*client, error) {
	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(opts.GoogleAPIKey))
	if err != nil {
		return nil, err
	}

	algoliaClient := algolia.NewClient(opts.AlgoliaApplicationID, opts.AlgoliaAPIKey)

	c := &client{
		algoliaClient:  algoliaClient,
		youtubeService: youtubeService,
	}

	return c, nil
}

func (c *client) Crawl(channelID string, all bool) (int, error) {
	videos := c.Scrape(channelID, all)

	if len(videos) < 1 {
		return 0, nil
	}

	index := c.algoliaClient.InitIndex(c.indexName)

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
