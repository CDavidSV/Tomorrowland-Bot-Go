package config

import (
	"log/slog"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/tmrlweb"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Logger      *slog.Logger
	LiveStreams *[]tmrlweb.YTVideo
}

type Content struct {
	Message     string
	Description string
}

var ErrorColorHex int = 0xe60000
var MainColorHex int = 0x7E22DE
var OneWorldRadioImg string = "https://d384fynlilbsl.cloudfront.net/one_world_radio_logo.png"
var ErrorImg string = "https://d384fynlilbsl.cloudfront.net/error-icon.png"
var StreamGif string = "https://d384fynlilbsl.cloudfront.net/stream.gif"

func (b *Bot) BotError(err error, commandName string) {
	b.Logger.Error(err.Error(), "command", commandName)
}

func (b *Bot) ErrorInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content Content, edit bool, ephemeral bool) {
	responseEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    content.Message,
			IconURL: ErrorImg,
		},
		Color: ErrorColorHex,
	}

	if content.Description != "" {
		responseEmbed.Description = content.Description
	}

	if edit {
		c := ""
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &c,
			Embeds:  &[]*discordgo.MessageEmbed{responseEmbed},
		})
		return
	}

	data := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			responseEmbed,
		},
	}

	if ephemeral {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func (b *Bot) ErrorMessageResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content Content) error {
	responseEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    content.Message,
			IconURL: ErrorImg,
		},
		Color: ErrorColorHex,
	}

	if content.Description != "" {
		responseEmbed.Description = content.Description
	}

	_, err := s.ChannelMessageSendEmbed(i.Interaction.ChannelID, responseEmbed)
	if err != nil {
		return err
	}

	return nil
}
