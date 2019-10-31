package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/api/option"
)

func TestScraper_Scrape(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/search.json")
	})
	mux.HandleFunc("/videos", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/videos.json")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	scraper, err := New(option.WithAPIKey("DUMMYAPIKEY"), option.WithEndpoint(server.URL))
	if err != nil {
		t.Fatal(err)
	}

	opts := &ScrapeOptions{
		All:             false,
		PublishedBefore: time.Date(2019, time.October, 31, 8, 0, 0, 0, time.UTC),
	}

	results := scraper.Scrape("UC0Owc36U9lOyi9Gx9Ic-4qg", opts)

	if got, want := len(results), 50; got != want {
		t.Errorf("results is %v video(s), want %v video(s)", got, want)
	}
}
