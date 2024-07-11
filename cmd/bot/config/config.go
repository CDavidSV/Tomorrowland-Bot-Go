package config

import (
	"log/slog"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/youtube"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Logger      *slog.Logger
	LiveStreams *[]youtube.YTVideo
}

type Content struct {
	Message     string
	Description string
}

func (b *Bot) BotError(err error, commandName string) {
	b.Logger.Error(err.Error(), "command", commandName)
}

func (b *Bot) ErrorInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content Content, ephemeral bool) {
	responseEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    content.Message,
			IconURL: "https://cdn.discordapp.com/attachments/1107660251286745108/1107661783457599619/error-icon.png?ex=6691306c&is=668fdeec&hm=1a4215ab4cff4a67bacb59e7e9989248b486d39921e9f0872b232f492602ba9e&",
		},
		Color: 0xe60000,
	}

	if content.Description != "" {
		responseEmbed.Description = content.Description
	}

	data := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			responseEmbed,
		},
	}

	if ephemeral {
		data.Flags = 1 << 6
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
			IconURL: "https://cdn.discordapp.com/attachments/1107660251286745108/1107661783457599619/error-icon.png?ex=6691306c&is=668fdeec&hm=1a4215ab4cff4a67bacb59e7e9989248b486d39921e9f0872b232f492602ba9e&",
		},
		Color: 0xe60000,
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
