package internetarchive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

// WaybackResponse represents the structure of the response from the Wayback Machine API
type WaybackResponse struct {
	ArchivedSnapshots struct {
		Closest struct {
			URL       string `json:"url"`
			Timestamp string `json:"timestamp"`
			Status    string `json:"status"`
		} `json:"closest"`
	} `json:"archived_snapshots"`
}

// GetLatestPageSnapshot fetches the latest snapshot of the given URL from the Wayback Machine
func GetLatestPageSnapshot(pageURL string, verbose bool) (string, error) {
	// Prepare the request URL with query parameters
	waybackResp, err := getWaybackAPIResponse(pageURL, verbose)
	if err != nil {
		return "", err
	}

	// Check if a snapshot is available
	if waybackResp.ArchivedSnapshots.Closest.URL == "" {
		// try again with either a slash at the end or without it and see if you get a hit
		if strings.HasSuffix(pageURL, "/") {
			pageURL = pageURL[:len(pageURL)-1]

			if verbose {
				logger.WriteInfo("Did not find original URL, so now trying for the URL without the ending slash.")
			}
		} else {
			pageURL += "/"

			if verbose {
				logger.WriteInfo("Did not find original URL, so now trying for the URL with an ending slash.")
			}
		}

		waybackResp, err = getWaybackAPIResponse(pageURL, verbose)
		if err != nil {
			return "", err
		}

		if waybackResp.ArchivedSnapshots.Closest.URL == "" {
			return "", fmt.Errorf("no snapshots available for this URL")
		}
	}

	// Return the URL of the latest snapshot
	return waybackResp.ArchivedSnapshots.Closest.URL, nil
}

func getWaybackAPIResponse(pageURL string, verbose bool) (WaybackResponse, error) {
	var waybackResp WaybackResponse

	reqURL, err := url.Parse(internetArchiveAPIUrl)
	if err != nil {
		return waybackResp, fmt.Errorf("failed to parse API URL: %v", err)
	}
	q := reqURL.Query()
	q.Set("url", pageURL)
	reqURL.RawQuery = q.Encode()

	if verbose {
		logger.WriteInfof("Trying to get latest snapshot of %q using url %q\n", pageURL, reqURL.String())
	}

	// Make the request to the Wayback Machine API
	resp, err := http.Get(reqURL.String())
	if err != nil {
		return waybackResp, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return waybackResp, fmt.Errorf("API request returned non-200 status: %s", resp.Status)
	}

	// Parse the JSON response
	if err := json.NewDecoder(resp.Body).Decode(&waybackResp); err != nil {
		return waybackResp, fmt.Errorf("failed to decode API response: %v", err)
	}

	return waybackResp, nil
}
