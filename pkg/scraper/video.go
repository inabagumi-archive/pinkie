package scraper

import (
	"fmt"
	"html"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Channel     *Channel   `json:"channel"`
	Duration    string     `json:"duration"`
	ID          string     `json:"id"`
	ObjectID    string     `json:"objectID"`
	PublishedAt int64      `json:"publishedAt"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
}

func NewVideo(item *youtube.Video) *Video {
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

	channel := &Channel{
		ID:    item.Snippet.ChannelId,
		Title: item.Snippet.ChannelTitle,
		URL:   fmt.Sprintf("https://www.youtube.com/channel/%s", item.Snippet.ChannelId),
	}

	b := fmt.Sprintf("https://i.ytimg.com/vi/%s/", item.Id)

	t, err := NewThumbnail(b + "maxresdefault.jpg")
	if t == nil {
		t, _ = NewThumbnail(b + "hqdefault.jpg")
	}

	video := &Video{
		Channel:     channel,
		Duration:    item.ContentDetails.Duration,
		ID:          item.Id,
		ObjectID:    item.Id,
		PublishedAt: publishedAt.Unix(),
		Thumbnail:   t,
		Title:       html.UnescapeString(item.Snippet.Title),
		URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id),
	}

	return video
}
