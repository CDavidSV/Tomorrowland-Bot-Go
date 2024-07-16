package components

import (
	"strconv"
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/player"
	"github.com/bwmarrin/discordgo"
)

func SelectLive(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
	selectedStreamIndex, err := strconv.Atoi(i.MessageComponentData().Values[0])
	if err != nil {
		return
	}

	if len(*bot.LiveStreams) <= selectedStreamIndex {
		return
	}

	selectedStream := (*bot.LiveStreams)[selectedStreamIndex]

	vcState, err := s.State.VoiceState(i.Interaction.GuildID, i.Interaction.Member.User.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Author: &discordgo.MessageEmbedAuthor{
							Name:    "You need to be inside a voice channel to execute this command",
							IconURL: config.ErrorImg,
						},
						Color: config.ErrorColorHex,
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	_, ok := s.VoiceConnections[i.Interaction.GuildID]
	if ok {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Author: &discordgo.MessageEmbedAuthor{
							Name:    "I'm already inside a voice channel, come listen!",
							IconURL: config.ErrorImg,
						},
						Color: config.ErrorColorHex,
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	// Attempt to join the user's vc
	connection, err := s.ChannelVoiceJoin(i.Interaction.GuildID, vcState.ChannelID, false, true)
	if err != nil {
		bot.BotError(err, "play")
		bot.ErrorInteractionResponse(s, i, config.Content{
			Message: "There was a error while attempting to join your voice channel. Please ensure that the bot has enough permissions to join the specified channel",
		}, false, true)
		return
	}

	responseEmbed := &discordgo.MessageEmbed{
		Title: selectedStream.Title,
		URL:   selectedStream.URL,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Tomorrowland Live",
			IconURL: config.StreamGif,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: selectedStream.ThumbnailURL,
		},
		Color: config.MainColorHex,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Tomorrowland Bot",
			IconURL: s.State.User.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Tomorrowland Live",
						IconURL: s.State.User.AvatarURL(""),
					},
					Title:       "Stream selected",
					Description: selectedStream.Title,
					Color:       config.MainColorHex,
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
	if err == nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				responseEmbed,
			},
		})
	}

	err = player.Play(connection, i.ChannelID, selectedStream.ManifestURL)
	if err != nil {
		bot.BotError(err, "play")
		bot.ErrorMessageResponse(s, i, config.Content{
			Message: "I'm sorry, something went wrong. Try again",
		})
	}
}
