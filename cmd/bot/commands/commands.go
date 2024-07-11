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
	DMPermissionFalse = false
	DMPermissionTrue  = true
	commands          = map[string]Command{}
)

func setCommand(c Command) {
	commands[c.Data.Name] = c
}

func LoadCommands(s *discordgo.Session, bot *config.Bot, reloadCommands bool, testGuildID string) {
	// Commands
	setCommand(PlayCommand)
	setCommand(StopCommand)
	setCommand(InviteCommand)

	defer s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if c, ok := commands[i.ApplicationCommandData().Name]; ok {
			// Log command in console
			bot.Logger.Info("Command executed", "command", i.ApplicationCommandData().Name, "user", i.Interaction.Member.User.ID, "guild", i.Interaction.GuildID)

			c.Callback(s, i, bot)
		}
	})

	if !reloadCommands {
		return
	}

	// Register commands
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v.Data)
		if err != nil {
			log.Panicf("Cannot create '%v' command:", err)
		}

		log.Println("Successfully created command:", v.Data.Name)
	}

	log.Println("Successfully created all application commands")
}
