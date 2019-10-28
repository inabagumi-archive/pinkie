package cli // import "github.com/inabagumi/ytc/v2/cli"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	algolia "github.com/algolia/algoliasearch-client-go/algolia/search"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

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

// A Client represents the client of CLI.
type Client struct {
	Name    string
	Version string

	youtubeService *youtube.Service
	algoliaClient  *algolia.Client
}

// NewClient returns a new CLI client.
func NewClient(name string, version string) (*Client, error) {
	httpClient := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("GOOGLE_API_KEY")},
	}

	youtubeService, err := youtube.New(httpClient)
	if err != nil {
		return nil, err
	}

	algoliaClient := algolia.NewClient(os.Getenv("ALGOLIA_APPLICATION_ID"), os.Getenv("ALGOLIA_API_KEY"))

	return &Client{
		Name:    name,
		Version: version,

		youtubeService: youtubeService,
		algoliaClient:  algoliaClient,
	}, nil
}

type channelList []string

func (i *channelList) String() string {
	return fmt.Sprint(*i)
}

func (i *channelList) Set(value string) error {
	*i = append(*i, value)

	return nil
}

type config struct {
	all          bool
	channels     channelList
	printVersion bool
}

// Run executes a client tasks and returns an exit code.
func (c *Client) Run(args []string) {
	fs := flag.NewFlagSet(c.Name, flag.ExitOnError)

	conf := config{}

	fs.BoolVar(&conf.all, "all", false, "Also index videos older than 60 days. Default: false.")
	fs.Var(&conf.channels, "channel", "Set the index target channel.")
	fs.BoolVar(&conf.printVersion, "version", false, "Print the version number.")

	if err := fs.Parse(args); err != nil {
		log.Fatalf("Error: %v", err)
	}

	if conf.printVersion {
		fmt.Printf("%s version %s\n", fs.Name(), c.Version)
		return
	}

	var results []*Video
	for _, channelID := range conf.channels {
		for _, video := range getVideosByChannelID(c.youtubeService, channelID, conf.all) {
			results = append(results, video)
		}
	}

	c.index(os.Getenv("ALGOLIA_INDEX_NAME"), results)
}

func (c *Client) index(indexName string, videos []*Video) {
	index := c.algoliaClient.InitIndex(indexName)

	res, err := index.SaveObjects(videos)
	if err != nil {
		log.Fatalf("Error indexing results: %v", err)
	}

	count := 0
	for _, batch := range res.Responses {
		count += len(batch.ObjectIDs)
	}

	log.Printf("Successfully indexed %d videos.", count)
}
