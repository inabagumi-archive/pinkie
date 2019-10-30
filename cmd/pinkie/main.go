package main

import (
	"log"
	"os"
	"path/filepath"

	pinkie "github.com/inabagumi/pinkie/pkg/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "dev"

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "")

	app.Version(version)
	app.VersionFlag.Short('v')

	app.HelpFlag.Short('h')

	all := app.Flag("all", "").
		Short('a').
		Bool()

	channels := app.Flag("channel", "A channel ID to scrape.").
		Short('c').
		Required().
		PlaceHolder("<id>").
		Strings()

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Usage(os.Args[1:])
		os.Exit(2)
	}

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

	for _, channel := range *channels {
		count, err := c.Crawl(channel, *all)
		if err != nil {
			log.Printf("error: %v", err)
		}

		if count > 0 {
			log.Printf("Successfully indexed %d videos.", count)
		} else {
			log.Print("Skipped index because there is no videos")
		}
	}
}
