package scraper

import (
	"context"
	"log"
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
	pageToken       string
	publishedBefore time.Time
}

func (s *Scraper) search(channelID string, opts *searchOptions) (*youtube.SearchListResponse, error) {
	publishedAfter := opts.publishedBefore.AddDate(0, 0, -60)

	call := s.service.Search.
		List("id").
		ChannelId(channelID).
		MaxResults(50).
		Order("date").
		PageToken(opts.pageToken).
		PublishedAfter(publishedAfter.Format(time.RFC3339)).
		PublishedBefore(opts.publishedBefore.Format(time.RFC3339)).
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
		List("liveStreamingDetails,snippet").
		Id(strIds).
		MaxResults(50)

	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

type ScrapeOptions struct {
	All             bool
	PublishedBefore time.Time
}

func (s *Scraper) Scrape(channelID string, opts *ScrapeOptions) []*Video {
	var (
		date      = opts.PublishedBefore
		pageToken = ""
		results   []*Video
	)

	for {
		log.Printf(`channel_id: "%s", published_before: "%s", page_token: "%s"`,
			channelID, date.Format(time.RFC3339), pageToken)

		searchOpts := &searchOptions{
			pageToken:       pageToken,
			publishedBefore: date,
		}

		searchRes, err := s.search(channelID, searchOpts)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		var ids []string
		for _, item := range searchRes.Items {
			ids = append(ids, item.Id.VideoId)
		}

		res, err := s.getVideoList(ids)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		for _, item := range res.Items {
			results = append(results, normalize(item))
		}

		if !opts.All || (pageToken == "" && len(searchRes.Items) < 1) {
			break
		}

		pageToken = searchRes.NextPageToken

		if pageToken == "" {
			date = date.AddDate(0, 0, -61)
		}
	}

	return results
}
