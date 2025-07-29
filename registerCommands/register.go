package registerCommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetCommands(minn *float64, maxx float64) []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "sala",
			Description: "Create a TEMPORARY voice channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Choose a voice channel name",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "🏢 Working", Value: "🏢 Working"},
						{Name: "🗣️ Only English", Value: "🗣️ Only English"},
						{Name: "😊 Chill Talk", Value: "😊 Chill Talk"},
						{Name: "💻 Study Chat", Value: "💻 Study Chat"},
						{Name: "🔴 LIVE", Value: "🔴 LIVE"},
						{Name: "🆘 HELP ME", Value: "🆘 HELP ME"},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "limit",
					Description: "Number of users to start a voice",
					Required:    true,
					MinValue:    minn,
					MaxValue:    40,
				},
			},
		},
		{
			Name:        "voice-limite",
			Description: "Change the user limit of your voice channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "limit",
					Description: "New user limit (0 for unlimited)",
					Required:    true,
					MinValue:    minn,
					MaxValue:    99,
				},
			},
		},
	}
}

func RegisterCommands(session *discordgo.Session, appID, guildID string, commands []*discordgo.ApplicationCommand) {
	for _, cmds := range commands {
		if _, err := session.ApplicationCommandCreate(appID, guildID, cmds); err != nil {
			fmt.Printf("Error creating command %s: %v\n", cmds.Name, err)
		}
	}
}
