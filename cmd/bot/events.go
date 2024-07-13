package main

import (
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/commands"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/player"
	"github.com/bwmarrin/discordgo"
)

type EventHandler interface{}

func LoadEvents(s *discordgo.Session, bot *config.Bot) {
	bot.Logger.Info("Loading events...")

	s.AddHandler(ReadyEvent(bot))
	s.AddHandler(InteractionCreate(bot))
	s.AddHandler(VoiceStateChangeEvent(bot))
}

func ReadyEvent(bot *config.Bot) EventHandler {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		bot.Logger.Info("Bot is online!", "ID", s.State.User.ID, "Username", r.User.Username)
	}
}

func InteractionCreate(bot *config.Bot) EventHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if c, ok := commands.SlashCommands[i.ApplicationCommandData().Name]; ok {
			// Log command in console
			bot.Logger.Info("Command executed", "command", i.ApplicationCommandData().Name, "user", i.Interaction.Member.User.ID, "guild", i.Interaction.GuildID)

			c.Callback(s, i, bot)
		}
	}
}

func VoiceStateChangeEvent(bot *config.Bot) EventHandler {
	return func(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
		// Check if the state update is for the bot and that it has disconnected
		if v.UserID == s.State.User.ID && v.BeforeUpdate != nil && v.ChannelID == "" {
			// Wait for 2 seconds before stopping the player
			time.Sleep(2 * time.Second)

			channelID, err := player.GetOriginalChannelID(v.GuildID)
			if err != nil {
				return
			}

			responseEmbed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Tomorrowland Live",
					IconURL: s.State.User.AvatarURL(""),
				},
				Title:     "Thanks for tuning in. See you next time!",
				Timestamp: time.Now().Format(time.RFC3339),
				Color:     config.MainColorHex,
			}

			s.ChannelMessageSendEmbed(channelID, responseEmbed)
		}
	}
}
