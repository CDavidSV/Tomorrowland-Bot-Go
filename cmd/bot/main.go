package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/commands"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/youtube"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// fetchLivestreamsJob fetches the live streams every 3 hours to prevent expired manifest urls
func fetchLivestreamsJob(list *[]youtube.YTVideo, stop chan struct{}) {
	ticker := time.NewTicker(3 * time.Hour)
	go func() {
		mu := sync.Mutex{}

		for {
			select {
			case <-ticker.C:
				mu.Lock()
				youtube.GetLiveStreams(list)
				mu.Unlock()
			case <-stop: // Stop the ticker
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Flags
	rc := flag.Bool("rc", true, "Reloads application commands")
	guildID := flag.String("guild", "", "Guild ID to test commands")
	flag.Parse()

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("No API Token provided")
	}

	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Error initializing bot: %s", err)
	}

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Fetch live streams
	list := []youtube.YTVideo{}
	youtube.GetLiveStreams(&list)
	fetchLivestreamsJob(&list, make(chan struct{}))

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	bot := &config.Bot{
		Logger:      logger,
		LiveStreams: &list,
	}

	err = session.Open()
	if err != nil {
		log.Fatalf("Error connecting to discord: %s", err)
	}
	log.Println("The Bot is online!")
	defer session.Close()

	commands.LoadCommands(session, bot, *rc, *guildID)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("Bot shutting down...")
}
