package tmrlweb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Performance struct {
	ArtistName string `json:"name,omitempty"`
	Stage      struct {
		Name string `json:"name,omitempty"`
	} `json:"stage,omitempty"`
	Date      string `json:"date,omitempty"`
	Day       string `json:"day,omitempty"`
	StartTime Date   `json:"startTime,omitempty"`
	EndTime   Date   `json:"endTime,omitempty"`
}

type YTVideo struct {
	ID           string `json:"id,omitempty"`
	Title        string `json:"fulltitle,omitempty"`
	ThumbnailURL string `json:"thumbnail,omitempty"`
	ChannelTitle string `json:"channel,omitempty"`
	ChannelURL   string `json:"channel_url,omitempty"`
	URL          string `json:"original_url,omitempty"`
	ManifestURL  string `json:"url,omitempty"`
	Live         bool   `json:"is_live,omitempty"`
}

type Date struct {
	time.Time
}

var DayNameMap = map[string]string{
	"2024-07-19": "WEEK 1 | FRI 19 JUL",
	"2024-07-20": "WEEK 1 | SAT 20 JUL",
	"2024-07-21": "WEEK 1 | SUN 21 JUL",
	"2024-07-26": "WEEK 2 | FRI 26 JUL",
	"2024-07-27": "WEEK 2 | SAT 27 JUL",
	"2024-07-28": "WEEK 2 | SUN 28 JUL",
}

var stages []string = []string{
	"MAINSTAGE",
	"FREEDOM BY BUD",
	"THE ROSE GARDEN",
	"ELIXIR",
	"CAGE",
	"THE RAVE CAVE",
	"PLANAXIS",
	"RISE BY COKE STUDIO",
	"ATMOSPHERE",
	"CORE",
	"CRYSTAL GARDEN",
	"THE LIBRARY",
	"MELODIA BY CORONA",
	"HOUSE OF FORTUNE BY JBL",
	"MOOSEBAR",
}

var performances map[string]map[string][]Performance = map[string]map[string][]Performance{}

var performanceURLS []string = []string{
	"https://artist-lineup-cdn.tomorrowland.com/TLBE24-W1-211903bb-da4c-445d-a1b3-6b17479a9fab.json",
	"https://artist-lineup-cdn.tomorrowland.com/TLBE24-W2-211903bb-da4c-445d-a1b3-6b17479a9fab.json",
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	// "2024-07-19 23:30:00+02:00"
	date, err := time.Parse(`"2006-01-02 15:04:05-07:00"`, string(b))
	if err != nil {
		return err
	}

	d.Time = date
	return nil
}

func GetTimetableImage(date string) (*bytes.Reader, error) {
	cmd := exec.Command("tmrl-web", "timetable", "-d", date)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)

	return r, err
}

func LoadPerformances() error {
	for _, url := range performanceURLS {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var performancesBody struct {
			Performances []Performance `json:"performances"`
		}
		err = json.Unmarshal(body, &performancesBody)
		if err != nil {
			return err
		}

		setPerformancesMap(performancesBody.Performances)
	}

	return nil
}

func setPerformancesMap(p []Performance) {
	// Fill performances map
	for _, v := range p {
		dayStages, ok := performances[v.Date]

		if !ok {
			stage := make(map[string][]Performance)

			stage[v.Stage.Name] = []Performance{v}
			performances[v.Date] = stage
			continue
		}

		performancesPerStage, ok := dayStages[v.Stage.Name]
		if !ok {
			dayStages[v.Stage.Name] = []Performance{v}
			continue
		}

		dayStages[v.Stage.Name] = append(performancesPerStage, v)
	}
}

func GetDayChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}

	for k, v := range DayNameMap {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  v,
			Value: k,
		})
	}

	return choices
}

func GetStageChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}

	for _, v := range stages {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  v,
			Value: v,
		})
	}

	return choices
}

func GetPerformances(date string, stage string) []Performance {
	if v, ok := performances[date][stage]; ok {
		return v
	}

	return nil
}

func GetStages(date string) map[string][]Performance {
	if v, ok := performances[date]; ok {
		return v
	}

	return nil
}

func GetLiveStreams(YTList *[]YTVideo) {
	log.Println("Fetching live streams...")

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

func CurrentlyPLaying(p Performance) bool {
	loc, _ := time.LoadLocation("Europe/Brussels")
	nowBelgium := time.Now().In(loc)

	return nowBelgium.Unix() >= p.StartTime.Unix() && nowBelgium.Unix() <= p.EndTime.Unix()
}
