package player

import (
	"fmt"
	"log"

	"github.com/CDavidSV/dca"
	"github.com/bwmarrin/discordgo"
)

// Right now using dca from "https://github.com/jonas747/dca", looking to make a custom implementation in the future

type Player struct {
	encoderSession    *dca.EncodeSession
	streamSession     *dca.StreamingSession
	originalChannelID string
}

var playerSessions map[string]Player = map[string]Player{}

func PlayerExists(guildID string) bool {
	if _, ok := playerSessions[guildID]; ok {
		return true
	}

	return false
}

func GetOriginalChannelID(guildID string) (string, error) {
	playerSesssion, ok := playerSessions[guildID]
	if !ok {
		return "", fmt.Errorf("there is no player session in guild %v", guildID)
	}

	return playerSesssion.originalChannelID, nil
}

func Play(vc *discordgo.VoiceConnection, channelID string, url string) error {
	// Options for encoding
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.Application = "lowdelay"
	options.Threads = 4

	encodingSession, err := dca.EncodeFile(url, options)
	if err != nil {
		return fmt.Errorf("error decoding stream: %v", err)
	}

	// Check if a current player session exists
	if ps, ok := playerSessions[vc.GuildID]; ok {
		// Stop the current player and delete the session
		ps.encoderSession.Cleanup()

		// Change the encoder session
		ps.encoderSession = encodingSession
	}

	err = vc.Speaking(true)
	if err != nil {
		encodingSession.Cleanup()
		return fmt.Errorf("error speaking: %v", err)
	}

	// Start the stream
	go startStream(encodingSession, vc, channelID)

	return nil
}

func startStream(encodingSession *dca.EncodeSession, vc *discordgo.VoiceConnection, channelID string) {
	// Wait for stream to be done
	streamDone := make(chan error)
	stream := dca.NewStream(encodingSession, vc, streamDone)

	// Add stream data to players sessions map
	playerSessions[vc.GuildID] = Player{
		encoderSession:    encodingSession,
		streamSession:     stream,
		originalChannelID: channelID,
	}

	defer encodingSession.Cleanup()

	if err := <-streamDone; err != nil {
		log.Println("stream error:", err)

		// Handle stream expiration
	}
}

func Stop(s *discordgo.Session, guildID string) error {
	// Stop the player
	player, ok := playerSessions[guildID]
	if !ok {
		return nil
	}

	// Clean the encoding session to prevent memory leaks
	player.encoderSession.Cleanup()

	// Remove from players map
	delete(playerSessions, guildID)

	vc, ok := s.VoiceConnections[guildID]
	if !ok {
		return nil
	}
	// Attempt to disconnect
	err := vc.Disconnect()
	if err != nil {
		return fmt.Errorf("error attempting to disconnect from voice channel: %v", err)
	}

	return nil
}
