package cli

import (
	"time"

	"github.com/inabagumi/ytc/internal/util"
	"google.golang.org/api/youtube/v3"
)

func search(service *youtube.Service, channelID string, publishedBefore time.Time, pageToken string) (*youtube.SearchListResponse, error) {
	publishedAfter := publishedBefore.AddDate(0, 0, -60)

	call := service.Search.
		List("id,snippet").
		ChannelId(channelID).
		MaxResults(50).
		Order("date").
		PageToken(pageToken).
		PublishedAfter(publishedAfter.Format(time.RFC3339)).
		PublishedBefore(publishedBefore.Format(time.RFC3339)).
		SafeSearch("none").
		Type("video")

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getVideosByChannelID(service *youtube.Service, channelID string, all bool) []*Video {
	date := time.Now()
	pageToken := ""

	var results []*Video
	for {
		response, err := search(service, channelID, date, pageToken)
		if err != nil {
			break
		}

		for _, item := range response.Items {
			channelID := item.Snippet.ChannelId
			videoID := item.Id.VideoId
			publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
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

		if !all || len(response.Items) < 1 {
			break
		}

		pageToken = response.NextPageToken

		if pageToken != "" {
			date = date.AddDate(0, 0, -61)
		}
	}

	return results
}
