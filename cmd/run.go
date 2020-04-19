// cmd is the cli interface
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/RustedDart/overrustle-scraper/overrustle"
)

type downloadBlock struct {
	File      *overrustle.File
	TargetDir string
	Timeout   time.Duration
}

// Run runs the cli app
func Run() {
	wg := &sync.WaitGroup{}
	var lastStreamer string

	streamers, err := overrustle.GetAllStreamers()

	if err != nil {
		log.Fatal(err)
	}

	downloadPipeline := make(chan *downloadBlock)
	defer close(downloadPipeline)

	startDownloader(downloadPipeline, wg)

	for key, streamer := range streamers {

		if string(streamer) != lastStreamer {
			lastStreamer = string(streamer)
			percentage := float64(key) / float64(len(streamers)) * 100.0
			fmt.Printf("[%d/%d(%.2f%%)]: %s\n", key+1, len(streamers), percentage, streamer)
		}

		var files []*overrustle.File

		for i := 0; i < 4; i++ {
			files, err = streamer.GetAllFileURLs()

			if err != nil {
				if i == 3 {
					log.Fatal(err)
				}

				<-time.After(time.Second * 3)
				continue
			} else {
				break
			}
		}

		wg.Add(len(files))

		for _, file := range files {
			dirName := filepath.Join("scrape", string(streamer), file.Month)
			if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
				log.Fatalln(err)
			}

			downloadBlock := &downloadBlock{
				File:      file,
				Timeout:   time.Duration(time.Second * 0),
				TargetDir: dirName,
			}

			downloadPipeline <- downloadBlock
		}

		wg.Wait()
		fmt.Println("Done with", streamer)
	}

	fmt.Println("Done! :)")
}

func startDownloader(pipeline chan *downloadBlock, wg *sync.WaitGroup) {

	for i := 0; i < 12; i++ {
		go func() {
			for block := range pipeline {
				block := block

				func() {
					defer wg.Done()

					<-time.After(block.Timeout)

					bytes, err := block.File.Download()

					if err != nil {

						if block.Timeout == time.Duration(time.Second*0) {
							block.Timeout = time.Duration(time.Second * 2)
						} else {
							block.Timeout = time.Duration(block.Timeout * 2)
						}

						if block.Timeout > time.Duration(time.Minute*1) {
							log.Fatalln("exit because https://overrustlelogs.net could not be reached for", block.Timeout)
						}

						log.Println("failed to download", block.File.URL, "timeout for:", block.Timeout)

						wg.Add(1)
						go func() { pipeline <- block }()

					}

					txtFilePath := filepath.Join(block.TargetDir, block.File.Filename)

					if _, err := os.Stat(txtFilePath); err != nil {
						err = ioutil.WriteFile(txtFilePath, bytes, os.ModePerm)

						if err != nil {
							log.Fatal(err)
						}
					}
				}()
			}
		}()
	}
}
