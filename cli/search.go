package cli // import "github.com/inabagumi/ytc/v2/cli"

import (
	"time"

	"github.com/inabagumi/ytc/v2/internal/util"
	"google.golang.org/api/youtube/v3"
)

func search(service *youtube.Service, channelID string, publishedBefore time.Time, pageToken string) (*youtube.SearchListResponse, error) {
	publishedAfter := publishedBefore.AddDate(0, 0, -60)

	call := service.Search.
		List("id").
		ChannelId(channelID).
		MaxResults(50).
		Order("date").
		PageToken(pageToken).
		PublishedAfter(publishedAfter.Format(time.RFC3339)).
		PublishedBefore(publishedBefore.Format(time.RFC3339)).
		SafeSearch("none").
		Type("video")

	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getVideosByChannelID(service *youtube.Service, channelID string, all bool) []*Video {
	date := time.Now()
	pageToken := ""

	var results []*Video
	for {
		searchRes, err := search(service, channelID, date, pageToken)
		if err != nil {
			break
		}

		var ids string

		for i, item := range searchRes.Items {
			ids += item.Id.VideoId
			if i != len(searchRes.Items)-1 {
				ids += ","
			}
		}

		call := service.Videos.
			List("liveStreamingDetails,snippet").
			Id(ids).
			MaxResults(50)

		res, err := call.Do()
		if err != nil {
			break
		}

		for _, item := range res.Items {
			channelID := item.Snippet.ChannelId
			videoID := item.Id

			rawPublishedAt := item.Snippet.PublishedAt

			if item.LiveStreamingDetails != nil {
				if item.LiveStreamingDetails.ActualStartTime != "" {
					rawPublishedAt = item.LiveStreamingDetails.ActualStartTime
				} else if item.LiveStreamingDetails.ScheduledStartTime != "" {
					rawPublishedAt = item.LiveStreamingDetails.ScheduledStartTime
				}
			}

			publishedAt, err := time.Parse(time.RFC3339, rawPublishedAt)
			if err != nil {
				publishedAt = time.Now()
			}

			video := &Video{
				Channel: &Channel{
					ID:    channelID,
					Title: item.Snippet.ChannelTitle,
					URL:   "https://www.youtube.com/channel/" + channelID,
				},
				ID:          videoID,
				ObjectID:    videoID,
				PublishedAt: publishedAt.Unix(),
				Title:       util.NormalizeTitle(item.Snippet.Title),
				URL:         "https://www.youtube.com/watch?v=" + videoID,
			}

			results = append(results, video)
		}

		if !all || len(res.Items) < 1 {
			break
		}

		pageToken = searchRes.NextPageToken

		if pageToken != "" {
			date = date.AddDate(0, 0, -61)
		}
	}

	return results
}
