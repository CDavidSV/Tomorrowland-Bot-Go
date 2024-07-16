package commands

import (
	"fmt"
	"time"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/CDavidSV/Tomorrowland-Bot-Go/internal/tmrlweb"
	"github.com/bwmarrin/discordgo"
)

var NowPlayingCommand Command = Command{
	Data: &discordgo.ApplicationCommand{
		Name:         "nowplaying",
		Description:  "Shows the artist currently playing in a stage",
		DMPermission: &dmPermissionTrue,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "stage",
				Description: "Select a stage",
				Choices:     tmrlweb.GetStageChoices(),
				Required:    false,
			},
		},
	},
	Callback: func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot) {
		loc, _ := time.LoadLocation("Europe/Brussels")
		nowBelgium := time.Now().In(loc)

		responseEmbed := &discordgo.MessageEmbed{
			Title: "Artist(s) currently playing",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Tomorrowland Live",
				IconURL: s.State.User.AvatarURL(""),
			},
			Color: config.MainColorHex,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: config.OneWorldRadioImg,
			},
		}

		fields := []*discordgo.MessageEmbedField{}

		if len(i.ApplicationCommandData().Options) > 0 {
			selectedStage := i.ApplicationCommandData().Options[0]

			responseEmbed.Description = selectedStage.StringValue()
			performances := tmrlweb.GetPerformances(nowBelgium.Format("2006-01-02"), selectedStage.StringValue())

			for _, performance := range performances {
				// Current performance is ongoing
				if tmrlweb.CurrentlyPLaying(performance) {
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:   performance.ArtistName,
						Value:  fmt.Sprintf("%v - %v", performance.StartTime.Format("15:04"), performance.EndTime.Format("15:04")),
						Inline: false,
					})
				}
			}
		} else {
			stages := tmrlweb.GetStages(nowBelgium.Format("2006-01-02"))
			// Loop over all stages and then check each performance to see i
			for _, stage := range stages {
				for _, performance := range stage {
					// Current performance is ongoing
					if tmrlweb.CurrentlyPLaying(performance) {
						fields = append(fields, &discordgo.MessageEmbedField{
							Name:   fmt.Sprintf("%v | %v", performance.ArtistName, performance.Stage.Name),
							Value:  fmt.Sprintf("%v to %v", performance.StartTime.Format("15:04"), performance.EndTime.Format("15:04")),
							Inline: false,
						})
					}
				}
			}
		}

		if len(fields) > 25 {
			fields = fields[:25]
		}

		if len(fields) == 0 {
			responseEmbed.Description = "There are no artists playing at the moment. Check the line-up"
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

		responseEmbed.Fields = fields
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
