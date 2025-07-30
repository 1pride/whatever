package main

import (
	cmd "bot/commands"
	regC "bot/registerCommands"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	_        = godotenv.Load()
	BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	AppID    = os.Getenv("DISCORD_APP_ID")
	GuildID  = os.Getenv("DISCORD_GUILD_ID")
	minn     = float64(2)
	maxx     = float64(100)
)

func main() {
	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	// Slash command handler
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		switch i.ApplicationCommandData().Name {
		case "voice-limite":
			cmd.VoiceLimitHandler(s, i)
		case "sala":
			cmd.StartVoiceHandler(s, i)
		}
	})

	// Handle select menu interaction
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "role_assign_select" {
			// Get the role ID from the original message
			originalMessage, err := s.ChannelMessage(i.ChannelID, i.Message.ID)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "❌ Error getting role information",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}

			roleMention := originalMessage.Embeds[0].Fields[0].Value
			roleID := strings.Trim(roleMention, "<@&>")

			// Assign a role to selected users
			for _, userID := range i.MessageComponentData().Values {
				err := s.GuildMemberRoleAdd(i.GuildID, userID, roleID)
				if err != nil {
					fmt.Printf("Error assigning role to user %s: %v\n", userID, err)
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "✅ Roles assigned successfully",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	})

	// Open connection
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	} else {
		fmt.Println("BOT ON!")
	}

	// Register slash commands
	commands := regC.GetCommands(&minn, maxx)
	regC.RegisterCommands(session, AppID, GuildID, commands)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	err = session.Close()
	if err != nil {
		fmt.Println("Error closing Discord session,", err)
		return
	}
}
