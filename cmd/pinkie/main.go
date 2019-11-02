package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	pinkie "github.com/inabagumi/pinkie/pkg/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "dev"

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "The Pinkie is a crawler that uses the YouTube Data API.")

	app.Version(version)
	app.VersionFlag.Short('v')

	app.HelpFlag.Short('h')

	all := app.Flag("all", "Fetch all videos of channel.").
		Short('a').
		Bool()

	channels := app.Flag("channel", "A channel ID to scrape.").
		Short('c').
		Required().
		PlaceHolder("<id>").
		Strings()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	opts := &pinkie.Options{
		AlgoliaAPIKey:        os.Getenv("ALGOLIA_API_KEY"),
		AlgoliaApplicationID: os.Getenv("ALGOLIA_APPLICATION_ID"),
		AlgoliaIndexName:     os.Getenv("ALGOLIA_INDEX_NAME"),
		GoogleAPIKey:         os.Getenv("GOOGLE_API_KEY"),
	}

	c, err := pinkie.New(opts)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	count := 0
	for _, channel := range *channels {
		wg.Add(1)

		go func(channel string) {
			res, err := c.Crawl(channel, *all)
			if err != nil {
				log.Printf("error: %v", err)
			}

			for _, batchRes := range res.Responses {
				count += len(batchRes.ObjectIDs)
			}

			wg.Done()
		}(channel)
	}

	wg.Wait()

	log.Printf("Successfully indexed %d videos.", count)
}
