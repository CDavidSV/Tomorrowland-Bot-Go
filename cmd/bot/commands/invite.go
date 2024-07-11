package commands

import (
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/bwmarrin/discordgo"
)

var InviteCommand Command = Command{
	Data: &discordgo.ApplicationCommand{
		Name:         "invite",
		Description:  "Invite me to your server",
		DMPermission: &DMPermissionTrue,
	},
	Callback: func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
		responseEmbed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Tomorrowland Bot Invite",
				IconURL: s.State.User.AvatarURL(""),
			},
			Description: "Bring the party to your server!",
			Color:       0xb917ff,
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		actionRow := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Invite Tomorrowland 24/7",
					URL:   "https://discord.com/api/oauth2/authorize?client_id=1000497170434236457&permissions=274914618624&scope=bot%20applications.commands",
					Style: discordgo.LinkButton,
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
			},
		})
	},
}
