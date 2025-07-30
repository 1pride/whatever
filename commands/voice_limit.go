package commands

import (
	"bot/config"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func VoiceLimitHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Validate member exists
	if i.Member == nil || i.Member.User == nil {
		respondWithError(s, i, "Cannot identify user")
		return
	}
	user := i.Member.User

	// Check if the user is in a voice channel
	guild, err := s.State.Guild(i.GuildID)
	if err != nil {
		respondWithError(s, i, "Error accessing guild information")
		return
	}

	// Find a user's voice state
	var userVoiceState *discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == user.ID {
			userVoiceState = vs
			break
		}
	}
	if userVoiceState == nil {
		respondWithError(s, i, "You must be in a voice channel to use this command")
		return
	}

	// Get the voice channel
	channel, err := s.Channel(userVoiceState.ChannelID)
	if err != nil {
		respondWithError(s, i, "Error accessing channel information")
		return
	}

	// Simplified permission checking for temporary voice channels
	// For channels created by the sala command, allow any user in the channel to modify the limit

	// Check if this is a temporary channel (in the same category as sala command)
	isTempChannel := channel.ParentID == config.CategoryID

	// Debug logging
	log.Printf("Channel ParentID: %s, Config CategoryID: %s, IsTempChannel: %v", channel.ParentID, config.CategoryID, isTempChannel)

	// For temporary channels, only allow the creator to modify the limit
	if isTempChannel {
		if !IsChannelOwner(channel.ID, user.ID) {
			respondWithError(s, i, "Only the creator of this temporary voice channel can change its limit.")
			return
		}
	} else {
		// If it's not a temporary channel, check for proper permissions
		// Try to check permissions, but don't fail if we can't
		perms, err := s.State.UserChannelPermissions(i.Member.User.ID, channel.ID)
		if err != nil {
			// Log the error for debugging
			log.Printf("Error checking permissions for user %s in channel %s: %v", user.ID, channel.ID, err)

			// For non-temporary channels, check if user is the channel owner
			if channel.OwnerID == user.ID {
				// User is the channel owner, allow it
			} else {
				// If we can't check permissions and user is not owner, deny access
				respondWithError(s, i, "Unable to verify permissions. Please ensure you have manage channel permissions.")
				return
			}
		} else {
			hasManagePerms := perms&discordgo.PermissionManageChannels != 0
			isOwner := channel.OwnerID == user.ID

			if !hasManagePerms && !isOwner {
				respondWithError(s, i, "You must be the owner or have manage channel permissions to change its limit")
				return
			}
		}
	}

	// Get command options directly since it's required
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		respondWithError(s, i, "Missing limit parameter")
		return
	}

	limit := int(options[0].IntValue())
	if limit < 0 || limit > 99 {
		respondWithError(s, i, "Limit must be between 0 and 99")
		return
	}

	// Update the channel with a new user limit
	_, err = s.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		UserLimit: limit,
	})
	if err != nil {
		respondWithError(s, i, fmt.Sprintf("Failed to update channel limit: %v", err))
		return
	}

	// Respond with a success message
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Voice channel limit has been set to %d", limit),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
