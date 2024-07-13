package commands

import (
	"log"

	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Data     *discordgo.ApplicationCommand
	Callback func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot)
}

var (
	dmPermissionFalse = false
	dmPermissionTrue  = true
	SlashCommands     = map[string]Command{}
)

func setCommand(c Command) {
	SlashCommands[c.Data.Name] = c
}

func LoadCommands(s *discordgo.Session, bot *config.Bot, reloadCommands bool, deleteOldCommands bool, testGuildID string) {
	// Commands
	setCommand(PlayCommand)
	setCommand(StopCommand)
	setCommand(InviteCommand)
	setCommand(LineupCommand)

	if deleteOldCommands {
		log.Println("Deleting old commands...")
		// Delete old commands
		oldCommands, err := s.ApplicationCommands(s.State.User.ID, testGuildID)
		if err != nil {
			log.Panicf("Cannot get old commands: %v", err)
		}

		for _, command := range oldCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, testGuildID, command.ID)
			if err != nil {
				log.Panicf("Failed to delete old command: %v", err)
			}

			bot.Logger.Info("Successfully deleted old command:", "command", command.Name)
		}

		bot.Logger.Info("Successfully deleted all old application commands")
	}

	if reloadCommands {
		log.Println("Loading commands...")
		// Register commands
		for _, v := range SlashCommands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v.Data)
			if err != nil {
				log.Panicf("Cannot create '%v' command:", err)
			}

			bot.Logger.Info("Successfully created command:", "command", v.Data.Name)
		}
		bot.Logger.Info("Successfully created all application commands")
	}
}
