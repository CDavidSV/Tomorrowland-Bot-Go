package commands

import (
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/player"
	"github.com/bwmarrin/discordgo"
)

var StopCommand Command = Command{
	Data: &discordgo.ApplicationCommand{
		Name:         "stop",
		Description:  "Stops the current live stream and disconnects from the voice channel",
		DMPermission: &dmPermissionFalse,
	},
	Callback: func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
		if i.Interaction.Member.Permissions&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "You do not have enough permissions to use this command",
			}, false, true)
			return
		}

		_, err := s.State.VoiceState(i.Interaction.GuildID, i.Interaction.Member.User.ID)
		if err != nil {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "You need to be inside a voice channel to execute this command",
			}, false, true)
			return
		}

		if yes := player.PlayerExists(i.Interaction.GuildID); !yes {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "I'm not inside any voice channel currently",
			}, false, true)
			return
		}

		err = player.Stop(s, i.Interaction.GuildID)
		if err != nil {
			bot.Logger.Error(err.Error(), "command", "stop")
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "I'm sorry, something went wrong. Try again",
			}, false, true)
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

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					responseEmbed,
				},
			},
		})
	},
}
