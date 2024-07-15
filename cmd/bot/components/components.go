package components

import (
	"github.com/CDavidSV/Tomorrowland-Bot-Go/cmd/bot/config"
	"github.com/bwmarrin/discordgo"
)

var Components = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, bot *config.Bot){
	"liveSelect": SelectLive,
}
