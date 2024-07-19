package commands

import (
	"fmt"
	"sort"
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
			URL: config.OneWorldRadioImg,
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
			URL: config.OneWorldRadioImg,
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

	// Sort performances by start time
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].StartTime.Unix() < performances[j].StartTime.Unix()
	})

	var fields []*discordgo.MessageEmbedField = make([]*discordgo.MessageEmbedField, len(performances))
	for i, p := range performances {
		fields[i] = &discordgo.MessageEmbedField{
			Name:   p.ArtistName,
			Value:  fmt.Sprintf("%v to %v", p.StartTime.Format("15:04"), p.EndTime.Format("15:04")),
			Inline: false,
		}
	}

	// if the length is greater than 25
	// TODO: Implement pagination in the future
	if len(fields) > 25 {
		fields = fields[:25]
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

	responseEmbed.Title = fmt.Sprintf("%v\n%v", selectedStage.StringValue(), tmrlweb.DayNameMap[selectedDay.StringValue()])
	responseEmbed.Fields = fields
	responseEmbed.Description = "All times are in GMT+2 (Brussels Time)"
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
				Description: "See the artists playing on a specific stage and day",
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
