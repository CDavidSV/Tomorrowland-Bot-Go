package commands

import (
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/player"
	"github.com/bwmarrin/discordgo"
)

var PlayCommand Command = Command{
	Data: &discordgo.ApplicationCommand{
		Name:         "play",
		Description:  "Play the tomorrowland livestream",
		DMPermission: &DMPermissionFalse,
	},
	Callback: func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
		if i.Interaction.Member.Permissions&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "You do not have enough permissions to use this command",
			}, false, true)
			return
		}

		// First check if the user is inside a vc
		vcState, err := s.State.VoiceState(i.Interaction.GuildID, i.Interaction.Member.User.ID)
		if err != nil {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "You need to be inside a voice channel to execute this command",
			}, false, true)
			return
		}

		// Now check if the bot is already inside a vc
		if yes := player.PlayerExists(i.Interaction.GuildID); yes {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "I'm already inside a voice channel, come listen!",
			}, false, false)
			return
		}

		// Attempt to join the users vc
		connection, err := s.ChannelVoiceJoin(i.Interaction.GuildID, vcState.ChannelID, false, true)
		if err != nil {
			bot.BotError(err, "play")
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "There was a error while attempting to join your voice channel. Please ensure that the bot has enough permissions to join the specified channel",
			}, false, true)
			return
		}

		selectedStream := (*bot.LiveStreams)[0]

		responseEmbed := &discordgo.MessageEmbed{
			Title: selectedStream.Title,
			URL:   selectedStream.URL,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Tomorrowland Live",
				IconURL: "https://cdn.discordapp.com/attachments/1107660251286745108/1131297918620532746/stream.gif?ex=6690d30a&is=668f818a&hm=4525fb060ce7c0e6e47dff23b51086d43b724e40d1ecd4185d6d90c8a0c42281&",
			},
			Image: &discordgo.MessageEmbedImage{
				URL: selectedStream.ThumbnailURL,
			},
			Color: 0x7E22DE,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Tomorrowland Bot",
				IconURL: s.State.User.AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					responseEmbed,
				},
			},
		})

		err = player.Play(connection, selectedStream.ManifestURL)
		if err != nil {
			bot.BotError(err, "play")
			bot.ErrorMessageResponse(s, i, config.Content{
				Message: "I'm sorry, something went wrong. Try again",
			})
		}
	},
}
