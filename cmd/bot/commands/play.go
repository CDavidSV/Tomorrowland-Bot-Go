package commands

import (
	"strconv"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/bwmarrin/discordgo"
)

var PlayCommand Command = Command{
	Data: &discordgo.ApplicationCommand{
		Name:         "play",
		Description:  "Play the tomorrowland livestream",
		DMPermission: &dmPermissionFalse,
	},
	Callback: func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
		if i.Interaction.Member.Permissions&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "You do not have enough permissions to use this command",
			}, false, true)
			return
		}

		// First check if the user is inside a vc
		_, err := s.State.VoiceState(i.Interaction.GuildID, i.Interaction.Member.User.ID)
		if err != nil {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "You need to be inside a voice channel to execute this command",
			}, false, true)
			return
		}

		// Now check if the bot is already inside a vc
		_, ok := s.VoiceConnections[i.Interaction.GuildID]
		if ok {
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message:     "I'm already inside a voice channel, come listen!",
				Description: "If you manually disconnected me, wait a few seconds before trying again",
			}, false, true)
			return
		}

		if len(*bot.LiveStreams) == 0 {
			bot.ErrorMessageResponse(s, i, config.Content{
				Message: "I'm sorry, there are no live streams at the moment. Try again later",
			})
			return
		}

		selectOptions := make([]discordgo.SelectMenuOption, len(*bot.LiveStreams))
		for i, v := range *bot.LiveStreams {
			selectOptions[i] = discordgo.SelectMenuOption{
				Label:       v.Title,
				Description: v.ChannelTitle,
				Value:       strconv.Itoa(i),
				Emoji: &discordgo.ComponentEmoji{
					Name: "🔴",
				},
			}
		}

		responseEmbed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Tomorrowland Live",
				IconURL: s.State.User.AvatarURL(""),
			},
			Title: "Select a live stream",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: config.OneWorldRadioImg,
			},
			Color: config.MainColorHex,
		}

		actionRow := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "liveSelect",
					Placeholder: "Choose the livestream you wish to play!",
					Options:     selectOptions,
				},
			},
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					responseEmbed,
				},
				Components: []discordgo.MessageComponent{
					actionRow,
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	},
}
