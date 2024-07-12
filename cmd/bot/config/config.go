package config

import (
	"bytes"
	"encoding/base64"
	"log/slog"
	"os/exec"

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

func (b *Bot) ErrorInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content Content, edit bool, ephemeral bool) {
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

func (bot *Bot) GetTimetable(date string) (*bytes.Reader, error) {
	cmd := exec.Command("tmrl-web", "-d", date)

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
