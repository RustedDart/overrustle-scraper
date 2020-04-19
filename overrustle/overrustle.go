// Package overrustle allows integration with the overrustlelogs.net API
package overrustle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	mainURL     = "https://overrustlelogs.net"
	apiVersion  = "/api/v1"
	channelsURL = mainURL + apiVersion + "/channels.json"
)

var (
	// ErrNotOK is returned when a request does not return the http status code 200
	ErrNotOK = errors.New("requested url did not return status code 200")
)

func doRequest(URL string) ([]byte, error) {
	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// overrustlelogs.net returns 500s if there was no event for subscribers in the streamer/month/subscribers.txt file
	// if resp.StatusCode != 200 {
	// 	return nil, errors.Wrap(ErrNotOK, "overrustle.doRequest.Get")
	// }

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Streamer represents a streamer
type Streamer string

// GetAllStreamers returns all streamers covered by overrustlelogs.net
func GetAllStreamers() ([]Streamer, error) {
	resp, err := doRequest(channelsURL)

	if err != nil {
		return nil, err
	}

	var streamers []Streamer

	if err := json.Unmarshal(resp, &streamers); err != nil {
		return nil, errors.Wrap(err, "overrustle.Streamer.GetAllStreamers.Unmarshal")
	}

	return streamers, nil
}

func (s Streamer) getMonths() ([]string, error) {
	url := fmt.Sprintf("%s%s/%s/months.json", mainURL, apiVersion, s)

	resp, err := doRequest(url)

	if err != nil {
		return nil, err
	}

	var months []string

	if err := json.Unmarshal(resp, &months); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("overrustle.Streamer.getMonths.Unmarshal, url: %s", url))
	}

	return months, nil
}

func (s Streamer) getDaysInMonth(month string) ([]string, error) {
	url := fmt.Sprintf("%s%s/%s/%s/days.json", mainURL, apiVersion, s, month)

	resp, err := doRequest(url)

	if err != nil {
		return nil, err
	}

	var days []string

	if err := json.Unmarshal(resp, &days); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("overrustle.Streamer.getDaysInMonth.Unmarshal, url: %s", url))
	}

	return days, nil
}

// GetAllFileURLs return all files which are associated with the streamer
func (s Streamer) GetAllFileURLs() ([]*File, error) {
	months, err := s.getMonths()

	if err != nil {
		return nil, err
	}

	urls := []*File{}

	for _, month := range months {
		days, err := s.getDaysInMonth(month)

		if err != nil {
			return nil, err
		}

		for _, filename := range days {
			url := fmt.Sprintf("%s/%s chatlog/%s/%s", mainURL, s, month, filename)
			file := &File{
				Day:      filename,
				Month:    month,
				URL:      url,
				Filename: filename,
			}
			urls = append(urls, file)
		}

	}

	return urls, nil
}

// File represents a file saved on overrustlelogs.net
type File struct {
	Month    string
	Day      string
	URL      string
	Filename string
}

// Download downloads the file and return its contents as bytes
func (f *File) Download() ([]byte, error) {
	bytes, err := doRequest(f.URL)

	if err != nil {
		return nil, errors.Wrap(err, "overrustle.File.Download.doRequest")
	}

	return bytes, nil
}
