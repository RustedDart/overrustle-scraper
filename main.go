package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/RustedDart/overrustle-scraper/overrustle"
)

func main() {
	streamers, err := overrustle.GetStreamers()

	if err != nil {
		log.Fatalln(err)
	}

	var lastStreamer string

	wg := sync.WaitGroup{}
	for key, streamer := range streamers {

		if streamer != lastStreamer {
			lastStreamer = streamer
			percentage := float64(key) / float64(len(streamers)) * 100.0
			fmt.Printf("[%d/%d(%.2f%%)]: %s\n", key+1, len(streamers), percentage, streamer)
		}

		months, err := overrustle.GetMonthsForStreamer(streamer)

		if err != nil {
			log.Fatalln(err)
		}

		for _, month := range months {
			txt, err := overrustle.GetTXTForMonthAndStreamer(streamer, month)

			if err != nil {
				log.Fatalln(err)
			}

			dirName := filepath.Join("scrape", streamer, month)
			if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
				log.Fatalln(err)
			}

			wg.Add(len(txt))

			for _, txtURL := range txt {
				txt := txtURL
				month := month
				streamer := streamer

				go func() {

					defer wg.Done()

					url := fmt.Sprintf("%s/%s chatlog/%s/%s", overrustle.MainURL, streamer, month, txt)
					bytes, err := overrustle.DownloadTXTFile(url)

					if err != nil {
						return
					}

					txtFilePath := filepath.Join(dirName, txt)
					ioutil.WriteFile(txtFilePath, bytes, os.ModePerm)

				}()
			}
		}
	}

	wg.Wait()
	log.Println("Done! :)")
}
