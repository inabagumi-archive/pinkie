package scraper

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Scraper struct {
	ctx     context.Context
	service *youtube.Service
}

func New(opts ...option.ClientOption) (*Scraper, error) {
	ctx := context.Background()
	s, err := youtube.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	scraper := &Scraper{ctx: ctx, service: s}

	return scraper, nil
}

type searchOptions struct {
	duration time.Duration
	token    string
	until    time.Time
}

func (s *Scraper) search(channelID string, opts *searchOptions) (*youtube.SearchListResponse, error) {
	since := opts.until.Add(-opts.duration)

	call := s.service.Search.
		List([]string{"id"}).
		ChannelId(channelID).
		MaxResults(50).
		Order("date").
		PageToken(opts.token).
		PublishedAfter(since.Format(time.RFC3339)).
		PublishedBefore(opts.until.Format(time.RFC3339)).
		SafeSearch("none").
		Type("video")

	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Scraper) getVideoList(ids []string) (*youtube.VideoListResponse, error) {
	var strIds string
	for i, id := range ids {
		strIds += id

		if i != len(ids)-1 {
			strIds += ","
		}
	}

	call := s.service.Videos.
		List([]string{"contentDetails", "liveStreamingDetails", "snippet"}).
		Id(strIds).
		MaxResults(50)

	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Scraper) scrape(channelID string, searchOpts *searchOptions) ([]*Video, string, error) {
	log.Printf(`channel_id: "%s", published_before: "%s", page_token: "%s"`,
		channelID, searchOpts.until.Format(time.RFC3339), searchOpts.token)

	searchRes, err := s.search(channelID, searchOpts)
	if err != nil {
		return nil, "", err
	}

	var ids []string
	for _, item := range searchRes.Items {
		ids = append(ids, item.Id.VideoId)
	}

	res, err := s.getVideoList(ids)
	if err != nil {
		return nil, "", err
	}

	var results []*Video

	var (
		mux sync.Mutex
		wg  sync.WaitGroup
	)

	for _, item := range res.Items {
		wg.Add(1)

		go func(item *youtube.Video) {
			defer wg.Done()

			v := NewVideo(item)

			mux.Lock()
			defer mux.Unlock()

			results = append(results, v)
		}(item)
	}

	wg.Wait()

	return results, searchRes.NextPageToken, nil
}

type ScrapeOptions struct {
	All   bool
	Until time.Time
}

func (s *Scraper) Scrape(channelID string, opts *ScrapeOptions) []*Video {
	var (
		days    = 3 * 24 * time.Hour
		date    = opts.Until
		token   = ""
		results []*Video
	)

	if opts.All {
		days = 60 * 24 * time.Hour
	}

	for {
		searchOpts := &searchOptions{
			duration: days,
			token:    token,
			until:    date,
		}
		items, nextToken, err := s.scrape(channelID, searchOpts)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		for _, item := range items {
			results = append(results, item)
		}

		if !opts.All || (token == "" && len(items) < 1) {
			break
		}

		token = nextToken

		if nextToken == "" {
			date = date.Add(-days - 1*time.Second)
		}
	}

	return results
}
