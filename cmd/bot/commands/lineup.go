package commands

import (
	"fmt"
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/tmrlweb"
	"github.com/bwmarrin/discordgo"
)

func getTimetable(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
	// Get the nested options for the selected dat for the command
	selectedDay := i.Interaction.ApplicationCommandData().Options[0].Options[0]

	responseEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Tomorrowland Live",
			IconURL: s.State.User.AvatarURL(""),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://d384fynlilbsl.cloudfront.net/one_world_radio_logo.png",
		},
		Title:     tmrlweb.DayNameMap[selectedDay.StringValue()],
		Color:     config.MainColorHex,
		Timestamp: time.Now().Format(time.RFC3339),
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://timetable.png",
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	r, err := tmrlweb.GetTimetableImage(selectedDay.StringValue())
	if err != nil {
		bot.BotError(err, "timetable")
		bot.ErrorInteractionResponse(s, i, config.Content{
			Message: "There was a problem generating the timetable. Please try again",
		}, true, false)
		return
	}

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label: "View Online",
				URL:   "https://belgium.tomorrowland.com/en/line-up/?page=timetable&day=" + selectedDay.StringValue(),
				Style: discordgo.LinkButton,
			},
		},
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
		Components: &[]discordgo.MessageComponent{
			actionRow,
		},
	})
	if err != nil {
		bot.BotError(err, "timetable")
	}
}

func getStage(s *discordgo.Session, i *discordgo.InteractionCreate, _ *config.Bot) {
	// Get the nested options for the subcommand
	selectedStage := i.Interaction.ApplicationCommandData().Options[0].Options[0]
	selectedDay := i.Interaction.ApplicationCommandData().Options[0].Options[1]

	performances := tmrlweb.GetPerformances(selectedDay.StringValue(), selectedStage.StringValue())

	responseEmbed := &discordgo.MessageEmbed{
		Title: "There are no performances for this stage on this day",
		Color: config.MainColorHex,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://d384fynlilbsl.cloudfront.net/one_world_radio_logo.png",
		},
	}

	if performances == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					responseEmbed,
				},
			},
		})
		return
	}

	fields := ""
	for _, p := range performances {
		fields += fmt.Sprintf("\n**%v**", p.ArtistName)
	}

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label: "View Online",
				URL:   "https://belgium.tomorrowland.com/en/line-up/?page=stages&day=" + selectedDay.StringValue(),
				Style: discordgo.LinkButton,
			},
		},
	}

	responseEmbed.Title = fmt.Sprintf("%v\n%v", selectedStage.StringValue(), selectedDay.StringValue())
	responseEmbed.Description = fields
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
}

var LineupCommand Command = Command{
	Data: &discordgo.ApplicationCommand{
		Name:         "lineup",
		Description:  "View the Tomorrowland lineup for a specific day",
		DMPermission: &dmPermissionTrue,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "timetable",
				Description: "Get the timetable for a specific day",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "day",
						Description: "Select a day",
						Choices:     tmrlweb.GetDayChoices(),
						Required:    true,
					},
				},
			},
			{
				Name:        "stage",
				Description: "See the artist playing on a specific stage and day",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "stage-name",
						Description: "Select a stage",
						Choices:     tmrlweb.GetStageChoices(),
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "day",
						Description: "Select a day",
						Choices:     tmrlweb.GetDayChoices(),
						Required:    true,
					},
				},
			},
		},
	},
	Callback: func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
		switch i.ApplicationCommandData().Options[0].Name {
		case "timetable":
			getTimetable(s, i, bot)
		case "stage":
			getStage(s, i, bot)
		}
	},
}
