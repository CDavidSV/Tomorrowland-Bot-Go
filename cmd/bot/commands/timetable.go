package commands

import (
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/bwmarrin/discordgo"
)

var TimetableCommand Command = Command{
	Data: &discordgo.ApplicationCommand{
		Name:         "timetable",
		Description:  "Get the timetable for a specific day",
		DMPermission: &DMPermissionTrue,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "day",
				Description: "Select a day",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "WEEK 1 - FRI 19 JUL",
						Value: "2024-07-19",
					},
					{
						Name:  "WEEK 1 - SAT 20 JUL",
						Value: "2024-07-20",
					},
					{
						Name:  "WEEK 1 - SUN 21 JUL",
						Value: "2024-07-21",
					},
					{
						Name:  "WEEK 2 - FRI 26 JUL",
						Value: "2024-07-26",
					},
					{
						Name:  "WEEK 2 - SAT 27 JUL",
						Value: "2024-07-27",
					},
					{
						Name:  "WEEK 2 - SUN 28 JUL",
						Value: "2024-07-28",
					},
				},
				Required: true,
			},
		},
	},
	Callback: func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
		selectedDay := i.Interaction.ApplicationCommandData().Options[0]

		var dayNameMap = map[string]string{
			"2024-07-19": "WEEK 1 | FRI 19 JUL",
			"2024-07-20": "WEEK 1 | SAT 20 JUL",
			"2024-07-21": "WEEK 1 | SUN 21 JUL",
			"2024-07-26": "WEEK 2 | FRI 26 JUL",
			"2024-07-27": "WEEK 2 | SAT 27 JUL",
			"2024-07-28": "WEEK 2 | SUN 28 JUL",
		}
		responseEmbed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Tomorrowland Live",
				IconURL: s.State.User.AvatarURL(""),
			},
			Title:     dayNameMap[selectedDay.StringValue()],
			Color:     0xb917ff,
			Timestamp: time.Now().Format(time.RFC3339),
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://timetable.png",
			},
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Generating timetable...",
			},
		})

		r, err := bot.GetTimetable(selectedDay.StringValue())
		if err != nil {
			bot.BotError(err, "timetable")
			bot.ErrorInteractionResponse(s, i, config.Content{
				Message: "There was a problem generating the timetable. Please try again",
			}, true, false)
			return
		}

		c := ""
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &c,
			Embeds: &[]*discordgo.MessageEmbed{
				responseEmbed,
			},
			Files: []*discordgo.File{
				{Name: "timetable.png", ContentType: "image/png", Reader: r},
			},
		})
		if err != nil {
			bot.BotError(err, "timetable")
		}
	},
}
