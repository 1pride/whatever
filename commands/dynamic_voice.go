package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func StartVoiceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userLimit := 10
	channelName := "🔊 Temp VC" // default name
	categoryID := "1380359231336611841"

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "limit":
			val := int(option.IntValue())
			if val < 2 || val > 40 {
				_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "❌ Please choose a user limit between **2** and **40**.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}
			userLimit = val

		case "name":
			channelName = option.StringValue()
		}
	}

	// Deferred response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Error sending deferred response:", err)
		return
	}

	channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:      channelName,
		Type:      discordgo.ChannelTypeGuildVoice,
		UserLimit: userLimit,
		ParentID:  categoryID,
	})
	if err != nil {
		fmt.Println("Error creating voice channel:", err)
		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Failed to create voice channel.",
		})
		return
	}

	content := fmt.Sprintf("✅ Temp voice channel created: **%s** (limit: %d users)\n\n➡️ https://discord.com/channels/%s/%s", channel.Name, userLimit, i.GuildID, channel.ID)
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	fmt.Println(content)
	if err != nil {
		fmt.Println("Error editing response:", err)
	}

	monitorAndDeleteWhenEmpty(s, i.GuildID, channel.ID, 5*time.Second, 10*time.Second)
}

func monitorAndDeleteWhenEmpty(s *discordgo.Session, guildID, channelID string, checkInterval time.Duration, emptyDuration time.Duration) {
	go func() {
		var emptySince time.Time

		for {
			// Get current voice states in the guild
			guild, err := s.State.Guild(guildID)
			if err != nil {
				fmt.Println("Error getting guild state:", err)
				return
			}

			active := false
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == channelID {
					active = true
					break
				}
			}

			if active {
				// Reset empty timer if anyone is in channel
				emptySince = time.Time{}
			} else {
				if emptySince.IsZero() {
					emptySince = time.Now()
				} else if time.Since(emptySince) >= emptyDuration {
					// Channel has been empty for long enough — delete it
					_, err := s.ChannelDelete(channelID)
					if err != nil {
						fmt.Println("Error deleting empty voice channel:", err)
					} else {
						fmt.Println("✅ Temp voice channel deleted due to sustained inactivity.")
					}
					return
				}
			}

			time.Sleep(checkInterval)
		}
	}()
}
