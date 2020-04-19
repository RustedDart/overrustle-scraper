package overrustle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	MainURL     = "https://overrustlelogs.net"
	apiVersion  = "/api/v1"
	channelsURL = MainURL + apiVersion + "/channels.json"
)

var (
	ErrNotOK = errors.New("requested url did not return status code 200")
)

func GetStreamers() ([]string, error) {
	resp, err := http.Get(channelsURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Wrap(ErrNotOK, "overrustle.GetStreamers.Get")
	}

	var streamers []string

	if err := json.NewDecoder(resp.Body).Decode(&streamers); err != nil {
		return nil, errors.Wrap(err, "overrustle.GetStreamers.Decode")
	}

	return streamers, nil
}

func GetMonthsForStreamer(streamer string) ([]string, error) {
	url := fmt.Sprintf("%s%s/%s/months.json", MainURL, apiVersion, streamer)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Wrap(ErrNotOK, "overrustle.GetMonthsForStreamer.Get")
	}

	var months []string

	if err := json.NewDecoder(resp.Body).Decode(&months); err != nil {
		return nil, errors.Wrap(err, "overrustle.GetMonthsForStreamer.Decode")
	}

	return months, nil
}

func GetTXTForMonthAndStreamer(streamer string, month string) ([]string, error) {
	url := fmt.Sprintf("%s%s/%s/%s/days.json", MainURL, apiVersion, streamer, month)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Wrap(ErrNotOK, "overrustle.GetTXTForMonthAndStreamer.Get")
	}

	var txtFiles []string

	if err := json.NewDecoder(resp.Body).Decode(&txtFiles); err != nil {
		return nil, errors.Wrap(err, "overrustle.GetTXTForMonthAndStreamer.Decode")
	}

	return txtFiles, nil
}

func DownloadTXTFile(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Wrap(ErrNotOK, "overrustle.DownloadTXTFile.Get")
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
