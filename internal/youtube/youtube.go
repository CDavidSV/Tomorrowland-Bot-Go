package youtube

import (
	"encoding/json"
	"log"
	"os/exec"
	"sync"
)

type YTVideo struct {
	ID           string `json:"id"`
	Title        string `json:"fulltitle"`
	ThumbnailURL string `json:"thumbnail"`
	ChannelTitle string `json:"channel"`
	ChannelURL   string `json:"channel_url"`
	URL          string `json:"original_url"`
	ManifestURL  string `json:"url"`
	Live         bool   `json:"is_live"`
}

func GetLiveStreams(YTList *[]YTVideo) {
	log.Println("Fetching live streams")

	livestreamsURLS := []string{
		"https://www.youtube.com/watch?v=wBgSH-CGPzg",
		"https://www.youtube.com/watch?v=aZT73SdhXok",
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Use yt-dl to fetch download urls from yt
	videoData := make(chan *YTVideo, len(livestreamsURLS))
	for _, v := range livestreamsURLS {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			cmd := exec.Command("yt-dlp", "--dump-json", v, "-f", "bestaudio/best")

			out, err := cmd.Output()
			if err != nil {
				log.Printf("Error fetching youtube stream url for %v: %v", url, err)
				return
			}

			videoMetadata := &YTVideo{}
			err = json.Unmarshal(out, videoMetadata)
			if err != nil {
				log.Printf("Error parsing video json dump")
				return
			}

			videoData <- videoMetadata
		}(v)
	}

	go func() {
		wg.Wait()
		close(videoData)
	}()

	for video := range videoData {
		log.Printf("Fetched live stream: %v", video.Title)
		mu.Lock()
		*YTList = append(*YTList, *video)
		mu.Unlock()
	}
}
