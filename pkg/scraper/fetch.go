package scraper

import (
	"fmt"
	"net/http"
)

func Fetch(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()

		return nil, fmt.Errorf("scraper: invalid status %d", resp.StatusCode)
	}

	return resp, nil
}
