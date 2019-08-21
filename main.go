package main // import "github.com/inabagumi/ytc"

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	algolia "github.com/algolia/algoliasearch-client-go/algolia/search"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type channelList []string

func (i *channelList) String() string {
	return fmt.Sprint(*i)
}

func (i *channelList) Set(value string) error {
	*i = append(*i, value)

	return nil
}

var all bool
var channels channelList

func init() {
	flag.BoolVar(&all, "all", false, "Default: false.")
	flag.Var(&channels, "channel", "Set the index target channel.")
}

// A Channel represents the channel of YouTube.
type Channel struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// A Video represents the video of YouTube.
type Video struct {
	Channel     *Channel `json:"channel"`
	ID          string   `json:"id"`
	ObjectID    string   `json:"objectID"`
	PublishedAt int64    `json:"publishedAt"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
}

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

func normalizeTitle(title string) string {
	re := regexp.MustCompile(`\s*【[^】]+\s*\/\s*あにまーれ】?\s*`)

	title = re.ReplaceAllString(title, " ")
	title = strings.TrimSpace(title)
	title = html.UnescapeString(title)

	return title
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
				Title:       normalizeTitle(item.Snippet.Title),
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

func main() {
	flag.Parse()

	client := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("GOOGLE_API_KEY")},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	var results []*Video
	for _, channelID := range channels {
		for _, video := range getVideosByChannelID(service, channelID, all) {
			results = append(results, video)
		}
	}

	searchClient := algolia.NewClient(os.Getenv("ALGOLIA_APPLICATION_ID"), os.Getenv("ALGOLIA_API_KEY"))
	index := searchClient.InitIndex(os.Getenv("ALGOLIA_INDEX_NAME"))

	_, err = index.SaveObjects(results)
	if err != nil {
		log.Fatalf("Error indexing results: %v", err)
	}
}
